#!/bin/bash
# PostToolUse hook: detect when source files change and remind about doc sync
# Triggers on Write/Edit tool calls to internal/ directory

CHANGED_FILE="${TOOL_INPUT_FILE_PATH:-}"

if [[ -z "$CHANGED_FILE" ]]; then
  exit 0
fi

# Walk up to find if this file is under a directory that has CLAUDE.md
check_dir=$(dirname "$CHANGED_FILE")
while [[ "$check_dir" != "." && "$check_dir" != "/" ]]; do
  if [[ -f "$check_dir/CLAUDE.md" ]]; then
    exit 0
  fi
  check_dir=$(dirname "$check_dir")
done

# Check if it's a source file change
if [[ "$CHANGED_FILE" == internal/* || "$CHANGED_FILE" == *.go ]]; then
  if [[ "$CHANGED_FILE" == *_test.go ]]; then
    exit 0
  fi
  echo "NOTICE: Source file changed ($CHANGED_FILE). Consider running /sync-docs if architecture changed."
fi
