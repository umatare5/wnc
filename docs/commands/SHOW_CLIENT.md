# 📱 wnc show client

Display detailed information about wireless clients connected to the infrastructure.

## ✨ Features

- Real-time client connectivity information
- RF metrics including RSSI, SNR, and throughput
- Traffic statistics (RX/TX bytes)
- Advanced filtering by radio band and SSID
- Flexible sorting options

## 📋 Syntax

```bash
wnc show client [options...]
```

**Aliases:** `s client`, `s c`

## ⚙️ Flags

| Flag            | Alias | Type   | Description                                                                                | Default     | Required | Environment Variable |
| --------------- | ----- | ------ | ------------------------------------------------------------------------------------------ | ----------- | -------- | -------------------- |
| `--controllers` | `-c`  | string | Controller-token pairs                                                                     | -           | Yes      | `WNC_CONTROLLERS`    |
| `--insecure`    | `-k`  | bool   | Skip TLS certificate verification                                                          | `false`     | No       | -                    |
| `--format`      | `-f`  | string | Output format: `json`, `table`                                                             | `table`     | No       | -                    |
| `--timeout`     | `-t`  | int    | HTTP client timeout in seconds                                                             | `60`        | No       | -                    |
| `--radio`       | `-r`  | string | Radio filter: `0` (2.4GHz), `1` (5GHz), `2` (5GHz/6GHz)                                    | -           | No       | -                    |
| `--ssid`        | `-s`  | string | ESSID name to filter results                                                               | -           | No       | -                    |
| `--sort-by`     | `-b`  | string | Sort field: `Hostname`, `IPAddress`, `RSSI`, `SNR`, `Throughput`, `RxTraffic`, `TxTraffic` | `IPAddress` | No       | -                    |
| `--sort-order`  | `-o`  | string | Sort order: `asc`, `desc`                                                                  | `desc`      | No       | -                    |

## 📝 Usage

```bash
# List all wireless clients
wnc show client --controllers "wnc.example.com:token"

# JSON format for scripting
wnc show client --format json --controllers "wnc.example.com:token"

# Filter by 5GHz radio only
wnc show client --controllers "wnc.example.com:token" --radio 1

# Filter by specific SSID
wnc show client --controllers "wnc.example.com:token" --ssid "CorpWiFi"

# Sort by signal strength (strongest first)
wnc show client --controllers "wnc.example.com:token" --sort-by RSSI --sort-order desc
```

## 📤 Example Output

### Table Format

```text
$ wnc show client

┌───────────────────┬──────────────┬───────────────────────────────┬──────────┬──────────┬──────────┬────────┬───────┬────────────┬─────────┬───────┬───────────┬───────────────┬───────────────┬────────────────────┬───────────────────────┐
│ MACAddress        │ IPAddress    │ Hostname                      │ Username │ SSID     │ Protocol │ Band   │ State │ Throughput │ RSSI    │ SNR   │ Stream    │ RxTraffic     │ TxTraffic     │ APName             │ Controller            │
├───────────────────┼──────────────┼───────────────────────────────┼──────────┼──────────┼──────────┼────────┼───────┼────────────┼─────────┼───────┼───────────┼───────────────┼───────────────┼────────────────────┼───────────────────────┤
│ cc:50:e3:00:00:00 │ 192.168.0.98 │ Unknown Device                │ N/A      │ labo1    │ 11n      │ 2.4GHz │ Run   │ 48 Mbps    │ -47 dBm │ 52 dB │ 0 Streams │ 46 KB         │ 4 KB          │ lab2-ap9166-06f-01 │ wnc1.example.internal │
│ 6c:b1:33:00:00:00 │ 192.168.0.96 │ MacBook Pro (14-inch, 2021)   │ N/A      │ labo3    │ dot11ax  │ 5GHz   │ Run   │ 516 Mbps   │ -57 dBm │ 36 dB │ 2 Streams │ 80,504 KB     │ 416,644 KB    │ lab2-ap9166-06f-01 │ wnc1.example.internal │
│ 50:d4:f7:00:00:00 │ 192.168.0.75 │ TP-LINK TECHNOLOGIES CO.,LTD. │ N/A      │ labo2    │ 11n      │ 2.4GHz │ Run   │ 72 Mbps    │ -49 dBm │ 50 dB │ 1 Streams │ 273 KB        │ 127 KB        │ lab2-ap9166-06f-01 │ wnc1.example.internal │
│ 0e:92:1c:00:00:00 │ 192.168.0.62 │ iPad Pro 3rd Gen (11 inch)    │ N/A      │ labo1    │ 11n      │ 2.4GHz │ Run   │ 144 Mbps   │ -25 dBm │ 74 dB │ 2 Streams │ 47,604 KB     │ 120,770 KB    │ lab2-ap9166-06f-01 │ wnc1.example.internal │
└───────────────────┴──────────────┴───────────────────────────────┴──────────┴──────────┴──────────┴────────┴───────┴────────────┴─────────┴───────┴───────────┴───────────────┴───────────────┴────────────────────┴───────────────────────┘

```

### JSON Format

```json
[
  {
    "Hostname": "laptop-01",
    "IPAddress": "192.168.1.10",
    "MACAddress": "aa:bb:cc:11:22:33",
    "SSID": "CorpWiFi",
    "Radio": "1",
    "RSSI": "-45",
    "SNR": "35",
    "APName": "AP-Floor1-Room101",
    "Throughput": "150 Mbps",
    "RxTraffic": "2.1 GB",
    "TxTraffic": "450 MB"
  },
  {
    "Hostname": "phone-02",
    "IPAddress": "192.168.1.15",
    "MACAddress": "dd:ee:ff:44:55:66",
    "SSID": "CorpWiFi",
    "Radio": "0",
    "RSSI": "-55",
    "SNR": "25",
    "APName": "AP-Floor1-Room102",
    "Throughput": "54 Mbps",
    "RxTraffic": "1.2 GB",
    "TxTraffic": "200 MB"
  }
]
```

## 📖 Related Commands

- [wnc show ap](SHOW_AP.md)
- [wnc show overview](SHOW_OVERVIEW.md)
- [wnc show wlan](SHOW_WLAN.md)
