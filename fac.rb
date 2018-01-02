class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.1/fac_1.0.1_darwin_amd64.tar.gz"
  version "1.0.1"
  sha256 "74e28c0400a71e4cc2bf3f0b54a74d3e98cdcedb048b36445fbf55d42fa7c8ca"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
