# wnc set policy-tag, site-tag, rf-tag

Create or update one tag on a controller.

```bash
wnc set rf-tag --name <name> [--profile-5ghz <profile>]
wnc set site-tag --name <name> [--ap-join-profile <profile>]
wnc set policy-tag --name <name> [--wlan <profile> --policy-profile <profile>]
```

```plaintext
Create RF tag test-rf-inside on WNC1? [y/N]: y
RF tag test-rf-inside on WNC1: created
```

A name the controller does not hold is created and one it holds is updated, so the same invocation may be repeated. A field no flag names is left as it is rather than cleared.

This tree and [delete](delete-tag.md) are the only commands that write a tag — an access point's [administrative state](enable-disable.md) is the other configuration this CLI writes.

## Flags

| Option              | Meaning                                                |
| :------------------ | :----------------------------------------------------- |
| `--name`            | The tag's name, at most 32 characters                  |
| `--description`     | Description, on all three kinds                        |
| `--wlan`            | WLAN profile to bind, `policy-tag` only                |
| `--policy-profile`  | Policy profile the WLAN is bound to, `policy-tag` only |
| `--ap-join-profile` | AP join profile to bind, `site-tag` only               |
| `--flex-profile`    | Flex profile to bind, `site-tag` only                  |
| `--local-site`      | Mark the site local, `site-tag` only                   |
| `--profile-24ghz`   | 2.4 GHz RF profile to bind, `rf-tag` only              |
| `--profile-5ghz`    | 5 GHz RF profile to bind, `rf-tag` only                |
| `--profile-6ghz`    | 6 GHz RF profile to bind, `rf-tag` only                |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--name` is the tag's name and the list key the controller stores it under. The key leaves declare the pattern `[!-~]([ -~]*[!-~])?` and no length, so a name is printable ASCII, may carry an inner space, and may not begin or end with one.

**The 32-character limit is the controller's and the schema does not declare it.** Measured on 17.12.8, one kind at a time: a 32-character name is accepted and a 33-character one answers `400 Validation failed ... <kind> Tag name should not exceed 32 characters` on all three kinds. The model is not the arbiter here, so the CLI applies the measured bound rather than the declared one.

A leading or trailing space does reach the check. The name arrives as a flag value, which the argument parser does not trim, so a padded name is refused at exit 2 with nothing sent — where a positional argument would have arrived already trimmed and been created without the padding.

## Guards

**`--wlan` and `--policy-profile` are a pair.** Both are keys of the WLAN-policy list, so one without the other binds nothing and is refused before any request.

**The name is checked before any request.** A name the pattern rejects is a usage fault at exit 2, so nothing reaches a controller.

**`--dry-run` reports whether the name would be created or updated** and writes nothing.

**An update naming no field writes nothing** and says so, rather than reporting a write it did not make. A field that cannot take effect on its own is refused rather than dropped, for the same reason.

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the read before the write failed.

## Notes

- **A write is lost on a reload** — run [`save-config`](save-config.md) to persist it
- **An update leaves unnamed fields alone** — every kind sends a merge `PATCH`
- **`--flex-profile` implies a non-local site** — so naming both is a usage fault and not a 400
- **The controller enforces no referential integrity** — a profile that exists nowhere persists
- **Read the matching view first** — it is what tells a create from an update before sending
- **An RF tag write never names the per-slot list** — the controller rejects it outright

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Create an RF tag bound to two band profiles:

```bash
wnc set rf-tag --name test-rf-inside \
  --description "inside coverage" \
  --profile-24ghz test-rf-24 \
  --profile-5ghz test-rf-5
```

Add a band to it later, leaving the rest as it is:

```bash
wnc set rf-tag --name test-rf-inside --profile-6ghz test-rf-6
```

Bind a WLAN to a policy profile:

```bash
wnc set policy-tag --name test-wlan --wlan test-corp --policy-profile test-corp-policy
```

Check what would happen without acting:

```bash
wnc --dry-run set site-tag --name test-site --ap-join-profile test-join
```
