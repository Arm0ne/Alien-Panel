package httpapi

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// profileEventCounts 返回指定用户的待办/已关闭事件数量。
func profileEventCounts(t *testing.T, db *sql.DB, userID string) (open, dismissed int) {
	t.Helper()
	if err := db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user' AND resource_id = ?
  AND COALESCE(event_status, 'open') NOT IN ('resolved', 'dismissed')`, userID).Scan(&open); err != nil {
		t.Fatalf("count open profile events: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_type = 'user' AND resource_id = ?
  AND COALESCE(event_status, 'open') = 'dismissed'`, userID).Scan(&dismissed); err != nil {
		t.Fatalf("count dismissed profile events: %v", err)
	}
	return open, dismissed
}

// TestNewUserProfileEventFollowsInboundAvailability 覆盖「测试入站停用后事件卡住」的问题：
//  1. 新建入站 → 生成待完善资料事件
//  2. 入站停用 → 事件自动关闭，不再卡在待处理
//  3. 入站恢复启用 → 用户仍未定价，事件重新提醒
//  4. 用户完成定价后再停用 → 不再产生新事件
func TestNewUserProfileEventFollowsInboundAvailability(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"name": "事件生命周期线路机", "type": "relay",
	})
	nodeData := created["data"].(map[string]any)
	nodeID := nodeData["nodeId"].(string)
	nodeKey := nodeData["nodeKey"].(string)
	agentToken := nodeData["token"].(string)

	syncInbound := func(syncID, observedAt string, enable bool) {
		payload := map[string]any{
			"node_key":    nodeKey,
			"sync_id":     syncID,
			"observed_at": observedAt,
			"status": map[string]any{
				"agent_version": "v1.0.4",
				"xray_running":  true,
			},
			"inbounds": []any{
				map[string]any{
					"remote_id": 8801, "tag": "event-test", "remark": "事件测试入站",
					"protocol": "vless", "port": 443, "listen": "0.0.0.0", "enable": enable,
					"expiry_time": 1900000000, "up": 10, "down": 20, "all_time": 30,
					"clients": []any{
						map[string]any{
							"remote_id": "event-client", "email": "event@example.com", "enable": enable,
							"expiry_time": 1900000000, "up": 10, "down": 20, "all_time": 30,
						},
					},
				},
			},
		}
		result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, payload)
		if result["code"] != successCode {
			t.Fatalf("agent sync %s = %#v", syncID, result)
		}
	}

	// 1) 新建启用中的入站
	syncInbound("event-sync-1", "2026-10-02T00:00:00Z", true)

	var userID string
	if err := database.QueryRow(`SELECT id FROM users WHERE display_name = '事件测试入站'`).Scan(&userID); err != nil {
		t.Fatalf("read synced user: %v", err)
	}
	open, dismissed := profileEventCounts(t, database, userID)
	if open != 1 || dismissed != 0 {
		t.Fatalf("after first sync open=%d dismissed=%d, want 1/0", open, dismissed)
	}

	// 2) 入站被停用：事件应立即关闭，不能停在待处理
	syncInbound("event-sync-2", "2026-10-02T00:01:00Z", false)
	open, dismissed = profileEventCounts(t, database, userID)
	if open != 0 {
		t.Fatalf("open events after inbound disabled = %d, want 0", open)
	}
	if dismissed != 1 {
		t.Fatalf("dismissed events after inbound disabled = %d, want 1", dismissed)
	}
	var requiresAction int
	if err := database.QueryRow(`SELECT requires_action FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ?`, userID).Scan(&requiresAction); err != nil {
		t.Fatalf("read dismissed event action flag: %v", err)
	}
	if requiresAction != 0 {
		t.Fatalf("dismissed event still requires_action = %d", requiresAction)
	}

	// 3) 入站恢复启用：用户仍未定价，事件应重新提醒
	syncInbound("event-sync-3", "2026-10-02T00:02:00Z", true)
	open, dismissed = profileEventCounts(t, database, userID)
	if open != 1 {
		t.Fatalf("open events after inbound re-enabled = %d, want 1", open)
	}
	if dismissed != 1 {
		t.Fatalf("dismissed history changed to %d, want 1", dismissed)
	}

	// 4) 完成定价后再次停用：不再产生待办事件
	if updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/users/"+userID, token, map[string]any{
		"billingAmount": 99,
	}); updated["code"] != successCode {
		t.Fatalf("update billing = %#v", updated)
	}
	syncInbound("event-sync-4", "2026-10-02T00:03:00Z", false)
	syncInbound("event-sync-5", "2026-10-02T00:04:00Z", true)
	open, _ = profileEventCounts(t, database, userID)
	if open != 0 {
		t.Fatalf("priced user still has %d open profile events, want 0", open)
	}
	_ = nodeID
}

// TestNewUserProfileEventTracksInboundRename 确认 Inbound 改名后事件文案同步更新，
// 事件中心不会一直显示旧名称。
func TestNewUserProfileEventTracksInboundRename(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"name": "改名事件线路机", "type": "relay",
	})
	nodeData := created["data"].(map[string]any)
	nodeKey := nodeData["nodeKey"].(string)
	agentToken := nodeData["token"].(string)

	syncTag := func(syncID, observedAt, tag, remark string) {
		payload := map[string]any{
			"node_key":    nodeKey,
			"sync_id":     syncID,
			"observed_at": observedAt,
			"status": map[string]any{
				"agent_version": "v1.0.4",
				"xray_running":  true,
			},
			"inbounds": []any{
				map[string]any{
					"remote_id": 8802, "tag": tag, "remark": remark,
					"protocol": "vless", "port": 443, "listen": "0.0.0.0", "enable": true,
					"expiry_time": 1900000000, "up": 10, "down": 20, "all_time": 30,
					"clients": []any{
						map[string]any{
							"remote_id": "rename-event-client", "email": "rename@example.com", "enable": true,
							"expiry_time": 1900000000, "up": 10, "down": 20, "all_time": 30,
						},
					},
				},
			},
		}
		result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, payload)
		if result["code"] != successCode {
			t.Fatalf("agent sync %s = %#v", syncID, result)
		}
	}

	syncTag("rename-event-1", "2026-10-03T00:00:00Z", "before-tag", "改名前")
	var userID string
	if err := database.QueryRow(`SELECT id FROM users WHERE display_name = '改名前'`).Scan(&userID); err != nil {
		t.Fatalf("read user: %v", err)
	}

	syncTag("rename-event-2", "2026-10-03T00:01:00Z", "after-tag", "改名后")

	var message string
	if err := database.QueryRow(`SELECT message FROM node_events
WHERE event_type = 'new_user_profile_required' AND resource_id = ?
  AND COALESCE(event_status, 'open') NOT IN ('resolved', 'dismissed')`, userID).Scan(&message); err != nil {
		t.Fatalf("read event message: %v", err)
	}
	if message == "" {
		t.Fatal("event message is empty")
	}
	want := "改名后"
	if !strings.Contains(message, want) {
		t.Fatalf("event message %q does not mention renamed user %q", message, want)
	}
	if strings.Contains(message, "改名前") {
		t.Fatalf("event message still shows the old name: %q", message)
	}
}
