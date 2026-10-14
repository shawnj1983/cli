#!/usr/bin/env bash
# Idempotent Cloud Agent install for the gh CLI.
# Safe on Cursor's default image and on a later boot that already has tools.
set -euo pipefail

printf '%s\n' 'unset NO_COLOR' | sudo tee /etc/profile.d/00-cursor-unset-no-color.sh >/dev/null

if ! command -v golangci-lint >/dev/null || ! golangci-lint version 2>/dev/null | grep -q '2.11.0'; then
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sudo sh -s -- -b /usr/local/bin v2.11.0
fi

GOTOOLCHAIN=auto go mod download
