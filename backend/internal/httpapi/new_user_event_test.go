package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRelayInboundCreatesProfileEventAndUpdateResolvesIt(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 9, 18, 1, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at)
VALUES ('new-user-node', 'new-user-node', '线路机 A', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at)
VALUES ('new-user-credential', 'new-user-node', ?, ?, ?)`, hashToken("new-user-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	payload := lifecycleSyncPayload("new-user-node", "new-user-sync-1", now, []map[string]any{
		{"remote_id": "new-client", "email": "new@example.com", "enable": true, "expiry_time": now.Add(30 * 24 * time.Hour).Unix()},
	})
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "new-user-token", payload); response["code"] != successCode {
		t.Fatalf("initial sync response = %#v", response)
	}

	var userID, eventID, eventStatus, actionType, dedupeKey string
	var requiresAction int
	if err := database.QueryRow(`SELECT u.id, e.id, e.event_status, e.requires_action, e.action_type, e.dedupe_key
FROM users u JOIN node_events e ON e.resource_type = 'user' AND e.resource_id = u.id
WHERE e.event_type = 'new_user_profile_required'`).
		Scan(&userID, &eventID, &eventStatus, &requiresAction, &actionType, &dedupeKey); err != nil {
		t.Fatal(err)
	}
	if eventStatus != "open" || requiresAction != 1 || actionType != "complete_user_profile" || dedupeKey != "new-user-profile:"+userID {
		t.Fatalf("new user event id=%q status=%q required=%d action=%q dedupe=%q", eventID, eventStatus, requiresAction, actionType, dedupeKey)
	}

	payload["sync_id"] = "new-user-sync-2"
	payload["observed_at"] = now.Add(time.Minute).Format(time.RFC3339Nano)
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "new-user-token", payload); response["code"] != successCode {
		t.Fatalf("repeat sync response = %#v", response)
	}
	var eventCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'new_user_profile_required' AND resource_id = ?`, userID).Scan(&eventCount); err != nil {
		t.Fatal(err)
	}
	if eventCount != 1 {
		t.Fatalf("new user profile events = %d, want 1", eventCount)
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	update := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/users/"+userID, token, map[string]any{
		"displayName": "正式客户 A", "billingType": "paid", "billingCycle": "monthly", "billingAmount": 100, "currency": "CNY", "notes": "已核对",
	})
	if update["code"] != successCode {
		t.Fatalf("user update response = %#v", update)
	}
	if err := database.QueryRow(`SELECT event_status, requires_action FROM node_events WHERE id = ?`, eventID).Scan(&eventStatus, &requiresAction); err != nil {
		t.Fatal(err)
	}
	if eventStatus != "resolved" || requiresAction != 0 {
		t.Fatalf("resolved profile event status=%q required=%d", eventStatus, requiresAction)
	}
	counts := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events/summary", token, nil)
	if counts["code"] != successCode || counts["data"].(map[string]any)["pendingCount"] != float64(0) {
		t.Fatalf("event summary after profile update = %#v", counts)
	}
}

func TestLandingInboundDoesNotCreateNewUserProfileEvent(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 9, 18, 1, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at)
VALUES ('landing-event-node', 'landing-event-node', '落地机 A', 'landing', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at)
VALUES ('landing-event-credential', 'landing-event-node', ?, ?, ?)`, hashToken("landing-event-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}
	payload := lifecycleSyncPayload("landing-event-node", "landing-event-sync", now, []map[string]any{
		{"remote_id": "landing-client", "enable": true},
	})
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "landing-event-token", payload); response["code"] != successCode {
		t.Fatalf("landing sync response = %#v", response)
	}
	var users, events int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&users); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'new_user_profile_required'`).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if users != 0 || events != 0 {
		t.Fatalf("landing sync users=%d profileEvents=%d, want 0", users, events)
	}
}
