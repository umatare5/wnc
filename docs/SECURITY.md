# Security

## Reporting

Open a private security advisory on the repository. Do not file a public issue for a vulnerability.

## What this CLI sends and stores

Every request goes to the controller's RESTCONF interface over TLS, and nothing is written to disk — the configuration file is only ever read. Seven of the nine command trees write to a controller, and the table below is the whole of what any of them sends.

| Tree                | What reaches the controller                                 |
| ------------------- | ----------------------------------------------------------- |
| `show`              | `GET` only                                                  |
| `reset`             | `GET`, then `POST` to `/restconf/operations/`               |
| `enable`, `disable` | `GET`, then `POST` to `/restconf/operations/`               |
| `set`               | `GET`, then `POST` to create, or `PATCH` to update          |
| `delete`            | `GET`, then `DELETE` on the keyed record                    |
| `save-config`       | `POST` to `/restconf/operations/`, with no read before it   |
| `deauth`            | `GET`, then `POST` to `/restconf/operations/`               |
| `generate-token`    | Nothing. It contacts no controller                          |

The `GET` in each acting row is the read that names the target, so a run that is then refused has still read the controller. `save-config` is the exception and has no such read, because its RPC takes no target. What those writes can change is the three tag lists, an access point's administrative state, whether the running configuration is persisted, and whether one client stays associated, and nothing else: every one asks before it acts unless `--yes` is given, and `--dry-run` stops before the write.

The access point's name is keyed straight into that read, so whatever is typed there reaches the controller as a path segment of the URL — measured against a local TLS stub. Nothing bounds it locally, because nothing measured bounds the name a controller accepts, so a secret pasted into that argument has left the host before anything reports it. `deauth` keys its client address the same way, and that one is validated as an address first, so a value that is not one is refused before any request goes out.

`deauth --username` is the one acting argument that reaches no URL. Its read is a filter applied to the collection locally, and the value leaves the host only in the RPC's request body — so nothing bounds it either, and a secret pasted there is disclosed just as fully, only later and somewhere a proxy log would not show it.

The token is sent as an `Authorization: Basic` header and appears nowhere else. No diagnostic quotes it — an error about a malformed `--controller` entry names the entry by its index and states the expected syntax instead of echoing what was typed, and a failed configuration decode reports the JSON pointer and the Go type rather than the value. A value a flag cannot parse is reported by the flag's name alone.

## Where to keep the token

In order of preference:

1. **A configuration file at mode `0600`.** The token never reaches the shell history or the process arguments. A file readable beyond its owner produces a warning naming the mode.
2. **`$WNC_ACCESS_TOKEN`.** Visible to anything that can read the process environment.
3. **`--access-token`.** Visible to every process on the host through the process list. Use it interactively only.

`--controller` carries no token, so naming a host on the command line costs nothing: the file's one token covers every host, whether or not the file lists it.

The same order applies to `generate-token`: prefer a piped password to `$WNC_PASSWORD`, and either to `-p`.

## Trusting a private certificate authority

A controller usually presents a self-signed or internally-issued certificate. `--insecure` turns verification off entirely, which also accepts an interception, so prefer trusting the issuer:

```bash
SSL_CERT_FILE=/path/to/ca-bundle.pem wnc show overview
```

> [!IMPORTANT]
> On macOS and Windows, Go 1.27 changed what this does: setting `SSL_CERT_FILE` or `SSL_CERT_DIR` now replaces
> the platform verifier rather than adding to it, so the file becomes the only trust root the process has.
> macOS has no default path for either variable, so nothing is inherited. Point it at a bundle carrying every
> root the process needs — for this CLI that is just the controller's issuer. `GODEBUG=x509sslcertoverrideplatform=0`
> restores the platform verifier and makes the variables inert again on those platforms.

On Linux the variable has always overridden the default bundle location, and the same rule applies: the file must carry what is needed.

## `--insecure`

> [!CAUTION]
> `--insecure` disables certificate verification, so the connection is no longer authenticated and a machine on
> the path can read the token. Use it against a lab controller, and never where the answer matters.

Every run that uses it logs a warning.

## Containers

The published image is `scratch` plus the binary, a certificate bundle, the licence and the third-party notices. It runs as uid 65534 and has no shell, so a compromise of the process has nothing to invoke.
