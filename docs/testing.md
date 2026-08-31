# Testing

How the suite is arranged, what it asserts, and the identities every fixture and sample uses.

```bash
make test-unit            # go test -race with coverage
make test-unit-coverage   # plus an HTML report under ./coverage
make lint                 # config verify, golangci-lint and go mod tidy
```

`gotestsum` and `golangci-lint` are required — `make help` names the install commands.

## How the suite is arranged

Tests sit next to the code they cover, in the same package, so an unexported rule can be asserted directly rather than through a public surface built for the test. Tables and `t.Parallel()` are the default. There is no assertion library and no golden file.

## Fixture Identities

A fixture takes its shape from a real controller response and none of its identities. Every value below is synthetic, and this section is the canon for the samples under [`commands/`](commands/) as well as for the test files.

Four kinds have a range reserved for exactly this, so nothing here is invented:

| Kind         | Range                                     | Reserved by           |
| :----------- | :---------------------------------------- | :-------------------- |
| MAC address  | `00:00:5e:00:53:00` – `00:00:5e:00:53:ff` | RFC 7042 §2.1.2       |
| IPv4 address | `192.168.0.0/16`                          | RFC 1918              |
| IPv6 address | `2001:db8::/32`                           | RFC 3849              |
| Domain name  | `example.internal`                        | ICANN private-use TLD |

The MAC block's first octet has the I/G bit clear and the U/L bit set, so no fixture address can be multicast or collide with a vendor assignment. Its last octet carries the role:

- **`:01`–`:0f`** — an access point's radio base address, `TEST-APnn` pairing with `:nn`
- **`:11`–`:1f`** — an access point's Ethernet address, `TEST-APnn` pairing with `:1n`
- **`:a1`–`:af`** — a client station

The rest have no standard to take them from, so they are this repository's own:

| Kind                 | Value                          |
| :------------------- | :----------------------------- |
| Controller name      | `WNC1` – `WNC4`                |
| Controller host      | `192.168.0.1` – `192.168.0.4`  |
| Access point name    | `TEST-AP01` – `TEST-AP99`      |
| Access point serial  | `TST0000AP01` – `TST0000AP99`  |
| Access point address | `192.168.0.11` onward          |
| Client address       | `192.168.0.21` onward          |
| Access token         | `TestToken0123456789ABCDEF==`  |
| Password             | `test-token-123`               |
| SSID                 | `test-essid01` onward          |
| Profile              | `test-<kind>-profile01` onward |
| Tag                  | a `test-` prefix               |

A profile's `<kind>` is `wlan`, `policy`, `rf`, `ap` or `flex`, so the kind a sample names is readable from the value alone.

A controller's name is what a configuration file supplies and what every prompt, report and `Controller` column carries, so a sample transcript names `WNC1` where the run behind it named `192.168.0.1` on `--controller`. [`examples/config.json`](../examples/config.json) is that file, and pairs the two.

A serial cannot take the `test-` prefix and keep its shape, so it reads `TST0000APnn` — a valid `[A-Z]{3}[0-9]{4}[A-Z0-9]{4}` that no site code and no week 00 of year 00 can collide with.

Two categories are deliberately outside the scheme.

- **A grammar case is not an identity** — `internal/config` keeps `a.example` and `[2001:db8::1]` for parse assertions
- **A dialled address must be unroutable** — RFC 1918 space can be live on a developer's own LAN

So a test whose host may be reached keeps an address reserved against reachability instead: `192.0.2.0/24` from RFC 5737, and `240.0.0.1` from RFC 5735 where the assertion is that nothing answers.

> [!IMPORTANT]
> Never paste a MAC address, serial number, hostname, username, SSID or tag name from a capture into a committed fixture or a sample transcript. Nothing in this repository redacts one, so a value pasted by hand reaches the tree unchanged.

## The RESTCONF layer

`internal/wnc` is tested against a TLS test server serving canned responses, routed on the last element of the request path. The SDK pins its own dialer, so no transport can be injected and the server has to be a real listener — one test asserts that much before the fixture-driven ones rely on it.

One fixture deliberately includes a credential leaf, to assert the hand-written struct drops it at decode.

## The fan-out

`internal/show` tests the fan-out with a fetch function that never reaches the network. Client construction still happens for real, against the unroutable range above, so the outcome classification, the reporting order and the "print nothing when everything failed" rule are exercised without a server.

## The command tree

`internal/cli` drives the real command tree end to end and asserts the exit-code contract: usage faults, an unknown command, the help paths, the settings rejections and the three ways `generate-token` takes a password.

Nothing in that file runs in parallel, and that is deliberate: urfave reads the `WNC_*` variables at parse time, so the suite clears them with `t.Setenv`, which forbids `t.Parallel`. `make test-unit` clears them again at the process level so a developer's own shell cannot change what the assertions see.

## Invariants

Some checks are about shape rather than behaviour, and they exist because the failure they catch is silent.

- **Three declarations, one order** — the sort-key list, the column list and the json tags must agree
- **Banned outright** — `omitempty` drops a reported zero, an empty string and a reported false alike
- **Allowed on a pointer only** — `omitzero` there means nil, which is what "not reported" is
- **Banned as well** — every value of the json `format` tag is rejected at run time, not at compile time
- **Every command carries the usage hook** — urfave consults the running command's own and no other

> [!NOTE]
> `json/v2` drops a tag a sibling field repeats, which would leave a column in the table and absent from the JSON with nothing failing. That silence is why the three declarations are asserted rather than reviewed.

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

`--dry-run` stops before the RPC, so it verifies everything except the write. The write itself was measured once on 17.18.4a for all four access-point RPCs, three of them through the ap-name arm. A dry run is not a substitute for repeating that on a release where it matters. `save-config` was measured on all three releases, `deauth --mac` on 17.18.4a and `deauth --username` on 17.15.6.

**A write measurement needs a stable target, not just a before and an after.** The `deauth --username` post was attributed to its effect because the target's association had been unchanged for 82 minutes across four snapshots, while two no-post control windows of 35 seconds each moved 1 and 0 of the 18 clients. Without the stability, one moved client is inside the estate's own churn.

**The 400 on 17.12.8 has no CLI-level measurement.** Every client on that controller carries an empty username, so `--username` is refused by the resolve at exit 1 and the RPC is never reached. Both arms' classification is pinned in `internal/wnc/deauth_test.go` and the re-wording in `internal/cli/deauth_test.go`.

An administrative state has no arbiter this CLI can compare against. Measured on 17.15.6, a `show running-config all` filtered on the access point's name returns nothing, and per-AP configuration is keyed by dotted MAC rather than by name, so that filter settles nothing either way. Read the state back with `wnc show ap` and `wnc show overview` instead — after an access-point-level disable the two disagree by design, which [`enable-disable.md`](./commands/enable-disable.md) explains.
