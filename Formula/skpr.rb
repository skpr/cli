class Skpr < Formula
  desc "CLI for the Skpr Hosting Platform"
  homepage "https://www.skpr.io"
  url "https://github.com/skpr/cli/releases/download/v0.7.0/skpr_darwin_amd64.tgz"
  version "v0.7.0"
  sha256 "9976c11c188747bb643cfae0b8952bef0c89822f93950dd2842a093cc39de7c9"

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
