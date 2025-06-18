#!/usr/bin/env node

/**
 * CodeGenius CLI - NPM Post-Install Script
 * Sets up the binary from the included package
 */

const fs = require('fs');
const path = require('path');
const os = require('os');

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

// Main installation function
async function install() {
    try {
        log('🤖 CodeGenius CLI - Post-Install Setup', 'cyan');
        log('=====================================', 'cyan');
        
        // Get platform info
        const { platform, arch } = getPlatform();
        log(`📋 Detected platform: ${platform}-${arch}`, 'blue');
        
        // Find the included binary
        const extension = platform === 'windows' ? '.exe' : '';
        const binaryName = platform === 'windows' ? 'codegenius-windows-amd64.exe' : 'codegenius';
        
        // Source binary path (included in package)
        const sourceBinaryPath = path.join(__dirname, '..', 'bin', binaryName);
        
        // Check if source binary exists
        if (!fs.existsSync(sourceBinaryPath)) {
            log(`❌ Binary not found at: ${sourceBinaryPath}`, 'red');
            log('This might be a packaging issue. Please report it.', 'yellow');
            process.exit(1);
        }
        
        log('✅ Found pre-built binary in package', 'green');
        
        // Make binary executable (Unix systems)
        if (platform !== 'windows') {
            fs.chmodSync(sourceBinaryPath, '755');
            log('✅ Binary permissions set', 'green');
        }
        
        log('', 'reset');
        log('🎉 CodeGenius CLI installed successfully!', 'green');
        log('', 'reset');
        log('📋 Next steps:', 'cyan');
        log('1. Get your Gemini API key: https://makersuite.google.com/app/apikey', 'reset');
        log('2. Set environment variable: export GEMINI_API_KEY="your-key"', 'yellow');
        log('3. Test installation: codegenius --help', 'blue');
        log('4. Start using: codegenius --tui', 'magenta');
        log('', 'reset');
        log('📚 Documentation: https://github.com/Shubhpreet-Rana/codegenius#readme', 'cyan');
        
    } catch (error) {
        log('', 'reset');
        log('❌ Installation failed!', 'red');
        log(`Error: ${error.message}`, 'red');
        log('', 'reset');
        log('🔧 Troubleshooting:', 'yellow');
        log('1. Check your platform is supported (macOS, Linux, Windows)', 'reset');
        log('2. Try: npm uninstall -g codegenius-cli && npm install -g codegenius-cli', 'reset');
        log('3. Report issues: https://github.com/Shubhpreet-Rana/codegenius/issues', 'reset');
        
        process.exit(1);
    }
}

// Run installation
install(); 