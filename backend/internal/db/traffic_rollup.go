package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	trafficRollupBatchSize = 5000
	trafficRollupGap       = 72 * time.Hour
)

type TrafficRollupProgress struct {
	ProcessedSnapshots int64
}

type trafficRollupSnapshot struct {
	inboundID string
	at        time.Time
	up        int64
	down      int64
	reset     bool
}

type trafficHourlyRollup struct {
	inboundID       string
	bucketStart     time.Time
	uploadBytes     int64
	downloadBytes   int64
	sampleCount     int64
	observedSeconds int64
	resetCount      int64
	hasGap          bool
	dataAt          time.Time
}

// TrafficHourStart returns the UTC hour containing the supplied timestamp.
func TrafficHourStart(value time.Time) time.Time {
	return value.UTC().Truncate(time.Hour)
}

func TrafficHourlyRollupsReady(database *sql.DB) (bool, error) {
	var status string
	err := database.QueryRow(`SELECT status FROM traffic_rollup_state WHERE id = 1`).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return status == "complete", err
}

// EnsureTrafficHourlyRollups builds historical hourly traffic data in small,
// independently committed batches. An interrupted build is safely restarted
// from an empty rollup table on the next launch.
func EnsureTrafficHourlyRollups(ctx context.Context, database *sql.DB, progress func(TrafficRollupProgress)) error {
	ready, err := TrafficHourlyRollupsReady(database)
	if err != nil {
		return fmt.Errorf("read traffic rollup state: %w", err)
	}
	if ready {
		return nil
	}

	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	transaction, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin traffic rollup rebuild: %w", err)
	}
	if _, err := transaction.ExecContext(ctx, `DELETE FROM traffic_hourly_rollups`); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("clear traffic rollups: %w", err)
	}
	if _, err := transaction.ExecContext(ctx, `UPDATE traffic_rollup_state
SET status = 'rebuilding', processed_snapshots = 0, completed_at = NULL, updated_at = ? WHERE id = 1`, nowText); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("mark traffic rollup rebuild: %w", err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("start traffic rollup rebuild: %w", err)
	}

	var lastInbound, lastCollected string
	var previous trafficRollupSnapshot
	var hasPrevious bool
	var processed int64
	for {
		rows, err := database.QueryContext(ctx, `SELECT inbound_id, collected_at, up, down, reset_detected
FROM traffic_snapshots
WHERE inbound_id > ? OR (inbound_id = ? AND collected_at > ?)
ORDER BY inbound_id, collected_at
LIMIT ?`, lastInbound, lastInbound, lastCollected, trafficRollupBatchSize)
		if err != nil {
			return fmt.Errorf("read traffic snapshot batch: %w", err)
		}

		batch := make([]trafficRollupSnapshot, 0, trafficRollupBatchSize)
		for rows.Next() {
			var snapshot trafficRollupSnapshot
			var collected string
			var reset int
			if err := rows.Scan(&snapshot.inboundID, &collected, &snapshot.up, &snapshot.down, &reset); err != nil {
				_ = rows.Close()
				return fmt.Errorf("scan traffic snapshot batch: %w", err)
			}
			parsed, err := time.Parse(time.RFC3339Nano, collected)
			if err != nil {
				_ = rows.Close()
				return fmt.Errorf("parse traffic snapshot time %q: %w", collected, err)
			}
			snapshot.at = parsed.UTC()
			snapshot.reset = reset == 1
			batch = append(batch, snapshot)
			lastInbound, lastCollected = snapshot.inboundID, collected
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return fmt.Errorf("iterate traffic snapshot batch: %w", err)
		}
		if err := rows.Close(); err != nil {
			return fmt.Errorf("close traffic snapshot batch: %w", err)
		}
		if len(batch) == 0 {
			break
		}

		rollups := make(map[string]*trafficHourlyRollup)
		for _, snapshot := range batch {
			if !hasPrevious || snapshot.inboundID != previous.inboundID {
				previous = snapshot
				hasPrevious = true
				continue
			}
			addTrafficRollupDelta(rollups, previous, snapshot)
			previous = snapshot
		}

		transaction, err := database.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin traffic rollup batch: %w", err)
		}
		for _, rollup := range rollups {
			if err := addTrafficHourlyRollup(ctx, transaction, rollup); err != nil {
				_ = transaction.Rollback()
				return err
			}
		}
		processed += int64(len(batch))
		if _, err := transaction.ExecContext(ctx, `UPDATE traffic_rollup_state
SET processed_snapshots = ?, updated_at = ? WHERE id = 1`, processed, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			_ = transaction.Rollback()
			return fmt.Errorf("save traffic rollup progress: %w", err)
		}
		if err := transaction.Commit(); err != nil {
			return fmt.Errorf("commit traffic rollup batch: %w", err)
		}
		if progress != nil {
			progress(TrafficRollupProgress{ProcessedSnapshots: processed})
		}
	}

	// Agent synchronization may have committed new current-hour snapshots while
	// older inbound IDs were already behind the keyset cursor. Recompute the two
	// live edge buckets before making the rollup readable.
	transaction, err = database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin traffic rollup edge repair: %w", err)
	}
	rows, err := transaction.QueryContext(ctx, `SELECT id FROM inbounds`)
	if err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("read traffic rollup edge inbounds: %w", err)
	}
	inboundIDs := make([]string, 0)
	for rows.Next() {
		var inboundID string
		if err := rows.Scan(&inboundID); err != nil {
			_ = rows.Close()
			_ = transaction.Rollback()
			return fmt.Errorf("scan traffic rollup edge inbound: %w", err)
		}
		inboundIDs = append(inboundIDs, inboundID)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		_ = transaction.Rollback()
		return fmt.Errorf("iterate traffic rollup edge inbounds: %w", err)
	}
	if err := rows.Close(); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("close traffic rollup edge inbounds: %w", err)
	}
	currentHour := TrafficHourStart(time.Now())
	rows, err = transaction.QueryContext(ctx, `SELECT inbound_id, bucket_start FROM traffic_rollup_dirty_buckets`)
	if err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("read dirty traffic rollup buckets: %w", err)
	}
	type dirtyBucket struct {
		inboundID string
		start     time.Time
	}
	dirtyBuckets := make([]dirtyBucket, 0)
	for rows.Next() {
		var inboundID, bucketText string
		if err := rows.Scan(&inboundID, &bucketText); err != nil {
			_ = rows.Close()
			_ = transaction.Rollback()
			return fmt.Errorf("scan dirty traffic rollup bucket: %w", err)
		}
		bucketStart, err := time.Parse(time.RFC3339Nano, bucketText)
		if err != nil {
			_ = rows.Close()
			_ = transaction.Rollback()
			return fmt.Errorf("parse dirty traffic rollup bucket: %w", err)
		}
		dirtyBuckets = append(dirtyBuckets, dirtyBucket{inboundID: inboundID, start: bucketStart})
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		_ = transaction.Rollback()
		return fmt.Errorf("iterate dirty traffic rollup buckets: %w", err)
	}
	if err := rows.Close(); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("close dirty traffic rollup buckets: %w", err)
	}
	seenDirty := make(map[string]struct{}, len(dirtyBuckets))
	for _, bucket := range dirtyBuckets {
		key := bucket.inboundID + "\x00" + bucket.start.Format(time.RFC3339Nano)
		if _, ok := seenDirty[key]; ok {
			continue
		}
		seenDirty[key] = struct{}{}
		if err := RebuildTrafficHourlyBucketTx(ctx, transaction, bucket.inboundID, bucket.start); err != nil {
			_ = transaction.Rollback()
			return fmt.Errorf("repair dirty traffic rollup bucket: %w", err)
		}
		if _, err := transaction.ExecContext(ctx, `DELETE FROM traffic_rollup_dirty_buckets WHERE inbound_id = ? AND bucket_start = ?`, bucket.inboundID, bucket.start.Format(time.RFC3339Nano)); err != nil {
			_ = transaction.Rollback()
			return fmt.Errorf("clear dirty traffic rollup bucket: %w", err)
		}
	}
	for _, inboundID := range inboundIDs {
		for _, bucketStart := range []time.Time{currentHour.Add(-time.Hour), currentHour} {
			if err := RebuildTrafficHourlyBucketTx(ctx, transaction, inboundID, bucketStart); err != nil {
				_ = transaction.Rollback()
				return fmt.Errorf("repair traffic rollup edge: %w", err)
			}
			if _, err := transaction.ExecContext(ctx, `DELETE FROM traffic_rollup_dirty_buckets WHERE inbound_id = ? AND bucket_start = ?`, inboundID, bucketStart.Format(time.RFC3339Nano)); err != nil {
				_ = transaction.Rollback()
				return fmt.Errorf("clear edge traffic rollup bucket: %w", err)
			}
		}
	}
	completedAt := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := transaction.ExecContext(ctx, `UPDATE traffic_rollup_state
SET status = 'complete', processed_snapshots = ?, completed_at = ?, updated_at = ? WHERE id = 1`, processed, completedAt, completedAt); err != nil {
		_ = transaction.Rollback()
		return fmt.Errorf("complete traffic rollup rebuild: %w", err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit traffic rollup edge repair: %w", err)
	}
	return nil
}

func addTrafficRollupDelta(rollups map[string]*trafficHourlyRollup, previous, snapshot trafficRollupSnapshot) {
	interval := snapshot.at.Sub(previous.at)
	if interval <= 0 {
		return
	}
	resetDetected := snapshot.reset || snapshot.up < previous.up || snapshot.down < previous.down
	uploadDelta := snapshot.up - previous.up
	downloadDelta := snapshot.down - previous.down
	if resetDetected {
		uploadDelta, downloadDelta = snapshot.up, snapshot.down
	}
	if uploadDelta < 0 {
		uploadDelta = 0
	}
	if downloadDelta < 0 {
		downloadDelta = 0
	}
	bucketStart := TrafficHourStart(snapshot.at)
	key := snapshot.inboundID + "\x00" + bucketStart.Format(time.RFC3339Nano)
	rollup := rollups[key]
	if rollup == nil {
		rollup = &trafficHourlyRollup{inboundID: snapshot.inboundID, bucketStart: bucketStart}
		rollups[key] = rollup
	}
	rollup.uploadBytes += uploadDelta
	rollup.downloadBytes += downloadDelta
	rollup.sampleCount++
	rollup.observedSeconds += int64(interval.Seconds())
	if resetDetected {
		rollup.resetCount++
	}
	rollup.hasGap = rollup.hasGap || interval > trafficRollupGap
	if snapshot.at.After(rollup.dataAt) {
		rollup.dataAt = snapshot.at
	}
}

type trafficRollupExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func addTrafficHourlyRollup(ctx context.Context, executor trafficRollupExecutor, rollup *trafficHourlyRollup) error {
	_, err := executor.ExecContext(ctx, `INSERT INTO traffic_hourly_rollups
(inbound_id, bucket_start, upload_bytes, download_bytes, sample_count, observed_seconds, reset_count, has_gap, data_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(inbound_id, bucket_start) DO UPDATE SET
  upload_bytes = upload_bytes + excluded.upload_bytes,
  download_bytes = download_bytes + excluded.download_bytes,
  sample_count = sample_count + excluded.sample_count,
  observed_seconds = observed_seconds + excluded.observed_seconds,
  reset_count = reset_count + excluded.reset_count,
  has_gap = CASE WHEN has_gap = 1 OR excluded.has_gap = 1 THEN 1 ELSE 0 END,
  data_at = CASE WHEN data_at IS NULL OR excluded.data_at > data_at THEN excluded.data_at ELSE data_at END,
  updated_at = excluded.updated_at`,
		rollup.inboundID, rollup.bucketStart.Format(time.RFC3339Nano), rollup.uploadBytes, rollup.downloadBytes,
		rollup.sampleCount, rollup.observedSeconds, rollup.resetCount, boolDatabaseValue(rollup.hasGap),
		rollup.dataAt.Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("write traffic rollup %s/%s: %w", rollup.inboundID, rollup.bucketStart.Format(time.RFC3339Nano), err)
	}
	return nil
}

// RebuildTrafficHourlyBucketTx recalculates one inbound/hour from raw data.
// It is used after each newly inserted Agent snapshot, and also handles a
// delayed out-of-order snapshot without double counting either adjacent delta.
func RebuildTrafficHourlyBucketTx(ctx context.Context, transaction *sql.Tx, inboundID string, bucketStart time.Time) error {
	bucketStart = TrafficHourStart(bucketStart)
	bucketEnd := bucketStart.Add(time.Hour)
	rows, err := transaction.QueryContext(ctx, `WITH baseline AS (
  SELECT inbound_id, collected_at, up, down, reset_detected
  FROM traffic_snapshots
  WHERE inbound_id = ? AND collected_at < ?
  ORDER BY collected_at DESC LIMIT 1
), samples AS (
  SELECT inbound_id, collected_at, up, down, reset_detected
  FROM traffic_snapshots
  WHERE inbound_id = ? AND collected_at >= ? AND collected_at < ?
)
SELECT inbound_id, collected_at, up, down, reset_detected FROM (
  SELECT * FROM baseline UNION ALL SELECT * FROM samples
) ORDER BY collected_at`, inboundID, bucketStart.Format(time.RFC3339Nano), inboundID,
		bucketStart.Format(time.RFC3339Nano), bucketEnd.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("read traffic rollup bucket: %w", err)
	}
	snapshots := make([]trafficRollupSnapshot, 0, 64)
	for rows.Next() {
		var snapshot trafficRollupSnapshot
		var collected string
		var reset int
		if err := rows.Scan(&snapshot.inboundID, &collected, &snapshot.up, &snapshot.down, &reset); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan traffic rollup bucket: %w", err)
		}
		parsed, err := time.Parse(time.RFC3339Nano, collected)
		if err != nil {
			_ = rows.Close()
			return fmt.Errorf("parse traffic rollup bucket time: %w", err)
		}
		snapshot.at = parsed.UTC()
		snapshot.reset = reset == 1
		snapshots = append(snapshots, snapshot)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate traffic rollup bucket: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close traffic rollup bucket: %w", err)
	}

	if _, err := transaction.ExecContext(ctx, `DELETE FROM traffic_hourly_rollups WHERE inbound_id = ? AND bucket_start = ?`, inboundID, bucketStart.Format(time.RFC3339Nano)); err != nil {
		return fmt.Errorf("clear traffic rollup bucket: %w", err)
	}
	rollups := make(map[string]*trafficHourlyRollup, 1)
	for index := 1; index < len(snapshots); index++ {
		addTrafficRollupDelta(rollups, snapshots[index-1], snapshots[index])
	}
	for _, rollup := range rollups {
		if err := addTrafficHourlyRollup(ctx, transaction, rollup); err != nil {
			return err
		}
	}
	return nil
}

func boolDatabaseValue(value bool) int {
	if value {
		return 1
	}
	return 0
}
