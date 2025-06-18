# 🚀 CodeGenius Release Guide

## Quick Release Process

### 1. **Authenticate with GitHub**
```bash
# First time only
gh auth login
# Follow the prompts to authenticate
```

### 2. **Create Release with Binaries**
```bash
# Run the automated release script
./scripts/create-release.sh
```

This will:
- ✅ Check all binaries exist in `dist/`
- ✅ Create GitHub release v1.1.2
- ✅ Upload all platform binaries
- ✅ Use release notes from `RELEASE_NOTES.md`

### 3. **Test Homebrew Installation**
```bash
# Direct formula install (should work after release)
brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb

# Test it works
codegenius --help
```

### 4. **Set Up Homebrew Tap (Optional)**
```bash
# Create the official tap
./scripts/setup-homebrew-tap.sh
```

## Manual Release Steps (if needed)

### **Step 1: Build Binaries**
```bash
make build-all
```

### **Step 2: Create Release Manually**
```bash
gh release create v1.1.2 \
  --title "CodeGenius CLI v1.1.2 - Homebrew Ready 🍺" \
  --notes-file RELEASE_NOTES.md \
  dist/codegenius-darwin-amd64 \
  dist/codegenius-darwin-arm64 \
  dist/codegenius-linux-amd64 \
  dist/codegenius-linux-arm64 \
  dist/codegenius-windows-amd64.exe
```

## Installation Methods After Release

### **🍺 Homebrew Direct Formula**
```bash
brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb
```

### **🍺 Homebrew via Tap**
```bash
brew tap Shubhpreet-Rana/codegenius
brew install codegenius
```

### **🌐 Curl Installer (Always Works)**
```bash
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/install.sh | bash
```

### **🔧 Go Install**
```bash
go install github.com/Shubhpreet-Rana/codegenius@latest
```

### **📦 NPM Package**
```bash
npm install -g codegenius-cli
```

## Verification Steps

### **1. Check Release**
```bash
# Verify release exists
gh release view v1.1.2

# Check binaries are uploaded
gh release view v1.1.2 --json assets -q '.assets[].name'
```

### **2. Test Homebrew Formula**
```bash
# Test SHA256 hashes match
shasum -a 256 dist/codegenius-darwin-amd64
# Should match the hash in Formula/codegenius.rb
```

### **3. Test Installation**
```bash
# Test direct formula
brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb

# Verify it works
codegenius --help
codegenius --tui

# Clean up test
brew uninstall codegenius
```

## Troubleshooting

### **"404 Not Found" Error**
- Make sure release v1.1.2 exists: `gh release view v1.1.2`
- Check binaries are uploaded to the release

### **"SHA256 Mismatch" Error**
- Regenerate hashes: `shasum -a 256 dist/*`
- Update `Formula/codegenius.rb` with new hashes

### **"Formula URL Not Found"**
- Ensure formula is pushed to `latest` branch
- Check URL in browser: https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb

## Current Status

- ✅ **Binaries Built**: All platforms ready in `dist/`
- ✅ **Formula Ready**: SHA256 hashes verified
- ✅ **Scripts Ready**: Automated release and tap setup
- 🔄 **Release Pending**: Run `./scripts/create-release.sh`
- 🔄 **Tap Pending**: Run `./scripts/setup-homebrew-tap.sh`

## Next Steps

1. **Authenticate**: `gh auth login`
2. **Create Release**: `./scripts/create-release.sh`
3. **Test Homebrew**: `brew install --formula https://...`
4. **Create Tap**: `./scripts/setup-homebrew-tap.sh`
5. **Update Documentation**: Mark Homebrew as fully available

After these steps, CodeGenius will be installable via:
- ✅ Curl installer
- ✅ Go install  
- ✅ NPM package
- ✅ Homebrew direct formula
- ✅ Homebrew tap

🎉 **Full deployment complete!** 