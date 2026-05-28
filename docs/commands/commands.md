---
sidebar_position: 6
sidebar_label: Commands
description: MMDB CLI Commands
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
---

# Commands

The MMDB CLI provides several commands to interact with MaxMind DB files. You can use these commands to inspect the contents of an MMDB file, extract metadata, and dump the data in a human-readable format. The following commands are available in MMDB CLI:

We will explore each of these commands in detail in the following sections.

You can use `--help` flag to get more information about available commands at any time:

```bash
mmdb-cli --help
```

## Core Commands

- [Metadata](./metadata): View the metadata of an MMDB file.
- [Inspect](./inspect): Look up IPs/CIDRs in an MMDB file.
- [Dump](./dump): Export a database to a file.
- [Generate](./generate): Generate an MMDB file from a JSON dataset.
- [Update](./update): Update an existing MMDB file with new data.
- [Verify](./verify): Verify that an MMDB file is valid.

## Built-in Utility Commands

- `completion`: Generate shell autocompletion scripts (`bash`, `zsh`, `fish`, `powershell`).
- `help`: Show help for any command.

### `completion` Usage

```bash
mmdb-cli completion [bash|zsh|fish|powershell]
```
