---
sidebar_position: 9
sidebar_label: Update
description: MMDB CLI Update Command
tags:
  - cli
  - mmdb
  - mmdb-cli
---

# Update Command 🧠

The `update` command updates existing MMDB files using dataset operations. It supports record removal, replacement, and merge strategies without recreating the database from scratch.

## Usage

```bash
mmdb-cli update -i <MMDB_FILE_PATH> -o <MMDB_OUTPUT_PATH> -d <JSON_FILE_PATH>
```

## Options

- `-i, --input <MMDB_FILE_PATH>`: The path to the MMDB file you want to update.
- `-o, --output <MMDB_OUTPUT_PATH>`: The path to the updated MMDB file where the changes will be saved. (must have a .mmdb extension)
- `-d, --dataset <JSON_FILE_PATH>`: The path to the dataset JSON file that contains update operations.
- `--disable-ipv4-aliasing`: Disable IPv4 aliasing for IPv6 networks. By default, IPv4 addresses are aliased to their IPv6 counterparts. Use this option to disable this feature.
- `--include-reserved-networks`: Include reserved networks in the generated MMDB file. By default, reserved networks are excluded from the output.
- `-v, --verbose`: Enable the verbose mode

## JSON Update Schema

The update data in JSON format must contain a top-level `dataset` array. Each object in that array represents one network CIDR operation.

| Field   | Type   | Description                           | Required |
| ------- | ------ | ------------------------------------- | -------- |
| network | string | The network CIDR block to update      | Yes      |
| method  | string | The action to perform on the network. | No (defaults to `deep_merge`) |
| data    | object | The data to insert, update, or delete | Yes      |

## Supported Methods

:::info[Supported Methods]

The supported methods are `remove`, `replace`, `top_level_merge`, and `deep_merge`. These methods are implemented based on the official MMDB writer library. For details, see the [MMDB Writer documentation](https://github.com/maxmind/mmdbwriter).

:::

|     Method      | Description                                                |
| :-------------: | ---------------------------------------------------------- |
|     remove      | Remove the network from the MMDB file                      |
|     replace     | Replace the data of the network with the new data          |
| top_level_merge | Merge the new data with the existing data at the top level |
|   deep_merge    | Deep merge the new data with the existing data             |

```json
{
  "version": "v1",
  "dataset": [
    {
        "network": "<NETWORK_CIDR>",
        "method": "<ACTION_METHOD>",
        "data": {
            <UPDATE_DATA>
        }
    },
    ...
  ]
}
```

## Examples

In the following example, we update the `GeoLite2-ASN.mmdb` file with new data from a JSON file:

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "1.1.1.1/32",
      "method": "deep_merge",
      "data": {
        "is_cloudflare": true,
        "attributes": {
          "country": "US",
          "city": "San Francisco"
        }
      }
    },
    {
      "network": "8.8.8.8/32",
      "method": "top_level_merge",
      "data": {
        "is_google": true
      }
    }
  ]
}
```

```bash
mmdb-cli update -i GeoLite2-ASN.mmdb -o GeoLite2-ASN-Updated.mmdb -d update-data.json
```

To verify the changes, you can use the `inspect` command to check the updated data in the MMDB file.

:::note[Disclaimer]

The data shown in the examples above is for demonstration purposes only and may not reflect the actual data in the MMDB file you are inspecting. The actual data may vary based on the version of the database and the IP addresses queried.

:::

```json
[
    {
        "query": "1.1.1.1",
        "records": [
            {
                "network": "1.1.1.1/32",
                "record": {
                    "attributes": {
                        "city": "San Francisco",
                        "country": "US"
                    },
                    "autonomous_system_number": 13335,
                    "autonomous_system_organization": "CLOUDFLARENET",
                    "is_cloudflare": true
                }
            }
        ]
    },
    {
        "query": "8.8.8.8",
        "records": [
            {
                "network": "8.8.8.8/32",
                "record": {
                    "autonomous_system_number": 15169,
                    "autonomous_system_organization": "GOOGLE",
                    "is_google": true
                }
            }
        ]
    }
]
```
