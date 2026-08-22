#!/usr/bin/env bash
# CI Go test floor — CGO required for tree-sitter analyzers.
set -euo pipefail
export CGO_ENABLED=1

echo "ci-go-test: go test (packages)…"
go test ./internal/... ./cmd/trace/... ./cmd/trace-mcp/... -count=1 -timeout=15m

echo "ci-go-test: evals (stable packages)…"
go test ./evals/capability/... ./evals/honesty/... ./evals/impact/... ./evals/p0x/... ./evals/replan/... ./evals/x0/... -count=1 -timeout=10m

echo "ci-go-test: go vet…"
go vet ./internal/... ./cmd/trace/... ./cmd/trace-mcp/... ./evals/...

echo "ci-go-test: build trace + trace-mcp…"
go build -o /dev/null ./cmd/trace
go build -o /dev/null ./cmd/trace-mcp

echo "ci-go-test: OK"
