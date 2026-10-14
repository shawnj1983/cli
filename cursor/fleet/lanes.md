# Fleet lanes

Run `bash cursor/fleet/dispatch.sh` or `bash cursor/fleet/dispatch.sh <pr>` before starting work.

| Lane | Paths | Model | Verify | Approve |
| --- | --- | --- | --- | --- |
| docs | Markdown that is not a policy | Cursor Models or Grok Bot | Skip full tests | Yes if Bugbot clean or missing |
| setup | `.cursor/*`, `cursor/fleet/*` | Cursor Models | `script/agent-verify.sh` | Yes if no secrets or egress weakening |
| code | Go, modules, Makefile, `script/` | Cursor Models | `script/agent-verify.sh` | Only if Bugbot completed clean |
| restricted | auth, attestation, workflows, approval policies, automation prompts | Human + stronger model if needed | `script/agent-verify.sh` | Never auto-approve |

One automation should own one lane. Do not spawn a second Cloud Agent on the same PR head unless the first is stuck.

The live router is [Pull Request Router](https://cursor.com/automations/2dd0a183-b81b-11f1-977f-f6b8f2fcf9b2). Paste templates from `cursor/automations/` into new Cursor automations; this repo cannot create dashboard automations for you.

Convert Ultra minutes into merges with `cursor/fleet/usage-to-ship.md` and the Daily Ship template.
