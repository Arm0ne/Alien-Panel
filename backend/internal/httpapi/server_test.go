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
		AdminUsername: "admin",
		AdminPassword: "test-password",
		SessionTTL:    time.Hour,
		CorsOrigins:   []string{"http://localhost:9527"},
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
