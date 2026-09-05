<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo_dark.png" width="115px" />
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo.png" width="115px" />
    <img alt="wnc" src="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/logo.png" width="115px" />
  </picture>

  <h1>wnc</h1>

  <p>A command-line interface for Cisco Catalyst 9800 Wireless Network Controllers.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/wnc?label=Latest%20version" />
    <a href="https://github.com/umatare5/wnc/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/wnc/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <a href="https://github.com/umatare5/wnc/actions/workflows/go-vulncheck.yml"><img alt="govulncheck" src="https://github.com/umatare5/wnc/actions/workflows/go-vulncheck.yml/badge.svg?branch=main" /></a><br>
    <img alt="Test Coverage" src="https://raw.githubusercontent.com/umatare5/wnc/main/docs/assets/coverage.svg" />
    <a href="https://www.bestpractices.dev/projects/10820"><img alt="OpenSSF Best Practices" src="https://www.bestpractices.dev/projects/10820/badge" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
    <a href="https://developer.cisco.com/codeexchange/github/repo/umatare5/wnc"><img alt="Published" src="https://static.production.devnetcloud.com/codeexchange/assets/images/devnet-published.svg" /></a>
  </p>

</div>

## Overview

This CLI reads [Catalyst 9800 controllers](https://www.cisco.com/site/us/en/products/networking/wireless/wireless-lan-controllers/catalyst-9800-series/index.html) over RESTCONF and prints their state as a table or JSON.

- 📤 **Shell-friendly**: A borderless table `awk` and `cut` read, and a JSON array keyed by the sort names
- 🌐 **Multi-controller**: Read concurrently, each row labelled with its own, and a partial read prints
- 🔭 **Joined views**: `show overview`, `show ap`, `show client` and `show wlan` join what the device splits
- 🎨 **Pretty output**: `--pretty` borders the table and glyphs the state columns, never the JSON

## Supported Environment

Cisco Catalyst 9800 Wireless Network Controller running on:

- **Cisco IOS-XE 17.12.x** — Verified on 17.12.8 (`deauth` unavailable)
- **Cisco IOS-XE 17.15.x** — Verified on 17.15.6
- **Cisco IOS-XE 17.18.x** — Verified on 17.18.4a

## Quick Start

Please enable RESTCONF and HTTPS on the Catalyst 9800 before using this CLI. Please see:

- [Cisco IOS XE 17.15 Programmability Configuration Guide — RESTCONF](https://www.cisco.com/c/en/us/td/docs/ios-xml/ios/prog/configuration/1715/b_1715_programmability_cg/restconf_protocol.html#id_125840)

### 1. Install the CLI

```bash
docker run --rm ghcr.io/umatare5/wnc:latest --help
```

> [!TIP]
> If you prefer using binaries, download them from the [Release](https://github.com/umatare5/wnc/releases).
>
> **Supported Platform:** `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64`

### 2. Generate a Basic Auth token

Encode your controller account as Base64.

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

> [!TIP]
> `wnc completion <shell>` writes the script to stdout, and `--help` names each shell and the line it needs.

## Syntax

`wnc --help` prints every flag, and [`docs/README.md`](docs/README.md) indexes the reference pages behind it.

The `show` commands read a controller and print:

| Command                                                   | What it does                                                     |
| :-------------------------------------------------------- | :--------------------------------------------------------------- |
| [`wnc show overview`](docs/commands/show-overview.md)     | One row per radio, with the RF settings and the load on it       |
| [`wnc show ap`](docs/commands/show-ap.md)                 | One row per access point                                         |
| [`wnc show ap-join`](docs/commands/show-ap-join.md)       | One row per access point the controller remembers, joined or not |
| [`wnc show ap-tag`](docs/commands/show-ap-tag.md)         | One row per access point, with the tags in force on it           |
| [`wnc show client`](docs/commands/show-client.md)         | One row per associated client                                    |
| [`wnc show wlan`](docs/commands/show-wlan.md)             | One row per WLAN and the policy profile bound to it              |
| [`wnc show policy-tag`](docs/commands/show-policy-tag.md) | One row per WLAN binding a policy tag carries                    |
| [`wnc show site-tag`](docs/commands/show-site-tag.md)     | One row per site tag, with the profiles it names                 |
| [`wnc show rf-tag`](docs/commands/show-rf-tag.md)         | One row per RF tag, with its profile on each band                |

The commands below **act on a controller**, in [the order they all keep](docs/README.md#acting-on-a-controller):

| Command                                                                | What it does                                       |
| :--------------------------------------------------------------------- | :------------------------------------------------- |
| [`wnc reset ap`](docs/commands/reset-ap.md)                            | Restart one access point                           |
| [`wnc reset capwap`](docs/commands/reset-capwap.md)                    | Reset one access point's controller session        |
| [`wnc (enable\|disable) (ap\|radio)`](docs/commands/enable-disable.md) | Set an access point's or one radio's admin state   |
| [`wnc set (policy\|site\|rf)-tag`](docs/commands/set-tag.md)           | Create or update one tag                           |
| [`wnc delete (policy\|site\|rf)-tag`](docs/commands/delete-tag.md)     | Delete one tag                                     |
| [`wnc deauth`](docs/commands/deauth.md)                                | Deauthenticate a client, by address or by username |

Two commands stand outside both groups, each for a reason of its own:

| Command                                                 | What makes it different                                        |
| :------------------------------------------------------ | :------------------------------------------------------------- |
| [`wnc generate-token`](docs/commands/generate-token.md) | Contacts no controller — it encodes an account and prints it   |
| [`wnc save-config`](docs/commands/save-config.md)       | Names no target, so it persists every change on the controller |

## Configuration

Three environment variables reach every command:

| Variable           | Description                                             |
| :----------------- | :------------------------------------------------------ |
| `WNC_CONTROLLER`   | Controller `host[:port]`, comma separated for several   |
| `WNC_ACCESS_TOKEN` | Basic auth token applied to every controller            |
| `WNC_CONFIG`       | Configuration file path, replacing the default location |

`wnc generate-token` reads `WNC_USERNAME` and `WNC_PASSWORD` instead, and contacts no controller. A configuration file keeps the token out of the shell history and the process arguments — see [`docs/configuration.md`](docs/configuration.md).

> [!CAUTION]
> `--insecure` disables TLS verification. **Never use it in production** — [trust the issuer](docs/troubleshooting.md#causetls) instead.

## Troubleshooting

A failed read is one stderr line ending in `(cause=…)`, and a usage fault or a refusal is one `wnc: …` line. [`docs/troubleshooting.md`](docs/troubleshooting.md) indexes both, and `--log-level debug` restores the logfmt form.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the `make` targets, the Docker build and the release process.

## Acknowledgement

I launched this project with the help of **GitHub Copilot Coding Agent**, and I am grateful to the global developer community for their contributions to open source projects and public repositories.

## License

MIT. The binary statically links MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
