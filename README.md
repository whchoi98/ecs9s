# ecs9s

A terminal UI for managing AWS ECS clusters and related services — inspired by [k9s](https://github.com/derailed/k9s), [e1s](https://github.com/keidarcy/e1s), and [tui-aws](https://github.com/whchoi98/tui-aws).

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- **21 resource views** — ECS (Cluster, Service, Task, Container, TaskDef) + operational services (Logs, ECR, ELB, ASG, VPC, IAM, Metrics, EC2, Events, Stopped Tasks, Resource Map, Cost, SSM, Secrets, Deploy History, Alarms)
- **Hybrid navigation** — Tab-based switching + k9s-style command mode (`:cluster`, `:service`, `:ecr`, etc.)
- **Drill-down** — Cluster → Service → Task → Container with breadcrumb history
- **ECS Exec** — Interactive shell access to containers (`x` key) with prerequisite validation
- **Actions** — Force deploy, scale, rollback, stop task, port forwarding, enable ECS Exec
- **3 themes** — Dark (Tokyo Night), Light, Blue (Navy)
- **Single binary** — `go build` produces one 18MB executable

## Quick Start

```bash
# Build
go build -o ecs9s .

# Run (uses default AWS profile and region)
./ecs9s

# With specific profile and region
./ecs9s --profile myprofile --region ap-northeast-2

# With theme
./ecs9s --theme light
```

### Prerequisites

- Go 1.24+
- AWS credentials configured (`~/.aws/credentials` or environment variables)
- `session-manager-plugin` (for ECS Exec shell access)

## Navigation

| Key | Action |
|-----|--------|
| `Tab` / `[` / `]` | Switch tabs |
| `:` + command | Command mode (`:cluster`, `:ecr`, `:cost`, etc.) |
| `Enter` | Drill down (Cluster → Service → Task → Container) |
| `Esc` / `Backspace` | Go back |
| `/` | Filter current table |
| `s` | Sort columns |
| `R` | Refresh |
| `?` | Help overlay |
| `q` | Quit |

## Page Actions

| Page | Key | Action |
|------|-----|--------|
| Service | `f` | Force new deployment |
| Service | `e` | Enable ECS Exec |
| Service | `S` | Scale desired count |
| Service | `b` | Rollback to previous task definition |
| Task | `Ctrl+d` | Stop task |
| Container | `x` | ECS Exec (interactive shell) |
| Container | `Ctrl+f` | Port forwarding |
| TaskDef | `Ctrl+d` | Deregister task definition |

## ECS Exec Shell Access

1. Select a service → press `e` to enable ECS Exec (sets `enableExecuteCommand` + force deploys)
2. Wait for new tasks to start with ExecuteCommandAgent
3. Drill down to Container → press `x` to open interactive shell

**Required IAM permissions** on Task Role:
```
ssmmessages:CreateControlChannel
ssmmessages:CreateDataChannel
ssmmessages:OpenControlChannel
ssmmessages:OpenDataChannel
```

## Configuration

Config file: `~/.ecs9s/config.yaml`

```yaml
aws:
  profile: default
  region: ap-northeast-2
theme: dark    # dark | light | blue
```

CLI flags override config file values.

## Project Structure

```
ecs9s/
├── main.go                     # Entry point
├── internal/
│   ├── app/                    # Root bubbletea model, routing
│   ├── ui/
│   │   ├── components/         # Table, tabs, commandbar, statusbar, help, confirm, logviewer, sparkline
│   │   ├── pages/              # 21 page models
│   │   └── styles/             # Lipgloss theme styles
│   ├── aws/                    # AWS SDK v2 clients (10 services)
│   ├── action/                 # ECS Exec, port forward, scale, deploy, rollback
│   ├── config/                 # YAML config loader
│   └── theme/                  # Dark, Light, Blue presets
├── tests/                      # TAP structure tests + hook behavior tests
└── .claude/                    # Claude Code hooks, skills, agents, commands
```

## Development

```bash
# Setup
bash scripts/setup.sh

# Test
go test ./... -v
bash tests/run-all.sh
bash tests/hooks/test-secret-scan.sh

# Build release binaries
bash -c 'source .claude/commands/deploy.md'  # or run commands manually
```

## Tech Stack

- [Bubbletea](https://github.com/charmbracelet/bubbletea) — Elm architecture TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (table, textinput, viewport)
- [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2) — AWS API access
