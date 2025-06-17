/**
 * CodeGenius CLI - NPM Package Entry Point
 * 
 * This package provides the CodeGenius CLI tool as an NPM package.
 * The actual CLI is implemented in Go and downloaded as a binary during installation.
 */

const path = require('path');
const { spawn } = require('child_process');

/**
 * Get the path to the CodeGenius binary
 * @returns {string} Path to the binary
 */
function getBinaryPath() {
    const os = require('os');
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
    
    const extension = platform === 'win32' ? '.exe' : '';
    const binaryName = `codegenius-${platformName}-${archName}${extension}`;
    
    return path.join(__dirname, 'bin', binaryName);
}

/**
 * Execute CodeGenius with the given arguments
 * @param {string[]} args - Command line arguments
 * @returns {Promise<number>} Exit code
 */
function executeCodeGenius(args = []) {
    return new Promise((resolve, reject) => {
        const binaryPath = getBinaryPath();
        
        const child = spawn(binaryPath, args, {
            stdio: 'inherit',
            windowsHide: false,
        });
        
        child.on('error', (error) => {
            reject(error);
        });
        
        child.on('exit', (code, signal) => {
            if (signal) {
                process.kill(process.pid, signal);
            } else {
                resolve(code || 0);
            }
        });
    });
}

module.exports = {
    getBinaryPath,
    executeCodeGenius,
    version: require('../package.json').version,
}; 