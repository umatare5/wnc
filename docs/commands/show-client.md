# wnc show client

One row per associated client, joined across the five collections that describe one.

```bash
wnc show client
```

```plaintext
MAC                IPv4        IPv6          Device          Username  SSID         AP Name    Slot  Band  Protocol  Channel  State           RSSI    SNR   Rate     Streams  Assoc  Rx        Tx       Controller
00:00:5e:00:53:11  192.0.2.11  2001:db8::11  Example Vendor  -         labo-p736b2  TEST-AP03  0     2.4   11ax      6ch      Run             -21dBm  78dB  143Mbps  1ss      2h15m  17.0MiB   19.6MiB  192.168.0.231
00:00:5e:00:53:12  192.0.2.12  -             Example Phone   lab-user  labo-p736b5  TEST-AP01  1     5     11ac      64ch     Run             -43dBm  56dB  866Mbps  2ss      17m    107.9KiB  30.7KiB  192.168.0.231
00:00:5e:00:53:13  192.0.2.13  -             Example Sensor  -         labo-t6c73d  TEST-AP03  2     6     11be      5ch      Authenticating  -55dBm  40dB  -        -        42s    -         -        192.168.0.231
```

## Options

Beyond the shared `show` options:

| Option          | Meaning                                |
| --------------- | -------------------------------------- |
| `--radio`, `-r` | Keep only clients on `2.4`, `5` or `6` |
| `--ssid`, `-s`  | Keep only clients on this SSID         |
| `--ap-name`     | Keep only clients on this access point |

## Columns

| Field             | Meaning                                                       |
| ----------------- | ------------------------------------------------------------- |
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

**`device` is a classification label, not a hostname.** The controller assigns it from traffic analysis, so the same value covers many clients. The leaf that does carry a hostname sits on a list this CLI cannot reach and which answers with no content, so it is not offered.

**`band` and `protocol` are two different facts.** The band says which radio the client is on and the protocol says which PHY generation it negotiated. The controller's own `show wireless client summary` pairs them as `11n(2.4)` — here they are separate columns so either can be filtered and sorted.

**`username` is `-` for an unauthenticated client.** The controller sends the key with an empty value rather than omitting it, so an empty username is reported as unreported.

**`ipv6` is one address chosen from up to eight.** Link-local addresses are dropped because every client has one and it identifies nothing, and the remaining addresses are compared numerically. The controller compresses some entries and not others, so a textual comparison would return a different address between polls.

**A zero is a reading only where zero is possible.** `snr_db` of 0 is a real margin and is shown as `0dB` — a channel, a rate, a stream count and an RSSI of exactly 0 are omissions and show `-`.

**`snr_db` is `-` when the traffic counters carried no row for the client.** The margin is read from a collection of its own, so a client the controller lists without one has no SNR to report rather than a margin of zero.

**The unit sits in the table and never in the JSON.** It is glued to the number so each cell stays one field for `awk` and `cut`, and an unreported value shows a bare `-` rather than a unit with nothing in front of it. Sorting reads the number, so `--sort-by channel` puts channel 6 before channel 11 rather than ordering the text.

**`assoc_seconds` is `-` when the controller reported no instant.** Several sibling timestamps on the same read return the Unix epoch to mean that, and an age computed from it would read as fifty-six years.

**A filter reading from a collection that failed produces no rows for that controller.** Reporting zero matches instead would claim the fleet holds no such client, when the truth is that nothing was read.

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
wnc show client -r 6 -s labo-t6c73d
```

The heaviest talkers:

```bash
wnc show client -b tx_bytes --sort-order desc
```
