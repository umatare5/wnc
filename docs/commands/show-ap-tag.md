# wnc show ap-tag

One row per access point: which tags are in force on it, and where they came from.

```bash
wnc show ap-tag
```

```plaintext
AP Name    AP MAC             Misconfigured  Misconfig Reason  Tag Source  Filter Name  Policy Tag      Site Tag        RF Tag        AP Profile         Flex Profile         Controller
TEST-AP01  00:00:5e:00:53:01  No             -                 Static      -            test-wlan-flex  test-site-flex  test-inside   test-ap-profile01  test-flex-profile01  WNC1
TEST-AP02  00:00:5e:00:53:02  No             -                 Static      -            test-wlan-flex  test-site-flex  test-outside  test-ap-profile01  test-flex-profile01  WNC1
TEST-AP03  00:00:5e:00:53:03  No             -                 Static      -            test-wlan-flex  test-site-flex  test-inside   test-ap-profile01  test-flex-profile01  WNC1
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field              | Meaning                                                             |
| :----------------- | :------------------------------------------------------------------ |
| `ap_name`          | Access point name                                                   |
| `ap_mac`           | Base radio address, as `show ap tag summary` names it               |
| `misconfigured`    | Whether the controller flags the tag assignment as broken           |
| `misconfig_reason` | Why it flagged it, on a release that publishes a reason             |
| `tag_source`       | How the tags were chosen: Static, Filter, AP-PnP, Default, Location |
| `filter_name`      | The filter that assigned the tags, when the source is Filter        |
| `policy_tag`       | Policy tag in force                                                 |
| `site_tag`         | Site tag in force                                                   |
| `rf_tag`           | RF tag in force                                                     |
| `ap_profile`       | AP join profile named by the configured site tag                    |
| `flex_profile`     | FlexConnect profile named by the configured site tag                |
| `controller`       | The controller this row was read from                               |

## Notes

- **These are the resolved tags** — the configured ones differ unless `tag_source` is Static
- **The two profiles come from the configured site tag** — no resolved form of either exists
- **So the profiles match Site Tag under Static alone** — any other source can leave them apart
- **A dash in `misconfig_reason` is not "none"** — the domain's own member for that renders `None`
- **This view is the outcome** — the three tag views hold a tag that is bound to nothing

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Access points the controller considers misconfigured:

```bash
wnc show ap-tag -f json | jq -r '.[] | select(.misconfigured == true) | .ap_name'
```

Access points whose tags were not assigned statically, which is where the profile columns and the tag columns can disagree:

```bash
wnc show ap-tag -f json | jq -r '.[] | select(.tag_source != "Static") | "\(.ap_name) \(.tag_source)"'
```

Grouped by RF tag:

```bash
wnc show ap-tag -b rf_tag
```
