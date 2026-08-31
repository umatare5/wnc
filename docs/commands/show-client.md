# wnc show client

One row per associated client, joined across the five collections that describe one.

```bash
wnc show client
```

```plaintext
MAC                IPv4          IPv6          Device          Username   SSID          AP Name    Slot  Band  Protocol  Channel  State           RSSI    SNR   Rate     Streams  Assoc  Rx        Tx       Controller
00:00:5e:00:53:a1  192.168.0.21  2001:db8::11  Example Vendor  -          test-essid01  TEST-AP03  0     2.4   11ax      6ch      Run             -21dBm  78dB  143Mbps  1ss      2h15m  17.0MiB   19.6MiB  WNC1
00:00:5e:00:53:a2  192.168.0.22  -             Example Phone   test-user  test-essid02  TEST-AP01  1     5     11ac      64ch     Run             -43dBm  56dB  866Mbps  2ss      17m    107.9KiB  30.7KiB  WNC1
00:00:5e:00:53:a3  192.168.0.23  -             Example Sensor  -          test-essid03  TEST-AP03  2     6     11be      5ch      Authenticating  -55dBm  40dB  -        -        42s    -         -        WNC1
```

## Flags

| Option          | Meaning                                |
| :-------------- | :------------------------------------- |
| `--radio`, `-r` | Keep only clients on `2.4`, `5` or `6` |
| `--ssid`, `-s`  | Keep only clients on this SSID         |
| `--ap-name`     | Keep only clients on this access point |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field             | Meaning                                                       |
| :---------------- | :------------------------------------------------------------ |
| `mac`             | Client address                                                |
| `ipv4`            | IPv4 address from the binding table                           |
| `ipv6`            | Lowest global IPv6 address from the binding table             |
| `device`          | The controller's device-classification label                  |
| `username`        | Authenticated username                                        |
| `ssid`            | SSID the client is on                                         |
| `ap_name`         | Access point serving the client                               |
| `slot`            | Radio slot serving the client                                 |
| `band`            | `2.4`, `5` or `6`                                             |
| `protocol`        | PHY generation: 11n, 11ac, 11ax, 11be                         |
| `channel`         | Channel number, shown as `6ch`                                |
| `state`           | Association phase the client last completed                   |
| `rssi_dbm`        | Last received signal strength, shown as `-21dBm`              |
| `snr_db`          | Last signal-to-noise margin, shown as `78dB`                  |
| `speed_mbps`      | Last negotiated rate, shown as `143Mbps`                      |
| `spatial_streams` | Spatial streams in use, shown as `1ss`                        |
| `assoc_seconds`   | Age of the association                                        |
| `rx_bytes`        | Octets received, raw in JSON and in binary units in the table |
| `tx_bytes`        | Octets transmitted, likewise                                  |
| `controller`      | The controller this row was read from                         |

## Notes

- **`device` is a classification label** — the leaf carrying a hostname answers with no content
- **Band and protocol are separate** — the device pairs them, so neither can be sorted there
- **One IPv6 of up to eight** — link-local dropped, the rest compared as numbers not as text
- **A zero is a reading only where zero is possible** — `0dB` is a margin, channel 0 an omission
- **`snr_db` reads a collection of its own** — a client with no row has no margin, not zero
- **A filter over a failed read prints nothing** — zero matches would claim an empty fleet

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Clients with a weak signal:

```bash
wnc show client -f json | jq -r '.[] | select(.rssi_dbm < -70) | "\(.mac) \(.ap_name) \(.rssi_dbm)"'
```

Clients stuck short of the run state:

```bash
wnc show client -f json | jq -r '.[] | select(.state != "Run") | "\(.mac) \(.state) \(.assoc_seconds)s"'
```

The 6 GHz clients on one SSID:

```bash
wnc show client -r 6 -s test-essid03
```

The heaviest talkers:

```bash
wnc show client -b tx_bytes --sort-order desc
```
