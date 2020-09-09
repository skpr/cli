class Skpr < Formula
  desc "CLI for the Skpr Hosting Platform"
  homepage "https://www.skpr.io"
  url "https://github.com/skpr/cli/releases/download/v0.7.1/skpr_darwin_amd64.tgz"
  version "v0.7.1"
  sha256 "8ea87999beb1187d6ced40a47df417ea1d466928d498b4d6a760d7a0797c978a"

  def install
    bin.install "skpr"
    bin.install "skpr-rsh"

    # Install bash completion
    output = Utils.safe_popen_read("#{bin}/skpr", "--completion-script-bash")
    (bash_completion/"skpr").write output

    # Install zsh completion
    output = Utils.safe_popen_read("#{bin}/skpr", "--completion-script-zsh")
    (zsh_completion/"_skpr").write output
  end
end
