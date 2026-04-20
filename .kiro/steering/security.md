# Security

## Credentials

- Never log AWS credentials or session tokens.
- Rely on the SDK's default credential chain: env vars → shared config → SSO → IMDS.
- Respect `--profile` / `--region` CLI flags; do not override them silently.

## Sensitive Data Handling

- SSM Parameter Store `SecureString` values must be fetched with `WithDecryption: false`. Display masked placeholders only.
- Secrets Manager values must never be printed to stdout, logs, or screenshots.
- Mask anything resembling a key, token, or password in UI tables and log viewers.

## Destructive / Blast-Radius Operations

Require explicit confirm dialog:

- Force new deployment
- Service scale (desired count change)
- Stop task
- Rollback to previous task definition
- Any future `DeleteService` / `DeregisterTaskDefinition` surface

Confirm dialog must display: resource identifier, action, and irreversibility note where applicable.

## ECS Exec / Port Forward

- Validate prerequisites before invoking: task has `enableExecuteCommand`, container is running, `session-manager-plugin` present on PATH.
- Build SSM target as `ecs:{cluster-name}_{task-id}_{runtime-id}`. Never pass ARNs directly.
- Do not persist session IDs or temporary credentials.

## Supply Chain

- Pin Go module versions in `go.mod` / `go.sum` (tooling default).
- Review transitive updates before `go mod tidy` on release branches.

## Secret Scanning

The legacy Claude Code harness ran a secret-scan hook on `Bash`/`Write`/`Edit` (see `.claude/hooks/secret-scan.sh`). Equivalent local scanning should remain part of the pre-commit or CI flow — do not remove it when migrating off Claude Code.
