# wnc show site-tag

One row per site tag: which profiles it names, and whether the site is local.

```bash
wnc show site-tag
```

```plaintext
Site Tag          Description                     AP Join Profile     Flex Profile          Local Site  Controller
default-site-tag  Preconfigured default site tag  default-ap-profile  default-flex-profile  No          192.168.0.233
labo-site-flex    -                               labo-common         labo-flex             No          192.168.0.233
```

## Columns

| Field             | Meaning                               |
| ----------------- | ------------------------------------- |
| `site_tag`        | Site tag name, the list key           |
| `description`     | Description the tag carries           |
| `ap_join_profile` | AP join profile the tag names         |
| `flex_profile`    | FlexConnect profile the tag names     |
| `local_site`      | Whether the tag marks the site local  |
| `controller`      | The controller this row was read from |

## Notes

**`local_site` has three readings and they are all different.** `Yes` and `No` are what the controller reported — `-` is the controller reporting nothing. The leaf carries a schema default, which is why this read asks for the values in force.

**A flex profile is in force on a non-local site only.** The `flex-profile` leaf declares `when "../is-local-site = 'false'"` on every release in scope, so the two columns are read side by side. `wnc set site-tag` refuses `--flex-profile` together with `--local-site` for the same reason: the pair answers 400 naming the failed `when` expression, measured on 17.12.8.

**`ap_join_profile` is spelt as the leaf and the flag spell it.** `wnc show ap-tag` calls the same value `ap_profile`, because there it comes from a leaf of that name inside the operational tag container rather than from the site tag itself. The device CLI is a third spelling: `show wireless tag site detailed <name>` prints `AP Profile`, and its configuration keyword is `ap-profile`.

**Five columns, not eleven.** The list carries six more leaves — the fabric control-plane and multicast pair, the image-download profile, ARP caching, DHCP broadcast and a load estimate. None is rendered, because none has been verified against a controller in scope.

**The read asks for the values in force.** This is a configuration read, so `with-defaults=report-all` is the right request: it is what separates a leaf left at its default from one the controller withheld. There is no fallback to a plain read — a plain answer would report a default as an absence.

**The whole view is one request.** No column comes from a second collection, so there is no secondary read to degrade. A failed read drops this controller's rows rather than printing a partial list.

## Examples

Every site tag before deciding what to delete:

```bash
wnc show site-tag --controller 192.168.0.231
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
