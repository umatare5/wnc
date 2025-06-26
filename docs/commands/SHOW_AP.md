# 📡 wnc show ap

Display comprehensive information about access points in the wireless infrastructure.

## ✨ Features

- Display all access points across multiple controllers
- Show AP status, location, and configuration details
- Support for both tabular and JSON output formats
- Real-time status information

## 📋 Syntax

```bash
wnc show ap [options...]
```

**Aliases:** `s ap`, `s a`

## ⚙️ Flags

| Flag            | Alias | Type   | Description                       | Default | Required | Environment Variable |
| --------------- | ----- | ------ | --------------------------------- | ------- | -------- | -------------------- |
| `--controllers` | `-c`  | string | Controller-token pairs            | -       | Yes      | `WNC_CONTROLLERS`    |
| `--insecure`    | `-k`  | bool   | Skip TLS certificate verification | `false` | No       | -                    |
| `--format`      | `-f`  | string | Output format: `json`, `table`    | `table` | No       | -                    |
| `--timeout`     | `-t`  | int    | HTTP client timeout in seconds    | `60`    | No       | -                    |

## 📝 Usage

```bash
# List all access points
wnc show ap --controllers "wnc.example.com:token"

# JSON format for scripting
wnc show ap --format json --controllers "wnc.example.com:token"

# Multiple controllers
wnc show ap --controllers "wnc1.example.com:token1,wnc2.example.com:token2"

# Using environment variable
export WNC_CONTROLLERS="wnc.example.com:token"
wnc show ap
```

## 📤 Example Output

### Table Format

```text
$ wnc show ap

┌────────────────────┬───────┬──────────────────┬─────────────┬───────────────────┬───────────────────┬──────────────┬────────┬────────────────┬────────────┬────────────┬─────────────────────────────────────┬──────────────┬────────────┬───────────────────────┐
│ AP Name            │ Slots │ Model            │ Serial      │ Ethernet MAC      │ Radio MAC         │ Country Code │ Domain │ IP Address     │ OS Version │ State      │ LLDP Neighbor                       │ Power Type   │ Power Mode │ Controller            │
├────────────────────┼───────┼──────────────────┼─────────────┼───────────────────┼───────────────────┼──────────────┼────────┼────────────────┼────────────┼────────────┼─────────────────────────────────────┼──────────────┼────────────┼───────────────────────┤
│ lab2-ap1815-06f-02 │ 2     │ AIR-AP1815I-Q-K9 │ 00000000000 │ 28:ac:9e:00:00:00 │ 28:ac:9e:00:00:00 │ J4           │ -Q     │ 192.168.255.11 │ 17.12.5.41 │ registered │ lab2-cat29c-06f-01.labo.local Gi0/2 │ Advanced PoE │ High       │ wnc1.example.internal │
│ lab2-ap9166-06f-01 │ 3     │ CW9166I-Q        │ 00000000000 │ c4:14:a2:00:00:00 │ f0:d8:05:00:00:00 │ J4           │ -Q     │ 192.168.255.12 │ 17.12.5.41 │ registered │ lab2-cat29c-06f-01.labo.local Gi0/3 │ Legacy PoE   │ High       │ wnc1.example.internal │
└────────────────────┴───────┴──────────────────┴─────────────┴───────────────────┴───────────────────┴──────────────┴────────┴────────────────┴────────────┴────────────┴─────────────────────────────────────┴──────────────┴────────────┴───────────────────────┘

```

### JSON Format

```json
[
  {
    "APName": "AP-Floor1-Room101",
    "MACAddress": "aa:bb:cc:dd:ee:01",
    "Model": "C9120AXI",
    "Status": "Up",
    "Location": "Floor1-Room101",
    "IPAddress": "192.168.1.101",
    "Uptime": "5d 12h"
  },
  {
    "APName": "AP-Floor1-Room102",
    "MACAddress": "aa:bb:cc:dd:ee:02",
    "Model": "C9120AXI",
    "Status": "Up",
    "Location": "Floor1-Room102",
    "IPAddress": "192.168.1.102",
    "Uptime": "5d 11h"
  }
]
```

## 📖 Related Commands

- [wnc show ap-tag](SHOW_AP_TAG.md)
- [wnc show overview](SHOW_OVERVIEW.md)
- [wnc show client](SHOW_CLIENT.md)
