package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRenewalConfirmHonoursEditedAmountAndPaidDate 校验运营在确认续费时可以改金额与收款日期，
// 并且收款日期按北京时间归属财务月份（月末 23:59 不会被算到下个月）。
func TestRenewalConfirmHonoursEditedAmountAndPaidDate(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at)
VALUES ('edit-relay', 'edit-relay', '编辑金额线路机', 'relay', '2026-09-01T00:00:00Z', '2026-09-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO users (id, display_name, status, monthly_fee, billing_cycle, billing_amount, expiry_time, created_at, updated_at)
VALUES ('edit-user', '编辑金额用户', 'active', 100, 'monthly', 100, '2026-10-05T00:00:00Z', '2026-09-01T00:00:00Z', '2026-09-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, expiry_time, first_seen_at, last_seen_at)
VALUES ('edit-inbound', 'edit-relay', '1', 'edit-user', 'user', '2026-10-05T00:00:00Z', '2026-09-01T00:00:00Z', '2026-09-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES ('edit-mapping', 'edit-user', 'edit-inbound', 1, '2026-09-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	// 建议金额 100，运营实际收 88.5，并在 9 月最后一天深夜收款
	if _, err := database.Exec(`INSERT INTO user_renewal_candidates
(id, user_id, inbound_id, old_expiry_at, new_expiry_at, detected_at, suggested_cycle, suggested_amount, currency, status)
VALUES ('edit-candidate', 'edit-user', 'edit-inbound', '2026-09-30T00:00:00Z', '2026-10-31T00:00:00Z',
        '2026-09-30T00:00:00Z', 'monthly', 100, 'CNY', 'pending')`); err != nil {
		t.Fatal(err)
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	confirmed := doJSON(t, ts.Client(), http.MethodPost,
		ts.URL+"/api/users/edit-user/renewals/edit-candidate/confirm", token, map[string]any{
			"amount":       88.5,
			"billingCycle": "monthly",
			"paidAt":       "2026-09-30T23:59:59+08:00",
		})
	if confirmed["code"] != successCode {
		t.Fatalf("renewal confirm = %#v", confirmed)
	}

	var amount float64
	var paidAt string
	if err := database.QueryRow(`SELECT amount, paid_at FROM user_billing_records
WHERE user_id = 'edit-user' AND renewal_candidate_id = 'edit-candidate'`).Scan(&amount, &paidAt); err != nil {
		t.Fatalf("read billing record: %v", err)
	}
	if amount != 88.5 {
		t.Fatalf("stored amount = %v, want 88.5", amount)
	}
	// 北京时间 9-30 23:59:59 对应 UTC 9-30 15:59:59
	if paidAt != "2026-09-30T15:59:59Z" {
		t.Fatalf("stored paid_at = %q, want 2026-09-30T15:59:59Z", paidAt)
	}

	// 财务按北京时间归月：这笔收款必须落在 9 月，10 月不应包含
	september := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-09", token, nil)
	if september["code"] != successCode {
		t.Fatalf("september summary = %#v", september)
	}
	sepCash, _ := september["data"].(map[string]any)["cashIncome"].(float64)
	if sepCash < 88.5 {
		t.Fatalf("september cash income = %v, want >= 88.5", sepCash)
	}
	october := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/costs/summary?period=2026-10", token, nil)
	if october["code"] != successCode {
		t.Fatalf("october summary = %#v", october)
	}
	octCash, _ := october["data"].(map[string]any)["cashIncome"].(float64)
	if octCash != 0 {
		t.Fatalf("october cash income = %v, want 0 (month-end payment must stay in September)", octCash)
	}
}
