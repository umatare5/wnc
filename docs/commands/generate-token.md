# wnc generate-token

Prints the Basic auth token every other command needs.

```bash
wnc generate-token -u admin -p 'the-password'
```

```plaintext
YWRtaW46dGhlLXBhc3N3b3Jk
```

## Flags

| Option             | Source order                                 |
| ------------------ | -------------------------------------------- |
| `--username`, `-u` | Flag, then `$WNC_USERNAME`                   |
| `--password`, `-p` | Flag, then `$WNC_PASSWORD`, then piped stdin |

## Notes

**Prefer a pipe or the environment to the flag.** A password given as `-p` is visible to every process on the host for as long as the command runs.

```bash
read -rs WNC_PASSWORD
printf '%s' "$WNC_PASSWORD" | wnc generate-token -u admin
```

**There is no interactive prompt.** Suppressing the echo needs a dependency this CLI does not take, and a prompt that echoes leaves the password on screen and in the scrollback. With no flag, no environment variable and a terminal on stdin, the command exits 2 and says which of the three to use.

**Only the line ending is stripped from stdin.** A password may contain spaces.

**A username may not contain a colon.** RFC 7617 gives the colon to the separator, so such a username cannot be encoded — the command says so rather than letting the controller truncate it.

## Storing the result

A configuration file at mode `0600` comes first, being the one place the token reaches neither the shell history nor the process arguments. See [`configuration.md`](../configuration.md#where-to-keep-the-token) for the whole order, and [`configuration.md`](../configuration.md#configuration-file) for the file itself.

```bash
export WNC_ACCESS_TOKEN="$(printf '%s' "$WNC_PASSWORD" | wnc generate-token -u admin)"
export WNC_CONTROLLER="wnc1.example.internal"
```
