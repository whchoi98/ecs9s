# ecs9s

ECS-focused Terminal UI for AWS — manage clusters, services, tasks, containers, and related operational services with a k9s-inspired keyboard-driven interface.

> Kiro agents: Read `.kiro/rules.md` and `.kiro/steering/*.md` for project rules, tech stack, structure, and conventions. See `.kiro/docs/architecture.md` for system design.

## Tech Stack

- **Language**: Go 1.24+
- **TUI Framework**: Bubbletea (Elm architecture) + Bubbles + Lipgloss
- **AWS SDK**: AWS SDK for Go v2
- **Config**: YAML (`~/.ecs9s/config.yaml`)
- **Build**: `go build` → single static binary

## Key Commands

```bash
go build -o ecs9s .          # Build
./ecs9s                       # Run (default profile/region)
./ecs9s --profile X --region Y --theme light
go test ./...                 # Test
go vet ./...                  # Lint
```

## Project Structure

```
ecs9s/
├── main.go                    # Entry point
├── internal/
│   ├── app/                   # Root bubbletea model, routing
│   ├── ui/{messages,components,pages,styles}
│   ├── aws/                   # AWS SDK v2 client wrappers
│   ├── action/                # Mutation operations
│   ├── config/                # YAML config
│   └── theme/                 # Preset themes
├── docs/                      # Design specs, architecture
├── tests/                     # Test suites
├── scripts/                   # Setup scripts
└── .kiro/                     # Kiro steering, rules, docs
```

## Architecture Highlights

- **Bubbletea Model-Update-View**: every page implements `Init()`, `Update(msg)`, `View()`.
- **Async AWS I/O**: wrap SDK calls in `tea.Cmd` returning `fooLoadedMsg{data, err}`.
- **Navigation**: Tab/Shift+Tab or `[`/`]`; command mode (`:cluster`, `:service`, `:ecr`…); drill-down via `Enter`, back via `Esc`/`Backspace`, `NavContext` carries ClusterARN/ServiceName.
- **Adding a page**: 5 edits in `internal/app/app.go` (field, `New()`, `resize()`, `initCurrentPage()`, `updateActivePage()`/`viewActivePage()`) + PageType constant and command string in `internal/ui/messages.go`.

## Conventions

- Lowercase file names; one file per AWS service client or page model.
- AWS errors surface via `ui.ErrorMsg` to status bar.
- Never decrypt SecureString; mask sensitive data in UI.
- Destructive actions require confirm dialog.
- Cost estimates use real `DescribeTaskDefinition` data, never name-based guessing.
- SSM target format: `ecs:{cluster-name}_{task-id}_{runtime-id}`.
- Rollback: explicitly find PRIMARY deployment, verify target task def exists.

## Auto-Sync Rules

After design changes:
1. Update `.kiro/docs/architecture.md` if structure changed.
2. Update this file if new patterns/conventions were established.
3. Update module-level notes if new directories were added.
