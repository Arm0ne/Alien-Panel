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

type billingRecordRequest struct {
	BillingCycle string  `json:"billingCycle"`
	Amount       float64 `json:"amount"`
	ServiceFrom  string  `json:"serviceFrom"`
	ServiceTo    string  `json:"serviceTo"`
	PaidAt       string  `json:"paidAt"`
	OrderType    string  `json:"orderType"`
	Notes        string  `json:"notes"`
}

type billingImportItem struct {
	UserID          string  `json:"userId"`
	NodeID          string  `json:"nodeId"`
	RemoteInboundID string  `json:"remoteInboundId"`
	BillingCycle    string  `json:"billingCycle"`
	Amount          float64 `json:"amount"`
	ServiceFrom     string  `json:"serviceFrom"`
	ServiceTo       string  `json:"serviceTo"`
	PaidAt          string  `json:"paidAt"`
	OrderType       string  `json:"orderType"`
	Verified        *bool   `json:"verified"`
	Notes           string  `json:"notes"`
}

type billingImportRequest struct {
	Records []billingImportItem `json:"records"`
	DryRun  bool                `json:"dryRun"`
}

type billingImportResult struct {
	DryRun     bool             `json:"dryRun"`
	Imported   int              `json:"imported"`
	Unverified int              `json:"unverified"`
	Records    []map[string]any `json:"records"`
	Errors     []map[string]any `json:"errors"`
}

type billingRecordActionRequest struct {
	PaidAt string `json:"paidAt"`
	Notes  string `json:"notes"`
}

// createBillingRecord records a manually confirmed payment. The service
// interval is explicit because a payment date alone cannot establish when a
// yearly or monthly service period began.
func (s *Server) createBillingRecord(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.PathValue("id"))
	if userID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id is required")
		return
	}
	var payload billingRecordRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid billing record payload")
		return
	}
	if err := s.validateUserExists(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "user not found")
			return
		}
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read user")
		return
	}
	cycle := strings.TrimSpace(payload.BillingCycle)
	if cycle != "monthly" && cycle != "annual" {
		writeFailure(w, http.StatusBadRequest, validationCode, "billingCycle must be monthly or annual")
		return
	}
	if payload.Amount < 0 || payload.Amount > 100000000 || payload.Amount != payload.Amount {
		writeFailure(w, http.StatusBadRequest, validationCode, "amount must be a non-negative number")
		return
	}
	orderType := strings.TrimSpace(payload.OrderType)
	if orderType == "" {
		orderType = "initial"
	}
	if orderType != "initial" && orderType != "renewal" && orderType != "recovery" {
		writeFailure(w, http.StatusBadRequest, validationCode, "orderType must be initial, renewal, or recovery")
		return
	}
	from, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(payload.ServiceFrom))
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "serviceFrom must be an RFC3339 timestamp")
		return
	}
	to, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(payload.ServiceTo))
	if err != nil || !to.After(from) {
		writeFailure(w, http.StatusBadRequest, validationCode, "serviceTo must be after serviceFrom")
		return
	}
	var overlap int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM user_billing_records
WHERE user_id = ? AND status <> 'cancelled' AND datetime(service_from) < datetime(?) AND datetime(service_to) > datetime(?)`, userID,
		to.UTC().Format(time.RFC3339Nano), from.UTC().Format(time.RFC3339Nano)).Scan(&overlap); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not check billing record interval")
		return
	}
	if overlap > 0 {
		writeFailure(w, http.StatusConflict, validationCode, "service interval overlaps an existing billing record")
		return
	}
	paidAt := time.Now().UTC()
	if strings.TrimSpace(payload.PaidAt) != "" {
		paidAt, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(payload.PaidAt))
		if err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "paidAt must be an RFC3339 timestamp")
			return
		}
		paidAt = paidAt.UTC()
	}
	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	orderNo := fmt.Sprintf("AP-%s-%s", time.Now().UTC().Format("20060102"), strings.ToUpper(newID()[:8]))
	_, err = s.db.Exec(`INSERT INTO user_billing_records
(id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, order_type, order_no, notes, created_at)
VALUES (?, ?, ?, ?, 'CNY', ?, ?, ?, 'confirmed', 'manual', ?, ?, ?, ?)`, newID(), userID, cycle, payload.Amount,
		from.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano), paidAt.Format(time.RFC3339Nano), orderType, orderNo, nullableDBString(strings.TrimSpace(payload.Notes)), nowText)
	if err != nil {
		s.logger.Error("create billing record", "user_id", userID, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create billing record")
		return
	}
	monthlyEquivalent := payload.Amount
	if cycle == "annual" {
		monthlyEquivalent = payload.Amount / 12
	}
	if _, err := s.db.Exec(`UPDATE users SET billing_type = 'paid', free_reason = NULL, billing_cycle = ?, billing_amount = ?, monthly_fee = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`, cycle, payload.Amount, monthlyEquivalent, nowText, userID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "billing record created but could not update user billing terms")
		return
	}
	s.writeAuditLog(r, "user.billing.create", "user", userID, nil, map[string]any{"orderNo": orderNo, "orderType": orderType, "amount": payload.Amount})
	result, err := s.readUserDetail(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "billing record created but could not read user")
		return
	}
	writeSuccess(w, result)
}

// importBillingRecords creates a reviewed batch of historical billing rows.
// Rows are matched by userId or by the stable nodeId + remoteInboundId pair.
// A dry run validates the complete batch without writing anything.
func (s *Server) importBillingRecords(w http.ResponseWriter, r *http.Request) {
	var payload billingImportRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid billing import payload")
		return
	}
	if len(payload.Records) == 0 || len(payload.Records) > 2000 {
		writeFailure(w, http.StatusBadRequest, validationCode, "records must contain between 1 and 2000 items")
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin billing import")
		return
	}
	defer tx.Rollback()

	result := billingImportResult{DryRun: payload.DryRun, Records: make([]map[string]any, 0, len(payload.Records)), Errors: make([]map[string]any, 0)}
	type interval struct {
		userID string
		from   time.Time
		to     time.Time
	}
	intervals := make([]interval, 0, len(payload.Records))
	for index, item := range payload.Records {
		userID, err := resolveBillingImportUser(tx, item)
		if err != nil {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": err.Error()})
			continue
		}
		cycle := strings.TrimSpace(item.BillingCycle)
		if cycle != "monthly" && cycle != "annual" {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "billingCycle must be monthly or annual"})
			continue
		}
		if item.Amount < 0 || item.Amount > 100000000 || item.Amount != item.Amount {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "amount must be a non-negative number"})
			continue
		}
		orderType := strings.TrimSpace(item.OrderType)
		if orderType == "" {
			orderType = "renewal"
		}
		if orderType != "initial" && orderType != "renewal" && orderType != "recovery" {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "orderType must be initial, renewal, or recovery"})
			continue
		}
		from, err := parseBillingDate(item.ServiceFrom)
		if err != nil {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "serviceFrom must be YYYY-MM-DD or RFC3339"})
			continue
		}
		to, err := parseBillingDate(item.ServiceTo)
		if err != nil || !to.After(from) {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "serviceTo must be after serviceFrom"})
			continue
		}
		verified := item.Verified == nil || *item.Verified
		paidAt := time.Time{}
		if strings.TrimSpace(item.PaidAt) != "" {
			paidAt, err = parseBillingDateTime(item.PaidAt)
			if err != nil {
				result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "paidAt must be YYYY-MM-DD or RFC3339"})
				continue
			}
		}
		if verified && paidAt.IsZero() {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "verified historical records require paidAt"})
			continue
		}
		rowOverlaps := false
		for _, existing := range intervals {
			if existing.userID == userID && from.Before(existing.to) && existing.from.Before(to) {
				result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "service interval overlaps another imported row"})
				rowOverlaps = true
				break
			}
		}
		if rowOverlaps {
			continue
		}
		var duplicate int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM user_billing_records
WHERE user_id = ? AND service_from = ? AND service_to = ? AND billing_cycle = ? AND amount = ? AND status <> 'cancelled'`, userID,
			from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano), cycle, item.Amount).Scan(&duplicate); err != nil {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "could not check duplicate billing record"})
			continue
		}
		if duplicate > 0 {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "an equivalent billing record already exists"})
			continue
		}
		var overlap int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM user_billing_records
WHERE user_id = ? AND status <> 'cancelled' AND datetime(service_from) < datetime(?) AND datetime(service_to) > datetime(?)`, userID,
			to.Format(time.RFC3339Nano), from.Format(time.RFC3339Nano)).Scan(&overlap); err != nil {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "could not check existing service intervals"})
			continue
		}
		if overlap > 0 {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "service interval overlaps an existing billing record"})
			continue
		}
		intervals = append(intervals, interval{userID: userID, from: from, to: to})
		originStatus := "verified"
		rowStatus := "confirmed"
		if !verified {
			originStatus = "unverified"
			rowStatus = "pending"
			result.Unverified++
		}
		result.Records = append(result.Records, map[string]any{
			"row": index + 1, "userId": userID, "billingCycle": cycle, "amount": item.Amount,
			"serviceFrom": from.Format(time.RFC3339Nano), "serviceTo": to.Format(time.RFC3339Nano),
			"paidAt": nullableString(formatOptionalBillingTime(paidAt)), "orderType": orderType,
			"status": rowStatus, "verificationStatus": originStatus,
		})
		if payload.DryRun {
			continue
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		orderNo := fmt.Sprintf("AP-%s-%s", time.Now().UTC().Format("20060102"), strings.ToUpper(newID()[:8]))
		if _, err := tx.Exec(`INSERT INTO user_billing_records
(id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, order_type, order_no, notes, created_at, origin, verification_status)
VALUES (?, ?, ?, ?, 'CNY', ?, ?, ?, ?, 'manual', ?, ?, ?, ?, 'historical_import', ?)`, newID(), userID, cycle, item.Amount,
			from.Format(time.RFC3339Nano), to.Format(time.RFC3339Nano), nullableDBString(formatOptionalBillingTime(paidAt)), rowStatus, orderType, orderNo,
			nullableDBString(strings.TrimSpace(item.Notes)), now, originStatus); err != nil {
			result.Errors = append(result.Errors, map[string]any{"row": index + 1, "message": "could not create historical billing record"})
			continue
		}
		result.Imported++
	}
	if len(result.Errors) > 0 {
		if payload.DryRun {
			_ = tx.Rollback()
			writeSuccess(w, result)
			return
		}
		writeFailure(w, http.StatusBadRequest, validationCode, fmt.Sprintf("billing import has %d invalid rows", len(result.Errors)))
		return
	}
	if payload.DryRun {
		_ = tx.Rollback()
		writeSuccess(w, result)
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit billing import")
		return
	}
	s.writeAuditLog(r, "user.billing.import", "billing_batch", newID(), nil, map[string]any{"imported": result.Imported, "unverified": result.Unverified})
	writeSuccess(w, result)
}

func resolveBillingImportUser(tx *sql.Tx, item billingImportItem) (string, error) {
	if value := strings.TrimSpace(item.UserID); value != "" {
		var exists int
		if err := tx.QueryRow(`SELECT 1 FROM users WHERE id = ? AND deleted_at IS NULL`, value).Scan(&exists); errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		} else if err != nil {
			return "", fmt.Errorf("could not read user")
		}
		return value, nil
	}
	nodeID := strings.TrimSpace(item.NodeID)
	remoteID := strings.TrimSpace(item.RemoteInboundID)
	if nodeID == "" || remoteID == "" {
		return "", fmt.Errorf("userId or nodeId + remoteInboundId is required")
	}
	var userID string
	err := tx.QueryRow(`SELECT i.user_id FROM inbounds i JOIN users u ON u.id = i.user_id AND u.deleted_at IS NULL
WHERE i.node_id = ? AND i.remote_inbound_id = ? AND i.user_id IS NOT NULL`, nodeID, remoteID).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("no business user matches nodeId + remoteInboundId")
	}
	if err != nil {
		return "", fmt.Errorf("could not resolve business user")
	}
	return userID, nil
}

func parseBillingDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed.UTC(), nil
	}
	return time.Parse(time.RFC3339Nano, value)
}

func parseBillingDateTime(value string) (time.Time, error) {
	parsed, err := parseBillingDate(value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func formatOptionalBillingTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func (s *Server) verifyBillingRecord(w http.ResponseWriter, r *http.Request) {
	s.processBillingRecordAction(w, r, true)
}

func (s *Server) cancelBillingRecord(w http.ResponseWriter, r *http.Request) {
	s.processBillingRecordAction(w, r, false)
}

func (s *Server) processBillingRecordAction(w http.ResponseWriter, r *http.Request, verify bool) {
	userID := strings.TrimSpace(r.PathValue("id"))
	recordID := strings.TrimSpace(r.PathValue("recordId"))
	if userID == "" || recordID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "user id and record id are required")
		return
	}
	var payload billingRecordActionRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := decodeJSON(r, &payload); err != nil {
			writeFailure(w, http.StatusBadRequest, validationCode, "invalid billing record action payload")
			return
		}
	}
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin billing record update")
		return
	}
	defer tx.Rollback()
	var status, origin, source, paidAt string
	if err := tx.QueryRow(`SELECT status, COALESCE(origin, 'live'), source, COALESCE(paid_at, '') FROM user_billing_records WHERE id = ? AND user_id = ?`, recordID, userID).Scan(&status, &origin, &source, &paidAt); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "billing record not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read billing record")
		return
	}
	if verify && origin != "historical_import" {
		writeFailure(w, http.StatusConflict, validationCode, "only historical records can be changed here")
		return
	}
	if !verify && source != "manual" {
		writeFailure(w, http.StatusConflict, validationCode, "only manually entered records can be cancelled here")
		return
	}
	if verify {
		if status != "pending" {
			writeFailure(w, http.StatusConflict, validationCode, "billing record is not waiting for verification")
			return
		}
		if strings.TrimSpace(payload.PaidAt) != "" {
			parsed, parseErr := parseBillingDateTime(payload.PaidAt)
			if parseErr != nil {
				writeFailure(w, http.StatusBadRequest, validationCode, "paidAt must be YYYY-MM-DD or RFC3339")
				return
			}
			paidAt = parsed.Format(time.RFC3339Nano)
		}
		if paidAt == "" {
			writeFailure(w, http.StatusBadRequest, validationCode, "verified billing records require paidAt")
			return
		}
		if _, err := tx.Exec(`UPDATE user_billing_records SET status = 'confirmed', verification_status = 'verified', paid_at = ?, notes = COALESCE(?, notes) WHERE id = ? AND user_id = ?`, paidAt, nullableDBString(strings.TrimSpace(payload.Notes)), recordID, userID); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not verify billing record")
			return
		}
	} else {
		if status == "cancelled" {
			writeFailure(w, http.StatusConflict, validationCode, "billing record is already cancelled")
			return
		}
		if _, err := tx.Exec(`UPDATE user_billing_records SET status = 'cancelled', notes = COALESCE(?, notes) WHERE id = ? AND user_id = ?`, nullableDBString(strings.TrimSpace(payload.Notes)), recordID, userID); err != nil {
			writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not cancel billing record")
			return
		}
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit billing record update")
		return
	}
	action := "user.billing.verify"
	if !verify {
		action = "user.billing.cancel"
	}
	s.writeAuditLog(r, action, "billing_record", recordID, map[string]any{"status": status}, map[string]any{"status": map[bool]string{true: "confirmed", false: "cancelled"}[verify]})
	result, err := s.readUserDetail(userID)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "billing record updated but could not read user")
		return
	}
	writeSuccess(w, result)
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
	err = tx.QueryRow(`SELECT c.old_expiry_at, c.new_expiry_at, c.suggested_cycle, c.suggested_amount, c.currency, c.status
FROM user_renewal_candidates c
JOIN users u ON u.id = c.user_id AND u.deleted_at IS NULL
WHERE c.id = ? AND c.user_id = ?`, candidateID, userID).
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
	orderType := "renewal"
	if !serviceFrom.Before(now) {
		orderType = "renewal"
	}
	if oldExpiryTime, parseErr := time.Parse(time.RFC3339Nano, oldExpiry); parseErr == nil && oldExpiryTime.Before(now) {
		orderType = "recovery"
	}
	orderNo := fmt.Sprintf("AP-%s-%s", now.Format("20060102"), strings.ToUpper(newID()[:8]))
	if _, err := tx.Exec(`INSERT INTO user_billing_records
(id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, renewal_candidate_id, order_type, order_no, notes, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'confirmed', 'agent_detected', ?, ?, ?, ?, ?)`, newID(), userID, cycle, amount, currency, serviceFrom.UTC().Format(time.RFC3339Nano), serviceTo.UTC().Format(time.RFC3339Nano), paidAt, candidateID, orderType, orderNo, nullableDBString(notes), nowText); err != nil {
		returnError := fmt.Errorf("create billing record: %w", err)
		s.logger.Error("create billing record", "user_id", userID, "error", returnError)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not create billing record")
		return
	}
	monthlyEquivalent := amount
	if cycle == "annual" {
		monthlyEquivalent = amount / 12
	}
	if _, err := tx.Exec(`UPDATE users SET billing_type = 'paid', free_reason = NULL, billing_cycle = ?, billing_amount = ?, monthly_fee = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`, cycle, amount, monthlyEquivalent, nowText, userID); err != nil {
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
	rows, err := s.db.Query(`SELECT id, billing_cycle, amount, currency, service_from, service_to, COALESCE(paid_at, ''), status, source, COALESCE(order_type, 'initial'), COALESCE(order_no, ''), COALESCE(notes, ''), created_at, COALESCE(origin, 'live'), COALESCE(verification_status, 'verified')
FROM user_billing_records WHERE user_id = ? ORDER BY service_from DESC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]map[string]any, 0)
	for rows.Next() {
		var id, cycle, currency, serviceFrom, serviceTo, paidAt, status, source, orderType, orderNo, notes, createdAt, origin, verificationStatus string
		var amount float64
		if err := rows.Scan(&id, &cycle, &amount, &currency, &serviceFrom, &serviceTo, &paidAt, &status, &source, &orderType, &orderNo, &notes, &createdAt, &origin, &verificationStatus); err != nil {
			return nil, err
		}
		items = append(items, map[string]any{
			"id": id, "billingCycle": cycle, "amount": amount, "currency": currency,
			"serviceFrom": serviceFrom, "serviceTo": serviceTo, "paidAt": nullableString(paidAt),
			"status": status, "source": source, "orderType": orderType, "orderNo": nullableString(orderNo), "notes": nullableString(notes), "createdAt": createdAt,
			"origin": origin, "verificationStatus": verificationStatus,
		})
	}
	return items, rows.Err()
}

// accruedUserIncome allocates confirmed payments over their service interval.
// Users without any non-cancelled billing history use the legacy monthly_fee
// fallback so existing installations keep their current financial totals after
// migration. A pending historical row intentionally blocks that fallback until
// the operator verifies or cancels it.
func (s *Server) accruedUserIncome(startText, endText string) (float64, error) {
	start, err := time.Parse(time.RFC3339Nano, startText)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse(time.RFC3339Nano, endText)
	if err != nil {
		return 0, err
	}
	// Confirmed records are commercial history. They remain accrued over their
	// service interval even after the associated operational user is removed.
	rows, err := s.db.Query(`SELECT amount, service_from, service_to
FROM user_billing_records
WHERE currency = 'CNY' AND status = 'confirmed' AND COALESCE(verification_status, 'verified') = 'verified'`)
	if err != nil {
		return 0, err
	}
	total := 0.0
	for rows.Next() {
		var amount float64
		var fromText, toText string
		if err := rows.Scan(&amount, &fromText, &toText); err != nil {
			_ = rows.Close()
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
			total += amount * overlapTo.Sub(overlapFrom).Seconds() / to.Sub(from).Seconds()
		}
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	// Installations that predate billing records keep their monthly fee
	// estimate. Deleted users are operationally gone, so only live users can
	// contribute this compatibility fallback.
	var legacyIncome float64
	if err := s.db.QueryRow(`SELECT COALESCE(SUM(u.monthly_fee), 0)
FROM users u
WHERE u.currency = 'CNY' AND u.deleted_at IS NULL AND u.status <> 'disabled'
  AND COALESCE(u.billing_type, 'paid') = 'paid'
  AND datetime(u.created_at) < datetime(?)
  AND (u.expiry_time IS NULL OR datetime(u.expiry_time) >= datetime(?))
  AND NOT EXISTS (
    SELECT 1 FROM user_billing_records b
    WHERE b.user_id = u.id AND b.currency = 'CNY' AND b.status <> 'cancelled'
  )`, end.UTC().Format(time.RFC3339Nano), start.UTC().Format(time.RFC3339Nano)).Scan(&legacyIncome); err != nil {
		return 0, err
	}
	return total + legacyIncome, nil
}
