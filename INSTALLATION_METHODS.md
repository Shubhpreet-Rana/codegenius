# 🚀 CodeGenius Installation Methods

Multiple ways to install CodeGenius CLI for maximum convenience!

## 🌟 Quick Install (Recommended)

### 🌐 One-Line Curl Install
```bash
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/main/install.sh | bash
```

### 📦 NPM Install
```bash
npm install -g codegenius-cli
```

### 🍺 Homebrew Install
```bash
brew tap yourusername/codegenius
brew install codegenius
```

---

## 📋 Setup After Installation

**All methods require the same setup:**

1. **Get Gemini API Key:** https://makersuite.google.com/app/apikey
2. **Set Environment Variable:**
   ```bash
   export GEMINI_API_KEY="your-api-key-here"
   ```
3. **Add to Shell Profile:** (~/.zshrc, ~/.bashrc, etc.)
4. **Verify Installation:**
   ```bash
   codegenius --help
   ```
5. **Start Using:**
   ```bash
   cd your-git-project
   codegenius --init
   codegenius --tui
   ```

---

## 🎯 Platform Support

| Platform | Curl | NPM | Homebrew | 
|----------|------|-----|----------|
| **macOS Intel** | ✅ | ✅ | ✅ |
| **macOS ARM64** | ✅ | ✅ | ✅ |
| **Linux x64** | ✅ | ✅ | ✅ |
| **Linux ARM64** | ✅ | ✅ | ✅ |
| **Windows** | ✅ | ✅ | ❌ |

---

## 🆘 Need Help?

- 📖 **Documentation:** [README.md](README.md)
- 🐛 **Issues:** [GitHub Issues](https://github.com/yourusername/codegenius/issues)

**Choose the method that works best for you and start using CodeGenius in seconds! 🚀** 