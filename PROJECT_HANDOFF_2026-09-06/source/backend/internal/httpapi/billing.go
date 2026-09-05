package httpapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type renewalActionRequest struct {
	BillingCycle *string  `json:"billingCycle"`
	Amount       *float64 `json:"amount"`
	PaidAt       *string  `json:"paidAt"`
	Notes        *string  `json:"notes"`
}

func (s *Server) listUserRenewals(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	if err := s.validateUserExists(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user renewals")
		return
	}
	items, err := s.readUserRenewalCandidates(id)
	if err != nil {
		s.logger.Error("read user renewal candidates", "user_id", id, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user renewals")
		return
	}
	writeSuccess(w, items)
}

func (s *Server) confirmUserRenewal(w http.ResponseWriter, r *http.Request) {
	s.processUserRenewal(w, r, true)
}

func (s *Server) rejectUserRenewal(w http.ResponseWriter, r *http.Request) {
	s.processUserRenewal(w, r, false)
}

func (s *Server) processUserRenewal(w http.ResponseWriter, r *http.Request, confirm bool) {
	userID := strings.TrimSpace(r.PathValue("id"))
	candidateID := strings.TrimSpace(r.PathValue("candidateId"))
	if userID == "" || candidateID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id and candidate id are required")
		return
	}
	var payload renewalActionRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "invalid renewal payload")
			return
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin renewal update")
		return
	}
	defer tx.Rollback()
	var oldExpiry, newExpiry, suggestedCycle, currency, status string
	var suggestedAmount float64
	err = tx.QueryRow(`SELECT old_expiry_at, new_expiry_at, suggested_cycle, suggested_amount, currency, status
FROM user_renewal_candidates WHERE id = ? AND user_id = ?`, candidateID, userID).
		Scan(&oldExpiry, &newExpiry, &suggestedCycle, &suggestedAmount, &currency, &status)
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "renewal candidate not found")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read renewal candidate")
		return
	}
	if status != "pending" {
		writeFailure(w, http.StatusConflict, validationCode, "renewal candidate has already been processed")
		return
	}
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	resolvedBy := ""
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok {
		resolvedBy = current.UserID
	}
	if !confirm {
		notes := ""
		if payload.Notes != nil {
			notes = strings.TrimSpace(*payload.Notes)
		}
		if _, err := tx.Exec(`UPDATE user_renewal_candidates SET status = 'rejected', processed_at = ?, notes = ? WHERE id = ?`, nowText, nullableDBString(notes), candidateID); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not reject renewal candidate")
			return
		}
		if err := resolveRenewalEventTx(tx, candidateID, "renewal_non_billable", map[string]any{"candidateId": candidateID, "userId": userID, "notes": notes}, now, resolvedBy); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close renewal event")
			return
		}
		if err := tx.Commit(); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit renewal decision")
			return
		}
		s.writeAuditLog(r, "user.renewal.reject", "user", userID, nil, map[string]any{"candidateId": candidateID})
		result, err := s.readUserDetail(userID)
		if err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "renewal rejected but could not read user")
			return
		}
		writeSuccess(w, result)
		return
	}

	cycle := suggestedCycle
	if payload.BillingCycle != nil {
		cycle = strings.TrimSpace(*payload.BillingCycle)
	}
	if cycle != "monthly" && cycle != "annual" {
		writeFailure(w, http.StatusBadRequest, validationCode, "billingCycle must be monthly or annual")
		return
	}
	amount := suggestedAmount
	if payload.Amount != nil {
		amount = *payload.Amount
	}
	if amount < 0 || amount > 100000000 || amount != amount {
		writeFailure(w, http.StatusBadRequest, validationCode, "amount must be a non-negative number")
		return
	}
	serviceFrom, err := time.Parse(time.RFC3339Nano, oldExpiry)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "candidate service start is invalid")
		return
	}
	serviceTo, err := time.Parse(time.RFC3339Nano, newExpiry)
	if err != nil || !serviceTo.After(serviceFrom) {
		writeFailure(w, http.StatusBadRequest, validationCode, "candidate service interval is invalid")
		return
	}
	notes := ""
	if payload.Notes != nil {
		notes = strings.TrimSpace(*payload.Notes)
	}
	paidAt := nowText
	if payload.PaidAt != nil && strings.TrimSpace(*payload.PaidAt) != "" {
		parsed, parseErr := time.Parse(time.RFC3339Nano, strings.TrimSpace(*payload.PaidAt))
		if parseErr != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "paidAt must be an RFC3339 timestamp")
			return
		}
		paidAt = parsed.UTC().Format(time.RFC3339Nano)
	}
	if _, err := tx.Exec(`UPDATE user_renewal_candidates SET status = 'confirmed', processed_at = ?, notes = ? WHERE id = ? AND status = 'pending'`, nowText, nullableDBString(notes), candidateID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not confirm renewal candidate")
		return
	}
	if err := resolveRenewalEventTx(tx, candidateID, "renewal_confirmed", map[string]any{
		"candidateId": candidateID, "userId": userID, "billingCycle": cycle, "amount": amount, "currency": currency,
		"serviceFrom": serviceFrom.UTC().Format(time.RFC3339Nano), "serviceTo": serviceTo.UTC().Format(time.RFC3339Nano),
	}, now, resolvedBy); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not close renewal event")
		return
	}
	if _, err := tx.Exec(`INSERT INTO user_billing_records
(id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, renewal_candidate_id, notes, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'confirmed', 'agent_detected', ?, ?, ?)`, newID(), userID, cycle, amount, currency, serviceFrom.UTC().Format(time.RFC3339Nano), serviceTo.UTC().Format(time.RFC3339Nano), paidAt, candidateID, nullableDBString(notes), nowText); err != nil {
		returnError := fmt.Errorf("create billing record: %w", err)
		s.logger.Error("create billing record", "user_id", userID, "error", returnError)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create billing record")
		return
	}
	monthlyEquivalent := amount
	if cycle == "annual" {
		monthlyEquivalent = amount / 12
	}
	if _, err := tx.Exec(`UPDATE users SET billing_cycle = ?, billing_amount = ?, monthly_fee = ?, updated_at = ? WHERE id = ?`, cycle, amount, monthlyEquivalent, nowText, userID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not update user billing terms")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit renewal confirmation")
		return
	}
	s.writeAuditLog(r, "user.renewal.confirm", "user", userID, nil, map[string]any{"candidateId": candidateID, "billingCycle": cycle, "amount": amount})
	result, err := s.readUserDetail(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "renewal confirmed but could not read user")
		return
	}
	writeSuccess(w, result)
}

func (s *Server) readUserRenewalCandidates(userID string) ([]map[string]any, error) {
	rows, err := s.db.Query(`SELECT id, COALESCE(inbound_id, ''), old_expiry_at, new_expiry_at, detected_at, suggested_cycle, suggested_amount, currency, status, COALESCE(processed_at, ''), COALESCE(notes, '')
FROM user_renewal_candidates WHERE user_id = ? ORDER BY detected_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, inboundID, oldExpiry, newExpiry, detectedAt, cycle, currency, status, processedAt, notes string
		var amount float64
		if err := rows.Scan(&id, &inboundID, &oldExpiry, &newExpiry, &detectedAt, &cycle, &amount, &currency, &status, &processedAt, &notes); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id": id, "inboundId": nullableString(inboundID), "oldExpiryAt": oldExpiry, "newExpiryAt": newExpiry,
			"detectedAt": detectedAt, "billingCycle": cycle, "suggestedAmount": amount, "currency": currency,
			"status": status, "processedAt": nullableString(processedAt), "notes": nullableString(notes),
		})
	}
	return items, rows.Err()
}

func (s *Server) readUserBillingRecords(userID string) ([]map[string]any, error) {
	rows, err := s.db.Query(`SELECT id, billing_cycle, amount, currency, service_from, service_to, COALESCE(paid_at, ''), status, source, COALESCE(notes, ''), created_at
FROM user_billing_records WHERE user_id = ? ORDER BY service_from DESC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, cycle, currency, serviceFrom, serviceTo, paidAt, status, source, notes, createdAt string
		var amount float64
		if err := rows.Scan(&id, &cycle, &amount, &currency, &serviceFrom, &serviceTo, &paidAt, &status, &source, &notes, &createdAt); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id": id, "billingCycle": cycle, "amount": amount, "currency": currency,
			"serviceFrom": serviceFrom, "serviceTo": serviceTo, "paidAt": nullableString(paidAt),
			"status": status, "source": source, "notes": nullableString(notes), "createdAt": createdAt,
		})
	}
	return items, rows.Err()
}

// accruedUserIncome allocates confirmed payments over their service interval.
// Users without a confirmed history use the legacy monthly_fee fallback so
// existing installations keep their current financial totals after migration.
func (s *Server) accruedUserIncome(startText, endText string) (float64, error) {
	start, err := time.Parse(time.RFC3339Nano, startText)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse(time.RFC3339Nano, endText)
	if err != nil {
		return 0, err
	}
	rows, err := s.db.Query(`SELECT id, monthly_fee
FROM users
WHERE currency = 'CNY' AND status <> 'disabled'
  AND datetime(created_at) < datetime(?)
  AND (expiry_time IS NULL OR datetime(expiry_time) >= datetime(?))`, end.UTC().Format(time.RFC3339Nano), start.UTC().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	users := make([]struct {
		id            string
		legacyMonthly float64
	}, 0)
	for rows.Next() {
		var item struct {
			id            string
			legacyMonthly float64
		}
		if err := rows.Scan(&item.id, &item.legacyMonthly); err != nil {
			return 0, err
		}
		users = append(users, item)
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	total := 0.0
	for _, user := range users {
		userID := user.id
		legacyMonthly := user.legacyMonthly
		recordRows, err := s.db.Query(`SELECT amount, service_from, service_to
FROM user_billing_records WHERE user_id = ? AND currency = 'CNY' AND status = 'confirmed'`, userID)
		if err != nil {
			return 0, err
		}
		hasRecords := false
		userIncome := 0.0
		for recordRows.Next() {
			hasRecords = true
			var amount float64
			var fromText, toText string
			if err := recordRows.Scan(&amount, &fromText, &toText); err != nil {
				_ = recordRows.Close()
				return 0, err
			}
			from, fromErr := time.Parse(time.RFC3339Nano, fromText)
			to, toErr := time.Parse(time.RFC3339Nano, toText)
			if fromErr != nil || toErr != nil || !to.After(from) {
				continue
			}
			overlapFrom := from
			if overlapFrom.Before(start) {
				overlapFrom = start
			}
			overlapTo := to
			if overlapTo.After(end) {
				overlapTo = end
			}
			if overlapTo.After(overlapFrom) {
				userIncome += amount * overlapTo.Sub(overlapFrom).Seconds() / to.Sub(from).Seconds()
			}
		}
		if err := recordRows.Close(); err != nil {
			return 0, err
		}
		if hasRecords {
			total += userIncome
		} else {
			total += legacyMonthly
		}
	}
	return total, nil
}
