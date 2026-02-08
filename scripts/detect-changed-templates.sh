#!/usr/bin/env bash
set -euo pipefail

# Detect which template directories have changed compared to the base branch.
# Outputs a JSON array of template names for use in a GitHub Actions matrix.
#
# Usage:
#   ./scripts/detect-changed-templates.sh [base-ref]
#   base-ref defaults to origin/main

BASE_REF="${1:-origin/main}"

# Get list of changed files under templates/
changed_files=$(git diff --name-only "$BASE_REF"...HEAD -- templates/ 2>/dev/null || true)

if [ -z "$changed_files" ]; then
  echo "[]"
  exit 0
fi

# Extract unique template directory names (first path component after templates/)
templates=$(echo "$changed_files" \
  | sed -n 's|^templates/\([^/]*\)/.*|\1|p' \
  | sort -u)

if [ -z "$templates" ]; then
  echo "[]"
  exit 0
fi

# Build JSON array
json="["
first=true
while IFS= read -r name; do
  if [ "$first" = true ]; then
    first=false
  else
    json+=","
  fi
  json+="\"$name\""
done <<< "$templates"
json+="]"

echo "$json"
