# wnc show site-tag

One row per site tag: which profiles it names, and whether the site is local.

```bash
wnc show site-tag
```

```plaintext
Site Tag          Description                     AP Join Profile     Flex Profile          Local Site  Controller
default-site-tag  Preconfigured default site tag  default-ap-profile  default-flex-profile  No          WNC3
test-site-flex    -                               test-ap-profile01   test-flex-profile01   No          WNC3
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field             | Meaning                               |
| :---------------- | :------------------------------------ |
| `site_tag`        | The list key                          |
| `description`     | Description the tag carries           |
| `ap_join_profile` | AP join profile the tag names         |
| `flex_profile`    | FlexConnect profile the tag names     |
| `local_site`      | Whether the tag marks the site local  |
| `controller`      | The controller this row was read from |

## Notes

- **`local_site` has three readings** — `Yes` and `No` are reported, `-` is the controller silent
- **A flex profile is in force on a non-local site alone** — `set site-tag` refuses the pair
- **Spelt as the leaf and the flag do** — [`show ap-tag`](show-ap-tag.md) reads another leaf and says `ap_profile`

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Every site tag before deciding what to delete:

```bash
wnc show site-tag --controller 192.168.0.1
```

Non-local sites, which are the ones a flex profile can be bound to:

```bash
wnc show site-tag --format json | jq -r '.[] | select(.local_site == false) | .site_tag'
```

Site tags naming no AP join profile:

```bash
wnc show site-tag --format json | jq -r '.[] | select(.ap_join_profile == null) | .site_tag'
```

Grouped by the AP join profile they name:

```bash
wnc show site-tag --sort-by ap_join_profile
```
