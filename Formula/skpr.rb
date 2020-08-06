class Notes < Formula
  desc "CLI for the Skpr Hosting Platform"
  homepage "https://www.skpr.io"
  url "https://github.com/skpr/cli/releases/${VERSION}/skpr-cli_darwin_amd64.tgz"
  version "${VERSION}"
  sha256 "${SHA_256-SUM}"

  def install
    bin.install "bin/amd64/darwin/skpr"
    bin.install "bin/amd64/darwin/skpr-rsh"

    # Install bash completion
    output = Utils.safe_popen_read("#{bin}/skpr", "--completion-script-bash")
    (bash_completion/"skpr").write output

    # Install zsh completion
    output = Utils.safe_popen_read("#{bin}/skpr", "--completion-script-zsh")
    (zsh_completion/"_skpr").write output
  end
end