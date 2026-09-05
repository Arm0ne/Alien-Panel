-- Preserve the commercial terms of a user independently from the X-Panel
-- snapshot.  Existing users are treated as monthly users for compatibility.
ALTER TABLE users ADD COLUMN billing_cycle TEXT NOT NULL DEFAULT 'monthly'
  CHECK (billing_cycle IN ('monthly', 'annual'));
ALTER TABLE users ADD COLUMN billing_amount REAL NOT NULL DEFAULT 0;

UPDATE users
SET billing_amount = monthly_fee
WHERE billing_amount = 0 AND monthly_fee <> 0;

CREATE INDEX IF NOT EXISTS idx_users_billing_cycle ON users(billing_cycle);

-- A confirmed billing record is immutable commercial history.  The service
-- interval is half-open: [service_from, service_to).
CREATE TABLE IF NOT EXISTS user_billing_records (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  billing_cycle TEXT NOT NULL CHECK (billing_cycle IN ('monthly', 'annual')),
  amount REAL NOT NULL CHECK (amount >= 0),
  currency TEXT NOT NULL DEFAULT 'CNY',
  service_from TEXT NOT NULL,
  service_to TEXT NOT NULL,
  paid_at TEXT,
  status TEXT NOT NULL DEFAULT 'confirmed' CHECK (status IN ('pending', 'confirmed', 'rejected', 'cancelled')),
  source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'agent_detected')),
  renewal_candidate_id TEXT,
  notes TEXT,
  created_at TEXT NOT NULL,
  UNIQUE(renewal_candidate_id)
);
CREATE INDEX IF NOT EXISTS idx_user_billing_records_user_service
  ON user_billing_records(user_id, service_from, service_to);
CREATE INDEX IF NOT EXISTS idx_user_billing_records_paid_at
  ON user_billing_records(paid_at);

-- Agent only suggests a renewal when an existing Inbound expiry moves
-- forward.  Operators confirm whether that change represented a payment.
CREATE TABLE IF NOT EXISTS user_renewal_candidates (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  inbound_id TEXT REFERENCES inbounds(id),
  old_expiry_at TEXT NOT NULL,
  new_expiry_at TEXT NOT NULL,
  detected_at TEXT NOT NULL,
  suggested_cycle TEXT NOT NULL CHECK (suggested_cycle IN ('monthly', 'annual')),
  suggested_amount REAL NOT NULL DEFAULT 0 CHECK (suggested_amount >= 0),
  currency TEXT NOT NULL DEFAULT 'CNY',
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'rejected')),
  processed_at TEXT,
  notes TEXT,
  UNIQUE(user_id, inbound_id, new_expiry_at)
);
CREATE INDEX IF NOT EXISTS idx_user_renewal_candidates_user_status
  ON user_renewal_candidates(user_id, status, detected_at);
