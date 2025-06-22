#!/bin/bash

set -e

# Compile TypeSpec to OpenAPI
cd tsp
tsp format main.tsp
npm run build
cd ../

# Generate Go code from OpenAPI
ogen --package oas --target internal/presentation/oas --clean ./static/oas/openapi.yaml --convenient-errors=on

sqlc generate
oapi-codegen -config oapi.yaml ./api/analysis.openapi.json

find . -name "*.go" | grep -v "vendor/\|.git/\|_test.go" | xargs -n 1 -t otelinji -template "./internal/infrastructure/telemetry/otelinji.template" -w -filename &> /dev/null

cd admin-ui
bun run build
cd ../

