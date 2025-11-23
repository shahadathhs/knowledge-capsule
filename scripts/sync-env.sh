#!/usr/bin/env bash
set -e

ENV_FILE=".env"

if [ ! -f "$ENV_FILE" ]; then
    echo "⚠️  No .env file found in current path"
    exit 1
fi

# Detect OS type
OS_TYPE=$(uname | tr '[:upper:]' '[:lower:]')

# Function to export for Bash/Zsh
export_bash() {
    while IFS='=' read -r key value || [ -n "$key" ]; do
        # Skip comments and empty lines
        [[ "$key" =~ ^#.* ]] && continue
        [[ -z "$key" ]] && continue

        # Trim whitespace
        key=$(echo $key | xargs)
        value=$(echo $value | xargs)

        export "$key=$value"
    done < "$ENV_FILE"
    echo "✅ Environment variables loaded for Bash/Zsh"
}

# Function to print PowerShell commands
print_powershell() {
    echo "⚡ To load variables in PowerShell, run:"
    while IFS='=' read -r key value || [ -n "$key" ]; do
        [[ "$key" =~ ^#.* ]] && continue
        [[ -z "$key" ]] && continue
        key=$(echo $key | xargs)
        value=$(echo $value | xargs)
        echo "\$env:$key = '$value'"
    done < "$ENV_FILE"
}

# Function to print CMD commands
print_cmd() {
    echo "⚡ To load variables in CMD, run:"
    while IFS='=' read -r key value || [ -n "$key" ]; do
        [[ "$key" =~ ^#.* ]] && continue
        [[ -z "$key" ]] && continue
        key=$(echo $key | xargs)
        value=$(echo $value | xargs)
        echo "set $key=$value"
    done < "$ENV_FILE"
}

# Determine what to do based on OS
case "$OS_TYPE" in
    linux*|darwin*|gnu*|msys*|cygwin*)
        # For Linux/macOS or Git Bash/WSL, export directly
        export_bash
        ;;
    *)
        # For unknown, fallback to printing commands
        echo "⚠️  Unknown OS type: $OS_TYPE"
        echo "Run the appropriate export commands manually:"
        print_powershell
        print_cmd
        ;;
esac
