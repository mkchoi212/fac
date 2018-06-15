class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v2.0.0/fac_2.0.0_darwin_amd64.tar.gz"
  version "2.0.0"
  sha256 "8e33b3375169f19b5b6a508cd69287ce16d75ad6dc06c1d09d9684de5010c0eb"

  def install
    bin.install "fac"
  end
end
