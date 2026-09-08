#!/usr/bin/env bash
set -euo pipefail

PREFIX="${PREFIX:-/usr/local/bin}"
CONFIG_DIR="${CONFIG_DIR:-/etc/xpanel-agent}"
DATA_DIR="${DATA_DIR:-/var/lib/xpanel-agent}"
SERVICE_NAME="xpanel-agent"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "run install.sh as root" >&2
  exit 1
fi

install -d -m 0755 "$PREFIX" "$CONFIG_DIR" "$DATA_DIR"
if ! id -u xpanel-agent >/dev/null 2>&1; then
  useradd --system --home-dir "$DATA_DIR" --no-create-home --shell /usr/sbin/nologin xpanel-agent
fi
chown xpanel-agent:xpanel-agent "$DATA_DIR"

if [[ -f "$CONFIG_DIR/agent.yaml" ]]; then
  chown xpanel-agent:xpanel-agent "$CONFIG_DIR/agent.yaml"
  chmod 0600 "$CONFIG_DIR/agent.yaml"
fi

if [[ ! -f "$PREFIX/xpanel-agent" ]]; then
  echo "missing binary: $PREFIX/xpanel-agent" >&2
  exit 1
fi
install -m 0644 "$(dirname "$0")/xpanel-agent.service" "/etc/systemd/system/$SERVICE_NAME.service"
systemctl daemon-reload
systemctl enable --now "$SERVICE_NAME.service"
echo "installed and started $SERVICE_NAME"
