<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo_dark.png" width="115px" />
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo.png" width="115px" />
    <img alt="wnc" src="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo.png" width="115px" />
  </picture>

  <h1>wnc</h1>

  <p>A command-line interface for Cisco C9800 Wireless Network Controller.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/wnc?label=Latest%20version" />
    <a href="https://github.com/umatare5/wnc/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/wnc/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <img alt="Test Coverage" src="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/coverage.svg" />
    <a href="https://goreportcard.com/report/github.com/umatare5/wnc"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/umatare5/wnc" /></a><br>
    <a href="https://www.bestpractices.dev/projects/10820"><img alt="OpenSSF Best Practices" src="https://www.bestpractices.dev/projects/10820/badge" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
    <a href="https://developer.cisco.com/codeexchange/github/repo/umatare5/wnc"><img alt="Published" src="https://static.production.devnetcloud.com/codeexchange/assets/images/devnet-published.svg" /></a>
  </p>

</div>

## Overview

This CLI reads [Cisco Catalyst 9800 Wireless Controllers](https://www.cisco.com/site/us/en/products/networking/wireless/wireless-lan-controllers/catalyst-9800-series/index.html) over RESTCONF and prints their state as a table or as JSON.

- 🌐 **Multi-controller**: One invocation reads several controllers concurrently and labels every row with its own
- 🔍 **Absence is visible**: A value the controller did not report shows as `-`, never as a zero
- 📤 **Shell-friendly**: A borderless table for reading, and a flat JSON array whose fields are the sort keys
- 🔒 **RESTCONF only**: Token authentication over TLS, and every command that acts asks before it does

💡 This CLI is a lightweight alternative to parts of [Cisco Catalyst Center](https://www.cisco.com/site/us/en/products/networking/catalyst-center/index.html).

## Supported Environment

Verified against Cisco Catalyst 9800 controllers running IOS-XE `17.12.8`, `17.15.6` and `17.18.4a`.

Older and newer releases are expected to work. An enum spelling this CLI does not know is passed through as the controller sent it rather than being blanked, so a release that adds a value stays readable.

## Quick Start

Please enable RESTCONF and HTTPS on the C9800 before using this CLI. Please see:

- [Cisco IOS XE 17.15 Programmability Configuration Guide — RESTCONF](https://www.cisco.com/c/en/us/td/docs/ios-xml/ios/prog/configuration/1715/b_1715_programmability_cg/restconf_protocol.html#id_125840)

### 1. Install the CLI

```bash
docker run --rm ghcr.io/umatare5/wnc:latest --help
```

> [!TIP]
> If you prefer using binaries, download them from the [release page](https://github.com/umatare5/wnc/releases).
>
> Supported Platforms are: `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64`

### 2. Generate a Basic Auth token

Encode your controller credentials as Base64.

```bash
read -rs WNC_PASSWORD && export WNC_PASSWORD
export WNC_ACCESS_TOKEN="$(wnc generate-token -u admin)"
```

### 3. Set required environment variables

```bash
export WNC_CONTROLLER="wnc1.example.internal"
```

`--controller` is repeatable, so several controllers need no separator at all:

```bash
wnc show overview -c wnc1.example.internal -c wnc2.example.internal
```

### 4. Read the controller

```bash
wnc show overview
```

### 5. Enable shell completion

`wnc completion` writes a script for `bash`, `zsh`, `fish` or `pwsh` to stdout.

```bash
source <(wnc completion bash)
```

Fish reads a file rather than a sourced script, so write it to `~/.config/fish/completions/wnc.fish`. `wnc completion --help` prints the line for each shell.

## Syntax

`wnc --help` prints every flag, and [`docs/CLI_REFERENCE.md`](docs/CLI_REFERENCE.md) carries the same list.

Every `show` command accepts `--controller`, `--access-token`, `--insecure`, `--format`, `--pretty`, `--timeout`, `--sort-by`, `--sort-keys` and `--sort-order`:

| Command                                                             | What it does                                                     |
| :------------------------------------------------------------------ | :--------------------------------------------------------------- |
| [`wnc generate-token`](docs/commands/generate-token.md)             | Print the Basic auth token every other command needs             |
| [`wnc show overview`](docs/commands/show-overview.md)               | One row per radio, with the RF settings and the load on it       |
| [`wnc show ap`](docs/commands/show-ap.md)                           | One row per access point                                         |
| [`wnc show ap-join`](docs/commands/show-ap-join.md)                 | One row per access point the controller remembers, joined or not |
| [`wnc show ap-tag`](docs/commands/show-ap-tag.md)                   | One row per access point, with the tags in force on it           |
| [`wnc show client`](docs/commands/show-client.md)                   | One row per associated client                                    |
| [`wnc show wlan`](docs/commands/show-wlan.md)                       | One row per WLAN and the policy profile bound to it              |
| [`wnc show policy-tag`](docs/commands/show-policy-tag.md)           | One row per WLAN binding a policy tag carries                    |
| [`wnc show site-tag`](docs/commands/show-site-tag.md)               | One row per site tag, with the profiles it names                 |
| [`wnc show rf-tag`](docs/commands/show-rf-tag.md)                   | One row per RF tag, with its profile on each band                |
| [`wnc reset ap`](docs/commands/reset-ap.md)                         | Restart one access point                                         |
| [`wnc reset capwap`](docs/commands/reset-capwap.md)                 | Reset one access point's controller session                      |
| [`wnc enable …` / `wnc disable …`](docs/commands/enable-disable.md) | Set an access point's or one radio's admin state                 |
| [`wnc set …-tag`](docs/commands/set-tag.md)                         | Create or update one policy, site or RF tag                      |
| [`wnc delete …-tag`](docs/commands/delete-tag.md)                   | Delete one policy, site or RF tag                                |
| [`wnc save-config`](docs/commands/save-config.md)                   | Copy the running configuration to the startup configuration      |
| [`wnc deauth`](docs/commands/deauth.md)                             | Deauthenticate a client, by address or by username               |

The `reset`, `enable`, `disable`, `set`, `delete` and `deauth` trees act on a controller. Each names one target, resolves it against the controller before acting, and asks first unless `--yes` is given. `--dry-run` reports and changes nothing.

`wnc save-config` acts on a controller too, and differs in naming no target: it persists the whole running configuration, so a write of any kind survives a reload only after it runs.

Configuration is limited to the three tag kinds and an access point's administrative state. Use [telee](https://github.com/umatare5/telee) for anything else — `wnc deauth` and the `reset` tree are operational actions and configure nothing.

> [!CAUTION]
> `--insecure` disables TLS certificate verification. **Never use it in production.** [`docs/SECURITY.md`](docs/SECURITY.md) covers trusting a private certificate authority instead.

## Configuration

This CLI reads five environment variables:

| Environment Variable | Description                                             |
| :------------------- | :------------------------------------------------------ |
| `WNC_CONTROLLER`     | Controller `host[:port]`, comma separated for several   |
| `WNC_ACCESS_TOKEN`   | Basic auth token applied to every controller            |
| `WNC_CONFIG`         | Configuration file path, replacing the default location |
| `WNC_USERNAME`       | Controller username, read by `wnc generate-token` only  |
| `WNC_PASSWORD`       | Controller password, read by `wnc generate-token` only  |

### Configuration File

A configuration file avoids repeating the controller list, and keeps the token out of the shell history and out of the process arguments. The path is `--config`, then `$WNC_CONFIG`, then `$XDG_CONFIG_HOME/wnc/config.json`, then `~/.config/wnc/config.json`.

```json
{
  "note": "lab controllers",
  "timeout": "30s",
  "insecure": false,
  "format": "table",
  "pretty": false,
  "log_level": "warning",
  "token": "...",
  "controllers": [
    {
      "name": "lab-17.12",
      "host": "192.168.0.231",
      "note": "retiring in 2026-09"
    }
  ]
}
```

The token is file-wide rather than per entry: one run reads every controller with one credential, so the file holds one secret however many hosts it lists. The read is strict: an unknown key, a duplicated key, a key differing only in case, a comment and a trailing comma are all rejected with the JSON pointer that located the fault. JSON has no comments, so the `note` fields exist to carry an operational remark instead.

Check a hand-edited file without contacting anything:

```bash
wnc --config ./config.json --dry-run
```

> [!IMPORTANT]
> Prefer the configuration file at mode `0600` over an environment variable, and an environment variable over `--access-token`: a token on the command line is visible to every process on the host. A file readable beyond its owner produces a warning.

A host named by `--controller` with no token takes the file's own `token`, so the file stays usable for a one-off read against a host it does not list. The entry naming that host is not consulted — there is nothing per entry to consult:

```bash
wnc show ap -c 192.168.0.231
```

## Output

The table is borderless so `awk` and `cut` work on it, and the first column starts at column zero. A cell the controller did not report is `-`.

The JSON is a flat array with one object per table row. Its field names are exactly the values `--sort-by` accepts, and `--sort-keys` prints that list for any command without contacting a controller. Numbers stay numbers, and a field the controller did not report is absent rather than null. An empty result is `[]`.

`--pretty` rules the table and puts a glyph in the state columns. It styles the table only: the JSON is unaffected, and so is anything piped from the default output.

```bash
wnc show client -f json | jq -r '.[] | select(.rssi_dbm < -70) | .mac'
```

## Exit Codes

| Code | Meaning                                                        |
| :--- | :------------------------------------------------------------- |
| 0    | Every controller answered. An empty fleet is a success         |
| 1    | No controller answered, or the CLI failed internally           |
| 2    | Usage or configuration fault. Nothing was sent to a controller |
| 3    | Partial: at least one read failed and at least one succeeded   |
| 130  | Interrupted. No partial table is printed                       |

## Troubleshooting

A read that fails is one line on stderr ending in `(cause=…)`, and [`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md) is indexed by that token. A usage fault and a refusal from a command that acts are one `wnc: …` line, which that page indexes by its wording instead. `--log-level debug` restores the logfmt form, which carries the controller and the HTTP status as fields of their own.

## Contributing

See [`CONTRIBUTING.md`](https://github.com/umatare5/wnc/blob/main/CONTRIBUTING.md) for the `make` targets, the Docker build, the release process and how to open a pull request.

## Acknowledgement

I launched this project with the help of **GitHub Copilot Coding Agent**, and I am grateful to the global developer community for their contributions to open source projects and public repositories.

## Licence

[MIT](LICENSE).
