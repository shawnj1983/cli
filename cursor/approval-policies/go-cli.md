# Go CLI policy

Applies to `*.go`, `go.mod`, `go.sum`, `Makefile`, and `script/`.

Do not approve authentication, credential, attestation, or release-verification changes without a human reviewer.

Other Go changes may be approved only when Cursor Bugbot completed and reported no unresolved findings that need human attention. If Bugbot is missing, skipped, failed, or still running after the wait window, do not approve; require human review.
