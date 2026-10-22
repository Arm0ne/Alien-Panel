-- X-Panel exposes a per-client traffic quota (totalGB, in GB). Store it as
-- bytes so the panel can render usage/quota without unit conversion, and keep
-- 0 as "unknown / unlimited" for older Agents that do not report it.
ALTER TABLE clients ADD COLUMN quota_bytes INTEGER NOT NULL DEFAULT 0;
