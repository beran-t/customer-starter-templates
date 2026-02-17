#!/usr/bin/env bash
set -euo pipefail

# Run Python and TypeScript example tests for a single template.
#
# Usage:
#   ./scripts/run-tests.sh <template-name> [tag]
#
# If tag is provided, it is passed via E2B_BUILD_TAG (for build) and
# E2B_TEMPLATE_TAG (for examples).

TEMPLATE_NAME="${1:?Usage: run-tests.sh <template-name> [tag]}"
TAG="${2:-}"
TEMPLATE_DIR="templates/$TEMPLATE_NAME"

if [ ! -d "$TEMPLATE_DIR" ]; then
  echo "Error: Template directory '$TEMPLATE_DIR' does not exist."
  exit 1
fi

if [ -n "$TAG" ]; then
  export E2B_BUILD_TAG="$TAG"
  export E2B_TEMPLATE_TAG="$TAG"
fi

if [ -f "$TEMPLATE_DIR/e2b.Dockerfile" ]; then
  echo "Building template via CLI for $TEMPLATE_NAME..."
  (cd "$TEMPLATE_DIR" && e2b template build)
fi

if [ -f "$TEMPLATE_DIR/build.ts" ]; then
  echo "Building template via SDK for $TEMPLATE_NAME..."
  npx tsx "$TEMPLATE_DIR/build.ts"
fi

echo "Running Python example for $TEMPLATE_NAME..."
python "$TEMPLATE_DIR/example.py"

echo "Running TypeScript example for $TEMPLATE_NAME..."
npx tsx "$TEMPLATE_DIR/example.ts"

echo "All tests passed for $TEMPLATE_NAME."
