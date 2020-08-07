class Skpr < Formula
  desc "CLI for the Skpr Hosting Platform"
  homepage "https://www.skpr.io"
  url "https://github.com/skpr/cli/releases/download/v0.6.4-kimtest5/skpr_darwin_amd64.tgz"
  version "v0.6.4-kimtest5"
  sha256 "be1233571f7a7e9e2fa78d6ab82f192a6b91936e59bf007860eeb5fb5ae65187"

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
