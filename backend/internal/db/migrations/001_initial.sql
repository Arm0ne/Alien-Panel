CREATE TABLE IF NOT EXISTS admin_users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  totp_secret TEXT,
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  last_login_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  admin_user_id TEXT NOT NULL REFERENCES admin_users(id),
  token_hash TEXT NOT NULL UNIQUE,
  refresh_hash TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  revoked_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_hash ON sessions(refresh_hash);

CREATE TABLE IF NOT EXISTS nodes (
  id TEXT PRIMARY KEY,
  node_key TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'unknown',
  hostname TEXT,
  public_ip TEXT,
  region TEXT,
  provider TEXT,
  panel_base_path TEXT,
  agent_version TEXT,
  xpanel_version TEXT,
  xray_version TEXT,
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  health_status TEXT NOT NULL DEFAULT 'unknown',
  last_seen_at TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS node_credentials (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  token_hash TEXT NOT NULL UNIQUE,
  last_rotated_at TEXT NOT NULL,
  revoked_at TEXT,
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_node_credentials_node_id ON node_credentials(node_id);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'unknown',
  monthly_fee REAL NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'CNY',
  expiry_time TEXT,
  notes TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS inbounds (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  remote_inbound_id TEXT NOT NULL,
  user_id TEXT REFERENCES users(id),
  kind TEXT NOT NULL DEFAULT 'unknown',
  tag TEXT,
  remark TEXT,
  protocol TEXT,
  port INTEGER,
  listen TEXT,
  enable INTEGER NOT NULL DEFAULT 1 CHECK (enable IN (0, 1)),
  expiry_time TEXT,
  up INTEGER NOT NULL DEFAULT 0,
  down INTEGER NOT NULL DEFAULT 0,
  all_time INTEGER NOT NULL DEFAULT 0,
  client_count INTEGER NOT NULL DEFAULT 0,
  config_hash TEXT,
  first_seen_at TEXT NOT NULL,
  last_seen_at TEXT,
  missing_since TEXT,
  deleted_at TEXT,
  UNIQUE(node_id, remote_inbound_id)
);
CREATE INDEX IF NOT EXISTS idx_inbounds_user_id ON inbounds(user_id);
CREATE INDEX IF NOT EXISTS idx_inbounds_node_id ON inbounds(node_id);

CREATE TABLE IF NOT EXISTS clients (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  inbound_id TEXT NOT NULL REFERENCES inbounds(id),
  remote_client_id TEXT NOT NULL,
  email TEXT,
  enable INTEGER NOT NULL DEFAULT 1 CHECK (enable IN (0, 1)),
  expiry_time TEXT,
  up INTEGER NOT NULL DEFAULT 0,
  down INTEGER NOT NULL DEFAULT 0,
  all_time INTEGER NOT NULL DEFAULT 0,
  last_online TEXT,
  limit_ip INTEGER,
  last_seen_at TEXT,
  UNIQUE(node_id, inbound_id, remote_client_id)
);
CREATE INDEX IF NOT EXISTS idx_clients_inbound_id ON clients(inbound_id);

CREATE TABLE IF NOT EXISTS user_inbounds (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  inbound_id TEXT NOT NULL REFERENCES inbounds(id),
  is_primary INTEGER NOT NULL DEFAULT 1 CHECK (is_primary IN (0, 1)),
  active_from TEXT,
  active_to TEXT,
  UNIQUE(user_id, inbound_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_inbounds_one_primary ON user_inbounds(user_id) WHERE is_primary = 1 AND active_to IS NULL;

CREATE TABLE IF NOT EXISTS routes (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  relay_node_id TEXT NOT NULL REFERENCES nodes(id),
  landing_node_id TEXT NOT NULL REFERENCES nodes(id),
  relay_outbound_tag TEXT,
  landing_inbound_id TEXT REFERENCES inbounds(id),
  landing_inbound_tag TEXT,
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  valid_from TEXT,
  valid_to TEXT,
  notes TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_routes (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  route_id TEXT NOT NULL REFERENCES routes(id),
  is_primary INTEGER NOT NULL DEFAULT 1 CHECK (is_primary IN (0, 1)),
  active_from TEXT,
  active_to TEXT,
  UNIQUE(user_id, route_id)
);

CREATE TABLE IF NOT EXISTS exit_ips (
  id TEXT PRIMARY KEY,
  landing_node_id TEXT NOT NULL REFERENCES nodes(id),
  ip TEXT NOT NULL,
  family INTEGER NOT NULL DEFAULT 4,
  provider TEXT,
  monthly_cost REAL NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'CNY',
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  valid_from TEXT,
  valid_to TEXT,
  notes TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  UNIQUE(landing_node_id, ip)
);

CREATE TABLE IF NOT EXISTS route_exit_ips (
  id TEXT PRIMARY KEY,
  route_id TEXT NOT NULL REFERENCES routes(id),
  exit_ip_id TEXT NOT NULL REFERENCES exit_ips(id),
  allocation_weight REAL NOT NULL DEFAULT 1,
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  UNIQUE(route_id, exit_ip_id)
);

CREATE TABLE IF NOT EXISTS node_costs (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  category TEXT NOT NULL,
  monthly_amount REAL NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'CNY',
  effective_from TEXT NOT NULL,
  effective_to TEXT,
  notes TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS other_costs (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  category TEXT NOT NULL,
  monthly_amount REAL NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'CNY',
  effective_from TEXT NOT NULL,
  effective_to TEXT,
  notes TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS traffic_snapshots (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  inbound_id TEXT NOT NULL REFERENCES inbounds(id),
  collected_at TEXT NOT NULL,
  up INTEGER NOT NULL DEFAULT 0,
  down INTEGER NOT NULL DEFAULT 0,
  all_time INTEGER NOT NULL DEFAULT 0,
  source TEXT NOT NULL DEFAULT 'xpanel',
  reset_detected INTEGER NOT NULL DEFAULT 0 CHECK (reset_detected IN (0, 1)),
  sync_run_id TEXT,
  UNIQUE(inbound_id, collected_at)
);
CREATE INDEX IF NOT EXISTS idx_traffic_snapshots_collected_at ON traffic_snapshots(collected_at);

CREATE TABLE IF NOT EXISTS sync_runs (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  sync_id TEXT NOT NULL UNIQUE,
  started_at TEXT NOT NULL,
  finished_at TEXT,
  status TEXT NOT NULL,
  inbound_count INTEGER NOT NULL DEFAULT 0,
  client_count INTEGER NOT NULL DEFAULT 0,
  error_message TEXT
);

CREATE TABLE IF NOT EXISTS node_events (
  id TEXT PRIMARY KEY,
  node_id TEXT REFERENCES nodes(id),
  event_type TEXT NOT NULL,
  severity TEXT NOT NULL DEFAULT 'info',
  message TEXT NOT NULL,
  created_at TEXT NOT NULL,
  acknowledged INTEGER NOT NULL DEFAULT 0 CHECK (acknowledged IN (0, 1))
);
CREATE INDEX IF NOT EXISTS idx_node_events_created_at ON node_events(created_at);

CREATE TABLE IF NOT EXISTS audit_logs (
  id TEXT PRIMARY KEY,
  admin_user_id TEXT REFERENCES admin_users(id),
  action TEXT NOT NULL,
  resource_type TEXT,
  resource_id TEXT,
  request_id TEXT,
  before_json TEXT,
  after_json TEXT,
  ip TEXT,
  created_at TEXT NOT NULL
);
