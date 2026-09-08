package main

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestReadOnlySQLiteDSN(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "panel.db")
	writable, err := sql.Open("sqlite", databasePath)
	if err != nil {
		t.Fatalf("open writable database: %v", err)
	}
	if _, err := writable.Exec(`CREATE TABLE check_rows (value TEXT NOT NULL); INSERT INTO check_rows (value) VALUES ('ok')`); err != nil {
		writable.Close()
		t.Fatalf("seed writable database: %v", err)
	}
	if err := writable.Close(); err != nil {
		t.Fatalf("close writable database: %v", err)
	}
	readonly, err := sql.Open("sqlite", readOnlySQLiteDSN(databasePath))
	if err != nil {
		t.Fatalf("open read-only database: %v", err)
	}
	defer readonly.Close()
	var value string
	if err := readonly.QueryRow(`SELECT value FROM check_rows`).Scan(&value); err != nil {
		t.Fatalf("query read-only database: %v", err)
	}
	if value != "ok" {
		t.Fatalf("value = %q", value)
	}
}
