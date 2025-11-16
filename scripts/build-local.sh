#!/usr/bin/env bash
set -e

echo "🔨 Building local Go binary..."

# Variables
BUILD_DIR="tmp"
BINARY_NAME="server"
BUILD_OUT="$BUILD_DIR/$BINARY_NAME"

# Create build dir if not exists
mkdir -p "$BUILD_DIR"

# Build binary
go build -o "$BUILD_OUT" ./main.go

echo "✅ Built local binary: $BUILD_OUT"
