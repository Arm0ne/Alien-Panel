package runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"xpanel-central/agent/internal/central"
	"xpanel-central/agent/internal/collector"
	"xpanel-central/agent/internal/xpanel"
)

func TestRunCollectsAndSendsHeartbeatAndSync(t *testing.T) {
	panel := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/panel/api/inbounds/list":
			_, _ = writer.Write([]byte(`{"success":true,"obj":[{"id":15,"tag":"user-15","enable":true}]}`))
		case "/panel/api/server/status":
			_, _ = writer.Write([]byte(`{"success":true,"obj":{"xray_running":true}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer panel.Close()

	var mu sync.Mutex
	paths := make([]string, 0, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	centralServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mu.Lock()
		paths = append(paths, request.URL.Path)
		mu.Unlock()
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"code":"0000","msg":"ok","data":{"status":"accepted"}}`))
		if request.URL.Path == "/agent/v1/sync" {
			cancel()
		}
	}))
	defer centralServer.Close()

	xpanelClient, err := xpanel.NewClient(panel.URL, "/", "admin", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	snapshotCollector, err := collector.New(xpanelClient, "relay-001")
	if err != nil {
		t.Fatal(err)
	}
	centralClient, err := central.NewClient(centralServer.URL, "node-token", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	agentRunner, err := New(snapshotCollector, centralClient, time.Hour, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := agentRunner.Run(ctx); err != context.Canceled {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(paths) != 2 || paths[0] != "/agent/v1/heartbeat" || paths[1] != "/agent/v1/sync" {
		t.Fatalf("central paths = %v", paths)
	}
}
