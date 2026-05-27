---
sidebar_position: 1
sidebar_label: JSONPath output format
description: Use JSONPath templates with the -f flag on inspect and dump
tags:
  - cli
  - mmdb
  - mmdb-cli
  - jsonpath
---

# JSONPath output format

Some commands accept a **format** value of `jsonpath={TEMPLATE}` on the `-f` / `--format` flag. This is not a separate command—it is an alternative output mode for:

| Command   | Default `-f` | JSONPath value        |
| :-------- | :----------- | :-------------------- |
| `inspect` | `yaml`       | `jsonpath={TEMPLATE}` |
| `dump`    | `json`       | `jsonpath={TEMPLATE}` |

`metadata`, `generate`, `update`, and `verify` do not support this format.

Templates use the [Kubernetes JSONPath engine](https://kubernetes.io/docs/reference/kubectl/jsonpath/) (Go templates in `{...}` braces), not the [Goessner JSONPath](https://goessner.net/articles/JsonPath/) query syntax.

```bash
mmdb-cli inspect -i db.mmdb \
  -f 'jsonpath={range .items[*]}{.network}{"\n"}{end}' \
  1.1.1.1

mmdb-cli dump -i db.mmdb -o networks.txt \
  -f 'jsonpath={range .items[*]}{.network}{"\n"}{end}'
```

:::tip[Shell quoting]

Wrap the entire `-f` value in single quotes so the shell does not expand `$`, `` ` ``, or `?` inside the template.

:::

## Behavior

- **`inspect`**: Rendered output goes to stdout (not wrapped in YAML/JSON/XML/CSV).
- **`dump`**: Rendered bytes are written to `-o`. The file does **not** need a `.json` extension.
- Invalid templates fail **before** iteration starts.

## Syntax (quick reference)

| Construct | Example |
| :-------- | :------ |
| Field access | `{.metadata.NodeCount}` |
| Wildcard | `{.items[*].network}` |
| `range` loop | `{range .items[*]}{.network}{"\n"}{end}` |
| Filter in `range` | `{range .items[?(@.record.country.iso_code=="AU")]}{.network}{end}` |
| Missing key | `{.items[0].missing}` → empty string, no error |

See [kubectl JSONPath](https://kubernetes.io/docs/reference/kubectl/jsonpath/) for the full grammar.

## Root data: `inspect`

```json
{
  "apiVersion": "mmdb-cli/v1",
  "kind": "InspectList",
  "items": [
    { "query": "1.1.1.1", "network": "1.1.1.0/24", "record": { } }
  ],
  "queries": [
    { "query": "1.1.1.1", "records": [ { "network": "1.1.1.0/24", "record": { } } ] }
  ]
}
```

Use **`items`** for flat `range` loops and filters. **`queries`** mirrors normal inspect output (grouped by input).

## Root data: `dump`

```json
{
  "apiVersion": "mmdb-cli/v1",
  "kind": "DumpList",
  "metadata": { },
  "items": [ { "network": "1.1.1.0/24", "record": { } } ],
  "dataset": [ ]
}
```

**`items`** and **`dataset`** are the same list. **`metadata`** uses Go field names (`NodeCount`, `DatabaseType`, …).

## Examples

Filter Australian networks during dump:

```bash
mmdb-cli dump -i GeoLite2-Country.mmdb -o au.txt \
  -f 'jsonpath={range .items[?(@.record.country.iso_code=="AU")]}{.network}{"\n"}{end}'
```

Inspect: network and country code, tab-separated:

```bash
mmdb-cli inspect -i GeoLite2-Country.mmdb \
  -f 'jsonpath={range .items[*]}{.network}{"\t"}{.record.country.iso_code}{"\n"}{end}' \
  1.1.1.1/24
```

Command-specific usage: [inspect](../commands/inspect), [dump](../commands/dump). Record field shapes: [dataset schema](../dataset-schema).
