- product: documentation
  boundary: Markdown docs such as AGENTS.md, README files, and other non-policy *.md files
  policies:
    - cursor/approval-policies/docs.md
- product: go-cli
  boundary: Go sources, go.mod, go.sum, Makefile, and script/
  policies:
    - cursor/approval-policies/go-cli.md
- product: cloud-setup
  boundary: .cursor/environment.json, .cursor/install.sh, and .cursor/rules/
  policies:
    - cursor/approval-policies/cloud-setup.md
