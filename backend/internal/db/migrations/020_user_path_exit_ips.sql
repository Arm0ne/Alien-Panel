-- Allow one business path to use an ordered set of fixed exit IP assets.
-- user_paths.exit_ip_id remains the first item for legacy readers and writes.
CREATE TABLE IF NOT EXISTS user_path_exit_ips (
  user_path_id TEXT NOT NULL REFERENCES user_paths(id) ON DELETE CASCADE,
  exit_ip_id TEXT NOT NULL REFERENCES exit_ips(id),
  position INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
  PRIMARY KEY (user_path_id, exit_ip_id)
);

CREATE INDEX IF NOT EXISTS idx_user_path_exit_ips_exit_ip_id
  ON user_path_exit_ips(exit_ip_id);

INSERT OR IGNORE INTO user_path_exit_ips (user_path_id, exit_ip_id, position, created_at)
SELECT id, exit_ip_id, 0, COALESCE(created_at, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
FROM user_paths;
