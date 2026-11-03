package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// agentSyncPayload 构造一个最小可用的 Agent 同步请求体。
func agentSyncPayload(nodeKey, syncID, observedAt, tag, remark string, remoteInboundID int64, inboundID string) map[string]any {
	return map[string]any{
		"node_key":    nodeKey,
		"sync_id":     syncID,
		"observed_at": observedAt,
		"status": map[string]any{
			"agent_version": "v1.0.4",
			"xray_running":  true,
			"xray_version":  "1.8.0",
		},
		"inbounds": []any{
			map[string]any{
				"remote_id":   remoteInboundID,
				"tag":         tag,
				"remark":      remark,
				"protocol":    "vless",
				"port":        443,
				"listen":      "0.0.0.0",
				"enable":      true,
				"expiry_time": 1900000000,
				"up":          100,
				"down":        200,
				"all_time":    300,
				"clients": []any{
					map[string]any{
						"remote_id":   "c-" + inboundID,
						"email":       "user@example.com",
						"enable":      true,
						"expiry_time": 1900000000,
						"up":          100,
						"down":        200,
						"all_time":    300,
					},
				},
			},
		},
	}
}

// TestAgentSyncFollowsInboundRenameUntilOperatorLocksName 覆盖两种情况：
// Inbound 在 X-Panel 改名后，未人工编辑过的用户名跟随更新；
// 运营在面板里改过名称后，后续同步不再覆盖。
func TestAgentSyncFollowsInboundRenameUntilOperatorLocksName(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	// 创建节点并取回 Agent Token
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"name": "改名测试线路机", "type": "relay",
	})
	if created["code"] != successCode {
		t.Fatalf("create node = %#v", created)
	}
	nodeData := created["data"].(map[string]any)
	nodeID := nodeData["nodeId"].(string)
	nodeKey := nodeData["nodeKey"].(string)
	agentToken := nodeData["token"].(string)

	syncOnce := func(syncID, tag, remark string) {
		payload := agentSyncPayload(nodeKey, syncID, "2026-10-01T00:00:00Z", tag, remark, 9001, nodeID)
		result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, payload)
		if result["code"] != successCode {
			t.Fatalf("agent sync %s = %#v", syncID, result)
		}
	}

	// 1) 首次同步：用户名取自 Inbound 备注
	syncOnce("rename-sync-1", "inbound-tag-a", "原始名称")
	var userID, displayName string
	if err := database.QueryRow(`SELECT id, display_name FROM users WHERE display_name = '原始名称'`).Scan(&userID, &displayName); err != nil {
		t.Fatalf("read first user name: %v", err)
	}

	// 2) X-Panel 改名后，名称应跟随
	syncOnce("rename-sync-2", "inbound-tag-b", "改名后的名称")
	if err := database.QueryRow(`SELECT display_name FROM users WHERE id = ?`, userID).Scan(&displayName); err != nil {
		t.Fatalf("read renamed user: %v", err)
	}
	if displayName != "改名后的名称" {
		t.Fatalf("display_name = %q, want 改名后的名称", displayName)
	}

	// 3) 运营在面板里手工改名：应加锁
	updated := doJSON(t, ts.Client(), http.MethodPatch, ts.URL+"/api/users/"+userID, token, map[string]any{
		"displayName": "张三-月付",
	})
	if updated["code"] != successCode {
		t.Fatalf("update user = %#v", updated)
	}
	var locked int
	if err := database.QueryRow(`SELECT display_name_locked FROM users WHERE id = ?`, userID).Scan(&locked); err != nil {
		t.Fatalf("read lock flag: %v", err)
	}
	if locked != 1 {
		t.Fatalf("display_name_locked = %d, want 1", locked)
	}

	// 4) 再次同步改名：不应覆盖运营手工名称
	syncOnce("rename-sync-3", "inbound-tag-c", "系统再次改名")
	if err := database.QueryRow(`SELECT display_name FROM users WHERE id = ?`, userID).Scan(&displayName); err != nil {
		t.Fatalf("read locked user: %v", err)
	}
	if displayName != "张三-月付" {
		t.Fatalf("locked display_name = %q, want 张三-月付", displayName)
	}
}

// TestNodeDetailKeepsMissingInboundVisible 确认 Agent 连续漏报 1、2 次时，
// 节点详情仍能看到该 Inbound 并标记为缺失中，避免出现「入站消失但用户还在」。
func TestNodeDetailKeepsMissingInboundVisible(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/nodes", token, map[string]any{
		"name": "缺失测试线路机", "type": "relay",
	})
	nodeData := created["data"].(map[string]any)
	nodeID := nodeData["nodeId"].(string)
	nodeKey := nodeData["nodeKey"].(string)
	agentToken := nodeData["token"].(string)

	// 首次同步：含 Inbound
	payload := agentSyncPayload(nodeKey, "missing-sync-1", "2026-10-01T00:00:00Z", "missing-tag", "缺失测试入站", 9002, nodeID)
	if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, payload); result["code"] != successCode {
		t.Fatalf("initial sync = %#v", result)
	}

	// 后续两次同步：不含任何 Inbound（模拟在 X-Panel 删除）
	for _, syncID := range []string{"missing-sync-2", "missing-sync-3"} {
		empty := map[string]any{
			"node_key":    nodeKey,
			"sync_id":     syncID,
			"observed_at": "2026-10-01T00:0" + syncID[len(syncID)-1:] + ":00Z",
			"status": map[string]any{
				"agent_version": "v1.0.4",
				"xray_running":  true,
			},
			"inbounds": []any{},
		}
		if result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/agent/v1/sync", agentToken, empty); result["code"] != successCode {
			t.Fatalf("empty sync %s = %#v", syncID, result)
		}
	}

	detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/nodes/"+nodeID, token, nil)
	if detail["code"] != successCode {
		t.Fatalf("node detail = %#v", detail)
	}
	inbounds := detail["data"].(map[string]any)["inbounds"].([]any)
	if len(inbounds) != 1 {
		t.Fatalf("missing inbound hidden: got %d inbounds, want 1", len(inbounds))
	}
	first := inbounds[0].(map[string]any)
	if first["status"] != "missing" {
		t.Fatalf("inbound status = %v, want missing", first["status"])
	}
	if first["pendingArchive"] != true {
		t.Fatalf("pendingArchive = %v, want true", first["pendingArchive"])
	}
	if count, ok := first["missingSyncCount"].(float64); !ok || count < 1 {
		t.Fatalf("missingSyncCount = %v, want >= 1", first["missingSyncCount"])
	}

	// 用户此时应仍然存在（还没到归档阈值），与页面显示一致
	var remaining int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&remaining); err != nil {
		t.Fatalf("count users: %v", err)
	}
	if remaining != 1 {
		t.Fatalf("active users = %d, want 1", remaining)
	}
}
