-- Track the last complete Client identity set independently from the live
-- clients table.  This lets the Agent safely purge deleted Client rows while
-- still detecting a customer replacement that happened during an empty sync.
ALTER TABLE inbounds ADD COLUMN client_set_json TEXT NOT NULL DEFAULT '[]';

-- Xray keeps Inbound counters cumulative even when its Clients are replaced.
-- These values define the current customer's traffic baseline.
ALTER TABLE inbounds ADD COLUMN traffic_baseline_up INTEGER NOT NULL DEFAULT 0;
ALTER TABLE inbounds ADD COLUMN traffic_baseline_down INTEGER NOT NULL DEFAULT 0;
ALTER TABLE inbounds ADD COLUMN traffic_baseline_all_time INTEGER NOT NULL DEFAULT 0;
ALTER TABLE inbounds ADD COLUMN traffic_baseline_at TEXT;
