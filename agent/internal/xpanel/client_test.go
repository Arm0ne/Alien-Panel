package xpanel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
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

func writeJSON(writer http.ResponseWriter, value Response) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(value)
}
