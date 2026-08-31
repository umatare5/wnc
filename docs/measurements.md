# Measurements

Every reading taken on a live Catalyst 9800 sits here once, and a page under [`commands/`](commands/) links here.

`All three` is `17.12.8`, `17.15.6` and `17.18.4a` — `Unrecorded` is a reading whose release went unwritten.

## Device heading map

Columns the device heads differently, with the command that prints it.

| Column                   | Device heading                     | Device command                        |
| :----------------------- | :--------------------------------- | :------------------------------------ |
| `ap_join_profile`        | `AP Profile`                       | `show wireless tag site detailed`     |
| `assoc_uptime_seconds`   | `Association Up Time`              | `show ap uptime`                      |
| `band` and `protocol`    | `11n(2.4)`, the pair in one field  | `show wireless client summary`        |
| `channel`                | `(64,60)`, a pair per radio        | `show ap dot11 5ghz summary`          |
| `misconfig_reason`       | none at all on the device          | `show ap tag summary`                 |
| `power_type`             | `PoE`, both mechanisms collapsed   | `show ap config general`              |
| `profile_24ghz`          | `2.4ghz RF Policy`                 | `show wireless tag rf detailed`       |
| `profile_5ghz`           | `5ghz RF Policy`                   | `show wireless tag rf detailed`       |
| `profile_6ghz`           | `6ghz RF Policy`                   | `show wireless tag rf detailed`       |
| `radio_mac`              | `Base MAC`                         | `show wireless stats ap join summary` |
| `uptime_seconds`         | `AP Up Time`                       | `show ap uptime`                      |
| `wlan`, `policy_profile` | `WLAN Profile Name`, `Policy Name` | `show wireless tag policy detailed`   |

## Access point

The access-point collections and the four RPCs that name one.

| Fact                                                                                                     | Release           | Condition                                                                    |
| :------------------------------------------------------------------------------------------------------- | :---------------- | :--------------------------------------------------------------------------- |
| `ap-name-mac-map=<name>` returns one row: the name, the base radio and the Ethernet address              | All three         | `404` for a name no access point holds, and a `200` with no row reads alike  |
| The keyed read `404`s a 256-character name, a space, a slash and a multi-byte character alike            | 17.18.4a          | Four spellings probed on the one release, each of them answering alike       |
| `reset ap`: out of `capwap-data` within 16s, rejoin at most 285s, clients down throughout                | 17.12.8, 17.18.4a | AIR-AP1815I, the model the bound was taken on and the one it moves with      |
| `Not Joined` with `Wtp reset config cmd sent`, then `reboot_reason` `ap-reboot-reason-reboot-cmd`        | 17.12.8           | Through and after a `reset ap`, readable in `wnc show ap-join` alone         |
| Name-arm `reset ap`: out of `capwap-data`, a new boot time, and only the access point named moved at all | 17.18.4a          | Through the `ap-name` arm rather than the base-radio address arm             |
| `reset capwap`: rejoin within 10s, `uptime_seconds` unchanged, `assoc_uptime_seconds` restarted          | 17.12.8, 17.18.4a | Against a `reset ap` at 285s, where both quantities restart together         |
| `num-join-req-recvd`, `num-config-req-recvd` and `ctrl-dtls-setup-req` accumulate across CAPWAP sessions | 17.12.8, 17.18.4a | On `ap-join-stats`, so a session reset shows as a counter delta              |
| An access-point-level disable takes `ap-state/ap-admin-state` to `adminstate-disabled`                   | 17.15.6, 17.18.4a | Both radios stayed `admin-state` `enabled` at `oper-state` `radio-down`      |
| A disabled access point stays `Registered` in `capwap-data` with `uptime_seconds` climbing               | 17.15.6, 17.18.4a | No reboot at all, which is the difference from a `reset ap`                  |
| `set-ap-admin-state`: the name arm produced the readings above, exactly as the address arm did           | 17.18.4a, 17.15.6 | Name arm at 17.18.4a, address arm at 17.15.6                                 |
| An unjoined access point leaves `capwap-data` and the controller still counts it                         | Unrecorded        | During a mode change: `show ap summary` read 2 where the join summary read 3 |
| Six of the eleven instant leaves on the join record read `1970-01-01T00:00:00+00:00`                     | All three         | Every record of the lab fleet, the sentinel for an event with no value       |
| Nineteen counters on the join record, and a clear time declared on none of them                          | All three         | On `ap-join-stats`, the list `wnc show ap-join` reads its columns from       |
| Join, config, discovery and reboot reason domains: 42, 14, 17 and 59 members                             | All three         | No display string printed for any, where all seven failure phases print      |
| The enum disconnect reason read `unkown` on 5 of 7 records, the free text on all 7                       | Unrecorded        | The controller's own spelling, which is why the free text is the column      |
| A second slot count includes the remote-LAN port, so it reads one radio high                             | Unrecorded        | Only on a model carrying one, where `show ap summary` prints the radios      |

## Radio

`radio-oper-data`, the band number the slot RPC takes, and the RRM measurement beside it.

| Fact                                                                                        | Release           | Condition                                                                      |
| :------------------------------------------------------------------------------------------ | :---------------- | :----------------------------------------------------------------------------- |
| The band number follows `radio-type`, not the band served — the served band's number `400`s | All three         | A dual-band or XOR radio takes 3 on either of the bands it serves              |
| XOR slot 2 serving 6 GHz: band 4 `400`s "AP slot: 2 does not have a dedicated radio"        | 17.15.6           | `TEST-AP03`, a 5-or-6 GHz XOR radio reporting `dot11-6-ghz-band`               |
| Band 3 on that same radio `204`s and takes slot 2 down while slots 0 and 1 stay up          | 17.15.6           | The second direction of the pair, taken on the one radio                       |
| Band 3 in slot 0 `400`s "AP does not support the specified radio type"                      | Unrecorded        | A dedicated 2.4 GHz layout, which takes band 1 — one of two access points      |
| Band 3 in slot 0 is accepted where band 1, the band served, `400`s first                    | 17.12.8           | A 2.4-or-5 GHz XOR layout, the other of the two access points                  |
| `must` accepts band 1 on slot 0 and no other slot                                           | All three         | A dedicated 2.4 GHz radio, `radio-80211bg`                                     |
| `must` accepts band 2 on slot 1 or slot 2                                                   | All three         | A dedicated 5 GHz radio, `radio-80211a`                                        |
| `must` accepts band 3 on slot 0 or slot 2                                                   | All three         | Dual band and XOR alike: `radio-80211abgn` and `radio-80211-xor-5-6ghz`        |
| `must` accepts band 4 on slot 2 or slot 3                                                   | All three         | A dedicated 6 GHz radio, `radio-80211-6ghz`, and never once written            |
| Band 1 with slot 1 `400`s, the controller naming that `must` clause                         | 17.15.6           | The one pair the clause forbids that was probed on the wire                    |
| `--slot 3` is refused only after the controller answers, so it exits 1                      | 17.12.8, 17.15.6  | An access point holding slots 0 and 1 only, so exit 2 never applies            |
| `enm-radio-type`: eight members at 17.12.8 and 17.15.6, nine at 17.18.4a                    | All three         | The ninth is `radio-80211-xor-24-6ghz`, which fits both 3 and 4                |
| `radio-80211-xor-24-6ghz`, `radio-uwb` and `radio-invalid` take no band number              | All three         | No leaf says which of 3 or 4 a 2.4-or-6 GHz XOR radio wants                    |
| Radio disable then enable: slot 1 `Enabled/Up` to `Disabled/Down` and back                  | 17.15.6, 17.18.4a | Slot 0 untouched throughout, and `--slot 0` behaved as the mirror image        |
| An XOR radio sends one power entry per band it can use, the first reading 22 dBm            | Unrecorded        | The device's own summary says 18 on the band currently selected                |
| `curr-freq` is guarded on the radio mode and absent on Monitor and Sniffer radios           | All three         | Width and power carry no guard and arrive, the device printing all three `N/A` |
| A remote-LAN port arrives as a radio entry with no mode, band, state, channel or power      | All three         | Listed among the radios, so the slot exists to be named                        |
| The RRM measurement list is shorter than the radio list it describes                        | Unrecorded        | So a radio can hold no channel-utilization row of its own                      |
| `cca-util-percentage` is a `uint16` and the RRM module declares no `range` statement at all | All three         | 100 is no ceiling, and a spare capacity must not be 100 minus it               |
| `units "dBm"` on transmit power and `units "percentage"` on channel utilization             | All three         | Width a bare `uint8`, channel no unit, and no Hz unit in any module            |

## Client

The client collections and `apf-ms-delete-all`, whose blast radius no schema states.

| Fact                                                                                                     | Release           | Condition                                                                  |
| :------------------------------------------------------------------------------------------------------- | :---------------- | :------------------------------------------------------------------------- |
| Every `client-mac` a controller serves is already lowercase                                              | 17.12.8, 17.15.6  | Already the form the SDK normalizes to, so nothing in this CLI case folds  |
| `Cisco-IOS-XE-wireless-client-rpc` at revision `2023-03-01` declares two `clear-sisf-binding` operations | 17.12.8           | Read 2026-08-29, and that revision declares no client delete at all        |
| The same module at `2024-03-01` declares `apf-ms-delete-all` beside those two                            | 17.15.6, 17.18.4a | Read 2026-08-29, which is why `deauth` needs 17.15.6 or later              |
| The RPC's own description reads "Delete wireless client based on client MAC or IP address or username"   | Unrecorded        | The name `apf-ms-delete-all` reads as a purge its own model contradicts    |
| A post to the release without it answers `400`, `error-tag: malformed-message`, `"invalid path"`         | 17.12.8           | The CLI re-words that status rather than reporting it bare                 |
| `204` for an identifier associated to nothing exactly as for a session it dropped, both under 330ms      | 17.18.4a          | So the answer alone cannot confirm that a client was ever there            |
| `--mac` post: `assoc_seconds` 6102 to 13, `state` `Run` to `IP Learning`                                 | 17.18.4a          | Recovery to `ipv4` at most 210s, with the other three clients untouched    |
| Two `--mac` posts each dropped two of eighteen: the target and one other station each                    | 17.15.6           | A 25s no-post window moved none of eighteen, so the second was no accident |
| The station that went with it shared the access point, BSSID, radio and WLAN                             | 17.15.6           | Another vendor, and thirteen others on that BSS stayed put                 |
| `--username` post: the client's row gone within 1.9s, deleted rather than reset in place                 | 17.15.6           | Target stable 82 minutes at `assoc_seconds` 4943 across four snapshots     |
| The station re-associated as a new record within 1.5s, back at `Run` with its username inside 42s        | 17.15.6           | So a read seconds later sees a young association rather than an absence    |
| That post moved none of the other 17 clients on the controller                                           | 17.15.6           | Two 35s no-post control windows moved 1 and 0 of eighteen                  |
| Sixteen of eighteen clients carried an empty username                                                    | 17.15.6           | So an empty `--username` would select nearly the whole fleet of eighteen   |
| The RPC's `ip-addr` arm answers `204` and the controller does not refuse it                              | 17.15.6           | Every SISF binding measured carries zone 0, the payload naming no other    |
| The username key arrives carrying an empty value rather than omitted                                     | Unrecorded        | So an empty username is reported as unreported rather than as blank        |
| Up to eight IPv6 addresses per client, some entries compressed and some not                              | Unrecorded        | A textual compare returns a different address between two polls            |
| Sibling instant leaves on the client read return the Unix epoch                                          | Unrecorded        | An age computed off one of them would read as fifty-six years              |
| The list carrying a client hostname answers with no content                                              | Unrecorded        | So the device-classification label is all that can be offered              |

## Tag

The three tag configuration lists and the operational container an AP resolves to.

| Fact                                                                                                 | Release           | Condition                                                                     |
| :--------------------------------------------------------------------------------------------------- | :---------------- | :---------------------------------------------------------------------------- |
| 32 characters accepted, 33 answering `400 Validation failed ... Tag name should not exceed 32`       | 17.12.8           | Per kind, the key leaves declaring a `pattern` and no `length` at all         |
| No `leafref` and no `require-instance` in the three tag configuration modules                        | All three         | Tag names typed as plain strings, so a dangling reference persists on the lab |
| A merge `PATCH` leaves a leaf the payload omits at the value the controller already holds            | 17.12.8, 17.15.6  | A second write naming one profile kept the description and the other profile  |
| A create naming a flex profile without `is-local-site` false: `400 the 'when' expression ... failed` | 17.12.8           | `when "../is-local-site = 'false'"` declared on all three, defaulting to true |
| `rf-tag-radio-profiles` with a null list answers `400 invalid value for: rf-tag-radio-profile`       | 17.12.8, 17.15.6  | So the per-slot radio profile list is out of reach of this CLI                |
| A plain read of `rf-tags` omits all three per-band profiles on `default-rf-tag`                      | All three         | `report-all` returns them, so no view falls back to a plain read              |
| An omitted description and an omitted RF profile name arrive absent rather than empty                | Unrecorded        | The mirror of the filter-name leaf, which arrives as an empty string instead  |
| `is-ap-misconfigured` arrives as an explicit `false` on a healthy access point                       | 17.12.8           | `No` is a reading there rather than a substitute for silence                  |
| `ap-misconfig` on 3 of 3 records, its domain going three members to four at 17.18.4a                 | 17.15.6, 17.18.4a | Undeclared at 17.12.8, and only `apmgr-no-misconfig` ever seen on it          |
| The filter-name leaf arrives as an empty string rather than absent, on 3 of 3 records                | 17.12.8           | With no `ap-filter-configs` configured on that controller                     |
| No resolved counterpart of the AP join profile or the flex profile in the schema                     | All three         | Both columns come from the configured site tag instead                        |

## Configuration

What a write reaches, what a save persists, and the transport's own floor.

| Fact                                                                                                     | Release           | Condition                                                                   |
| :------------------------------------------------------------------------------------------------------- | :---------------- | :-------------------------------------------------------------------------- |
| `writable-running` advertised, with `:startup` and `:candidate` advertised by none                       | All three         | So the startup configuration is unreachable over RESTCONF                   |
| One RF tag over RESTCONF: running only, both after a save, startup only after a delete                   | All three         | So an unsaved delete restores a deleted tag on the next reload              |
| `cisco-ia:save-config` returns `Save running-config successful` in a 78-byte body                        | All three         | Identical over six posts, 17.12.8 serving `cisco-ia` three years behind     |
| A save takes at most 3.7s against about 0.13s for a container read                                       | All three         | The one request a read-sized `--timeout` can still refuse                   |
| `ap-cfg` holds no admin-state leaf, and a `show running-config all` filtered on the name returns nothing | 17.15.6           | Per-AP configuration keys on the dotted MAC, so that filter settles nothing |
| The legacy WLAN band setting is obsolete and reads all bands where the list holds one                    | All three         | Gone at 17.18.4a, so `bands` comes from the per-band list alone             |
| The SDK pins the dialer at 30s, `--timeout 60s` and `120s` giving up at 30                               | Not release-bound | Against an unroutable address, so a 30s failure has no route at all         |

## Not measured

Known gaps, each with its reason and what an operator carries for it.

| Gap                                                      | Why                                                                  | Consequence                                                                 |
| :------------------------------------------------------- | :------------------------------------------------------------------- | :-------------------------------------------------------------------------- |
| Deleting a tag an access point resolves to               | It would mean deleting a tag in use on the lab fleet                 | The controller keeps a dangling reference rather than refusing one          |
| A dedicated 6 GHz radio on any slot                      | The lab holds no access point carrying one                           | The band 4 pairing rests on the clause and never on a write                 |
| A username selecting more than one client session        | No lab controller holds two sessions under one username              | The prompt's count is the whole of the operator's warning                   |
| The 17.12.8 `deauth` `400` at the CLI                    | Every client there carries an empty username, so the resolve exits 1 | Pinned on fixtures under `internal/wnc` and `internal/cli` instead          |
| The client impact of a `reset capwap` teardown           | The lab access point carried no clients when this was measured       | A control teardown is not a radio reset on a locally switching access point |
| The `reset capwap` name arm on the wire                  | No write through it recorded, nor which arm carried the readings     | So neither arm settles anything at all about the other                      |
| The `ip-addr` arm of `apf-ms-delete-all`                 | Every SISF binding measured carries zone 0, the payload naming none  | Left unimplemented, an address being what `--mac` already selects by        |
| Six site-tag leaves this CLI never renders               | None verified in scope: fabric pair, image download, ARP, DHCP, load | `wnc show site-tag` renders five columns rather than eleven                 |
| An independent arbiter for an access point's admin state | `show running-config all` on the name returns nothing at 17.15.6     | Read it back with `wnc show ap` and `wnc show overview` instead             |
| The age of the channel-utilization measurement           | No module declares a timestamp for it on any release in scope        | RRM's own cycle is invisible, so a stale read looks like a fresh one        |
