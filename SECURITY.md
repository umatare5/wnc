# Security Policy

## Supported Versions

Only the most recent tagged release carries fixes — reproduce a finding against it before reporting.

## Reporting a Vulnerability

Report privately through [GitHub Security Advisories](https://github.com/umatare5/wnc/security/advisories/new). **Please do not report a vulnerability through a public GitHub issue or a pull request.**

The response is best effort, with no promised window. The advisory goes out once the fix ships, carries a CVE request, and credits the reporter unless they ask otherwise.

## What to Include

**Redact these first.** None of them belongs in a report.

- The access token, from a log line, a command's output or a configuration file
- A device serial number, an access point's or the controller's
- A client-derived hostname, username or IPv6 address

Then include the following:

- **Affected versions** (required): The `wnc` release, and the controller's IOS-XE version
- **Reproduction steps** (required): The command and its flags, and whether `--dry-run` reproduces it
- **Output** (required): The table or the JSON, with every value above removed
- **Impact assessment** (required): The exploit scenario, and what it reaches
- **Suggested fix** (optional): Proposed remediation, if any
- **Disclosure status** (required): Whether it is shared elsewhere, and your plan for sharing it

## Scope

In scope:

- The access token reaching a log line, a help text, or the output of any command but `generate-token`
- The access token reaching the process table other than through `--access-token`, whose cost is documented
- Certificate verification weakened by anything but `--insecure` or the configuration file's `insecure` key
- A request reaching a controller, access point, radio or tag the command did not name
- The published container image

Out of scope:

- A dependency advisory with no path reachable from `./cmd` — show the path, or a `govulncheck` finding
- A controller-side or IOS-XE defect, which belongs to Cisco PSIRT
- An operator's own configuration, which [`docs/configuration.md`](docs/configuration.md) covers
