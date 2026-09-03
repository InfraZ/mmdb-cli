# MMDB CLI

<img src="https://docs.infraz.io/img/logo/logo_transparent_white.png" width="200">

![GitHub License](https://img.shields.io/github/license/InfraZ/mmdb-cli)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/InfraZ/mmdb-cli/release.yaml)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/InfraZ/mmdb-cli/tests.yaml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FInfraZ%2Fmmdb-cli.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FInfraZ%2Fmmdb-cli?ref=badge_shield&issueType=license)

MMDB CLI is a command-line open-source project, developed to simplify working with MMDB files. It is a command-line tool that allows you to generate, modify, export, and inspect MMDB files, among other functionalities.

## Supported Platforms

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

## Documentation

The official documentation for MMDB CLI is available on the
 [InfraZ Documentation Website](https://docs.infraz.io/docs/mmdb-cli).

> [!TIP]
> We recommend reading the documentation to get a better understanding of the MMDB CLI and its functionalities.

## Installation (Homebrew)

On macOS, you can install MMDB CLI via the [InfraZ Homebrew tap](https://github.com/InfraZ/homebrew-tap):

```bash
brew tap infraz/tap
brew trust infraz/tap
brew install infraz/tap/mmdb-cli
```

## Installation (Pre-compiled Binaries)

### Linux and macOS

1. Choose the version and platform you want to install from the [GitHub releases page](https://github.com/InfraZ/mmdb-cli/releases).

    ```bash
    export MMDB_CLI_VERSION=0.5.0
    export MMDB_CLI_PLATFORM=linux_amd64
    ```

2. Download the MMDB CLI tarball using `curl`, `wget`, or any other tool.

    ```bash
    curl -LO "https://github.com/InfraZ/mmdb-cli/releases/download/v${MMDB_CLI_VERSION}/mmdb-cli_${MMDB_CLI_VERSION}_${MMDB_CLI_PLATFORM}.tar.gz"
    ```

3. Extract the downloaded tarball.

    ```bash
    tar -xzf mmdb-cli_${MMDB_CLI_VERSION}_${MMDB_CLI_PLATFORM}.tar.gz
    ```

4. Move the extracted binary file to a directory in your PATH.

    ```bash
    sudo mv mmdb-cli /usr/local/bin/
    ```

5. Verify the installation by running the following command.

    ```bash
    mmdb-cli --version
    ```

## Container Image

Pre-built multi-arch images (`linux/amd64`, `linux/arm64`) are published to
the GitHub Container Registry:

| Reference | Points to |
| :-- | :-- |
| `ghcr.io/infraz/mmdb-cli:latest` | Latest commit on `main` |
| `ghcr.io/infraz/mmdb-cli:1.2.3` | Tagged release `v1.2.3` |

The image `ENTRYPOINT` is `mmdb-cli` and the working directory is `/data`, so
mount your files there:

```bash
# Inspect a database
docker run --rm -v "$PWD:/data" ghcr.io/infraz/mmdb-cli:latest \
  metadata -i GeoLite2-City.mmdb

# Generate a database from a JSON dataset (writes back to the host)
docker run --rm -v "$PWD:/data" --user "$(id -u):$(id -g)" \
  ghcr.io/infraz/mmdb-cli:latest \
  generate -i dataset.json -o result.mmdb
```

The container runs as a non-root user (UID 1000). When writing files to a
bind mount, add `--user "$(id -u):$(id -g)"` so the output is owned by you.

### Building the image locally

```bash
docker build -t mmdb-cli --build-arg VERSION=v0.0.0-dev .
docker run --rm mmdb-cli version
```

## Development

To get started, clone the repository and run the following commands to download the dependencies:

```bash
go mod download -x # Download dependencies
```

To build the project, run the following command:

```bash
go build -o mmdb-cli # Build the project
```

Then, you can run the project with the following command:

```bash
./mmdb-cli # Run the project
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FInfraZ%2Fmmdb-cli.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FInfraZ%2Fmmdb-cli?ref=badge_large&issueType=license)
