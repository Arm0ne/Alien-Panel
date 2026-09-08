package httpapi

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestDashboardTrafficUsesBusinessInboundScope(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Truncate(time.Second)
	from := now.Add(-3 * time.Hour)
	baseline := from.Add(-time.Hour)
	first := from.Add(time.Hour)
	second := from.Add(2 * time.Hour)
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
}
