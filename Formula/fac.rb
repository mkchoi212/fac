class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.0/fac_1.0.0_darwin_amd64.tar.gz"
  version "1.0.0"
  sha256 "5e3efea7bed1d35a1bec25bb5af1d2c03e1c8fe0ece4a9fe23c180ce7f060066"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
