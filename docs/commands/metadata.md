---
sidebar_position: 1
sidebar_label: Metadata
description: MMDB CLI Metadata Command
tags:
  - cli
  - mmdb
  - mmdb-cli
---

# Metadata Command 🪨

The `metadata` command prints MMDB metadata including database type, description, languages, IP version, node count, and record size.

:::note[Disclaimer]

The data shown in the examples above is for demonstration purposes only and may not reflect the actual data in the MMDB file you are inspecting. The actual data may vary based on the version of the database and the IP addresses queried.

:::

## Usage

```bash
mmdb-cli metadata -i <MMDB_FILE_PATH> [-f <FORMAT>]
```

## Options

- `-i, --input <MMDB_FILE_PATH>`: Path to the MMDB file.
- `-f, --format <FORMAT>`: Output format. Supported values: `yaml`, `json`, `json-pretty`, `xml`, `csv` (default: `yaml`).

## Examples

In the following example, we extract metadata from the `GeoLite2-ASN.mmdb` file in YAML format:

```bash
mmdb-cli metadata -i GeoLite2-ASN.mmdb
```

### Output (YAML):

```yaml
binary_format_major_version: 2
binary_format_minor_version: 0
build_epoch: 1728288921
database_type: GeoLite2-ASN
description:
    en: GeoLite2 ASN database
ip_version: 6
languages:
    - en
node_count: 1056663
record_size: 24
```

### Output (JSON)

```bash
mmdb-cli metadata -i GeoLite2-ASN.mmdb -f json
```

```json
{
    "description": {
        "en": "GeoLite2 ASN database"
    },
    "database_type": "GeoLite2-ASN",
    "languages": [
        "en"
    ],
    "binary_format_major_version": 2,
    "binary_format_minor_version": 0,
    "build_epoch": 1728288921,
    "ip_version": 6,
    "node_count": 1056663,
    "record_size": 24
}
```

### Output (XML / CSV)

```bash
mmdb-cli metadata -i GeoLite2-ASN.mmdb -f xml
mmdb-cli metadata -i GeoLite2-ASN.mmdb -f csv
```
