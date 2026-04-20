# Product

ecs9s is a keyboard-driven terminal UI for AWS ECS operations. It unifies management of clusters, services, tasks, containers, and related AWS services (CloudWatch, ECR, ELB, ASG, VPC, IAM, SSM, Secrets Manager) into a single binary.

## Target Users

- Platform / DevOps engineers who operate ECS daily.
- SREs investigating deployments, rollbacks, scaling events, and logs.
- Developers needing quick `ecs exec` shell access or port forwarding without context-switching to the AWS console.

## Inspirations

- [k9s](https://github.com/derailed/k9s) — Kubernetes TUI
- [e1s](https://github.com/keidarcy/e1s) — ECS TUI
- [tui-aws](https://github.com/whchoi98/tui-aws) — multi-service AWS TUI

## Feature Surface (21 views)

- **ECS core**: Cluster, Service, Task, Container, TaskDef
- **Operational**: Logs, ECR, ELB, ASG, VPC, IAM, Metrics, EC2, Events, Stopped Tasks, Resource Map, Cost, SSM, Secrets, Deploy History, Alarms

## UX Principles

- Keyboard-first. Mouse is optional.
- Hybrid navigation: Tab-switching for discoverability, command mode (`:cluster`, `:ecr`…) for power users.
- Drill-down with breadcrumb history (`NavContext`): Cluster → Service → Task → Container.
- Destructive actions always require explicit confirm.
- Sensitive data (SecureString, secrets) is masked; never decrypted client-side.
