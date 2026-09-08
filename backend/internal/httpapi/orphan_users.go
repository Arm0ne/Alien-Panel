package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// nodeAffectedUserIDsTx returns users whose current inbound, direct path, or
// legacy route is about to be removed with a node. The final orphan decision
// is made only after those records have been deleted.
func nodeAffectedUserIDsTx(tx *sql.Tx, nodeID string) ([]string, error) {
	rows, err := tx.Query(`SELECT DISTINCT user_id FROM (
SELECT ui.user_id
FROM user_inbounds ui
JOIN inbounds i ON i.id = ui.inbound_id
WHERE i.node_id = ?
UNION
SELECT p.user_id
FROM user_paths p
WHERE p.relay_node_id = ? OR p.landing_node_id = ?
   OR p.landing_inbound_id IN (SELECT id FROM inbounds WHERE node_id = ?)
   OR p.exit_ip_id IN (SELECT id FROM exit_ips WHERE owner_node_id = ? OR landing_node_id = ?)
   OR EXISTS (SELECT 1 FROM user_path_exit_ips upi JOIN exit_ips e ON e.id = upi.exit_ip_id WHERE upi.user_path_id = p.id AND (e.owner_node_id = ? OR e.landing_node_id = ?))
UNION
SELECT ur.user_id
FROM user_routes ur
JOIN routes r ON r.id = ur.route_id
WHERE r.relay_node_id = ? OR r.landing_node_id = ?
)`, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userIDs, nil
}

// autoDeleteOrphanedUsersTx logically deletes only the users made orphaned by
// the current node deletion. Billing records keep their user foreign key and
// therefore remain immutable financial history.
func (s *Server) autoDeleteOrphanedUsersTx(tx *sql.Tx, r *http.Request, userIDs []string, triggerNodeID, triggerNodeName string, now time.Time) error {
	seen := make(map[string]struct{}, len(userIDs))
	nowText := now.UTC().Format(time.RFC3339Nano)
	for _, userID := range userIDs {
		if _, duplicate := seen[userID]; duplicate {
			continue
		}
		seen[userID] = struct{}{}

		var displayName, status, expiry, billingCycle, currency string
		var monthlyFee, billingAmount float64
		err := tx.QueryRow(`SELECT u.display_name, u.status, COALESCE(u.expiry_time, ''), u.monthly_fee,
COALESCE(u.billing_cycle, 'monthly'), COALESCE(u.billing_amount, u.monthly_fee), u.currency
FROM users u
WHERE u.id = ? AND u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN nodes relay ON relay.id = i.node_id
    WHERE ui.user_id = u.id AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL
      AND relay.type = 'relay' AND relay.deleted_at IS NULL
  )
  AND NOT EXISTS (
    SELECT 1
    FROM user_paths p
    LEFT JOIN nodes relay ON relay.id = p.relay_node_id
    LEFT JOIN nodes landing ON landing.id = p.landing_node_id
    WHERE p.user_id = u.id AND p.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  )
  AND NOT EXISTS (
    SELECT 1
    FROM user_routes ur
    JOIN routes route ON route.id = ur.route_id
    LEFT JOIN nodes relay ON relay.id = route.relay_node_id
    LEFT JOIN nodes landing ON landing.id = route.landing_node_id
    WHERE ur.user_id = u.id AND ur.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  )`, userID).Scan(&displayName, &status, &expiry, &monthlyFee, &billingCycle, &billingAmount, &currency)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return fmt.Errorf("find orphaned user %s: %w", userID, err)
		}

		var pendingCandidates, billingRecordCount int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM user_renewal_candidates WHERE user_id = ? AND status = 'pending'`, userID).Scan(&pendingCandidates); err != nil {
			return fmt.Errorf("count pending renewal candidates for user %s: %w", userID, err)
		}
		if err := tx.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = ?`, userID).Scan(&billingRecordCount); err != nil {
			return fmt.Errorf("count billing records for user %s: %w", userID, err)
		}

		// Close operator work before changing candidate state so only the
		// pending candidate events are resolved. Existing confirmed/rejected
		// events remain their original immutable history.
		if _, err := tx.Exec(`UPDATE node_events
SET requires_action = 0, event_status = 'resolved', acknowledged = 1,
    read_at = COALESCE(read_at, ?), resolved_at = COALESCE(resolved_at, ?)
WHERE event_type = 'renewal_candidate_detected' AND resource_type = 'renewal'
  AND resource_id IN (
    SELECT id FROM user_renewal_candidates WHERE user_id = ? AND status = 'pending'
  )`, nowText, nowText, userID); err != nil {
			return fmt.Errorf("resolve renewal events for user %s: %w", userID, err)
		}
		if _, err := tx.Exec(`UPDATE user_renewal_candidates
SET status = 'rejected', processed_at = ?,
    notes = CASE WHEN COALESCE(notes, '') = '' THEN ? ELSE notes || char(10) || ? END
WHERE user_id = ? AND status = 'pending'`, nowText, "自动关闭：用户已无当前节点关联", "自动关闭：用户已无当前节点关联", userID); err != nil {
			return fmt.Errorf("close renewal candidates for user %s: %w", userID, err)
		}
		if _, err := tx.Exec(`UPDATE users
SET deleted_at = ?, status = 'disabled', updated_at = ?
WHERE id = ? AND deleted_at IS NULL`, nowText, nowText, userID); err != nil {
			return fmt.Errorf("logically delete user %s: %w", userID, err)
		}

		before := map[string]any{
			"id": userID, "displayName": displayName, "status": status,
			"expiresAt": nullableString(expiry), "monthlyFee": monthlyFee,
			"billingCycle": billingCycle, "billingAmount": billingAmount, "currency": currency,
		}
		after := map[string]any{
			"deleted": true, "deletedAt": nowText, "reason": "no_current_node_association",
			"triggerNode":             map[string]any{"id": triggerNodeID, "name": triggerNodeName},
			"closedRenewalCandidates": pendingCandidates, "billingRecordCount": billingRecordCount,
		}
		if err := s.writeAuditLogTx(tx, r, "user.auto_delete", "user", userID, before, after, now); err != nil {
			return fmt.Errorf("write user auto-delete audit for %s: %w", userID, err)
		}
	}
	return nil
}

// writeAuditLogTx is the transactional counterpart of writeAuditLog. Auto
// deletion changes the node graph and user lifecycle together, so its audit
// entry is part of the same atomic operation rather than a best-effort write.
func (s *Server) writeAuditLogTx(tx *sql.Tx, r *http.Request, action, resourceType, resourceID string, before, after any, createdAt time.Time) error {
	beforeJSON, err := marshalAuditState(before)
	if err != nil {
		return err
	}
	afterJSON, err := marshalAuditState(after)
	if err != nil {
		return err
	}

	var adminUserID any
	if current, ok := r.Context().Value(principalContextKey{}).(principal); ok && current.UserID != "" {
		adminUserID = current.UserID
	}
	var requestID any
	if value, ok := r.Context().Value(requestIDContextKey{}).(string); ok && value != "" {
		requestID = value
	}
	var ip any
	if value := clientIP(r.RemoteAddr); value != "" {
		ip = value
	}
	_, err = tx.Exec(`INSERT INTO audit_logs (id, admin_user_id, action, resource_type, resource_id, request_id, before_json, after_json, ip, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		newID(), adminUserID, action, resourceType, resourceID, requestID, string(beforeJSON), string(afterJSON), ip, createdAt.UTC().Format(time.RFC3339Nano))
	return err
}
