class Fac < Formula
  desc "Command line User Interface for fixing git conflicts"
  homepage "https://github.com/mkchoi212/fac"
  url "https://github.com/mkchoi212/fac/releases/download/v1.0.2/fac_1.0.2_darwin_amd64.tar.gz"
  version "1.0.2"
  sha256 "d5369b5a8532b8be25b061f32ace7531868e361258141f6ac44932006eb4e451"
  
  depends_on "git"
  depends_on "go"

  def install
    bin.install "fac"
  end

  test do
    
  end
end
