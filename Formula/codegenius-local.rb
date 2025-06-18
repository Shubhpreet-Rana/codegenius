class CodegeniusLocal < Formula
  desc "AI-powered Git commit message generator and code reviewer with beautiful TUI"
  homepage "https://github.com/Shubhpreet-Rana/codegenius"
  version "1.1.2"
  license "MIT"

  def install
    # Determine the correct binary for current platform
    if OS.mac?
      if Hardware::CPU.arm?
        binary_name = "codegenius-darwin-arm64"
      else
        binary_name = "codegenius-darwin-amd64"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        binary_name = "codegenius-linux-arm64"
      else
        binary_name = "codegenius-linux-amd64"
      end
    end

    # Path to the binary in our dist directory
    binary_path = "#{File.dirname(__FILE__)}/../dist/#{binary_name}"
    
    if File.exist?(binary_path)
      bin.install binary_path => "codegenius"
    else
      odie "Binary not found at #{binary_path}. Please run 'make build-all' first."
    end
  end

  def caveats
    <<~EOS
      🤖 CodeGenius CLI installed successfully via Homebrew!

      ✅ CodeGenius is ready to use globally!
      🚀 Try it now: codegenius --tui

      📋 Next steps:
      1. Get your Gemini API key: https://makersuite.google.com/app/apikey
      2. Set environment variable:
         export GEMINI_API_KEY="your-api-key-here"
      3. Add to your shell profile (~/.zshrc, ~/.bashrc, etc.)
      4. Initialize in any Git repository:
         cd your-project && codegenius --init
      5. Start using CodeGenius:
         codegenius --tui

      💡 Quick commands:
         codegenius --help     # Show help
         codegenius --tui      # Beautiful interface  
         codegenius --review   # Code review

      🔗 Documentation: https://github.com/Shubhpreet-Rana/codegenius#readme
    EOS
  end

  test do
    # Test that the binary exists and is executable
    assert_predicate bin/"codegenius", :exist?
    assert_predicate bin/"codegenius", :executable?
    
    # Test help command (doesn't require Git repo or API key)
    output = shell_output("#{bin}/codegenius --help")
    assert_match "CodeGenius CLI", output
  end
end 