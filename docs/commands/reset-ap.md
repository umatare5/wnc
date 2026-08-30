# wnc reset ap

Restart one access point. Of the seven trees that act on a controller, this is the only one that restarts the access point itself.

```bash
wnc reset ap --ap-name <ap-name>
```

```plaintext
Reset TEST-AP01 on 192.168.0.231? Its clients disconnect for about four minutes. [y/N]: y
TEST-AP01 on 192.168.0.231: reset sent
```

## Options

| Option           | Meaning                                            |
| ---------------- | -------------------------------------------------- |
| `--ap-name`      | The access point's name, as shown by `wnc show ap` |
| `--controller`   | The one controller serving the access point        |
| `--access-token` | Basic auth token for that controller               |
| `--insecure`     | Skip TLS certificate verification                  |
| `--timeout`      | Request timeout                                    |
| `--yes`          | Act without the confirmation prompt                |
| `--dry-run`      | Name the target and change nothing                 |

## The target

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap`. Nothing about it is checked here beyond being present and not blank: measured on 17.18, the keyed read answers `404` for a 256-character name and for one holding a space, a slash or a multi-byte character alike, so the controller distinguishes no grammar a local check could enforce.

The controller resolves it before the RPC. The read is `ap-name-mac-map=<name>`, a list keyed on the name that answers one row carrying the name, the base radio address and the Ethernet address on every release in scope. A name no access point holds answers `404`, reported on 17.12 as `wnc: 192.168.0.231 holds no access point named NO-SUCH-AP` and followed by no RPC.

Whatever is typed reaches the controller as a path segment of that read before it is reported.

## Guards

**One access point per invocation.** `--ap-name` names one, and a second occurrence is rejected rather than silently replacing the first. There is no glob, no repeat and no fleet selector.

**One controller per invocation.** An access point lives on one controller and this CLI does not guess which, so naming two is a usage fault. Nothing is sent.

**The target is resolved before you are asked.** The controller is read for the name first, so a name it holds no access point under is refused before the prompt. What the prompt still catches is the controller beside the name.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It names what would be restarted and sends no RPC. Reading is not changing, so it does contact the controller.

## Exit codes

| Code | Meaning                                                       |
| :--- | :------------------------------------------------------------ |
| 0    | The controller accepted the instruction, or you answered no   |
| 1    | The controller refused, or holds no access point of that name |
| 2    | Usage fault. Nothing was sent to a controller                 |

## Notes

**"reset sent" is not "rebooted".** `ap-reset` declares no output container, so a `204 No Content` establishes that the controller took the instruction and nothing further. What happened next is a separate reading.

**`wnc show ap-join` is that reading, and it is the only one.** The controller drops a restarting access point from `capwap-data`, so `wnc show ap` loses it while `wnc show ap-join` keeps it. Measured on 17.12.8: the access point leaves within 15 seconds, `show ap-join` reports `Not Joined` with `Wtp reset config cmd sent` as the disconnect reason, and it rejoins about 285 seconds later with `reboot_reason` reading `ap-reboot-reason-reboot-cmd`.

**Clients are disconnected for the whole of that.** Roughly four minutes on an AIR-AP1815I, of which about 216 seconds is the join itself.

**An access point that is not joined cannot be restarted.** The instruction travels over the CAPWAP session, so a name the controller holds no access point under is refused before anything is sent.

**A `404` and an answer carrying no row are reported alike.** The keyed read answers `404` for a name no access point holds on 17.12, 17.15 and 17.18, and a `200` carrying no row is treated as the same absence rather than as an access point with no address.

**The name arm has been measured.** This leaf posts the access point's name rather than an address, and a write through that arm on 17.18 restarted the access point named and nothing else: it left `capwap-data` within sixteen seconds, came back with a new boot time, and rejoined four minutes and thirty-nine seconds after the post on an AIR-AP1815I.

## Examples

Restart the access point a weak-signal client is on — the `ap_name` column of `wnc show client`:

```bash
wnc reset ap --ap-name TEST-AP01
```

Check the target without acting:

```bash
wnc --dry-run reset ap --ap-name TEST-AP01
```

Watch it leave and come back:

```bash
wnc reset ap --ap-name TEST-AP01 --yes
watch -n 10 'wnc show ap-join -f json | jq -r ".[] | select(.ap_name==\"TEST-AP01\") | .status"'
```
