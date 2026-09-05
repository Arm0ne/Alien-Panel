-- Direct per-user network paths.  The relay node is copied from the user's
-- primary inbound at assignment time; landing and exit assets are selected by
-- an administrator.  The old routes/user_routes tables remain available as
-- read-only history and for backwards-compatible API clients.
CREATE TABLE IF NOT EXISTS user_paths (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  relay_node_id TEXT NOT NULL REFERENCES nodes(id),
  landing_node_id TEXT REFERENCES nodes(id),
  landing_inbound_id TEXT REFERENCES inbounds(id),
  exit_ip_id TEXT NOT NULL REFERENCES exit_ips(id),
  mode TEXT NOT NULL CHECK (mode IN ('relay', 'landing', 'external')),
  notes TEXT,
  active_from TEXT NOT NULL,
  active_to TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_paths_one_active
  ON user_paths(user_id) WHERE active_to IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_paths_relay_node_id ON user_paths(relay_node_id);
CREATE INDEX IF NOT EXISTS idx_user_paths_landing_node_id ON user_paths(landing_node_id);
CREATE INDEX IF NOT EXISTS idx_user_paths_exit_ip_id ON user_paths(exit_ip_id);

-- Migrate only unambiguous fixed legacy assignments.  Pool assignments are
-- intentionally left for manual confirmation instead of choosing an IP.
INSERT INTO user_paths (
  id, user_id, relay_node_id, landing_node_id, landing_inbound_id,
  exit_ip_id, mode, notes, active_from, active_to, created_at, updated_at
)
SELECT
  'legacy-path-' || ur.id,
  ur.user_id,
  r.relay_node_id,
  CASE WHEN rei.scope = 'landing' THEN r.landing_node_id ELSE NULL END,
  CASE WHEN rei.scope = 'landing' THEN r.landing_inbound_id ELSE NULL END,
  rei.exit_ip_id,
  rei.scope,
  '从旧线路分配迁移',
  COALESCE(ur.active_from, r.created_at),
  ur.active_to,
  r.created_at,
  r.updated_at
FROM user_routes ur
JOIN routes r ON r.id = ur.route_id
JOIN route_exit_ips rei ON rei.id = ur.route_exit_ip_id
JOIN exit_ips e ON e.id = rei.exit_ip_id
WHERE ur.is_primary = 1
  AND rei.enabled = 1
  AND e.enabled = 1
  AND ur.active_to IS NULL
  AND NOT EXISTS (SELECT 1 FROM user_paths p WHERE p.user_id = ur.user_id AND p.active_to IS NULL);
