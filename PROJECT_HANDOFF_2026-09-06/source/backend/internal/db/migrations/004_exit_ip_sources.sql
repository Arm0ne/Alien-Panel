-- Exit IP assets can be hosted by either a relay/landing node or purchased
-- independently as an external S5 endpoint.  Rebuild the two small tables so
-- the legacy landing_node_id column becomes nullable while old data remains
-- readable by older clients.
CREATE TABLE exit_ips_new (
  id TEXT PRIMARY KEY,
  landing_node_id TEXT REFERENCES nodes(id),
  source_type TEXT NOT NULL DEFAULT 'node' CHECK (source_type IN ('node', 's5')),
  owner_node_id TEXT REFERENCES nodes(id),
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

INSERT INTO exit_ips_new (
  id, landing_node_id, source_type, owner_node_id, ip, family, provider,
  monthly_cost, currency, enabled, valid_from, valid_to, notes, created_at, updated_at
)
SELECT id, landing_node_id, 'node', landing_node_id, ip, family, provider,
  monthly_cost, currency, enabled, valid_from, valid_to, notes, created_at, updated_at
FROM exit_ips;

CREATE TABLE route_exit_ips_new (
  id TEXT PRIMARY KEY,
  route_id TEXT NOT NULL REFERENCES routes(id),
  exit_ip_id TEXT NOT NULL REFERENCES exit_ips_new(id),
  scope TEXT NOT NULL DEFAULT 'landing' CHECK (scope IN ('relay', 'landing', 'external')),
  allocation_weight REAL NOT NULL DEFAULT 1,
  enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  UNIQUE(route_id, exit_ip_id)
);

INSERT INTO route_exit_ips_new (id, route_id, exit_ip_id, scope, allocation_weight, enabled)
SELECT id, route_id, exit_ip_id, 'landing', allocation_weight, enabled
FROM route_exit_ips;

DROP TABLE route_exit_ips;
DROP TABLE exit_ips;
ALTER TABLE exit_ips_new RENAME TO exit_ips;
ALTER TABLE route_exit_ips_new RENAME TO route_exit_ips;

CREATE UNIQUE INDEX IF NOT EXISTS idx_exit_ips_owner_ip
  ON exit_ips(source_type, owner_node_id, ip);
CREATE INDEX IF NOT EXISTS idx_exit_ips_owner_node_id
  ON exit_ips(owner_node_id);
CREATE INDEX IF NOT EXISTS idx_exit_ips_source_type
  ON exit_ips(source_type);
CREATE INDEX IF NOT EXISTS idx_route_exit_ips_scope
  ON route_exit_ips(route_id, scope);
