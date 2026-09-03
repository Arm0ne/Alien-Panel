#!/usr/bin/env bash
set -Eeuo pipefail

# Read-only staging smoke check. Login is optional and is used only when
# XPANEL_ADMIN_PASSWORD is supplied; the access token is never printed.
base_url="${XPANEL_STAGING_URL:-http://127.0.0.1:8090}"
base_url="${base_url%/}"
admin_user="${XPANEL_ADMIN_USER:-admin}"

health="$(curl --fail --silent --show-error "$base_url/health/live")"
ready="$(curl --fail --silent --show-error "$base_url/health/ready")"
grep -q '"code":"0000"' <<<"$health" || { echo "health/live returned an unexpected envelope" >&2; exit 1; }
grep -q '"status":"ok"' <<<"$health" || { echo "health/live is not ok" >&2; exit 1; }
grep -q '"code":"0000"' <<<"$ready" || { echo "health/ready returned an unexpected envelope" >&2; exit 1; }
grep -q '"status":"ready"' <<<"$ready" || { echo "health/ready is not ready" >&2; exit 1; }

headers="$(curl --fail --silent --show-error --dump-header - --output /dev/null "$base_url/health/live")"
grep -qi '^x-content-type-options:[[:space:]]*nosniff' <<<"$headers" || { echo "missing X-Content-Type-Options" >&2; exit 1; }
grep -qi '^x-frame-options:[[:space:]]*DENY' <<<"$headers" || { echo "missing X-Frame-Options" >&2; exit 1; }
grep -qi '^cache-control:[[:space:]]*no-store' <<<"$headers" || { echo "missing Cache-Control: no-store" >&2; exit 1; }

if [[ -n "${XPANEL_ADMIN_PASSWORD:-}" ]]; then
  command -v jq >/dev/null 2>&1 || { echo "jq is required when XPANEL_ADMIN_PASSWORD is set" >&2; exit 1; }
  login_payload="$(jq -n --arg user "$admin_user" --arg pass "$XPANEL_ADMIN_PASSWORD" '{userName:$user,password:$pass}')"
  login="$(curl --fail --silent --show-error \
    -H 'Content-Type: application/json' \
    --data-binary @- "$base_url/api/auth/login" <<<"$login_payload")"
  token="$(jq -r 'select(.code == "0000") | .data.token // empty' <<<"$login")"
  [[ -n "$token" ]] || { echo "administrator login failed" >&2; exit 1; }
  dashboard="$(curl --fail --silent --show-error -H "Authorization: Bearer $token" "$base_url/api/dashboard")"
  grep -q '"code":"0000"' <<<"$dashboard" || { echo "protected dashboard request failed" >&2; exit 1; }
fi

echo "staging smoke check passed: $base_url"
