-- Short-lived, one-time credentials used by the online Agent installer.
-- The installer exchanges one for a normal node credential after the local
-- X-Panel credentials have been entered on the target machine.
CREATE TABLE IF NOT EXISTS node_install_tokens (
  id TEXT PRIMARY KEY,
  node_id TEXT NOT NULL REFERENCES nodes(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  used_at TEXT,
  created_by TEXT REFERENCES admin_users(id),
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_node_install_tokens_node_id
  ON node_install_tokens(node_id);
CREATE INDEX IF NOT EXISTS idx_node_install_tokens_active
  ON node_install_tokens(token_hash, used_at, expires_at);
