# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for the project overview, and [`docs/README.md`](docs/README.md) for the reference pages behind it.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`umatare5/cisco-ios-xe-wireless-go`](https://github.com/umatare5/cisco-ios-xe-wireless-go) v0.11+ — sole RESTCONF SDK for Cisco C9800 WNC
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ — command tree, flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) v1.10+ — process logger and the `slog` bridge the SDK takes
- [`olekukonko/tablewriter`](https://github.com/olekukonko/tablewriter) v1.1+ — borderless table writer
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); exits on what `internal/cli` returns
- `internal/cli/` — urfave command tree, flag definitions, exit codes, and the `version` string ldflags sets
- `internal/config/` — Flag, environment and file resolution in that precedence, and the three grammars
- `internal/log/` — logrus setup and the `*slog.Logger` the SDK takes
- `internal/wnc/` — Sole importer of the SDK; one `fetch_*.go` per read, and one file per action beside them
- `internal/show/` — Per-command row building, the concurrent controller fan-out, and the enum display tables
- `internal/render/` — `Column[T]`, shared by the table and the JSON writer so a column cannot exist in one only

## Setup and Commands

Install required tools (one-time):

- `go install gotest.tools/gotestsum@latest`
- `golangci-lint` - See <https://golangci-lint.run/docs/welcome/install/>
- `goreleaser` release builds (see [`.goreleaser.yml`](.goreleaser.yml))
- `make pre-commit-install` wires every hook in [`.pre-commit-config.yaml`](.pre-commit-config.yaml)

Make targets ([`Makefile`](Makefile)):

- `make build` — Build binary into `tmp/wnc`
- `make lint` — `golangci-lint config verify` + `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make snapshot` — Build a `goreleaser` snapshot
- `make clean` — Remove build artifacts and `.bak*` files
- `make image` — Build Docker image (`$USER/wnc`)
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the hooks

## Code Style

- Linting and formatting are `golangci-lint`'s, and its `formatters` block owns every formatter it runs
- A comment carries what the code cannot — a value another file must match, an order the device rejects
- One or two sentences, English, no emoji, and nothing a reader can derive from the code beside it

## Testing

- Run `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`) and in the same package.
- Coverage threshold is enforced by [`.github/workflows/go-test-coverage.yml`](.github/workflows/go-test-coverage.yml).

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, tests, CodeQL, govulncheck, actionlint, markdownlint and the link check.
- Take every fixture and sample identity from [`docs/testing.md`](docs/testing.md#fixture-identities), never a value read off a device.

## Domain Knowledge

### Verifying Values

- **A YANG model is a design document, not the implementation.** Units, ranges, enum spellings, and even the presence of a leaf can differ on a live controller, so confirm every value against a RESTCONF response from a real WNC, and take an enum's members and their numbers from the model the controller itself serves.
- **A configuration leaf missing from a response means its default is in force, not that nothing set it.** The default is often `true`, so decoding an omitted boolean as `false` inverts the reading: on a plain read of `wlan-cfg-entries`, `wpa2-enabled` and `auth-key-mgmt-dot1x` are absent from exactly the WLANs where they are enabled, and present only where they were explicitly switched off.
- **Ask for the values in force, and only on a configuration read — an operational absence is structural, so materialising a default there fabricates a reading.** `?with-defaults=report-all` returns the omitted leaves and a controller that rejects it answers `400`. It also materialises the PSK and the WEP key material, so the raw struct in [`internal/wnc/fetch_wlan.go`](internal/wnc/fetch_wlan.go) is an allow-list whose undeclared leaves are dropped at decode.
- **A present container can still be missing the leaf, and a present leaf must not be dropped.** Most SDK containers are non-pointer structs, so an absent leaf decodes to `0`: `curr-freq` sits under `when radio-mode != monitor/sniffer/invalid` while `phy-ht-cfg` around it stays present. The SDK types `curr-freq` and `chan-width` as pointers for exactly this reason, and a non-pointer sibling still needs a per-leaf zero guard; `chan-width` and `curr-tx-power-in-dbm` do arrive on those radios, so suppressing them to match the device CLI's `N/A` is fabrication inverted.
- **Arbitrate on the device with `show running-config all`, and only for configuration.** It prints the negated form for a feature that is off, so a WLAN with no such line has it on. A display string follows the controller's own wording instead — `show ap ... config general` prints `AP Mode : FlexConnect`, so `mode-flex-connect` renders as `FlexConnect`.

### RESTCONF Access Patterns

GET an operational collection:

```bash
curl -k -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        -H "Accept: application/yang-data+json" \
        "https://$WNC_CONTROLLER/restconf/data/Cisco-IOS-XE-wireless-access-point-oper:access-point-oper-data/capwap-data"
```

GET a configuration collection with the defaults in force:

```bash
curl -k -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        -H "Accept: application/yang-data+json" \
        "https://$WNC_CONTROLLER/restconf/data/Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-data/wlan-cfg-entries?with-defaults=report-all"
```

GET the model the controller implements, taking the revision from `ietf-yang-library:modules-state`:

```bash
curl -k -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        "https://$WNC_CONTROLLER/restconf/tailf/modules/Cisco-IOS-XE-wireless-access-point-oper/2023-08-01"
```

List the operations the controller publishes, which is a plain GET — only each RPC beneath it is POST-only:

```bash
curl -k -H "Authorization: Basic $WNC_ACCESS_TOKEN" \
        -H "Accept: application/yang-data+json" \
        "https://$WNC_CONTROLLER/restconf/operations"
```

### Admitting an RPC

The `reset` tree acts without persisting anything. **An RPC belongs there only if all three hold**, tested per RPC and never inherited from its module or from the `/restconf/operations` root.

- The controller declares no configuration-datastore twin for it
- `show running-config all` prints no line for it
- Its schema declares no `output` container

The same `-cmd-rpc` module as `ap-reset` also declares `set-ap-static-ip-enable` and `set-ap-reset-button`, both of which persist, and the operations root publishes `Cisco-IOS-XE-cli-rpc:config-ios-cli-rpc`, which enters configuration mode outright.

**`cisco-ia:save-config` is in the tree and did not pass this test.** It fails part three — its schema declares an `output` container — and it exists to persist, which the first line of this test forbids. It is a top-level command of its own rather than a `reset` leaf for exactly that reason, and it reaches a controller on the owner's decision below.

**`apf-ms-delete-all` passed the test and is still not a `reset` leaf.** It declares no `output`, the client domain publishes no `-cfg` module to hold a twin, and nothing persists — the client re-associates on its own. What keeps it out of `reset` is that tree's own two properties: it restarts something, and every leaf of it names one access point by a key the controller holds. A client has no such key: its address is one, and the username `--username` takes is a bare string that may select several sessions. So `deauth` is flat, and a pass here licenses a command rather than a placement.

Part two of that test is weak on its own and must not be used alone: per-AP configuration appears in `show running-config all` as `ap <dotted-mac>` blocks rather than as `ap name <x>`, and neither `set-rad-capwap-reset` nor the rejected `set-ap-reset-button` prints any line. What separates them is that a persisting RPC takes a state leaf readable back as AP state — `reset-button-state` at `capwap-data/*/device-detail/dynamic-info` — and a transient one takes only a target.

**`enable` and `disable` are the counter-example, and neither is a precedent.** Measured on 17.15.6, `set-ap-admin-state` and `set-ap-slot-admin-state` clear all three parts — `ap-cfg` carries no admin-state leaf, a `show running-config all` filtered on the access point's name returns nothing, and neither RPC declares an `output` — and the discriminator above still refuses them, because a disable is readable back at `ap-state/ap-admin-state` as `adminstate-disabled`. Part two is as weak here as it is above: a name filter cannot see the `ap <dotted-mac>` blocks either way. They reach a controller because the repository owner approved a second configuration surface, not because a test admitted them. So nothing enters `reset` on the three parts alone, and nothing enters `enable`, `disable`, `set` or `delete` on the strength of that approval.

### Writing configuration

The `set` and `delete` trees write the three tag lists, `enable` and `disable` write an admin state, and `save-config` persists the running configuration. Those are the whole of the configuration this CLI performs, and everything else belongs to [telee](https://github.com/umatare5/telee). Each was admitted by the repository owner, and admitting a fourth is that decision again rather than a consequence of any of them.

`deauth` reaches a controller on the same kind of decision and is not a fourth configuration surface: it configures nothing, and the client it drops comes back on its own.

`save-config` differs from the other two in naming no target, so it persists whatever else is in the running configuration. Nothing in this CLI can see that: the controller advertises no `:startup` capability on any release in scope, so the startup configuration is unreachable over RESTCONF and the prompt is the whole of the warning an operator gets.

Three facts measured on every release in scope govern a tag write.

- The key leaves declare a `pattern` and no `length`, yet the controller refuses a 33-character name on every kind — measured per kind on 17.12.8, so the 32-character cap is the device's and unreadable from the model
- The three configuration modules declare no `leafref` and no `require-instance`, so a dangling profile reference is accepted and persists
- All three tag services read the record and write it back with a merge PATCH, so a field the command did not name survives

An admin-state write rests on facts of its own, and the one that bites is the band number: it follows what occupies the slot rather than the band that radio is serving, so a dual-band radio takes 3 whichever band it is on and the served band's own number answers 400. That is why the CLI reads the number off `radio-type` and shows the operator `current-active-band`.
