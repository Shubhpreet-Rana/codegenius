# 🤖 CodeGenius Setup for Friends

Hey! Your friend shared CodeGenius with you - an AI-powered Git assistant that writes commit messages and reviews code. Here are **multiple easy ways** to install it:

## 🚀 Super Easy Installation (Pick One!)

### 🌐 Method 1: One-Line Install (Fastest)
```bash
curl -fsSL https://raw.githubusercontent.com/Shubhpreet-Rana/codegenius/main/install.sh | bash
```
**✅ Automatic everything - detects your platform, downloads, installs globally**

### 📦 Method 2: NPM Install (If you have Node.js)
```bash
npm install -g codegenius-cli
```
**✅ Easy to update later with `npm update -g codegenius-cli`**

### 💾 Method 3: Direct Binary (Manual)
Your friend should send you the right file:

- **macOS Intel**: `codegenius-darwin-amd64`
- **macOS Apple Silicon**: `codegenius-darwin-arm64` 
- **Linux x86_64**: `codegenius-linux-amd64`
- **Linux ARM64**: `codegenius-linux-arm64`
- **Windows**: `codegenius-windows-amd64.exe`

Then install:
```bash
# macOS/Linux: Make it executable and move to PATH
chmod +x codegenius-*
sudo mv codegenius-* /usr/local/bin/codegenius

# Windows: Move to a folder in your PATH or create one
```

### 🔧 Method 4: Build from Source (If you have Go)
```bash
# Get the source code from your friend or clone their repo
cd codegenius-project
go build -o codegenius cmd/main.go
sudo mv codegenius /usr/local/bin/  # macOS/Linux
```

---

## 🔑 Step 2: Get Your AI API Key

1. **Go to Google AI Studio**: https://makersuite.google.com/app/apikey
2. **Sign in** with your Google account
3. **Create API Key** - click "Create API Key"
4. **Copy the key** (looks like: `AIzaSyC7x...`)

## ⚙️ Step 3: Set Up Your Environment

**Add your API key to your shell profile:**

```bash
# For macOS/Linux (choose your shell)
echo 'export GEMINI_API_KEY="your-api-key-here"' >> ~/.zshrc     # Zsh
echo 'export GEMINI_API_KEY="your-api-key-here"' >> ~/.bashrc   # Bash

# Reload your shell
source ~/.zshrc  # or ~/.bashrc

# For Windows (Command Prompt)
setx GEMINI_API_KEY "your-api-key-here"

# For Windows (PowerShell)
$env:GEMINI_API_KEY = "your-api-key-here"
[Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "your-api-key-here", "User")
```

## 🎯 Step 4: Test Installation

```bash
# Check if it's working
codegenius --help

# You should see the help message with a robot emoji 🤖
```

## 🚀 Step 5: Start Using It!

### **In ANY Git Repository:**

```bash
# Navigate to your Git project
cd your-project

# Initialize CodeGenius (one-time setup per project)
codegenius --init

# Make some changes to your code
# Stage your changes
git add .

# Use the beautiful TUI (recommended)
codegenius --tui

# OR use traditional commands
codegenius          # Generate commit message
codegenius --review # Review your code
```

## 💡 Pro Tips

### **Create Aliases (Optional but Recommended)**
```bash
# Add these to ~/.zshrc or ~/.bashrc
alias cg="codegenius --tui"
alias cgr="codegenius --review"

# Then you can just type:
git add . && cg     # Quick commit
git add . && cgr    # Quick review
```

### **Quick Workflow**
1. Make code changes
2. `git add .`
3. `cg` (opens beautiful TUI)
4. Choose "Generate commit message"
5. Add context if needed
6. Commit!

## 🔍 What CodeGenius Can Do

### **🤖 Smart Commit Messages**
- Analyzes your code changes
- Generates meaningful commit messages
- Follows conventional commit format
- Considers your branch name and context

### **🔍 Code Reviews**
- **Security review**: Finds vulnerabilities
- **Performance review**: Identifies bottlenecks  
- **Style review**: Checks formatting and conventions
- **Structure review**: Evaluates architecture

### **📊 Work History**
- Tracks all your commits
- Shows monthly statistics
- Filters by date ranges

## 🆘 Troubleshooting

### **"Command not found"**
```bash
# Check if it's in your PATH
which codegenius

# If not found, make sure /usr/local/bin is in your PATH
echo $PATH

# Add to PATH if needed (add to ~/.zshrc or ~/.bashrc)
export PATH="/usr/local/bin:$PATH"
```

### **"API key not set"**
```bash
# Check if your API key is set
echo $GEMINI_API_KEY

# If empty, set it again and restart your terminal
```

### **"No Git repository"**
```bash
# Make sure you're in a Git repository
git status

# If not, initialize one
git init
```

### **Update CodeGenius**
```bash
# If installed via curl
curl -fsSL https://raw.githubusercontent.com/yourusername/codegenius/main/install.sh | bash

# If installed via NPM
npm update -g codegenius-cli

# If installed via Go
go install github.com/yourusername/codegenius/cmd@latest
```

## 🎉 You're All Set!

CodeGenius is now ready to:
- ✅ **Generate intelligent commit messages**
- ✅ **Review your code for issues**
- ✅ **Track your development progress**
- ✅ **Work across all your projects**

**Need help?** Ask your friend or check `codegenius --help`

---

**Happy coding with your new AI-powered Git assistant! 🚀** 