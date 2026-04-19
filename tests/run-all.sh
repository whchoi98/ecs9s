#!/bin/bash
set -e

PASS=0
FAIL=0
TOTAL=0

assert_eq() {
  TOTAL=$((TOTAL + 1))
  if [ "$1" = "$2" ]; then
    echo "ok $TOTAL - $3"
    PASS=$((PASS + 1))
  else
    echo "not ok $TOTAL - $3 (expected '$1', got '$2')"
    FAIL=$((FAIL + 1))
  fi
}

assert_file_exists() {
  TOTAL=$((TOTAL + 1))
  if [ -f "$1" ]; then
    echo "ok $TOTAL - $2"
    PASS=$((PASS + 1))
  else
    echo "not ok $TOTAL - $2 (file not found: $1)"
    FAIL=$((FAIL + 1))
  fi
}

assert_executable() {
  TOTAL=$((TOTAL + 1))
  if [ -x "$1" ]; then
    echo "ok $TOTAL - $2"
    PASS=$((PASS + 1))
  else
    echo "not ok $TOTAL - $2 (not executable: $1)"
    FAIL=$((FAIL + 1))
  fi
}

echo "TAP version 13"
echo "# ecs9s project structure tests"

# Structure tests
assert_file_exists "CLAUDE.md" "CLAUDE.md exists"
assert_file_exists "go.mod" "go.mod exists"
assert_file_exists "main.go" "main.go exists"
assert_file_exists ".gitignore" ".gitignore exists"
assert_file_exists ".claude/settings.json" "settings.json exists"
assert_file_exists "docs/architecture.md" "architecture.md exists"

# Hook tests
assert_file_exists ".claude/hooks/secret-scan.sh" "secret-scan hook exists"
assert_file_exists ".claude/hooks/check-doc-sync.sh" "doc-sync hook exists"
assert_file_exists ".claude/hooks/session-context.sh" "session-context hook exists"
assert_executable ".claude/hooks/secret-scan.sh" "secret-scan is executable"
assert_executable ".claude/hooks/check-doc-sync.sh" "doc-sync is executable"

# Module CLAUDE.md tests
for dir in internal/app internal/aws internal/ui internal/action internal/config internal/theme; do
  assert_file_exists "$dir/CLAUDE.md" "$dir/CLAUDE.md exists"
done

# Build test
TOTAL=$((TOTAL + 1))
if go build -o /dev/null . 2>/dev/null; then
  echo "ok $TOTAL - go build succeeds"
  PASS=$((PASS + 1))
else
  echo "not ok $TOTAL - go build fails"
  FAIL=$((FAIL + 1))
fi

# Vet test
TOTAL=$((TOTAL + 1))
if go vet ./... 2>/dev/null; then
  echo "ok $TOTAL - go vet passes"
  PASS=$((PASS + 1))
else
  echo "not ok $TOTAL - go vet fails"
  FAIL=$((FAIL + 1))
fi

echo ""
echo "1..$TOTAL"
echo "# pass $PASS / $TOTAL"
if [ $FAIL -gt 0 ]; then
  echo "# FAIL $FAIL test(s)"
  exit 1
fi
echo "# All tests passed"
