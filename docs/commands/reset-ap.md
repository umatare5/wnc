# wnc reset ap

Restart one access point.

```bash
wnc reset ap --ap-name <ap-name>
```

```plaintext
Reset TEST-AP01 on WNC1? Its clients disconnect for about four minutes. [y/N]: y
TEST-AP01 on WNC1: reset sent
```

It is the only command here that restarts the access point itself.

## Flags

| Option      | Meaning                                            |
| :---------- | :------------------------------------------------- |
| `--ap-name` | The access point's name, as shown by `wnc show ap` |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap`. Nothing about it is checked here beyond being present and not blank: measured on 17.18.4a, the keyed read answers `404` for a 256-character name and for one holding a space, a slash or a multi-byte character alike, so the controller distinguishes no grammar a local check could enforce.

The controller resolves it before the RPC. The read is `ap-name-mac-map=<name>`, a list keyed on the name that answers one row carrying the name, the base radio address and the Ethernet address on every release in scope. A name no access point holds answers `404`, reported on 17.12.8 as `wnc: WNC1 holds no access point named NO-SUCH-AP` and followed by no RPC.

Whatever is typed reaches the controller as a path segment of that read before it is reported.

## Guards

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the controller holds no access point of that name.

## Notes

- **[`show ap-join`](show-ap-join.md) is the only view that keeps it** — `show ap` loses it as it leaves
- **Clients are down until it rejoins** — the whole restart, not a control-session blip
- **An unjoined access point cannot be restarted** — the instruction rides the CAPWAP session

The readings behind these sit in [`measurements.md`](../measurements.md).

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
