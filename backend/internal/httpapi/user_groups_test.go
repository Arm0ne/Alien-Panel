package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"xpanel-central/backend/internal/config"
	"xpanel-central/backend/internal/db"
)

func TestUserGroupsExcludeUnassignedUsers(t *testing.T) {
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	server, err := NewServer(config.Config{
		AdminUsername: "admin", AdminPassword: "test-password", SessionTTL: time.Hour,
		CorsOrigins: []string{"http://localhost:9527"},
	}, database, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, sync_status, created_at, updated_at)
VALUES ('relay-1', 'relay-key-1', '香港线路 01', 'relay', 1, 'online', 'success', ?, ?)`, now, now); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, monthly_fee, currency, created_at, updated_at, billing_type)
VALUES ('user-1', '付费用户', 'active', 100, 'CNY', ?, ?, 'paid'),
       ('user-2', '免费用户', 'active', 0, 'CNY', ?, ?, 'free'),
       ('user-3', '未关联用户', 'expired', 100, 'CNY', ?, ?, 'paid')`, now, now, now, now, now, now); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, enable, up, down, client_count, first_seen_at, last_seen_at)
VALUES ('inbound-1', 'relay-1', 'remote-1', 'user', 'user-1', 1, 100, 200, 1, ?, ?),
       ('inbound-2', 'relay-1', 'remote-2', 'user', 'user-2', 1, 300, 400, 2, ?, ?)`, now, now, now, now); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES ('binding-1', 'user-1', 'inbound-1', 1, ?), ('binding-2', 'user-2', 'inbound-2', 1, ?)`, now, now); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	groups := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/groups", token, nil)
	if groups["code"] != successCode {
		t.Fatalf("groups code = %#v", groups["code"])
	}
	data := groups["data"].(map[string]any)
	items := data["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("group count = %d, want relay node only", len(items))
	}
	first := items[0].(map[string]any)
	if first["nodeId"] != "relay-1" || first["nodeName"] != "香港线路 01" {
		t.Fatalf("first group = %#v", first)
	}
	groupStats := first["stats"].(map[string]any)
	if groupStats["total"] != float64(2) || groupStats["active"] != float64(2) || groupStats["paid"] != float64(1) || groupStats["free"] != float64(1) {
		t.Fatalf("node group stats = %#v", groupStats)
	}

	filtered := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?node_id=relay-1&billing_type=free", token, nil)
	filteredData := filtered["data"].(map[string]any)
	if filteredData["total"] != float64(1) {
		t.Fatalf("filtered user count = %#v", filteredData["total"])
	}
}
