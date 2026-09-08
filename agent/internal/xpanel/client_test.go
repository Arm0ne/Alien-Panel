package xpanel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"xpanel-central/agent/internal/retry"
)

func TestClientLoginAndGetRetriesOnceAfterUnauthorized(t *testing.T) {
	var loginCalls atomic.Int32
	var getCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/xpanel/login":
			loginCalls.Add(1)
			if request.Method != http.MethodPost {
				t.Errorf("login method = %s, want POST", request.Method)
			}
			if got := request.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
				t.Errorf("login content type = %q", got)
			}
			if err := request.ParseForm(); err != nil {
				t.Errorf("ParseForm() error = %v", err)
			}
			if request.Form.Get("username") != "admin" || request.Form.Get("password") != "secret" {
				t.Errorf("unexpected login form: %v", request.Form)
			}
			http.SetCookie(writer, &http.Cookie{Name: "session", Value: "session-value", Path: "/xpanel"})
			writeJSON(writer, Response{Success: true, Msg: "ok"})
		case "/xpanel/panel/api/inbounds/list":
			getCalls.Add(1)
			if _, err := request.Cookie("session"); err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writeJSON(writer, Response{Success: true, Obj: json.RawMessage(`[{"id":15}]`)})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "/xpanel/", "admin", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := client.Get(context.Background(), "/inbounds/list")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if loginCalls.Load() != 1 || getCalls.Load() != 2 {
		t.Fatalf("login calls = %d, GET calls = %d, want 1 and 2", loginCalls.Load(), getCalls.Load())
	}
	if string(payload.Obj) != `[{"id":15}]` {
		t.Fatalf("Obj = %s", payload.Obj)
	}
}

func TestClientLoginAndGetHandlesRedirectToLogin(t *testing.T) {
	var loginCalls atomic.Int32
	var getCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/xpanel/login":
			loginCalls.Add(1)
			if err := request.ParseForm(); err != nil {
				t.Fatalf("ParseForm() error = %v", err)
			}
			http.SetCookie(writer, &http.Cookie{Name: "session", Value: "session-value", Path: "/xpanel"})
			writeJSON(writer, Response{Success: true, Msg: "ok"})
		case "/xpanel/panel/api/inbounds/list":
			getCalls.Add(1)
			if _, err := request.Cookie("session"); err != nil {
				http.Redirect(writer, request, "/xpanel/login", http.StatusFound)
				return
			}
			writeJSON(writer, Response{Success: true, Obj: json.RawMessage(`[{"id":15}]`)})
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "/xpanel/", "admin", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := client.Get(context.Background(), "/inbounds/list")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if loginCalls.Load() != 1 || getCalls.Load() != 2 {
		t.Fatalf("login calls = %d, GET calls = %d, want 1 and 2", loginCalls.Load(), getCalls.Load())
	}
	if string(payload.Obj) != `[{"id":15}]` {
		t.Fatalf("Obj = %s", payload.Obj)
	}
}

func TestClientLoginFailureEntersBackoff(t *testing.T) {
	var loginCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/login" {
			loginCalls.Add(1)
			writeJSON(writer, Response{Success: false, Msg: "invalid credentials"})
			return
		}
		writer.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client, err := NewClient(server.URL, "/", "admin", "wrong", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	firstErr := client.Login(context.Background())
	if firstErr == nil {
		t.Fatal("Login() error = nil, want authentication error")
	}
	secondErr := client.Login(context.Background())
	if secondErr == nil || secondErr.Error() == firstErr.Error() {
		t.Fatalf("second Login() error = %v, want distinct backoff error", secondErr)
	}
	if loginCalls.Load() != 1 {
		t.Fatalf("login calls = %d, want 1 while backoff is active", loginCalls.Load())
	}
}

func TestNewClientRejectsInvalidURL(t *testing.T) {
	if _, err := NewClient("127.0.0.1:2053", "/", "admin", "secret", time.Second); err == nil {
		t.Fatal("NewClient() error = nil, want invalid URL error")
	}
}

func TestClientRetriesTransientGetFailure(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/panel/api/server/status" {
			http.NotFound(writer, request)
			return
		}
		if calls.Add(1) < 3 {
			writer.WriteHeader(http.StatusBadGateway)
			return
		}
		writeJSON(writer, Response{Success: true, Obj: json.RawMessage(`{"xray_running":true}`)})
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "/", "admin", "secret", time.Second, retry.Policy{
		MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Get(context.Background(), "/server/status")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if calls.Load() != 3 || !response.Success {
		t.Fatalf("calls=%d response=%#v, want 3 successful attempts", calls.Load(), response)
	}
}

func TestClientDoesNotRetryApplicationFailure(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		writeJSON(writer, Response{Success: false, Msg: "permission denied"})
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "/", "admin", "secret", time.Second, retry.Policy{
		MaxAttempts: 3, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), "/server/status")
	if err == nil || !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("Get() error = %v, want application error", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("application failure calls = %d, want 1", calls.Load())
	}
}

func TestClientRetriesHTTPTimeout(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		time.Sleep(80 * time.Millisecond)
		writeJSON(writer, Response{Success: true, Obj: json.RawMessage(`{"ok":true}`)})
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "/", "admin", "secret", 15*time.Millisecond, retry.Policy{
		MaxAttempts: 2, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Get(context.Background(), "/server/status")
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "timeout") {
		t.Fatalf("Get() error = %v, want timeout", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("timeout attempts = %d, want 2", calls.Load())
	}
}

func TestClientRetriesLoginTimeout(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/login" {
			http.NotFound(writer, request)
			return
		}
		calls.Add(1)
		time.Sleep(80 * time.Millisecond)
		writeJSON(writer, Response{Success: true})
	}))
	defer server.Close()

	client, err := NewClientWithRetryPolicy(server.URL, "/", "admin", "secret", 15*time.Millisecond, retry.Policy{
		MaxAttempts: 2, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Login(context.Background())
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "timeout") {
		t.Fatalf("Login() error = %v, want timeout", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("login timeout attempts = %d, want 2", calls.Load())
	}
}

func writeJSON(writer http.ResponseWriter, value Response) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(value)
}
