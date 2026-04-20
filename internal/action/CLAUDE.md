# action

Mutation operations that interact with AWS or spawn subprocesses.

| File | Purpose |
|------|---------|
| exec.go | `ExecCommand()` returns `*exec.Cmd` for ECS Exec. Called via `tea.ExecProcess` (not cmd.Run). `CheckPrerequisites()` verifies aws CLI + session-manager-plugin. |
| portforward.go | SSM port forward. Target format: `ecs:{cluster-name}_{task-id}_{runtime-id}` (short names only, not ARNs). Has `ExtractClusterName()` and `ExtractTaskID()` helpers. |
| scale.go | Service desired count update. |
| deploy.go | Force new deployment. |
| rollback.go | Safe rollback — explicitly finds PRIMARY deployment, falls back to previous revision with verification. |
| taskdef.go | Task definition deregister. |

Tests: `portforward_test.go` (ExtractClusterName/TaskID), `rollback_test.go` (previousRevision edge cases).
