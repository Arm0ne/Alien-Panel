package collector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	if inbound.RemoteID != 15 || inbound.Port != 443 || inbound.AllTime != 579 || inbound.ExpiryTime != 1792022400 {
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

func TestParseInboundsNormalizesExpiryTimestampVariants(t *testing.T) {
	payload := json.RawMessage(`[
      {"id": 1, "expiryTime": "2026-09-05T00:00:00Z"},
      {"id": 2, "expiry_time": "1788566400000"},
      {"id": 3, "expiryTime": "2026-09-05"}
    ]`)
	inbounds, err := ParseInbounds(payload)
	if err != nil {
		t.Fatalf("ParseInbounds() error = %v", err)
	}
	if len(inbounds) != 3 {
		t.Fatalf("inbound count = %d", len(inbounds))
	}
	want := []int64{1788566400, 1788566400, 1788566400}
	for index, inbound := range inbounds {
		if inbound.ExpiryTime != want[index] {
			t.Fatalf("inbound %d expiry = %d, want %d", inbound.RemoteID, inbound.ExpiryTime, want[index])
		}
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

func TestParseInboundsNormalizesLastOnlineTimestampVariants(t *testing.T) {
	payload := json.RawMessage(`[
      {
        "id": 1,
        "settings": {"clients": [
          {"id": "rfc-client", "email": "rfc", "lastOnlineAt": "2026-09-04T00:00:00Z"}
        ]},
        "clientTraffics": [
          {"id": "rfc-client", "up": 1, "down": 2}
        ]
      },
      {
        "id": 2,
        "clientStats": [
          {"id": "ms-client", "email": "ms", "lastOnline": "1756944000000"}
        ]
      }
    ]`)
	inbounds, err := ParseInbounds(payload)
	if err != nil {
		t.Fatalf("ParseInbounds() error = %v", err)
	}
	if len(inbounds) != 2 || len(inbounds[0].Clients) != 1 || len(inbounds[1].Clients) != 1 {
		t.Fatalf("inbounds = %+v", inbounds)
	}
	if got, want := inbounds[0].Clients[0].LastOnline, int64(1788480000); got != want {
		t.Fatalf("RFC3339 last online = %d, want %d", got, want)
	}
	if got, want := inbounds[1].Clients[0].LastOnline, int64(1756944000); got != want {
		t.Fatalf("millisecond last online = %d, want %d", got, want)
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

func TestCompatibilityFixtures(t *testing.T) {
	tests := []struct {
		name             string
		inboundsFixture  string
		statusFixture    string
		inboundID        int64
		clientID         string
		clientAllTime    int64
		statusXray       string
		statusPanel      string
		statusCPU        float64
		statusMemoryUsed int64
		statusDiskUsed   int64
	}{
		{
			name:             "xpanel-v1-envelope-and-nested-metrics",
			inboundsFixture:  "xpanel-v1-inbounds.json",
			statusFixture:    "xpanel-v1-status.json",
			inboundID:        101,
			clientID:         "v1-client",
			clientAllTime:    30,
			statusXray:       "1.8.4",
			statusPanel:      "2.4.7",
			statusCPU:        6.25,
			statusMemoryUsed: 1048576,
			statusDiskUsed:   2147483648,
		},
		{
			name:             "xpanel-v2-data-envelope-and-string-fields",
			inboundsFixture:  "xpanel-v2-inbounds.json",
			statusFixture:    "xpanel-v2-status.json",
			inboundID:        202,
			clientID:         "v2-laptop",
			clientAllTime:    13,
			statusXray:       "1.8.5",
			statusPanel:      "2.5.0",
			statusCPU:        8.5,
			statusMemoryUsed: 2000,
			statusDiskUsed:   5000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inboundResponse := readResponseFixture(t, test.inboundsFixture)
			inbounds, err := ParseInbounds(responsePayload(inboundResponse))
			if err != nil {
				t.Fatalf("ParseInbounds() error = %v", err)
			}
			if len(inbounds) != 1 || inbounds[0].RemoteID != test.inboundID {
				t.Fatalf("inbounds = %+v", inbounds)
			}
			if len(inbounds[0].Clients) != 1 || inbounds[0].Clients[0].RemoteID != test.clientID || inbounds[0].Clients[0].AllTime != test.clientAllTime {
				t.Fatalf("clients = %+v", inbounds[0].Clients)
			}

			statusResponse := readResponseFixture(t, test.statusFixture)
			status, err := ParseStatus(responsePayload(statusResponse))
			if err != nil {
				t.Fatalf("ParseStatus() error = %v", err)
			}
			if !status.XrayRunning || status.XrayVersion != test.statusXray || status.XPanelVersion != test.statusPanel || status.CPUUsage != test.statusCPU || status.MemoryUsed != test.statusMemoryUsed || status.DiskUsed != test.statusDiskUsed {
				t.Fatalf("status = %+v", status)
			}
		})
	}
}

func readResponseFixture(t *testing.T, name string) xpanel.Response {
	t.Helper()
	path := filepath.Join("testdata", name)
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	var response xpanel.Response
	if err := json.Unmarshal(contents, &response); err != nil {
		t.Fatalf("decode fixture %s: %v", path, err)
	}
	if !response.Success {
		t.Fatalf("fixture %s is not successful: %s", path, response.Msg)
	}
	return response
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
