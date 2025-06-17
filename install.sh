#!/bin/bash

# CodeGenius CLI Installation Script

set -e

echo "🚀 Installing CodeGenius CLI..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19 or later first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
REQUIRED_VERSION="1.19"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please upgrade to Go 1.19 or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Build the application
echo "🔨 Building CodeGenius CLI..."
go mod tidy
go build -o codegenius cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

# Make executable
chmod +x codegenius

# Suggest installation location
echo ""
echo "🎉 CodeGenius CLI built successfully!"
echo ""
echo "To install globally, run one of these commands:"
echo "  sudo mv codegenius /usr/local/bin/          # Install system-wide"
echo "  mv codegenius ~/.local/bin/                 # Install for current user"
echo ""
echo "Or use from current directory:"
echo "  ./codegenius --help"
echo ""
echo "📋 Next steps:"
echo "1. Get your Gemini API key from: https://makersuite.google.com/app/apikey"
echo "2. Set environment variable: export GEMINI_API_KEY=\"your-key-here\""
echo "3. Initialize project config: ./codegenius --init"
echo "4. Stage some changes and run: ./codegenius"
echo ""
echo "📖 See README.md for detailed usage instructions." 