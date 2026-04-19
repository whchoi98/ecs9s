---
name: code-review
description: Quick inline review of staged/recent Go changes. For background review of large diffs, use the code-reviewer agent instead.
---

Quick review of the current git diff. For large multi-file reviews, dispatch the `code-reviewer` agent.

Review for:
1. **Security**: SecureString decryption, AWS credential exposure, SQL injection in SSM paths
2. **Bubbletea patterns**: Proper message handling, no blocking calls in Update(), correct tea.Cmd usage
3. **AWS SDK**: Context propagation, error wrapping, pagination handling
4. **Conventions**: File naming, page registration (5-point checklist), NavContext drill-down safety
5. **Naming**: Consistent with existing patterns (e.g., `fooLoadedMsg`, `NewFooPage`, `FooPage.fetchData`)

Run: `go vet ./...` and report any issues.
