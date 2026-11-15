package httpapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	centraldb "xpanel-central/backend/internal/db"
)

// pendingEventFilter returns the shared predicate for operator-actionable
// events. Older events may have requires_action=0 even though they carry an
// action_type, so both signals are considered while resolved and dismissed
// events remain excluded.
func pendingEventFilter(alias string) string {
	if strings.TrimSpace(alias) == "" {
		return "(requires_action = 1 OR NULLIF(TRIM(action_type), '') IS NOT NULL) AND event_status NOT IN ('resolved', 'dismissed')"
	}
	prefix := strings.TrimSpace(alias) + "."
	return "(" + prefix + "requires_action = 1 OR NULLIF(TRIM(" + prefix + "action_type), '') IS NOT NULL) AND " + prefix + "event_status NOT IN ('resolved', 'dismissed')"
}

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
	acknowledged := false
	readAt := ""
	if event.EventStatus == "resolved" || event.EventStatus == "dismissed" {
		// A completed event is a history record, not a new notification.
		event.RequiresAction = false
		acknowledged = true
		readAt = event.OccurredAt.UTC().Format(time.RFC3339Nano)
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
 resource_type, resource_id, action_type, payload_json, dedupe_key, source, correlation_id, acknowledged, read_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, nullableDBString(event.NodeID), event.EventType, event.Severity, event.Message,
		event.OccurredAt.UTC().Format(time.RFC3339Nano), event.Category, event.Title, event.Visibility,
		boolInt(event.RequiresAction), event.EventStatus, nullableDBString(event.ResourceType), nullableDBString(event.ResourceID),
		nullableDBString(event.ActionType), payloadJSON, nullableDBString(event.DedupeKey), event.Source, nullableDBString(event.CorrelationID),
		boolInt(acknowledged), nullableDBString(readAt))
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

func insertNewUserProfileEventTx(tx *sql.Tx, inboundID, userID, userName string, occurredAt time.Time) error {
	var nodeID, inboundName string
	if err := tx.QueryRow(`SELECT i.node_id,
COALESCE(NULLIF(i.remark, ''), NULLIF(i.tag, ''), i.remote_inbound_id)
FROM inbounds i WHERE i.id = ?`, inboundID).Scan(&nodeID, &inboundName); err != nil {
		return fmt.Errorf("read new user event context: %w", err)
	}
	if strings.TrimSpace(userName) == "" {
		userName = inboundName
	}
	payload := map[string]any{
		"userId": userID, "userName": userName,
		"inboundId": inboundID, "inboundName": inboundName,
	}
	return insertNodeEventTx(tx, nodeEventSpec{
		NodeID: nodeID, EventType: "new_user_profile_required", Category: "business", Severity: "warning",
		Title:          "发现新增用户，待完善资料",
		Message:        fmt.Sprintf("线路机新增用户「%s」（Inbound「%s」），请完善收费类型、金额和业务资料。", userName, inboundName),
		RequiresAction: true, EventStatus: "open", ResourceType: "user", ResourceID: userID, ActionType: "complete_user_profile",
		Payload: payload, DedupeKey: "new-user-profile:" + userID, Source: "agent", CorrelationID: inboundID, OccurredAt: occurredAt,
	})
}

func resolveNewUserProfileEventTx(tx *sql.Tx, userID, resolvedBy string, resolvedAt time.Time) error {
	now := resolvedAt.UTC().Format(time.RFC3339Nano)
	_, err := tx.Exec(`UPDATE node_events
SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
    read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?), resolved_by = COALESCE(resolved_by, ?)
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user' AND resource_id = ?
  AND event_status NOT IN ('resolved', 'dismissed')`, now, now, nullableDBString(resolvedBy), userID)
	return err
}

// refreshNewUserProfileEventTx 在 Inbound 改名或用户重命名后同步事件文案，
// 避免事件中心一直显示旧名称。
func refreshNewUserProfileEventTx(tx *sql.Tx, inboundID, userID, userName string) error {
	var nodeID, inboundName string
	if err := tx.QueryRow(`SELECT i.node_id,
COALESCE(NULLIF(i.remark, ''), NULLIF(i.tag, ''), i.remote_inbound_id)
FROM inbounds i WHERE i.id = ?`, inboundID).Scan(&nodeID, &inboundName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("read profile event context: %w", err)
	}
	if strings.TrimSpace(userName) == "" {
		userName = inboundName
	}
	payload, err := json.Marshal(map[string]any{
		"userId": userID, "userName": userName,
		"inboundId": inboundID, "inboundName": inboundName,
	})
	if err != nil {
		return fmt.Errorf("encode profile event payload: %w", err)
	}
	message := fmt.Sprintf("线路机新增用户「%s」（Inbound「%s」），请完善收费类型、金额和业务资料。", userName, inboundName)
	if _, err := tx.Exec(`UPDATE node_events
SET node_id = ?, message = ?, payload_json = ?
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user' AND resource_id = ?
  AND COALESCE(event_status, 'open') NOT IN ('resolved', 'dismissed')`,
		nullableDBString(nodeID), message, string(payload), userID); err != nil {
		return fmt.Errorf("refresh new user profile event: %w", err)
	}
	return nil
}

// dismissUnusableNewUserProfileEventsTx 关闭「主入站已不可用」的待办事件。
//
// 之前这些事件只有等用户因失去全部节点关联被回收时才关闭：如果 Inbound 只是
// 被停用、或者用户还挂着路径，用户不会被回收，事件就会永久停在待处理，点进去
// 还会因为「没有可用的主 Inbound」而报错。此处在每次同步后统一收敛，
// 并清空 dedupe_key，使 Inbound 恢复可用后仍能重新提醒。
func dismissUnusableNewUserProfileEventsTx(tx *sql.Tx, occurredAt time.Time) (int64, error) {
	nowText := occurredAt.UTC().Format(time.RFC3339Nano)
	result, err := tx.Exec(`UPDATE node_events
SET requires_action = 0, event_status = 'dismissed', acknowledged = 1,
    read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?),
    dedupe_key = NULL,
    message = message || '（主入站已停用或删除，事件已自动关闭）'
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user'
  AND COALESCE(event_status, 'open') NOT IN ('resolved', 'dismissed')
  AND NOT EXISTS (
    SELECT 1
    FROM users u
    JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
    JOIN inbounds i ON i.id = ui.inbound_id AND i.deleted_at IS NULL AND i.enable = 1 AND i.kind = 'user'
    JOIN nodes n ON n.id = i.node_id AND n.deleted_at IS NULL AND n.enabled = 1 AND n.type = 'relay'
    WHERE u.id = node_events.resource_id AND u.deleted_at IS NULL
  )`, nowText, nowText)
	if err != nil {
		return 0, fmt.Errorf("dismiss unusable new user profile events: %w", err)
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

// reopenDismissedNewUserProfileEventTx 在 Inbound 恢复可用后重新提醒。
// 只针对此前被系统因「主入站不可用」而关闭的事件，且用户仍未定价。
func reopenDismissedNewUserProfileEventTx(tx *sql.Tx, inboundID, userID, userName string, occurredAt time.Time) error {
	var dismissed int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user' AND resource_id = ?
  AND COALESCE(event_status, 'open') = 'dismissed' AND dedupe_key IS NULL`, userID).Scan(&dismissed); err != nil {
		return fmt.Errorf("read dismissed profile event: %w", err)
	}
	if dismissed == 0 {
		return nil
	}
	// 只有主入站确实可用时才重新提醒，否则会和 dismissUnusableNewUserProfileEventsTx
	// 形成「关闭 → 重开 → 再关闭」的每分钟循环。
	var unpriced int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM users u
JOIN user_inbounds ui ON ui.user_id = u.id AND ui.is_primary = 1 AND ui.active_to IS NULL
JOIN inbounds i ON i.id = ui.inbound_id AND i.deleted_at IS NULL AND i.enable = 1 AND i.kind = 'user'
JOIN nodes n ON n.id = i.node_id AND n.deleted_at IS NULL AND n.enabled = 1 AND n.type = 'relay'
WHERE u.id = ? AND u.deleted_at IS NULL
  AND COALESCE(u.billing_type, 'paid') = 'paid'
  AND COALESCE(u.billing_amount, u.monthly_fee, 0) <= 0
  AND NOT EXISTS (SELECT 1 FROM user_billing_records b WHERE b.user_id = u.id)
  AND NOT EXISTS (
    SELECT 1 FROM node_events e
    WHERE e.event_type = 'new_user_profile_required' AND e.resource_type = 'user' AND e.resource_id = u.id
      AND COALESCE(e.event_status, 'open') NOT IN ('resolved', 'dismissed')
  )`, userID).Scan(&unpriced); err != nil {
		return fmt.Errorf("check unpriced user: %w", err)
	}
	if unpriced == 0 {
		return nil
	}
	return insertNewUserProfileEventTx(tx, inboundID, userID, userName, occurredAt)
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
WHERE visibility = 'public' AND acknowledged = 0
  AND event_status NOT IN ('resolved', 'dismissed')
  AND requires_action = 0 AND NULLIF(TRIM(action_type), '') IS NULL`, now); err != nil {
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
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events WHERE visibility = 'public' AND ` + pendingEventFilter("")).Scan(&pending); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count pending events")
		return
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE visibility = 'public' AND acknowledged = 0 AND event_status NOT IN ('resolved', 'dismissed')
  AND requires_action = 0 AND NULLIF(TRIM(action_type), '') IS NULL`).Scan(&unread); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "could not count unread events")
		return
	}
	writeSuccess(w, map[string]any{"pendingCount": pending, "unreadCount": unread, "generatedAt": time.Now().UTC().Format(time.RFC3339Nano)})
}

// pruneEvents 清理已处理（resolved / dismissed）的历史事件。
// 运营可以用它主动回收事件中心的体积，未处理事件永远不会被删除。
func (s *Server) pruneEvents(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		OlderThanDays int `json:"olderThanDays"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid prune payload")
		return
	}
	days := payload.OlderThanDays
	if days < 1 || days > 3650 {
		days = 30
	}
	result, err := centraldb.PruneHistoricalData(r.Context(), s.db, centraldb.RetentionOptions{EventDays: days})
	if err != nil {
		s.logger.Error("prune settled events", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "清理事件失败")
		return
	}
	s.writeAuditLog(r, "events.prune", "event", "settled", nil, map[string]any{"olderThanDays": days, "deleted": result.Events})
	writeSuccess(w, map[string]any{"deleted": result.Events, "olderThanDays": days})
}
