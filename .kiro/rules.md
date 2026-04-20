# Project Rules — ecs9s

Condensed checklist for Kiro agents. Full context lives in `.kiro/steering/*.md` and `.kiro/docs/architecture.md`.

## Must Do

1. **Read first**: `AGENTS.md`, then `.kiro/steering/{product,tech,structure,conventions,security}.md`.
2. **Match structure**: one file per AWS service in `internal/aws/`; one file per page in `internal/ui/pages/`; lowercase filenames.
3. **Async AWS only**: AWS SDK calls go through `tea.Cmd`, results returned as typed `fooLoadedMsg{data, err}`. No blocking I/O in `Update` / `View`.
4. **Confirm destructive actions**: force deploy, scale, stop task, rollback, any delete — route through the `confirm` component with resource identifier shown.
5. **Mask secrets**: fetch SSM SecureString with `WithDecryption: false`; never print Secrets Manager values.
6. **Cost from real data**: use `DescribeTaskDefinition` CPU/memory; no name-based inference.
7. **Rollback safety**: find `PRIMARY` deployment, verify target task def exists before `UpdateService`.
8. **SSM target**: `ecs:{cluster-name}_{task-id}_{runtime-id}` — short names, not ARNs.
9. **Adding a page**: update `internal/ui/messages.go` (PageType + command), `internal/ui/pages/<name>.go`, and the 5 spots in `internal/app/app.go` (struct field, `New()`, `resize()`, `initCurrentPage()`, `updateActivePage()`/`viewActivePage()`).
10. **Keep docs in sync**: structural / pattern changes → update `.kiro/docs/architecture.md`, `docs/architecture.md`, `AGENTS.md`, and relevant `.kiro/steering/*`.

## Must Not

- Never commit AWS credentials, tokens, or `.env` files. Run secret scanning before commit.
- Never log or echo decrypted `SecureString` / Secrets Manager values.
- Never call AWS SDK from inside `Update` or `View` synchronously.
- Never bypass the confirm dialog for destructive operations.
- Never guess task definition cost from family names — always resolve real values.
- Never use ARNs where short ECS names are required (SSM exec / port forward targets).

## Commands

```bash
go build -o ecs9s .          # build
go test ./...                 # test
go vet ./...                  # lint
./ecs9s --profile X --region Y --theme light
```

## References

- [`AGENTS.md`](../AGENTS.md) — project entry doc
- [`.kiro/steering/product.md`](steering/product.md)
- [`.kiro/steering/tech.md`](steering/tech.md)
- [`.kiro/steering/structure.md`](steering/structure.md)
- [`.kiro/steering/conventions.md`](steering/conventions.md)
- [`.kiro/steering/security.md`](steering/security.md)
- [`.kiro/docs/architecture.md`](docs/architecture.md) → [`docs/architecture.md`](../docs/architecture.md)
