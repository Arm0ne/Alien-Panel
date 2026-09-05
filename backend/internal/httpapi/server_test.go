package httpapi

import (
	"bytes"
	"context"
	"crypto/tls"
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

func TestSecurityHeaders(t *testing.T) {
	server, _ := testServer(t)
	handler := server.Handler()

	local := httptest.NewRecorder()
	handler.ServeHTTP(local, httptest.NewRequest(http.MethodGet, "/health/live", nil))
	for header, want := range map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
		"Permissions-Policy":      "camera=(), microphone=(), geolocation=()",
		"Content-Security-Policy": "default-src 'none'; frame-ancestors 'none'",
		"Cache-Control":           "no-store",
	} {
		if got := local.Header().Get(header); got != want {
			t.Fatalf("local %s = %q, want %q", header, got, want)
		}
	}
	if got := local.Header().Get("Strict-Transport-Security"); got != "" {
		t.Fatalf("local HSTS = %q, want absent", got)
	}

	proxied := httptest.NewRecorder()
	proxiedRequest := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	proxiedRequest.Header.Set("X-Forwarded-Proto", "https")
	handler.ServeHTTP(proxied, proxiedRequest)
	if got := proxied.Header().Get("Strict-Transport-Security"); got != "max-age=31536000; includeSubDomains" {
		t.Fatalf("proxied HSTS = %q", got)
	}

	tlsRequest := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	tlsRequest.TLS = &tls.ConnectionState{}
	secure := httptest.NewRecorder()
	handler.ServeHTTP(secure, tlsRequest)
	if got := secure.Header().Get("Strict-Transport-Security"); got != "max-age=31536000; includeSubDomains" {
		t.Fatalf("TLS HSTS = %q", got)
	}
}

func TestAdministratorWritesValidateOrigin(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	payload := map[string]any{"nodeKey": "csrf-node", "name": "CSRF 测试节点", "type": "landing"}
	postWithOrigin := func(origin, referer string) (int, map[string]any) {
		encoded, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest(http.MethodPost, "/api/nodes", bytes.NewReader(encoded))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token)
		if origin != "" {
			request.Header.Set("Origin", origin)
		}
		if referer != "" {
			request.Header.Set("Referer", referer)
		}
		recorder := httptest.NewRecorder()
		server.Handler().ServeHTTP(recorder, request)
		var response map[string]any
		if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
			t.Fatalf("decode origin response: %v", err)
		}
		return recorder.Code, response
	}

	status, response := postWithOrigin("https://evil.example", "")
	if status != http.StatusForbidden || response["code"] != csrfCode {
		t.Fatalf("evil origin status=%d response=%#v", status, response)
	}
	status, response = postWithOrigin("", "https://evil.example/form")
	if status != http.StatusForbidden || response["code"] != csrfCode {
		t.Fatalf("evil referer status=%d response=%#v", status, response)
	}
	status, response = postWithOrigin("http://localhost:9527", "")
	if status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("allowed origin status=%d response=%#v", status, response)
	}
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM nodes WHERE node_key = 'csrf-node'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("allowed origin node count=%d err=%v", count, err)
	}
}

func TestSessionRefreshAndLogoutRevokeCredentials(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	loginData := login["data"].(map[string]any)
	oldToken := loginData["token"].(string)
	oldRefresh := loginData["refreshToken"].(string)
	refreshed := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/refreshToken", "", map[string]string{"refreshToken": oldRefresh})
	if refreshed["code"] != successCode {
		t.Fatalf("refresh response = %#v", refreshed)
	}
	refreshedData := refreshed["data"].(map[string]any)
	newToken := refreshedData["token"].(string)
	newRefresh := refreshedData["refreshToken"].(string)
	if newToken == oldToken || newRefresh == oldRefresh {
		t.Fatalf("refresh did not rotate credentials old=%q/%q new=%q/%q", oldToken, oldRefresh, newToken, newRefresh)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodGet, ts.URL+"/api/auth/me", oldToken, nil); status != http.StatusUnauthorized || response["code"] != unauthorizedCode {
		t.Fatalf("old access token status=%d response=%#v", status, response)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/refreshToken", "", map[string]string{"refreshToken": oldRefresh}); status != http.StatusUnauthorized || response["code"] != unauthorizedCode {
		t.Fatalf("old refresh token status=%d response=%#v", status, response)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodGet, ts.URL+"/api/auth/me", newToken, nil); status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("new access token status=%d response=%#v", status, response)
	}
	logout := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/logout", newToken, nil)
	if logout["code"] != successCode {
		t.Fatalf("logout response = %#v", logout)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodGet, ts.URL+"/api/auth/me", newToken, nil); status != http.StatusUnauthorized || response["code"] != unauthorizedCode {
		t.Fatalf("logged-out access token status=%d response=%#v", status, response)
	}
	var activeSessions int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sessions WHERE revoked_at IS NULL`).Scan(&activeSessions); err != nil || activeSessions != 0 {
		t.Fatalf("active sessions=%d err=%v, want 0", activeSessions, err)
	}
}

func TestAgentRegistrationRotatesPreviousToken(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	register := func(name string) string {
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/agent/v1/register", strings.NewReader(`{"node_key":"rotate-agent","node_name":"`+name+`","node_type":"relay"}`))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Agent-Registration-Token", "bootstrap-secret")
		response, err := ts.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		var body map[string]any
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != http.StatusOK || body["code"] != successCode {
			t.Fatalf("registration status=%d body=%#v", response.StatusCode, body)
		}
		return body["data"].(map[string]any)["token"].(string)
	}

	firstToken := register("轮换 Agent 1")
	secondToken := register("轮换 Agent 2")
	if firstToken == secondToken {
		t.Fatal("agent registration returned the same token")
	}
	payload := map[string]any{"node_key": "rotate-agent", "observed_at": time.Now().UTC().Format(time.RFC3339Nano), "status": map[string]any{"xray_running": true}}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", firstToken, payload); status != http.StatusUnauthorized || response["code"] != unauthorizedCode {
		t.Fatalf("revoked agent token status=%d response=%#v", status, response)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", secondToken, payload); status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("rotated agent token status=%d response=%#v", status, response)
	}
	var revoked int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_credentials WHERE node_id = (SELECT id FROM nodes WHERE node_key = 'rotate-agent') AND revoked_at IS NOT NULL`).Scan(&revoked); err != nil || revoked != 1 {
		t.Fatalf("revoked agent credentials=%d err=%v, want 1", revoked, err)
	}
}

func TestAuditLogsRedactSensitiveFields(t *testing.T) {
	server, database := testServer(t)
	request := httptest.NewRequest(http.MethodPatch, "/api/test", nil)
	server.writeAuditLog(request, "test.redaction", "test", "redaction-1",
		map[string]any{
			"token":  "access-secret",
			"nested": map[string]any{"password": "password-secret", "safe": "kept"},
		},
		map[string]any{"client_secret": "client-secret", "subscriptionUrl": "subscription-secret", "value": "kept"})

	var beforeJSON, afterJSON string
	if err := database.QueryRow(`SELECT before_json, after_json FROM audit_logs WHERE action = 'test.redaction'`).Scan(&beforeJSON, &afterJSON); err != nil {
		t.Fatalf("read redacted audit log: %v", err)
	}
	for _, secret := range []string{"access-secret", "password-secret", "client-secret", "subscription-secret"} {
		if strings.Contains(beforeJSON, secret) || strings.Contains(afterJSON, secret) {
			t.Fatalf("audit log contains secret %q: before=%s after=%s", secret, beforeJSON, afterJSON)
		}
	}
	if !strings.Contains(beforeJSON, "[REDACTED]") || !strings.Contains(afterJSON, "[REDACTED]") || !strings.Contains(beforeJSON, "kept") || !strings.Contains(afterJSON, "kept") {
		t.Fatalf("unexpected redacted audit state: before=%s after=%s", beforeJSON, afterJSON)
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

func TestLandingInboundMappingIsHiddenButUserIsKept(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	for _, statement := range []string{
		`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('legacy-landing-node', 'legacy-landing-node', '历史落地机', 'landing', 'online', '` + nowText + `', '` + nowText + `')`,
		`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('legacy-landing-user', '保留的预留用户', 'active', '` + nowText + `', '` + nowText + `')`,
		`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at, last_seen_at) VALUES ('legacy-landing-inbound', 'legacy-landing-node', '77', 'legacy-landing-user', 'user', '落地入口-旧记录', '` + nowText + `', '` + nowText + `')`,
		`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('legacy-landing-user-inbound', 'legacy-landing-user', 'legacy-landing-inbound', 1, '` + nowText + `')`,
	} {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("seed legacy landing user: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?page_size=20", token, nil)
	if list["code"] != successCode {
		t.Fatalf("user list response = %#v", list)
	}
	items := list["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected preserved user in list, got %#v", items)
	}
	item := items[0].(map[string]any)
	if item["id"] != "legacy-landing-user" || item["nodeId"] != "" || item["inboundTag"] != "" {
		t.Fatalf("landing inbound leaked into user summary: %#v", item)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/legacy-landing-user", token, nil)
	if detail["code"] != successCode {
		t.Fatalf("user detail response = %#v", detail)
	}
	inbound := detail["data"].(map[string]any)["inbound"].(map[string]any)
	node := detail["data"].(map[string]any)["node"].(map[string]any)
	if inbound["id"] != nil || inbound["tag"] != nil || node["id"] != nil {
		t.Fatalf("landing inbound leaked into user detail: inbound=%#v node=%#v", inbound, node)
	}

	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("begin cleanup sync: %v", err)
	}
	if err := server.ensureRelayInboundUser(tx, "landing", "legacy-landing-inbound", "77", agentInboundPayload{Enable: true}, "", time.Now().UTC()); err != nil {
		_ = tx.Rollback()
		t.Fatalf("clean legacy landing mapping: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit cleanup sync: %v", err)
	}
	var mappings, users int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_inbounds WHERE inbound_id = 'legacy-landing-inbound'`).Scan(&mappings); err != nil {
		t.Fatalf("count cleaned mapping: %v", err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM users WHERE id = 'legacy-landing-user'`).Scan(&users); err != nil {
		t.Fatalf("count preserved user: %v", err)
	}
	if mappings != 0 || users != 1 {
		t.Fatalf("runtime cleanup mappings=%d users=%d", mappings, users)
	}
}

func TestUserPathAssignmentLifecycle(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('path-relay', 'path-relay', '路径线路机', 'relay', 1, 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('path-landing', 'path-landing', '路径落地机', 'landing', 1, 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('path-other', 'path-other', '其他落地机', 'landing', 1, 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('path-disabled-landing', 'path-disabled-landing', '停用落地机', 'landing', 0, 'offline', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('path-user', '路径测试用户', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, protocol, port, enable, client_count, up, down, all_time, first_seen_at, last_seen_at) VALUES ('path-inbound', 'path-relay', 'path-1', 'path-user', 'user', 'path-user-inbound', 'vless', 443, 1, 1, 0, 0, 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, protocol, port, enable, client_count, up, down, all_time, first_seen_at, last_seen_at) VALUES ('path-landing-inbound', 'path-landing', 'landing-1', 'socks', 'path-landing-entry', 'socks', 10001, 1, 0, 0, 0, 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('path-user-inbound', 'path-user', 'path-inbound', 1, ?)`, []any{nowText}},
		{`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, enabled, created_at, updated_at) VALUES ('path-relay-ip', 'node', 'path-relay', '198.51.100.71', 4, 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, enabled, created_at, updated_at) VALUES ('path-landing-ip', 'node', 'path-landing', '198.51.100.72', 4, 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, enabled, created_at, updated_at) VALUES ('path-other-ip', 'node', 'path-other', '198.51.100.73', 4, 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, ip, family, enabled, created_at, updated_at) VALUES ('path-s5-ip', 's5', '198.51.100.74', 4, 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, enabled, created_at, updated_at) VALUES ('path-disabled-ip', 'node', 'path-relay', '198.51.100.75', 4, 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, owner_node_id, ip, family, enabled, created_at, updated_at) VALUES ('path-disabled-node-ip', 'node', 'path-disabled-landing', '198.51.100.76', 4, 1, ?, ?)`, []any{nowText, nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed user path data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)

	assets := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/path-user/path-assets", token, nil)
	if assets["code"] != successCode {
		t.Fatalf("path assets response = %#v", assets)
	}
	assetsData := assets["data"].(map[string]any)
	if assetsData["relay"].(map[string]any)["id"] != "path-relay" || len(assetsData["relayExitIps"].([]any)) != 1 || len(assetsData["externalExitIps"].([]any)) != 1 {
		t.Fatalf("unexpected relay path assets = %#v", assetsData)
	}
	landingAssets := assetsData["landingNodes"].([]any)
	if len(landingAssets) != 3 {
		t.Fatalf("expected all landing nodes in path assets, got %#v", landingAssets)
	}
	var foundLandingInbound bool
	var foundPendingLanding bool
	for _, raw := range landingAssets {
		item := raw.(map[string]any)
		if item["id"] == "path-other" && item["inboundState"] == "pending" && len(item["inbounds"].([]any)) == 0 {
			foundPendingLanding = true
		}
		if item["id"] != "path-landing" {
			continue
		}
		inbounds := item["inbounds"].([]any)
		if item["inboundState"] != "ready" || len(inbounds) != 1 || inbounds[0].(map[string]any)["purpose"] != "infrastructure" {
			t.Fatalf("unexpected landing inbound asset = %#v", item)
		}
		foundLandingInbound = true
	}
	if !foundLandingInbound {
		t.Fatalf("path landing asset not found: %#v", landingAssets)
	}
	if !foundPendingLanding {
		t.Fatalf("landing node without a successful sync should be marked pending: %#v", landingAssets)
	}
	var usersBeforeLandingInbound int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&usersBeforeLandingInbound); err != nil {
		t.Fatalf("count users before landing inbound classification: %v", err)
	}
	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("begin landing inbound classification: %v", err)
	}
	if err := server.ensureRelayInboundUser(tx, "landing", "path-landing-inbound", "landing-1", agentInboundPayload{Enable: true}, "", now); err != nil {
		_ = tx.Rollback()
		t.Fatalf("classify landing inbound: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit landing inbound classification: %v", err)
	}
	var usersAfterLandingInbound int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&usersAfterLandingInbound); err != nil || usersAfterLandingInbound != usersBeforeLandingInbound {
		t.Fatalf("landing inbound unexpectedly created a user: before=%d after=%d err=%v", usersBeforeLandingInbound, usersAfterLandingInbound, err)
	}

	direct := doJSON(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"exitIpId": "path-relay-ip", "notes": "线路机直出"})
	if direct["code"] != successCode {
		t.Fatalf("direct path response = %#v", direct)
	}
	directPath := direct["data"].(map[string]any)["path"].(map[string]any)
	if directPath["mode"] != "relay" || directPath["relayNodeId"] != "path-relay" || directPath["exitIpAddress"] != "198.51.100.71" {
		t.Fatalf("unexpected direct path = %#v", directPath)
	}

	missingLandingInboundStatus, missingLandingInbound := doJSONWithStatus(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"landingNodeId": "path-landing", "exitIpId": "path-landing-ip"})
	if missingLandingInboundStatus != http.StatusConflict || missingLandingInbound["code"] != validationCode {
		t.Fatalf("missing landing inbound status=%d response=%#v", missingLandingInboundStatus, missingLandingInbound)
	}

	landing := doJSON(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"landingNodeId": "path-landing", "landingInboundId": "path-landing-inbound", "exitIpId": "path-landing-ip", "notes": "落地机固定出口"})
	if landing["code"] != successCode {
		t.Fatalf("landing path response = %#v", landing)
	}
	landingData := landing["data"].(map[string]any)
	landingPath := landingData["path"].(map[string]any)
	if landingPath["mode"] != "landing" || landingPath["landingNodeId"] != "path-landing" || len(landingData["pathHistory"].([]any)) != 2 {
		t.Fatalf("unexpected landing path/history = %#v", landingData)
	}

	mismatchStatus, mismatch := doJSONWithStatus(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"landingNodeId": "path-landing", "landingInboundId": "path-landing-inbound", "exitIpId": "path-other-ip"})
	if mismatchStatus != http.StatusBadRequest || mismatch["code"] != validationCode {
		t.Fatalf("mismatch path status=%d response=%#v", mismatchStatus, mismatch)
	}
	mixedS5Status, mixedS5 := doJSONWithStatus(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"landingNodeId": "path-landing", "landingInboundId": "path-landing-inbound", "exitIpId": "path-s5-ip"})
	if mixedS5Status != http.StatusBadRequest || mixedS5["code"] != validationCode {
		t.Fatalf("mixed S5 path status=%d response=%#v", mixedS5Status, mixedS5)
	}
	disabledStatus, disabled := doJSONWithStatus(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"exitIpId": "path-disabled-ip"})
	if disabledStatus != http.StatusConflict || disabled["code"] != validationCode {
		t.Fatalf("disabled path status=%d response=%#v", disabledStatus, disabled)
	}
	disabledNodeStatus, disabledNode := doJSONWithStatus(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"landingNodeId": "path-disabled-landing", "exitIpId": "path-disabled-node-ip"})
	if disabledNodeStatus != http.StatusConflict || disabledNode["code"] != validationCode {
		t.Fatalf("disabled node path status=%d response=%#v", disabledNodeStatus, disabledNode)
	}

	s5 := doJSON(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/path-user/path", token, map[string]any{"exitIpId": "path-s5-ip"})
	if s5["code"] != successCode || s5["data"].(map[string]any)["path"].(map[string]any)["mode"] != "external" {
		t.Fatalf("S5 path response = %#v", s5)
	}

	cleared := doJSON(t, ts.Client(), http.MethodDelete, ts.URL+"/api/users/path-user/path", token, nil)
	clearedData := cleared["data"].(map[string]any)
	if cleared["code"] != successCode || clearedData["path"] != nil || len(clearedData["pathHistory"].([]any)) != 3 {
		t.Fatalf("cleared path response = %#v", cleared)
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
		{`INSERT INTO sync_runs (id, node_id, sync_id, started_at, finished_at, status, inbound_count, client_count) VALUES ('sync-list', 'n1', 'sync-list-1', ?, ?, 'success', 1, 2)`, []any{now, now}},
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
		data, ok := result["data"].(map[string]any)
		if !ok {
			t.Fatalf("%s data is not an object: %#v", path, result["data"])
		}
		if data["dataAt"] == nil || data["dataAt"] == "" {
			t.Fatalf("%s did not return latest sync dataAt: %#v", path, data)
		}
	}
}

func TestUserListAggregatesOneRowPerInbound(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	statements := []string{
		`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('aggregate-relay-1', 'aggregate-relay-1', '线路机 A', 'relay', 'online', '` + now + `', '` + now + `')`,
		`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('aggregate-relay-2', 'aggregate-relay-2', '线路机 B', 'relay', 'online', '` + now + `', '` + now + `')`,
		`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('aggregate-user-1', '用户 A', 'active', '` + now + `', '` + now + `')`,
		`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('aggregate-user-2', '用户 B', 'active', '` + now + `', '` + now + `')`,
		`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, client_count, up, down, first_seen_at, last_seen_at) VALUES ('aggregate-inbound-1', 'aggregate-relay-1', '101', 'aggregate-user-1', 'user', 'user-a', 2, 100, 200, '` + now + `', '` + now + `')`,
		`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, client_count, up, down, first_seen_at, last_seen_at) VALUES ('aggregate-inbound-2', 'aggregate-relay-2', '101', 'aggregate-user-2', 'user', 'user-b', 1, 300, 400, '` + now + `', '` + now + `')`,
		`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('aggregate-mapping-1', 'aggregate-user-1', 'aggregate-inbound-1', 1, '` + now + `')`,
		`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('aggregate-mapping-2', 'aggregate-user-2', 'aggregate-inbound-2', 1, '` + now + `')`,
		`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, all_time) VALUES ('aggregate-client-1', 'aggregate-relay-1', 'aggregate-inbound-1', 'phone-a', 'shared@example.com', 1, 100)`,
		`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, all_time) VALUES ('aggregate-client-2', 'aggregate-relay-1', 'aggregate-inbound-1', 'laptop-a', 'shared@example.com', 1, 200)`,
		`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, all_time) VALUES ('aggregate-client-3', 'aggregate-relay-2', 'aggregate-inbound-2', 'phone-b', 'shared@example.com', 1, 300)`,
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("seed aggregation data: %v", err)
		}
	}

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	result := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?page_size=20", token, nil)
	if result["code"] != successCode {
		t.Fatalf("user list response = %#v", result)
	}
	data := result["data"].(map[string]any)
	if data["total"] != float64(2) {
		t.Fatalf("user list total = %v, want 2", data["total"])
	}
	items := data["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("user list items = %d, want one row per Inbound", len(items))
	}

	byNode := make(map[string]map[string]any, len(items))
	for _, item := range items {
		row := item.(map[string]any)
		nodeName := row["nodeName"].(string)
		if _, exists := byNode[nodeName]; exists {
			t.Fatalf("duplicate user row for node %q: %#v", nodeName, row)
		}
		byNode[nodeName] = row
	}
	if byNode["线路机 A"]["clientCount"] != float64(2) || byNode["线路机 B"]["clientCount"] != float64(1) {
		t.Fatalf("client counts were not aggregated per Inbound: %#v", byNode)
	}
	if byNode["线路机 A"]["id"] == byNode["线路机 B"]["id"] {
		t.Fatalf("same Email across nodes incorrectly merged users: %#v", byNode)
	}

	nodeFiltered := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?node_id=aggregate-relay-1", token, nil)
	filteredData := nodeFiltered["data"].(map[string]any)
	if filteredData["total"] != float64(1) || len(filteredData["items"].([]any)) != 1 {
		t.Fatalf("node-filtered user list = %#v", filteredData)
	}
}

func TestNodeDetailAndManualSyncRequest(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, hostname, public_ip, region, provider, health_status, created_at, updated_at) VALUES ('node-detail', 'node-detail', '东京落地机', 'landing', 'landing.example', '203.0.113.20', '东京', 'Test ISP', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, protocol, port, enable, client_count, up, down, all_time, first_seen_at, last_seen_at) VALUES ('node-inbound', 'node-detail', '7', 'infrastructure', 'ss-entry', 'shadowsocks', 8443, 1, 4, 10, 20, 30, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, provider, monthly_cost, currency, created_at, updated_at) VALUES ('node-exit', 'node-detail', '198.51.100.10', 'Test ISP', 12, 'CNY', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO sync_runs (id, node_id, sync_id, started_at, finished_at, status, inbound_count, client_count) VALUES ('node-sync', 'node-detail', 'sync-1', ?, ?, 'success', 1, 4)`, []any{nowText, nowText}},
		{`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at) VALUES ('node-event', 'node-detail', 'sync_failed', 'error', 'sample failure', ?)`, []any{nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed node detail data: %v", err)
		}
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/node-detail", token, nil)
	if detail["code"] != successCode {
		t.Fatalf("node detail response: %#v", detail)
	}
	detailData := detail["data"].(map[string]any)
	if detailData["name"] != "东京落地机" || len(detailData["inbounds"].([]any)) != 1 || len(detailData["exitIps"].([]any)) != 1 {
		t.Fatalf("unexpected node detail: %#v", detailData)
	}
	queued := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes/node-detail/sync", token, nil)
	if queued["code"] != successCode || queued["data"].(map[string]any)["status"] != "queued" {
		t.Fatalf("manual sync response: %#v", queued)
	}
	var eventType string
	if err := database.QueryRow(`SELECT event_type FROM node_events WHERE id = ?`, queued["data"].(map[string]any)["requestId"].(string)).Scan(&eventType); err != nil || eventType != "sync_requested" {
		t.Fatalf("sync request event type = %q, err=%v", eventType, err)
	}
}

func TestRouteCRUDAndBindingProtection(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('route-relay', 'route-relay', '线路机 A', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('route-landing', 'route-landing', '落地机 A', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('route-other', 'route-other', '其他线路机', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at, last_seen_at) VALUES ('route-landing-inbound', 'route-landing', '99', 'infrastructure', 'ss-landing-a', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at, last_seen_at) VALUES ('route-other-inbound', 'route-other', '100', 'infrastructure', 'ss-other', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('route-user', '线路用户', 'active', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, created_at, updated_at) VALUES ('route-exit-ip', 'route-landing', '198.51.100.42', ?, ?)`, []any{nowText, nowText}},
	}
	for _, statement := range seed {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed route data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "东京线路 A", "relayNodeId": "route-relay", "landingNodeId": "route-landing",
		"relayOutboundTag": "to-landing", "landingInboundId": "route-landing-inbound", "landingInboundTag": "ss-landing-a",
		"validFrom": "2026-09-01", "validTo": "2026-09-30", "notes": "首条线路",
	})
	if created["code"] != successCode {
		t.Fatalf("create route response = %#v", created)
	}
	createdData := created["data"].(map[string]any)
	routeID := createdData["id"].(string)
	if createdData["name"] != "东京线路 A" || createdData["relayNodeName"] != "线路机 A" || createdData["landingInboundId"] != "route-landing-inbound" {
		t.Fatalf("created route data = %#v", createdData)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/routes/"+routeID, token, nil)
	if detail["code"] != successCode || detail["data"].(map[string]any)["id"] != routeID {
		t.Fatalf("route detail response = %#v", detail)
	}

	updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/routes/"+routeID, token, map[string]any{
		"name": "东京线路 A（更新）", "enabled": false, "validTo": "2026-10-01", "notes": "已更新",
	})
	if updated["code"] != successCode || updated["data"].(map[string]any)["name"] != "东京线路 A（更新）" || updated["data"].(map[string]any)["enabled"] != false {
		t.Fatalf("update route response = %#v", updated)
	}

	invalidRelay := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "非法线路", "relayNodeId": "route-landing", "landingNodeId": "route-other",
	})
	if invalidRelay["code"] != validationCode {
		t.Fatalf("invalid relay type response = %#v", invalidRelay)
	}
	invalidInbound := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "非法落地 Inbound", "relayNodeId": "route-relay", "landingNodeId": "route-landing", "landingInboundId": "route-other-inbound",
	})
	if invalidInbound["code"] != validationCode {
		t.Fatalf("invalid landing inbound response = %#v", invalidInbound)
	}
	invalidDates := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "非法日期", "relayNodeId": "route-relay", "landingNodeId": "route-landing", "validFrom": "2026-10-02", "validTo": "2026-10-01",
	})
	if invalidDates["code"] != validationCode {
		t.Fatalf("invalid route dates response = %#v", invalidDates)
	}

	deletable := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "可删除线路", "relayNodeId": "route-relay", "landingNodeId": "route-landing",
	})
	if deletable["code"] != successCode {
		t.Fatalf("create deletable route response = %#v", deletable)
	}
	deletableID := deletable["data"].(map[string]any)["id"].(string)
	deleteStatus, deleted := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/routes/"+deletableID, token, nil)
	if deleteStatus != http.StatusOK || deleted["code"] != successCode || deleted["data"].(map[string]any)["deleted"] != true {
		t.Fatalf("delete unbound route status=%d response=%#v", deleteStatus, deleted)
	}

	if _, err := database.Exec(`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('route-user-binding', 'route-user', ?, 1, ?)`, routeID, nowText); err != nil {
		t.Fatalf("bind route to user: %v", err)
	}
	conflictStatus, conflict := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/routes/"+routeID, token, nil)
	if conflictStatus != http.StatusConflict || conflict["code"] != validationCode {
		t.Fatalf("delete bound route status=%d response=%#v", conflictStatus, conflict)
	}

	exitBound := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes", token, map[string]any{
		"name": "出口 IP 绑定线路", "relayNodeId": "route-relay", "landingNodeId": "route-landing",
	})
	if exitBound["code"] != successCode {
		t.Fatalf("create exit-bound route response = %#v", exitBound)
	}
	exitRouteID := exitBound["data"].(map[string]any)["id"].(string)
	if _, err := database.Exec(`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, allocation_weight) VALUES ('route-exit-binding', ?, 'route-exit-ip', 1)`, exitRouteID); err != nil {
		t.Fatalf("bind route to exit IP: %v", err)
	}
	exitConflictStatus, exitConflict := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/routes/"+exitRouteID, token, nil)
	if exitConflictStatus != http.StatusConflict || exitConflict["code"] != validationCode {
		t.Fatalf("delete exit-bound route status=%d response=%#v", exitConflictStatus, exitConflict)
	}
}

func TestUserRouteAssignmentWithFixedExit(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC().Format(time.RFC3339Nano)
	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('assign-relay', 'assign-relay', '分配线路机', 'relay', 'online', ?, ?)`, []any{now, now}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('assign-landing', 'assign-landing', '分配落地机', 'landing', 'online', ?, ?)`, []any{now, now}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('assign-user', '待分配用户', 'active', ?, ?)`, []any{now, now}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('assign-route', '分配线路', 'assign-relay', 'assign-landing', 1, ?, ?)`, []any{now, now}},
		{`INSERT INTO exit_ips (id, landing_node_id, source_type, owner_node_id, ip, enabled, created_at, updated_at) VALUES ('assign-ip', NULL, 's5', NULL, '198.51.100.99', 1, ?, ?)`, []any{now, now}},
		{`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, scope, enabled) VALUES ('assign-binding', 'assign-route', 'assign-ip', 'external', 1)`, nil},
	}
	for _, statement := range seed {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed assignment data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	assigned := doJSON(t, ts.Client(), http.MethodPut, ts.URL+"/api/users/assign-user/route", token, map[string]any{
		"routeId": "assign-route", "routeExitIpId": "assign-binding",
	})
	if assigned["code"] != successCode {
		t.Fatalf("assign user route response = %#v", assigned)
	}
	assignedData := assigned["data"].(map[string]any)
	routes := assignedData["routes"].([]any)
	if len(routes) != 1 {
		t.Fatalf("assigned routes = %#v", routes)
	}
	routeData := routes[0].(map[string]any)
	if routeData["routeExitIpId"] != "assign-binding" || routeData["exitIpAddress"] != "198.51.100.99" || routeData["assignmentMode"] != "fixed" {
		t.Fatalf("assigned route details = %#v", routeData)
	}

	cleared := doJSON(t, ts.Client(), http.MethodDelete, ts.URL+"/api/users/assign-user/route", token, nil)
	if cleared["code"] != successCode || len(cleared["data"].(map[string]any)["routes"].([]any)) != 0 {
		t.Fatalf("clear user route response = %#v", cleared)
	}
	var active int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_routes WHERE user_id = 'assign-user' AND active_to IS NULL`).Scan(&active); err != nil || active != 0 {
		t.Fatalf("active user route count = %d, err=%v", active, err)
	}
	unboundStatus, unbound := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/routes/assign-route/exit-ips/assign-ip", token, nil)
	if unboundStatus != http.StatusOK || unbound["code"] != successCode {
		t.Fatalf("unbind released fixed assignment status=%d response=%#v", unboundStatus, unbound)
	}
}

func TestExitIPCRUDAndBindingProtection(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('exit-landing', 'exit-landing', '落地机 B', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('exit-relay', 'exit-relay', '线路机 B', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('exit-binding-route', '出口绑定线路', 'exit-relay', 'exit-landing', 1, ?, ?)`, []any{nowText, nowText}},
	}
	for _, statement := range seed {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed exit IP data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "203.0.113.42", "landingNodeId": "exit-landing", "family": 4,
		"provider": "Test ISP", "monthlyCost": 18.5, "validFrom": "2026-09-01", "validTo": "2026-12-31", "notes": "主出口",
	})
	if created["code"] != successCode {
		t.Fatalf("create exit IP response = %#v", created)
	}
	createdData := created["data"].(map[string]any)
	exitIPID := createdData["id"].(string)
	if createdData["address"] != "203.0.113.42" || createdData["landingNodeName"] != "落地机 B" || createdData["family"] != float64(4) {
		t.Fatalf("created exit IP data = %#v", createdData)
	}
	exitList := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips?page=1&page_size=20", token, nil)
	listData := exitList["data"].(map[string]any)
	listItems := listData["items"].([]any)
	var listedExitIP map[string]any
	for _, raw := range listItems {
		item := raw.(map[string]any)
		if item["id"] == exitIPID {
			listedExitIP = item
			break
		}
	}
	if exitList["code"] != successCode || listedExitIP == nil || listedExitIP["landingNodeId"] != "exit-landing" {
		t.Fatalf("exit IP list should include landingNodeId: %#v", exitList)
	}
	ipv6 := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "2001:db8::42", "landingNodeId": "exit-landing",
	})
	if ipv6["code"] != successCode || ipv6["data"].(map[string]any)["family"] != float64(6) {
		t.Fatalf("IPv6 family inference response = %#v", ipv6)
	}
	ipv6ID := ipv6["data"].(map[string]any)["id"].(string)
	if status, result := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/exit-ips/"+ipv6ID, token, nil); status != http.StatusOK || result["code"] != successCode {
		t.Fatalf("delete inferred IPv6 response status=%d response=%#v", status, result)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/"+exitIPID, token, nil)
	if detail["code"] != successCode || detail["data"].(map[string]any)["id"] != exitIPID {
		t.Fatalf("exit IP detail response = %#v", detail)
	}
	updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/exit-ips/"+exitIPID, token, map[string]any{
		"provider": "Updated ISP", "monthlyCost": 20, "enabled": false, "validTo": "2027-01-01", "notes": "已更新",
	})
	if updated["code"] != successCode || updated["data"].(map[string]any)["provider"] != "Updated ISP" || updated["data"].(map[string]any)["enabled"] != false {
		t.Fatalf("update exit IP response = %#v", updated)
	}

	invalidNode := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "203.0.113.43", "landingNodeId": "exit-relay", "family": 4,
	})
	if invalidNode["code"] != validationCode {
		t.Fatalf("invalid landing node response = %#v", invalidNode)
	}
	invalidAddress := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "not-an-ip", "landingNodeId": "exit-landing", "family": 4,
	})
	if invalidAddress["code"] != validationCode {
		t.Fatalf("invalid address response = %#v", invalidAddress)
	}
	invalidFamily := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "2001:db8::42", "landingNodeId": "exit-landing", "family": 4,
	})
	if invalidFamily["code"] != validationCode {
		t.Fatalf("invalid family response = %#v", invalidFamily)
	}
	invalidDates := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "203.0.113.44", "landingNodeId": "exit-landing", "family": 4, "validFrom": "2026-10-02", "validTo": "2026-10-01",
	})
	if invalidDates["code"] != validationCode {
		t.Fatalf("invalid dates response = %#v", invalidDates)
	}

	deletable := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, map[string]any{
		"address": "203.0.113.45", "landingNodeId": "exit-landing", "family": 4,
	})
	if deletable["code"] != successCode {
		t.Fatalf("create deletable exit IP response = %#v", deletable)
	}
	deletableID := deletable["data"].(map[string]any)["id"].(string)
	deleteStatus, deleted := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/exit-ips/"+deletableID, token, nil)
	if deleteStatus != http.StatusOK || deleted["code"] != successCode || deleted["data"].(map[string]any)["deleted"] != true {
		t.Fatalf("delete unbound exit IP status=%d response=%#v", deleteStatus, deleted)
	}

	if _, err := database.Exec(`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, allocation_weight) VALUES ('exit-binding', 'exit-binding-route', ?, 1)`, exitIPID); err != nil {
		t.Fatalf("bind exit IP to route: %v", err)
	}
	conflictStatus, conflict := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/exit-ips/"+exitIPID, token, nil)
	if conflictStatus != http.StatusConflict || conflict["code"] != validationCode {
		t.Fatalf("delete bound exit IP status=%d response=%#v", conflictStatus, conflict)
	}
}

func TestRouteExitIPBindingCRUD(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('binding-relay', 'binding-relay', '绑定线路机', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('binding-landing-1', 'binding-landing-1', '绑定落地机 1', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('binding-landing-2', 'binding-landing-2', '绑定落地机 2', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('binding-route', '绑定测试线路', 'binding-relay', 'binding-landing-1', 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, family, created_at, updated_at) VALUES ('binding-exit-1', 'binding-landing-1', '198.51.100.51', 4, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, family, created_at, updated_at) VALUES ('binding-exit-2', 'binding-landing-2', '198.51.100.52', 4, ?, ?)`, []any{nowText, nowText}},
	}
	for _, statement := range seed {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed route binding data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	bound := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes/binding-route/exit-ips", token, map[string]any{
		"exitIpId": "binding-exit-1", "allocationWeight": 2.5,
	})
	if bound["code"] != successCode {
		t.Fatalf("bind route exit IP response = %#v", bound)
	}
	boundData := bound["data"].(map[string]any)
	bindingID := boundData["id"].(string)
	if boundData["address"] != "198.51.100.51" || boundData["landingNodeName"] != "绑定落地机 1" || boundData["allocationWeight"] != 2.5 || boundData["enabled"] != true {
		t.Fatalf("bound route exit IP data = %#v", boundData)
	}

	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/routes/binding-route/exit-ips", token, nil)
	if list["code"] != successCode || len(list["data"].([]any)) != 1 {
		t.Fatalf("route exit IP list response = %#v", list)
	}
	duplicateStatus, duplicate := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes/binding-route/exit-ips", token, map[string]any{"exitIpId": "binding-exit-1"})
	if duplicateStatus != http.StatusConflict || duplicate["code"] != validationCode {
		t.Fatalf("duplicate route exit IP binding status=%d response=%#v", duplicateStatus, duplicate)
	}
	mismatchStatus, mismatch := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes/binding-route/exit-ips", token, map[string]any{"exitIpId": "binding-exit-2"})
	if mismatchStatus != http.StatusBadRequest || mismatch["code"] != validationCode {
		t.Fatalf("cross-landing route exit IP binding status=%d response=%#v", mismatchStatus, mismatch)
	}

	updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/routes/binding-route/exit-ips/binding-exit-1", token, map[string]any{
		"allocationWeight": 3, "enabled": false,
	})
	if updated["code"] != successCode || updated["data"].(map[string]any)["allocationWeight"] != float64(3) || updated["data"].(map[string]any)["enabled"] != false {
		t.Fatalf("update route exit IP response = %#v", updated)
	}
	routeDetail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/routes/binding-route", token, nil)
	if routeDetail["code"] != successCode || routeDetail["data"].(map[string]any)["exitIpCount"] != float64(0) {
		t.Fatalf("disabled route exit IP count response = %#v", routeDetail)
	}

	unboundStatus, unbound := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/routes/binding-route/exit-ips/binding-exit-1", token, nil)
	if unboundStatus != http.StatusOK || unbound["code"] != successCode || unbound["data"].(map[string]any)["id"] != bindingID {
		t.Fatalf("unbind route exit IP status=%d response=%#v", unboundStatus, unbound)
	}
	list = doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/routes/binding-route/exit-ips", token, nil)
	if list["code"] != successCode || len(list["data"].([]any)) != 0 {
		t.Fatalf("route exit IP list after unbind = %#v", list)
	}
}

func TestAllocationCountsOnlyEffectiveUsersAndUsableBindings(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	futureExpiry := now.Add(24 * time.Hour).Format(time.RFC3339Nano)
	pastExpiry := now.Add(-24 * time.Hour).Format(time.RFC3339Nano)
	futureCreated := now.Add(24 * time.Hour).Format(time.RFC3339Nano)
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('allocation-relay', 'allocation-relay', '归属线路机', 'relay', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('allocation-landing', 'allocation-landing', '归属落地机', 'landing', 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('allocation-route', '有效线路', 'allocation-relay', 'allocation-landing', 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('allocation-disabled-route', '停用线路', 'allocation-relay', 'allocation-landing', 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, family, enabled, created_at, updated_at) VALUES ('allocation-exit', 'allocation-landing', '198.51.100.60', 4, 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO exit_ips (id, landing_node_id, ip, family, enabled, created_at, updated_at) VALUES ('allocation-disabled-exit', 'allocation-landing', '198.51.100.61', 4, 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('allocation-valid', '有效用户', 'active', ?, ?, ?)`, []any{futureExpiry, nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('allocation-unlimited', '长期用户', 'active', NULL, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('allocation-expired', '过期用户', 'expired', ?, ?, ?)`, []any{pastExpiry, nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('allocation-disabled', '停用用户', 'disabled', ?, ?, ?)`, []any{futureExpiry, nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('allocation-future', '未来用户', 'active', ?, ?, ?)`, []any{futureExpiry, futureCreated, futureCreated}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-valid', 'allocation-valid', 'allocation-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-unlimited', 'allocation-unlimited', 'allocation-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-expired', 'allocation-expired', 'allocation-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-disabled', 'allocation-disabled', 'allocation-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-future', 'allocation-future', 'allocation-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO user_routes (id, user_id, route_id, is_primary, active_from) VALUES ('allocation-map-disabled-route', 'allocation-valid', 'allocation-disabled-route', 1, ?)`, []any{nowText}},
		{`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, enabled) VALUES ('allocation-binding', 'allocation-route', 'allocation-exit', 1)`, nil},
		{`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, enabled) VALUES ('allocation-disabled-binding', 'allocation-route', 'allocation-disabled-exit', 1)`, nil},
		{`INSERT INTO route_exit_ips (id, route_id, exit_ip_id, enabled) VALUES ('allocation-disabled-route-binding', 'allocation-disabled-route', 'allocation-exit', 1)`, nil},
		{`INSERT INTO user_paths (id, user_id, relay_node_id, landing_node_id, exit_ip_id, mode, active_from, created_at, updated_at) VALUES ('allocation-path-valid', 'allocation-valid', 'allocation-relay', 'allocation-landing', 'allocation-exit', 'landing', ?, ?, ?)`, []any{nowText, nowText, nowText}},
		{`INSERT INTO user_paths (id, user_id, relay_node_id, landing_node_id, exit_ip_id, mode, active_from, created_at, updated_at) VALUES ('allocation-path-unlimited', 'allocation-unlimited', 'allocation-relay', 'allocation-landing', 'allocation-exit', 'landing', ?, ?, ?)`, []any{nowText, nowText, nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed allocation data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	routeDetail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/routes/allocation-route", token, nil)
	if routeDetail["code"] != successCode {
		t.Fatalf("allocation route detail = %#v", routeDetail)
	}
	routeData := routeDetail["data"].(map[string]any)
	if routeData["allocatedUserCount"] != float64(2) || routeData["exitIpCount"] != float64(1) {
		t.Fatalf("allocation route counts = %#v, want users=2 exitIps=1", routeData)
	}

	exitDetail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/allocation-exit", token, nil)
	if exitDetail["code"] != successCode || exitDetail["data"].(map[string]any)["allocatedUserCount"] != float64(2) {
		t.Fatalf("allocation exit IP detail = %#v, want users=2", exitDetail)
	}
	disabledExit := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/allocation-disabled-exit", token, nil)
	if disabledExit["code"] != successCode || disabledExit["data"].(map[string]any)["allocatedUserCount"] != float64(0) {
		t.Fatalf("disabled exit IP allocation = %#v, want users=0", disabledExit)
	}
}

func TestNodeAdminRegistrationAndToggle(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	createdStatus, created := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"nodeKey": "admin-landing-1", "name": "管理创建落地机", "type": "landing",
		"hostname": "landing.example", "publicIp": "203.0.113.80", "region": "东京", "provider": "Test ISP", "panelBasePath": "/panel",
	})
	if createdStatus != http.StatusOK || created["code"] != successCode {
		t.Fatalf("create node status=%d response=%#v", createdStatus, created)
	}
	createdData := created["data"].(map[string]any)
	nodeID := createdData["nodeId"].(string)
	nodeToken := createdData["token"].(string)
	installerToken, _ := createdData["installerToken"].(string)
	if nodeID == "" || nodeToken == "" || installerToken == "" || createdData["nodeKey"] != "admin-landing-1" || createdData["type"] != "landing" {
		t.Fatalf("unexpected node registration data: %#v", createdData)
	}
	installerExpiry, _ := createdData["installerTokenExpiresAt"].(string)
	if _, err := time.Parse(time.RFC3339Nano, installerExpiry); err != nil {
		t.Fatalf("invalid installer token expiry: %q", installerExpiry)
	}

	bootstrapStatus, bootstrap := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/bootstrap", "", map[string]any{
		"install_token": installerToken, "hostname": "landing-host", "agent_version": "0.2.0",
	})
	if bootstrapStatus != http.StatusOK || bootstrap["code"] != successCode {
		t.Fatalf("Agent bootstrap status=%d response=%#v", bootstrapStatus, bootstrap)
	}
	bootstrapData := bootstrap["data"].(map[string]any)
	bootstrappedToken, _ := bootstrapData["central_token"].(string)
	if bootstrappedToken == "" || bootstrapData["node_key"] != "admin-landing-1" || bootstrapData["node_type"] != "landing" {
		t.Fatalf("unexpected Agent bootstrap response: %#v", bootstrapData)
	}
	secondBootstrapStatus, secondBootstrap := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/bootstrap", "", map[string]any{
		"install_token": installerToken,
	})
	if secondBootstrapStatus != http.StatusUnauthorized || secondBootstrap["code"] != unauthorizedCode {
		t.Fatalf("reused installer token status=%d response=%#v", secondBootstrapStatus, secondBootstrap)
	}

	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes?page_size=20", token, nil)
	items := list["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("node list after registration = %#v", list)
	}
	listed := items[0].(map[string]any)
	if listed["id"] != nodeID || listed["enabled"] != true || listed["status"] != "unknown" {
		t.Fatalf("unexpected listed node: %#v", listed)
	}

	duplicateStatus, duplicate := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"nodeKey": "admin-landing-1", "name": "重复节点", "type": "landing",
	})
	if duplicateStatus != http.StatusConflict || duplicate["code"] != validationCode {
		t.Fatalf("duplicate node status=%d response=%#v", duplicateStatus, duplicate)
	}

	disabledStatus, disabled := doJSONWithStatus(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID, token, map[string]any{"enabled": false})
	if disabledStatus != http.StatusOK || disabled["code"] != successCode || disabled["data"].(map[string]any)["enabled"] != false || disabled["data"].(map[string]any)["status"] != "disabled" {
		t.Fatalf("disable node status=%d response=%#v", disabledStatus, disabled)
	}

	disabledSyncStatus, disabledSync := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes/"+nodeID+"/sync", token, nil)
	if disabledSyncStatus != http.StatusConflict || disabledSync["code"] != validationCode {
		t.Fatalf("disabled node sync status=%d response=%#v", disabledSyncStatus, disabledSync)
	}
	disabledDetail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID, token, nil)
	if disabledDetail["code"] != successCode || disabledDetail["data"].(map[string]any)["status"] != "disabled" {
		t.Fatalf("disabled node detail response=%#v", disabledDetail)
	}

	enabledStatus, enabled := doJSONWithStatus(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID, token, map[string]any{"enabled": true, "region": "大阪"})
	if enabledStatus != http.StatusOK || enabled["code"] != successCode || enabled["data"].(map[string]any)["enabled"] != true || enabled["data"].(map[string]any)["region"] != "大阪" {
		t.Fatalf("enable node status=%d response=%#v", enabledStatus, enabled)
	}

	heartbeat := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", bootstrappedToken, map[string]any{
		"node_key": "admin-landing-1", "observed_at": time.Now().UTC().Format(time.RFC3339), "status": map[string]any{"xray_running": true},
	})
	if heartbeat["code"] != successCode {
		t.Fatalf("issued node token did not authenticate: %#v", heartbeat)
	}

	var credentialID string
	if err := database.QueryRow(`SELECT id FROM node_credentials WHERE node_id = ?`, nodeID).Scan(&credentialID); err != nil || credentialID == "" {
		t.Fatalf("read issued node credential id=%q err=%v", credentialID, err)
	}
}

func TestNodeManagementURLIsCanonicalAndEditable(t *testing.T) {
	server, _ := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"nodeKey": "management-url-node", "name": "管理地址节点", "type": "relay",
		"managementUrl": "https://node.example:18086/Alien/",
	})
	if created["code"] != successCode {
		t.Fatalf("create node response = %#v", created)
	}
	nodeID := created["data"].(map[string]any)["nodeId"].(string)

	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes?keyword=node.example", token, nil)
	items := list["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("management URL search returned %#v", list)
	}
	listed := items[0].(map[string]any)
	if listed["managementUrl"] != "https://node.example:18086/Alien" || listed["host"] != "https://node.example:18086/Alien" {
		t.Fatalf("unexpected canonical management URL in list: %#v", listed)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID, token, nil)
	detailData := detail["data"].(map[string]any)
	if detailData["managementUrl"] != "https://node.example:18086/Alien" || detailData["host"] != "https://node.example:18086/Alien" || detailData["panelBasePath"] != "/Alien" {
		t.Fatalf("unexpected canonical management URL in detail: %#v", detailData)
	}

	updatedStatus, updated := doJSONWithStatus(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID, token, map[string]any{
		"managementUrl": "http://new-node.example:2053/Panel/",
	})
	if updatedStatus != http.StatusOK || updated["code"] != successCode || updated["data"].(map[string]any)["managementUrl"] != "http://new-node.example:2053/Panel" || updated["data"].(map[string]any)["panelBasePath"] != "/Panel" {
		t.Fatalf("update management URL status=%d response=%#v", updatedStatus, updated)
	}

	invalidStatus, invalid := doJSONWithStatus(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID, token, map[string]any{
		"managementUrl": "https://user:pass@node.example/Panel",
	})
	if invalidStatus != http.StatusBadRequest || invalid["code"] != validationCode {
		t.Fatalf("invalid management URL status=%d response=%#v", invalidStatus, invalid)
	}
}

func TestNodeAdminGeneratesNodeKeyWhenOmitted(t *testing.T) {
	server, _ := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"name": "自动 Key 节点", "type": "relay",
	})
	if created["code"] != successCode {
		t.Fatalf("create node response = %#v", created)
	}
	nodeKey, _ := created["data"].(map[string]any)["nodeKey"].(string)
	if !strings.HasPrefix(nodeKey, "node-") || len(nodeKey) != len("node-")+16 {
		t.Fatalf("generated node key = %q", nodeKey)
	}
}

func TestNodeAdminRotatesInstallerToken(t *testing.T) {
	server, _ := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	adminToken := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", adminToken, map[string]any{
		"name": "安装 Token 轮换节点", "type": "relay",
	})
	if created["code"] != successCode {
		t.Fatalf("create node response = %#v", created)
	}
	createdData := created["data"].(map[string]any)
	nodeID := createdData["nodeId"].(string)
	firstToken := createdData["installerToken"].(string)
	rotatedStatus, rotated := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes/"+nodeID+"/install-token", adminToken, nil)
	if rotatedStatus != http.StatusOK || rotated["code"] != successCode {
		t.Fatalf("rotate installer token status=%d response=%#v", rotatedStatus, rotated)
	}
	rotatedData := rotated["data"].(map[string]any)
	secondToken := rotatedData["installerToken"].(string)
	if secondToken == "" || secondToken == firstToken || rotatedData["nodeId"] != nodeID {
		t.Fatalf("unexpected rotated installer token: %#v", rotatedData)
	}

	oldStatus, oldResponse := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/bootstrap", "", map[string]any{"install_token": firstToken})
	if oldStatus != http.StatusUnauthorized || oldResponse["code"] != unauthorizedCode {
		t.Fatalf("rotated installer token remained usable status=%d response=%#v", oldStatus, oldResponse)
	}
	newStatus, newResponse := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/bootstrap", "", map[string]any{"install_token": secondToken})
	if newStatus != http.StatusOK || newResponse["code"] != successCode {
		t.Fatalf("new installer token status=%d response=%#v", newStatus, newResponse)
	}
	reusedStatus, reusedResponse := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/bootstrap", "", map[string]any{"install_token": secondToken})
	if reusedStatus != http.StatusUnauthorized || reusedResponse["code"] != unauthorizedCode {
		t.Fatalf("rotated installer token reused status=%d response=%#v", reusedStatus, reusedResponse)
	}
}

func TestNodeAdminCreatesMultipleExitIPsAndDeletesNode(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"nodeKey": "multi-exit-node", "name": "多出口落地机", "type": "landing",
		"publicIp": "203.0.113.10", "exitIps": []string{"203.0.113.10", "198.51.100.20", "2001:db8::20"},
	})
	if created["code"] != successCode {
		t.Fatalf("create multi-exit node response = %#v", created)
	}
	createdData := created["data"].(map[string]any)
	nodeID := createdData["nodeId"].(string)
	if createdData["exitIpCount"] != float64(3) {
		t.Fatalf("created exit IP count = %#v", createdData["exitIpCount"])
	}
	var exitIPCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM exit_ips WHERE owner_node_id = ?`, nodeID).Scan(&exitIPCount); err != nil || exitIPCount != 3 {
		t.Fatalf("stored exit IP count=%d err=%v", exitIPCount, err)
	}
	legacyNow := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO routes (id, name, relay_node_id, landing_node_id, created_at, updated_at) VALUES ('legacy-node-route', '旧线路记录', ?, ?, ?, ?)`, nodeID, nodeID, legacyNow, legacyNow); err != nil {
		t.Fatalf("insert legacy route: %v", err)
	}

	deletedStatus, deleted := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/nodes/"+nodeID, token, nil)
	if deletedStatus != http.StatusOK || deleted["code"] != successCode || deleted["data"].(map[string]any)["deleted"] != true {
		t.Fatalf("delete node status=%d response=%#v", deletedStatus, deleted)
	}
	var remaining int
	if err := database.QueryRow(`SELECT COUNT(*) FROM exit_ips WHERE owner_node_id = ? OR landing_node_id = ?`, nodeID, nodeID).Scan(&remaining); err != nil || remaining != 0 {
		t.Fatalf("deleted node exit IPs remaining=%d err=%v", remaining, err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM routes WHERE relay_node_id = ? OR landing_node_id = ?`, nodeID, nodeID).Scan(&remaining); err != nil || remaining != 0 {
		t.Fatalf("deleted node legacy routes remaining=%d err=%v", remaining, err)
	}
	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes?page_size=20", token, nil)
	if list["code"] != successCode || len(list["data"].(map[string]any)["items"].([]any)) != 0 {
		t.Fatalf("deleted node remained in list: %#v", list)
	}
	detailStatus, detail := doJSONWithStatus(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID, token, nil)
	if detailStatus != http.StatusNotFound || detail["code"] != notFoundCode {
		t.Fatalf("deleted node detail status=%d response=%#v", detailStatus, detail)
	}

	heartbeatStatus, heartbeat := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", createdData["token"].(string), map[string]any{
		"node_key": "multi-exit-node", "observed_at": time.Now().UTC().Format(time.RFC3339), "status": map[string]any{"xray_running": true},
	})
	if heartbeatStatus != http.StatusUnauthorized || heartbeat["code"] != unauthorizedCode {
		t.Fatalf("deleted node token status=%d response=%#v", heartbeatStatus, heartbeat)
	}
}

func TestNodeCostCRUDAndHistoryProtection(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"nodeKey": "cost-node", "name": "成本测试节点", "type": "relay",
	})
	if created["code"] != successCode {
		t.Fatalf("create node response = %#v", created)
	}
	nodeID := created["data"].(map[string]any)["nodeId"].(string)

	cost := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes/"+nodeID+"/costs", token, map[string]any{
		"category": "服务器", "monthlyAmount": 45.5, "effectiveFrom": "2026-09-01", "notes": "主机月租",
	})
	if cost["code"] != successCode {
		t.Fatalf("create node cost response = %#v", cost)
	}
	costData := cost["data"].(map[string]any)
	costID := costData["id"].(string)
	if costData["nodeId"] != nodeID || costData["monthlyAmount"] != 45.5 || costData["effectiveFrom"] != "2026-09-01" {
		t.Fatalf("unexpected node cost = %#v", costData)
	}
	costList := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID+"/costs", token, nil)
	if costList["code"] != successCode || len(costList["data"].([]any)) != 1 {
		t.Fatalf("node cost list response = %#v", costList)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID, token, nil)
	if detail["code"] != successCode {
		t.Fatalf("node detail response = %#v", detail)
	}
	costs := detail["data"].(map[string]any)["costs"].([]any)
	if len(costs) != 1 || costs[0].(map[string]any)["id"] != costID {
		t.Fatalf("node detail costs = %#v", costs)
	}

	updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID+"/costs/"+costID, token, map[string]any{
		"monthlyAmount": 50, "notes": "主机月租已调整",
	})
	if updated["code"] != successCode || updated["data"].(map[string]any)["monthlyAmount"] != float64(50) {
		t.Fatalf("update node cost response = %#v", updated)
	}

	dateChangeStatus, dateChange := doJSONWithStatus(t, ts.Client(), http.MethodPatch, ts.URL+"/api/nodes/"+nodeID+"/costs/"+costID, token, map[string]any{
		"effectiveFrom": "2026-10-01",
	})
	if dateChangeStatus != http.StatusConflict || dateChange["code"] != validationCode {
		t.Fatalf("date change status=%d response=%#v", dateChangeStatus, dateChange)
	}

	finance := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-09", token, nil)
	if finance["code"] != successCode || finance["data"].(map[string]any)["monthCost"] != float64(50) {
		t.Fatalf("finance node cost response = %#v", finance)
	}
	var storedAmount float64
	if err := database.QueryRow(`SELECT monthly_amount FROM node_costs WHERE id = ?`, costID).Scan(&storedAmount); err != nil || storedAmount != 50 {
		t.Fatalf("stored node cost amount=%v err=%v", storedAmount, err)
	}
}

func TestFinanceUsesHistoricalEffectiveUsers(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	// Keep the fixtures independent of the wall clock. The request deliberately
	// asks for historical months while the maintenance refresh moves these
	// users' current operational statuses to expired where appropriate.
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, expiry_time, created_at, updated_at) VALUES ('revenue-history', '历史用户', 'active', 100, 'CNY', '2026-02-01T00:00:00Z', '2025-12-15T00:00:00Z', '2025-12-15T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, expiry_time, created_at, updated_at) VALUES ('revenue-before', '此前已过期', 'expired', 80, 'CNY', '2025-12-31T23:59:59Z', '2025-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, created_at, updated_at) VALUES ('revenue-unlimited', '长期用户', 'active', 40, 'CNY', '2026-01-15T00:00:00Z', '2026-01-15T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, created_at, updated_at) VALUES ('revenue-created-later', '后来创建', 'active', 70, 'CNY', '2026-02-01T00:00:00Z', '2026-02-01T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, expiry_time, created_at, updated_at) VALUES ('revenue-disabled', '停用用户', 'disabled', 90, 'CNY', '2026-12-31T00:00:00Z', '2025-12-01T00:00:00Z', '2025-12-01T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, expiry_time, created_at, updated_at) VALUES ('revenue-usd', '非人民币用户', 'active', 60, 'USD', '2026-12-31T00:00:00Z', '2025-12-01T00:00:00Z', '2025-12-01T00:00:00Z')`, nil},
		{`INSERT INTO users (id, display_name, status, monthly_fee, currency, expiry_time, created_at, updated_at) VALUES ('revenue-start-boundary', '起始日到期', 'active', 5, 'CNY', '2026-01-01T00:00:00Z', '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z')`, nil},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed revenue data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)

	january := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-01", token, nil)
	if january["code"] != successCode {
		t.Fatalf("january finance response = %#v", january)
	}
	januaryData := january["data"].(map[string]any)
	if januaryData["effectiveUserCount"] != float64(3) || januaryData["monthIncome"] != float64(145) {
		t.Fatalf("january effective revenue = %#v, want count=3 income=145", januaryData)
	}

	february := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-02", token, nil)
	februaryData := february["data"].(map[string]any)
	if february["code"] != successCode || februaryData["effectiveUserCount"] != float64(3) || februaryData["monthIncome"] != float64(210) {
		t.Fatalf("february effective revenue = %#v, want count=3 income=210", februaryData)
	}

	march := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-03", token, nil)
	marchData := march["data"].(map[string]any)
	if march["code"] != successCode || marchData["effectiveUserCount"] != float64(2) || marchData["monthIncome"] != float64(110) {
		t.Fatalf("march effective revenue = %#v, want count=2 income=110", marchData)
	}
}

func TestFinanceAggregatesTemporalCostsAndGrossProfit(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := "2026-09-03T00:00:00Z"
	statements := []string{
		`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('finance-relay', 'finance-relay', '财务线路机', 'relay', 'online', '` + now + `', '` + now + `')`,
		`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('finance-landing', 'finance-landing', '财务落地机', 'landing', 'online', '` + now + `', '` + now + `')`,
		`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-node-current', 'finance-relay', '主机', 30, 'CNY', '2026-09-01', NULL, '` + now + `')`,
		`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-node-next', 'finance-relay', '下月主机', 40, 'CNY', '2026-10-01', NULL, '` + now + `')`,
		`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-node-ended', 'finance-relay', '已结束', 20, 'CNY', '2026-08-01', '2026-08-31', '` + now + `')`,
		`INSERT INTO node_costs (id, node_id, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-node-boundary', 'finance-relay', '月初结束', 10, 'CNY', '2026-08-15', '2026-09-01', '` + now + `')`,
		`INSERT INTO other_costs (id, name, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-other-current', '监控', '服务', 5, 'CNY', '2026-09-15', NULL, '` + now + `')`,
		`INSERT INTO other_costs (id, name, category, monthly_amount, currency, effective_from, effective_to, created_at) VALUES ('finance-other-next', '下月服务', '服务', 50, 'CNY', '2026-10-01', NULL, '` + now + `')`,
		`INSERT INTO exit_ips (id, landing_node_id, ip, family, monthly_cost, currency, enabled, valid_from, valid_to, created_at, updated_at) VALUES ('finance-exit-current', 'finance-landing', '203.0.113.20', 4, 7, 'CNY', 1, '2026-09-01', '2026-09-30', '` + now + `', '` + now + `')`,
		`INSERT INTO exit_ips (id, landing_node_id, ip, family, monthly_cost, currency, enabled, valid_from, valid_to, created_at, updated_at) VALUES ('finance-exit-next', 'finance-landing', '203.0.113.21', 4, 9, 'CNY', 1, '2026-10-01', NULL, '` + now + `', '` + now + `')`,
		`INSERT INTO exit_ips (id, landing_node_id, ip, family, monthly_cost, currency, enabled, valid_from, valid_to, created_at, updated_at) VALUES ('finance-exit-ended', 'finance-landing', '203.0.113.22', 4, 11, 'CNY', 1, NULL, '2026-08-31', '` + now + `', '` + now + `')`,
		`INSERT INTO exit_ips (id, landing_node_id, ip, family, monthly_cost, currency, enabled, valid_from, valid_to, created_at, updated_at) VALUES ('finance-exit-disabled', 'finance-landing', '203.0.113.23', 4, 13, 'CNY', 0, '2026-09-01', NULL, '` + now + `', '` + now + `')`,
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("seed temporal cost data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	finance := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-09", token, nil)
	if finance["code"] != successCode {
		t.Fatalf("temporal finance response = %#v", finance)
	}
	data := finance["data"].(map[string]any)
	if data["monthIncome"] != float64(0) || data["monthCost"] != float64(52) || data["grossProfit"] != float64(-52) {
		t.Fatalf("temporal finance totals = %#v, want income=0 cost=52 grossProfit=-52", data)
	}
	breakdown := data["breakdown"].([]any)
	if len(breakdown) != 4 || breakdown[0].(map[string]any)["label"] != "用户月费收入" {
		t.Fatalf("temporal finance breakdown = %#v", breakdown)
	}
}

func doJSON(t *testing.T, client *http.Client, method, url, token string, payload any) map[string]any {
	t.Helper()
	return doJSONWithRequestID(t, client, method, url, token, "", payload)
}

func doJSONWithStatus(t *testing.T, client *http.Client, method, url, token string, payload any) (int, map[string]any) {
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
	return response.StatusCode, result
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
