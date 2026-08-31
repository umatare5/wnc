# Testing

```bash
make test-unit            # go test -race with coverage
make test-unit-coverage   # plus an HTML report under ./coverage
make lint                 # config verify, golangci-lint and go mod tidy
```

`gotestsum` and `golangci-lint` are required — `make help` names the install commands.

## How the suite is arranged

Tests sit next to the code they cover, in the same package, so an unexported rule can be asserted directly rather than through a public surface built for the test. Tables and `t.Parallel()` are the default. There is no assertion library and no golden file.

## The RESTCONF layer

`internal/wnc` is tested against a TLS test server serving canned responses, routed on the last element of the request path. The SDK pins its own dialer, so no transport can be injected and the server has to be a real listener — one test asserts that much before the fixture-driven ones rely on it.

Fixtures carry no value read off a device. Addresses come from the documentation range RFC 7042 reserves, and one fixture deliberately includes a credential leaf to assert the hand-written struct drops it at decode.

## The fan-out

`internal/show` tests the fan-out with a fetch function that never reaches the network. Client construction still happens for real, against hosts from the RFC 5737 test range, so the outcome classification, the reporting order and the "print nothing when everything failed" rule are exercised without a server.

## The command tree

`internal/cli` drives the real command tree end to end and asserts the exit-code contract: usage faults, an unknown command, the help paths, the settings rejections and the three ways `generate-token` takes a password.

Nothing in that file runs in parallel, and that is deliberate: urfave reads the `WNC_*` variables at parse time, so the suite clears them with `t.Setenv`, which forbids `t.Parallel`. `make test-unit` clears them again at the process level so a developer's own shell cannot change what the assertions see.

## Invariants

Some checks are about shape rather than behaviour, and they exist because the failure they catch is silent:

- A column is declared in three places — the sort-key list, the column list and the row struct's json tags — and all three must agree, in order. json/v2 drops a tag a sibling field repeats, which would otherwise leave a column in the table and absent from the JSON with nothing failing.
- `omitempty` is banned outright, because it also drops a reported zero, an empty string and a reported false. `omitzero` is allowed only on a pointer, where the zero value is nil and so genuinely means "not reported".
- The json `format` tag is banned: every value of it is rejected at run time, after passing both the compiler and the linter.
- Every command in the tree must carry the usage hook, because urfave consults only the running command's own.

## Coverage

CI enforces a floor. Run `make test-unit-coverage` and open `coverage/report.html` to see what a change left uncovered.

## Against a real controller

There is no integration-test target. A `show` command is verified by running it and comparing the result with the controller's own output:

```bash
wnc show overview -c "<host>" --insecure
```

`show ap dot11 5ghz summary`, `show ap uptime`, `show ap tag summary`, `show wlan id <n>` and `show wireless client summary` each cover one view. The three tag views compare against `show wireless tag {rf,site,policy} summary` and `detailed <name>`, and three of their headings deliberately differ from the device's: it labels the per-band RF profile `2.4ghz RF Policy`, the AP join profile `AP Profile` and the bound policy profile `Policy Name`, where these views follow the YANG leaf and the write flag instead.

The seven trees that act cannot be verified that way, because running one changes the controller. `--dry-run` exercises everything up to the write, and a `show` command reads the result back afterwards:

```bash
wnc --dry-run disable radio --ap-name "<ap-name>" --slot 1 -c "<host>"
```

`--dry-run` stops before the RPC, so it verifies everything except the write. The write itself was measured once on 17.18 for all four access-point RPCs, three of them through the ap-name arm. A dry run is not a substitute for repeating that on a release where it matters. `save-config` was measured on all three releases, `deauth --mac` on 17.18.4a and `deauth --username` on 17.15.6.

**A write measurement needs a stable target, not just a before and an after.** The `deauth --username` post was attributed to its effect because the target's association had been unchanged for 82 minutes across four snapshots, while two no-post control windows of 35 seconds each moved 1 and 0 of the 18 clients. Without the stability, one moved client is inside the estate's own churn.

**The 400 on 17.12.8 has no CLI-level measurement.** Every client on that controller carries an empty username, so `--username` is refused by the resolve at exit 1 and the RPC is never reached. Both arms' classification is pinned in `internal/wnc/deauth_test.go` and the re-wording in `internal/cli/deauth_test.go`.

An administrative state has no arbiter this CLI can compare against. Measured on 17.15.6, a `show running-config all` filtered on the access point's name returns nothing, and per-AP configuration is keyed by dotted MAC rather than by name, so that filter settles nothing either way. Read the state back with `wnc show ap` and `wnc show overview` instead — after an access-point-level disable the two disagree by design, which [`enable-disable.md`](./commands/enable-disable.md) explains.
