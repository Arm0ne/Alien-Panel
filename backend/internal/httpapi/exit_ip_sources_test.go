package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExitIPSourcesAndRouteScopes(t *testing.T) {
	server, database := testServer(t)
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	now := time.Now().UTC().Format(time.RFC3339Nano)
	for _, node := range []struct {
		id, key, name, kind string
	}{
		{"source-relay", "source-relay", "来源线路机", "relay"},
		{"source-landing", "source-landing", "来源落地机", "landing"},
	} {
		if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, enabled, health_status, created_at, updated_at) VALUES (?, ?, ?, ?, 1, 'online', ?, ?)`, node.id, node.key, node.name, node.kind, now, now); err != nil {
			t.Fatalf("insert node %s: %v", node.id, err)
		}
	}
	if _, err := database.Exec(`INSERT INTO routes (id, name, relay_node_id, landing_node_id, enabled, created_at, updated_at) VALUES ('source-route', '来源线路', 'source-relay', 'source-landing', 1, ?, ?)`, now, now); err != nil {
		t.Fatalf("insert route: %v", err)
	}
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{"userName": "admin", "password": "test-password"})
	token := login["data"].(map[string]any)["token"].(string)
	create := func(payload map[string]any) string {
		result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/exit-ips", token, payload)
		if result["code"] != successCode {
			t.Fatalf("create exit IP failed: %#v", result)
		}
		return result["data"].(map[string]any)["id"].(string)
	}
	relayIP := create(map[string]any{"address": "198.51.100.71", "sourceType": "node", "ownerNodeId": "source-relay"})
	s5IP := create(map[string]any{"address": "198.51.100.72", "sourceType": "s5", "provider": "S5 Supplier"})
	landingIP := create(map[string]any{"address": "198.51.100.73", "sourceType": "node", "ownerNodeId": "source-landing"})
	assertBinding := func(exitID, scope string) {
		result := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes/source-route/exit-ips", token, map[string]any{"exitIpId": exitID, "scope": scope})
		if result["code"] != successCode || result["data"].(map[string]any)["scope"] != scope {
			t.Fatalf("bind %s/%s failed: %#v", exitID, scope, result)
		}
	}
	assertBinding(relayIP, "relay")
	assertBinding(s5IP, "external")
	assertBinding(landingIP, "landing")
	if status, result := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/routes/source-route/exit-ips", token, map[string]any{"exitIpId": s5IP, "scope": "landing"}); status != http.StatusBadRequest || result["code"] != validationCode {
		t.Fatalf("S5 landing binding status=%d response=%#v", status, result)
	}
	relayDetail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/"+relayIP, token, nil)
	if relayDetail["data"].(map[string]any)["ownerNodeType"] != "relay" || relayDetail["data"].(map[string]any)["landingNodeId"] != nil {
		t.Fatalf("relay exit detail compatibility fields are wrong: %#v", relayDetail)
	}
	s5Detail := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/exit-ips/"+s5IP, token, nil)
	if s5Detail["data"].(map[string]any)["sourceType"] != "s5" || s5Detail["data"].(map[string]any)["ownerNodeId"] != nil {
		t.Fatalf("S5 exit detail is wrong: %#v", s5Detail)
	}
}
