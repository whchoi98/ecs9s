# Tech

## Language & Toolchain

- **Go 1.24+** (see `go.mod`)
- Build: `go build -o ecs9s .` → single static binary
- Test: `go test ./...`
- Lint: `go vet ./...`

## Libraries

| Purpose | Library |
|---------|---------|
| TUI framework (Elm-style MVU) | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| TUI widgets | [Bubbles](https://github.com/charmbracelet/bubbles) |
| Styling / layout | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| AWS APIs | AWS SDK for Go v2 |

## AWS Services Integrated

ECS, CloudWatch (Logs + Metrics + Alarms), ECR, ELBv2, EC2, IAM, Application Auto Scaling, SSM (Parameter Store + Session Manager), Secrets Manager.

## Runtime Prerequisites

- AWS credentials via `~/.aws/credentials`, env vars, or SSO.
- `session-manager-plugin` installed locally for ECS Exec shell and port forwarding.

## Config

- Path: `~/.ecs9s/config.yaml`
- Loader: `internal/config/config.go` (YAML → struct)
- CLI flag overrides: `--profile`, `--region`, `--theme`.

## Async Pattern

AWS calls wrapped in `tea.Cmd`, results returned as typed messages:

```go
func (p *Page) fetchData() tea.Cmd {
    return func() tea.Msg {
        data, err := client.ListFoo(context.Background())
        return fooLoadedMsg{data: data, err: err}
    }
}
```

## Build & Release

- Single-binary distribution; no runtime dependencies beyond `session-manager-plugin` for exec features.
- Cross-platform builds via `GOOS`/`GOARCH` (see `.claude/commands/deploy.md` legacy notes; replicate as needed).
