CREATE TABLE IF NOT EXISTS traffic_hourly_rollups (
  inbound_id TEXT NOT NULL REFERENCES inbounds(id) ON DELETE CASCADE,
  bucket_start TEXT NOT NULL,
  upload_bytes INTEGER NOT NULL DEFAULT 0,
  download_bytes INTEGER NOT NULL DEFAULT 0,
  sample_count INTEGER NOT NULL DEFAULT 0,
  observed_seconds INTEGER NOT NULL DEFAULT 0,
  reset_count INTEGER NOT NULL DEFAULT 0,
  has_gap INTEGER NOT NULL DEFAULT 0 CHECK (has_gap IN (0, 1)),
  data_at TEXT,
  updated_at TEXT NOT NULL,
  PRIMARY KEY (inbound_id, bucket_start)
);

CREATE INDEX IF NOT EXISTS idx_traffic_hourly_rollups_bucket
  ON traffic_hourly_rollups(bucket_start);

-- Backfill is deliberately performed by the application in small committed
-- batches. Keeping it out of this migration avoids one large startup
-- transaction on installations with years of minute-level snapshots.
CREATE TABLE IF NOT EXISTS traffic_rollup_state (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  status TEXT NOT NULL CHECK (status IN ('pending', 'rebuilding', 'complete')),
  processed_snapshots INTEGER NOT NULL DEFAULT 0,
  completed_at TEXT,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS traffic_rollup_dirty_buckets (
  inbound_id TEXT NOT NULL REFERENCES inbounds(id) ON DELETE CASCADE,
  bucket_start TEXT NOT NULL,
  PRIMARY KEY (inbound_id, bucket_start)
);

INSERT OR IGNORE INTO traffic_rollup_state (id, status, updated_at)
VALUES (1, 'pending', strftime('%Y-%m-%dT%H:%M:%fZ', 'now'));
