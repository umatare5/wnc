# wnc reset capwap

Reset one access point's controller session.

```bash
wnc reset capwap --ap-name <ap-name>
```

```plaintext
Reset the CAPWAP session of TEST-AP01 on WNC1? It rejoins within about ten seconds and does not reboot. [y/N]: y
TEST-AP01 on WNC1: capwap reset sent
```

The access point does not reboot — only the CAPWAP session between it and the controller is torn down and re-established.

## Flags

| Option      | Meaning                                            |
| :---------- | :------------------------------------------------- |
| `--ap-name` | The access point's name, as shown by `wnc show ap` |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap` — resolved on the controller exactly as [reset ap](reset-ap.md#the-target) resolves it.

## Guards

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the controller holds no access point of that name.

## Notes

- **The control session restarts and the access point does not** — its uptime keeps climbing
- **Read [`show ap`](show-ap.md) for the pair** — the association age returns to seconds, the uptime holds
- **Client impact was not measured** — a control teardown is not a radio reset on flex

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Reset the session of the access point a stuck join is suspected on:

```bash
wnc reset capwap --ap-name TEST-AP01
```

Check the target without acting:

```bash
wnc --dry-run reset capwap --ap-name TEST-AP01
```

Watch the association age reset while the uptime keeps climbing:

```bash
wnc reset capwap --ap-name TEST-AP01 --yes
watch -n 5 'wnc show ap -f json | jq -r ".[] | select(.ap_name==\"TEST-AP01\") | \"up \(.uptime_seconds)s assoc \(.assoc_uptime_seconds)s\""'
```
