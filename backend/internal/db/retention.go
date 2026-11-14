package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// RetentionOptions 定义历史数据的保留窗口。字段为 0 表示永久保留、不清理。
//
// 原始流量快照（traffic_snapshots）是数据库里增长最快的表：每个 Inbound 每次
// Agent 同步都会写一行。小时级汇总（traffic_hourly_rollups）已经覆盖了趋势图
// 与统计口径，所以原始快照只需要覆盖页面能查询的时间范围。
type RetentionOptions struct {
	// SnapshotDays 保留原始流量快照的天数。
	SnapshotDays int
	// SyncRunDays 保留 Agent 同步运行记录的天数。
	SyncRunDays int
	// EventDays 保留已处理事件的自然日天数；未处理事件永不清理。
	EventDays int
	// Now 便于测试注入固定时间。
	Now time.Time
}

// RetentionResult 记录本次清理删除的行数。
type RetentionResult struct {
	Snapshots int64
	SyncRuns  int64
	Events    int64
}

// Pruned 表示是否有任何数据被删除。
func (r RetentionResult) Pruned() bool {
	return r.Snapshots > 0 || r.SyncRuns > 0 || r.Events > 0
}

// Total 返回删除的总行数。
func (r RetentionResult) Total() int64 {
	return r.Snapshots + r.SyncRuns + r.Events
}

// pruneBatchSize 控制单次删除的行数。SQLite 单条大 DELETE 会长时间持有写锁，
// 分批删除可以让 Agent 同步在批次之间插入执行。
const pruneBatchSize = 5000

// PruneHistoricalData 按保留窗口删除过期历史数据。
//
// 注意：删除只释放页给后续写入复用，SQLite 文件本身不会变小；数据库大小会在
// 达到稳态后停止增长。需要回收磁盘空间请执行 VACUUM。
func PruneHistoricalData(ctx context.Context, database *sql.DB, options RetentionOptions) (RetentionResult, error) {
	result := RetentionResult{}
	now := options.Now
	if now.IsZero() {
		now = time.Now()
	}

	if options.SnapshotDays > 0 {
		cutoff := now.AddDate(0, 0, -options.SnapshotDays).UTC().Format(time.RFC3339Nano)
		affected, err := deleteBeforeBatch(ctx, database, `DELETE FROM traffic_snapshots WHERE rowid IN (
SELECT rowid FROM traffic_snapshots WHERE collected_at < ? LIMIT ?)`, cutoff)
		if err != nil {
			return result, fmt.Errorf("prune traffic snapshots: %w", err)
		}
		result.Snapshots = affected
	}

	if options.SyncRunDays > 0 {
		cutoff := now.AddDate(0, 0, -options.SyncRunDays).UTC().Format(time.RFC3339Nano)
		affected, err := deleteBeforeBatch(ctx, database, `DELETE FROM sync_runs WHERE rowid IN (
SELECT rowid FROM sync_runs WHERE started_at < ? LIMIT ?)`, cutoff)
		if err != nil {
			return result, fmt.Errorf("prune sync runs: %w", err)
		}
		result.SyncRuns = affected
	}

	if options.EventDays > 0 {
		cutoff := now.AddDate(0, 0, -options.EventDays).UTC().Format(time.RFC3339Nano)
		// 只清理已经处理完的事件；待处理事件是运营的待办，必须保留。
		affected, err := deleteBeforeBatch(ctx, database, `DELETE FROM node_events WHERE rowid IN (
SELECT rowid FROM node_events WHERE created_at < ? AND COALESCE(event_status, 'open') IN ('resolved', 'dismissed') LIMIT ?)`, cutoff)
		if err != nil {
			return result, fmt.Errorf("prune node events: %w", err)
		}
		result.Events = affected
	}

	return result, nil
}

// deleteBeforeBatch 分批执行删除，直到没有匹配行为止。
func deleteBeforeBatch(ctx context.Context, database *sql.DB, statement, cutoff string) (int64, error) {
	total := int64(0)
	for {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		outcome, err := database.ExecContext(ctx, statement, cutoff, pruneBatchSize)
		if err != nil {
			return total, err
		}
		affected, err := outcome.RowsAffected()
		if err != nil {
			return total, err
		}
		total += affected
		if affected < pruneBatchSize {
			return total, nil
		}
	}
}

// Vacuum 重建数据库文件，把 PruneHistoricalData 释放的页真正还给磁盘。
// 执行期间会持有写锁并需要约等于数据库大小的临时空间。
func Vacuum(ctx context.Context, database *sql.DB) error {
	if _, err := database.ExecContext(ctx, `VACUUM`); err != nil {
		return fmt.Errorf("vacuum database: %w", err)
	}
	return nil
}
