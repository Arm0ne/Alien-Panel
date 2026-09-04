// Command seed-demo adds a clearly labelled, repeatable dataset to a local
// central database so the user-management and aggregation pages can be tested
// before a real Agent is connected. It never touches non-demo records.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"xpanel-central/backend/internal/db"
)

func main() {
	databasePath := flag.String("database", "./data/panel.db", "SQLite database path")
	flag.Parse()

	database, err := db.Open(*databasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	if err := seed(database); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("demo dataset ready: user=demo-user-001, node=demo-relay-001, database=%s\n", *databasePath)
}

func seed(database *sql.DB) error {
	now := time.Now().UTC().Truncate(time.Second)
	nowText := now.Format(time.RFC3339)
	previousText := now.Add(-24 * time.Hour).Format(time.RFC3339)
	expiryText := now.Add(90 * 24 * time.Hour).Format(time.RFC3339)

	tx, err := database.Begin()
	if err != nil {
		return fmt.Errorf("begin demo seed: %w", err)
	}
	defer tx.Rollback()

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT OR IGNORE INTO nodes (id, node_key, name, type, hostname, public_ip, region, provider, panel_base_path, enabled, health_status, last_seen_at, created_at, updated_at)
VALUES ('demo-relay-001', 'demo-relay-001', '模拟线路机（测试）', 'relay', 'demo-relay.local', '192.0.2.10', '东京', 'Demo Provider', '/', 1, 'online', ?, ?, ?)`, []any{nowText, nowText, nowText}},
		{`INSERT OR IGNORE INTO nodes (id, node_key, name, type, hostname, public_ip, region, provider, panel_base_path, enabled, health_status, last_seen_at, created_at, updated_at)
VALUES ('demo-landing-001', 'demo-landing-001', '模拟落地机（测试）', 'landing', 'demo-landing.local', '192.0.2.20', '东京', 'Demo Provider', '/', 1, 'online', ?, ?, ?)`, []any{nowText, nowText, nowText}},
		{`INSERT OR IGNORE INTO users (id, display_name, status, monthly_fee, currency, expiry_time, notes, created_at, updated_at)
VALUES ('demo-user-001', '模拟用户-测试', 'active', 99, 'CNY', ?, '用于测试用户详情、编辑、设备与流量展示；不代表真实用户。', ?, ?)`, []any{expiryText, nowText, nowText}},
		{`INSERT OR IGNORE INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, remark, protocol, port, listen, enable, expiry_time, up, down, all_time, client_count, config_hash, first_seen_at, last_seen_at)
VALUES ('demo-inbound-001', 'demo-relay-001', 'demo-1001', 'demo-user-001', 'user', 'demo-reality-inbound', '模拟用户 Reality 入站', 'vless', 443, '0.0.0.0', 1, ?, 7340032, 12582912, 19922944, 2, 'demo-config-hash', ?, ?)`, []any{expiryText, previousText, nowText}},
		{`INSERT OR IGNORE INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from)
VALUES ('demo-user-inbound-001', 'demo-user-001', 'demo-inbound-001', 1, ?)`, []any{previousText}},
		{`INSERT OR IGNORE INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, expiry_time, up, down, all_time, last_online, last_seen_at)
VALUES ('demo-client-001', 'demo-relay-001', 'demo-inbound-001', 'demo-device-001', 'demo-phone@example.test', 1, ?, 2097152, 4194304, 6291456, ?, ?)`, []any{expiryText, nowText, nowText}},
		{`INSERT OR IGNORE INTO clients (id, node_id, inbound_id, remote_client_id, email, enable, expiry_time, up, down, all_time, last_online, last_seen_at)
VALUES ('demo-client-002', 'demo-relay-001', 'demo-inbound-001', 'demo-device-002', 'demo-laptop@example.test', 1, ?, 5242880, 8388608, 13631488, ?, ?)`, []any{expiryText, nowText, nowText}},
		{`INSERT OR IGNORE INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES ('demo-snapshot-001', 'demo-relay-001', 'demo-inbound-001', ?, 4194304, 8388608, 12582912, 'demo-seed', 0)`, []any{previousText}},
		{`INSERT OR IGNORE INTO traffic_snapshots (id, node_id, inbound_id, collected_at, up, down, all_time, source, reset_detected)
VALUES ('demo-snapshot-002', 'demo-relay-001', 'demo-inbound-001', ?, 7340032, 12582912, 19922944, 'demo-seed', 0)`, []any{nowText}},
		{`INSERT OR IGNORE INTO routes (id, name, relay_node_id, landing_node_id, relay_outbound_tag, landing_inbound_tag, enabled, valid_from, notes, created_at, updated_at)
VALUES ('demo-route-001', '模拟线路（测试）', 'demo-relay-001', 'demo-landing-001', 'demo-to-landing', 'demo-ss-inbound', 1, date('now'), '用于测试用户线路展示。', ?, ?)`, []any{nowText, nowText}},
		{`INSERT OR IGNORE INTO user_routes (id, user_id, route_id, is_primary, active_from)
VALUES ('demo-user-route-001', 'demo-user-001', 'demo-route-001', 1, ?)`, []any{previousText}},
		{`INSERT OR IGNORE INTO exit_ips (id, landing_node_id, source_type, owner_node_id, ip, family, provider, monthly_cost, currency, enabled, created_at, updated_at)
VALUES ('demo-exit-ip-001', 'demo-landing-001', 'node', 'demo-landing-001', '198.51.100.20', 4, 'Demo IP Provider', 15, 'CNY', 1, ?, ?)`, []any{nowText, nowText}},
		{`INSERT OR IGNORE INTO route_exit_ips (id, route_id, exit_ip_id, scope, allocation_weight, enabled)
VALUES ('demo-route-exit-001', 'demo-route-001', 'demo-exit-ip-001', 'landing', 1, 1)`, nil},
		{`INSERT OR IGNORE INTO user_paths (id, user_id, relay_node_id, landing_node_id, exit_ip_id, mode, notes, active_from, created_at, updated_at)
VALUES ('demo-user-path-001', 'demo-user-001', 'demo-relay-001', 'demo-landing-001', 'demo-exit-ip-001', 'landing', '模拟用户直接路径：线路机 → 落地机 → 固定出口 IP', ?, ?, ?)`, []any{previousText, nowText, nowText}},
		{`INSERT OR IGNORE INTO node_events (id, node_id, event_type, severity, message, created_at)
VALUES ('demo-node-event-001', 'demo-relay-001', 'demo_seed', 'info', '模拟数据已写入，用于本地功能测试', ?)`, []any{nowText}},
	}
	for _, statement := range statements {
		if _, err := tx.Exec(statement.query, statement.args...); err != nil {
			return fmt.Errorf("execute demo seed statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit demo seed: %w", err)
	}
	return nil
}
