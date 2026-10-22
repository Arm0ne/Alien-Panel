package httpapi

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"xpanel-central/backend/internal/db"
)

const (
	backupFileNamePrefix = "panel-"
	backupFileNameSuffix = ".sqlite3"
)

// backupEntry is the operator-facing description of one server-side backup.
type backupEntry struct {
	Name      string
	Size      int64
	CreatedAt string
}

// listBackups returns every server-side backup file in the configured backup
// directory. The directory is derived from the database path, so a caller can
// never point it somewhere else.
func (s *Server) listBackups(w http.ResponseWriter, r *http.Request) {
	directory := s.backupDirectory()
	entries, err := readBackupEntries(directory)
	if err != nil {
		s.logger.Error("list backups", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法读取备份目录")
		return
	}
	items := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		items = append(items, map[string]any{
			"name":      entry.Name,
			"size":      entry.Size,
			"createdAt": entry.CreatedAt,
		})
	}
	latest := ""
	if len(entries) > 0 {
		latest = entries[0].Name
	}
	writeSuccess(w, map[string]any{"items": items, "total": len(items), "directory": directory, "latest": latest})
}

// createBackup writes a manual backup into the backup directory without
// pruning anything: retention is managed explicitly from the UI.
func (s *Server) createBackup(w http.ResponseWriter, r *http.Request) {
	if s.cfg.DatabasePath == ":memory:" {
		writeFailure(w, http.StatusNotImplemented, internalErrorCode, "内存数据库不支持备份")
		return
	}
	directory := s.backupDirectory()
	if err := os.MkdirAll(directory, 0o750); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法创建备份目录")
		return
	}
	result, err := db.Backup(r.Context(), db.BackupOptions{
		DatabasePath:   s.cfg.DatabasePath,
		DestinationDir: directory,
		Retention:      0,
	})
	if err != nil {
		s.logger.Error("create manual backup", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "生成备份失败")
		return
	}
	name := filepath.Base(result.Path)
	s.writeAuditLog(r, "system.backup_create", "database", "panel", nil, map[string]any{"file": name, "size": result.Size})
	writeSuccess(w, map[string]any{"name": name, "size": result.Size})
}

// verifyBackup runs SQLite's integrity check against one stored backup.
func (s *Server) verifyBackup(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "invalid verify payload")
		return
	}
	path, err := s.resolveBackupPath(payload.Name)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	if err := db.Verify(r.Context(), path); err != nil {
		writeSuccess(w, map[string]any{"ok": false, "message": err.Error()})
		return
	}
	writeSuccess(w, map[string]any{"ok": true})
}

// deleteBackup removes one stored backup. The newest backup is protected so an
// operator cannot accidentally leave the panel without a recent snapshot.
func (s *Server) deleteBackup(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.PathValue("name"))
	path, err := s.resolveBackupPath(name)
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, err.Error())
		return
	}
	entries, err := readBackupEntries(s.backupDirectory())
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法读取备份目录")
		return
	}
	if len(entries) > 0 && entries[0].Name == name {
		writeFailure(w, http.StatusConflict, validationCode, "最新一份备份不可删除")
		return
	}
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeFailure(w, http.StatusNotFound, notFoundCode, "备份不存在")
			return
		}
		s.logger.Error("delete backup", "file", name, "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "删除备份失败")
		return
	}
	s.writeAuditLog(r, "system.backup_delete", "database", "panel", map[string]any{"file": name}, nil)
	writeSuccess(w, map[string]any{"deleted": name})
}

// resolveBackupPath validates a client-supplied file name and resolves it
// inside the backup directory. Only plain panel-*.sqlite3 names are accepted,
// which rules out path traversal and arbitrary file reads.
func (s *Server) resolveBackupPath(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed != filepath.Base(trimmed) || strings.ContainsAny(trimmed, `/\`) {
		return "", errors.New("备份文件名不合法")
	}
	if !strings.HasPrefix(trimmed, backupFileNamePrefix) || !strings.HasSuffix(trimmed, backupFileNameSuffix) {
		return "", errors.New("备份文件名不合法")
	}
	return filepath.Join(s.backupDirectory(), trimmed), nil
}

// readBackupEntries lists backup files newest first.
func readBackupEntries(directory string) ([]backupEntry, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []backupEntry{}, nil
		}
		return nil, err
	}
	items := make([]backupEntry, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, backupFileNamePrefix) || !strings.HasSuffix(name, backupFileNameSuffix) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, backupEntry{
			Name:      name,
			Size:      info.Size(),
			CreatedAt: info.ModTime().UTC().Format(time.RFC3339Nano),
		})
	}
	sort.Slice(items, func(left, right int) bool {
		if items[left].CreatedAt == items[right].CreatedAt {
			return items[left].Name > items[right].Name
		}
		return items[left].CreatedAt > items[right].CreatedAt
	})
	return items, nil
}
