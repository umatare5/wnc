# wnc enable, wnc disable

Set the administrative state of one access point, or of one of its radios.

```bash
wnc enable ap --ap-name <ap-name>
wnc disable ap --ap-name <ap-name>
wnc enable radio --ap-name <ap-name> --slot <n>
wnc disable radio --ap-name <ap-name> --slot <n>
```

```plaintext
Disable TEST-AP01 on WNC2? This sets the access point's admin state, not one radio's. [y/N]: y
TEST-AP01 on WNC2: disable sent
```

The `radio` leaves name the slot, and the band the controller reported for it:

```plaintext
Disable slot 1 (5 GHz) of TEST-AP01 on WNC2? [y/N]: y
slot 1 (5 GHz) of TEST-AP01 on WNC2: disable sent
```

## Flags

| Option      | Meaning                                            |
| :---------- | :------------------------------------------------- |
| `--ap-name` | The access point's name, as shown by `wnc show ap` |
| `--slot`    | Radio slot, on the `radio` leaves only             |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap` — resolved on the controller exactly as [reset ap](reset-ap.md#the-target) resolves it. The `radio` leaves keep the base radio address that resolve returns, because `radio-oper-data` is keyed on it and the slot RPC takes it.

## Guards

**One radio per invocation.** `--slot` names one slot and there is no "all radios" form. `disable ap` is the whole-access-point action, and it is a different RPC rather than a shorthand for every slot.

**The target is resolved before you are asked.** The controller is read for the name first, so a name it holds no access point under is refused before the prompt. The `radio` leaves read the radio as well, so the prompt carries the band the controller reported.

**A radio the RPC cannot express is refused.** A slot holding a remote-LAN port, a record carrying no radio type, a radio type the RPC has no band number for, a radio the controller reports no band for, a band spelling outside 2.4, 5 and 6 GHz, and a slot the `must` clause does not permit for that radio are each refused after the read and before the write.

**`--dry-run` names the band as well as the slot.** On the `radio` leaves it names the slot and the band the controller reported for it.

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the controller holds no such access point or radio.

A missing `--slot`, a `--slot` outside 0 to 3 and a missing name are each refused before a client exists, so they exit 2 with nothing sent. Measured on 17.12.8, `--slot 3` on an access point that has no slot 3 is refused only after the controller answers, which puts it at exit 1 rather than at a usage fault.

## Notes

- **`--slot` has no default** — zero is a real slot, so a missing flag must not read as one
- **The band number follows the radio type** — dual band is a kind of radio, not a frequency
- **The prompt names the band served** — that is the reading to check against `show overview`
- **After an access-point disable, [`show ap`](show-ap.md) is the authority** — `show overview` reads Enabled
- **A disabled access point stays joined** — the difference from [`reset ap`](reset-ap.md), which drops it

The readings behind these sit in [`measurements.md`](../measurements.md).

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
