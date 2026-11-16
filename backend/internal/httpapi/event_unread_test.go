package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUnreadSummaryMatchesUnreadFilter 保证铃铛的未读数与事件中心「未读」标签
// 使用同一口径：通知类未读（不含需要运营动作的事件，也不含已处理事件）。
func TestUnreadSummaryMatchesUnreadFilter(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	seed := []struct {
		id             string
		requiresAction int
		actionType     string
		status         string
		acknowledged   int
	}{
		{"info-unread-1", 0, "", "open", 0},
		{"info-unread-2", 0, "", "open", 0},
		{"info-read", 0, "", "open", 1},
		{"action-unread", 1, "inspect_node", "open", 0},
		{"action-typed", 0, "confirm_renewal", "open", 0},
		{"info-resolved", 0, "", "resolved", 0},
	}
	for _, item := range seed {
		if _, err := database.Exec(`INSERT INTO node_events
(id, event_type, severity, message, created_at, acknowledged, requires_action, event_status, action_type)
VALUES (?, 'sync_failed', 'warning', '未读口径测试', '2026-11-01T00:00:00Z', ?, ?, ?, NULLIF(?, ''))`,
			item.id, item.acknowledged, item.requiresAction, item.status, item.actionType); err != nil {
			t.Fatalf("seed %s: %v", item.id, err)
		}
	}

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	summary := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events/summary", token, nil)
	if summary["code"] != successCode {
		t.Fatalf("summary = %#v", summary)
	}
	unreadCount, _ := summary["data"].(map[string]any)["unreadCount"].(float64)
	if unreadCount != 2 {
		t.Fatalf("summary unreadCount = %v, want 2", unreadCount)
	}

	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events?status=unread&page=1&page_size=50", token, nil)
	if list["code"] != successCode {
		t.Fatalf("unread list = %#v", list)
	}
	total, _ := list["data"].(map[string]any)["total"].(float64)
	if total != unreadCount {
		t.Fatalf("unread list total = %v, summary unreadCount = %v; 两者必须一致", total, unreadCount)
	}

	// 待处理标签仍应包含需要动作的事件
	pending := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events?status=pending&page=1&page_size=50", token, nil)
	pendingTotal, _ := pending["data"].(map[string]any)["total"].(float64)
	if pendingTotal != 2 {
		t.Fatalf("pending total = %v, want 2", pendingTotal)
	}

	// 全部标为已读只影响通知类未读，不清除待处理事件的动作标记
	if marked := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/events/read-all", token, nil); marked["code"] != successCode {
		t.Fatalf("mark all read = %#v", marked)
	}
	after := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events/summary", token, nil)
	afterUnread, _ := after["data"].(map[string]any)["unreadCount"].(float64)
	if afterUnread != 0 {
		t.Fatalf("unreadCount after mark-all = %v, want 0", afterUnread)
	}
	afterPending := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/events?status=pending&page=1&page_size=50", token, nil)
	afterPendingTotal, _ := afterPending["data"].(map[string]any)["total"].(float64)
	if afterPendingTotal != 2 {
		t.Fatalf("pending total after mark-all = %v, want 2 (待处理不应被清除)", afterPendingTotal)
	}
}
