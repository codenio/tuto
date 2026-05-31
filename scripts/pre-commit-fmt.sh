#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

if ! command -v goimports >/dev/null 2>&1; then
  echo "Installing goimports..."
  go install golang.org/x/tools/cmd/goimports@latest
fi

make fmt

# Re-stage formatted Go files so the commit includes corrections automatically.
git diff --name-only | while read -r f; do
  case "$f" in
    *.go) git add "$f" ;;
  esac
done
