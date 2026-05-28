---
sidebar_position: 5
sidebar_label: Dataset Schema
description: MMDB CLI Dataset Schema
tags:
  - cli
  - mmdb
  - mmdb-cli
  - dataset
  - geoip
---

# Dataset Schema

MMDB CLI uses JSON dataset files for `generate`, `update`, and `dump` workflows.

## Top-level Structure

The dataset JSON supports these top-level fields:

| Field | Type | Description | Required |
| :---: | :---: | :--- | :---: |
| `version` | string | Dataset format version. Current supported value: `v1`. | No |
| `schema` | object | Optional dynamic schema object used by `generate`/`update`. | No |
| `metadata` | object | MMDB metadata used by `generate` (and included in `dump` output). | `generate`: Yes |
| `dataset` | array | Records or update operations (depends on command usage). | Yes |

```json
{
  "version": "v1",
  "schema": {},
  "metadata": {},
  "dataset": []
}
```

## `version` Field

If `version` is provided, it must be `v1`.

## `metadata` Field (for `generate`)

The `metadata` object defines MMDB writer options when generating a database.

| Field | Type | Description | Required |
| :---: | :---: | :--- | :---: |
| `DatabaseType` | string | Database type (for example `GeoLite2-ASN`). | Yes |
| `Description` | object | Localized description map (for example `{"en":"..."}`). | Yes |
| `IPVersion` | integer | Allowed values: `4`, `6`. Defaults to `6` if omitted. | No |
| `Languages` | array | Language list. Defaults to `["en"]` if omitted. | No |
| `RecordSize` | integer | Allowed values: `24`, `28`, `32`. Defaults to `28`. | No |

## `dataset` for `generate` / `dump`

For `generate`, each item in `dataset` must contain:

| Field | Type | Description | Required |
| :---: | :---: | :--- | :---: |
| `network` | string | Network address in CIDR notation (IPv4 or IPv6). | Yes |
| `record` | object | The record data associated with the network | Yes |

Example:

```json
{
  "network": "1.1.1.0/24",
  "record": {
    "autonomous_system_number": 13335,
    "autonomous_system_organization": "CLOUDFLARENET"
  }
}
```

## Supported Record Value Types

Record values are converted to MMDB types by MMDB CLI before writing.

### Default type inference (no `schema` provided)

When you do not provide a `schema`, MMDB CLI infers types from JSON values:

| JSON value | Stored MMDB type |
| :--- | :--- |
| string | `string` |
| boolean | `bool` |
| number | `float64` |
| object | map/object (recursively converted) |

Notes:

- JSON numbers are typically parsed as floating-point values in this mode.
- Nested objects are supported recursively.
- Arrays and `null` are not supported by the current default converter.

### Schema-driven typing (`schema` provided)

If you provide a top-level `schema` object, you can force specific types per field.

Supported schema type keywords:

- `string`
- `bool` / `boolean`
- `float` / `float64`
- `int` / `int32`
- `uint` / `uint32`
- `uint16`

Schema can be nested for nested record objects:

```json
{
  "schema": {
    "autonomous_system_number": "uint32",
    "autonomous_system_organization": "string",
    "is_cloudflare": "bool",
    "attributes": {
      "score": "float64"
    }
  }
}
```

If a value does not match the declared schema type, MMDB CLI falls back to a default value for that type (for example `0`, `false`, or empty string).

### End-to-end Examples

Example 1 - default inference (no `schema`):

```json
{
  "version": "v1",
  "metadata": {
    "DatabaseType": "GeoLite2-ASN-Custom",
    "Description": {
      "en": "Custom ASN Dataset"
    },
    "IPVersion": 6,
    "RecordSize": 24
  },
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "record": {
        "autonomous_system_number": 13335,
        "autonomous_system_organization": "CLOUDFLARENET",
        "is_anycast": true,
        "confidence": 99.5,
        "attributes": {
          "region": "global"
        }
      }
    }
  ]
}
```

Example 2 - schema-driven typing:

```json
{
  "version": "v1",
  "schema": {
    "autonomous_system_number": "uint32",
    "autonomous_system_organization": "string",
    "is_anycast": "bool",
    "confidence": "float64",
    "port": "uint16",
    "attributes": {
      "priority": "int32"
    }
  },
  "metadata": {
    "DatabaseType": "GeoLite2-ASN-Custom",
    "Description": {
      "en": "Typed ASN Dataset"
    },
    "IPVersion": 6,
    "RecordSize": 24
  },
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "record": {
        "autonomous_system_number": 13335,
        "autonomous_system_organization": "CLOUDFLARENET",
        "is_anycast": true,
        "confidence": 99.5,
        "port": 443,
        "attributes": {
          "priority": -1
        }
      }
    }
  ]
}
```

## `dataset` for `update`

For `update`, each item in `dataset` must contain:

| Field | Type | Description | Required |
| :---: | :---: | :--- | :---: |
| `network` | string | Target CIDR to update | Yes |
| `method` | string | One of `remove`, `replace`, `top_level_merge`, `deep_merge` | No (defaults to `deep_merge`) |
| `data` | object | Payload used by update method | Yes |

Example:

```json
{
  "network": "1.1.1.1/32",
  "method": "deep_merge",
  "data": {
    "is_cloudflare": true
  }
}
```
