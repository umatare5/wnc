# Configuration

Every setting this CLI reads, in precedence order: a flag, then the environment, then the configuration file.

## Flags

[`help.md`](help.md) transcribes `--help` for every command. What decides which flags a command takes is what that command does:

| Set                      | Flags                                                              |
| :----------------------- | :----------------------------------------------------------------- |
| Every command            | `--config`, `--log-level`                                          |
| Contacting a controller  | `--controller`, `--access-token`, `--insecure`, `--timeout`        |
| Reading, so `show` only  | `--format`, `--pretty`, `--sort-by`, `--sort-keys`, `--sort-order` |
| Acting, so asking first  | `--yes`                                                            |

`--dry-run` is a root flag rather than a per-command one. On the root it validates the configuration and contacts nothing, and on a command that acts it stops before the request.

> [!WARNING]
> `--ap-name` is keyed into the request URL as a path segment with nothing bounding its content, because nothing measured bounds the name a controller accepts — a secret pasted there has left the host before any diagnostic reports it. `deauth --mac` is validated as an address first and refused before a request goes out, and a tag `--name` is rejected locally at exit 2. `deauth --username` reaches no URL and leaves only in the RPC's request body, which discloses it just as fully, somewhere a proxy log would not show it.

## Environment Variables

Three variables reach every command:

| Variable           | Description                                             |
| :----------------- | :------------------------------------------------------ |
| `WNC_CONTROLLER`   | Controller `host[:port]`, comma separated for several   |
| `WNC_ACCESS_TOKEN` | Basic auth token applied to every controller            |
| `WNC_CONFIG`       | Configuration file path, replacing the default location |

Two more are read by `wnc generate-token` alone, which contacts no controller:

| Variable       | Description         |
| :------------- | :------------------ |
| `WNC_USERNAME` | Controller username |
| `WNC_PASSWORD` | Controller password |

## Configuration File

A file keeps the controller list out of every invocation and the token out of the shell history and the process arguments. [`examples/config.json`](../examples/config.json) is a working file whose `note` fields carry the remark JSON has no comment syntax for.

The path is `--config`, then `$WNC_CONFIG`, then `$XDG_CONFIG_HOME/wnc/config.json` where that variable is set and `~/.config/wnc/config.json` where it is not. The two defaults are a branch rather than a chain, so an unreadable XDG path is not followed by a look under the home directory.

The token is file-wide rather than per entry, so one run reads every controller with one credential however many hosts the file lists. A host named on `--controller` that the file does not list takes that same token, which is what keeps the file usable for a one-off read.

The read is strict: an unknown key, a duplicated key, a key differing only in case, a comment and a trailing comma are each rejected. An unknown or case-differing key is reported as the JSON pointer that located it, and the other three reach the operator with `jsontext`'s own locator instead — a duplicated key is not wrapped as a semantic fault, so it takes that path rather than the pointer one. In no case is the offending value echoed.

Check a hand-edited file without contacting anything:

```bash
wnc --config ./config.json --dry-run
```

## Where to Keep the Token

In order of preference:

1. **A configuration file at mode `0600`.** The token reaches neither the shell history nor the process arguments, and a file readable beyond its owner logs a warning naming the mode.
2. **`$WNC_ACCESS_TOKEN`.** Visible to anything that can read the process environment.
3. **`--access-token`.** Visible to every process on the host through the process list, so use it interactively only.

`--controller` carries no token, so naming a host on the command line costs nothing — the file's one token covers every host, listed there or not.

The same order applies to `wnc generate-token`, which prints the token it assembles: prefer a piped password to `$WNC_PASSWORD`, and either to `--password`.
