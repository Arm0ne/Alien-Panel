package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserListUsesLatestClientOnlineTime(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Truncate(time.Second)
	syncAt := now.Add(-2 * time.Hour).Format(time.RFC3339Nano)
	firstOnlineAt := now.Add(-90 * time.Minute).Format(time.RFC3339Nano)
	latestOnlineAt := now.Add(-15 * time.Minute).Format(time.RFC3339Nano)

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('activity-relay', 'activity-relay', '活动线路机', 'relay', 'online', ?, ?)`, []any{syncAt, syncAt}},
		{`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('activity-user', '活动用户', 'active', ?, ?)`, []any{syncAt, syncAt}},
		{`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, enable, first_seen_at, last_seen_at) VALUES ('activity-inbound', 'activity-relay', '1', 'activity-user', 'user', 'activity-user', 1, ?, ?)`, []any{syncAt, syncAt}},
		{`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('activity-map', 'activity-user', 'activity-inbound', 1, ?)`, []any{syncAt}},
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, last_online, last_seen_at) VALUES ('activity-client-old', 'activity-relay', 'activity-inbound', 'old', 'activity@example.com', 1, ?, ?)`, []any{firstOnlineAt, syncAt}},
		{`INSERT INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, last_online, last_seen_at) VALUES ('activity-client-latest', 'activity-relay', 'activity-inbound', 'latest', 'activity@example.com', 1, ?, ?)`, []any{latestOnlineAt, syncAt}},
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("seed user activity data: %v", err)
		}
	}

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	response := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?page_size=20", token, nil)
	if response["code"] != successCode {
		t.Fatalf("user list response = %#v", response)
	}
	items := response["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("user list items = %#v", items)
	}
	item := items[0].(map[string]any)
	if item["lastActivityAt"] != latestOnlineAt {
		t.Fatalf("last activity = %v, want latest Client online time %s", item["lastActivityAt"], latestOnlineAt)
	}
}
