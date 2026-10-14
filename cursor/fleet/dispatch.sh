#!/usr/bin/env bash
# Classify the current tree (or a PR) into a fleet lane and print the next step.
set -euo pipefail

root="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$root"

base="${FLEET_BASE:-origin/trunk}"
pr="${1:-}"

if [[ -n "$pr" ]]; then
  echo "==> checkout PR ${pr}"
  gh pr checkout "$pr"
fi

if ! git rev-parse --verify "$base" >/dev/null 2>&1; then
  git fetch origin trunk >/dev/null 2>&1 || true
  base="origin/trunk"
fi

mapfile -t files < <(git diff --name-only "${base}...HEAD" 2>/dev/null || git diff --name-only "$base")
if [[ ${#files[@]} -eq 0 ]]; then
  mapfile -t files < <(git diff --name-only)
fi

lane="docs"
reason="no risky paths"
for f in "${files[@]}"; do
  case "$f" in
    pkg/cmd/attestation/*|pkg/cmd/release/attestation/*|pkg/cmd/release/verify*|internal/oauth*|pkg/cmd/auth/*)
      lane="restricted"
      reason="$f"
      break
      ;;
    .github/workflows/*|APPROVAL_POLICY.md|cursor/approval-policies/*|cursor/approvals/*|cursor/automations/*)
      lane="restricted"
      reason="$f"
      break
      ;;
    *.go|go.mod|go.sum|Makefile|script/*)
      lane="code"
      reason="$f"
      ;;
    .cursor/*|cursor/fleet/*)
      if [[ "$lane" != "code" ]]; then
        lane="setup"
        reason="$f"
      fi
      ;;
  esac
done

echo "lane=${lane}"
echo "reason=${reason}"
echo "files=${#files[@]}"

case "$lane" in
  docs)
    echo "next: Cursor Models / Grok Bot. Approval policy: docs.md"
    echo "verify: skip full test suite"
    ;;
  setup)
    echo "next: Cursor Models. Run bash .cursor/install.sh && bash script/agent-verify.sh"
    echo "verify: bash script/agent-verify.sh"
    ;;
  code)
    echo "next: Cursor Models (Grok 4.6 / Composer 2.5). Run bash script/agent-verify.sh"
    echo "verify: bash script/agent-verify.sh"
    ;;
  restricted)
    echo "next: human review (shawnj1983). Do not auto-approve."
    echo "verify: bash script/agent-verify.sh"
    ;;
esac
