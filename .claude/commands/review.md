Review the current git diff for code quality, security, and ecs9s conventions.

Run `go vet ./...` first, then analyze:
- Bubbletea patterns (no blocking in Update, proper Cmd usage)
- AWS SDK usage (context, error wrapping)
- Security (no secret exposure, proper masking)
- Page registration completeness (messages.go + app.go sync)
