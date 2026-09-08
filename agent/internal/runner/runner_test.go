package runner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"xpanel-central/agent/internal/central"
	"xpanel-central/agent/internal/collector"
	"xpanel-central/agent/internal/retry"
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
	versions := make([]string, 0, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	centralServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mu.Lock()
		paths = append(paths, request.URL.Path)
		var payload struct {
			Status struct {
				AgentVersion string `json:"agent_version"`
			} `json:"status"`
		}
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Errorf("decode %s payload: %v", request.URL.Path, err)
		} else {
			versions = append(versions, payload.Status.AgentVersion)
		}
		mu.Unlock()
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"code":"0000","msg":"ok","data":{"status":"accepted"}}`))
		if request.URL.Path == "/agent/v1/sync" {
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
		}
	}))
	defer centralServer.Close()

	xpanelClient, err := xpanel.NewClient(panel.URL, "/", "admin", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	snapshotCollector, err := collector.New(xpanelClient, "relay-001", "test-agent")
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
	if len(versions) != 2 || versions[0] != "test-agent" || versions[1] != "test-agent" {
		t.Fatalf("Agent versions = %v, want both payloads to contain test-agent", versions)
	}
}

func TestRunRetriesXPanelAndCentralTransientFailures(t *testing.T) {
	var inboundCalls atomic.Int32
	var statusCalls atomic.Int32
	panel := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/panel/api/inbounds/list":
			if inboundCalls.Add(1) < 3 {
				writer.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			_, _ = writer.Write([]byte(`{"success":true,"obj":[{"id":15,"tag":"user-15","enable":true}]}`))
		case "/panel/api/server/status":
			if statusCalls.Add(1) < 2 {
				writer.WriteHeader(http.StatusGatewayTimeout)
				return
			}
			_, _ = writer.Write([]byte(`{"success":true,"obj":{"xray_running":true}}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer panel.Close()

	var heartbeatCalls atomic.Int32
	var syncCalls atomic.Int32
	var syncIDsMu sync.Mutex
	var syncIDs []string
	ctx, cancel := context.WithCancel(context.Background())
	centralServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/agent/v1/heartbeat":
			if heartbeatCalls.Add(1) < 3 {
				writer.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		case "/agent/v1/sync":
			if syncCalls.Add(1) < 3 {
				writer.WriteHeader(http.StatusBadGateway)
				return
			}
			var payload struct {
				SyncID string `json:"sync_id"`
			}
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Errorf("decode sync payload: %v", err)
			}
			syncIDsMu.Lock()
			syncIDs = append(syncIDs, payload.SyncID)
			syncIDsMu.Unlock()
		}
		_, _ = writer.Write([]byte(`{"code":"0000","msg":"ok","data":{"status":"success"}}`))
		if request.URL.Path == "/agent/v1/sync" {
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
		}
	}))
	defer centralServer.Close()

	policy := retry.Policy{MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond}
	xpanelClient, err := xpanel.NewClientWithRetryPolicy(panel.URL, "/", "admin", "secret", time.Second, policy)
	if err != nil {
		t.Fatal(err)
	}
	snapshotCollector, err := collector.New(xpanelClient, "relay-retry", "test-agent")
	if err != nil {
		t.Fatal(err)
	}
	centralClient, err := central.NewClientWithRetryPolicy(centralServer.URL, "node-token", time.Second, policy)
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
	if inboundCalls.Load() != 3 || statusCalls.Load() != 2 || heartbeatCalls.Load() != 3 || syncCalls.Load() != 3 {
		t.Fatalf("calls inbound=%d status=%d heartbeat=%d sync=%d; want 3,2,3,3", inboundCalls.Load(), statusCalls.Load(), heartbeatCalls.Load(), syncCalls.Load())
	}
	syncIDsMu.Lock()
	defer syncIDsMu.Unlock()
	if len(syncIDs) != 1 || strings.TrimSpace(syncIDs[0]) == "" {
		t.Fatalf("successful sync IDs = %v, want one non-empty ID", syncIDs)
	}
}

func TestRunStopsRetryBackoffWhenCanceled(t *testing.T) {
	panel := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer panel.Close()

	xpanelClient, err := xpanel.NewClientWithRetryPolicy(panel.URL, "/", "admin", "secret", time.Second, retry.Policy{
		MaxAttempts: 4, InitialBackoff: time.Hour, MaxBackoff: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshotCollector, err := collector.New(xpanelClient, "relay-cancel", "test-agent")
	if err != nil {
		t.Fatal(err)
	}
	centralClient, err := central.NewClientWithRetryPolicy("http://127.0.0.1:1", "node-token", time.Second, retry.Policy{
		MaxAttempts: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	agentRunner, err := New(snapshotCollector, centralClient, time.Hour, nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- agentRunner.Run(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("Run() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run() did not stop after context cancellation")
	}
}
