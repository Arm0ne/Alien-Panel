#!/usr/bin/env bash
set -Eeuo pipefail

# Daily backup entry point. Install the db-maintenance binary from the same
# release as the central service; it uses SQLite VACUUM INTO and handles WAL.
database="${XPANEL_DATABASE:-/var/lib/xpanel-central/panel.db}"
backup_dir="${XPANEL_BACKUP_DIR:-/var/backups/xpanel-central}"
retention="${XPANEL_BACKUP_RETENTION:-14}"
binary="${XPANEL_DB_MAINTENANCE_BIN:-/usr/local/bin/xpanel-db-maintenance}"

if [[ ! -x "$binary" ]]; then
  echo "db-maintenance binary is not executable: $binary" >&2
  exit 1
fi

exec "$binary" backup \
  --database "$database" \
  --backup-dir "$backup_dir" \
  --retention "$retention"
