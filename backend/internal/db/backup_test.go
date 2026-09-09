package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMigrateIsIdempotent(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if err := Migrate(database); err != nil {
		t.Fatalf("first migration: %v", err)
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("second migration: %v", err)
	}
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("count migrations: %v", err)
	}
	if count != 21 {
		t.Fatalf("migration count = %d, want 21", count)
	}
}

func TestLegacyAgentVersionMigrationClearsInstallerLabel(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		t.Fatalf("create migration table: %v", err)
	}
	for _, version := range legacyVersionMigrationNames {
		sqlBytes, readErr := migrationFiles.ReadFile("migrations/" + version)
		if readErr != nil {
			t.Fatalf("read %s: %v", version, readErr)
		}
		if _, execErr := database.Exec(string(sqlBytes)); execErr != nil {
			t.Fatalf("apply %s: %v", version, execErr)
		}
		if _, execErr := database.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, 'now')`, version); execErr != nil {
			t.Fatalf("record %s: %v", version, execErr)
		}
	}
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, agent_version, created_at, updated_at) VALUES ('legacy-node', 'legacy-node', 'legacy node', 'online', '2026-09-06T00:00:00Z', '2026-09-06T00:00:00Z')`); err != nil {
		t.Fatalf("insert legacy node: %v", err)
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("apply legacy-version migration: %v", err)
	}
	var version sql.NullString
	if err := database.QueryRow(`SELECT agent_version FROM nodes WHERE id = 'legacy-node'`).Scan(&version); err != nil {
		t.Fatalf("read migrated Agent version: %v", err)
	}
	if version.Valid {
		t.Fatalf("legacy Agent version = %q, want NULL", version.String)
	}
}

func TestAutoDeleteOrphanedUsersMigrationPreservesBillingHistory(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		t.Fatalf("create migration table: %v", err)
	}
	versions := append(append([]string{}, legacyVersionMigrationNames...), "015_clear_legacy_agent_version.sql", "016_client_expiry_and_sync_health.sql")
	for _, version := range versions {
		sqlBytes, readErr := migrationFiles.ReadFile("migrations/" + version)
		if readErr != nil {
			t.Fatalf("read %s: %v", version, readErr)
		}
		if _, execErr := database.Exec(string(sqlBytes)); execErr != nil {
			t.Fatalf("apply %s: %v", version, execErr)
		}
		if _, execErr := database.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, 'now')`, version); execErr != nil {
			t.Fatalf("record %s: %v", version, execErr)
		}
	}

	now := "2026-09-06T00:00:00Z"
	for _, statement := range []string{
		`INSERT INTO users (id, display_name, status, monthly_fee, billing_cycle, billing_amount, currency, created_at, updated_at) VALUES ('orphan-user', '旧孤立用户', 'active', 120, 'monthly', 120, 'CNY', '` + now + `', '` + now + `')`,
		`INSERT INTO users (id, display_name, status, monthly_fee, billing_cycle, billing_amount, currency, created_at, updated_at) VALUES ('current-user', '当前用户', 'active', 120, 'monthly', 120, 'CNY', '` + now + `', '` + now + `')`,
		`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('current-relay', 'current-relay', 'Current Relay', 'relay', '` + now + `', '` + now + `')`,
		`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, first_seen_at) VALUES ('current-inbound', 'current-relay', '1', 'current-user', 'user', '` + now + `')`,
		`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('current-mapping', 'current-user', 'current-inbound', 1, '` + now + `')`,
		`INSERT INTO user_billing_records (id, user_id, billing_cycle, amount, currency, service_from, service_to, paid_at, status, source, created_at) VALUES ('orphan-record', 'orphan-user', 'monthly', 120, 'CNY', '2026-09-01T00:00:00Z', '2026-10-01T00:00:00Z', '` + now + `', 'confirmed', 'manual', '` + now + `')`,
		`INSERT INTO user_renewal_candidates (id, user_id, old_expiry_at, new_expiry_at, detected_at, suggested_cycle, suggested_amount, currency, status) VALUES ('orphan-candidate', 'orphan-user', '2026-09-01T00:00:00Z', '2026-10-01T00:00:00Z', '` + now + `', 'monthly', 120, 'CNY', 'pending')`,
		`INSERT INTO node_events (id, event_type, severity, message, created_at, event_category, title, requires_action, event_status, resource_type, resource_id, source) VALUES ('orphan-event', 'renewal_candidate_detected', 'warning', 'pending renewal', '` + now + `', 'business', 'pending renewal', 1, 'open', 'renewal', 'orphan-candidate', 'agent')`,
	} {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("insert legacy auto-delete fixture: %v", err)
		}
	}

	if err := Migrate(database); err != nil {
		t.Fatalf("apply auto-delete migration: %v", err)
	}
	var orphanDeletedAt sql.NullString
	var orphanStatus string
	if err := database.QueryRow(`SELECT deleted_at, status FROM users WHERE id = 'orphan-user'`).Scan(&orphanDeletedAt, &orphanStatus); err != nil {
		t.Fatalf("read migrated orphan user: %v", err)
	}
	if !orphanDeletedAt.Valid || orphanStatus != "disabled" {
		t.Fatalf("migrated orphan state deleted_at=%#v status=%q", orphanDeletedAt, orphanStatus)
	}
	var currentDeletedAt sql.NullString
	if err := database.QueryRow(`SELECT deleted_at FROM users WHERE id = 'current-user'`).Scan(&currentDeletedAt); err != nil {
		t.Fatalf("read migrated current user: %v", err)
	}
	if currentDeletedAt.Valid {
		t.Fatalf("current user was incorrectly deleted: %#v", currentDeletedAt)
	}
	var billingRecords, auditRecords int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_billing_records WHERE user_id = 'orphan-user'`).Scan(&billingRecords); err != nil {
		t.Fatalf("count retained billing records: %v", err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM audit_logs WHERE action = 'user.auto_delete' AND resource_id = 'orphan-user'`).Scan(&auditRecords); err != nil {
		t.Fatalf("count auto-delete audit records: %v", err)
	}
	if billingRecords != 1 || auditRecords != 1 {
		t.Fatalf("retained billing records=%d audit records=%d", billingRecords, auditRecords)
	}
	var candidateStatus, eventStatus string
	var requiresAction int
	if err := database.QueryRow(`SELECT status FROM user_renewal_candidates WHERE id = 'orphan-candidate'`).Scan(&candidateStatus); err != nil {
		t.Fatalf("read migrated renewal candidate: %v", err)
	}
	if err := database.QueryRow(`SELECT event_status, requires_action FROM node_events WHERE id = 'orphan-event'`).Scan(&eventStatus, &requiresAction); err != nil {
		t.Fatalf("read migrated renewal event: %v", err)
	}
	if candidateStatus != "rejected" || eventStatus != "resolved" || requiresAction != 0 {
		t.Fatalf("migrated renewal candidate=%q event=%q requires_action=%d", candidateStatus, eventStatus, requiresAction)
	}
}

var legacyVersionMigrationNames = []string{
	"001_initial.sql", "002_inbound_missing_sync_count.sql", "003_node_metrics.sql", "004_exit_ip_sources.sql",
	"005_node_deletion.sql", "006_user_route_assignment.sql", "007_clear_inactive_route_exit_refs.sql", "008_user_paths.sql",
	"009_cleanup_non_relay_user_mappings.sql", "010_node_install_tokens.sql", "011_node_management_url.sql",
	"012_traffic_snapshot_inbound_time.sql", "013_user_billing.sql", "014_event_center.sql",
}

func TestExitIPSourceMigrationPreservesLegacyBindings(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		t.Fatalf("create migration table: %v", err)
	}
	for _, version := range []string{"001_initial.sql", "002_inbound_missing_sync_count.sql", "003_node_metrics.sql"} {
		sqlBytes, readErr := migrationFiles.ReadFile("migrations/" + version)
		if readErr != nil {
			t.Fatalf("read %s: %v", version, readErr)
		}
		if _, execErr := database.Exec(string(sqlBytes)); execErr != nil {
			t.Fatalf("apply %s: %v", version, execErr)
		}
		if _, execErr := database.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, 'now')`, version); execErr != nil {
			t.Fatalf("record %s: %v", version, execErr)
		}
	}
	now := "2026-09-03T00:00:00Z"
	for _, statement := range []string{
		`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('legacy-relay', 'legacy-relay', 'Legacy Relay', 'relay', '` + now + `', '` + now + `')`,
		`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('legacy-landing', 'legacy-landing', 'Legacy Landing', 'landing', '` + now + `', '` + now + `')`,
		`INSERT INTO routes (id, name, relay_node_id, landing_node_id, created_at, updated_at) VALUES ('legacy-route', 'Legacy Route', 'legacy-relay', 'legacy-landing', '` + now + `', '` + now + `')`,
		`INSERT INTO exit_ips (id, landing_node_id, ip, created_at, updated_at) VALUES ('legacy-exit', 'legacy-landing', '198.51.100.90', '` + now + `', '` + now + `')`,
		`INSERT INTO route_exit_ips (id, route_id, exit_ip_id) VALUES ('legacy-binding', 'legacy-route', 'legacy-exit')`,
	} {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("insert legacy fixture: %v", err)
		}
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("apply source migration: %v", err)
	}
	var sourceType, ownerNodeID, scope string
	if err := database.QueryRow(`SELECT source_type, owner_node_id FROM exit_ips WHERE id = 'legacy-exit'`).Scan(&sourceType, &ownerNodeID); err != nil {
		t.Fatalf("read migrated exit IP: %v", err)
	}
	if sourceType != "node" || ownerNodeID != "legacy-landing" {
		t.Fatalf("migrated source = %q/%q, want node/legacy-landing", sourceType, ownerNodeID)
	}
	if err := database.QueryRow(`SELECT scope FROM route_exit_ips WHERE id = 'legacy-binding'`).Scan(&scope); err != nil {
		t.Fatalf("read migrated binding: %v", err)
	}
	if scope != "landing" {
		t.Fatalf("migrated scope = %q, want landing", scope)
	}
}

func TestNonRelayUserMappingMigrationKeepsUsers(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()
	if _, err := database.Exec(`CREATE TABLE schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		t.Fatalf("create migration table: %v", err)
	}
	versions := []string{
		"001_initial.sql", "002_inbound_missing_sync_count.sql", "003_node_metrics.sql",
		"004_exit_ip_sources.sql", "005_node_deletion.sql", "006_user_route_assignment.sql",
		"007_clear_inactive_route_exit_refs.sql", "008_user_paths.sql",
	}
	for _, version := range versions {
		sqlBytes, readErr := migrationFiles.ReadFile("migrations/" + version)
		if readErr != nil {
			t.Fatalf("read %s: %v", version, readErr)
		}
		if _, execErr := database.Exec(string(sqlBytes)); execErr != nil {
			t.Fatalf("apply %s: %v", version, execErr)
		}
		if _, execErr := database.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, 'now')`, version); execErr != nil {
			t.Fatalf("record %s: %v", version, execErr)
		}
	}
	now := "2026-09-03T00:00:00Z"
	for _, statement := range []string{
		`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('legacy-landing', 'legacy-landing', 'Legacy Landing', 'landing', '` + now + `', '` + now + `')`,
		`INSERT INTO users (id, display_name, status, created_at, updated_at) VALUES ('kept-user', '预留用户', 'active', '` + now + `', '` + now + `')`,
		`INSERT INTO inbounds (id, node_id, remote_inbound_id, user_id, kind, tag, first_seen_at, last_seen_at) VALUES ('legacy-landing-inbound', 'legacy-landing', '101', 'kept-user', 'user', 'legacy-landing-entry', '` + now + `', '` + now + `')`,
		`INSERT INTO user_inbounds (id, user_id, inbound_id, is_primary, active_from) VALUES ('legacy-landing-mapping', 'kept-user', 'legacy-landing-inbound', 1, '` + now + `')`,
	} {
		if _, err := database.Exec(statement); err != nil {
			t.Fatalf("insert legacy landing mapping: %v", err)
		}
	}
	if err := Migrate(database); err != nil {
		t.Fatalf("apply cleanup migration: %v", err)
	}
	var mappings, users int
	if err := database.QueryRow(`SELECT COUNT(*) FROM user_inbounds WHERE inbound_id = 'legacy-landing-inbound'`).Scan(&mappings); err != nil {
		t.Fatalf("count cleaned mappings: %v", err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM users WHERE id = 'kept-user'`).Scan(&users); err != nil {
		t.Fatalf("count kept users: %v", err)
	}
	var userID, kind string
	if err := database.QueryRow(`SELECT COALESCE(user_id, ''), kind FROM inbounds WHERE id = 'legacy-landing-inbound'`).Scan(&userID, &kind); err != nil {
		t.Fatalf("read cleaned landing inbound: %v", err)
	}
	if mappings != 0 || users != 1 || userID != "" || kind != "infrastructure" {
		t.Fatalf("cleanup mappings=%d users=%d user_id=%q kind=%q", mappings, users, userID, kind)
	}
}

func TestBackupVerifyAndRestore(t *testing.T) {
	root := t.TempDir()
	databasePath := filepath.Join(root, "panel.db")
	backupDir := filepath.Join(root, "backups")
	database, err := Open(databasePath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := Migrate(database); err != nil {
		database.Close()
		t.Fatalf("migrate database: %v", err)
	}
	if _, err := database.Exec(`INSERT INTO nodes (id, node_key, name, type, created_at, updated_at) VALUES ('node-1', 'node-key', 'Node', 'relay', '2026-09-03T00:00:00Z', '2026-09-03T00:00:00Z')`); err != nil {
		database.Close()
		t.Fatalf("insert fixture: %v", err)
	}
	database.Close()

	ctx := context.Background()
	result, err := Backup(ctx, BackupOptions{
		DatabasePath:   databasePath,
		DestinationDir: backupDir,
		Retention:      2,
		Now:            func() time.Time { return time.Date(2026, 9, 3, 1, 2, 3, 0, time.UTC) },
	})
	if err != nil {
		t.Fatalf("backup: %v", err)
	}
	if err := Verify(ctx, result.Path); err != nil {
		t.Fatalf("verify backup: %v", err)
	}
	if result.Size == 0 {
		t.Fatal("backup is empty")
	}

	mutated, err := Open(databasePath)
	if err != nil {
		t.Fatalf("reopen database: %v", err)
	}
	if _, err := mutated.Exec(`UPDATE nodes SET name = 'Mutated' WHERE id = 'node-1'`); err != nil {
		mutated.Close()
		t.Fatalf("mutate fixture: %v", err)
	}
	mutated.Close()

	restored, err := Restore(ctx, result.Path, databasePath)
	if err != nil {
		t.Fatalf("restore: %v", err)
	}
	if restored.PreviousBackup == "" {
		t.Fatal("restore did not create previous database snapshot")
	}
	if err := Verify(ctx, databasePath); err != nil {
		t.Fatalf("verify restored database: %v", err)
	}
	check, err := Open(databasePath)
	if err != nil {
		t.Fatalf("open restored database: %v", err)
	}
	defer check.Close()
	var name string
	if err := check.QueryRow(`SELECT name FROM nodes WHERE id = 'node-1'`).Scan(&name); err != nil {
		t.Fatalf("read restored fixture: %v", err)
	}
	if name != "Node" {
		t.Fatalf("restored name = %q, want Node", name)
	}
	if _, err := os.Stat(restored.PreviousBackup); err != nil {
		t.Fatalf("previous snapshot missing: %v", err)
	}
}

func TestBackupRetention(t *testing.T) {
	root := t.TempDir()
	databasePath := filepath.Join(root, "panel.db")
	database, err := Open(databasePath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := Migrate(database); err != nil {
		database.Close()
		t.Fatalf("migrate database: %v", err)
	}
	database.Close()

	backupDir := filepath.Join(root, "backups")
	for index := 0; index < 3; index++ {
		stamp := time.Date(2026, 9, 3, 0, 0, index, 0, time.UTC)
		if _, err := Backup(context.Background(), BackupOptions{DatabasePath: databasePath, DestinationDir: backupDir, Retention: 2, Now: func() time.Time { return stamp }}); err != nil {
			t.Fatalf("backup %d: %v", index, err)
		}
	}
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("read backup directory: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("backup count = %d, want 2", len(entries))
	}
}

func TestVerifyMissingFileDoesNotCreateDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.db")
	if err := Verify(context.Background(), path); err == nil {
		t.Fatal("verify missing file unexpectedly succeeded")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("verify created missing file, stat error = %v", err)
	}
}
