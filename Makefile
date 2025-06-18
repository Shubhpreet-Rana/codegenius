# CodeGenius CLI Makefile
.PHONY: build install clean test lint fmt help tui demo dev-setup quick-commit build-all release dist npm-prepare npm-pack npm-publish npm-test homebrew-update

# Variables
BINARY_NAME=codegenius
BUILD_DIR=bin
DIST_DIR=dist
MAIN_FILE=main.go
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

# Platform targets
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
build-all:
	@echo "🔨 Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		platform_split=($${platform//\// }); \
		GOOS=$${platform_split[0]}; \
		GOARCH=$${platform_split[1]}; \
		output_name=$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then output_name+='.exe'; fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $(DIST_DIR)/$$output_name $(MAIN_FILE); \
		if [ $$? -ne 0 ]; then \
			echo "❌ Failed to build for $$GOOS/$$GOARCH"; \
			exit 1; \
		fi; \
	done
	@echo "✅ All platforms built successfully"

# Create distribution packages
dist: build-all
	@echo "📦 Creating distribution packages..."
	@cd $(DIST_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		if [[ $$file == *"windows"* ]]; then \
			zip $${file%.*}.zip $$file; \
		else \
			tar -czf $${file}.tar.gz $$file; \
		fi; \
		echo "Created package for $$file"; \
	done
	@echo "✅ Distribution packages created"

# Create a release
release: clean dist
	@echo "🚀 Creating release $(VERSION)..."
	@echo "📋 Release artifacts:"
	@ls -la $(DIST_DIR)/*.*
	@echo "✅ Release $(VERSION) ready for distribution"

# Install the application globally
install: build
	@echo "📦 Installing $(BINARY_NAME) globally..."
	@if [ -w /usr/local/bin ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"; \
	else \
		echo "🔐 Need sudo privileges to install to /usr/local/bin/"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"; \
	fi
	@echo "🎉 $(BINARY_NAME) is now available globally!"
	@echo "   Try: $(BINARY_NAME) --help"

# Install via Go (for users)
go-install:
	@echo "📥 Installing $(BINARY_NAME) via go install..."
	@go install github.com/Shubhpreet-Rana/codegenius@latest
	@echo "✅ Installation complete via Go modules"

# Uninstall the application
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@if [ -f /usr/local/bin/$(BINARY_NAME) ]; then \
		if [ -w /usr/local/bin ]; then \
			rm /usr/local/bin/$(BINARY_NAME); \
		else \
			sudo rm /usr/local/bin/$(BINARY_NAME); \
		fi; \
		echo "✅ $(BINARY_NAME) uninstalled from /usr/local/bin/"; \
	else \
		echo "⚠️  $(BINARY_NAME) not found in /usr/local/bin/"; \
	fi

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@go clean
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, skipping extended linting"; \
	fi

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod tidy

# Run the application locally
run:
	@go run $(MAIN_FILE) $(ARGS)

# Run TUI mode
tui:
	@echo "🚀 Launching CodeGenius TUI..."
	@go run $(MAIN_FILE) --tui

# Run interactive demo
demo:
	@echo "🎬 Running CodeGenius demo..."
	@echo "This will show the available features:"
	@go run $(MAIN_FILE) --help

# Development commands
dev-setup:
	@echo "⚙️  Setting up development environment..."
	@go mod tidy
	@echo "✅ Development environment ready"

# Quick commit (for development)
quick-commit:
	@echo "⚡ Quick commit with CodeGenius..."
	@go run $(MAIN_FILE)

# Global installation test
test-global:
	@echo "🧪 Testing global installation..."
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "✅ $(BINARY_NAME) is available globally"; \
		$(BINARY_NAME) --help | head -5; \
	else \
		echo "❌ $(BINARY_NAME) not found in PATH"; \
		echo "Run 'make install' or 'make go-install' first"; \
	fi

# Check version
version:
	@echo "📋 CodeGenius Version Information:"
	@echo "Version: $(VERSION)"
	@echo "Go Version: $(shell go version)"
	@echo "Git Commit: $(shell git rev-parse --short HEAD)"
	@echo "Build Date: $(shell date)"

# Show help
help:
	@echo "🤖 CodeGenius CLI - Available Commands:"
	@echo ""
	@echo "📋 Build Commands:"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for all platforms (cross-compile)"
	@echo "  dist         - Create distribution packages"
	@echo "  release      - Create a full release with all artifacts"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "📦 Installation Commands:"
	@echo "  install      - Install globally to /usr/local/bin"
	@echo "  go-install   - Install via 'go install' (for users)"
	@echo "  uninstall    - Remove global installation"
	@echo "  test-global  - Test if globally installed"
	@echo ""
	@echo "🧪 Development Commands:"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  dev-setup    - Setup development environment"
	@echo ""
	@echo "🚀 Usage Commands:"
	@echo "  run          - Run the application locally (use ARGS=... for arguments)"
	@echo "  tui          - Launch beautiful terminal UI"
	@echo "  demo         - Show application features"
	@echo "  quick-commit - Generate commit for staged changes"
	@echo ""
	@echo "📖 Examples:"
	@echo "  make build              # Build the application"
	@echo "  make install            # Install globally"
	@echo "  make build-all          # Build for all platforms"
	@echo "  make release            # Create full release"
	@echo "  make go-install         # Install via Go modules"
	@echo "  make run ARGS='--help'  # Show help"
	@echo ""
	@echo "🌍 Global Usage After Installation:"
	@echo "  codegenius --tui        # Use anywhere on your system"
	@echo "  codegenius --help       # Global help"
	@echo ""
	@echo "💡 For users: Use the one-line installer:"
	@echo "  curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/main/install.sh | bash"

# NPM packaging targets
npm-prepare: dist
	@echo "📦 Preparing NPM package..."
	mkdir -p lib/bin
	cp dist/* lib/bin/ 2>/dev/null || true
	chmod +x bin/codegenius.js
	chmod +x scripts/*.js

npm-pack: npm-prepare
	@echo "📦 Creating NPM package..."
	npm pack
	@echo "✅ NPM package created!"

npm-publish: npm-prepare
	@echo "🚀 Publishing to NPM..."
	npm publish --access public
	@echo "✅ Published to NPM!"

npm-test:
	@echo "🧪 Testing NPM package..."
	npm test

# Homebrew formula update
homebrew-update:
	@echo "🍺 Updating Homebrew formula..."
	@echo "⚠️  Remember to update SHA256 hashes in Formula/codegenius.rb"
	@echo "   Get hashes with: shasum -a 256 dist/*" 