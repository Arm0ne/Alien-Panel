-- Commercial classification and immutable order metadata used by the
-- operations finance view. Existing users remain paid for compatibility;
-- operators can explicitly mark friends/free accounts as free.
ALTER TABLE users ADD COLUMN billing_type TEXT NOT NULL DEFAULT 'paid'
  CHECK (billing_type IN ('paid', 'free'));
ALTER TABLE users ADD COLUMN free_reason TEXT;

ALTER TABLE user_billing_records ADD COLUMN order_type TEXT NOT NULL DEFAULT 'initial'
  CHECK (order_type IN ('initial', 'renewal', 'recovery'));
ALTER TABLE user_billing_records ADD COLUMN order_no TEXT;
UPDATE user_billing_records
SET order_type = 'renewal'
WHERE renewal_candidate_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_billing_records_order_no
  ON user_billing_records(order_no) WHERE order_no IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_user_billing_records_order_type
  ON user_billing_records(order_type, paid_at);
