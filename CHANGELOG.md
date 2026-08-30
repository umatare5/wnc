# Changelog

All notable changes to this project are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0]

First release of the rebuilt CLI. It shares no code with the previous implementation, and every interface below is new rather than changed.

### Added

- `wnc generate-token` prints a Basic auth token, reading the password from a flag, `WNC_PASSWORD`, or piped stdin.
- `wnc show overview|ap|ap-join|ap-tag|client|wlan|policy-tag|site-tag|rf-tag` read a controller over RESTCONF and print a borderless table or a flat JSON array.
- `--pretty` renders the bordered table with a glyph in the state columns.
- `--sort-keys` prints the field names `--sort-by` accepts, so no help line has to enumerate them.
- Several controllers are read concurrently. Each row carries the controller it came from, and a partial failure exits 3 with the rows that were read.
- A JSON configuration file supplies the controller list, the one token every controller is read with, the timeout and the output settings. The read is strict and a mode readable beyond its owner produces a warning.
- `--dry-run` reports what would happen and changes nothing.
- `wnc reset ap --ap-name <ap-name>` restarts one access point. It resolves the target and names it before asking, refuses more than one controller, and takes `--yes` where stdin cannot answer a prompt.
- `wnc reset capwap --ap-name <ap-name>` resets one access point's controller session without restarting the access point. Measured on 17.12.8 and 17.18, it rejoins within about ten seconds and its own uptime keeps climbing.
- `wnc enable|disable ap --ap-name <ap-name>` sets one access point's administrative state. Measured on 17.15.6 and 17.18, the access point stays joined and does not reboot, which is what separates it from `reset ap`.
- `wnc enable|disable radio --ap-name <ap-name> --slot <n>` sets one radio's administrative state. The band number the RPC needs is read from the controller and follows the radio type, so a dual-band radio takes one number whichever band it is serving; band 4, for a dedicated 6 GHz radio, is unverified because the lab holds none.
- `wnc set policy-tag|site-tag|rf-tag --name <name>` creates a tag the controller does not hold and updates one it does. A field no flag names is left as it is rather than cleared.
- `wnc delete policy-tag|site-tag|rf-tag --name <name>` removes a tag, reading it first so a name the controller does not hold is reported plainly.
- `wnc save-config` copies the running configuration to the startup configuration, which is what makes any write survive a reload. It names no target, so it persists everything on the controller and says so before it asks.
- `wnc deauth --mac <mac>` deauthenticates one client, reading it first so an address the controller holds no client at is reported plainly. The operation is absent before 17.15, where a post answers `400` and this command names the release instead of the status.
- `wnc deauth --username <username>` deauthenticates the sessions authenticated under one username, and the prompt says how many the controller holds. One invocation gives one of the two flags.

### Notes

- Every value a subcommand takes is named by a flag, and a leftover positional argument is a usage fault on every command. The word is never repeated back, because a leftover on `generate-token` can be a password. A group reports an unknown command and repeats the word only when it is spelt like a command name.
- One token covers every controller a run reads. Neither `--controller` nor a configuration-file entry carries a credential of its own.
- A value the controller did not report renders as `-` and is absent from the JSON, never as a zero.
- Every failure the controller or the transport answers with is reported in one clause this CLI writes itself, so a timeout, a refused connection, an unverified certificate and an HTTP status all read alike. What those clauses replace is a controller error document that can echo the configuration that was read, and a `*url.Error` naming the whole request URL. `cause=internal` still quotes its error, because a re-worded refusal this CLI writes itself reaches the operator that way.
- Both `reset` leaves report that the instruction was accepted rather than that anything restarted — neither RPC declares an output container, so the status is the whole answer. `wnc show ap-join` carries the outcome and is the only view that keeps a restarting access point listed.
- The three tag lists, an access point's administrative state and the running configuration's persistence are the whole configuration this CLI writes. Everything else about a controller's configuration stays with [telee](https://github.com/umatare5/telee). `deauth` and the `reset` tree write no configuration at all.
- `wnc deauth` reports that the instruction was accepted and not that a client was dropped. Its RPC answers `204` for an identifier associated to nothing exactly as it does for a client it deleted — measured on 17.18.4a, both in under 330ms — so the read that precedes it is the whole of what makes the report truthful. Two of its three arms are implemented: `ip-addr` is left out because every SISF binding measured carries zone 0 and the payload names no other zone.
- `wnc deauth --username` posts once and lets the controller resolve how many sessions that is. The leaf states no cardinality and the lab holds one session per username, so the count in the prompt is the operator's only warning of the blast radius — and the address arm has its own: measured on 17.15.6, two posts each moved the target and one other station on the same BSS.
- The RPC deletes the client's record rather than resetting its association. Measured on 17.15.6 through the username arm: the row was gone 1.9 seconds after the post, all seventeen other clients were untouched, and the station re-associated as a new record 1.5 seconds after the post returned.
- A tag write reaches the running configuration and nowhere else, so a reload loses a create and undoes a delete until `wnc save-config` has run. Measured on all three releases: the controller advertises `writable-running` and neither `:startup` nor `:candidate`, so nothing here can read the startup configuration back or warn that it differs.
- `enable` and `disable` report acceptance rather than a changed state, for the same reason — neither RPC declares an output container. Read it back with `wnc show ap`: measured on 17.15.6 and 17.18, after an access-point-level disable that view reports Admin `Disabled` while `wnc show overview` still reports Admin `Enabled` with Oper `Down`.
- `misconfig_reason` has been seen carrying only the domain's `apmgr-no-misconfig` member: the leaf appears at 17.15, and its display strings come from the model because the device CLI prints no heading for it.
- Verified against IOS-XE `17.12.8`, `17.15.6` and `17.18.4a`. An enum spelling the CLI does not know passes through as the controller sent it.

[v0.1.0]: https://github.com/umatare5/wnc/releases/tag/v0.1.0
