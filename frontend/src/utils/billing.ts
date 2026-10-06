/** Keep service boundaries as instants; the date picker displays them in local time. */
export function billingTimestamp(value?: string | null): number | null {
  if (!value) return null;
  const timestamp = Date.parse(value);
  return Number.isFinite(timestamp) ? timestamp : null;
}

export function billingTimestampISO(value: number | null): string {
  if (value === null || !Number.isFinite(value)) return '';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? '' : date.toISOString();
}

export function billingServiceInterval(
  records: ReadonlyArray<{ status: string; serviceTo: string }>,
  expiresAt?: string | null
) {
  let serviceFrom: number | null = null;
  for (const record of records) {
    // Match the backend overlap guard: pending history also reserves its interval.
    if (record.status === 'cancelled') continue;
    const end = billingTimestamp(record.serviceTo);
    if (end !== null && (serviceFrom === null || end > serviceFrom)) serviceFrom = end;
  }
  return { serviceFrom, serviceTo: billingTimestamp(expiresAt) };
}
