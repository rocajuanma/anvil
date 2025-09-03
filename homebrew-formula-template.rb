# This is a template for your Homebrew formula
# Create a new repository called "homebrew-anvil" and place this file at: Formula/anvil.rb

class Anvil < Formula
  desc "Complete macOS development environment automation tool"
  homepage "https://github.com/rocajuanma/anvil"
  url "https://github.com/rocajuanma/anvil/archive/v1.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256_OF_TARBALL"
  license "Apache-2.0"
  head "https://github.com/rocajuanma/anvil.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./main.go"
    
    # Generate shell completions (optional but recommended)
    generate_completions_from_executable(bin/"anvil", "completion")
  end

  def caveats
    <<~EOS
      🔨 Anvil has been installed!
      
      Get started:
        anvil init          # Initialize your environment
        anvil doctor        # Verify setup and dependencies
        anvil install dev   # Install development tools
      
      📚 Documentation: https://github.com/rocajuanma/anvil/docs
    EOS
  end

  test do
    assert_match "Anvil", shell_output("#{bin}/anvil --help")
    assert_match version.to_s, shell_output("#{bin}/anvil --version")
  end
end
