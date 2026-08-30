# wnc save-config

Copy the controller's running configuration to its startup configuration. Without it every tag this CLI writes is lost when the controller reloads.

```bash
wnc save-config
```

```plaintext
Save the running configuration of 192.168.0.231? Every change on the controller is persisted, including changes this CLI did not make. [y/N]: y
192.168.0.231: running configuration saved
```

## Options

| Option           | Meaning                                |
| ---------------- | -------------------------------------- |
| `--controller`   | The one controller to save             |
| `--access-token` | Basic auth token for that controller   |
| `--insecure`     | Skip TLS certificate verification      |
| `--timeout`      | Request timeout                        |
| `--yes`          | Act without the confirmation prompt    |
| `--dry-run`      | Name the controller and change nothing |

## Why this exists

A RESTCONF write lands in the running configuration and nowhere else. All three releases in scope advertise `writable-running` and advertise neither `:startup` nor `:candidate`, so there is no candidate datastore to commit from and no startup datastore to write to.

Measured on 17.12.8, 17.15.6 and 17.18.4a by creating one RF tag over RESTCONF:

| Step                      | Running | Startup |
| ------------------------- | ------- | ------- |
| After `wnc set rf-tag`    | present | absent  |
| After `wnc save-config`   | present | present |
| After `wnc delete rf-tag` | absent  | present |

**A delete is as unsaved as a create.** The third row is the half that surprises — an unsaved delete means a reload does not merely lose a new tag, it brings a deleted one back.

## What a save covers

| Written by                  | Persisted by a save |
| --------------------------- | ------------------- |
| `wnc set`, `wnc delete`     | Yes                 |
| `wnc enable`, `wnc disable` | No                  |

An access point's administrative state is no part of the configuration: `ap-cfg` carries no admin-state leaf and `show running-config all` prints no line for it, so there is nothing for a save to persist. It is readable back as access-point state instead, which is what `wnc show ap` reports.

## Guards

**One controller per invocation.** A save names one controller, so naming two is a usage fault and nothing is sent. A fleet is saved one at a time.

**The startup configuration is the only destination.** The RPC takes no argument at all, so no file may be named. Writing a configuration to a named file is an upgrade procedure rather than day-to-day operation, and it belongs to [telee](https://github.com/umatare5/telee).

**Everything on the controller is persisted.** The RPC names no scope, so a save commits whatever else is in the running configuration — another operator's change, a half-finished experiment, anything a reload was going to discard. The prompt says so, and this is the reason it does.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` names the controller and stops.** It cannot report whether anything needs saving, because the startup configuration is unreachable over RESTCONF.

## Exit codes

| Code | Meaning                                                    |
| :--- | :--------------------------------------------------------- |
| 0    | The controller saved the configuration, or you answered no |
| 1    | The controller refused, or reported no result              |
| 2    | Usage fault. Nothing was sent to a controller              |

## Notes

**"saved" is a completion, unlike every other write here.** `cisco-ia:save-config` is the one RPC this CLI posts whose schema declares an output container, and it answers with the controller's own account of the save. Six posts across the three releases each returned `Save running-config successful` in a 78-byte body, byte-identical even though 17.12 serves the `cisco-ia` module three years behind the other two. An answer carrying no result is reported as a failure rather than as a save.

**It takes two to four seconds.** Measured between 2.5s and 3.7s, against about 0.13s for a container read. A `--timeout` every read survives can still refuse this one, which is the only place in this CLI where that is true.

**This CLI cannot tell you whether a save is needed.** The startup configuration is not reachable over RESTCONF, so nothing here can compare the two. The controller itself can: `show running-config` heads its output with `Last configuration change` and `NVRAM config last updated`, and a running timestamp newer than the NVRAM one means there is something to save.

## Examples

Write a tag and make it survive a reload:

```bash
wnc set rf-tag --name labo-inside --profile-5ghz labo-rf-5gh-inside --yes
wnc save-config --yes
```

Check which controller would be saved without saving it:

```bash
wnc --dry-run save-config -c 192.168.0.231
```

Save every controller in a file, one at a time:

```bash
for host in 192.168.0.231 192.168.0.232 192.168.0.233; do
  wnc save-config -c "$host" --yes
done
```
