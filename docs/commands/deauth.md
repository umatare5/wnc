# wnc deauth

Deauthenticate a client on a controller, by address or by username. The controller deletes the client's record and the station re-associates on its own within seconds — nothing about the access point, its radios or its CAPWAP session is touched.

```bash
wnc deauth --mac <mac>
wnc deauth --username <username>
```

```plaintext
Deauthenticate aa:bb:cc:dd:ee:ff on 192.168.0.233? It is dropped and reconnects on its own. [y/N]: y
aa:bb:cc:dd:ee:ff on 192.168.0.233: deauthenticate sent
```

## Options

| Option           | Meaning                                              |
| ---------------- | ---------------------------------------------------- |
| `--mac`          | The client's address, as shown by `wnc show client`  |
| `--username`     | The client's username, as shown by `wnc show client` |
| `--controller`   | The one controller the client is associated to       |
| `--access-token` | Basic auth token for that controller                 |
| `--insecure`     | Skip TLS certificate verification                    |
| `--timeout`      | Request timeout                                      |
| `--yes`          | Act without the confirmation prompt                  |
| `--dry-run`      | Name the target and change nothing                   |

## The target

`--mac` and `--username` are the `mac` and `username` columns of `wnc show client`, and one invocation gives one of them. They are the RPC's own two usable arms rather than a convenience: a client carries no name on the wire, so an address or the identity it authenticated under is all there is to select by.

**`--username` is not `WNC_USERNAME`.** That variable is the controller login `wnc generate-token` reads, and this flag deliberately ignores it — a client's username and a controller account's are unrelated values that happen to share a spelling.

Either target is resolved on the controller before the prompt. For an address the row's own spelling is what the prompt, the report and the wire all name — measured on 17.12.8 and 17.15.6, every `client-mac` a controller serves is already lowercase, which is the form the SDK normalizes to. For a username the resolve adds the one thing the controller can tell you that you did not type: **how many sessions carry it.**

```plaintext
Deauthenticate 2 clients authenticated as someone@example.net on 192.168.0.233? Each is dropped and reconnects on its own. [y/N]:
```

## What it does

The RPC is `apf-ms-delete-all`, and its name reads as a fleet-wide purge that its schema description contradicts: "Delete wireless client based on client MAC or IP address or username". Measured on 17.18.4a against a real client:

| Reading         | Before  | After the post   |
| --------------- | ------- | ---------------- |
| `assoc_seconds` | 6102    | 13               |
| `state`         | Run     | IP Learning      |
| `ipv4`          | present | back within 210s |

**It is not promised to affect exactly one client.** On that controller the other three were untouched, but on 17.15.6 two posts at a 19-client controller each dropped the target and one other station — same access point, same BSSID, same radio, same WLAN, different vendor — while the other thirteen on that BSS stayed put. A 25-second window with no post moved none of the nineteen, so the second station was not a coincidence. Read `wnc show client` afterwards rather than assuming the blast radius.

**The record is deleted, not reset.** The RPC's name is accurate on this point. Measured on 17.15.6 through the username arm, against a client whose association had been unchanged for 82 minutes:

| Reading              | Before  | 1.9s after the post |
| -------------------- | ------- | ------------------- |
| the client's row     | present | **gone**            |
| `assoc_seconds`      | 4943    | —                   |
| the other 17 clients | —       | 0 moved             |

The station re-associated 1.5 seconds after the post returned, as a new record, and was back at `Run` with its username inside 42 seconds. So a read taken a few seconds later sees a young association rather than an absence, which is what the 17.18.4a row above caught.

**The username arm moved nothing else.** All seventeen other clients held their association time and their state across the post — including the one other client on that controller carrying a username. Two no-post windows of 35 seconds each moved 1 and 0 of 18 respectively, so the estate does have natural churn and a single moved client proves nothing on its own. What licenses the attribution here is the target's own 82 minutes of stability.

**How many sessions a username selects is the controller's decision, not this CLI's.** One post carries the username and the controller resolves it — there is no loop over the resolved clients. The count in the prompt is what the read is for.

## Release availability

| Release   | The operation                                 |
| --------- | --------------------------------------------- |
| `17.12.x` | Absent — the module declares no client delete |
| `17.15.x` | Served                                        |
| `17.18.x` | Served                                        |

Measured 2026-08-29. 17.12.8 serves `Cisco-IOS-XE-wireless-client-rpc` at revision `2023-03-01`, which declares two `clear-sisf-binding` operations and nothing else, while 17.15.6 and 17.18.4a serve `2024-03-01`, which adds this one. A post to 17.12.8 answers `400` with `error-tag: malformed-message` and `error-message: "invalid path"`, and this command re-words that status rather than reporting it bare.

## Guards

**One arm per invocation.** Naming both `--mac` and `--username` is a usage fault — the RPC's choice is mandatory and the controller resolves the first arm it finds, so sending both would let the controller pick which one you meant.

**One target per invocation.** Either flag given twice is rejected rather than silently replacing the first value. There is no glob, no repeat and no fleet selector.

**Neither target may be empty.** An empty address reaches the SDK's own not-found and would surface as a read failure. An empty username is worse: it is the value most clients carry — sixteen of eighteen on 17.15.6 — so it would select nearly the whole fleet.

**One controller per invocation.** A client is associated to one controller and this CLI does not guess which, so naming two is a usage fault. Nothing is sent.

**The target is resolved before you are asked.** This is the guard that matters here, not a courtesy: the RPC answers `204` for an identifier associated to nothing exactly as it does for a session it dropped, so without the read a reported deauthentication and a mistyped target would be the same output.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It doubles as an existence probe, which is the only thing it can report: the RPC's answer says nothing either way. On the username arm it also reports the session count.

## Exit codes

| Code | Meaning                                                     |
| :--- | :---------------------------------------------------------- |
| 0    | The controller accepted the instruction, or you answered no |
| 1    | The controller refused, or holds no client at that target   |
| 2    | Usage fault. Nothing was sent to a controller               |

## Notes

**"deauthenticate sent" is not "deauthenticated".** `apf-ms-delete-all` declares no output container, so a `204 No Content` establishes that the controller took the instruction and nothing further. Read `wnc show client` for the association age and the state.

**The IP arm is not implemented.** The RPC's choice also accepts `ip-addr`, which answered `204` on 17.15.6 and is not refused. It is left out because every SISF binding measured carries zone 0 and the payload names no other zone, so the arm cannot be exercised outside the one case — and because an address is what `--mac` already selects by.

**`zone-id` is not sent.** The schema declares it defaulting to `0`, so the default is in force. No lab controller has a second zone, and a leaf whose value cannot be checked against a device is not one to expose.

## Examples

Drop a client stuck short of `Run`:

```bash
wnc show client -f json | jq -r '.[] | select(.state != "Run") | .mac'
wnc deauth --mac aa:bb:cc:dd:ee:ff
```

Drop every session a user holds:

```bash
wnc show client -f json | jq -r '.[] | select(.username) | .username' | sort -u
wnc deauth --username someone@example.net
```

Check the target without acting, on either arm:

```bash
wnc --dry-run deauth --mac aa:bb:cc:dd:ee:ff
wnc --dry-run deauth --username someone@example.net
```

Drop a client and watch it come back:

```bash
wnc deauth --mac aa:bb:cc:dd:ee:ff --yes
watch -n 5 'wnc show client -f json | jq -r ".[] | select(.mac==\"aa:bb:cc:dd:ee:ff\") | \"\(.state) assoc \(.assoc_seconds)s\""'
```
