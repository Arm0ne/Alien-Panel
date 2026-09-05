#!/usr/bin/env bash
set -Eeuo pipefail

# Run migrations with a pre-migration snapshot. When systemd is available, an
# active service is stopped and restarted automatically.
database="${XPANEL_DATABASE:-/var/lib/xpanel-central/panel.db}"
backup_dir="${XPANEL_BACKUP_DIR:-/var/backups/xpanel-central}"
binary="${XPANEL_DB_MAINTENANCE_BIN:-/usr/local/bin/xpanel-db-maintenance}"
service="${XPANEL_SERVICE_NAME:-xpanel-central.service}"
skip_stop="${XPANEL_SKIP_SERVICE_STOP:-0}"
restart=0

if [[ ! -x "$binary" ]]; then
  echo "db-maintenance binary is not executable: $binary" >&2
  exit 1
fi

cleanup() {
  if [[ "$restart" == 1 ]]; then
    systemctl start "$service"
  fi
}
trap cleanup EXIT

if [[ "$skip_stop" != 1 ]] && command -v systemctl >/dev/null 2>&1; then
  if systemctl is-active --quiet "$service"; then
    systemctl stop "$service"
    restart=1
  fi
fi

"$binary" migrate \
  --database "$database" \
  --backup-dir "$backup_dir"
