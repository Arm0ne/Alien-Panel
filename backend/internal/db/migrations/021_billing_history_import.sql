-- Distinguish live billing records from one-time historical migration data.
-- Historical rows can be imported before their payment evidence is verified;
-- unverified rows remain outside cash and accrued-income totals until confirmed.
ALTER TABLE user_billing_records ADD COLUMN origin TEXT NOT NULL DEFAULT 'live'
  CHECK (origin IN ('live', 'historical_import'));
ALTER TABLE user_billing_records ADD COLUMN verification_status TEXT NOT NULL DEFAULT 'verified'
  CHECK (verification_status IN ('verified', 'unverified'));

CREATE INDEX IF NOT EXISTS idx_user_billing_records_verification
  ON user_billing_records(verification_status, status, paid_at);
