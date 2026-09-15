-- Migration 022 was already applied on some installations before the
-- redundant traffic index was removed from that migration. Keep this as a
-- new migration so those databases are upgraded as well.
DROP INDEX IF EXISTS idx_traffic_snapshots_inbound_collected;
