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

// TestBackupEndpointSurvivesShortServerWriteTimeout 防止回退：
// http.Server 的 WriteTimeout 会对整库 VACUUM INTO 生效，大库备份会被提前切断。
// 备份接口必须自行放宽写截止时间，否则真实部署里大库备份永远失败。
func TestBackupEndpointSurvivesShortServerWriteTimeout(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "panel.db")
	database, err := db.Open(databasePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}

	// 造出足够大的库，让 VACUUM INTO 明显超过下面设置的 60ms 写超时。
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at)
VALUES ('slow-node', 'slow-node', '超时测试节点', 'relay', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO inbounds (id, node_id, remote_inbound_id, kind, tag, first_seen_at, last_seen_at)
VALUES ('slow-inbound', 'slow-node', '1', 'user', 'slow-inbound', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	// 用递归 CTE 批量写入快照，避免在测试里跑几十万次单行插入。
	if _, err := database.Exec(`
WITH RECURSIVE seq(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 40000)
INSERT INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
SELECT 'snap-' || n, 'slow-node', 'slow-inbound',
       strftime('%Y-%m-%dT%H:%M:%SZ', '2026-01-01T00:00:00Z', '+' || n || ' seconds'),
       n, n * 2, n * 3, 'xpanel', 0
FROM seq`); err != nil {
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

	ts := httptest.NewUnstartedServer(server.Handler())
	// 60ms 远小于一次整库备份的耗时，模拟生产上 30s WriteTimeout 被大库突破的情形。
	ts.Config.WriteTimeout = 60 * time.Millisecond
	ts.Start()
	defer ts.Close()

	login := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/auth/login", "", map[string]string{
		"userName": "admin", "password": "test-password",
	})
	token := login["data"].(map[string]any)["token"].(string)

	started := time.Now()
	created := doJSON(t, ts.Client(), http.MethodPost, ts.URL+"/api/system/backups", token, nil)
	elapsed := time.Since(started)
	if created["code"] != successCode {
		t.Fatalf("backup under short WriteTimeout failed after %s: %#v", elapsed, created)
	}
	if elapsed < 60*time.Millisecond {
		t.Skipf("数据库太小，备份仅用 %s，未触发写超时；跳过（不影响结论）", elapsed)
	}
	name, _ := created["data"].(map[string]any)["name"].(string)
	if name == "" {
		t.Fatalf("created backup has no name: %#v", created)
	}
}
