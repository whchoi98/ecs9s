# pages

21 page models — one per AWS resource view. Each implements the bubbletea pattern:

- `Init()` / `Update(msg)` / `View()` / `SetSize(w, h)` / `Refresh() tea.Cmd`
- Context-dependent pages also have `SetContext(nav NavContext) tea.Cmd`
- Each page has a typed `xxxLoadedMsg` for async data results
- `HelpBindings()` returns page-specific keybindings for the help overlay

## Pages

| File | Resource | Context Required | Actions |
|------|----------|-----------------|---------|
| cluster.go | ECS Clusters | - | Enter → Service |
| service.go | ECS Services | ClusterARN | Enter → Task, f (deploy), e (enable exec), S (scale), b (rollback) |
| task.go | ECS Tasks | ClusterARN + ServiceName | Enter → Container, Ctrl+d (stop) |
| container.go | Containers | ClusterARN + TaskARN | x (ECS Exec shell) |
| taskdef.go | Task Definitions | - | Ctrl+d (deregister) |
| logs.go | CloudWatch Logs | log group + stream | F (follow), c (clear) |
| ecr.go | ECR Repositories | - | read-only |
| elb.go | Load Balancers | - | read-only |
| autoscaling.go | Auto Scaling | - | read-only |
| vpc.go | VPCs | - | read-only |
| iam.go | IAM Roles | - | read-only |
| metrics.go | CloudWatch Metrics | ClusterName + ServiceName | 1/3/6 (time range) |
| ec2.go | EC2 Instances | - | read-only |
| events.go | Service Events | ClusterARN + ServiceName | read-only |
| stopped.go | Stopped Tasks | ClusterARN | read-only |
| resmap.go | Resource Map | ClusterARN + ServiceName | read-only |
| cost.go | Cost Estimate | ClusterARN | read-only (SetSize reserves 6 lines for header/footer) |
| ssm.go | SSM Parameters | - | read-only |
| secrets.go | Secrets Manager | - | read-only |
| deploy.go | Deploy History | ClusterARN + ServiceName | read-only |
| alarms.go | CloudWatch Alarms | - | read-only |

## Convention

Use `findXxxByName/ID()` with `table.SelectedRow()` for row selection. Never index into the data array with `table.Cursor()`.
