# CodeGenius CLI - Global Distribution Guide 🌍

This guide explains how to distribute and use CodeGenius as a global CLI tool, similar to Firebase CLI, AWS CLI, or other popular developer tools.

## 🚀 For Users: How to Install CodeGenius Globally

### Method 1: One-Line Installer (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/codegenius/cli/main/install.sh | bash
```

### Method 2: Via Go (if you have Go installed)
```bash
go install github.com/codegenius/cli/cmd@latest
```

### Method 3: Download Pre-built Binary
1. Go to [GitHub Releases](https://github.com/codegenius/cli/releases)
2. Download the binary for your platform:
   - `codegenius-darwin-amd64` (macOS Intel)
   - `codegenius-darwin-arm64` (macOS Apple Silicon)
   - `codegenius-linux-amd64` (Linux x86_64)
   - `codegenius-linux-arm64` (Linux ARM64)
   - `codegenius-windows-amd64.exe` (Windows)

3. Install globally:
```bash
# macOS/Linux
sudo mv codegenius /usr/local/bin/
chmod +x /usr/local/bin/codegenius

# Or to user directory (add ~/.local/bin to PATH)
mkdir -p ~/.local/bin
mv codegenius ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

### Method 4: Via Package Managers (Coming Soon)
```bash
# Homebrew (macOS/Linux)
brew install codegenius

# Chocolatey (Windows)
choco install codegenius

# APT (Ubuntu/Debian)
sudo apt install codegenius

# YUM (CentOS/RHEL)
sudo yum install codegenius
```

## 🎯 Usage After Global Installation

Once installed globally, you can use CodeGenius from **any directory** on your system:

### Quick Start
```bash
# Navigate to any Git repository
cd ~/my-project

# Initialize CodeGenius (creates .codegenius.yaml)
codegenius --init

# Stage your changes
git add .

# Launch beautiful TUI
codegenius --tui

# Or use CLI commands
codegenius                # Generate commit message
codegenius --review       # Perform code review
codegenius --history      # View work history
```

### Pro Tips
```bash
# Create convenient aliases
echo 'alias cg="codegenius --tui"' >> ~/.zshrc
echo 'alias cgr="codegenius --review"' >> ~/.zshrc

# Quick workflow
git add . && cg          # Stage and commit
git add . && cgr         # Stage and review
```

## 🔧 For Developers: How to Distribute

### Setting Up Your Repository

1. **Update go.mod with public module path**:
```go
module github.com/yourusername/codegenius-cli
```

2. **Push to GitHub**:
```bash
git remote add origin https://github.com/yourusername/codegenius-cli.git
git push -u origin main
```

3. **Create a release**:
```bash
git tag v1.0.0
git push origin v1.0.0
```

### GitHub Actions for Automatic Releases

The included `.github/workflows/release.yml` automatically:
- Builds for all platforms (Linux, macOS, Windows)
- Creates GitHub releases with binaries
- Generates release notes
- Updates package manager formulas

### Manual Release Process

```bash
# Build for all platforms
make build-all

# Create distribution packages
make dist

# Create a complete release
make release

# Upload to GitHub releases manually or use:
gh release create v1.0.0 dist/* --generate-notes
```

## 📦 Distribution Channels

### 1. GitHub Releases
- Automatic via GitHub Actions
- Manual uploads via `gh` CLI or web interface
- Supports all platforms

### 2. Go Modules
- Automatic when you push tags
- Users install with `go install github.com/yourusername/codegenius-cli/cmd@latest`
- No additional setup required

### 3. Package Managers

#### Homebrew (macOS/Linux)
Create a Homebrew formula:
```ruby
class Codegenius < Formula
  desc "AI-powered Git commit message generator and code reviewer"
  homepage "https://github.com/yourusername/codegenius-cli"
  url "https://github.com/yourusername/codegenius-cli/archive/v1.0.0.tar.gz"
  sha256 "..."
  
  depends_on "go" => :build
  
  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd"
  end
end
```

#### Chocolatey (Windows)
Create a Chocolatey package:
```xml
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>codegenius</id>
    <version>1.0.0</version>
    <title>CodeGenius CLI</title>
    <authors>Your Name</authors>
    <description>AI-powered Git commit message generator</description>
    <projectUrl>https://github.com/yourusername/codegenius-cli</projectUrl>
  </metadata>
</package>
```

## 🌍 Global Usage Examples

### Individual Developers
```bash
# Work seamlessly across multiple projects
cd ~/work/project1 && codegenius --tui
cd ~/personal/project2 && codegenius --review
cd ~/opensource/project3 && codegenius --history
```

### Team Adoption
```bash
# Team lead shares configuration
cd team-project
codegenius --init
git add .codegenius.yaml
git commit -m "Add CodeGenius team configuration"

# Team members use globally
git pull
codegenius --tui  # Uses team configuration
```

### CI/CD Integration
```yaml
# GitHub Actions
- name: Code Review
  run: |
    codegenius --review > review-report.txt
    
# GitLab CI
script:
  - codegenius --review --output json > review.json
```

## 🔒 Security Considerations

### API Key Management
```bash
# Set globally (recommended)
export GEMINI_API_KEY="your-api-key"
echo 'export GEMINI_API_KEY="your-api-key"' >> ~/.zshrc

# Or per-project
cd project && echo "GEMINI_API_KEY=your-key" > .env
```

### Data Privacy
- Code never leaves your machine (except for AI analysis)
- Work history stored locally in `~/.codegenius/`
- Configuration files are project-specific
- All API calls use HTTPS

## 📊 Analytics & Telemetry (Optional)

To help improve CodeGenius, you can enable anonymous usage analytics:

```yaml
# .codegenius.yaml
analytics:
  enabled: true
  anonymous: true
  endpoint: "https://analytics.codegenius.dev"
```

## 🆘 Support & Troubleshooting

### Common Issues

1. **Command not found**: Ensure `/usr/local/bin` is in your PATH
2. **Permission denied**: Run `chmod +x /usr/local/bin/codegenius`
3. **API errors**: Verify your `GEMINI_API_KEY` is set correctly
4. **Git errors**: Ensure you're in a Git repository

### Getting Help
```bash
codegenius --help       # Built-in help
codegenius --version    # Version information
```

### Community Resources
- 📖 Documentation: https://github.com/yourusername/codegenius-cli#readme
- 🐛 Issues: https://github.com/yourusername/codegenius-cli/issues
- 💬 Discussions: https://github.com/yourusername/codegenius-cli/discussions
- 📧 Email: support@codegenius.dev

## 🎉 Success Stories

> "CodeGenius has transformed our team's commit quality. Now everyone writes meaningful commit messages!" - @developer1

> "I use it across 20+ repositories. The global installation saves me so much time." - @opensourcedev

> "The AI reviews caught security issues we missed in code review." - @teamlead

---

**Ready to make CodeGenius globally available? Follow this guide and let developers around the world improve their Git workflow! 🚀** 