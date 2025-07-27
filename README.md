# 📡 wnc - Wireless Network Controller CLI

![GitHub Tag](https://img.shields.io/github/v/tag/umatare5/wnc?label=Latest%20version)
[![Go Reference](https://pkg.go.dev/badge/umatare5/wnc.svg)](https://pkg.go.dev/github.com/umatare5/wnc)
[![Go Report Card](https://goreportcard.com/badge/github.com/umatare5/wnc)](https://goreportcard.com/report/github.com/umatare5/wnc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/umatare5/wnc/blob/main/LICENSE)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10820/badge)](https://www.bestpractices.dev/projects/10820)

A command-line interface tool for managing Cisco Catalyst 9800 Wireless Network Controllers.

- **🚀 Multi-Controller Support**: Efficiently operate several Wireless Network Controllers
- **🧠 Easy Operations**: Focus on key tasks without remembering many complex wireless commands or syntax
- **🔒 Secure API**: All operations use RESTCONF with token-based authentication and TLS encryption
- **🐚 Shell-Friendly Design**: Use shell features—piping, loops, scripting—for advanced automation workflows
- **📊 Clear Output**: Able to show data in table format optimized for easy reading and script processing

<img alt="Demo of wnc show overview" src="https://github.com/umatare5/wnc/blob/main/docs/demo/wnc_show_overview_demo.gif" />

💡 This CLI provides a lightweight and efficient alternative to certain features of [Cisco Catalyst Center](https://www.cisco.com/site/us/en/products/networking/catalyst-center/index.html).

## 📡 Supported Environment

Cisco Catalyst 9800 Wireless Network Controller running Cisco IOS-XE `17.12.x`.

## 📦 Installation

There are two ways to install this CLI:

### 🐳 Docker

```bash
docker run ghcr.io/umatare5/wnc
```

### 📥 Binary Download

Download binaries from the [release page](https://github.com/umatare5/wnc/releases).

**Supported Platforms:** `linux_amd64`, `linux_arm64`, `darwin_amd64` and `darwin_arm64`

## 🚀 Quick Start

### 🔧 Prerequisites

First, activate RESTCONF on the WNC. This CLI requires RESTCONF to communicate with the WNC.

> [!TIP]
>
> [Programmability Configuration Guide, Cisco IOS XE](https://www.cisco.com/c/en/us/td/docs/ios-xml/ios/prog/configuration/1712/b_1712_programmability_cg/m_1712_prog_restconf.html) is a good reference for enabling RESTCONF.

### 🔑 Creating Basic Auth Token

You must create a Basic Auth token using your Cisco WNC credentials before using the CLI.

```bash
# Create token for username:password
echo -n "admin:your-password" | base64
# Output: YWRtaW46eW91ci1wYXNzd29yZA==
```

### 🚀 Basic Usage

Start with this simple example to verify your WNC connection and credentials.

```bash
# Generate authentication token
wnc generate token -u <username> -p <password>

# View controller overview
wnc show overview --controllers "https://wnc1.example.internal:$WNC_ACCESS_TOKEN"
```

### ⚙️ Advanced Configuration

Customize CLI behavior using command-line flags to optimize for your specific environment and requirements.

```bash
# Increase timeout for slow networks
wnc show overview --controllers "https://wnc1.example.internal:$WNC_ACCESS_TOKEN" --timeout 30

# Skip certificate verification (development only)
wnc show overview --controllers "https://wnc1.example.internal:$WNC_ACCESS_TOKEN" --insecure
```

> [!CAUTION]
> The `--insecure` flag disables TLS certificate verification. This should only be used in development environments or when connecting to controllers with self-signed certificates. **Never use this option in production environments** as it compromises security.

## 🌐 CLI Reference

This CLI provides following commands for interacting with Cisco Catalyst 9800 WNC subsystems.

> [!Note]
> Currently, this CLI do not support enough APIs. This will be implemented by the future release `v1.0.0`.

### 🔐 Generate Commands

Generate secure authentication tokens for CLI operations.

| Command              | Description                                      | Documentation                                             |
| -------------------- | ------------------------------------------------ | --------------------------------------------------------- |
| `wnc generate token` | Generate basic auth token from username/password | [📖 GENERATE_TOKEN.md](./docs/commands/GENERATE_TOKEN.md) |

### 📊 Show Commands

Extend and enhance the native `show - summary` commands of C9800 WNC.

| Command             | Description                                       | Documentation                                           |
| ------------------- | ------------------------------------------------- | ------------------------------------------------------- |
| `wnc show overview` | Display the summary of 2.4 GHz, 5GHz and 6GHz.    | [📖 SHOW_OVERVIEW.md](./docs/commands/SHOW_OVERVIEW.md) |
| `wnc show ap`       | Display the summary of associated APs.            | [📖 SHOW_AP.md](./docs/commands/SHOW_AP.md)             |
| `wnc show ap-tag`   | Display the summary of tag names with the status. | [📖 SHOW_AP_TAG.md](./docs/commands/SHOW_AP_TAG.md)     |
| `wnc show client`   | Display the summary of associated clients.        | [📖 SHOW_CLIENT.md](./docs/commands/SHOW_CLIENT.md)     |
| `wnc show wlan`     | Display the summary of configured WLANs.          | [📖 SHOW_WLAN.md](./docs/commands/SHOW_WLAN.md)         |

### ⚡ Exec Commands

Please use [telee](https://github.com/umatare5/telee) as an alternative for executing commands on the WNC.

> [!Note]
>
> As of June 2025, the `exec` command is not yet implemented.

#### 🔧 Use Case Examples

**Reset an AP:**

```bash
telee -H wnc1 -C "ap name <apName> reset"
```

**Deauthenticate a client:**

```bash
# By MAC address
telee -H wnc1 -C "wireless client mac-address <clientMac> deauthenticate"

# By IP address
telee -H wnc1 -C "wireless client ip-address <clientIpAddr> deauthenticate"

# By username
telee -H wnc1 -C "wireless client username <userName> deauthenticate"
```

## 🧪 Testing

This repository includes comprehensive unit and integration tests to ensure reliability and compatibility with Cisco Catalyst 9800 controllers. For detailed testing information, please see **[TESTING.md](./docs/TESTING.md)**.

## 🛠️ Troubleshooting

If you encounter issues, please see the **[TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md)** for common problems and solutions.

## 🤝 Contributing

I welcome contributions to improve this CLI. Please follow these guidelines to ensure smooth collaboration.

1. **Fork the repository** and create a feature branch from `main`
2. **Make your changes** following existing code style and conventions
3. **Add comprehensive tests** for new functionality
4. **Update documentation** including README.md and code comments
5. **Ensure all tests pass** including unit and integration tests
6. **Submit a pull request** with a clear description of changes

## 🚀 Release

To release a new version:

1. **Update the version** in the `VERSION` file
2. **Submit a pull request** with the updated `VERSION` file

Once merged, the GitHub Workflow will automatically:

- **Create and push a new tag** based on the `VERSION` file

After that, manual release using [GitHub Actions: release workflow](https://github.com/umatare5/wnc/actions/workflows/release.yaml).

## 🙏 Acknowledgments

This code was developed with the assistance of **GitHub Copilot Agent Mode**. I extend our heartfelt gratitude to the global developer community who have contributed their knowledge, code, and expertise to open source projects and public repositories.

## 📄 License

Please see the [LICENSE](./LICENSE) file for details.
