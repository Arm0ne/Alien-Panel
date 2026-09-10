package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"xpanel-central/backend/internal/config"
	"xpanel-central/backend/internal/db"
)

func TestSystemBackupDownloadAndRestore(t *testing.T) {
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

	backupPath := filepath.Join(t.TempDir(), "source.sqlite3")
	if _, err := db.BackupTo(t.Context(), databasePath, backupPath); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(server.Handler())
	defer ts.Close()

	unauthorized, err := http.Get(ts.URL + "/api/system/backups/download")
	if err != nil {
		t.Fatal(err)
	}
	if unauthorized.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d", unauthorized.StatusCode)
	}
	_ = unauthorized.Body.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)
	download := requestWithToken(t, ts.Client(), http.MethodGet, ts.URL+"/api/system/backups/download", token, nil)
	if download.StatusCode != http.StatusOK {
		t.Fatalf("download status = %d", download.StatusCode)
	}
	body, err := io.ReadAll(download.Body)
	_ = download.Body.Close()
	if err != nil || !bytes.HasPrefix(body, []byte("SQLite format 3\x00")) {
		t.Fatalf("download is not sqlite: err=%v header=%q", err, body[:min(len(body), 16)])
	}

	var form bytes.Buffer
	writer := multipart.NewWriter(&form)
	part, err := writer.CreateFormFile("file", filepath.Base(backupPath))
	if err != nil {
		t.Fatal(err)
	}
	source, err := os.Open(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.Copy(part, source)
	_ = source.Close()
	_ = writer.Close()
	restoreRequest, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/system/restore", &form)
	restoreRequest.Header.Set("Authorization", "Bearer "+token)
	restoreRequest.Header.Set("Content-Type", writer.FormDataContentType())
	restoreResponse, err := ts.Client().Do(restoreRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer restoreResponse.Body.Close()
	var envelope struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(restoreResponse.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if restoreResponse.StatusCode != http.StatusOK || envelope.Code != successCode {
		t.Fatalf("restore response = %d %#v", restoreResponse.StatusCode, envelope)
	}
	oldSession := requestWithToken(t, ts.Client(), http.MethodGet, ts.URL+"/api/auth/me", token, nil)
	if oldSession.StatusCode != http.StatusUnauthorized {
		t.Fatalf("restored database should invalidate old session, status=%d", oldSession.StatusCode)
	}
}

func requestWithToken(t *testing.T, client *http.Client, method, url, token string, body io.Reader) *http.Response {
	t.Helper()
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
