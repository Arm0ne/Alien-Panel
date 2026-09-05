package httpapi

import (
	"context"
	"time"
)

const (
	userExpiringWindow         = 7 * 24 * time.Hour
	nodeOfflineAfter           = 5 * time.Minute
	missingInboundArchiveAfter = 3
)

// RunMaintenance keeps time-based operational states current even when no
// administrator is opening dashboard or list pages. It never removes data;
// Inbound archival is handled transactionally during a successful node sync.
func (s *Server) RunMaintenance(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	run := func() {
		s.refreshOperationalStatuses(time.Now().UTC())
	}
	run()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			run()
		}
	}
}

func (s *Server) refreshOperationalStatuses(now time.Time) {
	nowText := now.UTC().Format(time.RFC3339Nano)
	expiringText := now.UTC().Add(userExpiringWindow).Format(time.RFC3339Nano)
	if _, err := s.db.Exec(`UPDATE users SET status = CASE
WHEN expiry_time IS NULL THEN status
WHEN expiry_time <= ? THEN 'expired'
WHEN expiry_time <= ? THEN 'expiring'
ELSE 'active' END, updated_at = ?
WHERE status IN ('unknown', 'active', 'expiring', 'expired')`, nowText, expiringText, nowText); err != nil {
		s.logger.Warn("refresh user expiry statuses", "error", err)
	}
	offlineBefore := now.UTC().Add(-nodeOfflineAfter).Format(time.RFC3339Nano)
	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Warn("begin node health refresh", "error", err)
		return
	}
	defer tx.Rollback()
	rows, err := tx.Query(`SELECT id, name FROM nodes
WHERE deleted_at IS NULL AND enabled = 1 AND last_seen_at IS NOT NULL AND last_seen_at < ?
  AND health_status IN ('online', 'degraded', 'unknown')`, offlineBefore)
	if err != nil {
		s.logger.Warn("find offline nodes", "error", err)
		return
	}
	type staleNode struct{ id, name string }
	stale := make([]staleNode, 0)
	for rows.Next() {
		var node staleNode
		if scanErr := rows.Scan(&node.id, &node.name); scanErr != nil {
			_ = rows.Close()
			s.logger.Warn("read offline node", "error", scanErr)
			return
		}
		stale = append(stale, node)
	}
	if err := rows.Close(); err != nil {
		s.logger.Warn("close offline node query", "error", err)
		return
	}
	if _, err := tx.Exec(`UPDATE nodes SET health_status = 'offline', updated_at = ?
WHERE deleted_at IS NULL AND enabled = 1 AND last_seen_at IS NOT NULL AND last_seen_at < ? AND health_status IN ('online', 'degraded', 'unknown')`, nowText, offlineBefore); err != nil {
		s.logger.Warn("refresh node health statuses", "error", err)
		return
	}
	for _, node := range stale {
		if err := insertNodeEventTx(tx, nodeEventSpec{
			NodeID: node.id, EventType: "node_offline", Category: "node", Severity: "warning", Title: "节点已离线",
			Message: "节点「" + node.name + "」超过 5 分钟未收到 Agent 心跳", ResourceType: "node", ResourceID: node.id,
			ActionType: "inspect_node", Source: "system", OccurredAt: now,
		}); err != nil {
			s.logger.Warn("record node offline event", "node_id", node.id, "error", err)
		}
	}
	if err := tx.Commit(); err != nil {
		s.logger.Warn("commit node health refresh", "error", err)
	}
}
