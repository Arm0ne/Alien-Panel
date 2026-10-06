package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestRenewalDetectionAfterNodeRecovery(t *testing.T) {
	oldExpiry := time.Date(2026, 10, 6, 6, 30, 0, 0, time.UTC)
	newExpiry := oldExpiry.AddDate(0, 1, 0)
	for _, tc := range []struct {
		name         string
		intermediate string
		unchanged    bool
		replacement  bool
		wantRenewal  int
	}{
		{name: "offline only", wantRenewal: 1},
		{name: "disabled inbound", intermediate: "inbound", wantRenewal: 1},
		{name: "disabled clients", intermediate: "clients", wantRenewal: 1},
		{name: "empty client snapshot", intermediate: "empty", wantRenewal: 1},
		{name: "extension while disabled", intermediate: "extended-disabled", wantRenewal: 1},
		{name: "resume without extension", intermediate: "clients", unchanged: true},
		{name: "unlimited service clears old expiry", intermediate: "unlimited"},
		{name: "different customer after disabled snapshot", intermediate: "clients", replacement: true},
		{name: "different customer after empty snapshot", intermediate: "empty", replacement: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server, database := testServer(t)
			ts := httptest.NewServer(server.Handler())
			defer ts.Close()
			initialAt := oldExpiry.Add(-time.Hour)
			initialText := initialAt.Format(time.RFC3339Nano)
			if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('recovery-node', 'recovery-node', 'Recovery node', 'relay', ?, ?)`, initialText, initialText); err != nil {
				t.Fatal(err)
			}
			if _, err := database.Exec(`INSERT INTO node_credentials (id, node_id, token_hash, last_rotated_at, created_at) VALUES ('recovery-credential', 'recovery-node', ?, ?, ?)`, hashToken("recovery-token"), initialText, initialText); err != nil {
				t.Fatal(err)
			}
			syncNumber := 0
			syncInbound := func(at time.Time, enabled bool, clients []map[string]any) {
				t.Helper()
				syncNumber++
				payload := lifecycleSyncPayload("recovery-node", "recovery-"+strconv.Itoa(syncNumber), at, clients)
				payload["inbounds"].([]any)[0].(map[string]any)["enable"] = enabled
				response := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", "recovery-token", payload)
				if response["code"] != successCode {
					t.Fatalf("sync %d: %#v", syncNumber, response)
				}
			}
			clients := func(id string, enabled bool, expiry int64) []map[string]any {
				return []map[string]any{{"remote_id": id, "enable": enabled, "expiry_time": expiry}}
			}
			syncInbound(initialAt, true, clients("same-client", true, oldExpiry.Unix()))
			var userID string
			if err := database.QueryRow(`SELECT user_id FROM inbounds WHERE node_id = 'recovery-node'`).Scan(&userID); err != nil {
				t.Fatal(err)
			}
			recoveredAt := oldExpiry.Add(time.Hour)
			server.refreshOperationalStatuses(recoveredAt)
			var nodeStatus, status, expiry string
			if err := database.QueryRow(`SELECT health_status FROM nodes WHERE id = 'recovery-node'`).Scan(&nodeStatus); err != nil {
				t.Fatal(err)
			}
			if nodeStatus != "offline" {
				t.Fatalf("node status = %q, want offline", nodeStatus)
			}
			if tc.intermediate != "" {
				middleClients := clients("same-client", tc.intermediate != "clients", oldExpiry.Unix())
				middleEnabled := tc.intermediate != "inbound" && tc.intermediate != "extended-disabled"
				if tc.intermediate == "empty" {
					middleClients = nil
				} else if tc.intermediate == "unlimited" {
					middleClients = clients("same-client", true, 0)
				} else if tc.intermediate == "extended-disabled" {
					middleClients = clients("same-client", true, newExpiry.Unix())
				}
				// Repeated disabled snapshots must not consume the renewal baseline.
				syncInbound(recoveredAt, middleEnabled, middleClients)
				syncInbound(recoveredAt.Add(time.Minute), middleEnabled, middleClients)
				if err := database.QueryRow(`SELECT status, COALESCE(expiry_time, '') FROM users WHERE id = ?`, userID).Scan(&status, &expiry); err != nil {
					t.Fatal(err)
				}
				if tc.intermediate == "unlimited" {
					if status != "active" || expiry != "" {
						t.Fatalf("unlimited state = %q %q", status, expiry)
					}
				} else if status != "disabled" || expiry != oldExpiry.Format(time.RFC3339Nano) {
					t.Fatalf("disabled state = %q %q, want disabled with previous expiry", status, expiry)
				}
			}
			finalExpiry := newExpiry
			if tc.unchanged {
				finalExpiry = oldExpiry
			}
			clientID := "same-client"
			if tc.replacement {
				clientID = "new-customer-client"
			}
			syncInbound(recoveredAt.Add(2*time.Minute), true, clients(clientID, true, finalExpiry.Unix()))
			syncInbound(recoveredAt.Add(3*time.Minute), true, clients(clientID, true, finalExpiry.Unix()))
			var candidateCount, eventCount, billingCount int
			if err := database.QueryRow(`SELECT COUNT(*) FROM user_renewal_candidates WHERE user_id = ? AND status = 'pending'`, userID).Scan(&candidateCount); err != nil {
				t.Fatal(err)
			}
			if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'renewal_candidate_detected' AND requires_action = 1 AND event_status = 'open'`).Scan(&eventCount); err != nil {
				t.Fatal(err)
			}
			if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = ?`, userID).Scan(&billingCount); err != nil {
				t.Fatal(err)
			}
			if candidateCount != tc.wantRenewal || eventCount != tc.wantRenewal || billingCount != 0 {
				t.Fatalf("candidates=%d events=%d payments=%d, want %d candidates/events and no automatic payment", candidateCount, eventCount, billingCount, tc.wantRenewal)
			}
			if err := database.QueryRow(`SELECT expiry_time FROM users WHERE id = ?`, userID).Scan(&expiry); err != nil {
				t.Fatal(err)
			}
			if expiry != finalExpiry.Format(time.RFC3339Nano) {
				t.Fatalf("updated expiry = %q", expiry)
			}
			if tc.wantRenewal == 1 {
				var from, to string
				if err := database.QueryRow(`SELECT old_expiry_at, new_expiry_at FROM user_renewal_candidates WHERE user_id = ?`, userID).Scan(&from, &to); err != nil {
					t.Fatal(err)
				}
				if from != oldExpiry.Format(time.RFC3339Nano) || to != newExpiry.Format(time.RFC3339Nano) {
					t.Fatalf("renewal interval = %s -> %s", from, to)
				}
			}
		})
	}
}

func TestManualRenewalPreservesExactServiceBoundary(t *testing.T) {
	server, database := testServer(t)
	createdAt := "2026-09-06T06:30:00Z"
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('boundary-user', 'Boundary user', 'active', ?, ?)`, createdAt, createdAt); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	url := ts.URL + "/api/users/boundary-user/billing-records"
	order := map[string]any{
		"billingCycle": "monthly", "amount": 20, "orderType": "initial",
		"serviceFrom": createdAt, "serviceTo": "2026-10-06T06:30:00Z",
	}
	if result := doJSON(t, ts.Client(), http.MethodPost, url, token, order); result["code"] != successCode {
		t.Fatalf("initial order = %#v", result)
	}
	order["orderType"] = "renewal"
	order["serviceFrom"] = "2026-10-06T00:00:00Z"
	order["serviceTo"] = "2026-11-06T06:30:00.000Z"
	if status, _ := doJSONWithStatus(t, ts.Client(), http.MethodPost, url, token, order); status != http.StatusConflict {
		t.Fatalf("truncated boundary status = %d, want conflict", status)
	}
	// JavaScript serializes UTC with milliseconds; the instant must still
	// coincide with the old order's exact end rather than overlap it.
	order["serviceFrom"] = "2026-10-06T06:30:00.000Z"
	if result := doJSON(t, ts.Client(), http.MethodPost, url, token, order); result["code"] != successCode {
		t.Fatalf("adjacent renewal = %#v", result)
	}
	if status, _ := doJSONWithStatus(t, ts.Client(), http.MethodPost, url, token, order); status != http.StatusConflict {
		t.Fatalf("duplicate renewal status = %d, want conflict", status)
	}
	var count int
	var amount float64
	if err := database.QueryRow(`SELECT COUNT(*), SUM(amount) FROM user_billing_records WHERE user_id = 'boundary-user'`).Scan(&count, &amount); err != nil {
		t.Fatal(err)
	}
	if count != 2 || amount != 40 {
		t.Fatalf("orders=%d amount=%v, want 2 orders and 40 CNY", count, amount)
	}
}
