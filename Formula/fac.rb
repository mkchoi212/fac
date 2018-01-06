class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.3/fac_1.0.3_darwin_amd64.tar.gz"
  version "1.0.3"
  sha256 "a1ee8ed558fba5eeef74068a7666ec4c1b28f809ecd31e44055152c177dbe1b3"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
