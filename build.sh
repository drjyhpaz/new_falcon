#!/bin/bash

# Falcon Build Script

set -e

echo "🦅 Building Falcon..."

# Create output directory
mkdir -p bin

# Build for different platforms
echo "📦 Building for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/falcon-linux -ldflags "-s -w" .

echo "📦 Building for Windows..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/falcon-windows.exe -ldflags "-s -w" .

echo "📦 Building for macOS..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/falcon-darwin -ldflags "-s -w" .

echo "✅ Build complete!"
echo "📁 Binaries in ./bin/"
