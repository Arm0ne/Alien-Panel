package collector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"xpanel-central/agent/internal/xpanel"
)

func TestParseInboundsSupportsSettingsClientsAndTraffic(t *testing.T) {
	payload := json.RawMessage(`[
      {
        "id": 15,
        "tag": "user-15",
        "remark": "Customer A",
        "protocol": "vless",
        "port": "443",
        "enable": true,
        "expiryTime": 1792022400,
        "up": 123,
        "down": 456,
        "all_time": 579,
        "settings": {"clients": [
          {"id": "client-a", "email": "phone", "enable": true, "expiryTime": 1792022400},
          {"id": "client-b", "email": "laptop", "enable": false}
        ]},
        "clientTraffics": [
          {"id": "client-a", "email": "phone", "up": 10, "down": 20, "allTime": 30, "lastOnline": 1791000000}
        ]
      }
    ]`)
	inbounds, err := ParseInbounds(payload)
	if err != nil {
		t.Fatalf("ParseInbounds() error = %v", err)
	}
	if len(inbounds) != 1 {
		t.Fatalf("inbound count = %d", len(inbounds))
	}
	inbound := inbounds[0]
	if inbound.RemoteID != 15 || inbound.Port != 443 || inbound.AllTime != 579 {
		t.Fatalf("inbound = %+v", inbound)
	}
	if len(inbound.Clients) != 2 {
		t.Fatalf("client count = %d", len(inbound.Clients))
	}
	if inbound.Clients[0].RemoteID != "client-a" || inbound.Clients[0].AllTime != 30 || inbound.Clients[0].LastOnline != 1791000000 {
		t.Fatalf("client = %+v", inbound.Clients[0])
	}
	if inbound.ConfigHash == "" || len(inbound.ConfigHash) != 64 {
		t.Fatalf("ConfigHash = %q", inbound.ConfigHash)
	}
}

func TestParseInboundsSupportsWrappedPayloadAndStringSettings(t *testing.T) {
	payload := json.RawMessage(`{"items":[{"id":1,"settings":"{\"clients\":[{\"email\":\"phone\"}]}"}]}`)
	inbounds, err := ParseInbounds(payload)
	if err != nil {
		t.Fatalf("ParseInbounds() error = %v", err)
	}
	if len(inbounds) != 1 || len(inbounds[0].Clients) != 1 || inbounds[0].Clients[0].RemoteID != "phone" {
		t.Fatalf("inbounds = %+v", inbounds)
	}
}

func TestParseStatusSupportsNestedXrayObject(t *testing.T) {
	status, err := ParseStatus(json.RawMessage(`{"xray":{"state":"running","version":"1.2.3"},"cpuUsage":"12.5","memoryUsed":100,"memoryTotal":200}`))
	if err != nil {
		t.Fatalf("ParseStatus() error = %v", err)
	}
	if !status.XrayRunning || status.XrayVersion != "1.2.3" || status.CPUUsage != 12.5 || status.MemoryUsed != 100 || status.MemoryTotal != 200 {
		t.Fatalf("status = %+v", status)
	}
}

func TestCollectorCollectsInboundAndStatusWithoutReset(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		paths = append(paths, request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/panel/api/inbounds/list":
			_, _ = writer.Write([]byte(`{"success":true,"obj":[{"id":15,"tag":"user-15","protocol":"vless","enable":true,"up":1,"down":2}]}`))
		case "/panel/api/server/status":
			_, _ = writer.Write([]byte(`{"success":true,"obj":{"xray_running":true,"xray_version":"1.2.3"}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := xpanel.NewClient(server.URL, "/", "admin", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	collector, err := New(client, "relay-001")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := collector.Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if snapshot.NodeKey != "relay-001" || snapshot.SyncID == "" || snapshot.ObservedAt == "" || len(snapshot.Inbounds) != 1 || !snapshot.Status.XrayRunning {
		t.Fatalf("snapshot = %+v", snapshot)
	}
	if strings.Contains(strings.Join(paths, ","), "reset") {
		t.Fatalf("collector called reset endpoint: %v", paths)
	}
}
