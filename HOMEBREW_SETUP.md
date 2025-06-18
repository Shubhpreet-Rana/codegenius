# 🍺 Homebrew Setup for CodeGenius CLI

## **Current Status**
The Homebrew formula exists but isn't published to the main Homebrew repository yet. Here are the installation options:

## **Option 1: Install via Local Formula (Immediate)**

You can install directly from the repository using the local formula:

```bash
# Install directly from the repository
brew install --formula https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/Formula/codegenius.rb
```

**Note**: This requires a GitHub release with pre-built binaries.

## **Option 2: Create a Homebrew Tap (Recommended)**

### **Step 1: Create a Homebrew Tap Repository**

1. Create a new repository named `homebrew-codegenius` on GitHub
2. Clone it locally:
   ```bash
   git clone https://github.com/Shubhpreet-Rana/homebrew-codegenius.git
   cd homebrew-codegenius
   ```

3. Copy the formula:
   ```bash
   mkdir -p Formula
   cp /path/to/codegenius/Formula/codegenius.rb Formula/
   git add Formula/codegenius.rb
   git commit -m "Add CodeGenius formula"
   git push origin main
   ```

### **Step 2: Users Can Install via Tap**
```bash
# Add the tap
brew tap Shubhpreet-Rana/codegenius

# Install CodeGenius
brew install codegenius
```

## **Option 3: Alternative Installation Methods (Working Now)**

Since Homebrew setup requires additional steps, users can use these working methods:

### **🚀 Curl Installer (Recommended)**
```bash
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/install.sh | bash
```

### **📦 Go Install**
```bash
go install github.com/Shubhpreet-Rana/codegenius@latest
```

### **📱 NPM Package**
```bash
npm install -g codegenius-cli
```

## **Required Steps for Homebrew to Work**

### **1. Create GitHub Release with Binaries**
```bash
# Build binaries for all platforms
make build-all

# Create release on GitHub with binaries
# Upload: codegenius-darwin-amd64, codegenius-darwin-arm64, etc.
```

### **2. Update SHA256 Hashes**
After creating the release, update the formula with actual SHA256 hashes:
```bash
# Get SHA256 for each binary
shasum -a 256 dist/codegenius-darwin-amd64
shasum -a 256 dist/codegenius-darwin-arm64
# ... etc for all platforms
```

Update `Formula/codegenius.rb` with the real hashes.

## **Quick Test of Current Status**

To test what happens with current Homebrew command:
```bash
# This will fail with "No available formula" error
brew install codegenius

# This is expected until we complete the setup above
```

## **Recommended Immediate Action**

For now, update your documentation to recommend the **curl installer** as the primary method:

```bash
# Primary installation method (works immediately)
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/latest/install.sh | bash

# Alternative: Go install (for developers)
go install github.com/Shubhpreet-Rana/codegenius@latest
```

The Homebrew option can be added later once the tap is set up properly. 