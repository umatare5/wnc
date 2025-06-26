# 📶 wnc show wlan

Display information about WLANs (ESSIDs) configured in the wireless infrastructure.

## ✨ Features

- Display all configured WLANs/ESSIDs
- Show WLAN status and configuration
- Security policy information
- VLAN assignments

## 📋 Syntax

```bash
wnc show wlan [options...]
```

**Aliases:** `s wlan`, `s w`

## ⚙️ Flags

| Flag            | Alias | Type   | Description                       | Default | Required | Environment Variable |
| --------------- | ----- | ------ | --------------------------------- | ------- | -------- | -------------------- |
| `--controllers` | `-c`  | string | Controller-token pairs            | -       | Yes      | `WNC_CONTROLLERS`    |
| `--insecure`    | `-k`  | bool   | Skip TLS certificate verification | `false` | No       | -                    |
| `--format`      | `-f`  | string | Output format: `json`, `table`    | `table` | No       | -                    |
| `--timeout`     | `-t`  | int    | HTTP client timeout in seconds    | `60`    | No       | -                    |

## 📝 Usage

```bash
# List all WLANs
wnc show wlan --controllers "wnc.example.com:token"

# JSON format for scripting
wnc show wlan --format json --controllers "wnc.example.com:token"

# Using environment variable
export WNC_CONTROLLERS="wnc.example.com:token"
wnc show wlan
```

## 📤 Example Output

### Table Format

```text
$ wnc show wlan

┌────────┬──────────┬────┬───────────────────┬──────────────┬─────────────────┬───────────────┬────────────┬─────────────┬──────────────┬─────────────────────┬─────────────────┬──────────────┬─────────────┬───────────┬────────────────┬───────────────────────┐
│ Status │ ESSID    │ ID │ Profile Name      │ VLAN         │ Session Timeout │ DHCP Required │ Egress QoS │ Ingress QoS │ ATF Policies │ Auth Key Management │ mDNS Forwarding │ P2P Blocking │ Loadbalance │ Broadcast │ Tag Name       │ Controller            │
├────────┼──────────┼────┼───────────────────┼──────────────┼─────────────────┼───────────────┼────────────┼─────────────┼──────────────┼─────────────────────┼─────────────────┼──────────────┼─────────────┼───────────┼────────────────┼───────────────────────┤
│   ✅️   │ labo1     │ 1  │ labo-wlan-profile │ LAB-INTERNAL │ 43,200s         │       ✅️      │ platinum   │ platinum-up │ full         │ PSK                 │ Drop            │ Disabled     │      ⬜️     │     ⬜️    │ labo-wlan-flex │ wnc1.example.internal │
│   ✅️   │ labo3     │ 3  │ labo-wlan-profile │ LAB-INTERNAL │ 43,200s         │       ✅️      │ platinum   │ platinum-up │ full         │ Unknown             │ Drop            │ Disabled     │      ✅️     │     ⬜️    │ labo-wlan-flex │ wnc1.example.internal │
│   ✅️   │ labo2     │ 2  │ labo-wlan-profile │ LAB-INTERNAL │ 43,200s         │       ✅️      │ platinum   │ platinum-up │ full         │ PSK                 │ Drop            │ Disabled     │      ⬜️     │     ⬜️    │ labo-wlan-flex │ wnc1.example.internal │
└────────┴──────────┴────┴───────────────────┴──────────────┴─────────────────┴───────────────┴────────────┴─────────────┴──────────────┴─────────────────────┴─────────────────┴──────────────┴─────────────┴───────────┴────────────────┴───────────────────────┘
```

### JSON Format

```json
[
  {
    "WLANID": "1",
    "ESSID": "CorpWiFi",
    "Status": "Enabled",
    "Security": "WPA2-PSK",
    "VLAN": "100",
    "RadioPolicy": "All",
    "ClientCount": "45"
  },
  {
    "WLANID": "2",
    "ESSID": "GuestNet",
    "Status": "Enabled",
    "Security": "Open",
    "VLAN": "200",
    "RadioPolicy": "All",
    "ClientCount": "12"
  }
]
```

## 📖 Related Commands

- [wnc show client](SHOW_CLIENT.md)
- [wnc show ap](SHOW_AP.md)
- [wnc show overview](SHOW_OVERVIEW.md)
