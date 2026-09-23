package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestEnsureTrafficHourlyRollupsHandlesResetAndIsIdempotent(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if err := Migrate(database); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	now := time.Date(2026, 9, 23, 6, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('node', 'node', 'node', 'relay', ?, ?)`, nowText, nowText); err != nil {
		t.Fatalf("seed node: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('user', 'user', 'active', ?, ?)`, nowText, nowText); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at) VALUES ('inbound', 'node', '1', 'user', 'user', ?)`, nowText); err != nil {
		t.Fatalf("seed inbound: %v", err)
	}
	samples := []struct {
		at       time.Time
		up, down int64
		reset    int
	}{
		{now.Add(-2*time.Hour - 10*time.Minute), 100, 200, 0},
		{now.Add(-90 * time.Minute), 130, 250, 0},
		{now.Add(-70 * time.Minute), 160, 280, 0},
		{now.Add(-50 * time.Minute), 5, 7, 1},
		{now.Add(-10 * time.Minute), 25, 47, 0},
	}
	for index, sample := range samples {
		if _, err := database.Exec(`INSERT INTO traffic_snapshots
(id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES (?, 'node', 'inbound', ?, ?, ?, ?, 'xpanel', ?)`, fmt.Sprintf("sample-%d", index), sample.at.Format(time.RFC3339Nano), sample.up, sample.down, sample.up+sample.down, sample.reset); err != nil {
			t.Fatalf("insert traffic sample %d: %v", index, err)
		}
	}
	if err := EnsureTrafficHourlyRollups(context.Background(), database, nil); err != nil {
		t.Fatalf("build traffic rollups: %v", err)
	}
	assertTrafficRollupTotals(t, database, 85, 127, 4, 1)
	if err := EnsureTrafficHourlyRollups(context.Background(), database, nil); err != nil {
		t.Fatalf("repeat traffic rollup build: %v", err)
	}
	assertTrafficRollupTotals(t, database, 85, 127, 4, 1)
}

func assertTrafficRollupTotals(t *testing.T, database interface {
	QueryRow(string, ...any) *sql.Row
}, wantUpload, wantDownload, wantSamples, wantResets int64) {
	t.Helper()
	var upload, download, samples, resets int64
	if err := database.QueryRow(`SELECT COALESCE(SUM(upload_bytes), 0), COALESCE(SUM(download_bytes), 0),
COALESCE(SUM(sample_count), 0), COALESCE(SUM(reset_count), 0) FROM traffic_hourly_rollups`).Scan(&upload, &download, &samples, &resets); err != nil {
		t.Fatalf("read traffic rollup totals: %v", err)
	}
	if upload != wantUpload || download != wantDownload || samples != wantSamples || resets != wantResets {
		t.Fatalf("traffic rollups = upload %d download %d samples %d resets %d", upload, download, samples, resets)
	}
}
