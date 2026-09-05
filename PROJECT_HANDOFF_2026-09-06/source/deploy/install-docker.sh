#!/usr/bin/env bash
set -euo pipefail

# One-command installer for a public test server. It downloads only the
# release binary, frontend bundle and deployment files from a public GitHub
# repository; no Node.js or Go toolchain is required on the server.

INSTALL_DIR="${INSTALL_DIR:-/opt/xpanel-central}"
REPO_URL="${XPANEL_REPO_URL:-}"
REF="${XPANEL_REF:-main}"
PUBLIC_PORT="${XPANEL_PUBLIC_PORT:-18080}"
DOMAIN="${XPANEL_DOMAIN:-}"
ADMIN_USER="${XPANEL_ADMIN_USER:-admin}"
BIND_ADDRESS="${XPANEL_BIND_ADDRESS:-}"

usage() {
  cat <<'EOF'
Usage:
  install-docker.sh --repo https://github.com/OWNER/REPO[.git] [options]

Options:
  --repo URL       Public GitHub repository containing this project (required)
  --ref REF        Branch or tag to deploy (default: main)
  --domain HOST    Public hostname used to build the CORS allow-list
  --port PORT      Temporary HTTP port (default: 18080)
  --bind ADDRESS   Bind test port (default: 0.0.0.0 without domain, 127.0.0.1 with domain)
  --dir PATH       Installation directory (default: /opt/xpanel-central)
  --help           Show this help

The script creates a random administrator password and Agent registration
token on first install and writes them to <dir>/.env (mode 0600).
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
    --help|-h) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root, for example: sudo bash install-docker.sh ..." >&2
  exit 1
fi
if [[ -z "$REPO_URL" ]]; then
  echo "--repo is required" >&2
  usage >&2
  exit 2
fi
if [[ ! "$PUBLIC_PORT" =~ ^[0-9]+$ ]] || (( PUBLIC_PORT < 1024 || PUBLIC_PORT > 65535 )); then
  echo "--port must be an unprivileged TCP port between 1024 and 65535" >&2
  exit 2
fi
if [[ -z "$BIND_ADDRESS" ]]; then
  if [[ -n "$DOMAIN" ]]; then BIND_ADDRESS="127.0.0.1"; else BIND_ADDRESS="0.0.0.0"; fi
fi
if [[ ! "$BIND_ADDRESS" =~ ^(0\.0\.0\.0|127\.0\.0\.1)$ ]]; then
  echo "--bind must be 0.0.0.0 or 127.0.0.1" >&2
  exit 2
fi
case "$(uname -m)" in
  x86_64|amd64) ;;
  *)
    echo "This release bundle targets Linux amd64; server architecture is $(uname -m)." >&2
    echo "Build a matching bundle with deploy/build-bundle.ps1 -GoArch before installing." >&2
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

  if command -v apt-get >/dev/null 2>&1; then
    apt-get update
    apt-get install -y curl ca-certificates tar
    if ! apt-get install -y docker.io docker-compose-v2; then
      apt-get install -y docker.io docker-compose
    fi
    systemctl enable --now docker
  elif command -v dnf >/dev/null 2>&1; then
    dnf install -y curl ca-certificates tar docker docker-compose-plugin
    systemctl enable --now docker
  else
    echo "Docker Compose is missing and this OS has no supported package manager." >&2
    echo "Install Docker Engine + Compose plugin, then rerun this script." >&2
    exit 1
  fi
}

random_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 24
  else
    head -c 48 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c 48
  fi
}

repo_path="${REPO_URL#https://github.com/}"
repo_path="${repo_path#http://github.com/}"
repo_path="${repo_path%.git}"
repo_path="${repo_path%/}"
if [[ "$repo_path" == "$REPO_URL" || "$repo_path" != */* ]]; then
  echo "Only a GitHub repository URL is supported: $REPO_URL" >&2
  exit 2
fi

install_docker
stage="$(mktemp -d /tmp/xpanel-central.XXXXXX)"
cleanup() { rm -rf "$stage"; }
trap cleanup EXIT

echo "Downloading X-Panel Central $repo_path@$REF ..."
archive="$stage/source.tar.gz"
if ! curl -fsSL "https://codeload.github.com/$repo_path/tar.gz/refs/heads/$REF" -o "$archive"; then
  curl -fsSL "https://codeload.github.com/$repo_path/tar.gz/refs/tags/$REF" -o "$archive"
fi
tar -xzf "$archive" -C "$stage"
source_dir="$(find "$stage" -mindepth 1 -maxdepth 1 -type d | head -n 1)"
if [[ -z "$source_dir" ]]; then
  echo "Downloaded repository archive is empty" >&2
  exit 1
fi

required_files=(
  "$source_dir/deploy/docker-compose.yml"
  "$source_dir/deploy/Dockerfile.central"
  "$source_dir/deploy/nginx.docker.conf"
  "$source_dir/deploy/uninstall-docker.sh"
  "$source_dir/deploy/frontend-dist/index.html"
  "$source_dir/release/xpanel-central"
)
for required in "${required_files[@]}"; do
  if [[ ! -f "$required" ]]; then
    echo "Repository is missing $required" >&2
    echo "Build and commit the release bundle before installing." >&2
    exit 1
  fi
done

install -d -m 0755 "$INSTALL_DIR/deploy" "$INSTALL_DIR/release" "$INSTALL_DIR/frontend-dist"
cp -a "$source_dir/deploy/docker-compose.yml" "$source_dir/deploy/Dockerfile.central" "$source_dir/deploy/nginx.docker.conf" "$INSTALL_DIR/deploy/"
install -m 0755 "$source_dir/deploy/uninstall-docker.sh" "$INSTALL_DIR/deploy/uninstall-docker.sh"
cp -a "$source_dir/deploy/frontend-dist/." "$INSTALL_DIR/frontend-dist/"
install -m 0755 "$source_dir/release/xpanel-central" "$INSTALL_DIR/release/xpanel-central"

if [[ ! -f "$INSTALL_DIR/.env" ]]; then
  origins="http://localhost:$PUBLIC_PORT,http://127.0.0.1:$PUBLIC_PORT"
  if [[ -n "$DOMAIN" ]]; then
    origins="$origins,http://$DOMAIN,https://$DOMAIN"
  else
    for ip in $(hostname -I 2>/dev/null || true); do
      [[ "$ip" == *:* ]] && continue
      origins="$origins,http://$ip:$PUBLIC_PORT"
    done
    public_ip="$(curl -4 -fsSL --max-time 5 https://api.ipify.org 2>/dev/null || true)"
    if [[ "$public_ip" =~ ^[0-9.]+$ ]]; then
      origins="$origins,http://$public_ip:$PUBLIC_PORT"
    fi
  fi
  umask 077
  cat > "$INSTALL_DIR/.env" <<EOF
XPANEL_ADMIN_USER=$ADMIN_USER
XPANEL_ADMIN_PASSWORD=$(random_secret)
XPANEL_AGENT_REGISTRATION_TOKEN=$(random_secret)
XPANEL_CORS_ORIGINS=$origins
XPANEL_PUBLIC_PORT=$PUBLIC_PORT
XPANEL_BIND_ADDRESS=$BIND_ADDRESS
EOF
  chmod 0600 "$INSTALL_DIR/.env"
else
  echo "Keeping existing secrets in $INSTALL_DIR/.env"
  existing_port="$(sed -n 's/^XPANEL_PUBLIC_PORT=//p' "$INSTALL_DIR/.env" | tail -n 1)"
  if [[ -n "$existing_port" ]]; then
    PUBLIC_PORT="$existing_port"
  else
    printf '\nXPANEL_PUBLIC_PORT=%s\n' "$PUBLIC_PORT" >> "$INSTALL_DIR/.env"
  fi
  if ! grep -q '^XPANEL_BIND_ADDRESS=' "$INSTALL_DIR/.env"; then
    printf 'XPANEL_BIND_ADDRESS=%s\n' "$BIND_ADDRESS" >> "$INSTALL_DIR/.env"
  else
    existing_bind="$(sed -n 's/^XPANEL_BIND_ADDRESS=//p' "$INSTALL_DIR/.env" | tail -n 1)"
    [[ -n "$existing_bind" ]] && BIND_ADDRESS="$existing_bind"
  fi
fi

cd "$INSTALL_DIR"
if docker compose version >/dev/null 2>&1; then
  docker compose -p xpanel-central -f deploy/docker-compose.yml up -d --build
else
  docker-compose -p xpanel-central -f deploy/docker-compose.yml up -d --build
fi

echo
echo "X-Panel Central is running at: http://<server-ip>:$PUBLIC_PORT"
echo "Health check: http://<server-ip>:$PUBLIC_PORT/health/ready"
echo "Administrator credentials are stored in: $INSTALL_DIR/.env"
echo "For host Nginx + HTTPS, proxy your domain to 127.0.0.1:$PUBLIC_PORT."
