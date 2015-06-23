require 'formula'

class ComposeRepo < Formula
  desc "syncronize your docker-compose repos and rependancies"
  homepage "https://github.com/Assembla/compose-repo"
  url "https://github.com/Assembla/compose-repo/releases/download/0.0.4/compose-repo_darwin_amd64.zip"
  version "0.0.4"
  sha256 "9fa1d638cc5c2abb72db2ee3076189e940a6b818aebd7dfb22a8da624983aeb4"
  depends_on :arch => :x86_64

  def install
    bin.install "compose-repo"
  end
end
