---
sidebar_position: 2
sidebar_label: Capabilities
description: MMDB CLI Capabilities
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
  - capability
---

# Capabilities

MMDB CLI is a powerful command-line tool offering a wide range of features for working with MMDB files. Some of the key features include:

- **Generate MMDB**: You can generate MMDB files from JSON datasets using the `generate` command.
- **Dump MMDB**: Export a database to JSON (`dump -f json`) or custom text via the [`jsonpath` output format](./guides/jsonpath).
- **Check MMDB Metadata**: Print metadata in YAML, JSON, XML, or CSV using the `metadata` command.
- **Modify MMDB**: You can modify existing MMDB files by inserting, updating, or deleting records using the `update` command.
- **Inspect MMDB**: Look up IPs/CIDRs (YAML, JSON, XML, CSV, or [`jsonpath` output format](./guides/jsonpath)).
- **Verify MMDB**: Validate database integrity with the `verify` command.

<!---
  completion  Generate the autocompletion script for the specified shell
  dump        Dump MMDB data into a json dataset
  generate    Generate a MMDB database from a JSON dataset
  help        Help about any command
  inspect     Inspect an IP address or CIDR in the MMDB file
  metadata    Prints metadata of the MMDB file
  update      Update existing MMDB file
  version     Show version information for mmdb-cli

--->

```mermaid
mindmap
    root((MMDB CLI))
        metadata((Metadata))
            readMMDBMetadata(Read MMDB Metadata)
        inspect((Inspect))
            inspectIP(Inspect Single IP)
            inspectMultipleIPs(Inspect Multiple IPs)
            inspectNetworkCIDR(Inspect Network CIDR)
            inspectFormats(Output formats incl. jsonpath)
        generate((Generate))
            generateMMDB(Generate MMDB from JSON Dataset V1)
        dump((Dump))
            dumpDatasetV1(Dump JSON Dataset V1)
            dumpFormats(JSON export and jsonpath format)
        update((Update))
            updateMMDBRemove(Remove Record)
            updateMMDBInsert(Insert Record)
            updateMMDBReplace(Replace Record)
            updateMMDBMerge(Merge Records)
        verify((Verify))
            verifyMMDB(Verify MMDB File)
```

## Commands

MMDB CLI provides the following commands to work with MMDB files:

```bash

InfraZ MMDB CLI is a command line toolkit for working with MMDB
Complete documentation is available at https://docs.infraz.io/mmdb-cli

Usage:
  mmdb-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dump        Dump MMDB data into a json dataset
  generate    Generate a MMDB database from a JSON dataset
  help        Help about any command
  inspect     Inspect an IP address or CIDR in the MMDB file
  metadata    Prints metadata of the MMDB file
  update      Update existing MMDB file
  verify      Verify the MMDB file
  version     Show version information for mmdb-cli

Flags:
  -h, --help   help for mmdb-cli

Use "mmdb-cli [command] --help" for more information about a command.
```

## Implement new features

If you have any feature requests or suggestions for **MMDB CLI**, we encourage you to create an issue on our [GitHub repository](https://github.com/InfraZ/mmdb-cli/issues). We are continuously seeking ways to enhance MMDB CLI and make it even more valuable for our users.
