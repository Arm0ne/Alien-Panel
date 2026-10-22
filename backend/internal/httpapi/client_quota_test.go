package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestAgentSyncPersistsClientQuotaAndUserAPIExposesIt guards the quota pipeline:
// X-Panel totalGB -> Agent snapshot total_gb -> clients.quota_bytes -> user list.
func TestAgentSyncPersistsClientQuotaAndUserAPIExposesIt(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	const nodeID = "quota-node"
	const nodeKey = "quota-node"
	const agentToken = "quota-agent-token"
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES (?, ?, 'Quota node', 'relay', 'online', ?, ?)`, nodeID, nodeKey, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('quota-credential', ?, ?, ?, ?)`, nodeID, hashToken(agentToken), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	expiry := now.AddDate(1, 0, 0).Unix()
	payload := map[string]any{
		"node_key": nodeKey, "sync_id": "quota-sync-1", "observed_at": nowText,
		"status": map[string]any{"xray_running": true},
		"inbounds": []any{
			map[string]any{
				"remote_id": 11, "tag": "quota-user", "enable": true, "expiry_time": expiry,
				"clients": []any{
					map[string]any{"remote_id": "quota-phone", "email": "quota@example.com", "enable": true, "expiry_time": expiry, "total_gb": 100},
					map[string]any{"remote_id": "quota-laptop", "email": "quota@example.com", "enable": true, "expiry_time": expiry, "total_gb": 2.5},
					map[string]any{"remote_id": "quota-nodata", "email": "quota@example.com", "enable": true, "expiry_time": expiry},
				},
			},
		},
	}
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, payload); response["code"] != successCode {
		t.Fatalf("sync response = %#v", response)
	}

	for clientID, want := range map[string]int64{
		"quota-phone":  100 * (1 << 30),
		"quota-laptop": int64(2.5 * (1 << 30)),
		"quota-nodata": 0,
	} {
		var quota int64
		if err := database.QueryRow(`SELECT quota_bytes FROM clients WHERE remote_client_id = ?`, clientID).Scan(&quota); err != nil {
			t.Fatalf("read %s quota: %v", clientID, err)
		}
		if quota != want {
			t.Fatalf("%s quota = %d, want %d", clientID, quota, want)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin",
		"password": "test-password",
	})
	adminToken := login["data"].(map[string]any)["token"].(string)
	users := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?page=1&page_size=20", adminToken, nil)
	items := users["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("users = %d, want 1", len(items))
	}
	row := items[0].(map[string]any)
	quota, ok := row["quotaBytes"].(float64)
	if !ok {
		t.Fatalf("quotaBytes missing from user list row: %#v", row)
	}
	if int64(quota) != 100*(1<<30) {
		t.Fatalf("user quotaBytes = %v, want %d", quota, int64(100*(1<<30)))
	}
}
