#!/bin/bash
# PreToolUse hook: scan for secrets before tool execution
# Reads tool input from stdin (Claude Code hook JSON protocol)
# Exits 2 to block tool execution if a secret pattern is detected

# Read all stdin into a variable (Claude Code passes hook input via stdin)
CONTENT="$(cat)"

if [[ -z "$CONTENT" ]]; then
  exit 0
fi

PATTERNS=(
  'AKIA[0-9A-Z]{16}'
  'aws_secret_access_key\s*=\s*[A-Za-z0-9/+=]{40}'
  'password\s*[:=]\s*.{8,}'
  'ghp_[A-Za-z0-9]{36}'
  'sk-[A-Za-z0-9]{48}'
  'BEGIN\s+(RSA|DSA|EC|OPENSSH)\s+PRIVATE\s+KEY'
  'xox[bpors]-[A-Za-z0-9\-]{10,}'
  'AIza[0-9A-Za-z\-_]{35}'
  'IQoJb3JpZ2lu'
  'mongodb(\+srv)?://[^\s]+'
  'postgresql://[^\s]+'
)

for pattern in "${PATTERNS[@]}"; do
  if echo "$CONTENT" | grep -qEi "$pattern"; then
    echo "BLOCKED: Potential secret detected matching pattern: $pattern"
    echo "Please remove the secret before proceeding."
    exit 2
  fi
done

exit 0
