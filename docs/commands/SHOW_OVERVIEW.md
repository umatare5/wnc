# 📊 wnc show overview

Display a comprehensive overview of the wireless infrastructure with AP and client statistics.

## ✨ Features

- High-level infrastructure overview
- Per-AP client count and RF information
- Channel utilization and power levels
- Radio band filtering
- Executive summary view

## 📋 Syntax

```bash
wnc show overview [options...]
```

**Aliases:** `s overview`, `s o`

## ⚙️ Flags

| Flag            | Alias | Type   | Description                                                        | Default  | Required | Environment Variable |
| --------------- | ----- | ------ | ------------------------------------------------------------------ | -------- | -------- | -------------------- |
| `--controllers` | `-c`  | string | Controller-token pairs                                             | -        | Yes      | `WNC_CONTROLLERS`    |
| `--insecure`    | `-k`  | bool   | Skip TLS certificate verification                                  | `false`  | No       | -                    |
| `--format`      | `-f`  | string | Output format: `json`, `table`                                     | `table`  | No       | -                    |
| `--timeout`     | `-t`  | int    | HTTP client timeout in seconds                                     | `60`     | No       | -                    |
| `--radio`       | `-r`  | string | Radio filter: `0` (2.4GHz), `1` (5GHz), `2` (5GHz/6GHz)            | -        | No       | -                    |
| `--sort-by`     | `-b`  | string | Sort field: `APName`, `APMac`, `Channel`, `ClientCount`, `TxPower` | `APName` | No       | -                    |
| `--sort-order`  | `-o`  | string | Sort order: `asc`, `desc`                                          | `desc`   | No       | -                    |

## 📝 Usage

```bash
# List infrastructure overview
wnc show overview --controllers "wnc.example.com:token"

# JSON format for scripting
wnc show overview --format json --controllers "wnc.example.com:token"

# Filter by 5GHz radio
wnc show overview --controllers "wnc.example.com:token" --radio 1

# Sort by client count (busiest APs first)
wnc show overview --controllers "wnc.example.com:token" --sort-by ClientCount --sort-order desc

# Sort by AP name alphabetically
wnc show overview --controllers "wnc.example.com:token" --sort-by APName --sort-order asc
```

## 📤 Example Output

### Table Format

```text
$ wnc show overview

❯ ./wnc show overview --insecure
┌────────────────────┬───────────────────┬───────┬────────┬────────────────┬─────────┬─────────────┬────────────────────┬─────────────────────┬───────────────────────┐
│ APName             │ APMac             │ Radio │ Status │ Channel        │ TxPower │ ClientCount │ ChannelUtilization │ RFTagName           │ Controller            │
├────────────────────┼───────────────────┼───────┼────────┼────────────────┼─────────┼─────────────┼────────────────────┼─────────────────────┼───────────────────────┤
│ lab2-ap9166-06f-01 │ f0:d8:05:2c:41:20 │ 2     │   ✅️   │ 40 MHz (5,1)   │ 22 dBm  │ 0 clients   │ [          ] 1%    │ labo-rf-6gh         │ wnc1.example.internal │
│ lab2-ap9166-06f-01 │ f0:d8:05:2c:41:20 │ 1     │   ✅️   │ 40 MHz (64,60) │ 18 dBm  │ 2 clients   │ [          ] 2%    │ labo-rf-5gh-inside  │ wnc1.example.internal │
│ lab2-ap9166-06f-01 │ f0:d8:05:2c:41:20 │ 0     │   ✅️   │ 20 MHz (11)    │ 22 dBm  │ 16 clients  │ [          ] 7%    │ labo-rf-24gh        │ wnc1.example.internal │
│ lab2-ap1815-06f-02 │ 28:ac:9e:bb:3c:80 │ 1     │   ✅️   │ 40 MHz (36,40) │ 18 dBm  │ 3 clients   │ [####      ] 42%   │ labo-rf-5gh-outside │ wnc1.example.internal │
│ lab2-ap1815-06f-02 │ 28:ac:9e:bb:3c:80 │ 0     │   ✅️   │ 20 MHz (1)     │ 20 dBm  │ 0 clients   │ [#         ] 19%   │ labo-rf-24gh        │ wnc1.example.internal │
└────────────────────┴───────────────────┴───────┴────────┴────────────────┴─────────┴─────────────┴────────────────────┴─────────────────────┴───────────────────────┘
```

### JSON Format

```json
[
  {
    "APName": "AP-Floor1-Room101",
    "APMac": "aa:bb:cc:dd:ee:01",
    "Channel": "36",
    "ClientCount": "12",
    "TxPower": "17 dBm",
    "Radio": "1",
    "Status": "Up"
  },
  {
    "APName": "AP-Floor1-Room102",
    "APMac": "aa:bb:cc:dd:ee:02",
    "Channel": "44",
    "ClientCount": "8",
    "TxPower": "17 dBm",
    "Radio": "1",
    "Status": "Up"
  }
]
```

## 📖 Related Commands

- [wnc show ap](SHOW_AP.md)
- [wnc show client](SHOW_CLIENT.md)
- [wnc show wlan](SHOW_WLAN.md)
