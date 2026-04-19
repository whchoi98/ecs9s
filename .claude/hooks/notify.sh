#!/bin/bash
# Notification hook: handle notifications from Claude Code
# Customize this to send to Slack, email, etc.

NOTIFICATION="${1:-}"
echo "[ecs9s] $NOTIFICATION"
