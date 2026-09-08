package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	DefaultBackupRetention = 14
	backupPrefix           = "panel-"
	backupSuffix           = ".sqlite3"
)

// BackupOptions controls a generated backup. A backup is produced with
// SQLite's VACUUM INTO operation, so a running central service can be backed
// up without copying an in-progress WAL file.
type BackupOptions struct {
	DatabasePath   string
	DestinationDir string
	Retention      int
	Now            func() time.Time
}

// BackupResult describes a successfully created backup.
type BackupResult struct {
	Path string
	Size int64
}

// RestoreResult describes a restore and the safety snapshot made before it.
type RestoreResult struct {
	Path           string
	PreviousBackup string
}

// Backup creates a timestamped, self-contained SQLite backup and applies the
// requested retention policy. It does not run migrations before taking the
// snapshot, which makes it suitable as a pre-upgrade backup.
func Backup(ctx context.Context, options BackupOptions) (BackupResult, error) {
	if strings.TrimSpace(options.DatabasePath) == "" || options.DatabasePath == ":memory:" {
		return BackupResult{}, errors.New("database path must reference a file")
	}
	if strings.TrimSpace(options.DestinationDir) == "" {
		return BackupResult{}, errors.New("backup destination directory is required")
	}
	if options.Retention == 0 {
		options.Retention = DefaultBackupRetention
	}
	if options.Retention < 1 {
		return BackupResult{}, errors.New("backup retention must be at least 1")
	}
	if options.Now == nil {
		options.Now = time.Now
	}

	if _, err := os.Stat(options.DatabasePath); err != nil {
		return BackupResult{}, fmt.Errorf("stat database: %w", err)
	}
	if err := os.MkdirAll(options.DestinationDir, 0o750); err != nil {
		return BackupResult{}, fmt.Errorf("create backup directory: %w", err)
	}

	stamp := options.Now().UTC().Format("20060102T150405Z")
	destination, err := nextBackupPath(options.DestinationDir, backupPrefix+stamp, backupSuffix)
	if err != nil {
		return BackupResult{}, err
	}
	result, err := backupTo(ctx, options.DatabasePath, destination)
	if err != nil {
		return BackupResult{}, err
	}
	if err := pruneBackups(options.DestinationDir, options.Retention); err != nil {
		return BackupResult{}, err
	}
	return result, nil
}

// BackupTo writes a backup to an exact path. The destination must not already
// contain a non-empty file. It is useful for named pre-restore snapshots.
func BackupTo(ctx context.Context, databasePath, destination string) (BackupResult, error) {
	if strings.TrimSpace(databasePath) == "" || databasePath == ":memory:" {
		return BackupResult{}, errors.New("database path must reference a file")
	}
	if strings.TrimSpace(destination) == "" {
		return BackupResult{}, errors.New("backup destination is required")
	}
	return backupTo(ctx, databasePath, destination)
}

func backupTo(ctx context.Context, databasePath, destination string) (BackupResult, error) {
	sourceAbs, err := filepath.Abs(databasePath)
	if err != nil {
		return BackupResult{}, fmt.Errorf("resolve database path: %w", err)
	}
	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		return BackupResult{}, fmt.Errorf("resolve backup path: %w", err)
	}
	if samePath(sourceAbs, destinationAbs) {
		return BackupResult{}, errors.New("backup destination must differ from database")
	}
	if err := os.MkdirAll(filepath.Dir(destinationAbs), 0o750); err != nil {
		return BackupResult{}, fmt.Errorf("create backup parent: %w", err)
	}
	if info, err := os.Stat(destinationAbs); err == nil {
		if info.Size() > 0 {
			return BackupResult{}, fmt.Errorf("backup destination already exists: %s", destinationAbs)
		}
		if err := os.Remove(destinationAbs); err != nil {
			return BackupResult{}, fmt.Errorf("remove empty backup destination: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return BackupResult{}, fmt.Errorf("stat backup destination: %w", err)
	}

	database, err := Open(sourceAbs)
	if err != nil {
		return BackupResult{}, err
	}
	defer database.Close()
	if err := database.PingContext(ctx); err != nil {
		return BackupResult{}, fmt.Errorf("ping database: %w", err)
	}
	// Checkpointing first keeps the normal WAL small. VACUUM INTO still gives
	// the snapshot its own consistent file if another writer is active.
	if _, err := database.ExecContext(ctx, `PRAGMA wal_checkpoint(PASSIVE)`); err != nil {
		return BackupResult{}, fmt.Errorf("checkpoint database: %w", err)
	}
	statement := "VACUUM INTO " + quoteSQLiteString(destinationAbs)
	if _, err := database.ExecContext(ctx, statement); err != nil {
		return BackupResult{}, fmt.Errorf("vacuum database into backup: %w", err)
	}
	if err := Verify(ctx, destinationAbs); err != nil {
		_ = os.Remove(destinationAbs)
		return BackupResult{}, fmt.Errorf("verify backup: %w", err)
	}
	if err := os.Chmod(destinationAbs, 0o600); err != nil {
		return BackupResult{}, fmt.Errorf("restrict backup permissions: %w", err)
	}
	info, err := os.Stat(destinationAbs)
	if err != nil {
		return BackupResult{}, fmt.Errorf("stat completed backup: %w", err)
	}
	return BackupResult{Path: destinationAbs, Size: info.Size()}, nil
}

// Verify checks SQLite's structural and foreign-key integrity without
// changing the database schema or running application migrations.
func Verify(ctx context.Context, path string) error {
	if strings.TrimSpace(path) == "" || path == ":memory:" {
		return errors.New("database path must reference a file")
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat database: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return errors.New("database file is empty or not regular")
	}
	database, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return fmt.Errorf("open database for verification: %w", err)
	}
	defer database.Close()
	if err := database.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	var integrity string
	if err := database.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&integrity); err != nil {
		return fmt.Errorf("run integrity check: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(integrity), "ok") {
		return fmt.Errorf("sqlite integrity check failed: %s", integrity)
	}
	rows, err := database.QueryContext(ctx, `PRAGMA foreign_key_check`)
	if err != nil {
		return fmt.Errorf("run foreign key check: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		var table, rowID, parent, foreignKey any
		if err := rows.Scan(&table, &rowID, &parent, &foreignKey); err != nil {
			return fmt.Errorf("read foreign key check: %w", err)
		}
		return fmt.Errorf("foreign key violation table=%v row=%v parent=%v fk=%v", table, rowID, parent, foreignKey)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("finish foreign key check: %w", err)
	}
	return nil
}

// Restore replaces databasePath with a verified backup. Before replacing an
// existing database it creates a timestamped sibling snapshot, so a failed
// rollout can be reversed. The service must be stopped by the caller before
// invoking this function.
func Restore(ctx context.Context, source, databasePath string) (RestoreResult, error) {
	if strings.TrimSpace(source) == "" || strings.TrimSpace(databasePath) == "" || source == ":memory:" || databasePath == ":memory:" {
		return RestoreResult{}, errors.New("source and destination must reference files")
	}
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		return RestoreResult{}, fmt.Errorf("resolve restore source: %w", err)
	}
	databaseAbs, err := filepath.Abs(databasePath)
	if err != nil {
		return RestoreResult{}, fmt.Errorf("resolve restore destination: %w", err)
	}
	if samePath(sourceAbs, databaseAbs) {
		return RestoreResult{}, errors.New("restore source must differ from destination")
	}
	if err := Verify(ctx, sourceAbs); err != nil {
		return RestoreResult{}, fmt.Errorf("verify restore source: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(databaseAbs), 0o750); err != nil {
		return RestoreResult{}, fmt.Errorf("create database parent: %w", err)
	}

	var previous string
	if _, err := os.Stat(databaseAbs); err == nil {
		previous, err = nextBackupPath(filepath.Dir(databaseAbs), strings.TrimSuffix(filepath.Base(databaseAbs), filepath.Ext(databaseAbs))+"-pre-restore", backupSuffix)
		if err != nil {
			return RestoreResult{}, err
		}
		if _, err := backupTo(ctx, databaseAbs, previous); err != nil {
			return RestoreResult{}, fmt.Errorf("snapshot current database before restore: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return RestoreResult{}, fmt.Errorf("stat restore destination: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(databaseAbs), "."+filepath.Base(databaseAbs)+".restore-*")
	if err != nil {
		return RestoreResult{}, fmt.Errorf("create restore temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		return RestoreResult{}, fmt.Errorf("close restore temporary file: %w", err)
	}
	_ = os.Remove(temporaryPath)
	if err := copyFile(sourceAbs, temporaryPath); err != nil {
		_ = os.Remove(temporaryPath)
		return RestoreResult{}, fmt.Errorf("copy restore source: %w", err)
	}
	if err := os.Chmod(temporaryPath, 0o600); err != nil {
		_ = os.Remove(temporaryPath)
		return RestoreResult{}, fmt.Errorf("restrict restored database permissions: %w", err)
	}

	oldPath := ""
	if _, err := os.Stat(databaseAbs); err == nil {
		oldPath, err = nextTemporaryPath(filepath.Dir(databaseAbs), "."+filepath.Base(databaseAbs)+".previous-*")
		if err != nil {
			_ = os.Remove(temporaryPath)
			return RestoreResult{}, err
		}
		if err := os.Rename(databaseAbs, oldPath); err != nil {
			_ = os.Remove(temporaryPath)
			return RestoreResult{}, fmt.Errorf("move current database aside: %w", err)
		}
	}
	if err := os.Rename(temporaryPath, databaseAbs); err != nil {
		if oldPath != "" {
			_ = os.Rename(oldPath, databaseAbs)
		}
		_ = os.Remove(temporaryPath)
		return RestoreResult{}, fmt.Errorf("install restored database: %w", err)
	}
	// WAL and shared-memory files belong to the replaced database. They must
	// not be allowed to be replayed into the newly restored snapshot.
	_ = os.Remove(databaseAbs + "-wal")
	_ = os.Remove(databaseAbs + "-shm")
	if err := Verify(ctx, databaseAbs); err != nil {
		_ = os.Remove(databaseAbs)
		if previous != "" {
			// Restore from the verified standalone snapshot instead of moving
			// back a database that may have had an uncheckpointed WAL.
			_ = copyFile(previous, databaseAbs)
		} else if oldPath != "" {
			_ = os.Rename(oldPath, databaseAbs)
		}
		return RestoreResult{}, fmt.Errorf("verify restored database: %w", err)
	}
	if oldPath != "" {
		_ = os.Remove(oldPath)
	}
	return RestoreResult{Path: databaseAbs, PreviousBackup: previous}, nil
}

func nextBackupPath(directory, stem, suffix string) (string, error) {
	for index := 0; index < 10000; index++ {
		name := stem + suffix
		if index > 0 {
			name = fmt.Sprintf("%s-%d%s", stem, index, suffix)
		}
		path := filepath.Join(directory, name)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return path, nil
		} else if err != nil {
			return "", fmt.Errorf("check backup path: %w", err)
		}
	}
	return "", errors.New("could not allocate a unique backup path")
}

func nextTemporaryPath(directory, pattern string) (string, error) {
	file, err := os.CreateTemp(directory, pattern)
	if err != nil {
		return "", fmt.Errorf("allocate temporary database path: %w", err)
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("close temporary database path: %w", err)
	}
	if err := os.Remove(path); err != nil {
		return "", fmt.Errorf("remove temporary database placeholder: %w", err)
	}
	return path, nil
}

func pruneBackups(directory string, retention int) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("list backups: %w", err)
	}
	type item struct {
		name string
		path string
		mod  time.Time
	}
	items := make([]item, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), backupPrefix) || !strings.HasSuffix(entry.Name(), backupSuffix) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat backup %s: %w", entry.Name(), err)
		}
		items = append(items, item{name: entry.Name(), path: filepath.Join(directory, entry.Name()), mod: info.ModTime()})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].mod.Equal(items[j].mod) {
			return items[i].name > items[j].name
		}
		return items[i].mod.After(items[j].mod)
	})
	if len(items) <= retention {
		return nil
	}
	for _, old := range items[retention:] {
		if err := os.Remove(old.path); err != nil {
			return fmt.Errorf("remove expired backup %s: %w", old.name, err)
		}
	}
	return nil
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		return err
	}
	if err := output.Sync(); err != nil {
		_ = output.Close()
		return err
	}
	return output.Close()
}

func quoteSQLiteString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func samePath(first, second string) bool {
	if filepath.Clean(first) == filepath.Clean(second) {
		return true
	}
	return runtime.GOOS == "windows" && strings.EqualFold(filepath.Clean(first), filepath.Clean(second))
}
