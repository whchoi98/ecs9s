#!/bin/bash
# Behavioral tests for secret-scan.sh hook
# Uses TAP format

HOOK=".claude/hooks/secret-scan.sh"
PASS=0
FAIL=0
TOTAL=0

run_test() {
  TOTAL=$((TOTAL + 1))
  local desc="$1"
  local input="$2"
  local expect_exit="$3"

  actual_exit=0
  echo "$input" | bash "$HOOK" > /dev/null 2>&1 || actual_exit=$?

  if [ "$actual_exit" -eq "$expect_exit" ]; then
    echo "ok $TOTAL - $desc"
    PASS=$((PASS + 1))
  else
    echo "not ok $TOTAL - $desc (expected exit $expect_exit, got $actual_exit)"
    FAIL=$((FAIL + 1))
  fi
}

echo "TAP version 13"
echo "# secret-scan.sh behavioral tests"

# Positive tests (should block with exit 2)
run_test "blocks AWS access key" '{"input": "AKIAIOSFODNN7EXAMPLE"}' 2
run_test "blocks GitHub PAT" '{"input": "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghij"}' 2
run_test "blocks private key header" '{"input": "BEGIN RSA PRIVATE KEY"}' 2
run_test "blocks password assignment" '{"input": "password = \"supersecret123\""}' 2
run_test "blocks Slack token" '{"input": "xoxb-1234567890-abcdefghij"}' 2
run_test "blocks AWS session token" '{"input": "IQoJb3JpZ2lu"}' 2
run_test "blocks mongodb URI" '{"input": "mongodb+srv://user:pass@cluster.mongodb.net"}' 2
run_test "blocks postgresql URI" '{"input": "postgresql://user:pass@localhost:5432/db"}' 2

# Negative tests (should pass with exit 0)
run_test "allows normal Go code" '{"input": "func main() { fmt.Println(\"hello\") }"}' 0
run_test "allows normal text" '{"input": "This is a normal commit message"}' 0
run_test "allows empty input" '' 0
run_test "allows AWS service names" '{"input": "aws ecs list-clusters"}' 0

echo ""
echo "1..$TOTAL"
echo "# pass $PASS / $TOTAL"
if [ $FAIL -gt 0 ]; then
  echo "# FAIL $FAIL test(s)"
  exit 1
fi
echo "# All tests passed"
