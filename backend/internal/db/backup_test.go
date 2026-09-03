package db

import (
	"context"
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
	if count != 7 {
		t.Fatalf("migration count = %d, want 7", count)
	}
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
