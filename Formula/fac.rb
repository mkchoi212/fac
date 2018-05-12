class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.1.0/fac_1.1.0_darwin_amd64.tar.gz"
  version "1.1.0"
  sha256 "01e72c6147ea5d5cfb10f9c95e13a27ba764b7e71cdd409a96fd566ee16fbc4c"
  
  depends_on "git"

  def install
    bin.install "fac"
  end
end
