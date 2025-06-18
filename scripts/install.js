#!/usr/bin/env node

/**
 * CodeGenius CLI - NPM Post-Install Script
 * Downloads the appropriate binary for the user's platform
 */

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { promisify } = require('util');
const { pipeline } = require('stream');

const streamPipeline = promisify(pipeline);

// Configuration
const REPO_OWNER = 'Shubhpreet-Rana';
const REPO_NAME = 'codegenius';
const GITHUB_API = `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}`;

// Colors for output
const colors = {
    reset: '\x1b[0m',
    bright: '\x1b[1m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m',
    cyan: '\x1b[36m',
};

function colorize(text, color) {
    return `${colors[color]}${text}${colors.reset}`;
}

function log(message, color = 'reset') {
    console.log(colorize(message, color));
}

// Determine platform and architecture
function getPlatform() {
    const platform = os.platform();
    const arch = os.arch();
    
    let platformName;
    switch (platform) {
        case 'darwin':
            platformName = 'darwin';
            break;
        case 'linux':
            platformName = 'linux';
            break;
        case 'win32':
            platformName = 'windows';
            break;
        default:
            throw new Error(`Unsupported platform: ${platform}`);
    }
    
    let archName;
    switch (arch) {
        case 'x64':
            archName = 'amd64';
            break;
        case 'arm64':
            archName = 'arm64';
            break;
        default:
            throw new Error(`Unsupported architecture: ${arch}`);
    }
    
    return { platform: platformName, arch: archName };
}

// Get latest release information
async function getLatestRelease() {
    return new Promise((resolve, reject) => {
        const url = `${GITHUB_API}/releases/latest`;
        
        https.get(url, {
            headers: {
                'User-Agent': 'codegenius-npm-installer',
            },
        }, (response) => {
            let data = '';
            
            response.on('data', (chunk) => {
                data += chunk;
            });
            
            response.on('end', () => {
                try {
                    const release = JSON.parse(data);
                    resolve(release);
                } catch (error) {
                    reject(new Error(`Failed to parse release data: ${error.message}`));
                }
            });
        }).on('error', (error) => {
            reject(new Error(`Failed to fetch release info: ${error.message}`));
        });
    });
}

// Download binary
async function downloadBinary(url, outputPath) {
    return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(outputPath);
        
        https.get(url, (response) => {
            if (response.statusCode !== 200) {
                reject(new Error(`Download failed with status ${response.statusCode}`));
                return;
            }
            
            const totalSize = parseInt(response.headers['content-length'], 10);
            let downloadedSize = 0;
            
            response.on('data', (chunk) => {
                downloadedSize += chunk.length;
                const progress = ((downloadedSize / totalSize) * 100).toFixed(1);
                process.stdout.write(`\r${colorize('Downloading...', 'blue')} ${progress}%`);
            });
            
            response.pipe(file);
            
            file.on('finish', () => {
                file.close();
                console.log(''); // New line after progress
                resolve();
            });
            
        }).on('error', (error) => {
            fs.unlink(outputPath, () => {}); // Delete partial file
            reject(new Error(`Download failed: ${error.message}`));
        });
    });
}

// Main installation function
async function install() {
    try {
        log('🤖 CodeGenius CLI - Post-Install Setup', 'cyan');
        log('=====================================', 'cyan');
        
        // Get platform info
        const { platform, arch } = getPlatform();
        log(`📋 Detected platform: ${platform}-${arch}`, 'blue');
        
        // Create lib directory
        const libDir = path.join(__dirname, '..', 'lib');
        const binDir = path.join(libDir, 'bin');
        
        if (!fs.existsSync(libDir)) {
            fs.mkdirSync(libDir, { recursive: true });
        }
        
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true });
        }
        
        // Get latest release
        log('🔍 Fetching latest release...', 'blue');
        const release = await getLatestRelease();
        
        if (!release || !release.tag_name) {
            throw new Error('Could not determine latest version');
        }
        
        log(`📦 Latest version: ${release.tag_name}`, 'green');
        
        // Construct binary name and download URL
        const extension = platform === 'windows' ? '.exe' : '';
        const binaryName = `codegenius-${platform}-${arch}${extension}`;
        const downloadUrl = `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${release.tag_name}/${binaryName}`;
        const outputPath = path.join(binDir, binaryName);
        
        // Check if binary already exists
        if (fs.existsSync(outputPath)) {
            log('✅ Binary already exists, skipping download', 'yellow');
        } else {
            // Download binary
            log(`⬇️  Downloading ${binaryName}...`, 'blue');
            await downloadBinary(downloadUrl, outputPath);
            log('✅ Download complete!', 'green');
        }
        
        // Make binary executable (Unix systems)
        if (platform !== 'windows') {
            fs.chmodSync(outputPath, '755');
        }
        
        log('', 'reset');
        log('🎉 CodeGenius CLI installed successfully!', 'green');
        log('', 'reset');
        log('📋 Next steps:', 'cyan');
        log('1. Get your Gemini API key: https://makersuite.google.com/app/apikey', 'reset');
        log('2. Set environment variable: export GEMINI_API_KEY="your-key"', 'yellow');
        log('3. Start using: codegenius --tui', 'yellow');
        log('', 'reset');
        log('💡 Need help? Run: codegenius --help', 'blue');
        
    } catch (error) {
        log('', 'reset');
        log('❌ Installation failed!', 'red');
        log(`Error: ${error.message}`, 'red');
        log('', 'reset');
        log('🔧 Troubleshooting:', 'yellow');
        log('1. Check your internet connection', 'reset');
        log('2. Try again: npm uninstall -g codegenius-cli && npm install -g codegenius-cli', 'reset');
        log('3. Report issues: https://github.com/Shubhpreet-Rana/codegenius/issues', 'reset');
        
        process.exit(1);
    }
}

// Skip installation in CI environments
if (process.env.CI || process.env.NODE_ENV === 'test') {
    log('⏭️  Skipping binary download in CI environment', 'yellow');
    process.exit(0);
}

// Run installation
install(); 