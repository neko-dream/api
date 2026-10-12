#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

# A separate Compose project and DB port preserve the normal development DB.
docker compose -p kotohiro-local -f compose.local.yaml up -d --wait db mailpit

# Explicit local values take precedence over .env (godotenv.Load does not
# override existing variables). Never run migrations against the .env DB.
export ENV=local PORT=8080 DOMAIN=localhost
export DATABASE_URL='postgres://kotohiro:kotohiro@127.0.0.1:5433/kotohiro?sslmode=disable'
export GOOGLE_CALLBACK_URL='https://localhost:3000/api/auth/google/callback'
export LINE_CALLBACK_URL='https://localhost:3000/api/auth/line/callback'
export TOKEN_SECRET='local-development-only-session-secret'
export HASH_PEPPER='local-development-only-password-pepper'
export HASH_ITERATIONS=10
export ENCRYPTION_VERSION=v1 ENCRYPTION_SECRET='local-dev-key16!'
export ANALYSIS_API_DOMAIN='http://127.0.0.1:8090'
export SMTP_ADDR='127.0.0.1:1025'
export EMAIL_FROM='ことひろ <noreply@localhost.test>'
export WEBSITE_URL='https://localhost:3000' APP_NAME='ことひろ'
export EMAIL_VERIFICATION_PAGE_URL='https://localhost:3000/auth/email/verify'

cd server
exec go run ./cmd/server
