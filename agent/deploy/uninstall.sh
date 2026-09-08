#!/usr/bin/env bash
set -euo pipefail

PREFIX="${PREFIX:-/usr/local/bin}"
CONFIG_DIR="${CONFIG_DIR:-/etc/xpanel-agent}"
DATA_DIR="${DATA_DIR:-/var/lib/xpanel-agent}"
SERVICE_NAME="xpanel-agent"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "run uninstall.sh as root" >&2
  exit 1
fi

systemctl disable --now "$SERVICE_NAME.service" 2>/dev/null || true
rm -f "/etc/systemd/system/$SERVICE_NAME.service"
systemctl daemon-reload
rm -f "$PREFIX/xpanel-agent"
echo "removed service and binary; configuration/data were kept at $CONFIG_DIR and $DATA_DIR"
