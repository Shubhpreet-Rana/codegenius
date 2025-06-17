#!/usr/bin/env node

/**
 * CodeGenius CLI - NPM Pre-Uninstall Script
 * Cleans up when the package is removed
 */

const fs = require('fs');
const path = require('path');

// Colors for output
const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    cyan: '\x1b[36m',
};

function colorize(text, color) {
    return `${colors[color]}${text}${colors.reset}`;
}

function log(message, color = 'reset') {
    console.log(colorize(message, color));
}

// Remove lib directory and binaries
function cleanup() {
    try {
        log('🧹 CodeGenius CLI - Cleanup', 'cyan');
        log('==========================', 'cyan');
        
        const libDir = path.join(__dirname, '..', 'lib');
        
        if (fs.existsSync(libDir)) {
            log('🗑️  Removing downloaded binaries...', 'blue');
            
            // Recursively remove lib directory
            function removeDir(dirPath) {
                if (fs.existsSync(dirPath)) {
                    fs.readdirSync(dirPath).forEach((file) => {
                        const curPath = path.join(dirPath, file);
                        if (fs.lstatSync(curPath).isDirectory()) {
                            removeDir(curPath);
                        } else {
                            fs.unlinkSync(curPath);
                        }
                    });
                    fs.rmdirSync(dirPath);
                }
            }
            
            removeDir(libDir);
            log('✅ Cleanup complete!', 'green');
        } else {
            log('✅ No cleanup needed', 'yellow');
        }
        
        log('', 'reset');
        log('👋 Thank you for using CodeGenius CLI!', 'cyan');
        log('', 'reset');
        
    } catch (error) {
        log('⚠️  Cleanup failed (this is usually not a problem)', 'yellow');
        log(`Error: ${error.message}`, 'yellow');
    }
}

// Skip cleanup in CI environments
if (process.env.CI || process.env.NODE_ENV === 'test') {
    log('⏭️  Skipping cleanup in CI environment', 'yellow');
    process.exit(0);
}

// Run cleanup
cleanup(); 