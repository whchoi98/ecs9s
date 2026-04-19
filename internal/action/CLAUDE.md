# action

Mutation operations that interact with AWS or spawn subprocesses.

- exec.go: ECS Exec via `aws ecs execute-command`
- portforward.go: SSM port forward. Target format: `ecs:{cluster-name}_{task-id}_{runtime-id}` (short names only)
- scale.go, deploy.go: Service scaling and force deployment
- rollback.go: Safe rollback — finds PRIMARY deployment, selects verified previous task def
