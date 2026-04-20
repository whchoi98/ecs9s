# Conventions

## Bubbletea MVU

Every page implements:

- `Init() tea.Cmd` — usually dispatches the first data fetch.
- `Update(msg tea.Msg) (Model, tea.Cmd)` — pure state transitions.
- `View() string` — renders current state.

No blocking I/O inside `Update` / `View`. Use `tea.Cmd` goroutines for AWS calls and return a typed `fooLoadedMsg{data, err}`.

## Error Handling

- AWS errors returned through `fooLoadedMsg.err`.
- Surface user-visible errors via `ui.ErrorMsg` to the status bar. Do not panic.
- Log at caller level only when useful; avoid noisy logs in the hot path.

## Navigation & Commands

- Tab / Shift+Tab or `[` / `]` → switch pages.
- `:` opens command mode. Command strings are lowercase service aliases (`cluster`, `service`, `ecr`…) defined in `internal/ui/messages.go`.
- `Enter` drills down; `Esc` / `Backspace` pops `navStack`.
- `NavContext` carries `ClusterARN`, `ServiceName` etc. through drill-down chains. Never parse them back out of view strings — always pass via `NavContext`.

## Destructive Actions

- Every destructive action (force deploy, scale, stop task, rollback, delete) must route through the `confirm` component before executing.
- Include resource identifiers in the confirm dialog so the user can verify target.

## AWS-Specific Rules

- **SecureString**: call SSM with `WithDecryption: false`. Do not cache decrypted values in memory.
- **Cost estimates**: read actual CPU/memory from `DescribeTaskDefinition`. Do not infer from task family names.
- **SSM target format**: `ecs:{cluster-name}_{task-id}_{runtime-id}` — use short names, not ARNs.
- **Rollback safety**: locate the `PRIMARY` deployment first; verify the target task definition exists before calling `UpdateService`.

## Tests

- Unit tests live next to the file: `foo.go` ↔ `foo_test.go`.
- Run all: `go test ./...`
- Integration / shell tests live under `tests/`; run via `tests/run-all.sh`.

## Commits & Docs Sync

When design or structure changes:

1. Update `.kiro/docs/architecture.md`.
2. Update `AGENTS.md` if patterns/conventions changed.
3. Keep `CLAUDE.md` aligned with `AGENTS.md` until the legacy file is removed.
