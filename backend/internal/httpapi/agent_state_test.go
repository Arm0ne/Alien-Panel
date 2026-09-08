package httpapi

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserStateFromClientsUsesEnabledClientExpiry(t *testing.T) {
	now := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	earliest := now.Add(30 * 24 * time.Hour)
	later := now.Add(60 * 24 * time.Hour)

	state := userStateFromClients(agentInboundPayload{
		Enable:     true,
		ExpiryTime: later.Add(120 * 24 * time.Hour).Unix(), // Inbound expiry is not business expiry.
		Clients: []agentClientPayload{
			{RemoteID: "phone", Enable: true, ExpiryTime: earliest.Unix()},
			{RemoteID: "laptop", Enable: true, ExpiryTime: earliest.Unix()},
			{RemoteID: "old-device", Enable: false, ExpiryTime: now.Add(-time.Hour).Unix()},
		},
	}, now)
	if state.Status != "active" || state.ExpiryText != earliest.Format(time.RFC3339Nano) || state.ExpiryMismatch {
		t.Fatalf("same enabled Client expiries state = %#v", state)
	}

	state = userStateFromClients(agentInboundPayload{
		Enable: true,
		Clients: []agentClientPayload{
			{RemoteID: "phone", Enable: true, ExpiryTime: earliest.Unix()},
			{RemoteID: "laptop", Enable: true, ExpiryTime: later.Unix()},
			{RemoteID: "tablet", Enable: true},
		},
	}, now)
	if state.Status != "active" || state.ExpiryText != earliest.Format(time.RFC3339Nano) || !state.ExpiryMismatch {
		t.Fatalf("different enabled Client expiries state = %#v", state)
	}

	for name, inbound := range map[string]agentInboundPayload{
		"no enabled client": {Enable: true, Clients: []agentClientPayload{{RemoteID: "phone", Enable: false, ExpiryTime: earliest.Unix()}}},
		"disabled inbound":  {Enable: false, Clients: []agentClientPayload{{RemoteID: "phone", Enable: true, ExpiryTime: earliest.Unix()}}},
	} {
		state = userStateFromClients(inbound, now)
		if state.Status != "disabled" || state.ExpiryText != "" {
			t.Fatalf("%s state = %#v", name, state)
		}
	}
}

func TestAgentSyncKeepsInboundUsersSeparateAndReportsExpiryMismatch(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	const nodeID = "client-expiry-node"
	const nodeKey = "client-expiry-node"
	const token = "client-expiry-token"
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES (?, ?, 'Client expiry node', 'relay', 'online', ?, ?)`, nodeID, nodeKey, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('client-expiry-credential', ?, ?, ?, ?)`, nodeID, hashToken(token), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	earliest := now.AddDate(0, 1, 0).Unix()
	later := now.AddDate(0, 2, 0).Unix()
	payload := map[string]any{
		"node_key": nodeKey, "sync_id": "client-expiry-sync-1", "observed_at": nowText,
		"status": map[string]any{"xray_running": true},
		"inbounds": []any{
			map[string]any{
				"remote_id": 1, "tag": "customer-one", "enable": true, "expiry_time": now.AddDate(1, 0, 0).Unix(),
				"clients": []any{
					map[string]any{"remote_id": "one-phone", "email": "same@example.com", "enable": true, "expiry_time": earliest},
					map[string]any{"remote_id": "one-laptop", "email": "same@example.com", "enable": true, "expiry_time": later},
				},
			},
			map[string]any{
				"remote_id": 2, "tag": "customer-two", "enable": true, "expiry_time": now.AddDate(2, 0, 0).Unix(),
				"clients": []any{
					map[string]any{"remote_id": "two-phone", "email": "same@example.com", "enable": true, "expiry_time": later},
				},
			},
		},
	}
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", token, payload); response["code"] != successCode {
		t.Fatalf("first sync response = %#v", response)
	}

	rows, err := database.Query(`SELECT i.remote_inbound_id, u.expiry_time FROM inbounds i JOIN users u ON u.id = i.user_id WHERE i.node_id = ? ORDER BY i.remote_inbound_id`, nodeID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var records []struct{ inboundID, expiry string }
	for rows.Next() {
		var record struct{ inboundID, expiry string }
		if err := rows.Scan(&record.inboundID, &record.expiry); err != nil {
			t.Fatal(err)
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || records[0].inboundID != "1" || records[0].expiry != time.Unix(earliest, 0).UTC().Format(time.RFC3339Nano) || records[1].inboundID != "2" || records[1].expiry != time.Unix(later, 0).UTC().Format(time.RFC3339Nano) {
		t.Fatalf("inbound business users = %#v", records)
	}
	var mismatchEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_expiry_mismatch' AND event_status = 'open'`).Scan(&mismatchEvents); err != nil {
		t.Fatal(err)
	}
	if mismatchEvents != 1 {
		t.Fatalf("open mismatch events = %d, want 1", mismatchEvents)
	}

	payload["sync_id"] = "client-expiry-sync-2"
	payload["observed_at"] = now.Add(time.Minute).Format(time.RFC3339Nano)
	payload["inbounds"].([]any)[0].(map[string]any)["clients"].([]any)[1].(map[string]any)["expiry_time"] = earliest
	if response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", token, payload); response["code"] != successCode {
		t.Fatalf("consistent-client sync response = %#v", response)
	}
	var resolvedEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_expiry_mismatch' AND event_status = 'resolved'`).Scan(&resolvedEvents); err != nil {
		t.Fatal(err)
	}
	if resolvedEvents != 1 {
		t.Fatalf("resolved mismatch events = %d, want 1", resolvedEvents)
	}
}

func TestLivenessHeartbeatDoesNotHideXPanelSyncFailure(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	nowText := now.Format(time.RFC3339Nano)
	const nodeID = "liveness-node"
	const nodeKey = "liveness-node"
	const token = "liveness-token"
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES (?, ?, 'Liveness node', 'relay', 'unknown', ?, ?)`, nodeID, nodeKey, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('liveness-credential', ?, ?, ?, ?)`, nodeID, hashToken(token), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	failure := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/heartbeat", token, map[string]any{
		"node_key": nodeKey, "observed_at": nowText, "status_available": false, "sync_state": "failed",
		"sync_error": "collect inbounds: xpanel GET /inbounds/list returned HTTP 404",
		"status":     map[string]any{"agent_version": "v1.0.2"},
	})
	if failure["code"] != successCode {
		t.Fatalf("liveness heartbeat = %#v", failure)
	}
	var health, syncStatus, agentVersion string
	var syncError sql.NullString
	if err := database.QueryRow(`SELECT health_status, sync_status, last_sync_error, agent_version FROM nodes WHERE id = ?`, nodeID).Scan(&health, &syncStatus, &syncError, &agentVersion); err != nil {
		t.Fatal(err)
	}
	if health != "online" || syncStatus != "failed" || !syncError.Valid || agentVersion != "v1.0.2" {
		t.Fatalf("node after liveness heartbeat health=%q sync=%q error=%#v agent=%q", health, syncStatus, syncError, agentVersion)
	}

	recovery := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", token, map[string]any{
		"node_key": nodeKey, "sync_id": "liveness-recovery", "observed_at": now.Add(time.Minute).Format(time.RFC3339Nano),
		"status": map[string]any{"agent_version": "v1.0.2", "xray_running": true}, "inbounds": []any{},
	})
	if recovery["code"] != successCode {
		t.Fatalf("recovery sync = %#v", recovery)
	}
	if err := database.QueryRow(`SELECT sync_status, last_sync_error FROM nodes WHERE id = ?`, nodeID).Scan(&syncStatus, &syncError); err != nil {
		t.Fatal(err)
	}
	if syncStatus != "success" || syncError.Valid {
		t.Fatalf("node after recovery sync status=%q error=%#v", syncStatus, syncError)
	}
}
