# wnc show overview

One row per access point radio, with the RF settings and the load on it.

```bash
wnc show overview
```

```plaintext
AP Name    AP MAC             Slot  Mode         Band  Admin    Oper  Channel  Width  TxPower  Clients    ChUtil  RF Profile         Controller
TEST-AP01  00:00:5e:00:53:01  0     FlexConnect  2.4   Enabled  Up    11ch     20MHz  20dBm    1clients   23%     test-rf-profile01  WNC1
TEST-AP01  00:00:5e:00:53:01  1     FlexConnect  5     Enabled  Up    64ch     40MHz  18dBm    0clients   1%      test-rf-profile03  WNC1
TEST-AP02  00:00:5e:00:53:02  0     FlexConnect  2.4   Enabled  Up    1ch      20MHz  19dBm    2clients   16%     test-rf-profile01  WNC1
TEST-AP02  00:00:5e:00:53:02  1     FlexConnect  5     Enabled  Up    48ch     40MHz  17dBm    0clients   1%      test-rf-profile04  WNC1
TEST-AP03  00:00:5e:00:53:03  0     FlexConnect  2.4   Enabled  Up    6ch      20MHz  22dBm    14clients  10%     test-rf-profile01  WNC1
TEST-AP03  00:00:5e:00:53:03  1     FlexConnect  5     Enabled  Up    116ch    40MHz  22dBm    2clients   2%      test-rf-profile03  WNC1
TEST-AP03  00:00:5e:00:53:03  2     FlexConnect  6     Enabled  Up    5ch      40MHz  18dBm    1clients   2%      test-rf-profile05  WNC1
```

## Flags

| Option          | Meaning                                |
| :-------------- | :------------------------------------- |
| `--radio`, `-r` | Keep only `2.4`, `5` or `6` GHz radios |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field               | Meaning                                                            |
| :------------------ | :----------------------------------------------------------------- |
| `ap_name`           | Access point name                                                  |
| `ap_mac`            | The access point's base radio address, shared by all of its radios |
| `slot`              | This plus `ap_mac` identifies the row                              |
| `mode`              | Radio mode: Local, Monitor, FlexConnect, Sniffer and so on         |
| `band`              | `2.4`, `5` or `6`                                                  |
| `admin`             | Administrative state                                               |
| `oper`              | Operational state                                                  |
| `channel`           | Channel **number**, not a frequency. Shown as `11ch`               |
| `channel_width_mhz` | Shown as `20MHz`                                                   |
| `tx_power_dbm`      | Transmit power on the band the radio is on now, shown as `19dBm`   |
| `clients`           | Clients the controller holds in the run state, shown as `1clients` |
| `ch_util_percent`   | Channel utilization as last measured by RRM, shown as `28%`        |
| `rf_profile`        | The RF profile the tag in force supplies for this band             |
| `controller`        | The controller this row was read from                              |

## Notes

- **An XOR radio sends one power entry per band** — the one for the band now selected is shown
- **`clients` comes from the client list** — the RRM list is shorter, so a gap would read as zero
- **`ch_util_percent` is RRM's cycle, not a live reading** — no release declares its age
- **100 is not a declared ceiling** — nothing clamps it, and spare capacity is not 100 minus it
- **`MHz` rests on the enum descriptions** — the width leaf declares no unit, the channel none
- **A remote-LAN port is not a row** — it arrives with no mode, band, state, channel or power
- **A Monitor radio still shows width and power** — the guard is on the frequency leaf alone
- **`rf_profile` is the outcome** — picked by band, where [`show rf-tag`](show-rf-tag.md) lists the tags

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Radios above 50% utilization, worst first:

```bash
wnc show overview -f json | jq -r '.[] | select(.ch_util_percent > 50) | "\(.ap_name)/\(.slot) \(.ch_util_percent)%"'
```

The 6 GHz radios only:

```bash
wnc show overview -r 6
```

Radios carrying the most clients:

```bash
wnc show overview -b clients --sort-order desc
```
