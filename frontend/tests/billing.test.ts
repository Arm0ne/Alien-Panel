import assert from 'node:assert/strict';
import { test } from 'node:test';
import { billingServiceInterval, billingTimestamp, billingTimestampISO } from '../src/utils/billing';

test('renewal starts at the previous afternoon deadline and ends at the synced instant', () => {
  const interval = billingServiceInterval(
    [{ status: 'confirmed', serviceTo: '2026-10-06T14:30:45+08:00' }],
    '2026-11-06T14:30:45+08:00'
  );
  assert.equal(billingTimestampISO(interval.serviceFrom), '2026-10-06T06:30:45.000Z');
  assert.equal(billingTimestampISO(interval.serviceTo), '2026-11-06T06:30:45.000Z');
});

test('service boundaries survive local-midnight crossings and preserve milliseconds', () => {
  const interval = billingServiceInterval(
    [{ status: 'confirmed', serviceTo: '2026-10-06T00:15:20.123+08:00' }],
    '2026-11-06T00:15:20.123+08:00'
  );
  assert.equal(billingTimestampISO(interval.serviceFrom), '2026-10-05T16:15:20.123Z');
  assert.equal(billingTimestampISO(interval.serviceTo), '2026-11-05T16:15:20.123Z');
});

test('choose the latest non-cancelled end by instant across offsets and unordered history', () => {
  const interval = billingServiceInterval([
    { status: 'confirmed', serviceTo: '2026-10-06T14:00:00+08:00' },
    { status: 'cancelled', serviceTo: '2027-01-01T00:00:00Z' },
    { status: 'pending', serviceTo: '2026-10-06T07:00:00Z' },
    { status: 'confirmed', serviceTo: '2026-09-06T06:30:00Z' }
  ]);
  assert.equal(billingTimestampISO(interval.serviceFrom), '2026-10-06T07:00:00.000Z');
  assert.equal(interval.serviceTo, null);
});

test('missing or invalid timestamps require input instead of inventing a service period', () => {
  assert.deepEqual(billingServiceInterval([], null), { serviceFrom: null, serviceTo: null });
  assert.deepEqual(billingServiceInterval([{ status: 'cancelled', serviceTo: '2027-01-01T00:00:00Z' }]), {
    serviceFrom: null,
    serviceTo: null
  });
  assert.equal(billingTimestamp('invalid'), null);
  assert.equal(billingTimestampISO(null), '');
  assert.equal(billingTimestampISO(Number.NaN), '');
});
