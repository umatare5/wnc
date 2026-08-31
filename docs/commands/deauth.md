# wnc deauth

Deauthenticate a client on a controller, by address or by username.

```bash
wnc deauth --mac <mac>
wnc deauth --username <username>
```

```plaintext
Deauthenticate 00:00:5e:00:53:a1 on WNC3? It is dropped and reconnects on its own. [y/N]: y
00:00:5e:00:53:a1 on WNC3: deauthenticate sent
```

The controller deletes the client's record and the station re-associates on its own. Nothing about the access point, its radios or its CAPWAP session is touched.

## Flags

| Option       | Meaning                                              |
| :----------- | :--------------------------------------------------- |
| `--mac`      | The client's address, as shown by `wnc show client`  |
| `--username` | The client's username, as shown by `wnc show client` |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--mac` and `--username` are the `mac` and `username` columns of `wnc show client`, and one invocation gives one of them. They are the RPC's own two usable arms rather than a convenience: a client carries no name on the wire, so an address or the identity it authenticated under is all there is to select by.

**`--username` is not `WNC_USERNAME`.** That variable is the controller login `wnc generate-token` reads, and this flag deliberately ignores it — a client's username and a controller account's are unrelated values that happen to share a spelling.

Either target is resolved on the controller before the prompt. For an address the row's own spelling is what the prompt, the report and the wire all name — measured on 17.12.8 and 17.15.6, every `client-mac` a controller serves is already lowercase, which is the form the SDK normalizes to. For a username the resolve adds the one thing the controller can tell you that you did not type: **how many sessions carry it.**

```plaintext
Deauthenticate 2 clients authenticated as someone@example.net on WNC3? Each is dropped and reconnects on its own. [y/N]:
```

## Guards

**One arm per invocation.** Naming both `--mac` and `--username` is a usage fault — the RPC's choice is mandatory and the controller resolves the first arm it finds, so sending both would let the controller pick which one you meant.

**Neither target may be empty.** An empty address reaches the SDK's own not-found and would surface as a read failure. An empty username is worse: it is the value most clients carry — sixteen of eighteen on 17.15.6 — so it would select nearly the whole fleet.

**The target is resolved before you are asked.** This is the guard that matters here, not a courtesy: the RPC answers `204` for an identifier associated to nothing exactly as it does for a session it dropped — measured on 17.18.4a, both in under 330ms — so without the read a reported deauthentication and a mistyped target would be the same output.

**`--dry-run` doubles as an existence probe.** That is the only thing it can report: the RPC's answer says nothing either way. On the username arm it also reports the session count.

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the controller holds no client at that target.

## Notes

- **Not promised to affect exactly one client** — read [`show client`](show-client.md) rather than assuming
- **The record is deleted, not reset** — a read seconds later sees a young association
- **The controller resolves the username** — one post, and the prompt's count is the warning
- **Not served on 17.12.8** — that release's client RPC module declares no client delete at all
- **The IP arm is not implemented** — an address is what `--mac` already selects a client by
- **`zone-id` is not sent** — its schema default is 0, and no lab controller holds a second zone

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Drop a client stuck short of `Run`:

```bash
wnc show client -f json | jq -r '.[] | select(.state != "Run") | .mac'
wnc deauth --mac 00:00:5e:00:53:a1
```

Drop every session a user holds:

```bash
wnc show client -f json | jq -r '.[] | select(.username) | .username' | sort -u
wnc deauth --username someone@example.net
```

Check the target without acting, on either arm:

```bash
wnc --dry-run deauth --mac 00:00:5e:00:53:a1
wnc --dry-run deauth --username someone@example.net
```

Drop a client and watch it come back:

```bash
wnc deauth --mac 00:00:5e:00:53:a1 --yes
watch -n 5 'wnc show client -f json | jq -r ".[] | select(.mac==\"00:00:5e:00:53:a1\") | \"\(.state) assoc \(.assoc_seconds)s\""'
```
