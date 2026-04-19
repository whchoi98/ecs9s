Build release binaries for all target platforms.

```bash
# Verify first
go vet ./...
go test ./... -count=1

# Create dist directory
mkdir -p dist

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/ecs9s-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o dist/ecs9s-linux-arm64 .
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o dist/ecs9s-darwin-arm64 .
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/ecs9s-darwin-amd64 .
```

Report binary sizes and verify each builds successfully.

## Error Recovery

If a build step fails:
1. **go vet fails**: Fix the reported issues before proceeding. Do not build release binaries with vet warnings.
2. **go test fails**: Investigate test failures. Do not ship broken tests. Run `go test -v ./...` for details.
3. **Platform build fails**: Check `go version` supports the target GOOS/GOARCH. Minimum Go 1.22+. If a specific platform fails, build the others and report which platform failed with the exact error.
4. **dist/ permission error**: Run `mkdir -p dist` and verify write permissions.
5. **Binary too large**: Verify `-ldflags "-s -w"` is applied (strips debug info). Expected size: 15-20MB per binary.

## Rollback

If a release is found to be broken after distribution:
1. Remove the faulty binaries from dist/
2. Revert to the last known good git tag: `git checkout <previous-tag>`
3. Rebuild from that tag
