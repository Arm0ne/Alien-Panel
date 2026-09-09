#!/usr/bin/env bash
set -Eeuo pipefail

# Alien-Panel v1.0.0 production installer.
#
# This file is intentionally self-contained so it can be executed directly
# from a GitHub raw URL. It downloads the versioned production directory from
# the repository, validates the release files, installs Docker Compose files
# under /opt/alien-panel, and starts the two-container deployment.

PRODUCT_NAME="Alien-Panel"
PACKAGE_DIR_NAME="Alien-Panel"
DEFAULT_REPO="https://github.com/Arm0ne/Alien-Panel.git"
DEFAULT_REF="main"
DEFAULT_INSTALL_DIR="/opt/alien-panel"
DEFAULT_PUBLIC_PORT="18080"
COMPOSE_PROJECT="alien-panel"

REPO_URL="${XPANEL_REPO_URL:-$DEFAULT_REPO}"
REF="${XPANEL_REF:-$DEFAULT_REF}"
INSTALL_DIR="${XPANEL_INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
PUBLIC_PORT="${XPANEL_PUBLIC_PORT:-$DEFAULT_PUBLIC_PORT}"
DOMAIN="${XPANEL_DOMAIN:-}"
ADMIN_USER="${XPANEL_ADMIN_USER:-admin}"
BIND_ADDRESS="${XPANEL_BIND_ADDRESS:-}"

usage() {
  cat <<'EOF'
Alien-Panel production installer

Usage:
  install.sh [options]

Options:
  --repo URL       GitHub repository URL (default: Arm0ne/Alien-Panel)
  --ref REF        Branch, tag, or commit archive to deploy (default: main)
  --domain HOST    Public hostname; binds the temporary port to loopback
  --port PORT      Built-in HTTP port (default: 18080)
  --bind ADDRESS   0.0.0.0 or 127.0.0.1 (default depends on --domain)
  --dir PATH       Installation directory (default: /opt/alien-panel)
  --admin-user NAME
                   Initial administrator username (default: admin)
  --help           Show this help

The first install creates <dir>/.env with random administrator and Agent
registration secrets. Re-running the command keeps that file and the Docker
database volume so it can be used for upgrades.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo) REPO_URL="${2:-}"; shift 2 ;;
    --ref) REF="${2:-}"; shift 2 ;;
    --domain) DOMAIN="${2:-}"; shift 2 ;;
    --port) PUBLIC_PORT="${2:-}"; shift 2 ;;
    --bind) BIND_ADDRESS="${2:-}"; shift 2 ;;
    --dir) INSTALL_DIR="${2:-}"; shift 2 ;;
    --admin-user) ADMIN_USER="${2:-}"; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root, for example: sudo bash install.sh" >&2
  exit 1
fi

if [[ -z "$REPO_URL" || -z "$REF" || -z "$INSTALL_DIR" || -z "$ADMIN_USER" ]]; then
  echo "Repository, ref, installation directory, and admin user cannot be empty." >&2
  exit 2
fi
if [[ "$INSTALL_DIR" != /* || "$INSTALL_DIR" == "/" ]]; then
  echo "--dir must be an absolute path other than /." >&2
  exit 2
fi
if [[ ! "$PUBLIC_PORT" =~ ^[0-9]+$ ]] || (( PUBLIC_PORT < 1024 || PUBLIC_PORT > 65535 )); then
  echo "--port must be an unprivileged TCP port between 1024 and 65535." >&2
  exit 2
fi
if [[ -n "$DOMAIN" && "$DOMAIN" =~ [^A-Za-z0-9._-] ]]; then
  echo "--domain must contain only letters, numbers, dots, underscores, and hyphens." >&2
  exit 2
fi
if [[ -z "$BIND_ADDRESS" ]]; then
  if [[ -n "$DOMAIN" ]]; then BIND_ADDRESS="127.0.0.1"; else BIND_ADDRESS="0.0.0.0"; fi
fi
if [[ "$BIND_ADDRESS" != "0.0.0.0" && "$BIND_ADDRESS" != "127.0.0.1" ]]; then
  echo "--bind must be 0.0.0.0 or 127.0.0.1." >&2
  exit 2
fi
case "$(uname -m)" in
  x86_64|amd64) ;;
  *)
    echo "This release targets Linux amd64; detected $(uname -m)." >&2
    exit 1
    ;;
esac

install_docker() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    return
  fi
  if command -v docker-compose >/dev/null 2>&1; then
    return
  fi

  echo "Docker Compose is missing; installing Docker and Compose..."
  if command -v apt-get >/dev/null 2>&1; then
    apt-get update
    if ! apt-get install -y curl ca-certificates tar coreutils docker.io docker-compose-v2; then
      apt-get install -y curl ca-certificates tar coreutils docker.io docker-compose
    fi
    systemctl enable --now docker
  elif command -v dnf >/dev/null 2>&1; then
    dnf install -y curl ca-certificates tar coreutils docker docker-compose-plugin || \
      dnf install -y curl ca-certificates tar coreutils docker docker-compose
    systemctl enable --now docker
  else
    echo "No supported package manager found. Install Docker Engine and Compose, then rerun." >&2
    exit 1
  fi
}

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 24
  else
    od -An -N24 -tx1 /dev/urandom | tr -d ' \n'
  fi
}

repo_path="${REPO_URL#https://github.com/}"
repo_path="${repo_path#http://github.com/}"
repo_path="${repo_path%.git}"
repo_path="${repo_path%/}"
if [[ "$repo_path" == "$REPO_URL" || "$repo_path" != */* || "$repo_path" == *[^A-Za-z0-9_.\/-]* ]]; then
  echo "Only a GitHub repository URL is supported: $REPO_URL" >&2
  exit 2
fi

install_docker
stage="$(mktemp -d /tmp/alien-panel-install.XXXXXX)"
cleanup() { rm -rf -- "$stage"; }
trap cleanup EXIT

archive="$stage/source.tar.gz"
echo "Downloading $PRODUCT_NAME from $repo_path@$REF ..."
download_ok=0
for archive_url in \
  "https://codeload.github.com/$repo_path/tar.gz/refs/heads/$REF" \
  "https://codeload.github.com/$repo_path/tar.gz/refs/tags/$REF" \
  "https://codeload.github.com/$repo_path/tar.gz/$REF"; do
  if curl -fsSL --retry 3 --connect-timeout 10 --max-time 180 "$archive_url" -o "$archive"; then
    download_ok=1
    break
  fi
done
if (( download_ok == 0 )); then
  echo "Could not download $REPO_URL at ref $REF." >&2
  exit 1
fi
tar -xzf "$archive" -C "$stage"
# GitHub archives extract as owner-repo-ref/, not the package name
source_dir="$(find "$stage" -mindepth 1 -maxdepth 1 -type d -print -quit)"
if [[ -z "$source_dir" ]]; then
  echo "Could not find extracted repository directory." >&2
  exit 1
fi

required_files=(
  "$source_dir/VERSION"
  "$source_dir/RELEASE_MANIFEST.json"
  "$source_dir/deploy/docker-compose.yml"
  "$source_dir/deploy/Dockerfile.central"
  "$source_dir/deploy/nginx.docker.conf"
  "$source_dir/deploy/uninstall-docker.sh"
  "$source_dir/deploy/frontend-dist/index.html"
  "$source_dir/release/xpanel-central"
  "$source_dir/release/xpanel-agent"
)
for required in "${required_files[@]}"; do
  if [[ ! -f "$required" ]]; then
    echo "Package is missing: ${required#"$source_dir/"}" >&2
    exit 1
  fi
done

if [[ -f "$source_dir/checksums/release.sha256" ]]; then
  echo "Verifying release checksums..."
  (cd "$source_dir" && sha256sum -c checksums/release.sha256)
fi

assert_linux_amd64_elf() {
  local binary="$1"
  local magic class data machine
  magic="$(od -An -tx1 -N4 "$binary" | tr -d ' \n')"
  class="$(od -An -tu1 -j4 -N1 "$binary" | tr -d ' \n')"
  data="$(od -An -tu1 -j5 -N1 "$binary" | tr -d ' \n')"
  machine="$(od -An -tx1 -j18 -N2 "$binary" | tr -d ' \n')"
  if [[ "$magic" != "7f454c46" || "$class" != "2" || "$data" != "1" || "$machine" != "3e00" ]]; then
    echo "Release binary is not a Linux amd64 ELF executable: ${binary#"$source_dir/"}" >&2
    exit 1
  fi
}

for binary in \
  "$source_dir/release/xpanel-central" \
  "$source_dir/release/xpanel-agent" \
  "$source_dir/release/xpanel-db-maintenance" \
  "$source_dir/release/xpanel-traffic-check"; do
  assert_linux_amd64_elf "$binary"
done

echo "Installing to $INSTALL_DIR ..."
install -d -m 0755 "$INSTALL_DIR" "$INSTALL_DIR/deploy" "$INSTALL_DIR/release" "$INSTALL_DIR/frontend-dist"
cp -a "$source_dir/deploy/." "$INSTALL_DIR/deploy/"
rm -rf -- "$INSTALL_DIR/deploy/frontend-dist"
cp -a "$source_dir/deploy/frontend-dist/." "$INSTALL_DIR/frontend-dist/"
cp -a "$source_dir/release/." "$INSTALL_DIR/release/"
cp -a "$source_dir/VERSION" "$source_dir/RELEASE_MANIFEST.json" "$INSTALL_DIR/"
install -m 0755 "$source_dir/install.sh" "$INSTALL_DIR/install.sh"
chmod 0755 "$INSTALL_DIR/release/xpanel-central" "$INSTALL_DIR/release/xpanel-agent" 2>/dev/null || true
chmod 0755 "$INSTALL_DIR/deploy/"*.sh 2>/dev/null || true

env_file="$INSTALL_DIR/.env"
if [[ ! -f "$env_file" ]]; then
  origins="http://localhost:$PUBLIC_PORT,http://127.0.0.1:$PUBLIC_PORT"
  if [[ -n "$DOMAIN" ]]; then
    origins="$origins,http://$DOMAIN,https://$DOMAIN"
  else
    for ip in $(hostname -I 2>/dev/null || true); do
      [[ "$ip" == *:* ]] && continue
      origins="$origins,http://$ip:$PUBLIC_PORT"
    done
  fi
  umask 077
  cat > "$env_file" <<EOF
XPANEL_ADMIN_USER=$ADMIN_USER
XPANEL_ADMIN_PASSWORD=$(random_secret)
XPANEL_AGENT_REGISTRATION_TOKEN=$(random_secret)
XPANEL_CORS_ORIGINS=$origins
XPANEL_PUBLIC_PORT=$PUBLIC_PORT
XPANEL_BIND_ADDRESS=$BIND_ADDRESS
EOF
  chmod 0600 "$env_file"
else
  echo "Keeping existing credentials in $env_file"
  existing_port="$(sed -n 's/^XPANEL_PUBLIC_PORT=//p' "$env_file" | tail -n 1)"
  existing_bind="$(sed -n 's/^XPANEL_BIND_ADDRESS=//p' "$env_file" | tail -n 1)"
  [[ -n "$existing_port" ]] && PUBLIC_PORT="$existing_port"
  [[ -n "$existing_bind" ]] && BIND_ADDRESS="$existing_bind"
  if ! grep -q '^XPANEL_PUBLIC_PORT=' "$env_file"; then
    printf '\nXPANEL_PUBLIC_PORT=%s\n' "$PUBLIC_PORT" >> "$env_file"
  fi
  if ! grep -q '^XPANEL_BIND_ADDRESS=' "$env_file"; then
    printf 'XPANEL_BIND_ADDRESS=%s\n' "$BIND_ADDRESS" >> "$env_file"
  fi
  chmod 0600 "$env_file"
fi

cd "$INSTALL_DIR"
if docker compose version >/dev/null 2>&1; then
  docker compose -p "$COMPOSE_PROJECT" -f deploy/docker-compose.yml up -d --build
else
  docker-compose -p "$COMPOSE_PROJECT" -f deploy/docker-compose.yml up -d --build
fi

echo
echo "$PRODUCT_NAME $(cat "$source_dir/VERSION") is running."
if [[ "$BIND_ADDRESS" == "127.0.0.1" ]]; then
  echo "Web: http://127.0.0.1:$PUBLIC_PORT (use host Nginx/Caddy for HTTPS)"
else
  echo "Web: http://<server-ip>:$PUBLIC_PORT"
fi
echo "Health: http://127.0.0.1:$PUBLIC_PORT/health/ready"
echo "Credentials and Agent bootstrap token: $env_file"
echo "Status: docker compose -p $COMPOSE_PROJECT -f $INSTALL_DIR/deploy/docker-compose.yml ps"
