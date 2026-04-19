Run the full test suite and report results.

```bash
go vet ./...
go test ./... -v -count=1
go build -o /dev/null .
```

If tests fail, analyze the failure and suggest fixes.
If no test files exist yet, report which packages need tests.
