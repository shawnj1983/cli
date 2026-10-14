#!/usr/bin/env bash
# Fast correctness gate for Cloud Agents and fork CI.
# Unsets NO_COLOR so survey color tests are not false failures.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

echo "==> toolchain"
go version
golangci-lint version | head -n 1

echo "==> build"
make

echo "==> unit (NO_COLOR unset)"
env -u NO_COLOR go test -count=1 ./pkg/surveyext ./pkg/cmd/issue/list ./internal/agents

echo "==> lint"
make lint

echo "==> binary"
if [[ -x ./bin/gh ]]; then
  ./bin/gh --version
else
  echo "bin/gh missing after make" >&2
  exit 1
fi

echo "agent-verify: ok"
