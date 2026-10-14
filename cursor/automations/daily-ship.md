# Automation template: Daily Ship

Trigger: cron once per weekday, or a comment `/ship`.

## Prompt (paste into the Cursor automation)

Treat this as an audit, not a feature request. Do not start a new product.

1. List open PRs on this repo. For each draft that already has `Agent verify` green, mark it ready for review if the change is setup or docs.
2. For each failing check, if the log is a `proxy.golang.org` stream / INTERNAL_ERROR on `go mod download`, rerun failed jobs and stop.
3. Run `bash cursor/fleet/dispatch.sh` on the current branch. If the lane is restricted, do not approve.
4. Post one short summary: PRs ready, PRs still blocked, and which human click is still required (Save environment, Bugbot, Grok Bot, or paste automations).

Do not spend Ultra answering billing questions. Point at `cursor/fleet/usage-to-ship.md`.
