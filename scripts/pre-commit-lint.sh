#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

# Keep in sync with .github/workflows/integration.yml (golangci-lint v2).
GOLANGCI_LINT_VERSION=v2.12.2

go run "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}" run ./...
