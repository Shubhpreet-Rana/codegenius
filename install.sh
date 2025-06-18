#!/bin/bash

# CodeGenius CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Constants
REPO_OWNER="Shubhpreet-Rana"
REPO_NAME="codegenius"
BINARY_NAME="codegenius"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/codegenius-install"
GITHUB_API="https://api.github.com/repos"
GITHUB_RELEASES="https://github.com/repos"

# Functions
print_header() {
    echo -e "${PURPLE}"
    echo "┌─────────────────────────────────────────┐"
    echo "│           🤖 CodeGenius CLI            │"
    echo "│     AI-Powered Git Assistant            │"
    echo "└─────────────────────────────────────────┘"
    echo -e "${NC}"
}

print_step() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $arch in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) print_error "Unsupported architecture: $arch"; exit 1 ;;
    esac
    
    case $os in
        darwin) os="darwin" ;;
        linux) os="linux" ;;
        mingw*|msys*|cygwin*) os="windows" ;;
        *) print_error "Unsupported OS: $os"; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local api_url="${GITHUB_API}/${REPO_OWNER}/${REPO_NAME}/releases/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi
}

# Download file
download_file() {
    local url=$1
    local output=$2
    
    print_step "Downloading from $url..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$output"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$url" -O "$output"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi
}

# Check if Go is installed
check_go() {
    command -v go >/dev/null 2>&1
}

# Install via Go
install_via_go() {
    print_step "Installing CodeGenius via Go..."
    
    go install github.com/${REPO_OWNER}/${REPO_NAME}@latest
    
    local gobin=""
    if [ -n "$GOBIN" ]; then
        gobin="$GOBIN"
    elif [ -n "$GOPATH" ]; then
        gobin="$GOPATH/bin"
    else
        gobin="$(go env GOPATH)/bin"
    fi
    
    if [[ ":$PATH:" == *":$gobin:"* ]]; then
        print_success "CodeGenius installed successfully via Go!"
    else
        print_warning "Go binary path ($gobin) is not in your PATH."
        print_step "Add this to your shell profile (~/.bashrc, ~/.zshrc):"
        echo "export PATH=\"$gobin:\$PATH\""
    fi
}

# Install via binary download
install_via_binary() {
    local platform=$(detect_platform)
    local version=$(get_latest_version)
    
    if [ -z "$version" ]; then
        print_error "Could not determine latest version"
        exit 1
    fi
    
    print_step "Installing CodeGenius $version for $platform..."
    
    # Create temp directory
    mkdir -p "$TEMP_DIR"
    cd "$TEMP_DIR"
    
    # Construct download URL
    local binary_name="${BINARY_NAME}-${platform}"
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="${binary_name}.exe"
    fi
    
    local download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/${binary_name}"
    
    # Download binary
    download_file "$download_url" "$binary_name"
    
    # Make executable
    chmod +x "$binary_name"
    
    # Install to system directory
    local final_name="$BINARY_NAME"
    if [[ "$platform" == *"windows"* ]]; then
        final_name="${BINARY_NAME}.exe"
    fi
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_name" "$INSTALL_DIR/$final_name"
        print_success "CodeGenius installed to $INSTALL_DIR/$final_name"
    else
        print_step "Need sudo privileges to install to $INSTALL_DIR"
        sudo mv "$binary_name" "$INSTALL_DIR/$final_name"
        print_success "CodeGenius installed to $INSTALL_DIR/$final_name"
    fi
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
}

# Verify installation
verify_installation() {
    print_step "Verifying installation..."
    
    # First check if it's in PATH
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version_output=$("$BINARY_NAME" --help 2>/dev/null | head -1 || echo "CodeGenius CLI")
        print_success "CodeGenius is installed and working!"
        print_step "Version: $version_output"
        print_step "Location: $(which $BINARY_NAME)"
        return 0
    fi
    
    # If not in PATH, check Go bin directory
    local gobin=""
    if [ -n "$GOBIN" ]; then
        gobin="$GOBIN"
    elif [ -n "$GOPATH" ]; then
        gobin="$GOPATH/bin"
    else
        gobin="$(go env GOPATH)/bin"
    fi
    
    if [ -f "$gobin/$BINARY_NAME" ]; then
        local version_output=$("$gobin/$BINARY_NAME" --help 2>/dev/null | head -1 || echo "CodeGenius CLI")
        print_success "CodeGenius is installed and working!"
        print_step "Version: $version_output"
        print_step "Location: $gobin/$BINARY_NAME"
        print_warning "Note: $gobin is not in your PATH. Add it for global access."
        return 0
    fi
    
    print_error "Installation verification failed"
    return 1
}

# Setup instructions
show_setup_instructions() {
    echo ""
    echo -e "${GREEN}🎉 CodeGenius CLI is now installed!${NC}"
    echo ""
    echo -e "${BLUE}📋 Next Steps:${NC}"
    echo "1. Get your Gemini API key: https://makersuite.google.com/app/apikey"
    echo "2. Set environment variable:"
    echo -e "   ${YELLOW}export GEMINI_API_KEY=\"your-api-key-here\"${NC}"
    echo "3. Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)"
    echo "4. Initialize in any Git repository:"
    echo -e "   ${YELLOW}cd your-project && codegenius --init${NC}"
    echo "5. Start using CodeGenius:"
    echo -e "   ${YELLOW}codegenius --tui${NC}"
    echo ""
    echo -e "${BLUE}🚀 Quick Commands:${NC}"
    echo -e "   ${YELLOW}codegenius --help${NC}     # Show help"
    echo -e "   ${YELLOW}codegenius --tui${NC}      # Beautiful interface"
    echo -e "   ${YELLOW}codegenius --review${NC}   # Code review"
    echo ""
    echo -e "${BLUE}🔗 Resources:${NC}"
    echo "   Documentation: https://github.com/${REPO_OWNER}/${REPO_NAME}#readme"
    echo "   Issues: https://github.com/${REPO_OWNER}/${REPO_NAME}/issues"
    echo ""
}

# Main installation flow
main() {
    print_header
    
    # Check for existing installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_location=$(which $BINARY_NAME)
        print_warning "CodeGenius is already installed at $current_location"
        read -p "Do you want to update it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_step "Installation cancelled."
            exit 0
        fi
    fi
    
    # Choose installation method
    if check_go && [ "${FORCE_BINARY:-}" != "true" ]; then
        print_step "Go detected. Choose installation method:"
        echo "1) Install via Go (recommended for developers)"
        echo "2) Install pre-built binary"
        read -p "Choice (1 or 2): " -n 1 -r
        echo
        if [[ $REPLY == "1" ]]; then
            install_via_go
        else
            install_via_binary
        fi
    else
        print_step "Installing pre-built binary..."
        install_via_binary
    fi
    
    # Verify installation
    if verify_installation; then
        show_setup_instructions
    else
        print_error "Installation failed. Please check the error messages above."
        exit 1
    fi
}

# Handle command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --force-binary)
            FORCE_BINARY="true"
            shift
            ;;
        --help)
            echo "CodeGenius CLI Installer"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --force-binary    Force binary installation even if Go is available"
            echo "  --help           Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  REPO_OWNER       GitHub repository owner (default: yourusername)"
            echo "  REPO_NAME        GitHub repository name (default: codegenius)"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run main function
main "$@" 