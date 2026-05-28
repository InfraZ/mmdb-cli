---
sidebar_position: 1
sidebar_label: Create a Custom MMDB from Scratch
description: Step-by-step guide to creating a custom MMDB database using MMDB CLI
tags:
  - cli
  - mmdb
  - geoip
  - mmdb-cli
  - guides
  - generate
---

# Create a Custom MMDB from Scratch

This guide walks you through building a custom MMDB file from scratch — from designing your dataset to verifying the final output. By the end you will have a working `.mmdb` file you can query with any MaxMind-compatible library or tool.

## Prerequisites

- `mmdb-cli` installed ([Installation](../installation))
- Basic familiarity with JSON and CIDR notation

---

## Step 1 — Decide what data to store

An MMDB file maps IP networks (CIDR ranges) to arbitrary records. Before writing any JSON, answer these questions:

| Question | Example answer |
| :--- | :--- |
| What IP networks will I cover? | `1.1.1.0/24`, `10.0.0.0/8` |
| What fields does each record need? | ASN number, organization name, threat score |
| Do I need strict integer/float types? | Yes → use a `schema`; No → let MMDB CLI infer |
| IPv4-only, or dual-stack (IPv4 + IPv6)? | Dual-stack → `"IPVersion": 6` (recommended default) |

---

## Step 2 — Create the dataset JSON file

The dataset file is the single source of truth for your MMDB. Create a file named `my-database.json` with the following top-level structure:

```json
{
  "version": "v1",
  "schema": {},
  "metadata": {},
  "dataset": []
}
```

### 2a — Fill in the `metadata`

The `metadata` block controls how the MMDB file identifies itself. Every field shown below is required by the `generate` command:

```json
"metadata": {
  "DatabaseType": "My-Custom-DB",
  "Description": {
    "en": "My custom IP database"
  },
  "Languages": ["en"],
  "IPVersion": 6,
  "RecordSize": 24
}
```

| Field | Notes |
| :--- | :--- |
| `DatabaseType` | Free-form string — choose a name that is meaningful to your application. |
| `Description` | Map of language code → description text. |
| `IPVersion` | Use `6` for dual-stack (supports both IPv4 and IPv6 lookups). Use `4` for IPv4-only. |
| `RecordSize` | `24` is the most compact option and sufficient for most custom databases. |

### 2b — Add a `schema` (optional but recommended)

Without a `schema`, all JSON numbers are stored as `float64`. If your records contain integer fields (ASN numbers, port numbers, counts), add a `schema` to enforce the correct types:

```json
"schema": {
  "asn": "uint32",
  "organization": "string",
  "is_proxy": "bool",
  "threat_score": "float64"
}
```

Supported type keywords: `string`, `bool`, `float64`, `int32`, `uint32`, `uint16`.

Nested objects are supported:

```json
"schema": {
  "geo": {
    "latitude": "float64",
    "longitude": "float64"
  }
}
```

### 2c — Add records to `dataset`

Each entry in `dataset` maps one CIDR range to a record object:

```json
"dataset": [
  {
    "network": "1.1.1.0/24",
    "record": {
      "asn": 13335,
      "organization": "CLOUDFLARENET",
      "is_proxy": false,
      "threat_score": 0.1
    }
  },
  {
    "network": "8.8.8.0/24",
    "record": {
      "asn": 15169,
      "organization": "GOOGLE",
      "is_proxy": false,
      "threat_score": 0.0
    }
  },
  {
    "network": "185.220.101.0/24",
    "record": {
      "asn": 205100,
      "organization": "F3 Netze e.V.",
      "is_proxy": true,
      "threat_score": 8.5
    }
  }
]
```

### Complete example

```json
{
  "version": "v1",
  "schema": {
    "asn": "uint32",
    "organization": "string",
    "is_proxy": "bool",
    "threat_score": "float64"
  },
  "metadata": {
    "DatabaseType": "My-Custom-DB",
    "Description": {
      "en": "My custom IP database"
    },
    "Languages": ["en"],
    "IPVersion": 6,
    "RecordSize": 24
  },
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "record": {
        "asn": 13335,
        "organization": "CLOUDFLARENET",
        "is_proxy": false,
        "threat_score": 0.1
      }
    },
    {
      "network": "8.8.8.0/24",
      "record": {
        "asn": 15169,
        "organization": "GOOGLE",
        "is_proxy": false,
        "threat_score": 0.0
      }
    },
    {
      "network": "185.220.101.0/24",
      "record": {
        "asn": 205100,
        "organization": "F3 Netze e.V.",
        "is_proxy": true,
        "threat_score": 8.5
      }
    }
  ]
}
```

---

## Step 3 — Generate the MMDB file

Run the `generate` command, pointing it at your dataset file and choosing an output path:

```bash
mmdb-cli generate -i my-database.json -o my-database.mmdb
```

On success the command exits with code `0` and writes the binary MMDB file to `my-database.mmdb`. You can add `-v` to see verbose output:

```bash
mmdb-cli generate -i my-database.json -o my-database.mmdb -v
```

:::tip[IPv4-only datasets]

If every network in your dataset is IPv4 and you do **not** need IPv6 lookups, pass `--disable-ipv4-aliasing` to produce a leaner file:

```bash
mmdb-cli generate -i my-database.json -o my-database.mmdb --disable-ipv4-aliasing
```

:::

---

## Step 4 — Verify the file

Check that the generated file is a valid MMDB binary:

```bash
mmdb-cli verify -i my-database.mmdb
```

Expected output:

```text
The MMDB file is valid
```

If the file is corrupted or malformed, the command will exit with a non-zero code and print an error.

---

## Step 5 — Inspect and test lookups

Use the `inspect` command to query the new file and confirm your records were written correctly:

```bash
mmdb-cli inspect -i my-database.mmdb 1.1.1.1
```

Expected output (YAML):

```yaml
- query: 1.1.1.1
  records:
    - network: 1.1.1.0/24
      record:
        asn: 13335
        is_proxy: false
        organization: CLOUDFLARENET
        threat_score: 0.1
```

Test a few more IPs to cover each network you added:

```bash
mmdb-cli inspect -i my-database.mmdb 8.8.8.8 185.220.101.42
```

For a range query, pass a CIDR prefix:

```bash
mmdb-cli inspect -i my-database.mmdb 1.1.1.0/24
```

:::note

An IP address that falls outside every network in your dataset will return an empty `records` list — this is expected behavior.

:::

---

## Step 6 — Iterate and update

Need to add or change records without regenerating from scratch? Use the `update` command to merge changes into the existing file. See [Update Command](../commands/update) for details.

---

## Summary

| Step | Command |
| :--- | :--- |
| 1. Design your schema | _(plan on paper)_ |
| 2. Write `my-database.json` | _(text editor)_ |
| 3. Generate the MMDB | `mmdb-cli generate -i my-database.json -o my-database.mmdb` |
| 4. Verify integrity | `mmdb-cli verify -i my-database.mmdb` |
| 5. Test lookups | `mmdb-cli inspect -i my-database.mmdb <IP>` |
| 6. Update records | `mmdb-cli update ...` |

## See also

- [Dataset Schema](../dataset-schema) — full reference for the JSON format
- [Generate Command](../commands/generate) — all flags and options
- [Inspect Command](../commands/inspect) — output formats and CIDR queries
- [Verify Command](../commands/verify) — integrity checking
- [Update Command](../commands/update) — incremental updates to an existing MMDB
