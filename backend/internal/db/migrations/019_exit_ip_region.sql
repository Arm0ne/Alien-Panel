-- Store a manually assigned region for independent S5 exit IPs.  Node-owned
-- exit IPs continue to derive their region from the current owner node.
ALTER TABLE exit_ips ADD COLUMN region TEXT;

CREATE INDEX IF NOT EXISTS idx_exit_ips_region ON exit_ips(region);
