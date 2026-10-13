#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose -p kotohiro-local -f compose.local.yaml up -d --wait db mailpit

# DB tests create and drop their own schema; normal application rows are not
# modified. Real deliveries are confined to the local Mailpit receiver.
export MAIL_TEST_DATABASE_URL='postgres://kotohiro:kotohiro@127.0.0.1:5433/kotohiro?sslmode=disable'
export ADMIN_TEST_DATABASE_URL="$MAIL_TEST_DATABASE_URL"
export MAIL_TEST_SMTP_ADDR='127.0.0.1:1025'
export MAIL_TEST_HTTP_URL='http://127.0.0.1:8025'
cd server
go test -race -count=1 ./internal/infrastructure/email/... ./internal/infrastructure/external/aws/ses ./internal/domain/model/auth/password ./internal/domain/service ./internal/domain/service/organization ./internal/application/usecase/admin_usecase ./internal/application/usecase/auth_usecase ./internal/presentation/handler ./internal/infrastructure/persistence/query/user
