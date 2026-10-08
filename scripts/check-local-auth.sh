#!/usr/bin/env bash
set -euo pipefail
task_origin="${1:-https://localhost:3000}"
if [[ ! "$task_origin" =~ ^https://localhost:[0-9]+$ ]]; then
  echo 'Use the local HTTPS Web origin (https://localhost:3000).' >&2
  exit 1
fi
task_tmp=$(mktemp -d)
trap 'rm -f "$task_tmp/cookies" "$task_tmp/headers" "$task_tmp/body"; rmdir "$task_tmp"' EXIT

# --insecure is scoped to this localhost self-signed certificate.
task_status=$(curl --insecure --silent --show-error --request POST \
  --dump-header "$task_tmp/headers" --cookie-jar "$task_tmp/cookies" \
  --output "$task_tmp/body" --write-out '%{http_code}' \
  "$task_origin/api/auth/guest/login")
[[ "$task_status" == 200 ]] || { echo "Guest login failed: HTTP $task_status" >&2; exit 1; }
grep -qi 'set-cookie: SessionId=' "$task_tmp/headers"
grep -qi '; Secure' "$task_tmp/headers"
grep -qi '; HttpOnly' "$task_tmp/headers"
grep -qi '; SameSite=Lax' "$task_tmp/headers"
if grep -qi '; Domain=' "$task_tmp/headers"; then
  echo 'Local cookie must be host-only.' >&2
  exit 1
fi

task_status=$(curl --insecure --silent --show-error --cookie "$task_tmp/cookies" \
  --output "$task_tmp/body" --write-out '%{http_code}' "$task_origin/api/auth/token/info")
[[ "$task_status" == 200 ]] || { echo "Cookie authentication failed: HTTP $task_status" >&2; exit 1; }

task_status=$(curl --insecure --silent --show-error --request POST \
  --cookie "$task_tmp/cookies" --cookie-jar "$task_tmp/cookies" \
  --dump-header "$task_tmp/headers" --output "$task_tmp/body" --write-out '%{http_code}' \
  "$task_origin/api/auth/revoke")
[[ "$task_status" == 204 ]] || { echo "Logout failed: HTTP $task_status" >&2; exit 1; }
grep -qi 'Max-Age=0' "$task_tmp/headers"

task_status=$(curl --insecure --silent --show-error --cookie "$task_tmp/cookies" \
  --output "$task_tmp/body" --write-out '%{http_code}' "$task_origin/api/auth/token/info")
[[ "$task_status" == 403 ]] || { echo "Logged-out request must be forbidden: HTTP $task_status" >&2; exit 1; }
echo 'PASS: HTTPS proxy, host-only Secure/HttpOnly Cookie, authentication, logout and cookie expiration.'
