# config

Configuration loading from `~/.ecs9s/config.yaml`. Provides defaults: profile=default, region=ap-northeast-2, theme=dark.

Files: `config.go` (Config struct, Load, Save, DefaultConfig), `config_test.go` (save/load roundtrip, missing file fallback).
