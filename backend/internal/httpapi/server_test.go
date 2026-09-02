package httpapi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"xpanel-central/backend/internal/config"
	"xpanel-central/backend/internal/db"
)

func testServer(t *testing.T) (*Server, *sql.DB) {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		database.Close()
		t.Fatalf("migrate database: %v", err)
	}
	server, err := NewServer(config.Config{
		AdminUsername:          "admin",
		AdminPassword:          "test-password",
		AgentRegistrationToken: "bootstrap-secret",
		SessionTTL:             time.Hour,
		CorsOrigins:            []string{"http://localhost:9527"},
	}, database, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		database.Close()
		t.Fatalf("new server: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return server, database
}

func TestLoginAndDashboard(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	loginResponse := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin",
		"password": "test-password",
	})
	if loginResponse["code"] != successCode {
		t.Fatalf("login code = %v", loginResponse["code"])
	}
	loginData := loginResponse["data"].(map[string]any)
	token := loginData["token"].(string)
	if token == "" || token == "test-password" {
		t.Fatalf("unexpected token: %q", token)
	}

	userInfo := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/auth/me", token, nil)
	if userInfo["code"] != successCode || userInfo["data"].(map[string]any)["userName"] != "admin" {
		t.Fatalf("unexpected user info: %#v", userInfo)
	}

	dashboard := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard", token, nil)
	if dashboard["code"] != successCode {
		t.Fatalf("dashboard code = %v", dashboard["code"])
	}
	if dashboard["data"].(map[string]any)["nodes"].(map[string]any)["total"] != float64(0) {
		t.Fatalf("expected empty dashboard: %#v", dashboard)
	}

	var passwordHash string
	if err := database.QueryRow(`SELECT password_hash FROM admin_users WHERE username = 'admin'`).Scan(&passwordHash); err != nil {
		t.Fatalf("read password hash: %v", err)
	}
	if passwordHash == "test-password" || passwordHash == "" {
		t.Fatalf("password was not hashed: %q", passwordHash)
	}
}

func TestProtectedEndpointsRejectMissingToken(t *testing.T) {
	server, _ := testServer(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	server.Handler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestUserDetailAndBusinessFieldUpdate(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	unauthorized := httptest.NewRecorder()
	server.Handler().ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/api/users/user-detail", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized detail status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	expiryText := now.AddDate(0, 1, 0).Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('detail-relay', 'detail-relay', '东京线路机', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('detail-landing', 'detail-landing', '东京落地机', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, notes, expiry_time, created_at, updated_at) VALUES ('user-detail', '原始名称', 'active', 30, 'CNY', '原始备注', ?, ?, ?)`, []any{expiryText, nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, remark, protocol, port, enable, client_count, up, down, all_time, first_seen_at, last_seen_at) VALUES ('inbound-detail', 'detail-relay', '37', 'user-detail', 'user', 'reality-user-37', 'X-Panel 原始备注', 'vless', 443, 1, 2, 100, 200, 300, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('mapping-detail', 'user-detail', 'inbound-detail', 1, ?)`, []any{nowText}},
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, up, down, all_time) VALUES ('client-detail-1', 'detail-relay', 'inbound-detail', 'phone-uuid', 'phone@example.com', 1, 10, 20, 30)`, nil},
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, up, down, all_time) VALUES ('client-detail-2', 'detail-relay', 'inbound-detail', 'laptop-uuid', 'laptop@example.com', 1, 30, 40, 70)`, nil},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected) VALUES ('traffic-detail-1', 'detail-relay', 'inbound-detail', ?, 100, 200, 300, 'xpanel', 0)`, []any{nowText}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected) VALUES ('traffic-detail-2', 'detail-relay', 'inbound-detail', ?, 10, 20, 30, 'xpanel', 1)`, []any{now.Add(time.Minute).Format(time.RFC3339Nano)}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, landing_inbound_tag, enabled, created_at, updated_at) VALUES ('route-detail', '东京优化线路', 'detail-relay', 'detail-landing', 'ss-tokyo', 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('route-mapping-detail', 'user-detail', 'route-detail', 1, ?)`, []any{nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed user detail data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/user-detail", token, nil)
	if detail["code"] != successCode {
		t.Fatalf("user detail response = %#v", detail)
	}
	detailData := detail["data"].(map[string]any)
	traffic := detailData["traffic"].([]any)
	if detailData["displayName"] != "原始名称" || detailData["inbound"].(map[string]any)["tag"] != "reality-user-37" || len(detailData["clients"].([]any)) != 2 || len(detailData["routes"].([]any)) != 1 || len(traffic) != 2 || traffic[0].(map[string]any)["resetDetected"] != true {
		t.Fatalf("unexpected user detail = %#v", detailData)
	}

	negative := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/users/user-detail", token, map[string]any{"monthlyFee": -1})
	if negative["code"] != validationCode {
		t.Fatalf("negative fee response = %#v", negative)
	}

	updated := doJSONWithRequestID(t, ts.Client(), http.MethodPatch, ts.URL+"/api/users/user-detail", token, "user-detail-update-test", map[string]any{
		"displayName": "中央业务名称", "monthlyFee": 68.5, "currency": "CNY", "notes": "只在中央面板维护",
	})
	if updated["code"] != successCode || updated["data"].(map[string]any)["displayName"] != "中央业务名称" {
		t.Fatalf("update response = %#v", updated)
	}

	var displayName, currency string
	var monthlyFee float64
	var notes sql.NullString
	if err := database.QueryRow(`SELECT display_name, monthly_fee, currency, notes FROM users WHERE id = 'user-detail'`).Scan(&displayName, &monthlyFee, &currency, &notes); err != nil {
		t.Fatalf("read updated business fields: %v", err)
	}
	if displayName != "中央业务名称" || monthlyFee != 68.5 || currency != "CNY" || !notes.Valid || notes.String != "只在中央面板维护" {
		t.Fatalf("updated business fields = name=%q fee=%v currency=%q notes=%#v", displayName, monthlyFee, currency, notes)
	}

	var requestID, beforeJSON, afterJSON string
	if err := database.QueryRow(`SELECT request_id, before_json, after_json FROM audit_logs WHERE action = 'user.update' AND resource_id = 'user-detail'`).Scan(&requestID, &beforeJSON, &afterJSON); err != nil {
		t.Fatalf("read audit log: %v", err)
	}
	if requestID != "user-detail-update-test" || !strings.Contains(beforeJSON, "原始名称") || !strings.Contains(afterJSON, "中央业务名称") {
		t.Fatalf("unexpected audit record request_id=%q before=%s after=%s", requestID, beforeJSON, afterJSON)
	}

	// A later relay snapshot may refresh X-Panel-owned expiry/status, but must
	// never overwrite central-operated customer metadata.
	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("begin relay sync transaction: %v", err)
	}
	if err := server.ensureRelayInboundUser(tx, "relay", "inbound-detail", "37", agentInboundPayload{Enable: true, Remark: "X-Panel changed this"}, expiryText, now.Add(time.Minute)); err != nil {
		_ = tx.Rollback()
		t.Fatalf("refresh relay user: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit relay sync transaction: %v", err)
	}
	if err := database.QueryRow(`SELECT display_name, monthly_fee, currency, notes FROM users WHERE id = 'user-detail'`).Scan(&displayName, &monthlyFee, &currency, &notes); err != nil {
		t.Fatalf("read business fields after relay sync: %v", err)
	}
	if displayName != "中央业务名称" || monthlyFee != 68.5 || currency != "CNY" || !notes.Valid || notes.String != "只在中央面板维护" {
		t.Fatalf("relay sync overwrote business fields: name=%q fee=%v currency=%q notes=%#v", displayName, monthlyFee, currency, notes)
	}
}

func TestAgentRegistrationHeartbeatAndIdempotentSync(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	unauthorizedRegistration := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/register", "", map[string]any{
		"node_key": "relay-agent-unauthorized", "node_name": "unauthorized",
	})
	if unauthorizedRegistration["code"] != unauthorizedCode {
		t.Fatalf("registration without bootstrap token = %#v", unauthorizedRegistration)
	}

	registrationRequest, err := http.NewRequest(http.MethodPost, ts.URL+"/api/agent/v1/register", strings.NewReader(`{"node_key":"relay-agent-1","node_name":"线路机 Agent 1","node_type":"relay","hostname":"relay-1.example","panel_base_path":"/","agent_version":"0.1.0"}`))
	if err != nil {
		t.Fatal(err)
	}
	registrationRequest.Header.Set("Content-Type", "application/json")
	registrationRequest.Header.Set("X-Agent-Registration-Token", "bootstrap-secret")
	registrationResponse, err := ts.Client().Do(registrationRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer registrationResponse.Body.Close()
	var registration map[string]any
	if err := json.NewDecoder(registrationResponse.Body).Decode(&registration); err != nil {
		t.Fatal(err)
	}
	/*
		Registration is intentionally tested with a separate header because the
		bootstrap token is not a node credential and must never become a Bearer
		token.
	*/
	if registration["code"] != successCode {
		t.Fatalf("registration response = %#v", registration)
	}
	registrationData := registration["data"].(map[string]any)
	nodeToken := registrationData["token"].(string)
	if nodeToken == "" {
		t.Fatal("registration returned empty node token")
	}

	observedAt := "2026-09-02T12:00:00Z"
	heartbeat := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", nodeToken, map[string]any{
		"node_key": "relay-agent-1", "observed_at": observedAt,
		"status": map[string]any{"xray_running": true, "xray_version": "26.6.27", "xpanel_version": "2.4.0"},
	})
	if heartbeat["code"] != successCode {
		t.Fatalf("heartbeat response = %#v", heartbeat)
	}

	payload := map[string]any{
		"node_key": "relay-agent-1", "sync_id": "relay-agent-1-sync-001", "observed_at": observedAt,
		"status": map[string]any{"xray_running": true, "xray_version": "26.6.27", "xpanel_version": "2.4.0"},
		"inbounds": []any{map[string]any{
			"remote_id": 15, "tag": "user-15", "remark": "Customer A", "protocol": "vless", "port": 443,
			"enable": true, "expiry_time": 1792022400, "up": 100, "down": 200, "all_time": 300,
			"config_hash": "hash-1", "clients": []any{
				map[string]any{"remote_id": "client-a", "email": "phone", "enable": true, "all_time": 300},
			},
		}},
	}
	syncResponse := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, payload)
	if syncResponse["code"] != successCode || syncResponse["data"].(map[string]any)["status"] != "success" {
		t.Fatalf("sync response = %#v", syncResponse)
	}
	if _, err := database.Exec(`UPDATE users SET display_name = '运营名称', monthly_fee = 99, notes = 'central-only'`); err != nil {
		t.Fatalf("seed central business fields: %v", err)
	}
	duplicate := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, payload)
	if duplicate["code"] != successCode || duplicate["data"].(map[string]any)["idempotent"] != true {
		t.Fatalf("duplicate sync response = %#v", duplicate)
	}
	payload["sync_id"] = "relay-agent-1-sync-002"
	payload["observed_at"] = "2026-09-02T12:01:00Z"
	payload["inbounds"].([]any)[0].(map[string]any)["up"] = 50
	payload["inbounds"].([]any)[0].(map[string]any)["down"] = 50
	payload["inbounds"].([]any)[0].(map[string]any)["all_time"] = 100
	payload["inbounds"].([]any)[0].(map[string]any)["remark"] = "X-Panel changed remark"
	resetSync := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, payload)
	if resetSync["code"] != successCode {
		t.Fatalf("reset sync response = %#v", resetSync)
	}

	var nodes, inbounds, users, userInboundMappings, clients, snapshots, resetSnapshots, resetEvents, syncRuns int
	for query, target := range map[string]*int{
		"SELECT COUNT(*) FROM nodes WHERE node_key = 'relay-agent-1'":             &nodes,
		"SELECT COUNT(*) FROM inbounds WHERE remote_inbound_id = '15'":            &inbounds,
		"SELECT COUNT(*) FROM users":                                              &users,
		"SELECT COUNT(*) FROM user_inbounds":                                      &userInboundMappings,
		"SELECT COUNT(*) FROM clients WHERE remote_client_id = 'client-a'":        &clients,
		"SELECT COUNT(*) FROM traffic_snapshots WHERE all_time = 300":             &snapshots,
		"SELECT COUNT(*) FROM traffic_snapshots WHERE reset_detected = 1":         &resetSnapshots,
		"SELECT COUNT(*) FROM node_events WHERE event_type = 'traffic_reset'":     &resetEvents,
		"SELECT COUNT(*) FROM sync_runs WHERE sync_id = 'relay-agent-1-sync-001'": &syncRuns,
	} {
		if err := database.QueryRow(query).Scan(target); err != nil {
			t.Fatalf("query %q: %v", query, err)
		}
	}
	if nodes != 1 || inbounds != 1 || users != 1 || userInboundMappings != 1 || clients != 1 || snapshots != 1 || resetSnapshots != 1 || resetEvents != 1 || syncRuns != 1 {
		t.Fatalf("stored rows nodes=%d inbounds=%d users=%d userInboundMappings=%d clients=%d snapshots=%d resetSnapshots=%d resetEvents=%d syncRuns=%d", nodes, inbounds, users, userInboundMappings, clients, snapshots, resetSnapshots, resetEvents, syncRuns)
	}
	var displayName, notes string
	var monthlyFee float64
	var userExpiry sql.NullString
	if err := database.QueryRow(`SELECT display_name, monthly_fee, notes, expiry_time FROM users`).Scan(&displayName, &monthlyFee, &notes, &userExpiry); err != nil {
		t.Fatalf("read business user: %v", err)
	}
	if displayName != "运营名称" || monthlyFee != 99 || notes != "central-only" || !userExpiry.Valid {
		t.Fatalf("business fields were not preserved or expiry missing: name=%q fee=%v notes=%q expiry=%v", displayName, monthlyFee, notes, userExpiry)
	}
}

func TestInboundIsArchivedAfterThreeConsecutiveMissingSyncs(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('archive-node', 'archive-node', 'Archive Node', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatalf("create node: %v", err)
	}
	const nodeToken = "archive-node-token"
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('archive-credential', 'archive-node', ?, ?, ?)`, hashToken(nodeToken), nowText, nowText); err != nil {
		t.Fatalf("create node credential: %v", err)
	}

	initial := map[string]any{
		"node_key": "archive-node", "sync_id": "archive-initial", "observed_at": nowText,
		"status": map[string]any{"xray_running": true},
		"inbounds": []any{map[string]any{
			"remote_id": 88, "tag": "archive-me", "protocol": "vless", "port": 443,
			"enable": true, "all_time": 100, "config_hash": "archive-hash", "clients": []any{},
		}},
	}
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, initial); result["code"] != successCode {
		t.Fatalf("initial sync response = %#v", result)
	}

	for count := 1; count <= missingInboundArchiveAfter; count++ {
		observedAt := now.Add(time.Duration(count) * time.Minute).Format(time.RFC3339Nano)
		missing := map[string]any{
			"node_key": "archive-node", "sync_id": "archive-missing-" + strconv.Itoa(count), "observed_at": observedAt,
			"status": map[string]any{"xray_running": true}, "inbounds": []any{},
		}
		if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, missing); result["code"] != successCode {
			t.Fatalf("missing sync %d response = %#v", count, result)
		}

		var missingCount int
		var deletedAt sql.NullString
		if err := database.QueryRow(`SELECT missing_sync_count, deleted_at FROM inbounds WHERE node_id = 'archive-node' AND remote_inbound_id = '88'`).Scan(&missingCount, &deletedAt); err != nil {
			t.Fatalf("read inbound after missing sync %d: %v", count, err)
		}
		if missingCount != count {
			t.Fatalf("missing count after sync %d = %d", count, missingCount)
		}
		if count < missingInboundArchiveAfter && deletedAt.Valid {
			t.Fatalf("inbound archived too early after %d missing syncs", count)
		}
		if count == missingInboundArchiveAfter && !deletedAt.Valid {
			t.Fatalf("inbound was not archived after %d missing syncs", count)
		}
	}

	var archivedEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE node_id = 'archive-node' AND event_type = 'inbound_archived'`).Scan(&archivedEvents); err != nil {
		t.Fatalf("count archive events: %v", err)
	}
	if archivedEvents != 1 {
		t.Fatalf("archive event count = %d, want 1", archivedEvents)
	}
}

func TestDashboardUsesTrafficDeltas(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC()
	created := now.Format(time.RFC3339Nano)
	first := now.Add(-2 * time.Hour).Format(time.RFC3339Nano)
	second := now.Add(-time.Hour).Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('traffic-node', 'traffic-node', 'Traffic Node', 'relay', 'online', ?, ?)`, []any{created, created}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, tag, first_seen_at, last_seen_at) VALUES ('traffic-inbound', 'traffic-node', '15', 'user-15', ?, ?)`, []any{first, second}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('traffic-s1', 'traffic-node', 'traffic-inbound', ?, 40, 60, 100, 'xpanel')`, []any{first}},
		{`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source) VALUES ('traffic-s2', 'traffic-node', 'traffic-inbound', ?, 100, 150, 250, 'xpanel')`, []any{second}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed traffic data: %v", err)
		}
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	dashboard := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard", token, nil)
	if dashboard["code"] != successCode {
		t.Fatalf("dashboard response = %#v", dashboard)
	}
	traffic := dashboard["data"].(map[string]any)["traffic"].(map[string]any)
	if traffic["todayBytes"] != float64(150) || traffic["monthBytes"] != float64(150) {
		t.Fatalf("traffic = %#v, want delta 150", traffic)
	}
}

func TestOperationalStatusRefresh(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	oldText := now.Add(-10 * time.Minute).Format(time.RFC3339Nano)
	expiredText := now.Add(-time.Hour).Format(time.RFC3339Nano)
	expiringText := now.Add(24 * time.Hour).Format(time.RFC3339Nano)
	activeText := now.Add(30 * 24 * time.Hour).Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('expired-user', 'Expired', 'active', ?, ?, ?)`, []any{expiredText, nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('expiring-user', 'Expiring', 'active', ?, ?, ?)`, []any{expiringText, nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('active-user', 'Active', 'active', ?, ?, ?)`, []any{activeText, nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, last_seen_at, created_at, updated_at) VALUES ('stale-node', 'stale-node', 'Stale', 'relay', 'online', ?, ?, ?)`, []any{oldText, nowText, nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed status data: %v", err)
		}
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	_ = doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard", token, nil)
	var expiredStatus, expiringStatus, activeStatus, nodeStatus string
	queries := []struct {
		query  string
		target *string
	}{
		{"SELECT status FROM users WHERE id = 'expired-user'", &expiredStatus},
		{"SELECT status FROM users WHERE id = 'expiring-user'", &expiringStatus},
		{"SELECT status FROM users WHERE id = 'active-user'", &activeStatus},
		{"SELECT health_status FROM nodes WHERE id = 'stale-node'", &nodeStatus},
	}
	for _, item := range queries {
		if err := database.QueryRow(item.query).Scan(item.target); err != nil {
			t.Fatalf("query status: %v", err)
		}
	}
	if expiredStatus != "expired" || expiringStatus != "expiring" || activeStatus != "active" || nodeStatus != "offline" {
		t.Fatalf("statuses expired=%s expiring=%s active=%s node=%s", expiredStatus, expiringStatus, activeStatus, nodeStatus)
	}
}

func TestListEndpointsWithData(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, hostname, health_status, created_at, updated_at) VALUES ('n1', 'relay-1', '线路机 1', 'relay', 'relay.example', 'online', ?, ?)`, []any{now, now}},
		{`INSERT INTO nodes (id, node_key, name, type, hostname, health_status, created_at, updated_at) VALUES ('n2', 'landing-1', '落地机 1', 'landing', 'landing.example', 'offline', ?, ?)`, []any{now, now}},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, created_at, updated_at) VALUES ('u1', '测试用户', 'active', 100, 'CNY', ?, ?)`, []any{now, now}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, client_count, up, down, first_seen_at) VALUES ('i1', 'n1', '15', 'u1', 'user', 'user-15', 2, 100, 200, ?)`, []any{now}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('ui1', 'u1', 'i1', 1, ?)`, []any{now}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, landing_inbound_tag, created_at, updated_at) VALUES ('r1', '线路 A', 'n1', 'n2', 'ss-in', ?, ?)`, []any{now, now}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('ur1', 'u1', 'r1', 1, ?)`, []any{now}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, provider, monthly_cost, currency, created_at, updated_at) VALUES ('e1', 'n2', '203.0.113.10', 'Test ISP', 20, 'CNY', ?, ?)`, []any{now, now}},
		{`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, allocation_weight) VALUES ('rei1', 'r1', 'e1', 1)`, nil},
		{`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, created_at) VALUES ('c1', 'n1', 'server', 30, 'CNY', '2026-01-01', ?)`, []any{now}},
		{`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES ('ev1', 'n1', 'sync_failed', 'error', 'sample event', ?)`, []any{now}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed data: %v", err)
		}
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	headers := token
	paths := []string{"/api/users", "/api/nodes", "/api/routes", "/api/exit-ips", "/api/events", "/api/costs/summary"}
	for _, path := range paths {
		result := doJSON(t, ts.Client(), http.MethodGet, ts.URL+path, headers, nil)
		if result["code"] != successCode {
			t.Fatalf("%s response: %#v", path, result)
		}
	}
}

func doJSON(t *testing.T, client *http.Client, method, url, token string, payload any) map[string]any {
	t.Helper()
	return doJSONWithRequestID(t, client, method, url, token, "", payload)
}

func doJSONWithRequestID(t *testing.T, client *http.Client, method, url, token, requestID string, payload any) map[string]any {
	t.Helper()
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("encode request: %v", err)
		}
		body = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	if requestID != "" {
		request.Header.Set("X-Request-ID", requestID)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("perform request: %v", err)
	}
	defer response.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response (%d): %v", response.StatusCode, err)
	}
	return result
}
