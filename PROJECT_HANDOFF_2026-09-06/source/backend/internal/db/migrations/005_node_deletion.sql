-- Keep deleted nodes and their operational history for audit, while removing
-- them from active node management and invalidating their Agent credentials.
ALTER TABLE nodes ADD COLUMN deleted_at TEXT;
CREATE INDEX IF NOT EXISTS idx_nodes_deleted_at ON nodes(deleted_at);
