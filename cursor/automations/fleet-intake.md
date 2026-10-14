# Automation template: Fleet Intake

Trigger: issue labeled `agent` or a Slack/message that asks for implementation.

## Prompt (paste into the Cursor automation)

Classify the request into a lane from `cursor/fleet/lanes.md` using the file paths the user named. If paths are unclear, inspect the repo and pick the smallest lane.

- docs: answer or edit Markdown only. Do not open a PR unless asked.
- setup: edit `.cursor/` or `cursor/fleet/` and run `script/agent-verify.sh`.
- code: implement, run `script/agent-verify.sh`, open a PR.
- restricted: do not implement. Comment that `shawnj1983` must take it.

Use Cursor Models. One task, one branch, one PR. Do not start extra Cloud Agents for the same issue.
