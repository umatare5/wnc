# wnc show wlan

One row per WLAN and the policy profile bound to it.

```bash
wnc show wlan
```

```plaintext
ID  Profile              SSID          Status   Security            Bands  Broadcast  P2P Block  Policy Status  Switching  Interface      Session TO  DHCP Required  Policy Profile         Tags            Controller
5   test-wlan-profile01  test-essid01  Enabled  WPA2 PSK            2.4    Enabled    Disabled   Active         Local      TEST-INTERNAL  43200       Yes            test-policy-profile01  test-wlan-flex  WNC1
6   test-wlan-profile02  test-essid02  Enabled  WPA2 PSK            5      Enabled    Disabled   Active         Local      TEST-INTERNAL  43200       Yes            test-policy-profile01  test-wlan-flex  WNC1
7   test-wlan-profile03  test-essid03  Enabled  WPA3 802.1X-SHA256  5/6    Enabled    Disabled   Active         Local      TEST-INTERNAL  43200       Yes            test-policy-profile01  test-wlan-flex  WNC1
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field                     | Meaning                                            |
| :------------------------ | :------------------------------------------------- |
| `wlan_id`                 | WLAN identifier                                    |
| `profile`                 | WLAN profile name                                  |
| `ssid`                    | Broadcast network name                             |
| `status`                  | Whether the WLAN itself is enabled                 |
| `security`                | WPA generation, key management and fast transition |
| `bands`                   | Bands the WLAN is enabled on                       |
| `broadcast`               | Whether the SSID is advertised                     |
| `p2p_block`               | Peer-to-peer blocking action                       |
| `policy_status`           | Whether the bound policy profile is active or shut |
| `switching`               | Central or local switching                         |
| `interface`               | Interface name the policy profile attaches         |
| `session_timeout_seconds` | Session timeout                                    |
| `dhcp_required`           | Whether DHCP is required of every client           |
| `policy_profile`          | Policy profile bound to this WLAN                  |
| `tags`                    | Policy tags that produced this pairing             |
| `controller`              | The controller this row was read from              |

## Notes

- **A row is a WLAN and one policy profile** — the same WLAN under two profiles is two rows
- **A WLAN bound to nothing still gets a row** — only the policy half of it is unreported
- **A binding naming no WLAN is counted** — an inner join would hide the misconfiguration
- **Both status columns are needed** — an enabled WLAN under a shut profile reaches no radio
- **The `+FT` suffix follows the key-management flags** — the FT mode default would mark all

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

WLANs that are enabled but whose policy profile is shut, so nothing is on air:

```bash
wnc show wlan -f json | jq -r '.[] | select(.status == "Enabled" and .policy_status == "Shutdown") | .profile'
```

WLANs with no encryption:

```bash
wnc show wlan -f json | jq -r '.[] | select(.security == "Open") | "\(.wlan_id) \(.ssid)"'
```

WLANs on 6 GHz:

```bash
wnc show wlan -f json | jq -r '.[] | select(.bands | test("6")) | .ssid'
```
