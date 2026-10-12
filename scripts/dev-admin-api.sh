#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose -p kotohiro-admin-local -f compose.admin-local.yaml up -d --wait db

# Override any ancestor .env before config.LoadConfig reads it. Keep this DB
# separate from existing API/Web development and never migrate a remote DB.
export ENV=local PORT=8081 DOMAIN=localhost
export DATABASE_URL='postgres://kotohiro:kotohiro@127.0.0.1:5434/kotohiro?sslmode=disable'
export TOKEN_SECRET='local-admin-development-only-secret'
export ADMIN_JWT_SECRET='local-admin-jwt-only-secret-at-least-32-bytes'
export ADMIN_ALLOWED_ORIGINS='http://localhost:3300,http://localhost:8081'
export ENCRYPTION_VERSION=v1 ENCRYPTION_SECRET='local-dev-key16!'
export HASH_PEPPER='local-admin-development-only-pepper' HASH_ITERATIONS=4
export ANALYSIS_API_DOMAIN='http://127.0.0.1:8091'
export INITIAL_ADMIN_EMAIL="${LOCAL_ADMIN_EMAIL:-}"
export INITIAL_ADMIN_NAME='ローカル管理者'
export BASELIME_API_KEY='' SENTRY_DSN=''
cd server
exec go run ./cmd/server
