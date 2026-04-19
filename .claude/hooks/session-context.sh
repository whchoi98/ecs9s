#!/bin/bash
# SessionStart hook: load project context at session start

echo "Project: ecs9s (Go TUI for AWS ECS)"
echo "Build: go build -o ecs9s ."
echo "Test: go test ./..."
echo "Run: ./ecs9s [--profile X] [--region Y] [--theme Z]"

# Show current git state
if git rev-parse --git-dir > /dev/null 2>&1; then
  BRANCH=$(git branch --show-current 2>/dev/null)
  CHANGES=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
  echo "Git: branch=$BRANCH changes=$CHANGES"
fi
