package httpapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"xpanel-central/backend/internal/db"
)

const maxRestoreUploadSize int64 = 512 << 20

func (s *Server) backupDirectory() string {
	return filepath.Join(filepath.Dir(s.cfg.DatabasePath), "backups")
}

func (s *Server) downloadBackup(w http.ResponseWriter, r *http.Request) {
	if s.cfg.DatabasePath == ":memory:" {
		writeFailure(w, http.StatusNotImplemented, internalErrorCode, "内存数据库不支持下载备份")
		return
	}
	if err := os.MkdirAll(s.backupDirectory(), 0o750); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法创建备份目录")
		return
	}
	temporary, err := os.CreateTemp(s.backupDirectory(), ".panel-download-*.sqlite3")
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法创建临时备份")
		return
	}
	temporaryPath := temporary.Name()
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法准备临时备份")
		return
	}
	_ = os.Remove(temporaryPath)
	defer os.Remove(temporaryPath)

	result, err := db.BackupTo(r.Context(), s.cfg.DatabasePath, temporaryPath)
	if err != nil {
		s.logger.Error("create downloadable backup", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "生成备份失败")
		return
	}
	filename := "alien-panel-" + time.Now().UTC().Format("20060102T150405Z") + ".sqlite3"
	w.Header().Set("Content-Type", "application/vnd.sqlite3")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Header().Set("Content-Length", strconv.FormatInt(result.Size, 10))
	w.WriteHeader(http.StatusOK)
	file, err := os.Open(result.Path)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = io.Copy(w, file)
}

func (s *Server) restoreBackup(w http.ResponseWriter, r *http.Request) {
	if s.cfg.DatabasePath == ":memory:" {
		writeFailure(w, http.StatusNotImplemented, internalErrorCode, "内存数据库不支持恢复")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxRestoreUploadSize+1)
	if err := r.ParseMultipartForm(maxRestoreUploadSize + 1); err != nil {
		writeFailure(w, http.StatusRequestEntityTooLarge, validationCode, "备份文件超过 512 MB 限制")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "请选择 SQLite 备份文件")
		return
	}
	defer file.Close()
	if header.Size > maxRestoreUploadSize {
		writeFailure(w, http.StatusRequestEntityTooLarge, validationCode, "备份文件超过 512 MB 限制")
		return
	}
	if err := os.MkdirAll(s.backupDirectory(), 0o750); err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法创建备份目录")
		return
	}
	temporary, err := os.CreateTemp(s.backupDirectory(), ".panel-restore-*.sqlite3")
	if err != nil {
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法保存上传文件")
		return
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "无法限制上传文件权限")
		return
	}
	copied, err := io.Copy(temporary, io.LimitReader(file, maxRestoreUploadSize+1))
	if err != nil {
		_ = temporary.Close()
		writeFailure(w, http.StatusBadRequest, validationCode, "读取备份文件失败")
		return
	}
	if copied > maxRestoreUploadSize {
		_ = temporary.Close()
		writeFailure(w, http.StatusRequestEntityTooLarge, validationCode, "备份文件超过 512 MB 限制")
		return
	}
	if err := temporary.Close(); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "保存备份文件失败")
		return
	}
	if err := db.Verify(r.Context(), temporaryPath); err != nil {
		writeFailure(w, http.StatusBadRequest, validationCode, "备份文件校验失败，请选择有效的 SQLite 备份")
		return
	}

	s.dbMu.Lock()
	result, err := s.restoreDatabaseLocked(r.Context(), temporaryPath)
	s.dbMu.Unlock()
	if err != nil {
		s.logger.Error("restore database", "error", err)
		writeFailure(w, http.StatusInternalServerError, internalErrorCode, "恢复失败，原数据已保留")
		return
	}
	s.writeAuditLog(r, "system.database_restore", "database", "panel", nil, map[string]any{"restored": true})
	previousBackup := ""
	if result.PreviousBackup != "" {
		previousBackup = filepath.Base(result.PreviousBackup)
	}
	writeSuccess(w, map[string]any{
		"restored":       true,
		"previousBackup": previousBackup,
		"requiresLogin":  true,
	})
}

func (s *Server) restoreDatabaseLocked(ctx context.Context, source string) (db.RestoreResult, error) {
	if s.db == nil {
		return db.RestoreResult{}, errors.New("database is not open")
	}
	if err := s.db.Close(); err != nil {
		return db.RestoreResult{}, fmt.Errorf("close current database: %w", err)
	}
	result, err := db.Restore(ctx, source, s.cfg.DatabasePath)
	if err != nil {
		s.reopenDatabase()
		return db.RestoreResult{}, err
	}
	reopened, err := db.Open(s.cfg.DatabasePath)
	if err == nil {
		err = db.Migrate(reopened)
	}
	if err != nil {
		if reopened != nil {
			_ = reopened.Close()
		}
		if result.PreviousBackup != "" {
			_, _ = db.Restore(context.Background(), result.PreviousBackup, s.cfg.DatabasePath)
		}
		s.reopenDatabase()
		return db.RestoreResult{}, fmt.Errorf("open restored database: %w", err)
	}
	// Sessions from a backup are stale by definition. Revoking them also
	// ensures the browser cannot continue using a token from before restore.
	if _, err := reopened.Exec(`UPDATE sessions SET revoked_at = COALESCE(revoked_at, strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))`); err != nil {
		_ = reopened.Close()
		s.reopenDatabase()
		return db.RestoreResult{}, fmt.Errorf("revoke restored sessions: %w", err)
	}
	s.db = reopened
	return result, nil
}

func (s *Server) reopenDatabase() {
	reopened, err := db.Open(s.cfg.DatabasePath)
	if err != nil {
		s.logger.Error("reopen database after restore failure", "error", err)
		return
	}
	s.db = reopened
}
