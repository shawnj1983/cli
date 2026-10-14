# Automation template: CI Gate

Trigger: pull request synchronized, or a comment `/verify`.

## Prompt (paste into the Cursor automation)

Check out the PR head. Run `bash cursor/fleet/dispatch.sh` then `bash script/agent-verify.sh`. If verify fails, leave a non-blocking comment with the failing command and the last 30 log lines. Do not approve. Do not rewrite tests to hide `NO_COLOR`. If verify passes, comment only when the PR has no passing `Agent verify` check yet.
