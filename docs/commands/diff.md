---
sidebar_position: 11
sidebar_label: Diff
description: MMDB CLI Diff Command
tags:
  - cli
  - mmdb
  - mmdb-cli
  - geoip
---

# Diff Command

The `diff` command compares two MMDB files and reports **added**, **removed**, and **modified** networks. It is useful when validating database updates, reviewing changes between builds, or checking that a custom MMDB matches an expected baseline.

Output uses a unified-diff-inspired format on stdout. The command compares **network records** in each database (not MMDB metadata such as build epoch or node count).

## Usage

```bash
mmdb-cli diff <FILE_A.mmdb> <FILE_B.mmdb> [IP_OR_CIDR ...]
```

| Argument | Description |
| :------- | :---------- |
| `FILE_A.mmdb` | Baseline database (shown as `---` in output). |
| `FILE_B.mmdb` | Comparison database (shown as `+++` in output). |
| `IP_OR_CIDR` | Optional scope. When omitted, every network in both files is compared. When provided, only networks within the given IP addresses or CIDR prefixes are included. |

Both file paths must exist and use the `.mmdb` extension.

## Exit codes

| Code | Meaning |
| :--- | :------ |
| `0` | No differences found. |
| `1` | Differences exist, or a fatal error occurred. |

Use the exit code in CI or scripts to fail when two databases are not equivalent:

```bash
mmdb-cli diff old.mmdb new.mmdb || echo "Databases differ"
```

## Output format

Each run prints a header naming the two files, then one of:

- `[+] No differences found` when the scoped comparison is identical.
- A summary line with counts, followed by change lines:

| Prefix | Meaning |
| :----- | :------ |
| `-` | Network present in **file A** but not in **file B** (removed). |
| `+` | Network present in **file B** but not in **file A** (added). |
| `~` | Network exists in both files with a different record (modified). Modified entries show the before (`-`) and after (`+`) record as JSON. |

Networks are sorted before printing. Record values are JSON-encoded on a single line.

## Example — No differences

Compare a file to itself:

```bash
mmdb-cli diff GeoLite2-ASN.mmdb GeoLite2-ASN.mmdb
```

### Output

```text
--- GeoLite2-ASN.mmdb
+++ GeoLite2-ASN.mmdb
[+] No differences found
```

## Example — Full database diff

Compare two versions of a database:

```bash
mmdb-cli diff GeoLite2-ASN-old.mmdb GeoLite2-ASN-new.mmdb
```

### Output (illustrative)

:::note[Disclaimer]

Example records are illustrative.

:::

```text
--- GeoLite2-ASN-old.mmdb
+++ GeoLite2-ASN-new.mmdb
@@ Networks: 1 added, 0 removed, 2 modified @@
~ 1.0.0.0/24
  - {"autonomous_system_number":13335,"autonomous_system_organization":"CLOUDFLARENET"}
  + {"autonomous_system_number":13335,"autonomous_system_organization":"Cloudflare, Inc."}
+ 2.0.0.0/24 {"autonomous_system_number":8075,"autonomous_system_organization":"MICROSOFT-CORP-MSN-AS-BLOCK"}
```

## Example — Scoped to a CIDR

Limit the comparison to networks inside a prefix:

```bash
mmdb-cli diff file-a.mmdb file-b.mmdb 1.0.0.0/8
```

Only networks returned by `NetworksWithin` for that CIDR are compared. Changes outside the scope are ignored.

## Example — Scoped to an IP address

A bare IPv4 or IPv6 address is treated as a host route (`/32` or `/128`):

```bash
mmdb-cli diff file-a.mmdb file-b.mmdb 1.1.1.1
```

## Example — Multiple scopes

Pass several IPs or CIDRs to union the scopes:

```bash
mmdb-cli diff file-a.mmdb file-b.mmdb 1.1.1.1 8.8.8.0/24
```

Overlapping scopes are deduplicated by network.

## What is compared

- **Compared:** Each network prefix and its decoded record from the MMDB search tree.
- **Not compared:** Database metadata (for example `BuildEpoch`, `DatabaseType`, `NodeCount`).
- **Equality:** Records are compared by value. Two records with the same data but different internal representation may still appear as modified.

For a full export of either database, use the [`dump`](./dump) command. To look up specific addresses without comparing two files, use [`inspect`](./inspect).

## Notes

- **Memory:** A full diff (no scope filters) loads all matching networks and records from both files into memory. For very large production databases, prefer scoped diffs or export-and-compare workflows.
- **Tree structure:** MMDB stores data in a search tree. Two databases with equivalent lookup results but different internal prefix splits may report more modified networks than expected.
- **Output:** The format is human-readable and inspired by unified diff notation; it is not intended for use with `patch`.
