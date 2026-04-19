---
name: release
description: Automate ecs9s release process with semver tagging
---

1. Run `go build -o ecs9s .` and `go test ./...`
2. Determine version bump (major/minor/patch) from recent commits
3. Update version in main.go if hardcoded
4. Create git tag: `git tag v<version>`
5. Build release binaries: `GOOS=linux GOARCH=amd64 go build -o ecs9s-linux-amd64 .`
6. Generate changelog from commits since last tag
