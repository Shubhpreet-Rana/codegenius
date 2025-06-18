#!/usr/bin/env node

/**
 * CodeGenius CLI - NPM Wrapper
 * This script launches the native CodeGenius binary
 */

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');
const os = require('os');

// Get the path to the binary
function getBinaryPath() {
    const platform = os.platform();
    
    // Use the single unified binary name
    const binaryName = platform === 'win32' ? 'codegenius-windows-amd64.exe' : 'codegenius';
    const binaryPath = path.join(__dirname, binaryName);
    
    if (!fs.existsSync(binaryPath)) {
        console.error(`❌ CodeGenius binary not found at: ${binaryPath}`);
        console.error('Try reinstalling: npm uninstall -g codegenius-cli && npm install -g codegenius-cli');
        process.exit(1);
    }
    
    return binaryPath;
}

// Launch the binary
function launchCodeGenius() {
    const binaryPath = getBinaryPath();
    const args = process.argv.slice(2);
    
    // Spawn the binary with the same arguments
    const child = spawn(binaryPath, args, {
        stdio: 'inherit',
        windowsHide: false,
    });
    
    child.on('error', (error) => {
        console.error(`❌ Failed to start CodeGenius: ${error.message}`);
        process.exit(1);
    });
    
    child.on('exit', (code, signal) => {
        if (signal) {
            process.kill(process.pid, signal);
        } else {
            process.exit(code || 0);
        }
    });
    
    // Handle process termination
    process.on('SIGINT', () => {
        child.kill('SIGINT');
    });
    
    process.on('SIGTERM', () => {
        child.kill('SIGTERM');
    });
}

// Check if this is a help request and show wrapper info
if (process.argv.includes('--help') || process.argv.includes('-h')) {
    console.log('🤖 CodeGenius CLI - NPM Wrapper');
    console.log('');
    console.log('This is a Node.js wrapper for the CodeGenius binary.');
    console.log('All arguments are passed through to the native CodeGenius CLI.');
    console.log('');
    console.log('Examples:');
    console.log('  codegenius --tui       # Launch beautiful terminal UI');
    console.log('  codegenius --review    # Perform code review');
    console.log('  codegenius --help      # Show CodeGenius help');
    console.log('');
    console.log('For full documentation, visit:');
    console.log('https://github.com/Shubhpreet-Rana/codegenius#readme');
    console.log('');
}

// Launch CodeGenius
launchCodeGenius(); 