#!/usr/bin/env bash
set -euo pipefail

# Remove the X-Panel Central test deployment installed by this repository.
# The script is deliberately dry-run by default. Pass --yes to remove the
# service, binaries, configuration, database, backups, Docker data and upload
# directories listed below.

CONFIRM=0
case "${1:-}" in
  --yes) CONFIRM=1 ;;
  --dry-run|"") ;;
  --help|-h)
    cat <<'EOF'
Usage:
  sudo bash uninstall-all.sh             # show what would be removed
  sudo bash uninstall-all.sh --yes       # stop and remove the test install

This permanently removes the X-Panel Central service, database, backups,
Docker volume/image, Nginx site file, uploaded frontend and release binaries.
It does not remove X-Panel Agent files on node machines or TLS certificates.
EOF
    exit 0
    ;;
  *)
    echo "unknown option: $1 (use --yes, --dry-run or --help)" >&2
    exit 2
    ;;
esac

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root: sudo bash uninstall-all.sh --yes" >&2
  exit 1
fi

if (( CONFIRM == 0 )); then
  echo "Dry run only. Nothing will be deleted. Use --yes to confirm permanent removal."
fi

remove_path() {
  local path="$1"
  if [[ -e "$path" || -L "$path" ]]; then
    echo "remove: $path"
    if (( CONFIRM == 1 )); then
      rm -rf -- "$path"
    fi
  fi
}

if command -v systemctl >/dev/null 2>&1; then
  if systemctl is-active --quiet xpanel-central.service 2>/dev/null; then
    echo "stop: xpanel-central.service"
    if (( CONFIRM == 1 )); then
      systemctl disable --now xpanel-central.service || true
    fi
  elif systemctl is-enabled --quiet xpanel-central.service 2>/dev/null; then
    echo "disable: xpanel-central.service"
    if (( CONFIRM == 1 )); then
      systemctl disable xpanel-central.service || true
    fi
  fi
fi

if command -v docker >/dev/null 2>&1; then
  compose_file="/opt/xpanel-central/deploy/docker-compose.yml"
  if [[ -f "$compose_file" ]]; then
    echo "stop/remove Docker project: xpanel-central"
    if (( CONFIRM == 1 )); then
      if docker compose version >/dev/null 2>&1; then
        docker compose -p xpanel-central -f "$compose_file" down -v --remove-orphans || true
      elif command -v docker-compose >/dev/null 2>&1; then
        docker-compose -p xpanel-central -f "$compose_file" down -v --remove-orphans || true
      fi
    fi
  fi
  if docker volume inspect xpanel-central_central-data >/dev/null 2>&1; then
    echo "remove Docker volume: xpanel-central_central-data"
    if (( CONFIRM == 1 )); then docker volume rm xpanel-central_central-data || true; fi
  fi
  if docker image inspect xpanel-central:local >/dev/null 2>&1; then
    echo "remove Docker image: xpanel-central:local"
    if (( CONFIRM == 1 )); then docker image rm xpanel-central:local || true; fi
  fi
fi

remove_path /etc/systemd/system/xpanel-central.service
remove_path /etc/systemd/system/multi-user.target.wants/xpanel-central.service
remove_path /etc/xpanel-central
remove_path /var/lib/xpanel-central
remove_path /var/backups/xpanel-central
remove_path /var/log/xpanel-central
remove_path /opt/xpanel-central
remove_path /usr/local/bin/xpanel-central
remove_path /usr/local/bin/xpanel-db-maintenance
remove_path /usr/local/bin/xpanel-traffic-check
remove_path /tmp/xpanel-frontend
remove_path /tmp/dist
remove_path /etc/nginx/sites-enabled/xpanel-central
remove_path /etc/nginx/sites-available/xpanel-central

if (( CONFIRM == 1 )); then
  if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
    systemctl reset-failed xpanel-central.service || true
  fi
  if command -v nginx >/dev/null 2>&1 && nginx -t >/dev/null 2>&1; then
    systemctl reload nginx 2>/dev/null || true
  fi
  echo "X-Panel Central test deployment was removed. TLS certificates and node Agents were kept."
else
  echo "Dry run complete. Re-run with --yes to perform the removal."
fi
