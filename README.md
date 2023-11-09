# Skpr CLI

Skpr is a fully managed container-based hosting platform.

For more information see https://www.skpr.io

This repository is used soley for the releases of the skpr command line utilty.

Full documentation is available at https://docs.skpr.io

## Installation

### Debian/Ubuntu

Add the apt repository and public key to your config, then update and install.

```
wget -q https://packages.skpr.io/apt/packages.skpr.io.pub -O- | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/packages.skpr.io.pub > /dev/null
echo "deb [arch=amd64 signed-by=/etc/apt/trusted.gpg.d/packages.skpr.io.pub] https://packages.skpr.io/apt stable main" | sudo tee -a /etc/apt/sources.list.d/skpr.list > /dev/null
sudo apt update && sudo apt install skpr
```

### Homebrew

Homebrew is a package manager for MacOS. See https://brew.sh/ for further details.

To install via homebrew, run:

```
brew tap skpr/taps
brew install skpr
```

To upgrade, run:

```
brew upgrade skpr
```


### Manual installation

You can download the binaries from the [Releases](https://github.com/skpr/skpr/releases) section.

(Replace `VERSION` with a release version)

#### MacOS

```
curl -sSLO https://github.com/skpr/cli/releases/download/VERSION/skpr_darwin_amd64.tgz
sudo tar -zxf skpr_darwin_amd64.tgz -C /usr/local/bin/
```

#### Linux

```
curl -sSLO https://github.com/skpr/cli/releases/download/VERSION/skpr_linux_amd64.tgz
sudo tar -zxf skpr_linux_amd64.tgz -C /usr/local/bin/
```

## Documentation

Documentation can be found at https://docs.skpr.io/
