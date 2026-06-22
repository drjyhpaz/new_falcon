#!/bin/bash

# Falcon Installation Script

set -e

echo "🦅 Falcon Installation Script"
echo "================================"

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✅ Go version: $GO_VERSION"

# Clone repository if not already cloned
if [ ! -d ".git" ]; then
    echo "📥 Cloning repository..."
    git clone https://github.com/drjyhpaz/new_falcon.git
    cd new_falcon
fi

echo "📦 Downloading dependencies..."
go mod download

echo "📦 Building Falcon..."
mkdir -p bin
go build -o bin/falcon -ldflags "-s -w" .

echo "✅ Installation complete!"
echo "📁 Falcon binary: ./bin/falcon"
echo "🨶 Run './bin/falcon --help' for more information"
