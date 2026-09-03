package central

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

	"xpanel-central/agent/internal/retry"
)

func TestPostJSONSendsBearerAndRequestID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/agent/v1/heartbeat" {
			t.Errorf("request = %s %s", request.Method, request.URL.Path)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer node-token" {
			t.Errorf("Authorization = %q", got)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		if strings.TrimSpace(request.Header.Get("X-Request-ID")) == "" {
			t.Error("X-Request-ID is empty")
		}
		if strings.TrimSpace(request.Header.Get("X-Request-Timestamp")) == "" {
			t.Error("X-Request-Timestamp is empty")
		}
		var body map[string]string
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body["node_key"] != "relay-001" {
			t.Errorf("node_key = %q", body["node_key"])
		}
		writeEnvelope(writer, `{"accepted":true}`)
	}))
	defer server.Close()

	client, err := NewClient(server.URL+"/api", "node-token", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Accepted bool `json:"accepted"`
	}
	if err := client.PostJSON(context.Background(), "/agent/v1/heartbeat", map[string]string{"node_key": "relay-001"}, &result); err != nil {
		t.Fatalf("PostJSON() error = %v", err)
	}
	if !result.Accepted {
		t.Fatal("result.Accepted = false")
	}
}

func TestPostJSONReturnsCentralError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writeEnvelopeWithCode(writer, "E_SYNC", "sync rejected", "null")
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = client.PostJSON(context.Background(), "/agent/v1/sync", struct{}{}, nil)
	if err == nil || !strings.Contains(err.Error(), "sync rejected") {
		t.Fatalf("PostJSON() error = %v, want central error", err)
	}
}

func TestNewClientRejectsInvalidURL(t *testing.T) {
	if _, err := NewClient("central.example.test", "", time.Second); err == nil {
		t.Fatal("NewClient() error = nil, want invalid URL error")
	}
}

func TestPostJSONRetriesTransientCentralFailureWithSameRequestID(t *testing.T) {
	var calls atomic.Int32
	var requestIDs []string
	var requestIDsMu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		call := calls.Add(1)
		requestIDsMu.Lock()
		requestIDs = append(requestIDs, request.Header.Get("X-Request-ID"))
		requestIDsMu.Unlock()
		if call < 3 {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeEnvelope(writer, `{"accepted":true}`)
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "node-token", time.Second, retry.Policy{
		MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Accepted bool `json:"accepted"`
	}
	if err := client.PostJSON(context.Background(), "/agent/v1/heartbeat", map[string]string{"node_key": "relay-001"}, &result); err != nil {
		t.Fatalf("PostJSON() error = %v", err)
	}
	if !result.Accepted || calls.Load() != 3 {
		t.Fatalf("accepted=%v calls=%d, want accepted=true and 3 calls", result.Accepted, calls.Load())
	}
	requestIDsMu.Lock()
	defer requestIDsMu.Unlock()
	if len(requestIDs) != 3 || requestIDs[0] == "" || requestIDs[0] != requestIDs[1] || requestIDs[1] != requestIDs[2] {
		t.Fatalf("request IDs = %v, want one logical request ID", requestIDs)
	}
}

func TestPostJSONDoesNotRetryApplicationFailure(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls++
		writeEnvelopeWithCode(writer, "8888", "not authorized", "null")
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "expired-token", time.Second, retry.Policy{
		MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.PostJSON(context.Background(), "/agent/v1/sync", struct{}{}, nil)
	if err == nil || !strings.Contains(err.Error(), "not authorized") {
		t.Fatalf("PostJSON() error = %v, want application error", err)
	}
	if calls != 1 {
		t.Fatalf("application failure calls = %d, want 1", calls)
	}
}

func TestPostJSONStopsRetryWhenContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "node-token", time.Second, retry.Policy{
		MaxAttempts: 3, InitialBackoff: time.Hour, MaxBackoff: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = client.PostJSON(ctx, "/agent/v1/heartbeat", struct{}{}, nil)
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("PostJSON() error = %v, want context cancellation", err)
	}
}

func TestPostJSONRetriesHTTPTimeout(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		time.Sleep(80 * time.Millisecond)
		writeEnvelope(writer, `{"accepted":true}`)
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "node-token", 15*time.Millisecond, retry.Policy{
		MaxAttempts: 2, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.PostJSON(context.Background(), "/agent/v1/heartbeat", struct{}{}, nil)
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "timeout") {
		t.Fatalf("PostJSON() error = %v, want timeout", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("timeout attempts = %d, want 2", calls.Load())
	}
}

func writeEnvelope(writer http.ResponseWriter, data string) {
	writeEnvelopeWithCode(writer, "0000", "ok", data)
}

func writeEnvelopeWithCode(writer http.ResponseWriter, code, msg, data string) {
	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write([]byte(`{"code":"` + code + `","msg":"` + msg + `","data":` + data + `}`))
}
