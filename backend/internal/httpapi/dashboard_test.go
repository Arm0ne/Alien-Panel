package httpapi

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	centraldb "xpanel-central/backend/internal/db"
)

func TestParseDashboardRangeUsesBusinessMidnight(t *testing.T) {
	now := time.Date(2026, 9, 16, 0, 30, 0, 0, time.UTC)
	r := httptest.NewRequest(http.MethodGet, "/api/dashboard?range=today", nil)

	spec, from, to, err := parseDashboardRange(r, now)
	if err != nil {
		t.Fatalf("parse dashboard range: %v", err)
	}
	if spec.name != "today" {
		t.Fatalf("range name = %q, want today", spec.name)
	}
	wantFrom := time.Date(2026, 9, 15, 16, 0, 0, 0, time.UTC)
	if !from.Equal(wantFrom) {
		t.Fatalf("today from = %s, want %s", from.Format(time.RFC3339), wantFrom.Format(time.RFC3339))
	}
	if !to.Equal(now) {
		t.Fatalf("today to = %s, want %s", to.Format(time.RFC3339), now.Format(time.RFC3339))
	}
}

func TestDashboardHourlyRollupMatchesRawTraffic(t *testing.T) {
	server, database := testServer(t)
	to := time.Date(2026, 9, 23, 6, 53, 0, 0, time.UTC)
	from := to.Add(-7 * 24 * time.Hour)
	createdAt := from.Add(-time.Hour).Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at)
VALUES ('rollup-node', 'rollup-node', 'Rollup node', 'relay', 'online', ?, ?)`, createdAt, createdAt); err != nil {
		t.Fatalf("seed rollup node: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, created_at, updated_at)
VALUES ('rollup-user', 'Rollup user', 'active', ?, ?)`, createdAt, createdAt); err != nil {
		t.Fatalf("seed rollup user: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at)
VALUES ('rollup-inbound', 'rollup-node', '1', 'rollup-user', 'user', ?)`, createdAt); err != nil {
		t.Fatalf("seed rollup inbound: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES ('rollup-mapping', 'rollup-user', 'rollup-inbound', 1, ?)`, createdAt); err != nil {
		t.Fatalf("seed rollup mapping: %v", err)
	}

	transaction, err := database.Begin()
	if err != nil {
		t.Fatalf("begin traffic seed: %v", err)
	}
	start := from.Add(-33 * time.Minute)
	for index, at := 0, start; !at.After(to); index, at = index+1, at.Add(20*time.Minute) {
		up := int64(index * 17)
		down := int64(index * 29)
		reset := 0
		if index == 240 {
			up, down, reset = 3, 5, 1
		}
		if index > 240 {
			up = 3 + int64(index-240)*11
			down = 5 + int64(index-240)*19
		}
		if _, err := transaction.Exec(`INSERT INTO traffic_snapshots
(id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES (?, 'rollup-node', 'rollup-inbound', ?, ?, ?, ?, 'xpanel', ?)`,
			fmt.Sprintf("rollup-sample-%d", index), at.Format(time.RFC3339Nano), up, down, up+down, reset); err != nil {
			_ = transaction.Rollback()
			t.Fatalf("seed traffic sample: %v", err)
		}
	}
	if err := transaction.Commit(); err != nil {
		t.Fatalf("commit traffic seed: %v", err)
	}
	inbounds, err := server.dashboardInbounds()
	if err != nil {
		t.Fatalf("read dashboard inbounds: %v", err)
	}
	spec := dashboardRangeSpec{name: "7d", duration: to.Sub(from), bucket: 24 * time.Hour}
	raw, err := server.dashboardTraffic(from, to, spec, inbounds)
	if err != nil {
		t.Fatalf("read raw dashboard traffic: %v", err)
	}
	if err := centraldb.EnsureTrafficHourlyRollups(context.Background(), database, nil); err != nil {
		t.Fatalf("build traffic rollups: %v", err)
	}
	rolledUp, source, err := server.dashboardTrafficForRange(from, to, spec, inbounds)
	if err != nil {
		t.Fatalf("read rolled up dashboard traffic: %v", err)
	}
	if source != "rollup" {
		t.Fatalf("dashboard traffic source = %q, want rollup", source)
	}
	if raw.trend.Summary != rolledUp.trend.Summary {
		t.Fatalf("rollup summary = %#v, raw = %#v", rolledUp.trend.Summary, raw.trend.Summary)
	}
	if len(raw.trend.Points) != len(rolledUp.trend.Points) {
		t.Fatalf("rollup points = %d, raw = %d", len(rolledUp.trend.Points), len(raw.trend.Points))
	}
	for index := range raw.trend.Points {
		if raw.trend.Points[index] != rolledUp.trend.Points[index] {
			t.Fatalf("rollup point %d = %#v, raw = %#v", index, rolledUp.trend.Points[index], raw.trend.Points[index])
		}
	}
	if *raw.byInbound["rollup-inbound"] != *rolledUp.byInbound["rollup-inbound"] {
		t.Fatalf("rollup inbound traffic = %#v, raw = %#v", rolledUp.byInbound["rollup-inbound"], raw.byInbound["rollup-inbound"])
	}
}

func TestDashboardThirtyDayHourlyRollupScalesToProductionHistory(t *testing.T) {
	if testing.Short() {
		t.Skip("large production-scale dataset is covered by the non-short performance test")
	}
	server, database := testServer(t)
	const inboundCount = 68
	const sampleInterval = 10 * time.Minute
	to := time.Date(2026, 9, 23, 6, 50, 0, 0, time.UTC)
	from := to.Add(-30 * 24 * time.Hour)
	createdAt := from.Add(-time.Hour).Format(time.RFC3339Nano)
	transaction, err := database.Begin()
	if err != nil {
		t.Fatalf("begin production-scale seed: %v", err)
	}
	nodeStatement, err := transaction.Prepare(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES (?, ?, ?, 'relay', 'online', ?, ?)`)
	if err != nil {
		_ = transaction.Rollback()
		t.Fatalf("prepare production-scale node: %v", err)
	}
	userStatement, err := transaction.Prepare(`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES (?, ?, 'active', ?, ?)`)
	if err != nil {
		_ = nodeStatement.Close()
		_ = transaction.Rollback()
		t.Fatalf("prepare production-scale user: %v", err)
	}
	inboundStatement, err := transaction.Prepare(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at) VALUES (?, ?, ?, ?, 'user', ?)`)
	if err != nil {
		_ = nodeStatement.Close()
		_ = userStatement.Close()
		_ = transaction.Rollback()
		t.Fatalf("prepare production-scale inbound: %v", err)
	}
	mappingStatement, err := transaction.Prepare(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES (?, ?, ?, 1, ?)`)
	if err != nil {
		_ = nodeStatement.Close()
		_ = userStatement.Close()
		_ = inboundStatement.Close()
		_ = transaction.Rollback()
		t.Fatalf("prepare production-scale mapping: %v", err)
	}
	snapshotStatement, err := transaction.Prepare(`INSERT INTO traffic_snapshots
(id, node_id, inbound_id, collected_at, up, down, all_time, source)
VALUES (?, ?, ?, ?, ?, ?, ?, 'xpanel')`)
	if err != nil {
		_ = nodeStatement.Close()
		_ = userStatement.Close()
		_ = inboundStatement.Close()
		_ = mappingStatement.Close()
		_ = transaction.Rollback()
		t.Fatalf("prepare production-scale snapshot: %v", err)
	}
	for index := 0; index < inboundCount; index++ {
		nodeID := fmt.Sprintf("scale-node-%02d", index)
		userID := fmt.Sprintf("scale-user-%02d", index)
		inboundID := fmt.Sprintf("scale-inbound-%02d", index)
		if _, err := nodeStatement.Exec(nodeID, nodeID, nodeID, createdAt, createdAt); err != nil {
			t.Fatalf("insert production-scale node: %v", err)
		}
		if _, err := userStatement.Exec(userID, userID, createdAt, createdAt); err != nil {
			t.Fatalf("insert production-scale user: %v", err)
		}
		if _, err := inboundStatement.Exec(inboundID, nodeID, inboundID, userID, createdAt); err != nil {
			t.Fatalf("insert production-scale inbound: %v", err)
		}
		if _, err := mappingStatement.Exec("scale-map-"+inboundID, userID, inboundID, createdAt); err != nil {
			t.Fatalf("insert production-scale mapping: %v", err)
		}
		counter := int64(0)
		for sampleIndex, at := 0, from.Add(-sampleInterval); !at.After(to); sampleIndex, at = sampleIndex+1, at.Add(sampleInterval) {
			counter += int64(1000 + index*13)
			if _, err := snapshotStatement.Exec(fmt.Sprintf("scale-sample-%02d-%05d", index, sampleIndex), nodeID, inboundID,
				at.Format(time.RFC3339Nano), counter/3, counter-counter/3, counter); err != nil {
				t.Fatalf("insert production-scale snapshot: %v", err)
			}
		}
	}
	for _, statement := range []*sql.Stmt{nodeStatement, userStatement, inboundStatement, mappingStatement, snapshotStatement} {
		if err := statement.Close(); err != nil {
			_ = transaction.Rollback()
			t.Fatalf("close production-scale statement: %v", err)
		}
	}
	if err := transaction.Commit(); err != nil {
		t.Fatalf("commit production-scale seed: %v", err)
	}
	inbounds, err := server.dashboardInbounds()
	if err != nil {
		t.Fatalf("read production-scale inbounds: %v", err)
	}
	spec := dashboardRangeSpec{name: "30d", duration: to.Sub(from), bucket: 24 * time.Hour}
	raw, err := server.dashboardTraffic(from, to, spec, inbounds)
	if err != nil {
		t.Fatalf("read production-scale raw traffic: %v", err)
	}
	if raw.trend.Summary.SampleCount < 250_000 {
		t.Fatalf("seeded only %d traffic samples, want minute-level production history", raw.trend.Summary.SampleCount)
	}
	if err := centraldb.EnsureTrafficHourlyRollups(context.Background(), database, nil); err != nil {
		t.Fatalf("build production-scale rollups: %v", err)
	}
	started := time.Now()
	rolledUp, source, err := server.dashboardTrafficForRange(from, to, spec, inbounds)
	queryDuration := time.Since(started)
	if err != nil {
		t.Fatalf("read production-scale rolled up traffic: %v", err)
	}
	if source != "rollup" {
		t.Fatalf("traffic source = %q, want rollup", source)
	}
	if queryDuration >= time.Second {
		t.Fatalf("30-day rollup query took %s for %d snapshots", queryDuration, raw.trend.Summary.SampleCount)
	}
	t.Logf("30-day rollup query: %s for %d minute-level snapshots across %d Inbounds", queryDuration, raw.trend.Summary.SampleCount, inboundCount)
	if rolledUp.trend.Summary != raw.trend.Summary {
		t.Fatalf("production-scale rollup summary = %#v, raw = %#v", rolledUp.trend.Summary, raw.trend.Summary)
	}
}

func TestParseDashboardCustomDateUsesBusinessMidnight(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/dashboard?range=custom&from=2026-09-16&to=2026-09-17", nil)
	_, from, to, err := parseDashboardRange(r, time.Date(2026, 9, 16, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("parse custom dashboard range: %v", err)
	}
	wantFrom := time.Date(2026, 9, 15, 16, 0, 0, 0, time.UTC)
	wantTo := time.Date(2026, 9, 17, 16, 0, 0, 0, time.UTC)
	if !from.Equal(wantFrom) || !to.Equal(wantTo) {
		t.Fatalf("custom range = %s to %s, want %s to %s", from.Format(time.RFC3339), to.Format(time.RFC3339), wantFrom.Format(time.RFC3339), wantTo.Format(time.RFC3339))
	}
}

func TestDashboardTrafficUsesBusinessInboundScope(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-3 * time.Hour)
	baseline := from.Add(-time.Hour)
	first := from.Add(time.Hour)
	second := from.Add(2 * time.Hour)
	onlineAt := second.Add(30 * time.Minute)
	nowText := now.Format(time.RFC3339Nano)

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('dashboard-relay-a', 'dashboard-relay-a', '线路机 A', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('dashboard-relay-b', 'dashboard-relay-b', '线路机 B', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('dashboard-landing', 'dashboard-landing', '落地机', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('dashboard-user-a', '用户 A', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('dashboard-user-b', '用户 B', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('dashboard-user-first', '首个快照用户', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at) VALUES ('dashboard-inbound-a', 'dashboard-relay-a', '1', 'dashboard-user-a', 'user', 'business-a', ?)`, []any{nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at) VALUES ('dashboard-inbound-b', 'dashboard-relay-b', '2', 'dashboard-user-b', 'user', 'business-b', ?)`, []any{nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at) VALUES ('dashboard-inbound-first', 'dashboard-relay-b', '3', 'dashboard-user-first', 'user', 'first-snapshot', ?)`, []any{nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at) VALUES ('dashboard-landing-inbound', 'dashboard-landing', '4', 'infrastructure', 'landing-traffic', ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('dashboard-map-a', 'dashboard-user-a', 'dashboard-inbound-a', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('dashboard-map-b', 'dashboard-user-b', 'dashboard-inbound-b', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('dashboard-map-first', 'dashboard-user-first', 'dashboard-inbound-first', 1, ?)`, []any{nowText}},
		// Two Clients with the same Email still form one Inbound business user and
		// must not make the Inbound traffic count twice.
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable) VALUES ('dashboard-client-a-phone', 'dashboard-relay-a', 'dashboard-inbound-a', 'phone', 'same@example.com', 1)`, nil},
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable) VALUES ('dashboard-client-a-laptop', 'dashboard-relay-a', 'dashboard-inbound-a', 'laptop', 'same@example.com', 1)`, nil},
		{`UPDATE clients SET last_online = ?, last_seen_at = ? WHERE inbound_id = 'dashboard-inbound-a'`, []any{onlineAt.Format(time.RFC3339Nano), nowText}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-a-base', 'dashboard-relay-a', 'dashboard-inbound-a', ?, 10, 20, 30, 'xpanel')`, []any{baseline.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-a-first', 'dashboard-relay-a', 'dashboard-inbound-a', ?, 30, 50, 80, 'xpanel')`, []any{first.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-a-second', 'dashboard-relay-a', 'dashboard-inbound-a', ?, 50, 70, 120, 'xpanel')`, []any{second.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-b-base', 'dashboard-relay-b', 'dashboard-inbound-b', ?, 3, 2, 5, 'xpanel')`, []any{baseline.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-b-first', 'dashboard-relay-b', 'dashboard-inbound-b', ?, 8, 12, 20, 'xpanel')`, []any{first.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected) VALUES ('dashboard-b-reset', 'dashboard-relay-b', 'dashboard-inbound-b', ?, 2, 3, 5, 'xpanel', 1)`, []any{second.Format(time.RFC3339Nano)}},
		// The first known point is only a baseline, so this large counter is not usage.
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-first-only', 'dashboard-relay-b', 'dashboard-inbound-first', ?, 500, 500, 1000, 'xpanel')`, []any{first.Format(time.RFC3339Nano)}},
		// Landing traffic is infrastructure traffic and is outside the business scope.
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-landing-base', 'dashboard-landing', 'dashboard-landing-inbound', ?, 0, 0, 0, 'xpanel')`, []any{baseline.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-landing-usage', 'dashboard-landing', 'dashboard-landing-inbound', ?, 5000, 5000, 10000, 'xpanel')`, []any{first.Format(time.RFC3339Nano)}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed dashboard data: %v", err)
		}
	}

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	query := url.Values{"range": {"custom"}, "from": {from.Format(time.RFC3339Nano)}, "to": {now.Format(time.RFC3339Nano)}}
	response := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard?"+query.Encode(), token, nil)
	if response["code"] != successCode {
		t.Fatalf("dashboard response = %#v", response)
	}

	data := response["data"].(map[string]any)
	traffic := data["traffic"].(map[string]any)
	if traffic["totalBytes"] != float64(110) || traffic["uploadBytes"] != float64(47) || traffic["downloadBytes"] != float64(63) {
		t.Fatalf("traffic = %#v, want upload=47 download=63 total=110", traffic)
	}
	trend := data["trafficTrend"].(map[string]any)
	if trend["summary"].(map[string]any)["totalBytes"] != float64(110) {
		t.Fatalf("trend summary = %#v, want 110", trend["summary"])
	}
	points := trend["points"].([]any)
	if len(points) != 2 || !points[1].(map[string]any)["resetDetected"].(bool) {
		t.Fatalf("trend points = %#v, want reset in the second bucket", points)
	}

	nodes := data["nodeTrafficRanking"].([]any)
	if len(nodes) != 2 || nodes[0].(map[string]any)["nodeId"] != "dashboard-relay-a" || nodes[0].(map[string]any)["totalBytes"] != float64(90) {
		t.Fatalf("node traffic ranking = %#v", nodes)
	}
	users := data["userTrafficRanking"].([]any)
	if len(users) != 3 || users[0].(map[string]any)["inboundId"] != "dashboard-inbound-a" || users[0].(map[string]any)["clientCount"] != float64(2) || users[0].(map[string]any)["totalBytes"] != float64(90) {
		t.Fatalf("user traffic ranking = %#v", users)
	}
	if users[0].(map[string]any)["lastActivityAt"] != onlineAt.Format(time.RFC3339Nano) {
		t.Fatalf("dashboard last activity = %v, want latest Client online time %s", users[0].(map[string]any)["lastActivityAt"], onlineAt.Format(time.RFC3339Nano))
	}
}

func TestDashboardTrafficIgnoresSnapshotsWithoutLiveUserMetadata(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-3 * time.Hour)
	baseline := from.Add(-time.Hour)
	usage := from.Add(time.Hour)
	nowText := now.Format(time.RFC3339Nano)

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('dashboard-live-node', 'dashboard-live-node', '在线节点', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('dashboard-live-user', '在线用户', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, deleted_at, created_at, updated_at) VALUES ('dashboard-deleted-user', '已删除用户', 'disabled', ?, ?, ?)`, []any{nowText, nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at) VALUES ('dashboard-live-inbound', 'dashboard-live-node', '1', 'dashboard-live-user', 'user', ?)`, []any{nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at) VALUES ('dashboard-deleted-inbound', 'dashboard-live-node', '2', 'dashboard-deleted-user', 'user', ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('dashboard-live-mapping', 'dashboard-live-user', 'dashboard-live-inbound', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('dashboard-deleted-mapping', 'dashboard-deleted-user', 'dashboard-deleted-inbound', 1, ?)`, []any{nowText}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-live-base', 'dashboard-live-node', 'dashboard-live-inbound', ?, 10, 20, 30, 'xpanel')`, []any{baseline.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-live-usage', 'dashboard-live-node', 'dashboard-live-inbound', ?, 30, 50, 80, 'xpanel')`, []any{usage.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-deleted-base', 'dashboard-live-node', 'dashboard-deleted-inbound', ?, 1000, 2000, 3000, 'xpanel')`, []any{baseline.Format(time.RFC3339Nano)}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('dashboard-deleted-usage', 'dashboard-live-node', 'dashboard-deleted-inbound', ?, 9000, 12000, 21000, 'xpanel')`, []any{usage.Format(time.RFC3339Nano)}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed deleted-user dashboard data: %v", err)
		}
	}

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	query := url.Values{"range": {"custom"}, "from": {from.Format(time.RFC3339Nano)}, "to": {now.Format(time.RFC3339Nano)}}
	response := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard?"+query.Encode(), token, nil)
	if response["code"] != successCode {
		t.Fatalf("dashboard response = %#v", response)
	}

	data := response["data"].(map[string]any)
	traffic := data["traffic"].(map[string]any)
	if traffic["totalBytes"] != float64(50) {
		t.Fatalf("traffic = %#v, want deleted-user snapshot excluded", traffic)
	}
	users := data["userTrafficRanking"].([]any)
	if len(users) != 1 || users[0].(map[string]any)["userId"] != "dashboard-live-user" {
		t.Fatalf("user traffic ranking = %#v, want only live user", users)
	}
}
