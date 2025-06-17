#!/bin/bash

# CodeGenius CLI Deployment Script
# Deploys to GitHub, NPM, and prepares Homebrew

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
VERSION=${1:-"1.0.0"}
REPO_OWNER="yourusername"
REPO_NAME="codegenius"

echo -e "${BLUE}🚀 CodeGenius CLI Deployment${NC}"
echo -e "${BLUE}==============================${NC}"
echo "Version: $VERSION"
echo ""

# Step 1: Build all platforms
echo -e "${BLUE}📦 Step 1: Building all platforms...${NC}"
make build-all
echo -e "${GREEN}✅ All platforms built${NC}"
echo ""

# Step 2: Commit and tag
echo -e "${BLUE}📝 Step 2: Committing changes...${NC}"
git add .
git commit -m "feat: release v$VERSION with multiple installation methods" || echo "No changes to commit"

# Create tag
git tag "v$VERSION" || echo "Tag already exists"
echo -e "${GREEN}✅ Tagged version v$VERSION${NC}"
echo ""

# Step 3: GitHub Release
echo -e "${BLUE}🐙 Step 3: GitHub Release...${NC}"
echo "Push to GitHub:"
echo -e "${YELLOW}git push origin main${NC}"
echo -e "${YELLOW}git push origin v$VERSION${NC}"
echo ""
echo "Then create release manually or use GitHub CLI:"
echo -e "${YELLOW}gh release create v$VERSION dist/* --title \"CodeGenius CLI v$VERSION\" --notes \"Multi-platform release with curl, npm, and homebrew support\"${NC}"
echo ""

# Step 4: NPM Package
echo -e "${BLUE}📦 Step 4: NPM Package...${NC}"
make npm-prepare
echo ""
echo "To publish to NPM:"
echo -e "${YELLOW}npm login${NC}"
echo -e "${YELLOW}npm publish --access public${NC}"
echo ""

# Step 5: Homebrew
echo -e "${BLUE}🍺 Step 5: Homebrew Formula...${NC}"
echo "Update SHA256 hashes in Formula/codegenius.rb:"
echo ""
shasum -a 256 dist/* | while read hash file; do
    filename=$(basename "$file")
    echo -e "${YELLOW}$filename: $hash${NC}"
done
echo ""
echo "Then commit to homebrew-codegenius repository"
echo ""

# Installation URLs
echo -e "${GREEN}🎉 After deployment, users can install with:${NC}"
echo ""
echo -e "${BLUE}Curl:${NC}"
echo "curl -fsSL https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/main/install.sh | bash"
echo ""
echo -e "${BLUE}NPM:${NC}"
echo "npm install -g codegenius-cli"
echo ""
echo -e "${BLUE}Homebrew:${NC}"
echo "brew tap $REPO_OWNER/codegenius"
echo "brew install codegenius"
echo ""

echo -e "${GREEN}🚀 Deployment preparation complete!${NC}"
echo -e "${YELLOW}⚠️  Remember to:${NC}"
echo "1. Push to GitHub"
echo "2. Create GitHub release"
echo "3. Publish to NPM"
echo "4. Update Homebrew tap" 