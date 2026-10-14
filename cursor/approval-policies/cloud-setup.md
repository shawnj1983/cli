# Cloud setup policy

Applies to `.cursor/environment.json`, `.cursor/install.sh`, and `.cursor/rules/`.

These files may be approved when they only install tools, unset `NO_COLOR` for login shells, or add agent operating guidance, and Bugbot is clean or not present.

Do not approve if the change adds secrets, weakens egress, or runs untrusted remote scripts beyond the pinned golangci-lint installer already used here.
