# Usage to ship

Ultra spend is the cost. A merged PR is the return. Anything else is waste.

## Daily rule

1. Classify with `bash cursor/fleet/dispatch.sh`.
2. Use Cursor Models or Grok Bot for docs and setup. Do not open a Cloud Agent on Ultra for account questions.
3. Run `make agent-verify` before opening or updating a PR.
4. Mark the PR ready. Do not leave drafts sitting.
5. If CI fails on `proxy.golang.org` / `INTERNAL_ERROR` during `go mod download`, rerun failed jobs. That is not a code defect.

## What this fork sells

This repository is the `gh` CLI, not a storefront. Revenue here means: fewer wasted Ultra minutes per merged change. Track it as:

- minutes on Cloud Agents that produced a merge
- minutes on chats that produced no PR
- PRs left draft or unapproved for more than one day

## Human-only clicks

These cannot be done from an agent. They unblock everything else:

- Save the Cloud Agent environment
- Paste `cursor/automations/*.md` into Cursor Automations
- Enable Cursor Bugbot
- Install Grok Bot and sign in as the Ultra account
