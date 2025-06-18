#!/bin/bash

# 🍺 CodeGenius Homebrew Tap Setup Script
# This script automates the creation of a Homebrew tap for CodeGenius

set -e

GITHUB_USERNAME="Shubhpreet-Rana"
TAP_REPO_NAME="homebrew-codegenius"
MAIN_REPO_NAME="codegenius"

echo "🍺 Setting up Homebrew tap for CodeGenius..."
echo ""

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is required but not installed."
    echo "Install it with: brew install gh"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Please authenticate with GitHub CLI first:"
    echo "gh auth login"
    exit 1
fi

echo "✅ GitHub CLI is ready"

# Step 1: Create the tap repository
echo ""
echo "📦 Step 1: Creating tap repository..."
if gh repo view "$GITHUB_USERNAME/$TAP_REPO_NAME" &> /dev/null; then
    echo "⚠️  Repository $TAP_REPO_NAME already exists. Updating it..."
else
    echo "Creating new repository: $TAP_REPO_NAME"
    gh repo create "$TAP_REPO_NAME" --public --description "Homebrew tap for CodeGenius CLI" --clone
    cd "$TAP_REPO_NAME"
fi

# Step 2: Setup tap structure
echo ""
echo "📁 Step 2: Setting up tap structure..."
mkdir -p Formula
mkdir -p .github/workflows

# Step 3: Copy formula from main repo
echo ""
echo "📋 Step 3: Copying formula..."
if [ -f "../$MAIN_REPO_NAME/Formula/codegenius.rb" ]; then
    cp "../$MAIN_REPO_NAME/Formula/codegenius.rb" Formula/
    echo "✅ Formula copied successfully"
else
    echo "❌ Formula not found at ../$MAIN_REPO_NAME/Formula/codegenius.rb"
    echo "Please run this script from the directory containing both repositories"
    exit 1
fi

# Step 4: Create README for tap
echo ""
echo "📝 Step 4: Creating tap README..."
cat > README.md << 'EOF'
# 🍺 Homebrew Tap for CodeGenius CLI

This is the official Homebrew tap for [CodeGenius CLI](https://github.com/Shubhpreet-Rana/codegenius) - an AI-powered Git commit message generator and code reviewer.

## Installation

```bash
# Add the tap
brew tap Shubhpreet-Rana/codegenius

# Install CodeGenius
brew install codegenius
```

## Usage

```bash
# Get your Gemini API key from: https://makersuite.google.com/app/apikey
export GEMINI_API_KEY="your-gemini-api-key"

# Use anywhere in your Git projects
cd your-project
git add .
codegenius --tui
```

## Alternative Installation Methods

- **Curl installer**: `curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/install.sh | bash`
- **Go install**: `go install github.com/Shubhpreet-Rana/codegenius@latest`
- **NPM package**: `npm install -g codegenius-cli`

## Documentation

See the [main repository](https://github.com/Shubhpreet-Rana/codegenius) for full documentation and features.

## Support

If you encounter any issues with the Homebrew installation, please [open an issue](https://github.com/Shubhpreet-Rana/codegenius/issues) in the main repository.
EOF

# Step 5: Create GitHub Actions workflow for automated testing
echo ""
echo "🔄 Step 5: Creating CI workflow..."
cat > .github/workflows/test.yml << 'EOF'
name: Test Formula

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Homebrew
      id: set-up-homebrew
      uses: Homebrew/actions/setup-homebrew@master
      
    - name: Test formula
      run: |
        brew install --formula ./Formula/codegenius.rb
        codegenius --help
        
    - name: Test uninstall
      run: brew uninstall codegenius
EOF

# Step 6: Commit and push
echo ""
echo "🚀 Step 6: Committing and pushing..."
git add .
git commit -m "feat: initial Homebrew tap setup for CodeGenius CLI

- Add codegenius formula with real SHA256 hashes
- Add tap README with installation instructions
- Add GitHub Actions workflow for automated testing
- Support for macOS (Intel/ARM) and Linux (x64/ARM64)

Ready for: brew tap Shubhpreet-Rana/codegenius && brew install codegenius" || echo "No changes to commit"

git push origin main

echo ""
echo "🎉 Homebrew tap setup complete!"
echo ""
echo "📋 Next steps for users:"
echo "1. brew tap $GITHUB_USERNAME/codegenius"
echo "2. brew install codegenius"
echo ""
echo "🔗 Tap repository: https://github.com/$GITHUB_USERNAME/$TAP_REPO_NAME"
echo ""
echo "✅ CodeGenius is now available via Homebrew!" 