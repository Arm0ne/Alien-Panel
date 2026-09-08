#!/usr/bin/env bash
set -euo pipefail

# Backward-compatible entry point. The package-root install.sh is the single
# source of truth and is safe to run from a local checkout or via curl.
PACKAGE_ROOT="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -f "$PACKAGE_ROOT/install.sh" ]]; then
  exec bash "$PACKAGE_ROOT/install.sh" "$@"
fi

RAW_INSTALLER_URL="${ALIEN_PANEL_INSTALLER_URL:-https://raw.githubusercontent.com/Arm0ne/Alien-Panel/main/install.sh}"
exec curl -fsSL "$RAW_INSTALLER_URL" | bash -s -- "$@"
