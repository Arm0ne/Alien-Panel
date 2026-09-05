CREATE INDEX IF NOT EXISTS idx_traffic_snapshots_inbound_collected
ON traffic_snapshots(inbound_id, collected_at);
