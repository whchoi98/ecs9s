#!/bin/bash
set -e

echo "Setting up ecs9s development environment..."

# Check Go
if ! command -v go &> /dev/null; then
  echo "Error: Go is not installed. Install Go 1.22+ from https://go.dev/dl/"
  exit 1
fi

echo "Go version: $(go version)"

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Build
echo "Building ecs9s..."
go build -o ecs9s .

# Verify
echo "Build successful: ./ecs9s ($(ls -lh ecs9s | awk '{print $5}'))"

# Install hooks
if [ -f scripts/install-hooks.sh ]; then
  bash scripts/install-hooks.sh
fi

echo "Setup complete! Run ./ecs9s to start."
