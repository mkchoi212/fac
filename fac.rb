class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.1/fac_1.0.1_darwin_amd64.tar.gz"
  version "1.0.1"
  sha256 "f64e25dae39e9eb160afe225754b345e850489f7dcab707bc513311d7098d59a"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
