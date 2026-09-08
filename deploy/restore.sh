#!/usr/bin/env bash
set -Eeuo pipefail

usage() {
  echo "Usage: $0 --source /path/to/backup.sqlite3 [--database /path/to/panel.db] --yes" >&2
  exit 2
}

source_path=""
database="${XPANEL_DATABASE:-/var/lib/xpanel-central/panel.db}"
confirmed=0
while [[ $# -gt 0 ]]; do
  case "$1" in
    --source) [[ $# -ge 2 ]] || usage; source_path="$2"; shift 2 ;;
    --database) [[ $# -ge 2 ]] || usage; database="$2"; shift 2 ;;
    --yes) confirmed=1; shift ;;
    *) usage ;;
  esac
done
[[ -n "$source_path" && "$confirmed" == 1 ]] || usage

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

"$binary" restore --source "$source_path" --database "$database" --yes
