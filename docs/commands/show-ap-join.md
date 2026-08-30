# wnc show ap-join

One row per access point the controller remembers — **including one that is not joined**, which no other command can show.

```bash
wnc show ap-join
```

```plaintext
AP Name    Radio MAC          Ethernet MAC       IP Address      Status  Last Failure Phase  Last Join Failure        Last Config Failure  Last Discovery Failure  Last Disconnect Reason      Reboot Reason                 Last Join  Last Config  Last Discovery  Last Error  Controller
TEST-AP01  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.11  Joined  Join                jf-dtls-alert-from-peer  None                 None                    DTLS close alert from peer  ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h7m        192.168.0.231
TEST-AP02  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.12  Joined  Image-Download      None                     None                 None                    Image Download Success      ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h9m        192.168.0.231
TEST-AP03  aa:bb:cc:dd:ee:ff  aa:bb:cc:dd:ee:ff  192.168.255.13  Joined  Join                jf-dtls-alert-from-peer  None                 None                    DTLS close alert from peer  ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h7m        192.168.0.231
```

## Columns

| Field                    | Meaning                                                                              |
| ------------------------ | ------------------------------------------------------------------------------------ |
| `ap_name`                | Access point name                                                                    |
| `radio_mac`              | Base radio address, which the device's join summary calls `Base MAC`                 |
| `ethernet_mac`           | Wired interface address                                                              |
| `ip_address`             | Address the access point discovered the controller from                              |
| `status`                 | `Joined` or `Not Joined`, in the device's own words                                  |
| `last_failure_phase`     | Which phase failed last: Unknown, Discovery, DTLS, Join, Config, Image-Download, Run |
| `last_join_failure`      | Reason the last join attempt failed                                                  |
| `last_config_failure`    | Reason the last configuration attempt failed                                         |
| `last_disc_failure`      | Reason the last discovery attempt failed                                             |
| `disconnect_reason`      | Free text the controller records for the last disconnect                             |
| `reboot_reason`          | Reason the access point last rebooted                                                |
| `last_join_seconds`      | Age of the last successful join                                                      |
| `last_config_seconds`    | Age of the last successful configuration                                             |
| `last_discovery_seconds` | Age of the last successful discovery                                                 |
| `last_error_seconds`     | Age of the last error                                                                |
| `controller`             | The controller this row was read from                                                |

## Notes

**This is the only view that shows an access point that is not joined.** The controller drops an unjoined access point from `capwap-data`, so `wnc show ap`, `wnc show ap-tag` and `wnc show overview` all lose it while the controller still remembers why it left. Measured during a mode change: `show ap summary` on the device reported two access points and `show wireless stats ap join summary` reported three.

**No counter is a column.** The list carries nineteen of them — join requests received, successful responses sent, discovery errors, DTLS failures — and nothing on it declares a clear time. A cumulative total with no window answers no question a reader can act on, so none is rendered. Read the route directly if a total is wanted.

**An age of `-` means the event never happened.** The controller writes `1970-01-01T00:00:00+00:00` for an instant it has no value for, and six of the eleven such leaves read that way on every record of the lab fleet. A `-` under `Last Config` beside a joined access point therefore says the join completed and the configuration never did.

**`last_failure_phase` is the last phase that failed, not the current state.** A healthy access point reports the phase it last stumbled in even while joined, so `Join` beside `Joined` is normal: the phase is history and `status` is now.

**Only the healthy member of each reason domain is translated.** The join, config, discovery and reboot domains carry 42, 14, 17 and 59 members and the device prints a display string for none of them, so `None` is the one word mapped and every other spelling appears as the controller sent it. `Last Failure Phase` is the exception — its seven members are all mapped, because the device does print those.

**The disconnect reason is free text and is never translated.** The controller writes a sentence there. It is more reliable than the sibling enum leaf, which read `unkown` on five of seven lab records while the free text was populated on all seven.

## Examples

Access points the controller remembers but is not serving:

```bash
wnc show ap-join -f json | jq -r '.[] | select(.status != "Joined") | "\(.ap_name) \(.disconnect_reason)"'
```

Access points that joined but were never configured:

```bash
wnc show ap-join -f json | jq -r '.[] | select(.last_join_seconds and (.last_config_seconds | not)) | .ap_name'
```

Most recently joined first:

```bash
wnc show ap-join -b last_join_seconds
```
