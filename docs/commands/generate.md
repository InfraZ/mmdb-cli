---
sidebar_position: 8
sidebar_label: Generate
description: MMDB CLI Generate Command
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
---

# Generate Command

The `generate` command is one of amazing features of MMDB CLI that allows you to generate MMDB files from JSON datasets. You can use this command to convert your JSON datasets into MMDB files and create custom databases for your applications. This command is useful for creating custom databases for IP geolocation, ASN lookup, and other use cases.

:::info[Dataset Schema]

The JSON dataset must use `dataset` records and `metadata` fields supported by MMDB CLI. For details, see [Dataset Schema](../dataset-schema).

:::

## Usage

```bash
mmdb-cli generate -i <JSON_FILE_PATH> -o <MMDB_OUTPUT_PATH>
```

## Options

- `-i, --input <JSON_FILE_PATH>`: The path to the JSON file that contains the dataset you want to convert to an MMDB file.
- `-o, --output <MMDB_OUTPUT_PATH>`: The path to the MMDB file where the output will be saved. (must have a .mmdb extension)
- `--disable-ipv4-aliasing`: Disable IPv4 aliasing for IPv6 networks. By default, IPv4 addresses are aliased to their IPv6 counterparts. Use this option to disable this feature.
- `--include-reserved-networks`: Include reserved networks in the generated MMDB file. By default, reserved networks are excluded from the output.
- `-v, --verbose`: Enable the verbose mode

## Examples

In the following example, we generate an MMDB file from a JSON dataset:

:::note[Disclaimer]

The data shown in the examples above is for demonstration purposes only and may not reflect the actual data in the MMDB file you are inspecting. The actual data may vary based on the version of the database and the IP addresses queried.

:::

```json
{
  "dataset": [
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
  ],
  "metadata": {
    "Description": {
      "en": "GeoLite2 ASN Custom Database"
    },
    "DatabaseType": "GeoLite2-ASN-Custom",
    "Languages": [
      "en"
    ],
    "IPVersion": 6,
    "RecordSize": 24
  },
  "version": "v1"
}
```

```bash
mmdb-cli generate -i GeoLite2-ASN-Custom.json -o GeoLite2-ASN-Custom.mmdb
```

To verify the generated MMDB file, you can use the `inspect` command to view the contents of the file, for more information, see [Inspect Command](./inspect).
