# Structure

## Repository Layout

```
ecs9s/
├── main.go                       # CLI entry: flags → config → tea.Program
├── go.mod / go.sum
├── internal/
│   ├── app/
│   │   └── app.go                # Root bubbletea model, routing, page registry
│   ├── ui/
│   │   ├── messages.go           # Shared msg types, PageType enum, NavContext
│   │   ├── components/           # Reusable: table, tabs, commandbar, statusbar,
│   │   │                         # help, confirm, logviewer, sparkline
│   │   ├── pages/                # 21 page models (one per resource view)
│   │   └── styles/               # Lipgloss styles per theme
│   ├── aws/
│   │   ├── session.go            # Profile/region session bootstrap
│   │   ├── ecs.go, ecr.go, elb.go, ec2.go, iam.go,
│   │   ├── autoscaling.go, ssm.go, secrets.go, cloudwatch.go
│   ├── action/
│   │   ├── exec.go               # ECS Exec shell
│   │   ├── portforward.go        # SSM port forwarding
│   │   ├── scale.go, deploy.go, rollback.go, taskdef.go
│   ├── config/                   # YAML loader
│   └── theme/                    # Dark, Light, Blue presets
├── docs/
│   ├── architecture.md           # Canonical design doc
│   ├── runbooks/                 # Ops procedures
│   ├── decisions/                # ADRs
│   └── superpowers/              # Plans and specs
├── scripts/{setup,install-hooks}.sh
├── tests/{hooks,fixtures,structure}/ + run-all.sh
├── tools/prompts/                # Prompt assets
├── screenshots/                  # README images
├── AGENTS.md                     # Kiro entry doc
├── CLAUDE.md                     # Legacy Claude Code doc (keep in sync with AGENTS.md)
├── README.md / CHANGELOG.md / LICENSE
└── .kiro/                        # Kiro steering, rules, docs
    ├── rules.md
    ├── steering/{product,tech,structure,conventions,security}.md
    └── docs/architecture.md
```

## Layering

- **Presentation** (`internal/ui`, `internal/app`, `internal/theme`) — pure view/controller logic; no direct AWS calls.
- **Data** (`internal/aws`) — thin SDK wrappers returning domain structs; one file per AWS service.
- **Action** (`internal/action`) — mutation operations that may invoke external processes (`session-manager-plugin`) or combine multiple AWS calls (e.g. safe rollback).
- **Config / Theme** — pure, stateless.

## Adding a New Page

Edit these locations:

1. `internal/ui/messages.go` — add `PageType` constant and command string.
2. `internal/ui/pages/<name>.go` — implement `Init`, `Update`, `View`.
3. `internal/app/app.go` — 5 edits:
   - Add field to `App` struct.
   - Initialize in `New()`.
   - Handle in `resize()`.
   - Handle in `initCurrentPage()`.
   - Handle in `updateActivePage()` / `viewActivePage()`.
4. If it calls a new AWS service, add `internal/aws/<service>.go`.
5. Update `AGENTS.md` and `.kiro/docs/architecture.md` if the surface changed.

## File Naming

- Lowercase, no underscores in package paths.
- One file per AWS service client (`ecs.go`, `ecr.go`…).
- One file per page model under `internal/ui/pages/`.
