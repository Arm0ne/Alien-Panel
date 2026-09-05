package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAgentExpiryExtensionCreatesConfirmableRenewal(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	oldExpiry := now.Add(30 * 24 * time.Hour).Format(time.RFC3339Nano)
	newExpiry := now.Add(60 * 24 * time.Hour).Format(time.RFC3339Nano)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('billing-relay', 'billing-relay', '计费线路机', 'relay', ?, ?)`, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, monthly_fee, billing_cycle, billing_amount, expiry_time, created_at, updated_at) VALUES ('billing-user', '计费用户', 'active', 100, 'annual', 1080, ?, ?, ?)`, oldExpiry, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, expiry_time, first_seen_at, last_seen_at) VALUES ('billing-inbound', 'billing-relay', '1', 'billing-user', 'user', ?, ?, ?)`, oldExpiry, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)); err != nil {
		t.Fatal(err)
	}

	tx, err := database.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err := server.ensureRelayInboundUser(tx, "relay", "billing-inbound", "1", agentInboundPayload{Enable: true, ExpiryTime: newExpiryEpoch(newExpiry)}, newExpiry, now.Add(time.Minute)); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_renewal_candidates WHERE user_id = 'billing-user' AND status = 'pending'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("pending renewal count = %d, want 1", count)
	}
	var renewalEvents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM node_events WHERE event_type = 'renewal_candidate_detected' AND requires_action = 1`).Scan(&renewalEvents); err != nil {
		t.Fatal(err)
	}
	if renewalEvents != 1 {
		t.Fatalf("renewal event count = %d, want 1", renewalEvents)
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/billing-user/renewals", token, nil)
	items := list["data"].([]any)
	if list["code"] != successCode || len(items) != 1 {
		t.Fatalf("renewal list = %#v", list)
	}
	candidateID := items[0].(map[string]any)["id"].(string)
	events := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events?status=pending", token, nil)
	eventData := events["data"].(map[string]any)
	if events["code"] != successCode || eventData["total"] != float64(1) {
		t.Fatalf("pending events = %#v", events)
	}
	confirmed := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/users/billing-user/renewals/"+candidateID+"/confirm", token, map[string]any{"amount": 1080, "billingCycle": "annual"})
	if confirmed["code"] != successCode {
		t.Fatalf("renewal confirmation = %#v", confirmed)
	}
	var records, pending int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = 'billing-user' AND status = 'confirmed'`).Scan(&records); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_renewal_candidates WHERE user_id = 'billing-user' AND status = 'pending'`).Scan(&pending); err != nil {
		t.Fatal(err)
	}
	if records != 1 || pending != 0 {
		t.Fatalf("billing records=%d pending=%d", records, pending)
	}
	counts := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events/summary", token, nil)
	if counts["code"] != successCode || counts["data"].(map[string]any)["pendingCount"] != float64(0) {
		t.Fatalf("event summary after confirmation = %#v", counts)
	}
}

func TestAnnualBillingIsAllocatedAcrossServiceMonths(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	created := "2026-08-01T00:00:00Z"
	expiry := "2026-10-15T00:00:00Z"
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, monthly_fee, billing_cycle, billing_amount, expiry_time, created_at, updated_at) VALUES ('annual-finance-user', '年费用户', 'active', 100, 'annual', 1200, ?, ?, ?)`, expiry, created, created); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO user_billing_records (id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, created_at) VALUES ('annual-finance-record', 'annual-finance-user', 'annual', 1200, 'CNY', '2026-09-15T00:00:00Z', '2026-10-15T00:00:00Z', '2026-09-15T00:00:00Z', 'confirmed', 'manual', ?)`, created); err != nil {
		t.Fatal(err)
	}
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	finance := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-09", token, nil)
	if finance["code"] != successCode {
		t.Fatalf("finance response = %#v", finance)
	}
	data := finance["data"].(map[string]any)
	// September contributes 16/30 of the 1,200 CNY service interval (the
	// half-open month includes Sep 15 through Sep 30).
	if data["monthIncome"] != float64(640) || data["cashIncome"] != float64(1200) {
		t.Fatalf("annual finance = %#v, want accrued=640 cash=1200", data)
	}
}

func newExpiryEpoch(value string) int64 {
	parsed, _ := time.Parse(time.RFC3339Nano, value)
	return parsed.Unix()
}
