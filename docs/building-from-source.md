---
sidebar_position: 4
sidebar_label: Building from Source
description: MMDB CLI Building from Source
tags:
  - cli
  - mmdb
  - mmdb-cli
---

# Building MMDB CLI from Source 🏗️

If you want to build MMDB CLI from the source code, you can follow the instructions below.

## Prerequisites

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/) (version 1.26 or higher)

## Clone the Repository

First, you need to clone the MMDB CLI repository to your local machine.

```bash
git clone https://github.com/InfraZ/mmdb-cli.git
cd mmdb-cli
```

## Download Dependencies

Next, you need to download the dependencies using the `go mod download` command.

```bash
go mod download -x
```

## Build the Binary

Finally, you can build the MMDB CLI binary using the `go build` command.

```bash
go build -o mmdb-cli .
```

## Verify the Binary

You can verify the binary by running the following command.

```bash
./mmdb-cli version
```

**Note:** If you encounter any issues during the build process, please let us know by creating an issue on our GitHub repository.
