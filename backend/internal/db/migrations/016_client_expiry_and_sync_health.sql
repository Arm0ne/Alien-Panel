-- A relay Inbound remains the central business-user identity. Its enabled
-- Clients now supply the actual access-expiry state, because one business
-- user may have several device credentials under the same Inbound.
ALTER TABLE nodes ADD COLUMN sync_status TEXT NOT NULL DEFAULT 'unknown'
  CHECK (sync_status IN ('unknown', 'success', 'failed'));
ALTER TABLE nodes ADD COLUMN last_sync_error TEXT;

UPDATE nodes
SET sync_status = CASE WHEN EXISTS (
  SELECT 1 FROM sync_runs sr WHERE sr.node_id = nodes.id AND sr.status = 'success'
) THEN 'success' ELSE 'unknown' END;

-- Keep the earliest finite expiry among enabled Clients. A NULL Client expiry
-- means that device has no expiry; MIN ignores it, so a mixture of unlimited
-- and finite Clients remains bounded by the finite device. Runtime code also
-- records that mixed configuration as an operator-visible mismatch.
UPDATE users
SET expiry_time = (
  SELECT MIN(c.expiry_time)
  FROM user_inbounds ui
  JOIN inbounds i ON i.id = ui.inbound_id
  JOIN clients c ON c.inbound_id = i.id
  JOIN nodes n ON n.id = i.node_id
  WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
    AND i.kind = 'user' AND i.deleted_at IS NULL AND i.enable = 1 AND n.type = 'relay' AND n.deleted_at IS NULL
    AND c.enable = 1 AND c.expiry_time IS NOT NULL
),
status = CASE
  WHEN NOT EXISTS (
    SELECT 1
    FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN clients c ON c.inbound_id = i.id
    JOIN nodes n ON n.id = i.node_id
    WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL AND i.enable = 1
      AND n.type = 'relay' AND n.deleted_at IS NULL AND c.enable = 1
  ) THEN 'disabled'
  WHEN (
    SELECT MIN(c.expiry_time)
    FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN clients c ON c.inbound_id = i.id
    JOIN nodes n ON n.id = i.node_id
    WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
      AND c.enable = 1 AND c.expiry_time IS NOT NULL
  ) IS NULL THEN 'active'
  WHEN datetime((
    SELECT MIN(c.expiry_time)
    FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN clients c ON c.inbound_id = i.id
    JOIN nodes n ON n.id = i.node_id
    WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
      AND c.enable = 1 AND c.expiry_time IS NOT NULL
  )) <= datetime('now') THEN 'expired'
  WHEN datetime((
    SELECT MIN(c.expiry_time)
    FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN clients c ON c.inbound_id = i.id
    JOIN nodes n ON n.id = i.node_id
    WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
      AND c.enable = 1 AND c.expiry_time IS NOT NULL
  )) <= datetime('now', '+7 days') THEN 'expiring'
  ELSE 'active'
END,
updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE EXISTS (
  SELECT 1
  FROM user_inbounds ui
  JOIN inbounds i ON i.id = ui.inbound_id
  JOIN nodes n ON n.id = i.node_id
  WHERE ui.user_id = users.id AND ui.is_primary = 1 AND ui.active_to IS NULL
    AND i.kind = 'user' AND i.deleted_at IS NULL AND n.type = 'relay' AND n.deleted_at IS NULL
);
