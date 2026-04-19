# .claude/ Extension Guide

This directory contains Claude Code harness components for the ecs9s project.

## Directory Structure

```
.claude/
├── settings.json          # Hook registrations + deny list
├── hooks/                 # Lifecycle hooks (shell scripts)
├── skills/<name>/SKILL.md # Reusable skill definitions
├── commands/<name>.md     # Slash commands (/review, /test-all, /deploy)
├── agents/<name>.md       # Background agents (code-reviewer, security-auditor)
└── README.md              # This file
```

## Adding a New Hook

1. Create a shell script in `.claude/hooks/<name>.sh`
2. Make it executable: `chmod +x .claude/hooks/<name>.sh`
3. Register it in `settings.json` under the appropriate event:

```json
{
  "hooks": {
    "<Event>": [
      {
        "type": "command",
        "command": ".claude/hooks/<name>.sh",
        "toolNames": ["Bash", "Write", "Edit"]
      }
    ]
  }
}
```

**Hook events**: `SessionStart`, `PreToolUse`, `PostToolUse`, `Notification`, `Stop`

**Exit codes**: `0` = pass, `2` = block (PreToolUse only)

**Input**: Hook receives tool input via **stdin** as JSON. Read with `CONTENT="$(cat)"`.

## Adding a New Skill

1. Create `.claude/skills/<name>/SKILL.md` (subdirectory layout)
2. Also create `.claude/skills/<name>.md` (flat copy for tooling compatibility)
3. Include YAML frontmatter:

```yaml
---
name: <skill-name>
description: One-line description of what this skill does
---
```

## Adding a New Agent

1. Create `.claude/agents/<name>.md`
2. Include YAML frontmatter with `name`, `description`, `model`, `allowed-tools`
3. Include an **Output Format** section with a concrete template
4. Agents run in the background — use for large-scope analysis

## Adding a New Command

1. Create `.claude/commands/<name>.md`
2. Include copy-pasteable shell blocks
3. Include an **Error Recovery** section explaining what to do when steps fail

## Component Interaction

```
SessionStart → session-context.sh (loads project info)
PreToolUse   → secret-scan.sh (blocks secrets in Bash/Write/Edit)
PostToolUse  → check-doc-sync.sh (reminds about doc sync on Write/Edit)
/review      → quick inline code review of git diff
/test-all    → go vet + go test + go build
/deploy      → cross-platform binary builds
code-reviewer agent → background parallel review (large changes)
security-auditor agent → full codebase security scan
```
