---
sidebar_position: 2
sidebar_label: Inspect
description: MMDB CLI Inspect Command
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
---

# Inspect Command

The `inspect` command looks up one or more IP addresses or CIDR prefixes in an MMDB file and prints every matching network and its record.

:::note[Disclaimer]

The data shown in the examples above is for demonstration purposes only and may not reflect the actual data in the MMDB file you are inspecting. The actual data may vary based on the version of the database and the IP addresses queried.

:::

## Usage

```bash
mmdb-cli inspect -i <MMDB_FILE_PATH> [OPTIONS] <IP_OR_CIDR> [MORE_IPS_OR_CIDRS...]
```

## Options

| Flag | Description |
| :--- | :---------- |
| `-i, --input <PATH>` | Path to the MMDB file (**required**). |
| `-f, --format <FORMAT>` | Output format (see below). Default: `yaml`. |

### Output formats (`-f`)

| Format | Description |
| :----- | :---------- |
| `yaml` | Human-readable YAML (default). |
| `json` | Compact JSON. |
| `json-pretty` | Indented JSON. |
| `xml` | XML with a `<root>` wrapper. |
| `csv` | Tabular CSV (flattened fields). |
| `jsonpath={TEMPLATE}` | Template output to stdout (see [JSONPath output format](../guides/jsonpath)). |

## Example — Single IP

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb 1.1.1.1
```

### Output (YAML)

```yaml
- query: 1.1.1.1
  records:
    - network: 1.1.1.0/24
      record:
        autonomous_system_number: 13335
        autonomous_system_organization: CLOUDFLARENET
```

### Output (JSON)

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb -f json 1.1.1.1
```

```json
[
    {
        "query": "1.1.1.1",
        "records": [
            {
                "network": "1.1.1.0/24",
                "record": {
                    "autonomous_system_number": 13335,
                    "autonomous_system_organization": "CLOUDFLARENET"
                }
            }
        ]
    }
]
```

### Output (XML / CSV)

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb -f xml 1.1.1.1
mmdb-cli inspect -i GeoLite2-ASN.mmdb -f csv 1.1.1.1
```

## Example — CIDR Range

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb 1.1.1.0/20
```

Returns every network contained in the range, each with its record.

### Output YAML (CIDR Range)

```yaml
- query: 1.1.1.1/20
  records:
    - network: 1.1.1.0/24
      record:
        autonomous_system_number: 13335
        autonomous_system_organization: CLOUDFLARENET
    - network: 1.1.8.0/24
      record:
        autonomous_system_number: 138421
        autonomous_system_organization: China Unicom
```

### Output JSON (CIDR Range)

```json
[
    {
        "query": "1.1.1.1/20",
        "records": [
            {
                "network": "1.1.1.0/24",
                "record": {
                    "autonomous_system_number": 13335,
                    "autonomous_system_organization": "CLOUDFLARENET"
                }
            },
            {
                "network": "1.1.8.0/24",
                "record": {
                    "autonomous_system_number": 138421,
                    "autonomous_system_organization": "China Unicom"
                }
            }
        ]
    }
]
```

## Example — Multiple IPs

In the following example, we inspect the `GeoLite2-ASN.mmdb` file for multiple IP addresses:

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb 1.0.0.1 1.1.1.1
```

Each query appears as a separate top-level entry with its own `query` and `records` list.

### Output YAML

```yaml
- query: 1.0.0.1
  records:
    - network: 1.0.0.0/24
      record:
        autonomous_system_number: 13335
        autonomous_system_organization: CLOUDFLARENET
- query: 1.1.1.1
  records:
    - network: 1.1.1.0/24
      record:
        autonomous_system_number: 13335
        autonomous_system_organization: CLOUDFLARENET
```

### Output JSON

```json
[
    {
        "query": "1.0.0.1",
        "records": [
            {
                "network": "1.0.0.0/24",
                "record": {
                    "autonomous_system_number": 13335,
                    "autonomous_system_organization": "CLOUDFLARENET"
                }
            }
        ]
    },
    {
        "query": "1.1.1.1",
        "records": [
            {
                "network": "1.1.1.0/24",
                "record": {
                    "autonomous_system_number": 13335,
                    "autonomous_system_organization": "CLOUDFLARENET"
                }
            }
        ]
    }
]
```
