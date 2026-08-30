# wnc delete policy-tag, site-tag, rf-tag

Delete one tag from a controller. This tree and [set](set-tag.md) are the only commands that write a tag — an access point's [administrative state](enable-disable.md) is the other configuration this CLI writes.

```bash
wnc delete rf-tag --name <name>
wnc delete site-tag --name <name>
wnc delete policy-tag --name <name>
```

```plaintext
Delete RF tag lab-rf-inside on 192.168.0.231? [y/N]: y
RF tag lab-rf-inside on 192.168.0.231: deleted
```

## Options

| Option           | Meaning                               |
| ---------------- | ------------------------------------- |
| `--name`         | The tag's name, at most 32 characters |
| `--controller`   | The one controller holding the tag    |
| `--access-token` | Basic auth token for that controller  |
| `--insecure`     | Skip TLS certificate verification     |
| `--timeout`      | Request timeout                       |
| `--yes`          | Act without the confirmation prompt   |
| `--dry-run`      | Name the target and change nothing    |

## The target

`--name` is the tag's name — the `policy_tag`, `site_tag` or `rf_tag` column of the matching `wnc show` view. It is checked against the rule [set](set-tag.md#the-target) applies, because both trees run the same check: printable ASCII throughout, no leading or trailing space, and at most 32 characters.

That keeps two answers apart. A name the rule refuses is a usage fault at exit 2 and never becomes a request, while a well-formed name the controller does not hold exits 1 after the read below.

## Guards

**One tag per invocation.** `--name` names one, and a second occurrence is rejected rather than silently replacing the first.

**One controller per invocation.** A tag is deleted on the controller named, so naming two is a usage fault. Nothing is sent.

**The target is resolved before you are asked.** The controller is read for the name first, so a name it does not hold is reported as such and no delete is sent. Without the read, RESTCONF's own `404` would reach the operator as a status rather than as the plain fact.

**A prompt you cannot answer is a fault.** With stdin piped, `--yes` is required — the run is refused before any request goes out.

**`--dry-run` reads and stops.** It names what would be deleted and writes nothing.

## Exit codes

| Code | Meaning                                                |
| :--- | :----------------------------------------------------- |
| 0    | The controller accepted the delete, or you answered no |
| 1    | The controller refused, or holds no tag of that name   |
| 2    | Usage fault. Nothing was sent to a controller          |

## Notes

**A delete is lost on a reload until it is saved.** The record leaves the running configuration and stays in the startup configuration, so a controller that reloads brings the tag back. Measured on all three releases in scope. Run [wnc save-config](save-config.md) to persist the delete.

**Deleting a tag an access point resolves to was not measured.** The three configuration modules declare no `leafref` and no `require-instance` on any release measured, so the schema imposes no constraint that would refuse it — but what the controller actually does was not tested, because doing so would have meant deleting a tag in use on the lab. Read `wnc show ap-tag` first and move the access points yourself.

**The controller keeps a dangling reference rather than refusing it.** The access-point-to-tag binding types its tag names as plain strings, and the lab carries live examples of a tag naming a profile that exists nowhere. Absence of a complaint is not evidence that a delete was safe.

**`wnc show policy-tag`, `wnc show site-tag` and `wnc show rf-tag` list the tags that exist.** Read the matching view before a delete, and `wnc show ap-tag` for the ones access points resolve to.

**A default tag is still deletable as far as this command knows.** `default-rf-tag`, `default-site-tag` and `default-policy-tag` are Cisco's preconfigured tags. This command applies no special rule to them.

## Examples

Delete a tag nothing resolves to:

```bash
wnc show ap-tag -f json | jq -r '.[] | .rf_tag' | sort -u
wnc delete rf-tag --name lab-rf-retired
```

Check the target without acting:

```bash
wnc --dry-run delete policy-tag --name lab-wlan
```
