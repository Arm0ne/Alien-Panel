ALTER TABLE inbounds ADD COLUMN missing_sync_count INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_inbounds_missing_archive
ON inbounds(node_id, deleted_at, missing_sync_count);
