#!/bin/bash
set -euo pipefail

echo "=== Claude Code Template Test ==="

echo "Checking Docker..."
docker --version

echo "Checking Node.js..."
node --version

echo "Checking Python..."
python3 --version

echo "Checking Claude Code CLI..."
claude --version

echo "Checking mcp-gateway..."
mcp-gateway --help > /dev/null

echo "Checking uv..."
uv --version

echo "Checking poetry..."
poetry --version

echo "Checking jq..."
jq --version

echo "Checking git..."
git --version

echo ""
echo "All checks passed."
