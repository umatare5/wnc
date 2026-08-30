# Troubleshooting

A read that fails is logged to stderr as one line carrying a `cause` field, and the sections below are indexed by it. A usage fault and a refusal from a command that acts carry no such field: they are one `wnc: …` line, indexed here by its wording.

```plaintext
error: 192.168.0.231: the controller answered 401 Unauthorized (cause=auth)
```

The controller leads the sentence, and `endpoint=` joins the cause inside the parentheses when the failure cost only some cells rather than the whole row set.

**The sentence is this CLI's own, not the error's.** Every cause the controller or the transport answers with gets one clause of the same shape, so a failure reads alike whether the controller answered it, the transport raised it, or the SDK wrapped it in a sentence of its own. What those clauses replace is a controller error document that can echo the configuration that was read, and a `*url.Error` naming the whole request URL. `cause=internal` is the exception and quotes its error, because that is how this CLI's own re-worded refusals reach you.

`--log-level debug` renders the same record as logfmt, makes the controller a field of its own, adds the HTTP status beside it, and is where the Go error itself appears — as the SDK's own record, one line earlier:

```plaintext
level=debug msg="HTTP request failed" error="Get \"https://192.0.2.1/restconf/data/...\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)" url="https://192.0.2.1/restconf/data/..."
level=error msg="the controller did not answer in time" cause=timeout controller=192.0.2.1 status=0
```

The fatal line names no endpoint because it does not have to: each view has exactly one read it cannot proceed without, so the command names it. A read that cost only some cells is reported separately and carries `endpoint=`.

## cause=auth

The controller rejected the token with 401.

Regenerate it and check the account is still valid on the controller:

```bash
printf '%s' "$WNC_PASSWORD" | wnc generate-token -u admin
```

A token is `base64(username:password)`. A trailing newline inside it, from `echo` rather than `printf`, is the usual cause of a token that looks right and is not.

## cause=forbidden

The account authenticated but is not permitted to read. RESTCONF needs privilege level 15 on IOS-XE.

## cause=not-found

The controller has no such node. This is normal on a release that dropped a leaf, and the CLI treats it as an empty result where the read is optional. Seeing it as a failure means the whole collection is missing, which usually means RESTCONF is enabled but the wireless feature set is not.

## cause=timeout

No answer within `--timeout`, which defaults to 60 seconds per request.

**A failure at 30 seconds is the connect stage, and `--timeout` does not move it.** The SDK pins the dialer at 30 seconds and exports no option for it — measured, `--timeout 60s` and `--timeout 120s` both give up at 30 against an unroutable address. So a run that fails at 30 seconds every time has no route, and raising the flag will not change the outcome.

The wall-clock ceiling for one controller is the timeout multiplied by the number of sequential reads the command makes — five for `show overview` and `show client`, three for `show ap` and `show wlan`, one for every other view. Controllers are read concurrently, so adding one does not extend the ceiling.

`wnc save-config` is the one command a read-sized timeout refuses. It took 2.5s to 3.7s on every release measured, against about 0.13s for a container read, so a `--timeout` that every view survives can still fail the save.

A controller under load answers a whole-container read slowly. Raise the timeout before suspecting the network:

```bash
wnc show client --timeout 120s
```

## cause=tls

The certificate did not verify. See [SECURITY.md](./SECURITY.md#trusting-a-private-certificate-authority) for trusting the issuer, which is preferable to `--insecure`.

## cause=connection

The host did not resolve, or refused the connection. Check the authority is reachable on 443:

```bash
curl -kIsS "https://<host>/restconf/" -o /dev/null -w '%{http_code}\n'
```

## cause=http

The controller answered, and with a status none of the causes above name — a `400` from a rejected query parameter and a `5xx` from a controller under load both land here. The sentence carries the code, and `--log-level debug` puts it in a `status` field of its own:

```bash
curl -kisS -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        -H "Accept: application/yang-data+json" \
        "https://<host>/restconf/data/<the endpoint the log named>" | head -1
```

## cause=canceled

The run was interrupted. Nothing is printed, and the exit code is 130.

## cause=internal

A fault in the CLI rather than in the exchange. Re-run with `--log-level debug`, which logs every request with its host, endpoint, status and duration, and open an issue with that output. The debug log never carries the token or any header.

## Everything shows `-`

`-` means the controller reported nothing for that cell, which is not the same as zero. A Monitor radio has no channel, an unauthenticated client has no username, and a WLAN with no policy tag bound has no interface. If a whole column is `-` and you expect values, look for an `endpoint=` line: a secondary read that failed leaves exactly its own columns empty.

## The row set is smaller than expected

Two commands hide entries on purpose:

- `show overview` omits a remote-LAN port, which the controller reports among the radios but which has no RF
- `show wlan` omits a policy-tag binding whose WLAN profile does not exist, and reports the count as a warning

`show client -r` and `show overview -r` report how many rows they dropped because the controller reported no band for them. The `--ssid` and `--ap-name` filters report no such count.

## Exit code 3 with a full table

At least one read failed while at least one succeeded. The table holds everything that was read — the stderr lines say what was not.

## A command that acts says the controller holds no access point of that name

The name is resolved against the controller before anything is sent, through `ap-name-mac-map`, which is keyed on the access point's name and answers 404 for a name no access point holds. Every leaf of `reset`, `enable` and `disable` resolves its target through that read, so all six answer this way. Two things produce it, and `wnc show ap-join` separates them.

- The access point is on another controller, so name that one with `--controller`
- The name is not the one the controller holds — take it from the `ap_name` column of `wnc show ap`

## `enable radio` or `disable radio` says the controller reports no radio address for it

Different from the message above: the controller does hold that access point and sent no address in the row that resolved the name. `radio-oper-data` is keyed on that address, so the slot cannot be read and nothing is sent — read `wnc show ap` for the record as the controller has it.

## `enable radio` or `disable radio` says the controller holds no radio in that slot

The keyed read of `radio-oper-data` returned nothing for that slot. Read the `Slot` column of `wnc show overview` for the slots that exist — measured on 17.15.6, `--slot 3` against an access point with two radios exits 1, which is what separates it from the usage faults that exit 2 having read nothing.

## `enable radio` or `disable radio` says the slot is a remote-LAN port

The controller lists a remote-LAN port among the radios and it carries neither a band nor an admin state, so there is nothing to set. `wnc show overview` drops such an entry rather than printing a row of dashes, which is why the slot appears in no table.

## `enable radio` or `disable radio` says the RPC has no band number for a radio type

The number the RPC takes follows the radio type, and its domain holds four members. A slot reporting `radio-uwb`, `radio-invalid`, `radio-80211-xor-24-6ghz` — which fits both the dual-band 3 and the 6 GHz 4 — or a spelling a later release adds is refused rather than given a guessed number, because a wrong number reaches a radio the operator did not name. A record carrying no `radio-type` at all is refused separately, with `reports no radio type for slot <n>`.

## `enable radio` or `disable radio` says the band is unknown or unreported

That refusal is about the prompt rather than the wire. The band the prompt names is the one the radio is serving, and it is what an operator checks against the `Band` column of `wnc show overview` before answering — so a radio the controller reports no serving band for is refused even where the RPC has a number for it. Two messages cover it: `reports no band for slot <n>` where the leaf was absent, and `reports slot <n> of <ap> on an unknown band (<spelling>)` where it carried a spelling outside 2.4, 5 and 6 GHz.

## `enable radio` or `disable radio` says the radio is accepted on other slots only

The RPC's own `must` clause pairs each band number with a slot range: 1 on slot 0, 2 on slot 1 or 2, 3 on slot 0 or 2, and 4 on slot 2 or 3. Measured on 17.15.6, band 1 with slot 1 answers `400` naming that clause, so the CLI refuses the pair instead of sending it. The range follows the radio type and does not move when a dual-band radio changes the band it serves, so read the `Slot` column of `wnc show overview` and not its `Band` column to predict this one.

## A command that acts refuses a piped stdin

There is no terminal to answer the prompt on, so nothing was sent. Pass `--yes` to act, or `--dry-run` to name the target and stop. This covers every leaf of `reset`, `enable`, `disable`, `set` and `delete`, and `save-config` and `deauth` as well.

## `set` says the tag exists and there is nothing to change

The name is already on the controller and no binding flag was given, so there was nothing to send. Name a field — `wnc set rf-tag --help` lists the ones that kind takes.

## `delete` says the controller holds no tag of that name

The name was read before the delete, and the controller does not have it. Tag names are case sensitive and each kind is a separate list, so a policy tag and an RF tag may share a name without being related. `wnc show ap-tag` names the tags in force; the controller's own configuration holds the ones nothing resolves to.

## A tag write is refused with 400

The controller rejected the payload. The commonest cause is a name the device's own pattern refuses, but the CLI checks that first, so a 400 that reaches you is usually a binding the release does not accept. Read the node back to see what it holds:

```bash
curl -kisS -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        -H "Accept: application/yang-data+json" \
        "https://<host>/restconf/data/Cisco-IOS-XE-wireless-rf-cfg:rf-cfg-data/rf-tags/rf-tag=<name>"
```

A profile name that exists nowhere is **not** a cause: the configuration modules declare no `leafref`, so the controller accepts a dangling reference and keeps it.

## `save-config` says the controller reported no result

The controller answered the save but carried no `result` string, which is its whole account of what it did. Six posts across 17.12.8, 17.15.6 and 17.18.4a each returned `Save running-config successful`, so an answer without one is a shape no release produced, and reporting a save that may not have happened is the one answer this command must not give.

Whether the configuration was in fact saved is readable on the controller itself: `show running-config` heads its output with `Last configuration change` and `NVRAM config last updated`, and a running timestamp newer than the NVRAM one means there is still something to save.

## A command says it takes no positional arguments

Every value a subcommand takes is named by a flag, so a bare word is a fault wherever it lands. The message reports how many were given and never repeats them — a leftover word on `generate-token` can be the password a wrapper misplaced. Nothing was sent.

Where the leaf has a target flag the message names it: `use --ap-name` on the six leaves of `reset`, `enable` and `disable`, `use --name` on the six of `set` and `delete`, and `use --mac` on `deauth`. `wnc show client` takes `--ap-name` as a filter, so it names that one too. Elsewhere `--help` lists the flag the value belonged to.

## A tag name is refused for a leading or trailing space

The key leaf's own pattern refuses one, and the name now reaches that check with its spaces intact: a flag value is not trimmed, where a positional argument was. Quote what you meant, or drop the padding — an inner space is legal and passes.

## A message ends with `see 'wnc … --help'`

The path names the command whose help lists what was refused, which for a leaf is the leaf's own help rather than its group's. That help carries an `INHERITED OPTIONS:` section, so a connection flag the leaf parses through its parent is listed there too.

A fault the parse never reached carries no such suffix. `--slot 9: accepted values are 0 to 3` is decided by the command after its flags parsed, so it ends at the value.

## A message offers `(did you mean X?)`

The suggestion is one of this CLI's own names and never a repetition of what was typed. `--formt` scores `--format`, `resett` scores `reset` and `show cleint` scores `client`. A word too far from anything is left without a hint rather than given a wrong one, which is why `wnc help bogus` offers none.

## `invalid value for flag -t` names no value

That is deliberate and nothing an operator needs was dropped. `-t` is the short name of `--timeout` while `--access-token` has none, so `wnc show -t <token> ap` is a plausible slip — and the parser's own text would have quoted the value twice, so the flag is named and the value is not. Pass the token on `--access-token`, or better, keep it out of the argument list altogether.
