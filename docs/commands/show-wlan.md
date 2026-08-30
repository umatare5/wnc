# wnc show wlan

One row per WLAN and the policy profile bound to it.

```bash
wnc show wlan
```

```plaintext
ID  Profile      SSID         Status   Security            Bands  Broadcast  P2P Block  Policy Status  Switching  Interface     Session TO  DHCP Required  Policy Profile     Tags            Controller
5   labo-p736b2  labo-p736b2  Enabled  WPA2 PSK            2.4    Enabled    Disabled   Active         Local      LAB-INTERNAL  43200       Yes            labo-wlan-profile  labo-wlan-flex  192.168.0.231
6   labo-p736b5  labo-p736b5  Enabled  WPA2 PSK            5      Enabled    Disabled   Active         Local      LAB-INTERNAL  43200       Yes            labo-wlan-profile  labo-wlan-flex  192.168.0.231
7   labo-t6c73d  labo-t6c73d  Enabled  WPA3 802.1X-SHA256  5/6    Enabled    Disabled   Active         Local      LAB-INTERNAL  43200       Yes            labo-wlan-profile  labo-wlan-flex  192.168.0.231
```

## Columns

| Field                     | Meaning                                            |
| ------------------------- | -------------------------------------------------- |
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

**A row is a WLAN paired with one policy profile.** The same WLAN bound under two tags to two different profiles is two rows, because the policy half differs. Two tags naming the same pair produce one row listing both.

**A WLAN with no binding still gets a row.** Only the policy half is unreported: the WLAN exists and is worth seeing.

**A binding naming a WLAN that does not exist is counted and reported.** The model permits it, and an inner join on its own would hide the misconfiguration and show a clean estate.

**`status` and `policy_status` are both needed.** A WLAN can be enabled while the profile bound to it is shut, in which case no radio carries it.

**`security` is derived in a fixed order.** The master WPA switch is read first, then WEP, OSEN and shared-key authentication — none of which is open even with every key-management flag false — and only then the WPA generation and the key-management set. The fast-transition suffix comes from the FT key-management flags and never from the FT mode setting, whose default would put `+FT` on nearly every WLAN.

**`bands` comes from the per-band list only.** The controller also publishes a single legacy band setting, which is marked obsolete on every release in scope, stopped arriving at 17.18, and reports "all bands" on WLANs the per-band list confines to one.

**`interface` is an interface name, not a VLAN identifier.** That is what the controller calls it.

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
