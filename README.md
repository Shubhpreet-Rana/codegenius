# CodeGenius CLI 🤖✨

An intelligent Git commit message generator and code reviewer powered by AI with a **beautiful terminal user interface**. Use CodeGenius anywhere on your system - just like Firebase CLI!

## 🌟 Features

- **🎨 Beautiful Terminal UI**: Modern, interactive interface with multi-select, input fields, and elegant styling
- **🤖 AI-Powered Commit Messages**: Generate conventional, meaningful commit messages using Google's Gemini AI
- **🔍 Interactive Code Reviews**: Multi-select review types with additional context input (security, performance, style, structure)
- **📊 Work History Tracking**: Visual progress tracking with statistics and filtering
- **💬 Context-Aware**: Add custom context to both commits and reviews for better AI analysis
- **⚙️ Fully Configurable**: Customizable templates, review settings, and project configuration
- **🌍 Global CLI Tool**: Install once, use anywhere on your system

## 🚀 Installation Methods

Choose your preferred installation method:

### 🌐 One-Line Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/codegenius/main/install.sh | bash
```
**✅ Automatic platform detection, system-wide installation, instant setup**

### 📦 NPM (Node.js)
```bash
npm install -g codegenius-cli
```
**✅ Easy updates, clean uninstall, works everywhere Node.js does**

### 🍺 Homebrew (macOS/Linux)
```bash
brew tap yourusername/codegenius
brew install codegenius
```
**✅ Native package management, automatic updates**

### 🔧 Go Install (Developers)
```bash
go install github.com/yourusername/codegenius/cmd@latest
```
**✅ Build from source, latest features**

### 💾 Manual Download
Download the latest release for your platform:
- [macOS Intel](https://github.com/yourusername/codegenius/releases/latest/download/codegenius-darwin-amd64)
- [macOS Apple Silicon](https://github.com/yourusername/codegenius/releases/latest/download/codegenius-darwin-arm64)  
- [Linux x86_64](https://github.com/yourusername/codegenius/releases/latest/download/codegenius-linux-amd64)
- [Linux ARM64](https://github.com/yourusername/codegenius/releases/latest/download/codegenius-linux-arm64)
- [Windows x64](https://github.com/yourusername/codegenius/releases/latest/download/codegenius-windows-amd64.exe)

Then install:
```bash
# macOS/Linux
chmod +x codegenius-*
sudo mv codegenius-* /usr/local/bin/codegenius

# Windows: Move .exe to a folder in your PATH
```

### 📋 All Platforms Supported
| Platform | Curl | NPM | Homebrew | Go | Manual |
|----------|------|-----|----------|----|----|
| **macOS Intel** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **macOS ARM64** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Linux x64** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Linux ARM64** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Windows** | ✅ | ✅ | ❌ | ✅ | ✅ |

**💡 See [INSTALLATION_METHODS.md](INSTALLATION_METHODS.md) for detailed instructions**

## ⚡ Quick Start

### 1. Setup (One-time)
```bash
# Get your API key from: https://makersuite.google.com/app/apikey
export GEMINI_API_KEY="your-gemini-api-key"

# Add to your shell profile
echo 'export GEMINI_API_KEY="your-gemini-api-key"' >> ~/.zshrc
```

### 2. Use in Any Git Repository
```bash
# Navigate to any Git project
cd your-project

# Initialize CodeGenius (creates .codegenius.yaml)
codegenius --init

# Stage your changes
git add .

# Use the beautiful TUI
codegenius --tui

# Or traditional CLI
codegenius          # Generate commit message
codegenius --review # Perform code review
```

## 🎯 Global Usage Examples

### Smart Commit Generation
```bash
# In any Git repository
cd ~/my-awesome-project
git add .
codegenius --tui
```

### Code Review Anywhere
```bash
# Review staged changes in any project
cd ~/work/client-app
git add src/
codegenius --review
```

### Project History
```bash
# View your work history
codegenius --history "Dec 2024"
```

### Global Configuration
```bash
# Your preferences follow you everywhere
codegenius --init  # Creates .codegenius.yaml in current directory
```

## 🎯 Usage

### Basic Commands (Work Anywhere!)

```bash
# Generate commit message for staged changes
codegenius

# Perform code review
codegenius --review

# View work history
codegenius --history "Dec 2024"

# Interactive mode
codegenius --interactive

# Beautiful TUI mode (recommended)
codegenius --tui

# Initialize configuration
codegenius --init

# Show help
codegenius --help
```

### TUI Mode (Recommended)
```bash
# Launch beautiful terminal interface
codegenius --tui
```

The TUI provides:
- **🤖 Smart Commit Generation** with context input
- **🔍 Interactive Code Review** with multi-select options
- **📊 Visual History & Statistics**
- **⚙️ Configuration Management**

## 🌍 How It Works Globally

CodeGenius is designed to work seamlessly across your entire system:

1. **Install Once**: Single installation works everywhere
2. **Per-Project Configuration**: Each project can have its own `.codegenius.yaml`
3. **Global Settings**: Your API key and preferences travel with you
4. **Context Aware**: Automatically detects project language and Git status
5. **Cross-Platform**: Works on macOS, Linux, and Windows

## 📁 Project Structure

```
CLI_GO/
├── cmd/
│   └── main.go              # Global CLI entry point
├── internal/
│   ├── tui/                 # Beautiful terminal UI
│   ├── interfaces/          # Clean architecture
│   ├── container/           # Dependency injection
│   ├── ai/                  # AI integration
│   ├── config/              # Configuration management
│   ├── git/                 # Git operations
│   ├── history/             # Work history tracking
│   └── review/              # Code review functionality
├── .codegenius.yaml         # Project-specific configuration
├── install.sh              # Global installation script
├── go.mod                   # Go module definition
└── README.md               # This file
```

## ⚙️ Configuration

CodeGenius works with both global and project-specific configurations:

### Global Configuration
Your API key and global preferences:
```bash
export GEMINI_API_KEY="your-api-key"
```

### Project Configuration (`.codegenius.yaml`)
Each project can have its own settings:
```yaml
project:
  name: "Your Project"
  language: "go"
  overview: "Project description"
  scopes:
    - core
    - api
    - docs
  standards: "https://golang.org/doc/effective_go.html"

ai:
  model: "gemini-2.0-flash"
  max_tokens: 4000
  context_templates:
    default: "Standard commit message generation"
    bugfix: "Focus on bug fixes and impact"
    feature: "Emphasize new functionality"

review:
  enabled_types:
    - security
    - performance
    - style
    - structure
  text_only: true  # No code snippets in reviews
  security_patterns:
    - '(?i)(password|secret|key|token)\s*[:=]\s*["'"'"'][^"'"'"']+["'"'"']'
```

## 🛠️ Development & Contributing

### For Contributors
```bash
# Clone and setup
git clone https://github.com/codegenius/cli.git
cd CLI_GO
make dev-setup

# Build locally
make build

# Test globally
sudo cp bin/codegenius /usr/local/bin/
codegenius --help
```

### Building Releases
```bash
# Build for multiple platforms
make build-all

# Create release
make release
```

## 🌟 Use Cases

### Individual Developers
```bash
# Work on multiple projects seamlessly
cd ~/project1 && codegenius --tui
cd ~/project2 && codegenius --review
cd ~/project3 && codegenius --history
```

### Teams
```bash
# Consistent commit messages across team
codegenius --init  # Share .codegenius.yaml with team
git add .codegenius.yaml && git commit -m "Add CodeGenius config"
```

### CI/CD Integration
```bash
# Use in build scripts
codegenius --review > review-report.txt
```

## 📦 Distribution

CodeGenius is distributed through multiple channels:

- **GitHub Releases**: Pre-built binaries for all platforms
- **Go Modules**: `go install github.com/codegenius/cli/cmd@latest`
- **Homebrew**: `brew install codegenius` (coming soon)
- **Package Managers**: APT, YUM, Chocolatey (planned)

## 🔧 API Integration

Other applications can integrate CodeGenius:

```go
import "github.com/codegenius/cli/internal/interfaces"

// Use as a library
service := buildCodeGeniusService()
message, err := service.AI.GenerateCommitMessage(diff, files, branch, context)
```

## 💻 System Requirements

- **OS**: macOS, Linux, Windows
- **Git**: Any version (for Git operations)
- **Internet**: For AI features (Gemini API)
- **Go**: Optional (only for `go install` method)

## 🚀 Performance

- **Fast**: Typically generates commits in 2-3 seconds
- **Lightweight**: ~10MB binary, minimal memory usage
- **Offline**: Some features work without internet
- **Concurrent**: Multiple operations can run simultaneously

## 🔒 Privacy & Security

- **No Code Storage**: Your code never leaves your machine (except for AI analysis)
- **Secure API**: Uses HTTPS for all AI communications
- **Local History**: Work history stored locally
- **Configurable**: Control what data is sent to AI

## 📄 License

[MIT License](LICENSE)

## 🙏 Acknowledgments

- Google Gemini AI for intelligent code analysis
- The Charm team for amazing TUI libraries
- The Go community for excellent tooling
- Contributors and users of CodeGenius

---

**Ready to revolutionize your Git workflow? Install CodeGenius and use it anywhere! 🚀✨**

## 💡 Pro Tips

```bash
# Create aliases for faster access
echo 'alias cg="codegenius --tui"' >> ~/.zshrc
echo 'alias cgr="codegenius --review"' >> ~/.zshrc

# Quick commit workflow
git add . && cg

# Quick review workflow
git add . && cgr
```

**Experience the future of Git workflow with CodeGenius CLI - your global AI-powered Git assistant! 🤖⚡** 