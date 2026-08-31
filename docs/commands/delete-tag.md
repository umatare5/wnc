# wnc delete policy-tag, site-tag, rf-tag

Delete one tag from a controller.

```bash
wnc delete rf-tag --name <name>
wnc delete site-tag --name <name>
wnc delete policy-tag --name <name>
```

```plaintext
Delete RF tag test-rf-inside on WNC1? [y/N]: y
RF tag test-rf-inside on WNC1: deleted
```

The counterpart is [set](set-tag.md), which carries the tag-name grammar and what a write can change.

## Flags

| Option   | Meaning                               |
| :------- | :------------------------------------ |
| `--name` | The tag's name, at most 32 characters |

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

This command prints no table.

## The target

`--name` is the tag's name — the `policy_tag`, `site_tag` or `rf_tag` column of the matching `wnc show` view. It is checked against the rule [set](set-tag.md#the-target) applies, because both trees run the same check: printable ASCII throughout, no leading or trailing space, and at most 32 characters.

That keeps two answers apart. A name the rule refuses is a usage fault at exit 2 and never becomes a request, while a well-formed name the controller does not hold exits 1 after the read below.

## Guards

**The target is resolved before you are asked.** The controller is read for the name first, so a name it does not hold is reported as such and no delete is sent. Without the read, RESTCONF's own `404` would reach the operator as a status rather than as the plain fact.

The rules every write keeps are in [`README.md`](../README.md#acting-on-a-controller).

## Exit codes

The five codes are in [`README.md`](../README.md#exit-codes).

Here 1 also means the controller holds no tag of that name.

## Notes

- **A delete is lost on a reload** — the tag returns until [`save-config`](save-config.md) persists it
- **Deleting a tag in use was not measured** — move the access points off it yourself first
- **A dangling reference is kept, not refused** — silence is not evidence a delete was safe
- **Read the matching view first** — and [`show ap-tag`](show-ap-tag.md) for the tags in force
- **A default tag is deletable as far as this knows** — no special rule covers the three

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Delete a tag nothing resolves to:

```bash
wnc show ap-tag -f json | jq -r '.[] | .rf_tag' | sort -u
wnc delete rf-tag --name test-rf-retired
```

Check the target without acting:

```bash
wnc --dry-run delete policy-tag --name test-wlan
```
