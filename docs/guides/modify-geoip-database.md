---
sidebar_position: 2
sidebar_label: Modify a GeoIP Database
description: Step-by-step guide to adding, editing, replacing, and removing records in an existing GeoIP MMDB database
tags:
  - cli
  - mmdb
  - geoip
  - mmdb-cli
  - guides
  - update
---

# Modify a GeoIP Database

This guide walks through updating an existing GeoIP `.mmdb` database with MMDB CLI. You can add a new network record, edit fields on an existing record, replace a record completely, or remove a network from the database.

MMDB files map IP networks to records. When you modify a database, you are changing the record attached to a CIDR network such as `1.1.1.0/24`, `8.8.8.8/32`, or `2001:db8::/32`.

## Prerequisites

- `mmdb-cli` installed ([Installation](../installation))
- An existing `.mmdb` file, for example `GeoLite2-ASN.mmdb`
- Basic familiarity with JSON and CIDR notation

:::tip

Write updates to a new output file first. Keep the original database unchanged until you inspect and verify the updated file.

:::

---

## Step 1 - Inspect the current record

Before changing a network, inspect the current data so you know which fields already exist:

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb 1.1.1.1
```

Example output:

```yaml
- query: 1.1.1.1
  records:
    - network: 1.1.1.0/24
      record:
        autonomous_system_number: 13335
        autonomous_system_organization: CLOUDFLARENET
```

If you want machine-readable output while preparing a change, use JSON:

```bash
mmdb-cli inspect -i GeoLite2-ASN.mmdb -f json 1.1.1.1
```

---

## Step 2 - Create an update file

Create a file named `geoip-updates.json`. The `update` command expects a top-level `dataset` array. Each item describes one operation:

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "method": "deep_merge",
      "data": {
        "is_anycast": true
      }
    }
  ]
}
```

Use `data`, not `record`, in update files. The `record` field is used by datasets for `generate`; update operations use `data`.

If your update data contains integers, add a top-level `schema` so MMDB CLI stores those values with the expected MMDB type. Without a schema, JSON numbers are stored as `float64`.

```json
{
  "version": "v1",
  "schema": {
    "autonomous_system_number": "uint32",
    "geoname_id": "uint32",
    "confidence": "float64"
  },
  "dataset": []
}
```

## Step 3 - Choose the right method

| Method | Use it when you want to |
| :--- | :--- |
| `replace` | Add a new network record or replace the full record for an existing network |
| `deep_merge` | Add or edit nested fields while keeping existing fields |
| `top_level_merge` | Add or edit top-level fields; nested objects under the same top-level key are replaced as a whole |
| `remove` | Remove the target network from the database |

If you are unsure which edit method to use, start with `deep_merge` for partial updates and `replace` for full record changes.

---

## Step 4 - Add a new record

To add a record for a new network, use `replace`. For a single IPv4 address, use `/32`. For a single IPv6 address, use `/128`.

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "203.0.113.10/32",
      "method": "replace",
      "data": {
        "country": {
          "iso_code": "US",
          "names": {
            "en": "United States"
          }
        },
        "custom": {
          "source": "internal",
          "confidence": 95.5
        }
      }
    }
  ]
}
```

Apply the update:

```bash
mmdb-cli update -i GeoLite2-Country.mmdb -o GeoLite2-Country-Updated.mmdb -d geoip-updates.json
```

---

## Step 5 - Edit an existing record

Use `deep_merge` when you want to add or change fields without losing the rest of the record.

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "method": "deep_merge",
      "data": {
        "custom": {
          "is_cloudflare": true,
          "reviewed_by": "network-team"
        }
      }
    }
  ]
}
```

Run the update:

```bash
mmdb-cli update -i GeoLite2-ASN.mmdb -o GeoLite2-ASN-Updated.mmdb -d geoip-updates.json
```

After the update, the original ASN fields remain and the new `custom` fields are added.

---

## Step 6 - Replace a full record

Use `replace` when the new data should become the complete record for the network.

```json
{
  "version": "v1",
  "schema": {
    "autonomous_system_number": "uint32"
  },
  "dataset": [
    {
      "network": "1.1.1.0/24",
      "method": "replace",
      "data": {
        "autonomous_system_number": 13335,
        "autonomous_system_organization": "Cloudflare, Inc.",
        "custom": {
          "is_cloudflare": true
        }
      }
    }
  ]
}
```

Fields that are not included in `data` are removed from that network's record.

---

## Step 7 - Remove a record

Use `remove` to remove the target network from the database. The update schema still requires a `data` object, so pass an empty object.

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "203.0.113.10/32",
      "method": "remove",
      "data": {}
    }
  ]
}
```

Apply the removal:

```bash
mmdb-cli update -i GeoLite2-Country.mmdb -o GeoLite2-Country-Updated.mmdb -d geoip-updates.json
```

:::note

`remove` removes a network record. To remove only one field from a record, use `replace` and provide the complete record without that field.

:::

---

## Step 8 - Apply multiple changes at once

You can include multiple operations in one update file. They are processed in order.

```json
{
  "version": "v1",
  "dataset": [
    {
      "network": "203.0.113.10/32",
      "method": "replace",
      "data": {
        "country": {
          "iso_code": "US"
        },
        "custom": {
          "source": "internal"
        }
      }
    },
    {
      "network": "1.1.1.0/24",
      "method": "deep_merge",
      "data": {
        "custom": {
          "is_cloudflare": true
        }
      }
    },
    {
      "network": "198.51.100.0/24",
      "method": "remove",
      "data": {}
    }
  ]
}
```

Apply all changes:

```bash
mmdb-cli update -i GeoLite2-Country.mmdb -o GeoLite2-Country-Updated.mmdb -d geoip-updates.json
```

---

## Step 9 - Verify and inspect the result

Verify that the updated file is a valid MMDB database:

```bash
mmdb-cli verify -i GeoLite2-Country-Updated.mmdb
```

Expected output:

```text
The MMDB file is valid
```

Inspect every changed network before replacing the original file:

```bash
mmdb-cli inspect -i GeoLite2-Country-Updated.mmdb 203.0.113.10 1.1.1.1 198.51.100.1
```

An IP address with no matching network returns an empty `records` list. If the removed network was inside a broader network that still exists, lookup may return that broader network instead.

---

## Common mistakes

| Problem | Fix |
| :--- | :--- |
| The update file uses `record` | Use `data` for update operations |
| A single IP was updated without a prefix | Use `/32` for IPv4 or `/128` for IPv6 |
| Existing fields disappeared | Use `deep_merge` for partial updates instead of `replace` |
| Integer fields became floats | Add a top-level `schema` for numeric fields |
| A field should be deleted, not the whole network | Use `replace` with the complete desired record |
| The output file overwrote expectations | Write to a new filename, verify it, then promote it |

## Summary

| Task | Method | Example |
| :--- | :--- | :--- |
| Add a new record | `replace` | `203.0.113.10/32` with full `data` |
| Edit part of a record | `deep_merge` | Add `custom.is_cloudflare` |
| Edit top-level fields only | `top_level_merge` | Add `is_datacenter` |
| Replace the full record | `replace` | Provide the complete desired record |
| Remove a network record | `remove` | Use `data: {}` |

## See also

- [Update Command](../commands/update) - full command reference
- [Inspect Command](../commands/inspect) - query records before and after updates
- [Verify Command](../commands/verify) - validate the updated MMDB file
- [Dataset Schema](../dataset-schema) - JSON schema reference
