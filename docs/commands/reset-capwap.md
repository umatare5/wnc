# wnc reset capwap

Reset one access point's controller session. The access point does not reboot — only the CAPWAP session between it and the controller is torn down and re-established.

```bash
wnc reset capwap --ap-name <ap-name>
```

```plaintext
Reset the CAPWAP session of TEST-AP01 on 192.168.0.231? It rejoins within about ten seconds and does not reboot. [y/N]: y
TEST-AP01 on 192.168.0.231: capwap reset sent
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

`--ap-name` is the access point's name — the `ap_name` column of `wnc show ap` — resolved on the controller exactly as [reset ap](reset-ap.md#the-target) resolves it.

## How this differs from reset ap

| Reading                | `reset capwap` | `reset ap`   |
| ---------------------- | -------------- | ------------ |
| Rejoin time            | under 10s      | 260s to 285s |
| `uptime_seconds`       | unchanged      | restarts     |
| `assoc_uptime_seconds` | restarts       | restarts     |

Measured on 17.12.8 against an AIR-AP1815I. `wnc show ap` reports both quantities, and the pair is the evidence: after a CAPWAP reset the access point's own uptime keeps climbing while its association age returns to a few seconds. After a reset ap both restart together.

The controller's `ap-join-stats` counters also move. `num-join-req-recvd`, `num-config-req-recvd` and `ctrl-dtls-setup-req` each accumulate across sessions, so a reset is visible as a delta rather than only as a change of last-state fields.

## Guards

**One access point per invocation.** `--ap-name` names one, and a second occurrence is rejected rather than silently replacing the first. There is no glob, no repeat and no fleet selector.

**One controller per invocation.** An access point lives on one controller and this CLI does not guess which, so naming two is a usage fault. Nothing is sent.

**The target is resolved before you are asked.** The controller is read for the name first, so a name it holds no access point under is refused before the prompt.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It names what would be reset and sends no RPC.

## Exit codes

| Code | Meaning                                                       |
| :--- | :------------------------------------------------------------ |
| 0    | The controller accepted the instruction, or you answered no   |
| 1    | The controller refused, or holds no access point of that name |
| 2    | Usage fault. Nothing was sent to a controller                 |

## Notes

**"capwap reset sent" is not "rejoined".** `set-rad-capwap-reset` declares no output container, so a `204 No Content` establishes that the controller took the instruction and nothing further. Read `wnc show ap` for the association age, or `wnc show ap-join` for the join state.

**Client impact was not measured.** The lab access point carried no clients when this was measured, so the CLI does not state what happens to them. A control-session teardown is not a radio reset, and on a locally switching access point the two are not the same thing.

**The name arm has no recorded write.** This leaf posts the access point's name rather than an address, and no write through that arm is recorded as sent on any release. Which arm carried the readings above is not recorded either, so neither settles anything about the other.

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
