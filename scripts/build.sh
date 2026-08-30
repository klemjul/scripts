#!/usr/bin/env bash
set -euo pipefail

# Build all binaries under cmd/
# Usage: ./scripts/build.sh [output-dir]

OUTPUT_DIR="${1:-./bin}"
mkdir -p "$OUTPUT_DIR"

echo "Building all scripts to $OUTPUT_DIR ..."

for dir in cmd/*/; do
    name=$(basename "$dir")
    echo "  -> $name"
    go build -o "$OUTPUT_DIR/$name" "./$dir"
done

echo "Done."
