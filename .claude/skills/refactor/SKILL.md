---
name: refactor
description: Refactor Go code for clarity, DRY, and maintainability
---

Identify refactoring opportunities in recently changed files:
1. Extract shared page patterns (filter handling, table setup) into base helpers if 3+ pages duplicate logic
2. Consolidate AWS client initialization if constructors are repetitive
3. Ensure single responsibility per file
4. Check for dead code or unused exports

Run `go build ./...` after each change to verify compilation.
