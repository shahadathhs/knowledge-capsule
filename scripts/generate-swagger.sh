#!/bin/bash
set -e

# Set GOBIN if not already set
GOBIN="${GOBIN:-$(pwd)/.bin}"

# Ensure swag is installed
if ! command -v "$GOBIN/swag" &> /dev/null; then
    echo "⚠️  swag not found, installing..."
    mkdir -p "$GOBIN"
    GO111MODULE=on go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
echo "📝 Generating Swagger documentation..."
"$GOBIN/swag" init -g main.go --output docs
echo "✅ Swagger docs generated in docs/"
