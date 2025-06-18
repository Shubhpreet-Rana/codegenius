# 🚀 CodeGenius CLI v1.1.2 Release Notes

## 🎯 **Homebrew Formula Ready!**

This release makes CodeGenius fully compatible with Homebrew installation. The formula is now updated with actual SHA256 hashes and ready for tap deployment.

## 🔧 **What's New**

### ✅ **Homebrew Support Complete**
- **Real SHA256 hashes**: Formula updated with actual binary hashes
- **Multi-platform support**: macOS (Intel/ARM) and Linux (x64/ARM64)
- **Improved testing**: Better formula validation and testing
- **Ready for tap**: Can now be deployed to Homebrew tap

### 🛠 **Enhanced Installation**
- **Four installation methods**: Curl, Go install, NPM, Manual download
- **Automatic PATH setup**: Installer detects shell and adds to PATH
- **Non-interactive mode**: Works with `curl | bash` pipelines
- **Global accessibility**: Install once, use anywhere

### 📦 **Built Binaries**
All platforms pre-built and tested:
- `codegenius-darwin-amd64` (macOS Intel)
- `codegenius-darwin-arm64` (macOS Apple Silicon)  
- `codegenius-linux-amd64` (Linux x86_64)
- `codegenius-linux-arm64` (Linux ARM64)
- `codegenius-windows-amd64.exe` (Windows x64)

## 📋 **Installation Methods**

### 🌐 **Primary: One-Line Install**
```bash
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/install.sh | bash
```

### 🔧 **Go Install**
```bash
go install github.com/Shubhpreet-Rana/codegenius@latest
```

### 📦 **NPM Package**
```bash
npm install -g codegenius-cli
```

### 🍺 **Homebrew (Ready for Deployment)**
```bash
# Option 1: Direct formula install
brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb

# Option 2: Via tap (once set up)
brew tap Shubhpreet-Rana/codegenius
brew install codegenius
```

## 🔐 **SHA256 Checksums**

For security verification:
```
6f73e4fa9c7c1f610ec9cc86acbb4effce926b257dc781fefde463262ef00047  codegenius-darwin-amd64
ed4c797036f42d028a2a733586be2e51f838544c30d3b9641c4c5210e0d4dd81  codegenius-darwin-arm64
692d20c1fe050799d10637863226c63a57585d5b4a2ffa4c9c9db6c62afd381b  codegenius-linux-amd64
05d62bf6c86bc7c2812e768de6141a69357cc992b6c0494086d9a9cfe3b4dc56  codegenius-linux-arm64
```

## 🎯 **Quick Start**
```bash
# Get API key from: https://makersuite.google.com/app/apikey
export GEMINI_API_KEY="your-gemini-api-key"

# Use anywhere
cd your-project
git add .
codegenius --tui
```

## 🛠 **For Homebrew Tap Setup**

To complete Homebrew deployment:

1. **Create tap repository**: `homebrew-codegenius`
2. **Copy formula**: `Formula/codegenius.rb` → tap repo
3. **Users install via**: `brew tap your-username/codegenius && brew install codegenius`

## 🐛 **Bug Fixes**
- Fixed installer input handling in non-interactive environments
- Improved PATH detection and setup across different shells
- Enhanced binary download and verification process

## 📈 **Improvements**
- Better error messages and user guidance
- Enhanced documentation with clear installation paths
- Streamlined release process with automated binary building

---

**🎉 Ready to revolutionize your Git workflow? Choose your preferred installation method and start using CodeGenius globally!** 