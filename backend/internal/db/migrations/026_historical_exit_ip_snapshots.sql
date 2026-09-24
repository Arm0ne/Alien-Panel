-- Keep the address of retired exit assets available in historical paths.
-- The main exit_ip_id remains a live assignment reference; historical rows
-- are not physically rewritten by this migration.
ALTER TABLE user_paths ADD COLUMN exit_ip_address_snapshot TEXT;

UPDATE user_paths
SET exit_ip_address_snapshot = (
  SELECT e.ip FROM exit_ips e WHERE e.id = user_paths.exit_ip_id
)
WHERE exit_ip_address_snapshot IS NULL;
