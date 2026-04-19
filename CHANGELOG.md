# Changelog

## [Unreleased]

### Added
- Initial TUI implementation with 21 page views
- ECS core: Cluster, Service, Task, Container, TaskDef pages
- Operational: Logs, ECR, ELB, ASG, VPC, IAM, Metrics, EC2 pages
- Management: Events, Stopped Tasks, Resource Map, Cost, SSM, Secrets, Deploy History, Alarms pages
- Actions: ECS Exec, Port Forwarding, Scale, Force Deploy, Rollback, TaskDef management
- Hybrid navigation: Tab-based + k9s-style command mode (`:cluster`, `:service`, etc.)
- 3 preset themes: Dark (Tokyo Night), Light, Blue (Navy)
- Config: `~/.ecs9s/config.yaml` with profile/region/theme settings
- Claude Code project structure: CLAUDE.md hierarchy, hooks, skills, agents, commands
- CI/CD: GitHub Actions workflow (vet, test, build, structure tests)
- Security: Secret scanning hook (11 patterns), deny list (20 entries)
- Tests: Go unit tests (config, action, theme, ui/messages) + TAP structure tests + hook behavior tests
