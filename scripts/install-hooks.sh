#!/bin/bash
set -e

HOOKS_DIR=".git/hooks"

if [ ! -d ".git" ]; then
  echo "Not a git repository. Skipping hook installation."
  exit 0
fi

mkdir -p "$HOOKS_DIR"

# Install commit-msg hook: removes Co-Authored-By lines (AI contributor exclusion)
cat > "$HOOKS_DIR/commit-msg" << 'HOOK'
#!/bin/bash
# Remove all Co-Authored-By lines from commit messages
# Prevents Claude and other AI assistants from appearing as contributors
TEMP=$(mktemp)
grep -vi "^Co-Authored-By:" "$1" > "$TEMP" || true
# Remove trailing blank lines
sed -e :a -e '/^\n*$/{$d;N;ba' -e '}' "$TEMP" > "$1"
rm -f "$TEMP"
HOOK

chmod +x "$HOOKS_DIR/commit-msg"
echo "Installed commit-msg hook (AI contributor exclusion)"
