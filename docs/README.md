# Documentation

Reference pages for `wnc`, carrying the detail behind the [README](../README.md).

| Page                                       | Focus                                                  |
| :----------------------------------------- | :----------------------------------------------------- |
| [`help.md`](help.md)                       | The help text of every command but `completion`        |
| [`configuration.md`](configuration.md)     | The environment variables and the configuration file   |
| [`troubleshooting.md`](troubleshooting.md) | A failure's `cause=` token, then a refusal's wording   |
| [`measurements.md`](measurements.md)       | Every reading taken on a live controller, and the gaps |
| [`testing.md`](testing.md)                 | How the suite is arranged and what it asserts          |

The per-command pages live under [`commands/`](commands/), each linked from the table in the [README](../README.md).

## Technical Information

### Output

Every `show` command prints either a borderless table or a flat JSON array, and both carry the same cells.

- **Borderless** — no rules, and the first column starts at column zero, so `awk` and `cut` work
- **Absence** — a cell the controller did not report is `-`, and absent from the JSON rather than null
- **Sort keys** — the JSON field names are exactly what `--sort-by` takes, and `--sort-keys` prints them
- **Types survive** — a number stays a number, and an empty result is `[]`
- **Units are the table's** — glued to the number so a cell stays one field, and never in the JSON
- **Sorting reads the number** — `--sort-by channel` puts 6 before 11 rather than ordering the text
- **Styling is the table's** — `--pretty` borders it and glyphs the state columns, leaving the JSON alone

> [!NOTE]
> `-` is not zero, and reading it as one inverts the reading: a Monitor radio has no channel, an unauthenticated client has no username, and a WLAN with no policy tag bound has no interface. A whole column of `-` usually means a secondary read failed, which [`troubleshooting.md`](troubleshooting.md) indexes by its `endpoint=` field.

### Exit codes

One invocation answers with one of five codes, and the fleet decides which.

| Code | Meaning                                                        |
| :--- | :------------------------------------------------------------- |
| 0    | Every controller answered, or you declined a prompt            |
| 1    | No controller answered, or the CLI failed internally           |
| 2    | Usage or configuration fault. Nothing was sent to a controller |
| 3    | Partial: at least one read failed and at least one succeeded   |
| 130  | Interrupted. No partial table is printed                       |

- **2 is decided locally** — the fault is found before a request goes out, so no controller saw it
- **3 still prints** — the table holds every row that was read, and stderr says what was not
- **130 prints nothing** — a partial table from an interrupted fan-out would be a reading nobody took

### Acting on a controller

Every command that acts obeys the same order before it does.

- **Local checks first** — a missing, empty or repeated target flag is exit 2, nothing sent
- **One controller** — naming a second is exit 2, and nothing is sent either
- **Answerable first** — with stdin piped, `--yes` or `--dry-run` is required
- **Resolved before asked** — the target is read, so a declined run has still read it
- **Asked before written** — `--yes` answers the prompt, and nothing writes before it

`--dry-run` is a root flag rather than a per-command one, so it precedes the subcommand: `wnc --dry-run reset ap --ap-name TEST-AP01`. It resolves the target and posts nothing, so it cannot report whether anything needed doing.

Half the acting leaves post an RPC that declares no output, where a `204 No Content` establishes only that the instruction was accepted. The other half read a result back and say what happened.
