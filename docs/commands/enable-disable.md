# wnc enable, wnc disable

Set the administrative state of one access point, or of one of its radios.

```bash
wnc enable ap --ap-name <ap-name>
wnc disable ap --ap-name <ap-name>
wnc enable radio --ap-name <ap-name> --slot <n>
wnc disable radio --ap-name <ap-name> --slot <n>
```

```plaintext
Disable TEST-AP01 on 192.168.0.232? This sets the access point's admin state, not one radio's. [y/N]: y
TEST-AP01 on 192.168.0.232: disable sent
```

The `radio` leaves name the slot, and the band the controller reported for it:

```plaintext
Disable slot 1 (5 GHz) of TEST-AP01 on 192.168.0.232? [y/N]: y
slot 1 (5 GHz) of TEST-AP01 on 192.168.0.232: disable sent
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
| `--slot`         | Radio slot, on the `radio` leaves only             |

## The target

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap` — resolved on the controller exactly as [reset ap](reset-ap.md#the-target) resolves it. The `radio` leaves keep the base radio address that resolve returns, because `radio-oper-data` is keyed on it and the slot RPC takes it.

## The slot

`--slot` is the `Slot` column of `wnc show overview`, and omitting it is a usage fault. It has no default: zero is a slot the controller reports, so a missing flag and a named zero must not read alike, and the help suppresses the `(default: 0)` an integer flag would otherwise print.

The band is read from the controller rather than asked for, and there are two of them. The number the RPC takes follows the **radio type** — dual band is a kind of radio and not a frequency, so a 5-or-6 GHz XOR radio takes 3 whichever band it is serving. The band the prompt names is the one the radio is **serving**, because that is the reading an operator can check against the `Band` column of `wnc show overview` before answering.

Measured in both directions, which is what settles that the number is the radio's and not the frequency's:

| Radio                       | Band sent | Answer                                                |
| --------------------------- | --------- | ----------------------------------------------------- |
| XOR in slot 2 serving 6 GHz | 4         | `400`, "AP slot: 2 does not have a dedicated radio"   |
| The same radio              | 3         | `204`, slot 2 alone went down                         |
| Dedicated 2.4 GHz in slot 0 | 3         | `400`, "AP does not support the specified radio type" |

The RPC's own `must` clause pairs the two, and a pair it would reject is refused here instead:

| Band sent | Slots the RPC accepts | The radio it names |
| --------- | --------------------- | ------------------ |
| 1         | 0                     | dedicated 2.4 GHz  |
| 2         | 1 or 2                | dedicated 5 GHz    |
| 3         | 0 or 2                | dual band or XOR   |
| 4         | 2 or 3                | dedicated 6 GHz    |

## Guards

**One access point per invocation.** `--ap-name` names one, and a second occurrence is rejected rather than silently replacing the first. There is no glob, no repeat and no fleet selector.

**One radio per invocation.** `--slot` names one slot and there is no "all radios" form. `disable ap` is the whole-access-point action, and it is a different RPC rather than a shorthand for every slot.

**One controller per invocation.** An access point lives on one controller and this CLI does not guess which, so naming two is a usage fault. Nothing is sent.

**The target is resolved before you are asked.** The controller is read for the name first, so a name it holds no access point under is refused before the prompt. The `radio` leaves read the radio as well, so the prompt carries the band the controller reported.

**A radio the RPC cannot express is refused.** A slot holding a remote-LAN port, a record carrying no radio type, a radio type the RPC has no band number for, a radio the controller reports no band for, a band spelling outside 2.4, 5 and 6 GHz, and a slot the `must` clause does not permit for that radio are each refused after the read and before the write.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It names what would change and sends no RPC. On the `radio` leaves it names the slot and the band the controller reported for it.

## Exit codes

| Code | Meaning                                                        |
| :--- | :------------------------------------------------------------- |
| 0    | The controller accepted the instruction, or you answered no    |
| 1    | The controller refused, or holds no such access point or radio |
| 2    | Usage fault. Nothing was sent to a controller                  |

A missing `--slot`, a `--slot` outside 0 to 3 and a missing name are each refused before a client exists, so they exit 2 with nothing sent. Measured on 17.12, `--slot 3` on an access point that has no slot 3 is refused only after the controller answers, which puts it at exit 1 rather than at a usage fault.

## Notes

**"disable sent" is not "disabled".** The `ap` leaves post `set-ap-admin-state` and the `radio` leaves post `set-ap-slot-admin-state`, and neither declares an output container — so a `204 No Content` establishes that the controller took the instruction and nothing more. Read the state back afterwards.

**The name arm has been measured.** The `ap` leaves post the access point's name rather than an address, and a write through that arm on 17.18 had the same effect as the earlier one through the address arm on 17.15.6 — the readings below hold for both.

**After an access-point-level disable, `wnc show ap` is the authority and `wnc show overview` is not.** Measured on 17.15.6 and again on 17.18: `ap-admin-state` flipped to `adminstate-disabled` while both radios' own `admin-state` stayed `enabled` and their `oper-state` went to `radio-down`. So `wnc show ap` reports Admin `Disabled` and `wnc show overview` still reports Admin `Enabled` with Oper `Down`. Neither reading is wrong — they are different leaves, and the access point is what was disabled.

**A disabled access point stays joined.** Measured on 17.15.6 and again on 17.18: it stayed `Registered`, stayed in `capwap-data`, and its `uptime_seconds` kept climbing across the exchange, so it did not reboot. That is the difference from [reset ap](reset-ap.md), which drops the access point out of every view except `wnc show ap-join`.

**A radio-level change touches one radio.** Measured on 17.15.6 and again on 17.18: `disable radio --slot 1` then `enable radio --slot 1` took slot 1 from `Enabled/Up` to `Disabled/Down` and back, and slot 0 was unaffected throughout. `--slot 0` behaved as the mirror image.

**Band 4 is unverified, and band 3 no longer is.** Measured on 17.15 against `TEST-AP03` slot 2, a 5 or 6 GHz XOR radio reporting `dot11-6-ghz-band`: band 3 answered `204` and took slot 2 down while slots 0 and 1 stayed up, and band 4 on the same slot answered `400`. Band 3 was measured again on 17.12 against a 2.4/5 GHz XOR radio in slot 0, where the served band's own number answered `400` first. No dedicated 6 GHz radio exists in the lab, so the band-4 row of the table above comes from the `must` clause and not from a write.

## Examples

Take one radio out of service:

```bash
wnc disable radio --ap-name TEST-AP01 --slot 1
```

Check the target and the band without acting:

```bash
wnc --dry-run disable radio --ap-name TEST-AP01 --slot 1
```

Bring an access point back and read both views:

```bash
wnc enable ap --ap-name TEST-AP01 --yes
wnc show ap -f json | jq -r '.[] | select(.ap_name=="TEST-AP01") | .admin'
wnc show overview -f json | jq -r '.[] | select(.ap_name=="TEST-AP01") | "\(.slot) \(.admin)/\(.oper)"'
```
