# Architecture (Kiro pointer)

Canonical document: [`../../docs/architecture.md`](../../docs/architecture.md)

This file is a short entry point for Kiro agents. Prefer the canonical doc for diagrams, layer tables (Presentation / Data / Action), and design-decision rationale.

## TL;DR

- **Stack**: Go 1.24 + Bubbletea (Elm MVU) + Bubbles + Lipgloss + AWS SDK v2.
- **Entry**: `main.go` parses flags, loads `~/.ecs9s/config.yaml`, starts `tea.Program` with `app.New()`.
- **Shell**: `internal/app/app.go` is the root model; owns tab bar, status bar, command bar, help, and the active page (1 of 21).
- **Pages**: `internal/ui/pages/*` — one model per AWS resource view; route via `PageType` enum in `internal/ui/messages.go`.
- **AWS**: `internal/aws/*.go` — one file per service (ecs, cloudwatch, ecr, elb, ec2, iam, autoscaling, ssm, secrets). Called only from `tea.Cmd` goroutines.
- **Actions**: `internal/action/*.go` — mutations (exec, portforward, scale, deploy, rollback, taskdef). May shell out to `session-manager-plugin`.

## Data Flow

```
User Input → App.Update() → Page.Update() → tea.Cmd(AWS API) → fooLoadedMsg → Page.View()
```

## Key Design Decisions (short)

| Decision | Why |
|---------|-----|
| Bubbletea MVU | Clean state management for 21+ async views |
| Hybrid navigation (Tab + `:cmd`) | Discoverable for new users, fast for power users |
| `NavContext` drill-down | Carries `ClusterARN` / `ServiceName` through Cluster → Service → Task → Container |
| `WithDecryption: false` for SSM SecureString | Secrets never reach process memory |
| Real `DescribeTaskDefinition` for cost | Accurate vs. name-based heuristics |

## When to Update

Any change to the layer structure, page count, navigation model, or action set must be reflected in `../../docs/architecture.md` and summarized here. Keep `AGENTS.md` aligned.
