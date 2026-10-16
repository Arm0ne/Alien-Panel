package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserPathAssetsIncludesBoundUnavailableExitIPs(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	pastDate := now.Add(-24 * time.Hour).Format("2006-01-02")
	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('bound-relay', 'bound-relay', '绑定线路机', 'relay', 1, 'online', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('bound-disabled-user', '停用用户', 'disabled', ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, protocol, port, enable, first_seen_at) VALUES ('bound-inbound', 'bound-relay', 'bound-1', 'bound-disabled-user', 'user', 'bound-user', 'vless', 443, 0, ?)`, []any{nowText}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('bound-user-inbound', 'bound-disabled-user', 'bound-inbound', 1, ?)`, []any{nowText}},
		{`INSERT INTO exit_ips (id, source_type, ip, family, enabled, valid_to, created_at, updated_at) VALUES ('bound-expired-exit', 's5', '198.51.100.81', 4, 1, ?, ?, ?)`, []any{pastDate, nowText, nowText}},
		{`INSERT INTO exit_ips (id, source_type, ip, family, enabled, created_at, updated_at) VALUES ('bound-disabled-exit', 's5', '198.51.100.82', 4, 0, ?, ?)`, []any{nowText, nowText}},
		{`INSERT INTO user_paths (id, user_id, relay_node_id, exit_ip_id, mode, active_from, created_at, updated_at) VALUES ('bound-path', 'bound-disabled-user', 'bound-relay', 'bound-expired-exit', 'external', ?, ?, ?)`, []any{nowText, nowText, nowText}},
		{`INSERT INTO user_path_exit_ips (user_path_id, exit_ip_id, position, created_at) VALUES ('bound-path', 'bound-expired-exit', 0, ?)`, []any{nowText}},
		{`INSERT INTO user_path_exit_ips (user_path_id, exit_ip_id, position, created_at) VALUES ('bound-path', 'bound-disabled-exit', 1, ?)`, []any{nowText}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed bound unavailable exit IP data: %v", err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)

	assets := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/bound-disabled-user/path-assets", token, nil)
	if assets["code"] != successCode {
		t.Fatalf("path assets response = %#v", assets)
	}
	external := assets["data"].(map[string]any)["externalExitIps"].([]any)
	byID := make(map[string]map[string]any, len(external))
	for _, raw := range external {
		item := raw.(map[string]any)
		byID[item["id"].(string)] = item
	}
	expired, ok := byID["bound-expired-exit"]
	if !ok || expired["address"] != "198.51.100.81" || expired["status"] != "active" || expired["validTo"] != pastDate || expired["allocatedUserCount"] != float64(1) {
		t.Fatalf("expired bound exit IP asset = %#v", expired)
	}
	disabled, ok := byID["bound-disabled-exit"]
	if !ok || disabled["address"] != "198.51.100.82" || disabled["enabled"] != false || disabled["status"] != "disabled" || disabled["allocatedUserCount"] != float64(1) {
		t.Fatalf("disabled bound exit IP asset = %#v", disabled)
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/bound-disabled-user", token, nil)
	if detail["code"] != successCode {
		t.Fatalf("user detail response = %#v", detail)
	}
	path := detail["data"].(map[string]any)["path"].(map[string]any)
	exitItems := path["exitIps"].([]any)
	if len(exitItems) != 2 {
		t.Fatalf("user detail exit IPs = %#v", path)
	}
	if path["exitIpAddress"] != "198.51.100.81" || path["exitIpValidTo"] != pastDate {
		t.Fatalf("user detail primary exit IP metadata = %#v", path)
	}
	first := exitItems[0].(map[string]any)
	second := exitItems[1].(map[string]any)
	if first["address"] != "198.51.100.81" || first["validTo"] != pastDate || first["enabled"] != true || second["address"] != "198.51.100.82" || second["enabled"] != false {
		t.Fatalf("user detail exit IP metadata = %#v", exitItems)
	}

	users := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/bound-expired-exit/users", token, nil)
	usersData := users["data"].(map[string]any)
	if users["code"] != successCode || usersData["total"] != float64(1) {
		t.Fatalf("bound expired exit IP users = %#v", users)
	}
	userItems := usersData["items"].([]any)
	if len(userItems) != 1 || userItems[0].(map[string]any)["id"] != "bound-disabled-user" || userItems[0].(map[string]any)["status"] != "disabled" {
		t.Fatalf("bound expired exit IP user details = %#v", users)
	}
}
