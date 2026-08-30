# wnc show ap

One row per access point: what it is, how it is powered, what it is plugged into and how long it has been up.

```bash
wnc show ap
```

```plaintext
AP Name    Model             Serial       Ethernet MAC       Radio MAC          IP Address      SW Version  Slots  Country  Mode         Admin    State       LLDP Neighbor                        Power Type    Power Mode  Uptime  Assoc  Controller
TEST-AP01  AIR-AP1815I-Q-K9  XXXXXXXXXXX  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.11  17.12.7.13  2      J4       FlexConnect  Enabled  Registered  lab2-cat29c-06f-01.labo.local:Gi0/2  PoE+          Full Power  1d14h   1d14h  192.168.0.231
TEST-AP02  AIR-AP2802I-Q-K9  XXXXXXXXXXX  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.12  17.12.7.13  2      J4       FlexConnect  Enabled  Registered  lab2-cat29c-06f-01.labo.local:Gi0/4  PoE+          Full Power  1d1h    1d1h   192.168.0.231
TEST-AP03  CW9166I-Q         XXXXXXXXXXX  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.13  17.12.7.13  3      J4       FlexConnect  Enabled  Registered  lab2-cat29c-06f-01.labo.local:Gi0/3  PoE (legacy)  Full Power  5d20h   5d20h  192.168.0.231
```

## Columns

| Field                  | Meaning                                                      |
| ---------------------- | ------------------------------------------------------------ |
| `ap_name`              | Access point name                                            |
| `model`                | Hardware model                                               |
| `serial`               | Chassis serial number                                        |
| `ethernet_mac`         | Wired interface address                                      |
| `radio_mac`            | Base radio address, which is the controller's key for the AP |
| `ip_address`           | Management address                                           |
| `sw_version`           | Running software version                                     |
| `slots`                | Number of radio slots                                        |
| `country`              | Regulatory country code                                      |
| `mode`                 | AP mode, with the sub-mode in brackets where one is set      |
| `admin`                | Administrative state                                         |
| `state`                | Operational state                                            |
| `lldp_neighbor`        | Neighbour system name and port, comma separated if several   |
| `power_type`           | Power source                                                 |
| `power_mode`           | Power the AP is drawing                                      |
| `uptime_seconds`       | How long the access point itself has been up                 |
| `assoc_uptime_seconds` | How long the current CAPWAP association has lasted           |
| `controller`           | The controller this row was read from                        |

## Notes

**`uptime_seconds` and `assoc_uptime_seconds` are different quantities.** The first is the access point's own uptime and the second is the age of its current association. A controller switchover renews only the second, so a view carrying one of them would read a switchover as a fleet reboot. The controller's own `show ap uptime` prints both, as "AP Up Time" and "Association Up Time".

**`mode` includes the sub-mode because it changes how the rest of the row reads.** A Monitor access point with no clients is healthy — the same row on a Local one is an outage. A WIPS sub-mode renders as `Monitor (WIPS)`.

Changing the mode reboots the access point, so a fleet that has just been reconfigured shows a short uptime and a shorter association age. Until the access point rejoins it appears in neither this view nor any other — the controller still remembers it, and `show wireless stats ap join summary` on the device names the reason.

**`slots` counts radios and not ports.** The controller also publishes a slot count that includes a remote-LAN port, which reads one high on a model that has one. This column uses the radio count, which is what `show ap summary` prints.

**`power_type` stays inside what the controller says.** Its schema describes a local supply, an injector, a "legacy" PoE mechanism and an "advanced" one, and names no 802.3 clause anywhere, so neither does this column. `show ap config general` collapses the last two to a bare `PoE` — keeping them apart shows which one was negotiated.

**`admin` is the state as read, which is spelt differently from the state as set.** The RPC that changes it uses a separate value domain with the same numbers, so the two are never mapped onto one another.

**`lldp_neighbor` can hold several entries.** The controller lists one row per neighbour, and they are joined onto the single access point row rather than duplicating it. An access point with no neighbour shows `-`.

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
