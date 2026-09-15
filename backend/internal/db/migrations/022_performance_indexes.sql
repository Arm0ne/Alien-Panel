-- The UNIQUE(inbound_id, collected_at) constraint already owns an equivalent
-- index. Keep the constraint and remove the redundant explicit index.
DROP INDEX IF EXISTS idx_traffic_snapshots_inbound_collected;

-- Common node detail/list queries find the most recent successful sync for a
-- node. These indexes also support the global latest-sync timestamp.
CREATE INDEX IF NOT EXISTS idx_sync_runs_node_status_time
  ON sync_runs(node_id, status, COALESCE(finished_at, started_at) DESC);
CREATE INDEX IF NOT EXISTS idx_sync_runs_status_time
  ON sync_runs(status, COALESCE(finished_at, started_at) DESC);
