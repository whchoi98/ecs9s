# ecs9s Design Specification

**Date:** 2026-04-16
**Status:** Approved

## Overview

ecs9s is a Terminal User Interface (TUI) application for managing AWS ECS clusters and related services. Inspired by [e1s](https://github.com/keidarcy/e1s) (ECS-focused TUI), [tui-aws](https://github.com/whchoi98/tui-aws) (broad AWS TUI), and [k9s](https://github.com/derailed/k9s) (Kubernetes TUI), it combines ECS-centric functionality with operational support services in a modern, keyboard-driven interface.

## Scope

**ECS-focused with related services** — ECS is the primary domain, supplemented by services essential for ECS operations: CloudWatch, ECR, ALB/NLB, Auto Scaling, VPC/Subnet/SG, IAM, EC2.

## Architecture

### Layer Diagram

```
┌─────────────────────────────────────────────────┐
│                   main.go                        │
│              (CLI flags, config load)            │
├─────────────────────────────────────────────────┤
│                   internal/                      │
│  ┌───────────┐  ┌───────────┐  ┌─────────────┐ │
│  │    ui/     │  │    aws/   │  │   config/   │ │
│  │ bubbletea  │  │  service  │  │  theme,     │ │
│  │ models &   │←→│  clients  │  │  profile,   │ │
│  │ components │  │  (SDK v2) │  │  keymap     │ │
│  └───────────┘  └───────────┘  └─────────────┘ │
│  ┌───────────┐  ┌───────────┐                   │
│  │  action/   │  │   theme/  │                   │
│  │ exec, port │  │ dark/light│                   │
│  │ fwd, stop  │  │ presets   │                   │
│  └───────────┘  └───────────┘                   │
├─────────────────────────────────────────────────┤
│              AWS SDK for Go v2                   │
└─────────────────────────────────────────────────┘
```

### Layers

- **`ui/`** — Bubbletea Model-Update-View pattern. Per-resource page models (cluster, service, task, etc.) and shared components (table, command bar, log viewer, tabs, status bar).
- **`aws/`** — AWS SDK v2 client wrappers. One file per service (ECS, CloudWatch, ECR, ELB, EC2, IAM, Auto Scaling). Shared session/profile/region management.
- **`config/`** — User configuration loaded from `~/.ecs9s/config.yaml`. Stores profile, region, theme preference, and keybindings.
- **`action/`** — Mutation operations: ECS Exec, port forwarding, service scale/deploy/rollback, task stop, task definition register.
- **`theme/`** — Lipgloss-based preset themes (Dark, Light, Blue).

## UI Layout

```
┌─ ecs9s ──────────────────────────────────────────────┐
│ [Cluster] [Service] [Task] [Container] [TaskDef] ... │  ← Tab Bar
│ Context: my-cluster | Region: ap-northeast-2         │  ← Status Bar
├──────────────────────────────────────────────────────┤
│ NAME          STATUS    TASKS   CPU    MEM    AGE    │
│ ▸ web-api     ACTIVE    3/3     45%    62%    5d     │  ← Resource Table
│   worker      ACTIVE    2/2     30%    41%    5d     │
│   scheduler   ACTIVE    1/1     12%    25%    3d     │
├──────────────────────────────────────────────────────┤
│ :service  /filter  ?help  q:quit                     │  ← Command Bar
└──────────────────────────────────────────────────────┘
```

### Navigation

| Input | Action |
|-------|--------|
| `Tab` / `Shift+Tab` or `[` / `]` | Switch between tabs |
| `:` + resource name | Command mode — jump to any resource view (`:cluster`, `:service`, `:task`, `:log`, `:ecr`, etc.) |
| `Enter` | Drill down — Cluster → Service → Task → Container |
| `Esc` / `Backspace` | Go back to parent resource |
| `/` | Filter current table |
| `s` | Cycle column sort |
| `?` | Show keybinding help overlay |
| `p` | Switch AWS profile |
| `r` | Switch AWS region |
| `R` | Refresh current view |
| `q` | Quit |
| `j` / `k` or `↑` / `↓` | Row navigation |

## Resource Views & Actions

### ECS Core

| Resource | View Columns | Actions |
|----------|-------------|---------|
| **Cluster** | Name, Status, Services count, Tasks count, Instances count | Select → drill down to Services |
| **Service** | Name, Status, Desired/Running, Load Balancer, CPU%, Mem% | Force Deploy (`f`), Scale (`S`), Rollback (`b`) |
| **Task** | Task ID, Status, Started At, Container count, Private IP | Stop Task (`Ctrl+d`) |
| **Container** | Name, Image, Status, Ports, Health | ECS Exec (`x`), Port Forward (`Ctrl+f`) |
| **Task Definition** | Family, Revision, CPU/Mem, Compatibility | Register New (`n`), Deregister (`Ctrl+d`) |

### Operational Services

| Resource | View Columns | Actions |
|----------|-------------|---------|
| **CloudWatch Logs** | Real-time log streaming (tail -f style) | Filter text, Adjust time range |
| **ECR** | Repository, Tag, Size, Pushed At | Read-only |
| **ALB/NLB** | LB Name, DNS, Target Groups, Healthy/Unhealthy count | Read-only |
| **Auto Scaling** | Policy Name, Min/Max/Desired, Metric | Adjust scale (`S`) |
| **VPC/Subnet/SG** | Network configuration linked to ECS tasks | Read-only |
| **IAM Roles** | Task Role, Execution Role, Attached Policies | Read-only |
| **CloudWatch Metrics** | CPU/Memory sparkline charts per service/task | Adjust time range |
| **EC2 Instances** | Container Instance ID, AMI, Status, Instance Type | Read-only |

## Technology Stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.22+ |
| TUI Framework | Bubbletea + Bubbles + Lipgloss |
| AWS SDK | AWS SDK for Go v2 |
| Config File | `~/.ecs9s/config.yaml` (YAML) |
| Themes | Dark (default), Light, Blue presets |
| Build Output | Single static binary |

## Project Structure

```
ecs9s/
├── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── app/
│   │   └── app.go              # Root bubbletea model, initialization
│   ├── ui/
│   │   ├── components/
│   │   │   ├── table.go        # Sortable, filterable resource table
│   │   │   ├── tabs.go         # Tab bar component
│   │   │   ├── commandbar.go   # Command mode input (: prefix)
│   │   │   ├── statusbar.go    # Context/region/info display
│   │   │   ├── logviewer.go    # Scrollable log streaming view
│   │   │   ├── help.go         # Keybinding help overlay
│   │   │   ├── confirm.go      # Confirmation dialog for destructive actions
│   │   │   └── sparkline.go    # Inline metric sparkline chart
│   │   ├── pages/
│   │   │   ├── cluster.go
│   │   │   ├── service.go
│   │   │   ├── task.go
│   │   │   ├── container.go
│   │   │   ├── taskdef.go
│   │   │   ├── logs.go
│   │   │   ├── ecr.go
│   │   │   ├── elb.go
│   │   │   ├── autoscaling.go
│   │   │   ├── vpc.go
│   │   │   ├── iam.go
│   │   │   ├── metrics.go
│   │   │   └── ec2.go
│   │   └── styles/
│   │       └── styles.go       # Lipgloss style constants per theme
│   ├── aws/
│   │   ├── session.go          # AWS session, profile, region management
│   │   ├── ecs.go
│   │   ├── cloudwatch.go
│   │   ├── ecr.go
│   │   ├── elb.go
│   │   ├── ec2.go
│   │   ├── iam.go
│   │   └── autoscaling.go
│   ├── action/
│   │   ├── exec.go             # ECS Exec (interactive shell)
│   │   ├── portforward.go      # SSM port forwarding
│   │   ├── scale.go            # Service desired count adjustment
│   │   ├── deploy.go           # Force new deployment
│   │   ├── rollback.go         # Service rollback to previous task def
│   │   └── taskdef.go          # Register/deregister task definitions
│   ├── config/
│   │   └── config.go           # ~/.ecs9s/config.yaml parsing
│   └── theme/
│       ├── theme.go            # Theme interface and loader
│       ├── dark.go
│       ├── light.go
│       └── blue.go
├── docs/
│   └── superpowers/
│       └── specs/
│           └── 2026-04-16-ecs9s-design.md
└── README.md
```

## Configuration

Config file at `~/.ecs9s/config.yaml`:

```yaml
aws:
  profile: default
  region: ap-northeast-2

theme: dark        # dark | light | blue

keybindings:
  drill_down: enter
  go_back: esc
  filter: /
  command: ":"
  quit: q
  help: "?"
  refresh: R
  profile: p
  region: r
  sort: s
```

## Data Flow

1. **Startup**: Load config → Initialize AWS session → Fetch cluster list → Render cluster page.
2. **Navigation**: User selects cluster → Fetch services for that cluster → Render service page. Drill-down continues through Task → Container.
3. **Command Mode**: User types `:ecr` → Switch to ECR page → Fetch ECR repositories → Render.
4. **Actions**: User presses `x` on a container → Confirm dialog → Spawn ECS Exec subprocess (replaces TUI temporarily) → Return to TUI on exit.
5. **Refresh**: Manual (`R`) or auto-refresh on a configurable interval for the active view.

## Error Handling

- AWS API errors display in the status bar with the error message.
- Network timeouts show a retry prompt.
- Missing IAM permissions display a clear message indicating which permission is needed.
- Destructive actions (stop task, deregister task def) require confirmation via dialog.

## Non-Goals (Out of Scope for MVP)

- YAML-based custom skin system (future enhancement)
- Plugin/extension system
- Multi-cluster simultaneous view
- Cost estimation
- CloudFormation/CDK integration
