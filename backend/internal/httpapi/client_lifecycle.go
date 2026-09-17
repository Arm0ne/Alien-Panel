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

type clientReplacementEventPayload struct {
	InboundID      string `json:"inboundId"`
	UserID         string `json:"userId"`
	ClientSetHash  string `json:"clientSetHash"`
	OldClientCount int    `json:"oldClientCount"`
	NewClientCount int    `json:"newClientCount"`
}

type inboundUserResetResult struct {
	InboundID             string
	OldUserID             string
	NewUserID             string
	RemovedTrafficSamples int64
	ClosedPathCount       int64
	ClosedRouteCount      int64
}

// resetUserFromReplacementEvent confirms an Agent-detected customer change.
// The event payload is only a pointer to the Inbound and the previous user;
// all mutable state is re-read inside the same transaction before changing it.
func (s *Server) resetUserFromReplacementEvent(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimSpace(r.PathValue("id"))
	if eventID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "event id is required")
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin user reset")
		return
	}
	defer tx.Rollback()

	var eventType, eventStatus, resourceType, resourceID, payloadJSON string
	if err := tx.QueryRow(`SELECT event_type, event_status, COALESCE(resource_type, ''), COALESCE(resource_id, ''), COALESCE(payload_json, '')
FROM node_events WHERE id = ? AND visibility = 'public'`, eventID).
		Scan(&eventType, &eventStatus, &resourceType, &resourceID, &payloadJSON); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "event not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read replacement event")
		return
	}
	if eventType != "client_set_replacement_detected" || resourceType != "inbound" {
		writeFailure(w, http.StatusConflict, validationCode, "event is not a client replacement event")
		return
	}
	if eventStatus == "resolved" || eventStatus == "dismissed" {
		writeSuccess(w, map[string]any{"id": eventID, "resolved": true, "alreadyProcessed": true})
		return
	}
	var payload clientReplacementEventPayload
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "replacement event payload is invalid")
		return
	}
	inboundID := strings.TrimSpace(resourceID)
	if inboundID == "" {
		inboundID = strings.TrimSpace(payload.InboundID)
	}
	if inboundID == "" || strings.TrimSpace(payload.UserID) == "" {
		writeFailure(w, http.StatusConflict, validationCode, "replacement event is missing its user or inbound")
		return
	}

	result, err := s.resetInboundUserTx(tx, r, inboundID, payload.UserID, time.Now().UTC())
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusConflict, validationCode, "the inbound or its current user has changed; refresh and review it")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not reset inbound user")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	resolvedBy := ""
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok {
		resolvedBy = current.UserID
	}
	if _, err := tx.Exec(`UPDATE node_events SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
read_at = COALESCE(read_at, ?), resolved_at = ?, resolved_by = ? WHERE id = ?`, now, now, nullableDBString(resolvedBy), eventID); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not resolve replacement event")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit user reset")
		return
	}
	s.clearDashboardCache()
	writeSuccess(w, map[string]any{
		"id": eventID, "resolved": true, "inboundId": result.InboundID,
		"oldUserId": result.OldUserID, "newUserId": result.NewUserID,
		"removedTrafficSamples": result.RemovedTrafficSamples,
	})
}

// resetInboundUser is the explicit fallback for cases where the remote Client
// IDs were reused and no automatic replacement event could be generated.
func (s *Server) resetInboundUser(w http.ResponseWriter, r *http.Request) {
	inboundID := strings.TrimSpace(r.PathValue("id"))
	if inboundID == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "inbound id is required")
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not begin user reset")
		return
	}
	defer tx.Rollback()
	var oldUserID, inboundKind, nodeType string
	if err := tx.QueryRow(`SELECT COALESCE(i.user_id, ''), COALESCE(i.kind, ''), COALESCE(n.type, '')
FROM inbounds i JOIN nodes n ON n.id = i.node_id
WHERE i.id = ? AND i.deleted_at IS NULL`, inboundID).Scan(&oldUserID, &inboundKind, &nodeType); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "inbound not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read inbound")
		return
	}
	if nodeType != "relay" || inboundKind != "user" {
		writeFailure(w, http.StatusConflict, validationCode, "only a relay user inbound can be reset")
		return
	}
	if oldUserID == "" {
		writeFailure(w, http.StatusConflict, validationCode, "inbound has no current user")
		return
	}
	result, err := s.resetInboundUserTx(tx, r, inboundID, oldUserID, time.Now().UTC())
	if errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusConflict, validationCode, "the inbound or its current user has changed; refresh and review it")
		return
	}
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not reset inbound user")
		return
	}
	if err := tx.Commit(); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not commit user reset")
		return
	}
	s.clearDashboardCache()
	writeSuccess(w, map[string]any{
		"inboundId": result.InboundID, "oldUserId": result.OldUserID, "newUserId": result.NewUserID,
		"removedTrafficSamples": result.RemovedTrafficSamples,
	})
}

func (s *Server) resetInboundUserTx(tx *sql.Tx, r *http.Request, inboundID, expectedUserID string, now time.Time) (inboundUserResetResult, error) {
	result := inboundUserResetResult{InboundID: inboundID, OldUserID: expectedUserID}
	nowText := now.UTC().Format(time.RFC3339Nano)
	var currentUserID, inboundName, inboundKind, nodeType string
	var rawUp, rawDown, rawAllTime int64
	if err := tx.QueryRow(`SELECT COALESCE(i.user_id, ''), COALESCE(NULLIF(i.remark, ''), NULLIF(i.tag, ''), i.remote_inbound_id),
COALESCE(i.kind, ''), COALESCE(n.type, ''), i.up, i.down, i.all_time
FROM inbounds i JOIN nodes n ON n.id = i.node_id
WHERE i.id = ? AND i.deleted_at IS NULL`, inboundID).
		Scan(&currentUserID, &inboundName, &inboundKind, &nodeType, &rawUp, &rawDown, &rawAllTime); err != nil {
		return result, err
	}
	if currentUserID == "" || currentUserID != expectedUserID || inboundKind != "user" || nodeType != "relay" {
		return result, sql.ErrNoRows
	}

	if _, err := tx.Exec(`UPDATE user_inbounds SET active_to = ? WHERE user_id = ? AND inbound_id = ? AND active_to IS NULL`, nowText, expectedUserID, inboundID); err != nil {
		return result, fmt.Errorf("close inbound user mapping: %w", err)
	}

	// A user normally owns one relay Inbound, but the schema allows an
	// administrator to associate more than one.  Only retire user-level paths,
	// routes and state when this was the last live association; otherwise the
	// remaining Inbounds keep the old user's operational configuration.
	var remainingAssociations int
	if err := tx.QueryRow(`SELECT COUNT(*)
FROM user_inbounds ui JOIN inbounds i ON i.id = ui.inbound_id
WHERE ui.user_id = ? AND ui.active_to IS NULL AND i.id <> ? AND i.deleted_at IS NULL AND i.kind = 'user'`, expectedUserID, inboundID).Scan(&remainingAssociations); err != nil {
		return result, fmt.Errorf("count remaining user inbounds: %w", err)
	}
	lastAssociation := remainingAssociations == 0
	if lastAssociation {
		if count, err := tx.Exec(`UPDATE user_paths SET active_to = ?, updated_at = ? WHERE user_id = ? AND active_to IS NULL`, nowText, nowText, expectedUserID); err != nil {
			return result, fmt.Errorf("close user paths: %w", err)
		} else if affected, rowsErr := count.RowsAffected(); rowsErr == nil {
			result.ClosedPathCount = affected
		}
		if count, err := tx.Exec(`UPDATE user_routes SET is_primary = 0, active_to = ?, route_exit_ip_id = NULL WHERE user_id = ? AND active_to IS NULL`, nowText, expectedUserID); err != nil {
			return result, fmt.Errorf("close user routes: %w", err)
		} else if affected, rowsErr := count.RowsAffected(); rowsErr == nil {
			result.ClosedRouteCount = affected
		}
		resolvedBy := ""
		if current, ok := r.Context().Value(principalContextKey{}).(principal); ok {
			resolvedBy = current.UserID
		}
		if err := resolveNewUserProfileEventTx(tx, expectedUserID, resolvedBy, now); err != nil {
			return result, fmt.Errorf("resolve previous new user profile event: %w", err)
		}
	}

	// Billing rows remain immutable and continue to reference the retired user.
	candidateScope := `user_id = ?`
	candidateArgs := []any{expectedUserID}
	if !lastAssociation {
		candidateScope += ` AND inbound_id = ?`
		candidateArgs = append(candidateArgs, inboundID)
	}
	if _, err := tx.Exec(`UPDATE user_renewal_candidates SET status = 'rejected', processed_at = ?,
notes = CASE WHEN COALESCE(notes, '') = '' THEN '自动关闭：客户已更换' ELSE notes || char(10) || '自动关闭：客户已更换' END

WHERE `+candidateScope+` AND status = 'pending'`, append([]any{nowText}, candidateArgs...)...); err != nil {
		return result, fmt.Errorf("close renewal candidates: %w", err)
	}
	if _, err := tx.Exec(`UPDATE node_events SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?)
WHERE event_type = 'renewal_candidate_detected' AND resource_type = 'renewal' AND event_status NOT IN ('resolved', 'dismissed')
  AND resource_id IN (SELECT id FROM user_renewal_candidates WHERE `+candidateScope+` AND status = 'rejected')`, append([]any{nowText, nowText}, candidateArgs...)...); err != nil {
		return result, fmt.Errorf("close renewal events: %w", err)
	}

	removedTrafficSamples, err := deleteInboundTrafficTx(tx, inboundID)
	if err != nil {
		return result, fmt.Errorf("purge inbound traffic: %w", err)
	}
	result.RemovedTrafficSamples = removedTrafficSamples
	if _, err := tx.Exec(`UPDATE inbounds SET traffic_baseline_up = ?, traffic_baseline_down = ?, traffic_baseline_all_time = ?, traffic_baseline_at = ? WHERE id = ?`, rawUp, rawDown, rawAllTime, nowText, inboundID); err != nil {
		return result, fmt.Errorf("save traffic baseline: %w", err)
	}

	var newState businessUserState
	newState, err = clientStateFromDBTx(tx, inboundID, now)
	if err != nil {
		return result, fmt.Errorf("derive new user state: %w", err)
	}
	newUserID := newID()
	if _, err := tx.Exec(`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, newUserID, inboundName, newState.Status, nullableDBString(newState.ExpiryText), nowText, nowText); err != nil {
		return result, fmt.Errorf("create replacement user: %w", err)
	}
	if lastAssociation {
		if _, err := tx.Exec(`UPDATE users SET deleted_at = ?, status = 'disabled', updated_at = ? WHERE id = ? AND deleted_at IS NULL`, nowText, nowText, expectedUserID); err != nil {
			return result, fmt.Errorf("retire previous user: %w", err)
		}
	}
	if _, err := tx.Exec(`UPDATE inbounds SET user_id = ?, kind = 'user' WHERE id = ?`, newUserID, inboundID); err != nil {
		return result, fmt.Errorf("link replacement user: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES (?, ?, ?, 1, ?)`, newID(), newUserID, inboundID, nowText); err != nil {
		return result, fmt.Errorf("save replacement user mapping: %w", err)
	}
	if err := insertNewUserProfileEventTx(tx, inboundID, newUserID, inboundName, now); err != nil {
		return result, fmt.Errorf("record replacement user profile event: %w", err)
	}
	if err := s.writeAuditLogTx(tx, r, "user.replace_from_client_sync", "user", expectedUserID,
		map[string]any{"inboundId": inboundID, "displayName": inboundName},
		map[string]any{"newUserId": newUserID, "reason": "client_set_replaced", "trafficSamplesRemoved": result.RemovedTrafficSamples, "previousUserRetired": lastAssociation}, now); err != nil {
		return result, fmt.Errorf("write replacement audit: %w", err)
	}
	result.NewUserID = newUserID
	return result, nil
}

func deleteInboundTrafficTx(tx *sql.Tx, inboundID string) (int64, error) {
	result, err := tx.Exec(`DELETE FROM traffic_snapshots WHERE inbound_id = ?`, inboundID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func clientStateFromDBTx(tx *sql.Tx, inboundID string, observedAt time.Time) (businessUserState, error) {
	var inboundEnabled int
	if err := tx.QueryRow(`SELECT enable FROM inbounds WHERE id = ?`, inboundID).Scan(&inboundEnabled); err != nil {
		return businessUserState{}, err
	}
	rows, err := tx.Query(`SELECT enable, COALESCE(expiry_time, '') FROM clients WHERE inbound_id = ?`, inboundID)
	if err != nil {
		return businessUserState{}, err
	}
	defer rows.Close()
	clients := make([]agentClientPayload, 0)
	for rows.Next() {
		var enabled int
		var expiry string
		if err := rows.Scan(&enabled, &expiry); err != nil {
			return businessUserState{}, err
		}
		var expiryTime int64
		if expiry != "" {
			if parsed, parseErr := time.Parse(time.RFC3339Nano, expiry); parseErr == nil {
				expiryTime = parsed.Unix()
			}
		}
		clients = append(clients, agentClientPayload{Enable: enabled == 1, ExpiryTime: expiryTime, RemoteID: "replacement"})
	}
	if err := rows.Err(); err != nil {
		return businessUserState{}, err
	}
	return userStateFromClients(agentInboundPayload{Enable: inboundEnabled == 1, Clients: clients}, observedAt), nil
}
