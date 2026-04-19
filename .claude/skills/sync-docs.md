---
name: sync-docs
description: Synchronize documentation with current code state
---

1. Scan `internal/ui/pages/` for all page files → update CLAUDE.md page count
2. Scan `internal/aws/` for all client files → update CLAUDE.md service list
3. Scan `internal/ui/messages.go` for PageType constants → verify tab items match in app.go
4. Check `docs/architecture.md` reflects current component list
5. Verify module CLAUDE.md files exist for all `internal/` subdirectories
6. Report quality score: files checked, issues found, auto-fixed
