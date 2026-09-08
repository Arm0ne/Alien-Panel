package httpapi

import (
	"database/sql"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func insertOrphanTestNode(t *testing.T, database *sql.DB, id, nodeType string, now string) {
	t.Helper()
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`, id, id, id+" name", nodeType, now, now); err != nil {
		t.Fatalf("insert node %s: %v", id, err)
	}
}

func insertOrphanTestUser(t *testing.T, database *sql.DB, userID, nodeID, inboundID, now string, primary bool) {
	t.Helper()
	if _, err := database.Exec(`INSERT INTO users
(id, display_name, status, monthly_fee, billing_cycle, billing_amount, currency, created_at, updated_at)
VALUES (?, ?, 'active', 100, 'monthly', 100, 'CNY', ?, ?)`, userID, userID+" name", now, now); err != nil {
		t.Fatalf("insert user %s: %v", userID, err)
	}
	insertOrphanTestInbound(t, database, userID, nodeID, inboundID, now, primary)
}

func insertOrphanTestInbound(t *testing.T, database *sql.DB, userID, nodeID, inboundID, now string, primary bool) {
	t.Helper()
	if _, err := database.Exec(`INSERT INTO inbounds
(id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at, last_seen_at)
VALUES (?, ?, ?, ?, 'user', ?, ?, ?)`, inboundID, nodeID, inboundID, userID, inboundID+" tag", now, now); err != nil {
		t.Fatalf("insert inbound %s: %v", inboundID, err)
	}
	primaryValue := 0
	if primary {
		primaryValue = 1
	}
	if _, err := database.Exec(`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES (?, ?, ?, ?, ?)`, "mapping-"+inboundID, userID, inboundID, primaryValue, now); err != nil {
		t.Fatalf("insert user inbound mapping %s: %v", inboundID, err)
	}
}

func orphanTestToken(t *testing.T, ts *httptest.Server) string {
	t.Helper()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	return login["data"].(map[string]any)["token"].(string)
}

func TestNodeDeletionAutoDeletesOrphanedUserAndPreservesFinancialHistory(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339Nano)
	period := now.Format("2006-01")
	serviceFrom, serviceTo, err := monthBounds(period)
	if err != nil {
		t.Fatal(err)
	}
	insertOrphanTestNode(t, database, "orphan-relay", "relay", nowText)
	insertOrphanTestUser(t, database, "orphan-user", "orphan-relay", "orphan-inbound", nowText, true)
	if _, err := database.Exec(`INSERT INTO user_billing_records
(id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, created_at)
VALUES ('orphan-billing', 'orphan-user', 'monthly', 300, 'CNY', ?, ?, ?, 'confirmed', 'manual', ?)`, serviceFrom, serviceTo, nowText, nowText); err != nil {
		t.Fatalf("insert billing record: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO user_renewal_candidates
(id, user_id, inbound_id, old_expiry_at, new_expiry_at, detected_at, suggested_cycle, suggested_amount, currency, status)
VALUES ('orphan-candidate', 'orphan-user', 'orphan-inbound', ?, ?, ?, 'monthly', 100, 'CNY', 'pending')`, serviceFrom, serviceTo, nowText); err != nil {
		t.Fatalf("insert renewal candidate: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO node_events
(id, node_id, event_type, severity, message, created_at, event_category, title, requires_action, event_status, resource_type, resource_id, source)
VALUES ('orphan-event', 'orphan-relay', 'renewal_candidate_detected', 'warning', 'pending renewal', ?, 'business', 'pending renewal', 1, 'open', 'renewal', 'orphan-candidate', 'agent')`, nowText); err != nil {
		t.Fatalf("insert renewal event: %v", err)
	}

	token := orphanTestToken(t, ts)
	status, deleted := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/nodes/orphan-relay", token, nil)
	if status != http.StatusOK || deleted["code"] != successCode {
		t.Fatalf("delete node status=%d response=%#v", status, deleted)
	}

	var deletedAt sql.NullString
	var userStatus string
	if err := database.QueryRow(`SELECT deleted_at, status FROM users WHERE id = 'orphan-user'`).Scan(&deletedAt, &userStatus); err != nil {
		t.Fatalf("read deleted user: %v", err)
	}
	if !deletedAt.Valid || deletedAt.String == "" || userStatus != "disabled" {
		t.Fatalf("user lifecycle state deleted_at=%#v status=%q", deletedAt, userStatus)
	}
	var billingCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = 'orphan-user' AND status = 'confirmed'`).Scan(&billingCount); err != nil {
		t.Fatalf("count billing records: %v", err)
	}
	if billingCount != 1 {
		t.Fatalf("confirmed billing records = %d, want 1", billingCount)
	}
	var candidateStatus string
	var candidateProcessedAt, candidateNotes sql.NullString
	if err := database.QueryRow(`SELECT status, processed_at, notes FROM user_renewal_candidates WHERE id = 'orphan-candidate'`).Scan(&candidateStatus, &candidateProcessedAt, &candidateNotes); err != nil {
		t.Fatalf("read renewal candidate: %v", err)
	}
	if candidateStatus != "rejected" || !candidateProcessedAt.Valid || !candidateNotes.Valid {
		t.Fatalf("candidate after user deletion status=%q processed_at=%#v notes=%#v", candidateStatus, candidateProcessedAt, candidateNotes)
	}
	var eventStatus string
	var requiresAction int
	var eventNodeID sql.NullString
	if err := database.QueryRow(`SELECT event_status, requires_action, node_id FROM node_events WHERE id = 'orphan-event'`).Scan(&eventStatus, &requiresAction, &eventNodeID); err != nil {
		t.Fatalf("read renewal event: %v", err)
	}
	if eventStatus != "resolved" || requiresAction != 0 || eventNodeID.Valid {
		t.Fatalf("renewal event after node deletion status=%q requires_action=%d node_id=%#v", eventStatus, requiresAction, eventNodeID)
	}
	var auditCount int
	if err := database.QueryRow(`SELECT COUNT(*) FROM audit_logs
WHERE action = 'user.auto_delete' AND resource_id = 'orphan-user'
  AND after_json LIKE '%no_current_node_association%'`).Scan(&auditCount); err != nil {
		t.Fatalf("read auto-delete audit: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("auto-delete audit count = %d, want 1", auditCount)
	}

	users := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/users?page_size=20", token, nil)
	userData := users["data"].(map[string]any)
	if users["code"] != successCode || userData["total"] != float64(0) {
		t.Fatalf("deleted user remained in list: %#v", users)
	}
	detailStatus, detail := doJSONWithStatus(t, ts.Client(), http.MethodGet, ts.URL+"/api/users/orphan-user", token, nil)
	if detailStatus != http.StatusNotFound || detail["code"] != notFoundCode {
		t.Fatalf("deleted user detail status=%d response=%#v", detailStatus, detail)
	}
	dashboard := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/dashboard", token, nil)
	if dashboard["data"].(map[string]any)["users"].(map[string]any)["active"] != float64(0) {
		t.Fatalf("deleted user remained in dashboard: %#v", dashboard)
	}
	finance, err := server.financeSummary(period)
	if err != nil {
		t.Fatalf("finance after user deletion: %v", err)
	}
	if finance.EffectiveUserCount != 0 || math.Abs(finance.CashIncome-300) > 0.00001 || math.Abs(finance.MonthIncome-300) > 0.00001 {
		t.Fatalf("financial history after user deletion = %#v", finance)
	}
}

func TestNodeDeletionRetainsUserUntilLastRelayAssociationIsRemoved(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	insertOrphanTestNode(t, database, "first-relay", "relay", nowText)
	insertOrphanTestNode(t, database, "second-relay", "relay", nowText)
	insertOrphanTestUser(t, database, "multi-relay-user", "first-relay", "first-inbound", nowText, true)
	insertOrphanTestInbound(t, database, "multi-relay-user", "second-relay", "second-inbound", nowText, false)

	token := orphanTestToken(t, ts)
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/nodes/first-relay", token, nil); status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("delete first relay status=%d response=%#v", status, response)
	}
	var afterFirst sql.NullString
	if err := database.QueryRow(`SELECT deleted_at FROM users WHERE id = 'multi-relay-user'`).Scan(&afterFirst); err != nil {
		t.Fatalf("read user after first deletion: %v", err)
	}
	if afterFirst.Valid {
		t.Fatalf("user deleted despite a second relay inbound: %#v", afterFirst)
	}
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/nodes/second-relay", token, nil); status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("delete second relay status=%d response=%#v", status, response)
	}
	var afterSecond sql.NullString
	if err := database.QueryRow(`SELECT deleted_at FROM users WHERE id = 'multi-relay-user'`).Scan(&afterSecond); err != nil {
		t.Fatalf("read user after final deletion: %v", err)
	}
	if !afterSecond.Valid {
		t.Fatal("user was not deleted after its last relay association was removed")
	}
}

func TestLandingDeletionDoesNotDeleteUserWithRelayInbound(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	nowText := time.Now().UTC().Format(time.RFC3339Nano)
	insertOrphanTestNode(t, database, "path-relay", "relay", nowText)
	insertOrphanTestNode(t, database, "path-landing", "landing", nowText)
	insertOrphanTestUser(t, database, "landing-path-user", "path-relay", "path-inbound", nowText, true)
	if _, err := database.Exec(`INSERT INTO exit_ips
(id, landing_node_id, owner_node_id, ip, source_type, created_at, updated_at)
VALUES ('path-exit', 'path-landing', 'path-landing', '198.51.100.80', 'node', ?, ?)`, nowText, nowText); err != nil {
		t.Fatalf("insert path exit IP: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO user_paths
(id, user_id, relay_node_id, landing_node_id, exit_ip_id, mode, active_from, created_at, updated_at)
VALUES ('landing-path', 'landing-path-user', 'path-relay', 'path-landing', 'path-exit', 'landing', ?, ?, ?)`, nowText, nowText, nowText); err != nil {
		t.Fatalf("insert user path: %v", err)
	}

	token := orphanTestToken(t, ts)
	if status, response := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/nodes/path-landing", token, nil); status != http.StatusOK || response["code"] != successCode {
		t.Fatalf("delete landing status=%d response=%#v", status, response)
	}
	var deletedAt sql.NullString
	if err := database.QueryRow(`SELECT deleted_at FROM users WHERE id = 'landing-path-user'`).Scan(&deletedAt); err != nil {
		t.Fatalf("read user after landing deletion: %v", err)
	}
	if deletedAt.Valid {
		t.Fatalf("user deleted despite an active relay inbound: %#v", deletedAt)
	}
}
