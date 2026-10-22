package httpapi

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"xpanel-central/backend/internal/config"
	"xpanel-central/backend/internal/db"
)

// TestBackupManageLifecycle covers the manual backup workflow: create, list,
// verify, delete an older file, and refuse to delete the newest one.
func TestBackupManageLifecycle(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "panel.db")
	database, err := db.Open(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	server, err := NewServer(config.Config{
		AdminUsername: "admin", AdminPassword: "test-password", DatabasePath: databasePath,
		SessionTTL: time.Hour, CorsOrigins: []string{"http://localhost:9527"},
	}, database, slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })

	ts := httptest.NewServer(server.Handler())
	defer ts.Close()
	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	empty := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/system/backups", token, nil)
	if empty["code"] != successCode || len(empty["data"].(map[string]any)["items"].([]any)) != 0 {
		t.Fatalf("initial list = %#v", empty)
	}

	first := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/system/backups", token, nil)
	if first["code"] != successCode {
		t.Fatalf("create backup = %#v", first)
	}
	firstName := first["data"].(map[string]any)["name"].(string)
	if firstName == "" {
		t.Fatal("created backup has no name")
	}

	list := doJSON(t, ts.Client(), http.MethodGet, ts.URL+"/api/system/backups", token, nil)
	items := list["data"].(map[string]any)["items"].([]any)
	if len(items) != 1 || list["data"].(map[string]any)["latest"] != firstName {
		t.Fatalf("list after create = %#v", list)
	}

	verified := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/system/backups/verify", token, map[string]any{"name": firstName})
	if verified["code"] != successCode || verified["data"].(map[string]any)["ok"] != true {
		t.Fatalf("verify = %#v", verified)
	}

	if status, _ := doJSONWithStatus(t, ts.Client(), http.MethodPost, ts.URL+"/api/system/backups/verify", token, map[string]any{"name": "../" + firstName}); status != http.StatusBadRequest {
		t.Fatalf("traversal verify status = %d", status)
	}

	time.Sleep(1100 * time.Millisecond)
	second := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/system/backups", token, nil)
	secondName := second["data"].(map[string]any)["name"].(string)
	if secondName == firstName {
		t.Fatalf("second backup reused name %s", firstName)
	}

	deleted := doJSON(t, ts.Client(), http.MethodDelete, ts.URL+"/api/system/backups/"+firstName, token, nil)
	if deleted["code"] != successCode {
		t.Fatalf("delete older backup = %#v", deleted)
	}

	if status, _ := doJSONWithStatus(t, ts.Client(), http.MethodDelete, ts.URL+"/api/system/backups/"+secondName, token, nil); status != http.StatusConflict {
		t.Fatalf("delete newest backup status = %d, want 409", status)
	}
}
