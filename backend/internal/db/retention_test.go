package db

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestPruneHistoricalDataHonoursRetentionWindows 校验保留策略只删除窗口外的历史行，
// 并且不触碰小时级汇总与未处理事件。
func TestPruneHistoricalDataHonoursRetentionWindows(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := Migrate(database); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	stamp := func(daysAgo int) string {
		return now.AddDate(0, 0, -daysAgo).UTC().Format(time.RFC3339Nano)
	}

	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at)
VALUES ('n1', 'n1', '测试节点', 'relay', ?, ?)`, stamp(400), stamp(400)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at, last_seen_at)
VALUES ('i1', 'n1', '1', 'user', 'test', ?, ?)`, stamp(400), stamp(0)); err != nil {
		t.Fatal(err)
	}

	// 快照：两条超出 90 天，两条在窗口内
	for index, daysAgo := range []int{120, 100, 30, 1} {
		if _, err := database.Exec(`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES (?, 'n1', 'i1', ?, 1, 2, 3, 'xpanel', 0)`, fmt.Sprintf("s%d", index), stamp(daysAgo)); err != nil {
			t.Fatal(err)
		}
	}

	// 同步记录：一条超出 30 天
	for index, daysAgo := range []int{45, 5} {
		if _, err := database.Exec(`INSERT INTO sync_runs (id, node_id, sync_id, started_at, status)
VALUES (?, 'n1', ?, ?, 'success')`, fmt.Sprintf("r%d", index), fmt.Sprintf("sync-%d", index), stamp(daysAgo)); err != nil {
			t.Fatal(err)
		}
	}

	// 事件：一条 200 天前已处理（应删）、一条 200 天前未处理（必须保留）、一条近期已处理（保留）
	if _, err := database.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at, event_status)
VALUES ('e-old-resolved', 'n1', 'sync_failed', 'warning', '旧已处理', ?, 'resolved')`, stamp(200)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at, event_status)
VALUES ('e-old-open', 'n1', 'sync_failed', 'warning', '旧未处理', ?, 'open')`, stamp(200)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at, event_status)
VALUES ('e-new-resolved', 'n1', 'sync_failed', 'warning', '新已处理', ?, 'resolved')`, stamp(5)); err != nil {
		t.Fatal(err)
	}

	// 小时汇总不应被保留策略影响
	if _, err := database.Exec(`INSERT INTO traffic_hourly_rollups (inbound_id, bucket_start, upload_bytes, download_bytes, sample_count, observed_seconds, reset_count, has_gap, updated_at)
VALUES ('i1', ?, 10, 20, 60, 3600, 0, 0, ?)`, stamp(400), stamp(400)); err != nil {
		t.Fatal(err)
	}

	result, err := PruneHistoricalData(context.Background(), database, RetentionOptions{
		SnapshotDays: 90,
		SyncRunDays:  30,
		EventDays:    180,
		Now:          now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Snapshots != 2 || result.SyncRuns != 1 || result.Events != 1 {
		t.Fatalf("pruned snapshots=%d syncRuns=%d events=%d, want 2/1/1", result.Snapshots, result.SyncRuns, result.Events)
	}

	count := func(query string) int {
		t.Helper()
		var value int
		if err := database.QueryRow(query).Scan(&value); err != nil {
			t.Fatalf("count %q: %v", query, err)
		}
		return value
	}
	if got := count(`SELECT COUNT(*) FROM traffic_snapshots`); got != 2 {
		t.Fatalf("remaining snapshots = %d, want 2", got)
	}
	if got := count(`SELECT COUNT(*) FROM traffic_snapshots WHERE collected_at < '2026-03-03'`); got != 0 {
		t.Fatalf("old snapshots survived: %d", got)
	}
	if got := count(`SELECT COUNT(*) FROM sync_runs`); got != 1 {
		t.Fatalf("remaining sync runs = %d, want 1", got)
	}
	if got := count(`SELECT COUNT(*) FROM node_events`); got != 2 {
		t.Fatalf("remaining events = %d, want 2", got)
	}
	if got := count(`SELECT COUNT(*) FROM node_events WHERE id = 'e-old-open'`); got != 1 {
		t.Fatalf("unresolved event was pruned, it must be kept")
	}
	if got := count(`SELECT COUNT(*) FROM traffic_hourly_rollups`); got != 1 {
		t.Fatalf("hourly rollups were touched: %d", got)
	}
}

// TestPruneHistoricalDataDisabledWindowsKeepEverything 确认 0 表示永久保留。
func TestPruneHistoricalDataDisabledWindowsKeepEverything(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := Migrate(database); err != nil {
		t.Fatal(err)
	}

	old := time.Now().AddDate(-2, 0, 0).UTC().Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at)
VALUES ('n1', 'n1', '节点', 'relay', ?, ?)`, old, old); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at, last_seen_at)
VALUES ('i1', 'n1', '1', 'user', 'test', ?, ?)`, old, old); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES ('s1', 'n1', 'i1', ?, 1, 2, 3, 'xpanel', 0)`, old); err != nil {
		t.Fatal(err)
	}

	result, err := PruneHistoricalData(context.Background(), database, RetentionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Pruned() {
		t.Fatalf("expected no pruning when all windows are 0, got %+v", result)
	}
	var remaining int
	if err := database.QueryRow(`SELECT COUNT(*) FROM traffic_snapshots`).Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if remaining != 1 {
		t.Fatalf("snapshot removed with retention disabled: %d", remaining)
	}
}
