# wnc show ap

One row per access point: what it is, how it is powered, what it is plugged into and how long it has been up.

```bash
wnc show ap
```

```plaintext
AP Name    Model             Serial       Ethernet MAC       Radio MAC          IP Address    SW Version  Slots  Country  Mode         Admin    State       LLDP Neighbor                     Power Type    Power Mode  Uptime  Assoc  Controller
TEST-AP01  AIR-AP1815I-Q-K9  TST0000AP01  00:00:5e:00:53:11  00:00:5e:00:53:01  192.168.0.11  17.12.7.13  2      J4       FlexConnect  Enabled  Registered  test-sw01.example.internal:Gi0/2  PoE+          Full Power  1d14h   1d14h  WNC1
TEST-AP02  AIR-AP2802I-Q-K9  TST0000AP02  00:00:5e:00:53:12  00:00:5e:00:53:02  192.168.0.12  17.12.7.13  2      J4       FlexConnect  Enabled  Registered  test-sw01.example.internal:Gi0/4  PoE+          Full Power  1d1h    1d1h   WNC1
TEST-AP03  CW9166I-Q         TST0000AP03  00:00:5e:00:53:13  00:00:5e:00:53:03  192.168.0.13  17.12.7.13  3      J4       FlexConnect  Enabled  Registered  test-sw01.example.internal:Gi0/3  PoE (legacy)  Full Power  5d20h   5d20h  WNC1
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field                  | Meaning                                                                |
| :--------------------- | :--------------------------------------------------------------------- |
| `ap_name`              | Access point name                                                      |
| `model`                | Hardware model                                                         |
| `serial`               | Chassis serial number                                                  |
| `ethernet_mac`         | Wired interface address                                                |
| `radio_mac`            | Base radio address, which is the controller's key for the access point |
| `ip_address`           | Management address                                                     |
| `sw_version`           | Running software version                                               |
| `slots`                | Number of radio slots                                                  |
| `country`              | Regulatory country code                                                |
| `mode`                 | AP mode, with the sub-mode in brackets where one is set                |
| `admin`                | Administrative state                                                   |
| `state`                | Operational state                                                      |
| `lldp_neighbor`        | Neighbour system name and port, comma separated if several             |
| `power_type`           | Power source                                                           |
| `power_mode`           | Power the access point is drawing                                      |
| `uptime_seconds`       | How long the access point itself has been up                           |
| `assoc_uptime_seconds` | How long the current CAPWAP association has lasted                     |
| `controller`           | The controller this row was read from                                  |

## Notes

- **Two uptimes, two quantities** — a controller switchover renews the association age alone
- **Mode carries the sub-mode** — a Monitor row with no clients is healthy, a Local one is not
- **A mode change reboots** — until it rejoins, [`show ap-join`](show-ap-join.md) is the only view left
- **`slots` counts radios** — the controller's other slot count includes a remote-LAN port
- **`power_type` stays inside the schema** — it names no 802.3 clause, so neither does the column
- **`admin` is the state as read** — the RPC that sets it uses a separate domain, never mapped
- **Several neighbours share one row** — the controller's rows are joined, not duplicated

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Access points that are not registered:

```bash
wnc show ap -f json | jq -r '.[] | select(.state != "Registered") | "\(.ap_name) \(.state)"'
```

Access points whose association is much younger than their uptime, which is what a controller switchover looks like:

```bash
wnc show ap -f json | jq -r '.[] | select(.uptime_seconds - .assoc_uptime_seconds > 600) | .ap_name'
```

The most recently rebooted access points first:

```bash
wnc show ap -b uptime_seconds
```
