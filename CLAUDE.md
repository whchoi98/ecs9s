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
│   │   ├── messages.go        # Shared message types, PageType enum (21 types), NavContext
│   │   ├── components/        # 8 reusable widgets: table, tabs, commandbar, statusbar, help, confirm, logviewer, sparkline
│   │   ├── pages/             # 21 page models (one per AWS resource view)
│   │   └── styles/            # Lipgloss style definitions per theme
│   ├── aws/                   # 10 AWS SDK v2 client wrappers: ecs, cloudwatch, ecr, elb, ec2, iam, autoscaling, ssm, secrets, session
│   ├── action/                # Mutation operations: exec (tea.ExecProcess), portforward, scale, deploy, rollback, taskdef
│   ├── config/                # Config loading from ~/.ecs9s/config.yaml
│   └── theme/                 # Preset themes (dark, light, blue)
├── .claude/                   # Claude Code hooks, skills, agents, commands
├── .kiro/                     # Kiro CLI steering docs, rules, architecture
├── .github/workflows/         # CI pipeline (vet, test, build)
├── docs/                      # Design specs, architecture docs
├── tests/                     # Go unit tests + TAP structure tests + hook behavior tests
├── scripts/                   # setup.sh, install-hooks.sh
└── screenshots/               # TUI screenshots for README
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
go test ./...                     # Go unit tests
go vet ./...                      # Lint
bash tests/run-all.sh             # TAP structure tests (19 checks)
bash tests/hooks/test-secret-scan.sh  # Hook behavior tests (12 checks)
```

## Architecture Patterns

### Bubbletea Model-Update-View

Every page implements: `Init()` → `Update(msg)` → `View()`. The root App model in `internal/app/app.go` composes all pages and components.

### Async Data Loading

AWS API calls return `tea.Cmd` closures. Results arrive as typed messages (e.g., `clusterLoadedMsg`).

### Navigation

- **Tab switching**: `Tab`/`Shift+Tab` or `[`/`]` — preserves drill-down context
- **Command mode**: `:cluster`, `:service`, `:ecr`, etc.
- **Drill-down**: `Enter` (Cluster → Service → Task → Container) — only way to set NavContext
- **Go back**: `Esc`/`Backspace` (restores from navStack)
- **Global pages** (ECR, ELB, VPC, IAM, etc.): status bar hides cluster display
- **Context pages** (Service, Task, Events, Metrics, etc.): show data scoped to drill-down context

### Page Registration (5-step checklist)

Adding a new page requires changes in `internal/app/app.go`:
1. Add field to `App` struct
2. Initialize in `New()` constructor
3. Add to `resize()`
4. Add case to `initCurrentPage()`, `updateActivePage()`, `viewActivePage()`

Also add PageType constant + command string in `internal/ui/messages.go`, and tab entry in `New()`.

### Row Selection (SelectedRow, not Cursor)

Always use `p.table.SelectedRow()` to get the selected row data, then find the original data with a `findXxxByName/ID()` helper. Never use `p.table.Cursor()` as an index into the data array — it breaks when the table is filtered or sorted.

### ECS Exec Flow

Container page `x` key → `ExecRequestMsg` → App checks prerequisites + `CheckExecEnabled()` → `execReadyMsg` → `tea.ExecProcess(cmd)` suspends TUI → shell runs → TUI restores.

## Conventions

- **File naming**: lowercase, one file per AWS service client or page model
- **Error handling**: AWS errors surface via `ui.ErrorMsg` to status bar
- **Security**: Never decrypt SecureString values (`WithDecryption: false`); mask with `"****"`
- **Destructive actions**: Require confirmation dialog (confirm component)
- **Cost estimates**: Use actual `DescribeTaskDefinition` data, never guess from names
- **SSM target format**: `ecs:{cluster-name}_{task-id}_{runtime-id}` (short names, not ARNs)
- **Rollback safety**: Explicitly find PRIMARY deployment, verify target task def exists
- **Empty states**: All pages must show a helpful message when context is missing or results are empty
- **Table height**: Pages with headers/footers (Cost, Metrics) must subtract those lines from table height in `SetSize()`

## Auto-Sync Rules

When exiting Plan mode after creating or updating a design document:
1. Update `docs/architecture.md` if project structure changed
2. Update this CLAUDE.md if new patterns or conventions were established
3. Update module CLAUDE.md files if new directories were added
