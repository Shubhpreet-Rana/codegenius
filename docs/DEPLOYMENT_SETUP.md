# 🚀 Automated Deployment Pipeline Setup

This guide explains how to set up the automated CI/CD pipeline that deploys CodeGenius CLI to all distribution channels when pushing to the `latest` branch.

## 🎯 **Pipeline Overview**

The pipeline automatically:
1. **Calculates semantic version** (patch/minor/major)
2. **Builds binaries** for all platforms
3. **Creates GitHub release** with binaries and checksums
4. **Updates NPM package** and publishes
5. **Updates Homebrew formula** in tap repository
6. **Updates curl installer** script
7. **Provides deployment summary**

## 🔧 **Required GitHub Secrets**

You need to configure these secrets in your GitHub repository:

### 1. NPM Publishing (`NPM_TOKEN`)
```bash
# Create NPM access token
npm login
npm token create --type=automation

# Add to GitHub Secrets:
# Name: NPM_TOKEN
# Value: npm_1a2b3c4d5e6f7g8h9i0j...
```

### 2. Homebrew Tap Access (`TAP_GITHUB_TOKEN`)
```bash
# Create GitHub Personal Access Token with repo permissions
# Go to: https://github.com/settings/tokens
# Select: repo (Full control of private repositories)

# Add to GitHub Secrets:
# Name: TAP_GITHUB_TOKEN  
# Value: ghp_1a2b3c4d5e6f7g8h9i0j...
```

### 3. Default GitHub Token
The `GITHUB_TOKEN` is automatically provided by GitHub Actions.

## 📋 **Setup Steps**

### Step 1: Configure Repository Secrets

1. Go to your repository settings
2. Navigate to **Secrets and variables > Actions**
3. Add the following secrets:

| Secret Name | Description | How to Get |
|-------------|-------------|------------|
| `NPM_TOKEN` | NPM automation token | `npm token create --type=automation` |
| `TAP_GITHUB_TOKEN` | GitHub token for tap repo | GitHub Settings > Developer settings > Personal access tokens |

### Step 2: Ensure Repository Structure

```
your-repo/
├── .github/workflows/deploy.yml  # ✅ Created
├── Formula/codegenius.rb         # ✅ Exists
├── install.sh                    # ✅ Exists  
├── package.json                  # ✅ Exists
├── bin/codegenius.js            # ✅ Exists
└── main.go                      # ✅ Exists
```

### Step 3: Verify Homebrew Tap Repository

Ensure you have a separate repository: `homebrew-codegenius`
```bash
# Create if it doesn't exist
gh repo create homebrew-codegenius --public
cd homebrew-codegenius
mkdir Formula
# Copy your current formula
cp ../CLI_GO/Formula/codegenius.rb Formula/
git add . && git commit -m "Initial tap setup"
git push
```

## 🚀 **How to Trigger Deployment**

### Automatic Deployment (Patch Version)
```bash
# Push to latest branch triggers patch version bump (1.0.0 → 1.0.1)
git checkout latest
git push origin latest
```

### Manual Deployment with Version Control
```bash
# Trigger manually with specific version bump
# Go to: https://github.com/your-username/codegenius/actions
# Click "🚀 Auto Deploy to All Channels"
# Click "Run workflow"
# Select version bump: patch | minor | major
```

## 🔢 **Version Bumping Logic**

| Current Version | Bump Type | New Version |
|----------------|-----------|-------------|
| `1.0.0` | patch | `1.0.1` |
| `1.0.9` | patch | `1.0.10` |
| `1.0.0` | minor | `1.1.0` |
| `1.9.0` | minor | `1.10.0` |
| `1.0.0` | major | `2.0.0` |
| `1.9.9` | major | `2.0.0` |

## 📦 **What Gets Updated**

### 1. GitHub Release
- Creates new release with calculated version tag
- Uploads binaries for all platforms
- Generates SHA256 checksums
- Auto-generated release notes

### 2. NPM Package
- Updates `package.json` version
- Includes latest binary in package
- Publishes to NPM registry as `codegenius-cli`

### 3. Homebrew Formula
- Updates main repository `Formula/codegenius.rb`
- Updates tap repository `Shubhpreet-Rana/homebrew-codegenius`
- New SHA256 hashes for all platforms
- Version bump in formula

### 4. Curl Installer
- Updates version reference in `install.sh`
- Commits changes to main repository

## 🔍 **Monitoring Deployment**

### View Pipeline Status
```bash
# Watch the pipeline
gh run list --workflow="deploy.yml"
gh run watch

# View specific run
gh run view <run-id>
```

### Verify Deployments
```bash
# Check NPM
npm view codegenius-cli version

# Check Homebrew  
brew search codegenius

# Check GitHub releases
gh release list

# Test installation
npm install -g codegenius-cli@latest
```

## 🛠️ **Troubleshooting**

### Common Issues

#### NPM Token Expired
```bash
# Generate new token
npm token create --type=automation
# Update GitHub secret: NPM_TOKEN
```

#### Homebrew Push Failed
```bash
# Check TAP_GITHUB_TOKEN permissions
# Ensure token has 'repo' scope
# Verify homebrew-codegenius repository exists
```

#### Binary Build Failed
```bash
# Check Go version in workflow (currently 1.21)
# Verify main.go compiles locally
go build .
```

#### Version Calculation Error
```bash
# Check if git tags exist
git tag -l
# Create initial tag if needed
git tag v1.0.0
git push --tags
```

### Debug Commands
```bash
# Test version calculation locally
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.0")
echo "Latest: $LATEST_TAG"

# Test binary build locally  
GOOS=linux GOARCH=amd64 go build -o test-binary .
file test-binary

# Test SHA256 generation
sha256sum test-binary
```

## 🎯 **Best Practices**

### Release Strategy
1. **Development**: Work on feature branches
2. **Testing**: Merge to `main` for testing
3. **Release**: Push to `latest` for deployment
4. **Hotfixes**: Use manual trigger with patch bump

### Version Management
- Use **patch** for bug fixes and small improvements
- Use **minor** for new features
- Use **major** for breaking changes
- Follow [Semantic Versioning](https://semver.org/)

### Security
- Rotate tokens regularly
- Use least-privilege access
- Monitor deployment logs
- Verify published packages

## 📈 **Pipeline Performance**

Typical deployment times:
- **Version Calculation**: ~30 seconds
- **Binary Building**: ~2-3 minutes
- **GitHub Release**: ~1 minute
- **NPM Publishing**: ~1 minute
- **Homebrew Update**: ~1 minute
- **Total**: ~6-8 minutes

## 🎉 **Success Indicators**

After successful deployment, verify:

1. ✅ **GitHub Release**: New tag with binaries
2. ✅ **NPM**: `npm view codegenius-cli version`
3. ✅ **Homebrew**: Updated formula with new SHA256s
4. ✅ **Install Test**: All methods work

```bash
# Quick verification
brew install codegenius
npm install -g codegenius-cli
curl -fsSL https://raw.githubusercontent.com/.../install.sh | bash
```

---

**🚀 Your automated deployment pipeline is ready! Push to `latest` branch to trigger deployment to all channels.** 