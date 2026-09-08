#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/opt/alien-panel}"
COMPOSE_PROJECT="${COMPOSE_PROJECT:-alien-panel}"
COMPOSE_FILE="$INSTALL_DIR/deploy/docker-compose.yml"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "X-Panel Docker installation was not found at $INSTALL_DIR" >&2
  exit 1
fi

cd "$INSTALL_DIR"
if docker compose version >/dev/null 2>&1; then
  docker compose -p "$COMPOSE_PROJECT" -f deploy/docker-compose.yml down
elif command -v docker-compose >/dev/null 2>&1; then
  docker-compose -p "$COMPOSE_PROJECT" -f deploy/docker-compose.yml down
else
  echo "Docker Compose is not installed; containers were not changed." >&2
  exit 1
fi

echo "Containers stopped. Persistent database volume was kept."
echo "To remove data too, inspect and remove the named volume ${COMPOSE_PROJECT}_central-data explicitly."
