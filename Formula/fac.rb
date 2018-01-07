class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.4/fac_1.0.4_darwin_amd64.tar.gz"
  version "1.0.4"
  sha256 "5aaab92e82e02594126f54ab973b2c9eaab422ac13b87d7c79fe7f1afdc42965"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
