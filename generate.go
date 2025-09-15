//go:build generate
// +build generate

package main

//go:generate sh -c "cd tsp && tsp format main.tsp && npm run build"
//go:generate ogen --package oas --target internal/presentation/oas --clean ./static/oas/openapi.yaml --convenient-errors=on
//go:generate sqlc generate
//go:generate oapi-codegen -config oapi.yaml ./api/analysis.openapi.json
//go:generate sh -c "find . -name '*.go' | grep -v 'vendor/\\|.git/\\|_test.go' | xargs -n 1 -t otelinji -template ./internal/infrastructure/telemetry/otelinji.template -w -filename &>/dev/null || true"
