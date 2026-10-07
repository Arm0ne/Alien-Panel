-- Resolved and dismissed events are historical records and must not remain
-- unread after an upgrade from versions that did not enforce this invariant.
UPDATE node_events
SET requires_action = 0,
    acknowledged = 1,
    read_at = COALESCE(read_at, resolved_at, created_at)
WHERE event_status IN ('resolved', 'dismissed')
  AND (acknowledged = 0 OR requires_action = 1 OR read_at IS NULL);
