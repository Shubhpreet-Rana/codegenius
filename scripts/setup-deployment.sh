#!/bin/bash

# 🚀 CodeGenius CLI - Deployment Pipeline Setup Script
# This script helps configure the automated deployment pipeline

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Icons
ROCKET="🚀"
CHECK="✅"
WARNING="⚠️"
INFO="ℹ️"
GEAR="⚙️"

echo -e "${BLUE}${ROCKET} CodeGenius CLI - Deployment Pipeline Setup${NC}"
echo "=============================================="
echo ""

# Check if GitHub CLI is installed
check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        echo -e "${RED}${WARNING} GitHub CLI (gh) is not installed.${NC}"
        echo "Please install it first:"
        echo "  macOS: brew install gh"
        echo "  Linux: https://github.com/cli/cli/blob/trunk/docs/install_linux.md"
        echo "  Windows: https://github.com/cli/cli/releases"
        exit 1
    fi
}

# Check if user is logged into GitHub CLI
check_gh_auth() {
    if ! gh auth status &> /dev/null; then
        echo -e "${YELLOW}${WARNING} Not logged into GitHub CLI.${NC}"
        echo "Please login first:"
        echo "  gh auth login"
        echo ""
        read -p "Press Enter after logging in..."
    fi
}

# Get repository information
get_repo_info() {
    REPO_OWNER=$(gh repo view --json owner --jq '.owner.login' 2>/dev/null || echo "")
    REPO_NAME=$(gh repo view --json name --jq '.name' 2>/dev/null || echo "")
    
    if [[ -z "$REPO_OWNER" || -z "$REPO_NAME" ]]; then
        echo -e "${RED}${WARNING} Could not determine repository info.${NC}"
        echo "Make sure you're in a Git repository with GitHub remote."
        exit 1
    fi
    
    echo -e "${GREEN}${INFO} Repository: ${REPO_OWNER}/${REPO_NAME}${NC}"
}

# Setup NPM token
setup_npm_token() {
    echo ""
    echo -e "${BLUE}${GEAR} Setting up NPM Token${NC}"
    echo "----------------------------"
    echo ""
    echo "1. First, login to NPM if you haven't already:"
    echo "   ${CYAN}npm login${NC}"
    echo ""
    echo "2. Create an automation token:"
    echo "   ${CYAN}npm token create --type=automation${NC}"
    echo ""
    echo "3. Copy the token that starts with 'npm_...'"
    echo ""
    
    read -p "Have you created the NPM token? (y/n): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        read -s -p "Enter your NPM token: " NPM_TOKEN
        echo ""
        
        if [[ -n "$NPM_TOKEN" ]]; then
            gh secret set NPM_TOKEN --body "$NPM_TOKEN" --repo "$REPO_OWNER/$REPO_NAME"
            echo -e "${GREEN}${CHECK} NPM_TOKEN secret added successfully!${NC}"
        else
            echo -e "${RED}${WARNING} No token provided. Skipping NPM token setup.${NC}"
        fi
    else
        echo -e "${YELLOW}${WARNING} Skipping NPM token setup. You can add it later.${NC}"
    fi
}

# Setup GitHub token for Homebrew tap
setup_tap_token() {
    echo ""
    echo -e "${BLUE}${GEAR} Setting up Homebrew Tap Token${NC}"
    echo "-----------------------------------"
    echo ""
    echo "1. Go to: ${CYAN}https://github.com/settings/tokens${NC}"
    echo "2. Click 'Generate new token (classic)'"
    echo "3. Select scope: ${CYAN}repo${NC} (Full control of private repositories)"
    echo "4. Copy the token that starts with 'ghp_...'"
    echo ""
    
    read -p "Have you created the GitHub token? (y/n): " -n 1 -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo ""
        read -s -p "Enter your GitHub token: " TAP_TOKEN
        echo ""
        
        if [[ -n "$TAP_TOKEN" ]]; then
            gh secret set TAP_GITHUB_TOKEN --body "$TAP_TOKEN" --repo "$REPO_OWNER/$REPO_NAME"
            echo -e "${GREEN}${CHECK} TAP_GITHUB_TOKEN secret added successfully!${NC}"
        else
            echo -e "${RED}${WARNING} No token provided. Skipping tap token setup.${NC}"
        fi
    else
        echo -e "${YELLOW}${WARNING} Skipping tap token setup. You can add it later.${NC}"
    fi
}

# Check if homebrew tap repository exists
check_tap_repo() {
    echo ""
    echo -e "${BLUE}${GEAR} Checking Homebrew Tap Repository${NC}"
    echo "------------------------------------"
    
    TAP_REPO="homebrew-codegenius"
    TAP_FULL_NAME="$REPO_OWNER/$TAP_REPO"
    
    if gh repo view "$TAP_FULL_NAME" &> /dev/null; then
        echo -e "${GREEN}${CHECK} Homebrew tap repository exists: ${TAP_FULL_NAME}${NC}"
    else
        echo -e "${YELLOW}${WARNING} Homebrew tap repository doesn't exist.${NC}"
        echo ""
        read -p "Create homebrew tap repository? (y/n): " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Creating repository: $TAP_FULL_NAME..."
            gh repo create "$TAP_REPO" --public --description "Homebrew tap for CodeGenius CLI"
            
            # Clone and setup the tap repository
            TEMP_DIR=$(mktemp -d)
            cd "$TEMP_DIR"
            gh repo clone "$TAP_FULL_NAME"
            cd "$TAP_REPO"
            
            # Create Formula directory and copy formula
            mkdir -p Formula
            if [[ -f "../../Formula/codegenius.rb" ]]; then
                cp "../../Formula/codegenius.rb" Formula/
            else
                echo -e "${YELLOW}${WARNING} Formula not found. You'll need to add it manually.${NC}"
            fi
            
            # Create README
            cat > README.md << 'EOF'
# Homebrew Tap for CodeGenius CLI

## Installation

```bash
brew tap REPO_OWNER/codegenius
brew install codegenius
```

## About CodeGenius

AI-powered Git commit message generator and code reviewer with beautiful TUI.

## Repository

https://github.com/REPO_OWNER/codegenius
EOF
            
            sed -i.bak "s/REPO_OWNER/$REPO_OWNER/g" README.md && rm README.md.bak
            
            git add .
            git commit -m "Initial tap setup for CodeGenius CLI"
            git push
            
            cd - > /dev/null
            rm -rf "$TEMP_DIR"
            
            echo -e "${GREEN}${CHECK} Homebrew tap repository created and configured!${NC}"
        else
            echo -e "${YELLOW}${WARNING} Skipping tap repository creation.${NC}"
        fi
    fi
}

# Verify current deployment workflow
check_workflow() {
    echo ""
    echo -e "${BLUE}${GEAR} Checking Deployment Workflow${NC}"
    echo "-------------------------------"
    
    if [[ -f ".github/workflows/deploy.yml" ]]; then
        echo -e "${GREEN}${CHECK} Deployment workflow exists!${NC}"
    else
        echo -e "${RED}${WARNING} Deployment workflow not found.${NC}"
        echo "The deploy.yml file should be at: .github/workflows/deploy.yml"
    fi
}

# Show summary and next steps
show_summary() {
    echo ""
    echo -e "${PURPLE}${ROCKET} Setup Summary${NC}"
    echo "================"
    echo ""
    
    # Check secrets
    echo "Repository Secrets:"
    
    if gh secret list --repo "$REPO_OWNER/$REPO_NAME" | grep -q "NPM_TOKEN"; then
        echo -e "  ${GREEN}${CHECK} NPM_TOKEN${NC}"
    else
        echo -e "  ${RED}❌ NPM_TOKEN${NC} - Add with: gh secret set NPM_TOKEN"
    fi
    
    if gh secret list --repo "$REPO_OWNER/$REPO_NAME" | grep -q "TAP_GITHUB_TOKEN"; then
        echo -e "  ${GREEN}${CHECK} TAP_GITHUB_TOKEN${NC}"
    else
        echo -e "  ${RED}❌ TAP_GITHUB_TOKEN${NC} - Add with: gh secret set TAP_GITHUB_TOKEN"
    fi
    
    echo ""
    echo "Next Steps:"
    echo "1. ${CYAN}Push to 'latest' branch to trigger deployment${NC}"
    echo "   git checkout latest && git push origin latest"
    echo ""
    echo "2. ${CYAN}Or trigger manually with version control:${NC}"
    echo "   Go to: https://github.com/$REPO_OWNER/$REPO_NAME/actions"
    echo "   Run: '🚀 Auto Deploy to All Channels'"
    echo ""
    echo "3. ${CYAN}Monitor deployment:${NC}"
    echo "   gh run list --workflow=\"deploy.yml\""
    echo "   gh run watch"
    echo ""
    echo "4. ${CYAN}Verify installation methods:${NC}"
    echo "   npm install -g codegenius-cli"
    echo "   brew tap $REPO_OWNER/codegenius && brew install codegenius"
    echo ""
    echo -e "${GREEN}${ROCKET} Deployment pipeline ready!${NC}"
}

# Main execution
main() {
    echo -e "${INFO} Starting deployment pipeline setup...\n"
    
    check_gh_cli
    check_gh_auth
    get_repo_info
    setup_npm_token
    setup_tap_token
    check_tap_repo
    check_workflow
    show_summary
}

# Run main function
main "$@" 