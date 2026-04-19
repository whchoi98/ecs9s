# ecs9s

ECS-focused Terminal UI for AWS — manage clusters, services, tasks, containers, and related operational services with a k9s-inspired keyboard-driven interface.

## Tech Stack

- **Language**: Go 1.24+
- **TUI Framework**: Bubbletea (Elm architecture) + Bubbles + Lipgloss
- **AWS SDK**: AWS SDK for Go v2
- **Config**: YAML (`~/.ecs9s/config.yaml`)
- **Build**: `go build` → single static binary

## Project Structure

```
ecs9s/
├── main.go                    # Entry point (CLI flags, config load, tea.Program)
├── internal/
│   ├── app/app.go             # Root bubbletea model, routing, page management
│   ├── ui/
│   │   ├── messages.go        # Shared message types, PageType enum, NavContext
│   │   ├── components/        # Reusable UI: table, tabs, commandbar, statusbar, help, confirm, logviewer, sparkline
│   │   ├── pages/             # 21 page models (one per AWS resource view)
│   │   └── styles/            # Lipgloss style definitions per theme
│   ├── aws/                   # AWS SDK v2 client wrappers (one file per service)
│   ├── action/                # Mutation operations (exec, portforward, scale, deploy, rollback)
│   ├── config/                # Config loading from ~/.ecs9s/config.yaml
│   └── theme/                 # Preset themes (dark, light, blue)
├── docs/                      # Design specs, architecture docs
├── tests/                     # Test suites
└── scripts/                   # Setup and utility scripts
```

## Key Commands

```bash
# Build
go build -o ecs9s .

# Run
./ecs9s                           # default profile/region
./ecs9s --profile myprofile       # specific AWS profile
./ecs9s --region us-east-1        # specific region
./ecs9s --theme light             # theme override

# Test
go test ./...

# Lint
go vet ./...
```

## Architecture Patterns

### Bubbletea Model-Update-View

Every page implements the Bubbletea pattern:
- `Init()` → initial tea.Cmd (typically fetches data)
- `Update(msg)` → handles messages, returns new model + cmd
- `View()` → renders the current state as a string

### Async Data Loading

AWS API calls use `tea.Cmd` for non-blocking I/O:
```go
func (p *Page) fetchData() tea.Cmd {
    return func() tea.Msg {
        data, err := client.ListFoo(context.Background())
        return fooLoadedMsg{data: data, err: err}
    }
}
```

### Navigation

- **Tab switching**: `Tab`/`Shift+Tab` or `[`/`]`
- **Command mode**: `:cluster`, `:service`, `:ecr`, etc.
- **Drill-down**: `Enter` (Cluster → Service → Task → Container)
- **Go back**: `Esc`/`Backspace` (uses navStack)
- **NavContext**: carries ClusterARN/ServiceName through drill-down chain

### Page Registration

Adding a new page requires changes in 5 places in `internal/app/app.go`:
1. Add field to `App` struct
2. Initialize in `New()` constructor
3. Add to `resize()`
4. Add to `initCurrentPage()` / `updateActivePage()` / `viewActivePage()`

Also add the PageType constant and command string in `internal/ui/messages.go`.

## Conventions

- **File naming**: lowercase, one file per AWS service client or page model
- **Error handling**: AWS errors surface via `ui.ErrorMsg` to status bar
- **Security**: Never decrypt SecureString values; mask sensitive data in UI
- **Destructive actions**: Require confirmation dialog (confirm component)
- **Cost estimates**: Use actual DescribeTaskDefinition data, never guess from names
- **SSM target format**: `ecs:{cluster-name}_{task-id}_{runtime-id}` (short names, not ARNs)
- **Rollback safety**: Explicitly find PRIMARY deployment, verify target task def exists

## Auto-Sync Rules

When exiting Plan mode after creating or updating a design document:
1. Update `docs/architecture.md` if project structure changed
2. Update this CLAUDE.md if new patterns or conventions were established
3. Update module CLAUDE.md files if new directories were added
