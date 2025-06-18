class Codegenius < Formula
  desc "AI-powered Git commit message generator and code reviewer with beautiful TUI"
  homepage "https://github.com/Shubhpreet-Rana/codegenius"
  version "1.1.2"
  license "MIT"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v1.1.2/codegenius-darwin-arm64"
      sha256 "ed4c797036f42d028a2a733586be2e51f838544c30d3b9641c4c5210e0d4dd81"
    else
      url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v1.1.2/codegenius-darwin-amd64"
      sha256 "6f73e4fa9c7c1f610ec9cc86acbb4effce926b257dc781fefde463262ef00047"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v1.1.2/codegenius-linux-arm64"
      sha256 "05d62bf6c86bc7c2812e768de6141a69357cc992b6c0494086d9a9cfe3b4dc56"
    else
      url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v1.1.2/codegenius-linux-amd64"
      sha256 "692d20c1fe050799d10637863226c63a57585d5b4a2ffa4c9c9db6c62afd381b"
    end
  end

  def install
    bin.install Dir["codegenius*"].first => "codegenius"
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