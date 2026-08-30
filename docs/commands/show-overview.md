# wnc show overview

One row per access point radio, with the RF settings and the load on it.

```bash
wnc show overview
```

```plaintext
AP Name    AP MAC             Slot  Mode         Band  Admin    Oper  Channel  Width  TxPower  Clients    ChUtil  RF Profile           Controller
TEST-AP01  aa:bb:cc:dd:ee:ff  0     FlexConnect  2.4   Enabled  Up    11ch     20MHz  20dBm    1clients   23%     labo-rf-24gh         192.168.0.231
TEST-AP01  aa:bb:cc:dd:ee:ff  1     FlexConnect  5     Enabled  Up    64ch     40MHz  18dBm    0clients   1%      labo-rf-5gh-inside   192.168.0.231
TEST-AP02  aa:bb:cc:dd:ee:ff  0     FlexConnect  2.4   Enabled  Up    1ch      20MHz  19dBm    2clients   16%     labo-rf-24gh         192.168.0.231
TEST-AP02  aa:bb:cc:dd:ee:ff  1     FlexConnect  5     Enabled  Up    48ch     40MHz  17dBm    0clients   1%      labo-rf-5gh-outside  192.168.0.231
TEST-AP03  aa:bb:cc:dd:ee:ff  0     FlexConnect  2.4   Enabled  Up    6ch      20MHz  22dBm    14clients  10%     labo-rf-24gh         192.168.0.231
TEST-AP03  aa:bb:cc:dd:ee:ff  1     FlexConnect  5     Enabled  Up    116ch    40MHz  22dBm    2clients   2%      labo-rf-5gh-inside   192.168.0.231
TEST-AP03  aa:bb:cc:dd:ee:ff  2     FlexConnect  6     Enabled  Up    5ch      40MHz  18dBm    1clients   2%      labo-rf-6gh          192.168.0.231
```

## Options

Beyond the shared `show` options:

| Option            | Meaning                                 |
| ----------------- | --------------------------------------- |
| `--radio`, `-r`   | Keep only `2.4`, `5` or `6` GHz radios  |
| `--sort-by`, `-b` | Any field name below. Default `ap_name` |

## Columns

| Field               | Meaning                                                            |
| ------------------- | ------------------------------------------------------------------ |
| `ap_name`           | Access point name                                                  |
| `ap_mac`            | The access point's base radio address, shared by all of its radios |
| `slot`              | Radio slot. This plus `ap_mac` identifies the row                  |
| `mode`              | Radio mode: Local, Monitor, FlexConnect, Sniffer and so on         |
| `band`              | `2.4`, `5` or `6`                                                  |
| `admin`             | Administrative state                                               |
| `oper`              | Operational state                                                  |
| `channel`           | Channel **number**, not a frequency. Shown as `11ch`               |
| `channel_width_mhz` | Channel width in MHz, shown as `20MHz`                             |
| `tx_power_dbm`      | Transmit power on the band the radio is on now, shown as `19dBm`   |
| `clients`           | Clients the controller holds in the run state, shown as `1clients` |
| `ch_util_percent`   | Channel utilization as last measured by RRM, shown as `28%`        |
| `rf_profile`        | The RF profile the tag in force supplies for this band             |
| `controller`        | The controller this row was read from                              |

## Notes

**The row identity is `ap_mac` plus `slot`.** `ap_mac` is the access point's base radio address, so every radio of one access point reports the same value. The column is not called a radio MAC for that reason, and neither does the controller's own per-radio summary.

**`channel` is a channel number and not a frequency.** `show ap dot11 5ghz summary` renders the same value as `(64,60)`, and the schema declares no unit for it.

**`tx_power_dbm` is the power on the band the radio is on now.** An XOR radio reports one entry per band it can use, and the one that matches the band currently selected is the one shown. Reading the first entry instead would report 22 dBm where the controller's own summary says 18.

**`clients` comes from the client list, not from the RRM measurement.** The measurement list is shorter than the radio list, so a radio without a measurement row would report zero clients rather than an unknown number. The count therefore agrees with the `wnc show client` rows whose State is `Run`, read at the same moment.

**`ch_util_percent` is a measurement, not a live reading.** RRM updates it on its own cycle. The schema declares no timestamp for it on any release in scope, so its age cannot be shown.

**100 is not a declared ceiling.** `cca-util-percentage` is a `uint16` and the RRM module carries no `range` statement at all, on any release in scope. The value is not clamped here and nothing should compute a spare capacity as 100 minus it.

**Two of the four units are the controller's and two are this view's.** It declares `units "dBm"` on the transmit-power leaf and `units "percentage"` on the utilization leaf. The width is a bare `uint8` and no module declares an Hz unit anywhere, so `MHz` rests on the controller's own enum descriptions — and the channel carries no unit at all, which is why its cell says `ch`.

**The unit sits in the table and never in the JSON.** It is glued to the number so each cell stays one field for `awk` and `cut`, and an unreported value shows a bare `-`. Sorting reads the number, so `--sort-by channel` puts channel 6 before channel 11 rather than ordering the text.

**A remote-LAN port is not listed.** The controller reports one as a radio entry with no mode, band, state, channel or power — it is dropped rather than rendered as a row of dashes.

**A Monitor or Sniffer radio shows no channel, and still shows a width and a power.** The schema guards `curr-freq` on the radio mode and so sends nothing, while the width and power leaves carry no such guard and are sent. The device's own `show ap dot11 <band> summary` renders all three as `N/A` — dropping a value the controller sent is the same fault as reporting one it withheld. The Mode column is what says the radio is not serving.

**`-` is not zero.** A radio whose state the controller omitted shows `-` rather than Down, and a radio with no clients on an idle channel shows `0clients` and `0%` because those were reported.

**`rf_profile` is the outcome, not the tag.** It is the profile the RF tag in force supplies for the band this radio is on, picked by band rather than by slot. `wnc show rf-tag` lists the tags themselves, including one that names no profile on a band.

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
