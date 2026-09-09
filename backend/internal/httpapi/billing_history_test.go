package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHistoricalBillingImportPreviewAndVerification(t *testing.T) {
	server, database := testServer(t)
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('history-node', 'history-node', '历史线路机', 'relay', '2026-03-01T00:00:00Z', '2026-03-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, expiry_time, created_at, updated_at) VALUES ('history-user', '历史用户', 'active', '2026-09-01T00:00:00Z', '2026-03-01T00:00:00Z', '2026-03-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at, last_seen_at) VALUES ('history-inbound', 'history-node', '17', 'history-user', 'user', 'history-user', '2026-03-01T00:00:00Z', '2026-03-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	rows := []any{
		map[string]any{"userId": "history-user", "billingCycle": "monthly", "amount": 100, "serviceFrom": "2026-03-01", "serviceTo": "2026-04-01", "paidAt": "2026-03-01", "orderType": "initial", "verified": true},
		map[string]any{"userId": "history-user", "billingCycle": "monthly", "amount": 100, "serviceFrom": "2026-04-01", "serviceTo": "2026-05-01", "orderType": "renewal", "verified": false},
	}
	preview := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/billing/import", token, map[string]any{"records": rows, "dryRun": true})
	if preview["code"] != successCode || preview["data"].(map[string]any)["imported"] != float64(0) || len(preview["data"].(map[string]any)["records"].([]any)) != 2 {
		t.Fatalf("preview = %#v", preview)
	}
	imported := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/billing/import", token, map[string]any{"records": rows})
	if imported["code"] != successCode || imported["data"].(map[string]any)["imported"] != float64(2) || imported["data"].(map[string]any)["unverified"] != float64(1) {
		t.Fatalf("imported = %#v", imported)
	}
	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/history-user", token, nil)
	if detail["code"] != successCode || len(detail["data"].(map[string]any)["billingRecords"].([]any)) != 2 {
		t.Fatalf("detail = %#v", detail)
	}
	var pendingID string
	for _, item := range detail["data"].(map[string]any)["billingRecords"].([]any) {
		record := item.(map[string]any)
		if record["verificationStatus"] == "unverified" {
			pendingID = record["id"].(string)
		}
	}
	if pendingID == "" {
		t.Fatalf("unverified record missing: %#v", detail)
	}
	finance := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-04", token, nil)
	if finance["code"] != successCode || finance["data"].(map[string]any)["cashIncome"] != float64(0) || finance["data"].(map[string]any)["monthIncome"] != float64(0) {
		t.Fatalf("unverified income included: %#v", finance)
	}
	verified := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/users/history-user/billing-records/"+pendingID+"/verify", token, map[string]any{"paidAt": "2026-04-01"})
	if verified["code"] != successCode {
		t.Fatalf("verify = %#v", verified)
	}
	finance = doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-04", token, nil)
	if finance["code"] != successCode || finance["data"].(map[string]any)["cashIncome"] != float64(100) || finance["data"].(map[string]any)["monthIncome"] != float64(100) {
		t.Fatalf("verified income = %#v", finance)
	}
}
