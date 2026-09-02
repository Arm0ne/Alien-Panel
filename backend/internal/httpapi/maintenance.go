package httpapi

import "time"

const (
	userExpiringWindow = 7 * 24 * time.Hour
	nodeOfflineAfter   = 5 * time.Minute
)

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
	if _, err := s.db.Exec(`UPDATE nodes SET health_status = 'offline', updated_at = ?
WHERE enabled = 1 AND last_seen_at IS NOT NULL AND last_seen_at < ? AND health_status IN ('online', 'degraded', 'unknown')`, nowText, offlineBefore); err != nil {
		s.logger.Warn("refresh node health statuses", "error", err)
	}
}
