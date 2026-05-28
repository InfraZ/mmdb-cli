---
sidebar_position: 7
sidebar_label: Dump
description: MMDB CLI Dump Command
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
---

# Dump Command

The `dump` command exports an entire MMDB database to a file. By default it writes the [dataset schema v1](../dataset-schema) JSON document.

:::tip[Dataset schema]

Default JSON output includes `version`, `metadata`, and `dataset` sections. See [Dataset Schema](../dataset-schema).

:::

## Usage

```bash
mmdb-cli dump -i <MMDB_FILE_PATH> -o <OUTPUT_PATH> [OPTIONS]
```

## Options

| Flag | Description |
| :--- | :---------- |
| `-i, --input <PATH>` | Path to the MMDB file (**required**). |
| `-o, --output <PATH>` | Output file path (**required**). |
| `-f, --format <FORMAT>` | `json` (default) or `jsonpath={TEMPLATE}`. |
| `-v, --verbose` | Log each record while dumping (JSON mode only). |

### Output formats (`-f`)

| Format | Output file | Description |
| :----- | :------------ | :---------- |
| `json` | Must end with `.json` | Full dataset JSON (`version`, `metadata`, `dataset`). |
| `jsonpath={TEMPLATE}` | Any extension (e.g. `.txt`) | Template output (see [JSONPath output format](../guides/jsonpath)). |

## Example — Full JSON dump

```bash
mmdb-cli dump -i GeoLite2-ASN.mmdb -o GeoLite2-ASN.json
```

### Output structure (JSON)

:::note[Disclaimer]

Example records are illustrative.

:::

```json
{
  "version": "v1",
  "metadata": {
    "Description": {
      "en": "GeoLite2 ASN database"
    },
    "DatabaseType": "GeoLite2-ASN",
    "Languages": ["en"],
    "BinaryFormatMajorVersion": 2,
    "BinaryFormatMinorVersion": 0,
    "BuildEpoch": 1728288921,
    "IPVersion": 6,
    "NodeCount": 1056663,
    "RecordSize": 24
  },
  "dataset": [
    {
      "network": "::100:0/120",
      "record": {
        "autonomous_system_number": 13335,
        "autonomous_system_organization": "CLOUDFLARENET"
      }
    },
    {
      "network": "::100:400/118",
      "record": {
        "autonomous_system_number": 38803,
        "autonomous_system_organization": "Gtelecom Pty Ltd"
      }
    },
    ...
}
```

Verbose progress:

```bash
mmdb-cli dump -i GeoLite2-ASN.mmdb -o GeoLite2-ASN.json -v
```
