# Approval policy

Closest policy under `cursor/approval-policies/` wins for matching files. This file is the fallback.

## Default

- Do not approve when Cursor Bugbot reported unresolved findings, failed analysis, or incomplete results.
- A missing or skipped Bugbot check is not by itself a reason to withhold approval.
- Do not assign GitHub teams from `.github/CODEOWNERS` (`@cli/*`). Those teams are upstream-only and are not reviewers on this fork.
- When a human reviewer is required and the assignment tool is available, assign at most `shawnj1983`.

## Require human review

Require human review (do not approve) when the PR changes any of:

- authentication, tokens, credentials, or HTTP signing
- release, attestation, or verification code
- GitHub workflow files under `.github/workflows/`
- `APPROVAL_POLICY.md`, `cursor/approval-policies/`, or `cursor/approvals/`
- more than a documentation or Cloud Agent setup change and Bugbot is not present

## Documentation and setup

Documentation-only changes (`*.md` that is not a policy file) and Cloud Agent setup under `.cursor/` may be approved when Bugbot is clean or not present, and no required-human-review path above is included.
