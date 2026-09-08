// Command db-maintenance provides the operational backup, restore, integrity
// check and migration workflow for the central SQLite database.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"xpanel-central/backend/internal/db"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "db-maintenance:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: db-maintenance <backup|verify|restore|migrate> [flags]")
	}
	ctx := context.Background()
	switch args[0] {
	case "backup":
		flags := flag.NewFlagSet("backup", flag.ContinueOnError)
		database := flags.String("database", envOr("XPANEL_DATABASE", "./data/panel.db"), "SQLite database path")
		backupDir := flags.String("backup-dir", envOr("XPANEL_BACKUP_DIR", "./data/backups"), "backup directory")
		retention := flags.Int("retention", 14, "number of generated backups to retain")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		result, err := db.Backup(ctx, db.BackupOptions{DatabasePath: *database, DestinationDir: *backupDir, Retention: *retention})
		if err != nil {
			return err
		}
		fmt.Printf("backup created: %s (%d bytes)\n", result.Path, result.Size)
		return nil

	case "verify":
		flags := flag.NewFlagSet("verify", flag.ContinueOnError)
		path := flags.String("path", "", "SQLite database or backup path")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *path == "" {
			return errors.New("verify requires --path")
		}
		if err := db.Verify(ctx, *path); err != nil {
			return err
		}
		fmt.Printf("database is valid: %s\n", *path)
		return nil

	case "restore":
		flags := flag.NewFlagSet("restore", flag.ContinueOnError)
		source := flags.String("source", "", "verified backup path")
		database := flags.String("database", envOr("XPANEL_DATABASE", "./data/panel.db"), "SQLite database path")
		confirmed := flags.Bool("yes", false, "confirm replacement of the database")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *source == "" {
			return errors.New("restore requires --source")
		}
		if !*confirmed {
			return errors.New("restore replaces the database; pass --yes after stopping the central service")
		}
		result, err := db.Restore(ctx, *source, *database)
		if err != nil {
			return err
		}
		if result.PreviousBackup != "" {
			fmt.Printf("restore complete: %s; previous database snapshot: %s\n", result.Path, result.PreviousBackup)
		} else {
			fmt.Printf("restore complete: %s\n", result.Path)
		}
		return nil

	case "migrate":
		flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
		database := flags.String("database", envOr("XPANEL_DATABASE", "./data/panel.db"), "SQLite database path")
		backupDir := flags.String("backup-dir", envOr("XPANEL_BACKUP_DIR", "./data/backups"), "pre-migration backup directory")
		noBackup := flags.Bool("no-backup", false, "skip the pre-migration backup (not recommended)")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if !*noBackup {
			if _, err := os.Stat(*database); err == nil {
				result, err := db.Backup(ctx, db.BackupOptions{DatabasePath: *database, DestinationDir: *backupDir, Retention: 14})
				if err != nil {
					return fmt.Errorf("create pre-migration backup: %w", err)
				}
				fmt.Printf("pre-migration backup: %s\n", result.Path)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("stat database: %w", err)
			}
		}
		databaseHandle, err := db.Open(*database)
		if err != nil {
			return err
		}
		defer databaseHandle.Close()
		if err := db.Migrate(databaseHandle); err != nil {
			return err
		}
		fmt.Printf("migrations are up to date: %s\n", filepath.Clean(*database))
		return nil

	default:
		return fmt.Errorf("unknown command %q; use backup, verify, restore or migrate", args[0])
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
