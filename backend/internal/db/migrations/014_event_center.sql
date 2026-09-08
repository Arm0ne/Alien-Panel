-- Turn node_events into an operator-facing event stream while retaining the
-- existing table and event identifiers for backwards compatibility.
ALTER TABLE node_events ADD COLUMN event_category TEXT NOT NULL DEFAULT 'system';
ALTER TABLE node_events ADD COLUMN title TEXT NOT NULL DEFAULT '';
ALTER TABLE node_events ADD COLUMN visibility TEXT NOT NULL DEFAULT 'public';
ALTER TABLE node_events ADD COLUMN requires_action INTEGER NOT NULL DEFAULT 0;
ALTER TABLE node_events ADD COLUMN event_status TEXT NOT NULL DEFAULT 'open';
ALTER TABLE node_events ADD COLUMN resource_type TEXT;
ALTER TABLE node_events ADD COLUMN resource_id TEXT;
ALTER TABLE node_events ADD COLUMN action_type TEXT;
ALTER TABLE node_events ADD COLUMN payload_json TEXT;
ALTER TABLE node_events ADD COLUMN read_at TEXT;
ALTER TABLE node_events ADD COLUMN resolved_at TEXT;
ALTER TABLE node_events ADD COLUMN resolved_by TEXT;
ALTER TABLE node_events ADD COLUMN dedupe_key TEXT;
ALTER TABLE node_events ADD COLUMN source TEXT NOT NULL DEFAULT 'system';
ALTER TABLE node_events ADD COLUMN correlation_id TEXT;

UPDATE node_events
SET title = CASE event_type
  WHEN 'traffic_reset' THEN '流量累计值发生回退'
  WHEN 'inbound_missing' THEN 'Inbound 暂时缺失'
  WHEN 'inbound_archived' THEN 'Inbound 已归档'
  WHEN 'sync_requested' THEN '立即同步请求'
  WHEN 'sync_failed' THEN '节点同步失败'
  ELSE event_type
END
WHERE title = '';

UPDATE node_events
SET event_category = CASE
  WHEN event_type LIKE 'renewal_%' THEN 'business'
  WHEN event_type IN ('traffic_reset', 'inbound_missing', 'inbound_archived') THEN 'node'
  WHEN event_type LIKE 'sync_%' THEN 'sync'
  ELSE 'system'
END;

UPDATE node_events
SET visibility = 'internal', source = 'admin'
WHERE event_type = 'sync_requested';

CREATE INDEX IF NOT EXISTS idx_node_events_visibility_created
  ON node_events(visibility, created_at);
CREATE INDEX IF NOT EXISTS idx_node_events_action_status
  ON node_events(requires_action, event_status, created_at);
CREATE INDEX IF NOT EXISTS idx_node_events_resource
  ON node_events(resource_type, resource_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_node_events_dedupe
  ON node_events(dedupe_key)
  WHERE dedupe_key IS NOT NULL;
