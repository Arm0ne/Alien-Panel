-- Full browser/X-Panel management URL for a node, including scheme, port and
-- optional panel base path.  The legacy public_ip column remains available for
-- old clients and historical records; new UI flows use management_url.
ALTER TABLE nodes ADD COLUMN management_url TEXT;

CREATE INDEX IF NOT EXISTS idx_nodes_management_url
  ON nodes(management_url);
