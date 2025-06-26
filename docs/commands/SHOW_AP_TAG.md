# 🏷️ wnc show ap-tag

Display access point tag configurations and assignments.

## ✨ Features

- Display AP tag policies and assignments
- Show RF profile mappings
- Policy inheritance information
- Support for both tabular and JSON output formats

## 📋 Syntax

```bash
wnc show ap-tag [options...]
```

## ⚙️ Flags

| Flag            | Alias | Type   | Description                       | Default | Required | Environment Variable |
| --------------- | ----- | ------ | --------------------------------- | ------- | -------- | -------------------- |
| `--controllers` | `-c`  | string | Controller-token pairs            | -       | Yes      | `WNC_CONTROLLERS`    |
| `--insecure`    | `-k`  | bool   | Skip TLS certificate verification | `false` | No       | -                    |
| `--format`      | `-f`  | string | Output format: `json`, `table`    | `table` | No       | -                    |
| `--timeout`     | `-t`  | int    | HTTP client timeout in seconds    | `60`    | No       | -                    |

## 📝 Usage

```bash
# List all AP tags
wnc show ap-tag --controllers "wnc.example.com:token"

# JSON format for scripting
wnc show ap-tag --format json --controllers "wnc.example.com:token"

# Multiple controllers
wnc show ap-tag --controllers "wnc1.example.com:token1,wnc2.example.com:token2"
```

## 📤 Example Output

### Table Format

```text
$ wnc show ap-tag

┌────────────────────┬────────┬─────────────────┬──────────────┬────────────────┬─────────────┬──────────────┬───────────────────┐
│ AP Name            │ Config │ Policy Tag Name │ RF Tag Name  │ Site Tag Name  │ AP Profile  │ Flex Profile │ Tag Source        │
├────────────────────┼────────┼─────────────────┼──────────────┼────────────────┼─────────────┼──────────────┼───────────────────┤
│ lab2-ap1815-06f-02 │   ✅️   │ labo-wlan-flex  │ labo-outside │ labo-site-flex │ labo-common │ labo-flex    │ tag-source-static │
│ lab2-ap9166-06f-01 │   ✅️   │ labo-wlan-flex  │ labo-inside  │ labo-site-flex │ labo-common │ labo-flex    │ tag-source-static │
└────────────────────┴────────┴─────────────────┴──────────────┴────────────────┴─────────────┴──────────────┴───────────────────┘
```

### JSON Format

```json
$ wnc show ap-tag --format json

[
  {
    "TagName": "default-ap-tag",
    "PolicyProfile": "default",
    "RFProfile": "default",
    "SiteTag": "default",
    "Description": "Default AP tag"
  },
  {
    "TagName": "floor1-tag",
    "PolicyProfile": "floor1-policy",
    "RFProfile": "high-dense",
    "SiteTag": "building",
    "Description": "Floor 1 APs"
  }
]
```

## 📖 Related Commands

- [wnc show ap](SHOW_AP.md)
- [wnc show overview](SHOW_OVERVIEW.md)
- [wnc show wlan](SHOW_WLAN.md)
