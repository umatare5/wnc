# wnc show ap-join

One row per access point the controller remembers, joined or not.

```bash
wnc show ap-join
```

```plaintext
AP Name    Radio MAC          Ethernet MAC       IP Address    Status  Last Failure Phase  Last Join Failure        Last Config Failure  Last Discovery Failure  Last Disconnect Reason      Reboot Reason                 Last Join  Last Config  Last Discovery  Last Error  Controller
TEST-AP01  00:00:5e:00:53:01  00:00:5e:00:53:11  192.168.0.11  Joined  Join                jf-dtls-alert-from-peer  None                 None                    DTLS close alert from peer  ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h7m        WNC1
TEST-AP02  00:00:5e:00:53:02  00:00:5e:00:53:12  192.168.0.12  Joined  Image-Download      None                     None                 None                    Image Download Success      ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h9m        WNC1
TEST-AP03  00:00:5e:00:53:03  00:00:5e:00:53:13  192.168.0.13  Joined  Join                jf-dtls-alert-from-peer  None                 None                    DTLS close alert from peer  ap-reboot-reason-img-upgrade  3h3m       3h3m         3h3m            3h7m        WNC1
```

**It is the only view that shows an access point which is not joined.**

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field                    | Meaning                                                                              |
| :----------------------- | :----------------------------------------------------------------------------------- |
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

- **The only view of an unjoined access point** — every other collection drops it as it leaves
- **No counter is a column** — the list declares no clear time, so a total answers nothing
- **A `-` age is an event that never happened** — `-` under Last Config beside `Joined` says so
- **The failure phase is history** — `Join` beside `Joined` is normal, `status` being what is now
- **Only `None` is translated** — the device prints no string for these domains, bar the phase
- **The disconnect reason is free text** — more reliable than the sibling enum, so never mapped

The readings behind these sit in [`measurements.md`](../measurements.md).

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
