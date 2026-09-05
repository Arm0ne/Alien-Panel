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

func (s *Server) ensurePendingRenewalEvents() error {
	rows, err := s.db.Query(`SELECT c.id, c.user_id, COALESCE(u.display_name, ''), COALESCE(c.inbound_id, ''),
COALESCE(i.node_id, ''), c.old_expiry_at, c.new_expiry_at, c.suggested_cycle, c.suggested_amount, c.currency, c.detected_at
FROM user_renewal_candidates c
LEFT JOIN users u ON u.id = c.user_id
LEFT JOIN inbounds i ON i.id = c.inbound_id
WHERE c.status = 'pending'`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type pendingRenewal struct {
		candidateID, userID, userName, inboundID, nodeID, oldExpiry, newExpiry, cycle, currency, detectedAt string
		amount                                                                                              float64
	}
	pending := make([]pendingRenewal, 0)
	for rows.Next() {
		var item pendingRenewal
		if err := rows.Scan(&item.candidateID, &item.userID, &item.userName, &item.inboundID, &item.nodeID, &item.oldExpiry, &item.newExpiry, &item.cycle, &item.amount, &item.currency, &item.detectedAt); err != nil {
			return err
		}
		pending = append(pending, item)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range pending {
		userName := item.userName
		if userName == "" {
			userName = item.userID
		}
		occurredAt, parseErr := time.Parse(time.RFC3339Nano, item.detectedAt)
		if parseErr != nil {
			occurredAt = time.Now().UTC()
		}
		if err := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: item.nodeID, EventType: "renewal_candidate_detected", Category: "business", Severity: "warning",
			Title: "检测到续费变更，待确认", Message: fmt.Sprintf("用户「%s」到期时间由 %s 延长至 %s，请确认是否计入收费", userName, item.oldExpiry, item.newExpiry),
			RequiresAction: true, EventStatus: "open", ResourceType: "renewal", ResourceID: item.candidateID,
			ActionType: "confirm_renewal", Payload: map[string]any{"candidateId": item.candidateID, "userId": item.userID, "userName": userName, "inboundId": nullableString(item.inboundID), "oldExpiryAt": item.oldExpiry, "newExpiryAt": item.newExpiry, "billingCycle": item.cycle, "suggestedAmount": item.amount, "currency": item.currency},
			DedupeKey: "renewal-candidate:" + item.candidateID, Source: "agent", CorrelationID: item.candidateID, OccurredAt: occurredAt,
		}); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// nodeEventSpec is the small, operator-facing event contract. Technical
// sync_runs remain the source of truth for scheduler/execution diagnostics;
// node_events is reserved for meaningful state changes and operator work.
type nodeEventSpec struct {
	ID             string
	NodeID         string
	EventType      string
	Category       string
	Severity       string
	Title          string
	Message        string
	Visibility     string
	RequiresAction bool
	EventStatus    string
	ResourceType   string
	ResourceID     string
	ActionType     string
	Payload        any
	DedupeKey      string
	Source         string
	CorrelationID  string
	OccurredAt     time.Time
}

func insertNodeEventTx(tx *sql.Tx, event nodeEventSpec) error {
	if event.ID == "" {
		event.ID = newID()
	}
	if event.Category == "" {
		event.Category = "system"
	}
	if event.Severity == "" {
		event.Severity = "info"
	}
	if event.Title == "" {
		event.Title = event.EventType
	}
	if event.Visibility == "" {
		event.Visibility = "public"
	}
	if event.EventStatus == "" {
		event.EventStatus = "open"
	}
	if event.Source == "" {
		event.Source = "system"
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	var payloadJSON any
	if event.Payload != nil {
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}
		payloadJSON = string(payload)
	}
	_, err := tx.Exec(`INSERT OR IGNORE INTO node_events
(id, node_id, event_type, severity, message, created_at, event_category, title, visibility, requires_action, event_status,
 resource_type, resource_id, action_type, payload_json, dedupe_key, source, correlation_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, nullableDBString(event.NodeID), event.EventType, event.Severity, event.Message,
		event.OccurredAt.UTC().Format(time.RFC3339Nano), event.Category, event.Title, event.Visibility,
		boolInt(event.RequiresAction), event.EventStatus, nullableDBString(event.ResourceType), nullableDBString(event.ResourceID),
		nullableDBString(event.ActionType), payloadJSON, nullableDBString(event.DedupeKey), event.Source, nullableDBString(event.CorrelationID))
	return err
}

func insertNodeEvent(database *sql.DB, event nodeEventSpec) error {
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := insertNodeEventTx(tx, event); err != nil {
		return err
	}
	return tx.Commit()
}

func resolveRenewalEventTx(tx *sql.Tx, candidateID, eventType string, payload any, occurredAt time.Time, resolvedBy string) error {
	now := occurredAt.UTC().Format(time.RFC3339Nano)
	if _, err := tx.Exec(`UPDATE node_events
SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
    read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?), resolved_by = COALESCE(resolved_by, ?)
WHERE event_type = 'renewal_candidate_detected' AND resource_type = 'renewal' AND resource_id = ?`, now, now, nullableDBString(resolvedBy), candidateID); err != nil {
		return err
	}
	var nodeID string
	_ = tx.QueryRow(`SELECT COALESCE(node_id, '') FROM node_events
WHERE event_type = 'renewal_candidate_detected' AND resource_type = 'renewal' AND resource_id = ? LIMIT 1`, candidateID).Scan(&nodeID)
	title := "续费已确认"
	message := "续费候选已确认并计入财务"
	if eventType == "renewal_non_billable" {
		title = "已标记为非收费变更"
		message = "到期时间变更已标记为非收费变更"
	}
	return insertNodeEventTx(tx, nodeEventSpec{
		NodeID: nodeID, EventType: eventType, Category: "business", Severity: "info", Title: title, Message: message,
		EventStatus: "resolved", ResourceType: "renewal", ResourceID: candidateID, Payload: payload,
		DedupeKey: eventType + ":" + candidateID, Source: "admin", CorrelationID: candidateID, OccurredAt: occurredAt,
	})
}

func (s *Server) markEventRead(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "event id is required")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := s.db.Exec(`UPDATE node_events SET acknowledged = 1, read_at = COALESCE(read_at, ?)
WHERE id = ? AND visibility = 'public'`, now, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not mark event as read")
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		writeFailure(w, http.StatusNotFound, notFoundCode, "event not found")
		return
	}
	writeSuccess(w, map[string]any{"id": id, "read": true})
}

func (s *Server) markAllEventsRead(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE node_events SET acknowledged = 1, read_at = COALESCE(read_at, ?)
WHERE visibility = 'public' AND acknowledged = 0`, now); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not mark events as read")
		return
	}
	writeSuccess(w, map[string]any{"read": true})
}

func (s *Server) resolveEvent(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeFailure(w, http.StatusBadRequest, validationCode, "event id is required")
		return
	}
	var eventType string
	if err := s.db.QueryRow(`SELECT event_type FROM node_events WHERE id = ? AND visibility = 'public'`, id).Scan(&eventType); errors.Is(err, sql.ErrNoRows) {
		writeFailure(w, http.StatusNotFound, notFoundCode, "event not found")
		return
	} else if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not read event")
		return
	}
	if eventType == "renewal_candidate_detected" || strings.HasPrefix(eventType, "renewal_") {
		writeFailure(w, http.StatusConflict, validationCode, "renewal events must be confirmed or marked as non-billable")
		return
	}
	current, ok := r.Context().Value(principalContextKey{}).(principal)
	resolvedBy := ""
	if ok {
		resolvedBy = current.UserID
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := s.db.Exec(`UPDATE node_events SET event_status = 'resolved', resolved_at = ?, resolved_by = ?, acknowledged = 1, read_at = COALESCE(read_at, ?)
WHERE id = ? AND visibility = 'public' AND event_status NOT IN ('resolved', 'dismissed')`, now, nullableDBString(resolvedBy), now, id)
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not resolve event")
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		writeSuccess(w, map[string]any{"id": id, "resolved": true})
		return
	}
	s.writeAuditLog(r, "event.resolve", "event", id, nil, map[string]any{"eventType": eventType})
	writeSuccess(w, map[string]any{"id": id, "resolved": true})
}

func (s *Server) eventSummary(w http.ResponseWriter, _ *http.Request) {
	var pending, unread int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE visibility = 'public' AND requires_action = 1 AND event_status NOT IN ('resolved', 'dismissed')`).Scan(&pending); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count pending events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE visibility = 'public' AND acknowledged = 0 AND event_status NOT IN ('resolved', 'dismissed')`).Scan(&unread); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count unread events")
		return
	}
	writeSuccess(w, map[string]any{"pendingCount": pending, "unreadCount": unread, "generatedAt": time.Now().UTC().Format(time.RFC3339Nano)})
}
