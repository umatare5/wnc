# wnc set policy-tag, site-tag, rf-tag

Create or update one tag on a controller. This tree and [delete](delete-tag.md) are the only commands that write a tag — an access point's [administrative state](enable-disable.md) is the other configuration this CLI writes.

```bash
wnc set rf-tag --name <name> [--profile-5ghz <profile>]
wnc set site-tag --name <name> [--ap-join-profile <profile>]
wnc set policy-tag --name <name> [--wlan <profile> --policy-profile <profile>]
```

```plaintext
Create RF tag lab-rf-inside on 192.168.0.231? [y/N]: y
RF tag lab-rf-inside on 192.168.0.231: created
```

A name the controller does not hold is created and one it holds is updated, so the same invocation may be repeated. A field no flag names is left as it is rather than cleared.

## Options

| Option              | Meaning                                                |
| ------------------- | ------------------------------------------------------ |
| `--name`            | The tag's name, at most 32 characters                  |
| `--controller`      | The one controller the tag is written on               |
| `--access-token`    | Basic auth token for that controller                   |
| `--insecure`        | Skip TLS certificate verification                      |
| `--timeout`         | Request timeout                                        |
| `--yes`             | Act without the confirmation prompt                    |
| `--dry-run`         | Report what would be written and change nothing        |
| `--description`     | Description, on all three kinds                        |
| `--wlan`            | WLAN profile to bind, `policy-tag` only                |
| `--policy-profile`  | Policy profile the WLAN is bound to, `policy-tag` only |
| `--ap-join-profile` | AP join profile to bind, `site-tag` only               |
| `--flex-profile`    | Flex profile to bind, `site-tag` only                  |
| `--local-site`      | Mark the site local, `site-tag` only                   |
| `--profile-24ghz`   | 2.4 GHz RF profile to bind, `rf-tag` only              |
| `--profile-5ghz`    | 5 GHz RF profile to bind, `rf-tag` only                |
| `--profile-6ghz`    | 6 GHz RF profile to bind, `rf-tag` only                |

`--wlan` and `--policy-profile` are a pair. Both are keys of the WLAN-policy list, so one without the other binds nothing and is refused before any request.

## The target

`--name` is the tag's name and the list key the controller stores it under. The key leaves declare the pattern `[!-~]([ -~]*[!-~])?` and no length, so a name is printable ASCII, may carry an inner space, and may not begin or end with one.

**The 32-character limit is the controller's and the schema does not declare it.** Measured on 17.12.8, one kind at a time: a 32-character name is accepted and a 33-character one answers `400 Validation failed ... <kind> Tag name should not exceed 32 characters` on all three kinds. The model is not the arbiter here, so the CLI applies the measured bound rather than the declared one.

A leading or trailing space does reach the check. The name arrives as a flag value, which the argument parser does not trim, so a padded name is refused at exit 2 with nothing sent — where a positional argument would have arrived already trimmed and been created without the padding.

## Guards

**One tag per invocation.** `--name` names one, and a second occurrence is rejected rather than silently replacing the first.

**One controller per invocation.** A tag is written on the controller named, so naming two is a usage fault. Nothing is sent.

**The name is checked before any request.** A name the pattern rejects is a usage fault at exit 2, so nothing reaches a controller.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It reports whether the name would be created or updated and writes nothing.

**An update naming no field writes nothing** and says so, rather than reporting a write it did not make. A field that cannot take effect on its own is refused rather than dropped, for the same reason.

## Exit codes

| Code | Meaning                                               |
| :--- | :---------------------------------------------------- |
| 0    | The controller accepted the write, or you answered no |
| 1    | The controller refused, or the read before it failed  |
| 2    | Usage fault. Nothing was sent to a controller         |

## Notes

**A write is lost on a reload until it is saved.** A RESTCONF write reaches the running configuration and nowhere else, so a controller that reloads comes back without it. Measured on all three releases in scope. Run [wnc save-config](save-config.md) to persist it.

**An update leaves unnamed fields alone.** Every kind sends a merge `PATCH`, so a leaf the payload omits keeps the value the controller holds. Measured on 17.12.8 and again on 17.15.6: a second write naming only one profile left the description and the other profile from the first one in place.

**`--flex-profile` clears `--local-site`, and naming both is refused.** The leaf declares `when "../is-local-site = 'false'"` on every release in scope and `is-local-site` defaults to **true**, so a body naming a flex profile without the flag answers `400 the 'when' expression ... failed` — measured on 17.12.8. A create therefore sends `is-local-site: false` alongside the profile, exactly as an update does, and `--local-site --flex-profile` together is a usage fault rather than a 400.

**The controller enforces no referential integrity.** The three configuration modules declare no `leafref` and no `require-instance` on any release measured, so a profile name that does not exist is accepted and persists. Read `wnc show wlan` and the controller's own profile lists before binding.

**`wnc show policy-tag`, `wnc show site-tag` and `wnc show rf-tag` list what exists.** Read the matching view first — it is what tells a create from an update before either is sent.

**An RF tag write never names `rf-tag-radio-profiles`.** The controller rejects that container with a null list — `400 invalid value for: rf-tag-radio-profile`, measured on 17.12.8 and again on 17.15.6 — so the per-slot radio profile list is out of reach of this command.

**Deleting a tag an access point resolves to was not measured.** See [delete](delete-tag.md#notes).

## Examples

Create an RF tag bound to two band profiles:

```bash
wnc set rf-tag --name lab-rf-inside \
  --description "inside coverage" \
  --profile-24ghz lab-rf-24 \
  --profile-5ghz lab-rf-5
```

Add a band to it later, leaving the rest as it is:

```bash
wnc set rf-tag --name lab-rf-inside --profile-6ghz lab-rf-6
```

Bind a WLAN to a policy profile:

```bash
wnc set policy-tag --name lab-wlan --wlan lab-corp --policy-profile lab-corp-policy
```

Check what would happen without acting:

```bash
wnc --dry-run set site-tag --name lab-site --ap-join-profile lab-join
```
