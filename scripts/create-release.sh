#!/bin/bash

# 🚀 CodeGenius GitHub Release Creation Script
# This script creates a GitHub release v1.1.2 with all platform binaries

set -e

VERSION="v1.1.2"
RELEASE_TITLE="CodeGenius CLI v1.1.2 - Homebrew Ready 🍺"

echo "🚀 Creating GitHub release $VERSION for CodeGenius CLI..."
echo ""

# Check if GitHub CLI is installed and authenticated
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is required but not installed."
    echo "Install it with: brew install gh"
    exit 1
fi

if ! gh auth status &> /dev/null; then
    echo "❌ Please authenticate with GitHub CLI first:"
    echo "gh auth login"
    echo ""
    echo "After authentication, run this script again."
    exit 1
fi

echo "✅ GitHub CLI is ready"

# Check if binaries exist
echo ""
echo "🔍 Checking for built binaries..."
BINARIES=(
    "dist/codegenius-darwin-amd64"
    "dist/codegenius-darwin-arm64"
    "dist/codegenius-linux-amd64"
    "dist/codegenius-linux-arm64"
    "dist/codegenius-windows-amd64.exe"
)

for binary in "${BINARIES[@]}"; do
    if [[ ! -f "$binary" ]]; then
        echo "❌ Binary not found: $binary"
        echo "Please run 'make build-all' first to build all binaries."
        exit 1
    else
        echo "✅ Found: $binary"
    fi
done

# Check if release already exists
echo ""
echo "🔍 Checking if release $VERSION already exists..."
if gh release view $VERSION &> /dev/null; then
    echo "⚠️  Release $VERSION already exists!"
    echo ""
    read -p "Do you want to delete and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🗑️  Deleting existing release..."
        gh release delete $VERSION --yes
        echo "✅ Existing release deleted"
    else
        echo "❌ Aborted. Please delete the existing release manually or use a different version."
        exit 1
    fi
fi

# Create the release
echo ""
echo "📦 Creating release $VERSION..."

gh release create $VERSION \
    --title "$RELEASE_TITLE" \
    --notes-file RELEASE_NOTES.md \
    --draft=false \
    --prerelease=false \
    "${BINARIES[@]}"

echo ""
echo "🎉 Release $VERSION created successfully!"
echo ""
echo "🔗 Release URL: https://github.com/$(gh repo view --json owner,name -q '.owner.login + "/" + .name')/releases/tag/$VERSION"
echo ""
echo "✅ Now you can install via Homebrew:"
echo "   brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb"
echo ""
echo "🍺 Or set up the Homebrew tap:"
echo "   ./scripts/setup-homebrew-tap.sh"
echo ""
echo "🚀 CodeGenius CLI v1.1.2 is now live and ready for global installation!" 