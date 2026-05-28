---
sidebar_position: 3
sidebar_label: Installation
description: MMDB CLI Installation
tags:
  - cli
  - mmdb
  - mmdb-cli
  - installation
  - geoip
---

# MMDB CLI Installation

MMDB CLI ships as a single static binary. You can install it with [Homebrew](#installation-macos-homebrew) on macOS, Linux package managers, pre-built archives from [GitHub Releases](https://github.com/InfraZ/mmdb-cli/releases), or [build from source](./building-from-source).

## Supported Platforms

Pre-built binaries are published for the following OS and architecture combinations:

| Platform |     Architecture      | Supported |
| :------: | :-------------------: | :-------: |
|  Linux   |         amd64         |    ✅     |
|  Linux   |         arm64         |    ✅     |
|  macOS   |         amd64         |    ✅     |
|  macOS   | arm64 (Apple Silicon) |    ✅     |
| FreeBSD  |         amd64         |    ✅     |
| FreeBSD  |         arm64         |    ✅     |
| Windows  |         amd64         |    ✅     |
| Windows  |         arm64         |    ✅     |

**Note:** If your platform is not listed above, you can [build MMDB CLI from source](./building-from-source).

**Note:** We mainly test MMDB CLI on Linux (amd64) and macOS (arm64). If you encounter issues on other platforms, please [open an issue](https://github.com/InfraZ/mmdb-cli/issues).

:::tip[Linux packages]

We also publish `deb`, `rpm`, `apk`, and Arch Linux packages. Download them from [GitHub Releases](https://github.com/InfraZ/mmdb-cli/releases).

:::

| Package Format | Architecture | Supported |
| :------------: | :----------: | :-------: |
|      deb       |    amd64     |    ✅     |
|      deb       |    arm64     |    ✅     |
|      rpm       |    amd64     |    ✅     |
|      rpm       |    arm64     |    ✅     |
|      apk       |    amd64     |    ✅     |
|      apk       |    arm64     |    ✅     |
|   Arch Linux   |    amd64     |    ✅     |
|   Arch Linux   |    arm64     |    ✅     |

## Installation (macOS — Homebrew)

On macOS, install MMDB CLI from the [InfraZ Homebrew tap](https://github.com/InfraZ/homebrew-tap):

```bash
brew tap infraz/tap
brew install infraz/tap/mmdb-cli
```

Verify:

```bash
mmdb-cli version
```

:::caution[macOS Gatekeeper]

The Homebrew cask is not Apple-notarized. If macOS blocks the binary, remove the quarantine attribute:

```bash
xattr -dr com.apple.quarantine "$(which mmdb-cli)"
```

:::

## Installation Instructions (Linux Package Manager)

### Debian and Ubuntu

1. Resolve the latest version and architecture:

```bash
export MMDB_CLI_VERSION=$(curl -s https://api.github.com/repos/InfraZ/mmdb-cli/releases/latest | grep tag_name | cut -d '"' -f 4 | sed 's/v//')
export ARCH=$(dpkg --print-architecture)
```

2. Download the `.deb` package using `curl` or `wget`
```bash
curl -LO "https://github.com/InfraZ/mmdb-cli/releases/download/v${MMDB_CLI_VERSION}/mmdb-cli_${MMDB_CLI_VERSION}_linux_${ARCH}.deb"
```

3. Install `.deb` package using `apt`
```bash
apt install ./mmdb-cli_${MMDB_CLI_VERSION}_linux_${ARCH}.deb
```

4. Verify the installation by running the following command.
```bash
mmdb-cli version
```

### RH-based (RHEL, Fedora, CentOS, Rocky Linux, AlmaLinux)

1. Resolve the latest version:

```bash
export MMDB_CLI_VERSION=$(curl -s https://api.github.com/repos/InfraZ/mmdb-cli/releases/latest | grep tag_name | cut -d '"' -f 4 | sed 's/v//')
```

2. Install the `.rpm` package using `yum` or `dnf`
```bash
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64; [ "$ARCH" = "aarch64" ] && ARCH=arm64
dnf install "https://github.com/InfraZ/mmdb-cli/releases/download/v${MMDB_CLI_VERSION}/mmdb-cli_${MMDB_CLI_VERSION}_linux_${ARCH}.rpm"
```

3. Verify the installation by running the following command.
```bash
mmdb-cli version
```

### Alpine Linux

1. Resolve the latest version:

```bash
export MMDB_CLI_VERSION=$(curl -s https://api.github.com/repos/InfraZ/mmdb-cli/releases/latest | grep tag_name | cut -d '"' -f 4 | sed 's/v//')
```

2. Download the `.apk` package using `curl` or `wget`
```bash
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64; [ "$ARCH" = "aarch64" ] && ARCH=arm64
curl -LO "https://github.com/InfraZ/mmdb-cli/releases/download/v${MMDB_CLI_VERSION}/mmdb-cli_${MMDB_CLI_VERSION}_linux_${ARCH}.apk"
```

3. Install `.apk` package using `apk`
```bash
apk add --allow-untrusted ./mmdb-cli_${MMDB_CLI_VERSION}_linux_${ARCH}.apk
```

4. Verify the installation by running the following command.
```bash
mmdb-cli version
```

## Installation Instructions (Pre-built Binary)

### Linux and macOS

1. Pick a version and platform from [GitHub Releases](https://github.com/InfraZ/mmdb-cli/releases):

```bash
export MMDB_CLI_VERSION=0.5.0
export MMDB_CLI_PLATFORM=linux_amd64   # e.g. darwin_arm64, freebsd_amd64
```

2. Download the MMDB CLI tarball using `curl`, `wget`, or any other tool.

```bash
curl -LO "https://github.com/InfraZ/mmdb-cli/releases/download/v${MMDB_CLI_VERSION}/mmdb-cli_${MMDB_CLI_VERSION}_${MMDB_CLI_PLATFORM}.tar.gz"
```

3. Extract the downloaded tarball.

```bash
tar -xzf mmdb-cli_${MMDB_CLI_VERSION}_${MMDB_CLI_PLATFORM}.tar.gz
```

4. Move the extracted binary file to a directory in your `$PATH`.

```bash
mv mmdb-cli /usr/local/bin/
```

5. Verify the installation by running the following command.

```bash
mmdb-cli version
```

### FreeBSD

Use the same tarball flow as Linux/macOS with a FreeBSD platform label, for example `freebsd_amd64` or `freebsd_arm64`.

### Windows

1. Download `mmdb-cli_<version>_windows_<arch>.zip` from [GitHub Releases](https://github.com/InfraZ/mmdb-cli/releases).
2. Extract `mmdb-cli.exe` to a directory on your `PATH`.
3. Open a new terminal and run:

```powershell
mmdb-cli version
```
