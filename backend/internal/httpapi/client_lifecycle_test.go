package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestAgentSyncPurgesDeletedClientsAndDetectsReplacement(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC().Truncate(time.Second)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('lifecycle-node', 'lifecycle-node', 'Lifecycle Node', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('lifecycle-credential', 'lifecycle-node', ?, ?, ?)`, hashToken("lifecycle-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	initial := lifecycleSyncPayload("lifecycle-node", "lifecycle-1", now, []map[string]any{
		{"remote_id": "old-phone", "email": "old-phone@example.com", "enable": true, "expiry_time": now.Add(30 * 24 * time.Hour).Unix(), "last_online": now.Add(-time.Hour).Unix()},
		{"remote_id": "old-laptop", "email": "old-laptop@example.com", "enable": true, "expiry_time": now.Add(30 * 24 * time.Hour).Unix(), "last_online": now.Add(-2 * time.Hour).Unix()},
	})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "lifecycle-token", initial); result["code"] != successCode {
		t.Fatalf("initial sync response = %#v", result)
	}
	var oldUserID, inboundID string
	if err := database.QueryRow(`SELECT u.id, i.id FROM users u JOIN inbounds i ON i.user_id = u.id WHERE i.remote_inbound_id = '42'`).Scan(&oldUserID, &inboundID); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO user_billing_records (id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, created_at) VALUES ('lifecycle-order', ?, 'monthly', 100, 'CNY', ?, ?, ?, 'confirmed', 'manual', ?)`, oldUserID, nowText, now.Add(30*24*time.Hour).Format(time.RFC3339Nano), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	replacementAt := now.Add(time.Minute)
	replacement := lifecycleSyncPayload("lifecycle-node", "lifecycle-2", replacementAt, []map[string]any{
		{"remote_id": "new-phone", "email": "new-phone@example.com", "enable": true, "expiry_time": replacementAt.Add(30 * 24 * time.Hour).Unix()},
		{"remote_id": "new-laptop", "email": "new-laptop@example.com", "enable": true, "expiry_time": replacementAt.Add(30 * 24 * time.Hour).Unix()},
	})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "lifecycle-token", replacement); result["code"] != successCode {
		t.Fatalf("replacement sync response = %#v", result)
	}
	var oldClients, newClients, replacementEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients WHERE inbound_id = ? AND remote_client_id LIKE 'old-%'`, inboundID).Scan(&oldClients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients WHERE inbound_id = ? AND remote_client_id LIKE 'new-%'`, inboundID).Scan(&newClients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_set_replacement_detected'`).Scan(&replacementEvents); err != nil {
		t.Fatal(err)
	}
	if oldClients != 0 || newClients != 2 || replacementEvents != 1 {
		t.Fatalf("client lifecycle counts old=%d new=%d events=%d", oldClients, newClients, replacementEvents)
	}
	var renewalCandidates int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_renewal_candidates`).Scan(&renewalCandidates); err != nil {
		t.Fatal(err)
	}
	if renewalCandidates != 0 {
		t.Fatalf("replacement created %d renewal candidates, want 0", renewalCandidates)
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	events := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events?status=pending", token, nil)
	items := events["data"].(map[string]any)["items"].([]any)
	var eventID string
	for _, item := range items {
		candidate := item.(map[string]any)
		if candidate["type"] == "client_set_replacement_detected" {
			eventID = candidate["id"].(string)
			break
		}
	}
	if eventID == "" {
		t.Fatalf("replacement event missing from pending events: %#v", items)
	}
	reset := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/events/"+eventID+"/reset-user", token, nil)
	if reset["code"] != successCode {
		t.Fatalf("reset response = %#v", reset)
	}
	resetData := reset["data"].(map[string]any)
	newUserID := resetData["newUserId"].(string)
	if newUserID == "" || newUserID == oldUserID {
		t.Fatalf("replacement user id = %q, old=%q", newUserID, oldUserID)
	}
	var oldProfileEvents, newProfileEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ? AND event_status = 'resolved' AND requires_action = 0`, oldUserID).Scan(&oldProfileEvents); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ? AND event_status = 'open' AND requires_action = 1`, newUserID).Scan(&newProfileEvents); err != nil {
		t.Fatal(err)
	}
	if oldProfileEvents != 1 || newProfileEvents != 1 {
		t.Fatalf("profile events after replacement oldResolved=%d newOpen=%d", oldProfileEvents, newProfileEvents)
	}
	var oldDeleted, oldStatus string
	if err := database.QueryRow(`SELECT COALESCE(deleted_at, ''), status FROM users WHERE id = ?`, oldUserID).Scan(&oldDeleted, &oldStatus); err != nil {
		t.Fatal(err)
	}
	var billingCount, trafficCount, activeMappings int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = ?`, oldUserID).Scan(&billingCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM traffic_snapshots WHERE inbound_id = ?`, inboundID).Scan(&trafficCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_inbounds WHERE user_id = ? AND active_to IS NULL`, newUserID).Scan(&activeMappings); err != nil {
		t.Fatal(err)
	}
	if oldDeleted == "" || oldStatus != "disabled" || billingCount != 1 || trafficCount != 0 || activeMappings != 1 {
		t.Fatalf("reset state deleted=%q status=%q billing=%d traffic=%d mappings=%d", oldDeleted, oldStatus, billingCount, trafficCount, activeMappings)
	}
	var baselineAllTime int64
	if err := database.QueryRow(`SELECT traffic_baseline_all_time FROM inbounds WHERE id = ?`, inboundID).Scan(&baselineAllTime); err != nil {
		t.Fatal(err)
	}
	if baselineAllTime != 100 {
		t.Fatalf("baseline all time = %d, want current inbound total 100", baselineAllTime)
	}
	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/"+newUserID, token, nil)
	if detail["code"] != successCode {
		t.Fatalf("replacement user detail response = %#v", detail)
	}
	detailInbound := detail["data"].(map[string]any)["inbound"].(map[string]any)
	if detailInbound["allTime"] != float64(0) || detailInbound["up"] != float64(0) || detailInbound["down"] != float64(0) {
		t.Fatalf("replacement inbound traffic = %#v, want zeroed baseline", detailInbound)
	}
}

func TestAgentSyncDetectsReplacementAfterEmptyCompleteSync(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC().Truncate(time.Second)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('empty-node', 'empty-node', 'Empty Node', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('empty-credential', 'empty-node', ?, ?, ?)`, hashToken("empty-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}
	initial := lifecycleSyncPayload("empty-node", "empty-1", now, []map[string]any{{"remote_id": "old", "email": "old@example.com", "enable": true}})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "empty-token", initial); result["code"] != successCode {
		t.Fatalf("initial sync response = %#v", result)
	}
	empty := lifecycleSyncPayload("empty-node", "empty-2", now.Add(time.Minute), []map[string]any{})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "empty-token", empty); result["code"] != successCode {
		t.Fatalf("empty sync response = %#v", result)
	}
	var clients, events int
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients`).Scan(&clients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_set_replacement_detected'`).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if clients != 0 || events != 0 {
		t.Fatalf("empty sync state clients=%d events=%d", clients, events)
	}
	newSync := lifecycleSyncPayload("empty-node", "empty-3", now.Add(2*time.Minute), []map[string]any{{"remote_id": "new", "email": "new@example.com", "enable": true}})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "empty-token", newSync); result["code"] != successCode {
		t.Fatalf("new sync response = %#v", result)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_set_replacement_detected'`).Scan(&events); err != nil {
		t.Fatal(err)
	}
	if events != 1 {
		t.Fatalf("replacement events after empty sync = %d", events)
	}
}

func TestAgentSyncCreatesNewUserWhenArchivedInboundReappears(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC().Truncate(time.Second)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('reactivated-node', 'reactivated-node', 'Reactivated Node', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('reactivated-credential', 'reactivated-node', ?, ?, ?)`, hashToken("reactivated-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	initial := lifecycleSyncPayload("reactivated-node", "reactivated-1", now, []map[string]any{
		{"remote_id": "old-client", "email": "old@example.com", "enable": true, "expiry_time": now.Add(30 * 24 * time.Hour).Unix()},
	})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "reactivated-token", initial); result["code"] != successCode {
		t.Fatalf("initial sync response = %#v", result)
	}
	var oldUserID, inboundID string
	if err := database.QueryRow(`SELECT u.id, i.id FROM users u JOIN inbounds i ON i.user_id = u.id WHERE i.node_id = 'reactivated-node' AND i.remote_inbound_id = '42'`).Scan(&oldUserID, &inboundID); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO user_billing_records (id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, created_at) VALUES ('reactivated-order', ?, 'monthly', 100, 'CNY', ?, ?, ?, 'confirmed', 'manual', ?)`, oldUserID, nowText, now.Add(30*24*time.Hour).Format(time.RFC3339Nano), nowText, nowText); err != nil {
		t.Fatal(err)
	}

	for count := 1; count <= missingInboundArchiveAfter; count++ {
		observedAt := now.Add(time.Duration(count) * time.Minute)
		missing := map[string]any{
			"node_key": "reactivated-node", "sync_id": "reactivated-missing-" + strconv.Itoa(count), "observed_at": observedAt.Format(time.RFC3339Nano),
			"status": map[string]any{"xray_running": true}, "inbounds": []any{},
		}
		if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "reactivated-token", missing); result["code"] != successCode {
			t.Fatalf("missing sync %d response = %#v", count, result)
		}
	}

	var oldDeleted string
	if err := database.QueryRow(`SELECT COALESCE(deleted_at, '') FROM users WHERE id = ?`, oldUserID).Scan(&oldDeleted); err != nil {
		t.Fatal(err)
	}
	if oldDeleted == "" {
		t.Fatalf("old user was not retired after inbound archive")
	}

	reappearedAt := now.Add(4 * time.Minute)
	reappeared := lifecycleSyncPayload("reactivated-node", "reactivated-2", reappearedAt, []map[string]any{
		{"remote_id": "new-client", "email": "new@example.com", "enable": true, "expiry_time": reappearedAt.Add(30 * 24 * time.Hour).Unix()},
	})
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "reactivated-token", reappeared); result["code"] != successCode {
		t.Fatalf("reappeared sync response = %#v", result)
	}

	var newUserID, inboundUserID string
	if err := database.QueryRow(`SELECT id FROM users WHERE id <> ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1`, oldUserID).Scan(&newUserID); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT user_id FROM inbounds WHERE id = ?`, inboundID).Scan(&inboundUserID); err != nil {
		t.Fatal(err)
	}
	if newUserID == "" || newUserID == oldUserID || inboundUserID != newUserID {
		t.Fatalf("reactivated user old=%q new=%q inbound=%q", oldUserID, newUserID, inboundUserID)
	}

	var activeMapping, oldClients, newClients, replacementEvents, billingRecords, oldProfileEvents, newProfileEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_inbounds WHERE user_id = ? AND inbound_id = ? AND active_to IS NULL`, newUserID, inboundID).Scan(&activeMapping); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients WHERE inbound_id = ? AND remote_client_id = 'old-client'`, inboundID).Scan(&oldClients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients WHERE inbound_id = ? AND remote_client_id = 'new-client'`, inboundID).Scan(&newClients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'client_set_replacement_detected'`).Scan(&replacementEvents); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = ?`, oldUserID).Scan(&billingRecords); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ? AND event_status = 'resolved' AND requires_action = 0`, oldUserID).Scan(&oldProfileEvents); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ? AND event_status = 'open' AND requires_action = 1`, newUserID).Scan(&newProfileEvents); err != nil {
		t.Fatal(err)
	}
	if activeMapping != 1 || oldClients != 0 || newClients != 1 || replacementEvents != 0 || billingRecords != 1 || oldProfileEvents != 1 || newProfileEvents != 1 {
		t.Fatalf("reactivated lifecycle mapping=%d oldClients=%d newClients=%d replacementEvents=%d billingRecords=%d oldProfileEvents=%d newProfileEvents=%d", activeMapping, oldClients, newClients, replacementEvents, billingRecords, oldProfileEvents, newProfileEvents)
	}
}

func TestDirectInboundResetKeepsDisabledInboundState(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC().Truncate(time.Second)
	nowText := now.Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('disabled-node', 'disabled-node', 'Disabled Node', 'relay', 'online', ?, ?)`, nowText, nowText); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('disabled-credential', 'disabled-node', ?, ?, ?)`, hashToken("disabled-token"), nowText, nowText); err != nil {
		t.Fatal(err)
	}
	payload := lifecycleSyncPayload("disabled-node", "disabled-1", now, []map[string]any{{"remote_id": "disabled-client", "email": "disabled@example.com", "enable": true}})
	payload["inbounds"].([]any)[0].(map[string]any)["enable"] = false
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "disabled-token", payload); result["code"] != successCode {
		t.Fatalf("sync response = %#v", result)
	}
	var inboundID string
	if err := database.QueryRow(`SELECT id FROM inbounds WHERE node_id = 'disabled-node' AND remote_inbound_id = '42'`).Scan(&inboundID); err != nil {
		t.Fatal(err)
	}
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	reset := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/inbounds/"+inboundID+"/reset-user", token, nil)
	if reset["code"] != successCode {
		t.Fatalf("reset response = %#v", reset)
	}
	newUserID := reset["data"].(map[string]any)["newUserId"].(string)
	var status string
	if err := database.QueryRow(`SELECT status FROM users WHERE id = ?`, newUserID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "disabled" {
		t.Fatalf("replacement user status = %q, want disabled", status)
	}
}

func lifecycleSyncPayload(nodeKey, syncID string, observedAt time.Time, clients []map[string]any) map[string]any {
	return map[string]any{
		"node_key": nodeKey, "sync_id": syncID, "observed_at": observedAt.Format(time.RFC3339Nano),
		"status": map[string]any{"xray_running": true},
		"inbounds": []any{map[string]any{
			"remote_id": 42, "tag": "lifecycle-inbound", "remark": "Lifecycle Inbound", "protocol": "vless", "port": 443,
			"enable": true, "up": 40, "down": 60, "all_time": 100, "clients_complete": true, "clients": clients,
		}},
	}
}
