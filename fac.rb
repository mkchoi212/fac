class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.1/fac_1.0.1_darwin_amd64.tar.gz"
  version "1.0.1"
  sha256 "6a51d3b363a84243309bfc050c68c5dd3bc1732bea343b32386630ecb2e89421"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
