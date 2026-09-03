package httpapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestOfflineRefreshPreservesLastSnapshotAndHistory(t *testing.T) {
	server, database := testServer(t)
	now := time.Now().UTC()
	oldSeen := now.Add(-10 * time.Minute).Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, last_seen_at, created_at, updated_at) VALUES ('offline-node', 'offline-node', 'Offline node', 'relay', 'online', ?, ?, ?)`, oldSeen, oldSeen, oldSeen); err != nil {
		t.Fatalf("insert node: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, tag, enable, all_time, last_seen_at, first_seen_at) VALUES ('offline-inbound', 'offline-node', '1', 'offline-user', 1, 1234, ?, ?)`, oldSeen, oldSeen); err != nil {
		t.Fatalf("insert inbound: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, all_time, source) VALUES ('offline-snapshot', 'offline-node', 'offline-inbound', ?, 1234, 'xpanel')`, oldSeen); err != nil {
		t.Fatalf("insert snapshot: %v", err)
	}

	server.refreshOperationalStatuses(now)
	var health string
	var inboundCount, snapshotCount int
	if err := database.QueryRow(`SELECT health_status FROM nodes WHERE id = 'offline-node'`).Scan(&health); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM inbounds WHERE id = 'offline-inbound'`).Scan(&inboundCount); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM traffic_snapshots WHERE id = 'offline-snapshot'`).Scan(&snapshotCount); err != nil {
		t.Fatal(err)
	}
	if health != "offline" || inboundCount != 1 || snapshotCount != 1 {
		t.Fatalf("offline refresh health=%q inboundCount=%d snapshotCount=%d; history must remain", health, inboundCount, snapshotCount)
	}
}

func TestConcurrentDuplicateSyncIsIdempotentAndNonNegative(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	const nodeKey = "stress-node"
	const nodeToken = "stress-node-token"
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES ('stress-node', ?, 'Stress node', 'relay', 'unknown', ?, ?)`, nodeKey, now, now); err != nil {
		t.Fatalf("insert node: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('stress-credential', 'stress-node', ?, ?, ?)`, hashToken(nodeToken), now, now); err != nil {
		t.Fatalf("insert credential: %v", err)
	}
	payload := stressSyncPayload(nodeKey, "stress-sync-1", now, 100)
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, payload); result["code"] != successCode {
		t.Fatalf("initial sync = %#v", result)
	}

	const duplicateRequests = 24
	var wait sync.WaitGroup
	errors := make(chan string, duplicateRequests)
	for index := 0; index < duplicateRequests; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, payload)
			if result["code"] != successCode {
				errors <- "duplicate response: " + toJSONText(result)
			}
		}()
	}
	wait.Wait()
	close(errors)
	for message := range errors {
		t.Error(message)
	}

	var syncRuns, snapshots, users, clients, archived, missingCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sync_runs WHERE sync_id = 'stress-sync-1'`).Scan(&syncRuns); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM traffic_snapshots WHERE sync_run_id IN (SELECT id FROM sync_runs WHERE sync_id = 'stress-sync-1')`).Scan(&snapshots); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&users); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM clients`).Scan(&clients); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM inbounds WHERE remote_inbound_id = '1' AND deleted_at IS NOT NULL`).Scan(&archived); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COALESCE(missing_sync_count, 0) FROM inbounds WHERE remote_inbound_id = '1'`).Scan(&missingCount); err != nil {
		t.Fatal(err)
	}
	if syncRuns != 1 || snapshots != 1 || users != 1 || clients != 1 || archived != 0 || missingCount != 0 {
		t.Fatalf("duplicate sync changed state: syncRuns=%d snapshots=%d users=%d clients=%d archived=%d missingCount=%d", syncRuns, snapshots, users, clients, archived, missingCount)
	}

	resetPayload := stressSyncPayload(nodeKey, "stress-sync-2", time.Now().UTC().Add(time.Second).Format(time.RFC3339Nano), 50)
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeToken, resetPayload); result["code"] != successCode {
		t.Fatalf("reset sync = %#v", result)
	}
	var resetEvents, negativeSnapshots int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE node_id = 'stress-node' AND event_type = 'traffic_reset'`).Scan(&resetEvents); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM traffic_snapshots WHERE all_time < 0 OR up < 0 OR down < 0`).Scan(&negativeSnapshots); err != nil {
		t.Fatal(err)
	}
	if resetEvents != 1 || negativeSnapshots != 0 {
		t.Fatalf("reset handling resetEvents=%d negativeSnapshots=%d", resetEvents, negativeSnapshots)
	}
}

func TestConcurrentDifferentNodeSyncsCompleteIndependently(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	const nodeCount = 4
	for index := 0; index < nodeCount; index++ {
		nodeID := "parallel-node-" + strconv.Itoa(index)
		now := time.Now().UTC().Format(time.RFC3339Nano)
		if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, health_status, created_at, updated_at) VALUES (?, ?, ?, 'relay', 'unknown', ?, ?)`, nodeID, nodeID, nodeID, now, now); err != nil {
			t.Fatalf("insert node %d: %v", index, err)
		}
		if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES (?, ?, ?, ?, ?)`, nodeID+"-credential", nodeID, hashToken(nodeID+"-token"), now, now); err != nil {
			t.Fatalf("insert credential %d: %v", index, err)
		}
	}

	var wait sync.WaitGroup
	errors := make(chan string, nodeCount)
	for index := 0; index < nodeCount; index++ {
		index := index
		wait.Add(1)
		go func() {
			defer wait.Done()
			nodeID := "parallel-node-" + strconv.Itoa(index)
			result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", nodeID+"-token", stressSyncPayload(nodeID, nodeID+"-sync", time.Now().UTC().Format(time.RFC3339Nano), int64(index+1)))
			if result["code"] != successCode {
				errors <- nodeID + ": " + toJSONText(result)
			}
		}()
	}
	wait.Wait()
	close(errors)
	for message := range errors {
		t.Error(message)
	}
	var successful int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sync_runs WHERE status = 'success' AND node_id LIKE 'parallel-node-%'`).Scan(&successful); err != nil {
		t.Fatal(err)
	}
	if successful != nodeCount {
		t.Fatalf("successful parallel syncs = %d, want %d", successful, nodeCount)
	}
}

func stressSyncPayload(nodeKey, syncID, observedAt string, allTime int64) map[string]any {
	return map[string]any{
		"node_key": nodeKey, "sync_id": syncID, "observed_at": observedAt,
		"status": map[string]any{"xray_running": true},
		"inbounds": []any{map[string]any{
			"remote_id": 1, "tag": "stress-inbound", "remark": "Stress inbound", "protocol": "vless", "enable": true,
			"up": allTime / 2, "down": allTime - allTime/2, "all_time": allTime, "clients": []any{
				map[string]any{"remote_id": "stress-client", "email": "stress@example.com", "enable": true, "all_time": allTime},
			},
		}},
	}
}

func toJSONText(value map[string]any) string {
	return fmt.Sprintf("%v", value["msg"])
}
