package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEventCenterHidesInternalSyncRequests(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES ('event-node', 'event-node', '事件节点', 'relay', 1, 'online', ?, ?)`, now, now); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at, event_category, title, visibility, requires_action, event_status) VALUES ('internal-sync', 'event-node', 'sync_requested', 'info', 'internal', ?, 'sync', '立即同步请求', 'internal', 0, 'resolved')`, now); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_events (id, node_id, event_type, severity, message, created_at, event_category, title, visibility, requires_action, event_status, action_type) VALUES ('public-error', 'event-node', 'sync_failed', 'error', 'timeout', ?, 'sync', '节点同步失败', 'public', 0, 'open', 'retry_sync')`, now); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events", token, nil)
	if list["code"] != successCode {
		t.Fatalf("events response = %#v", list)
	}
	if list["data"].(map[string]any)["total"] != float64(1) {
		t.Fatalf("visible event total = %#v", list)
	}

	read := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/events/public-error/read", token, nil)
	if read["code"] != successCode {
		t.Fatalf("mark read response = %#v", read)
	}
	resolved := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/events/public-error/resolve", token, nil)
	if resolved["code"] != successCode {
		t.Fatalf("resolve response = %#v", resolved)
	}
	counts := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events/summary", token, nil)
	if counts["code"] != successCode || counts["data"].(map[string]any)["pendingCount"] != float64(0) {
		t.Fatalf("event summary = %#v", counts)
	}
}
