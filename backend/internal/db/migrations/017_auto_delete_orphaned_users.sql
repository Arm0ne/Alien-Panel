-- Removing a node can leave a central business user without any current
-- relay or landing association. Keep the user row for immutable billing
-- history, but mark it deleted so operational pages no longer show it.
ALTER TABLE users ADD COLUMN deleted_at TEXT;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Resolve operator work before closing the candidate. This migration also
-- cleans up pre-existing orphaned users from older central installations.
UPDATE node_events
SET requires_action = 0,
    event_status = 'resolved',
    acknowledged = 1,
    read_at = COALESCE(read_at, strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    resolved_at = COALESCE(resolved_at, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
WHERE event_type = 'renewal_candidate_detected'
  AND resource_type = 'renewal'
  AND resource_id IN (
    SELECT c.id
    FROM user_renewal_candidates c
    JOIN users u ON u.id = c.user_id
    WHERE c.status = 'pending'
      AND u.deleted_at IS NULL
      AND NOT EXISTS (
        SELECT 1 FROM user_inbounds ui
        JOIN inbounds i ON i.id = ui.inbound_id
        JOIN nodes relay ON relay.id = i.node_id
        WHERE ui.user_id = u.id AND ui.active_to IS NULL
          AND i.kind = 'user' AND i.deleted_at IS NULL
          AND relay.type = 'relay' AND relay.deleted_at IS NULL
      )
      AND NOT EXISTS (
        SELECT 1 FROM user_paths p
        LEFT JOIN nodes relay ON relay.id = p.relay_node_id
        LEFT JOIN nodes landing ON landing.id = p.landing_node_id
        WHERE p.user_id = u.id AND p.active_to IS NULL
          AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
            OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
      )
      AND NOT EXISTS (
        SELECT 1 FROM user_routes ur
        JOIN routes route ON route.id = ur.route_id
        LEFT JOIN nodes relay ON relay.id = route.relay_node_id
        LEFT JOIN nodes landing ON landing.id = route.landing_node_id
        WHERE ur.user_id = u.id AND ur.active_to IS NULL
          AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
            OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
      )
  );

UPDATE user_renewal_candidates
SET status = 'rejected',
    processed_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    notes = CASE WHEN COALESCE(notes, '') = '' THEN '自动关闭：用户已无当前节点关联'
                 ELSE notes || char(10) || '自动关闭：用户已无当前节点关联' END
WHERE status = 'pending'
  AND user_id IN (
    SELECT u.id
    FROM users u
    WHERE u.deleted_at IS NULL
      AND NOT EXISTS (
        SELECT 1 FROM user_inbounds ui
        JOIN inbounds i ON i.id = ui.inbound_id
        JOIN nodes relay ON relay.id = i.node_id
        WHERE ui.user_id = u.id AND ui.active_to IS NULL
          AND i.kind = 'user' AND i.deleted_at IS NULL
          AND relay.type = 'relay' AND relay.deleted_at IS NULL
      )
      AND NOT EXISTS (
        SELECT 1 FROM user_paths p
        LEFT JOIN nodes relay ON relay.id = p.relay_node_id
        LEFT JOIN nodes landing ON landing.id = p.landing_node_id
        WHERE p.user_id = u.id AND p.active_to IS NULL
          AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
            OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
      )
      AND NOT EXISTS (
        SELECT 1 FROM user_routes ur
        JOIN routes route ON route.id = ur.route_id
        LEFT JOIN nodes relay ON relay.id = route.relay_node_id
        LEFT JOIN nodes landing ON landing.id = route.landing_node_id
        WHERE ur.user_id = u.id AND ur.active_to IS NULL
          AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
            OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
      )
  );

INSERT INTO audit_logs (id, admin_user_id, action, resource_type, resource_id, request_id, before_json, after_json, ip, created_at)
SELECT lower(hex(randomblob(16))), NULL, 'user.auto_delete', 'user', u.id, NULL, NULL,
       '{"deleted":true,"reason":"no_current_node_association","trigger":"migration"}', NULL,
       strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
FROM users u
WHERE u.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN nodes relay ON relay.id = i.node_id
    WHERE ui.user_id = u.id AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL
      AND relay.type = 'relay' AND relay.deleted_at IS NULL
  )
  AND NOT EXISTS (
    SELECT 1 FROM user_paths p
    LEFT JOIN nodes relay ON relay.id = p.relay_node_id
    LEFT JOIN nodes landing ON landing.id = p.landing_node_id
    WHERE p.user_id = u.id AND p.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  )
  AND NOT EXISTS (
    SELECT 1 FROM user_routes ur
    JOIN routes route ON route.id = ur.route_id
    LEFT JOIN nodes relay ON relay.id = route.relay_node_id
    LEFT JOIN nodes landing ON landing.id = route.landing_node_id
    WHERE ur.user_id = u.id AND ur.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  );

UPDATE users
SET deleted_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
    status = 'disabled',
    updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM user_inbounds ui
    JOIN inbounds i ON i.id = ui.inbound_id
    JOIN nodes relay ON relay.id = i.node_id
    WHERE ui.user_id = users.id AND ui.active_to IS NULL
      AND i.kind = 'user' AND i.deleted_at IS NULL
      AND relay.type = 'relay' AND relay.deleted_at IS NULL
  )
  AND NOT EXISTS (
    SELECT 1 FROM user_paths p
    LEFT JOIN nodes relay ON relay.id = p.relay_node_id
    LEFT JOIN nodes landing ON landing.id = p.landing_node_id
    WHERE p.user_id = users.id AND p.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  )
  AND NOT EXISTS (
    SELECT 1 FROM user_routes ur
    JOIN routes route ON route.id = ur.route_id
    LEFT JOIN nodes relay ON relay.id = route.relay_node_id
    LEFT JOIN nodes landing ON landing.id = route.landing_node_id
    WHERE ur.user_id = users.id AND ur.active_to IS NULL
      AND ((relay.type = 'relay' AND relay.deleted_at IS NULL)
        OR (landing.type = 'landing' AND landing.deleted_at IS NULL))
  );
