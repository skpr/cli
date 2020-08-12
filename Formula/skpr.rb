class Skpr < Formula
  desc "CLI for the Skpr Hosting Platform"
  homepage "https://www.skpr.io"
  url "https://github.com/skpr/cli/releases/download/v0.6.5/skpr_darwin_amd64.tgz"
  version "v0.6.5"
  sha256 "71b773e7a5e039dc6b5d40256b9e82061a8f283f67ce26b5a57a36fcc07c852a"

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
