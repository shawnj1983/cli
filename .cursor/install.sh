#!/usr/bin/env bash
# Idempotent Cloud Agent install for the gh CLI.
set -euo pipefail

if ! command -v sudo >/dev/null; then
  echo "sudo is required" >&2
  exit 1
fi

printf '%s\n' 'unset NO_COLOR' | sudo tee /etc/profile.d/00-cursor-unset-no-color.sh >/dev/null

if ! command -v golangci-lint >/dev/null || ! golangci-lint version 2>/dev/null | grep -q '2.11.0'; then
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v2.11.0/install.sh | sudo sh -s -- -b /usr/local/bin v2.11.0
fi

GOTOOLCHAIN=auto go mod download
GOTOOLCHAIN=auto go version >/dev/null
