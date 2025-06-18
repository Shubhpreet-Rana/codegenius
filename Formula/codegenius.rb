class Codegenius < Formula
  desc "AI-powered Git commit message generator and code reviewer with beautiful terminal UI"
  homepage "https://github.com/Shubhpreet-Rana/codegenius"
  version "1.1.1"
  license "MIT"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v#{version}/codegenius-darwin-amd64"
    sha256 "your-sha256-hash-here"  # Update this with actual hash
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v#{version}/codegenius-darwin-arm64"
    sha256 "your-sha256-hash-here"  # Update this with actual hash
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v#{version}/codegenius-linux-amd64"
    sha256 "your-sha256-hash-here"  # Update this with actual hash
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/Shubhpreet-Rana/codegenius/releases/download/v#{version}/codegenius-linux-arm64"
    sha256 "your-sha256-hash-here"  # Update this with actual hash
  end

  def install
    bin.install Dir["*"].first => "codegenius"
  end

  def caveats
    <<~EOS
      🤖 CodeGenius CLI installed successfully!

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
    system "#{bin}/codegenius", "--help"
  end
end 