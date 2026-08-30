# wnc show ap-tag

One row per access point: which tags are in force on it, and where they came from.

```bash
wnc show ap-tag
```

```plaintext
AP Name    AP MAC             Misconfigured  Misconfig Reason  Tag Source  Filter Name  Policy Tag      Site Tag        RF Tag        AP Profile   Flex Profile  Controller
TEST-AP01  aa:bb:cc:dd:ee:ff  No             -                 Static      -            labo-wlan-flex  labo-site-flex  labo-inside   labo-common  labo-flex     192.168.0.231
TEST-AP02  aa:bb:cc:dd:ee:ff  No             -                 Static      -            labo-wlan-flex  labo-site-flex  labo-outside  labo-common  labo-flex     192.168.0.231
TEST-AP03  aa:bb:cc:dd:ee:ff  No             -                 Static      -            labo-wlan-flex  labo-site-flex  labo-inside   labo-common  labo-flex     192.168.0.231
```

## Columns

| Field              | Meaning                                                             |
| ------------------ | ------------------------------------------------------------------- |
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

**The three tags are the resolved ones — what is in force.** The controller also publishes the configured tags, which differ whenever the source is not Static.

**The two profiles are not resolved, because no resolved form exists.** The schema publishes no resolved counterpart of either on any release measured, so both come from the configured site tag. They describe the same site tag as the Site Tag column only while the configured and resolved site tags agree — which is exactly when `tag_source` is Static.

**`misconfigured` shows `-` when the controller said nothing.** On 17.12 it sends an explicit `false` for a healthy access point, so `No` there is a reading and not a substitute for silence.

**`misconfig_reason` has been seen carrying only the domain's "no misconfiguration" member.** The leaf appears at 17.15 — present on 3 of 3 records on 17.15.6 — and its domain grows from three members to four at 17.18, while 17.12 does not declare it. The display strings come from the model's own descriptions rather than from the device CLI, which prints no heading for this quantity at all.

**A dash in `misconfig_reason` is not "no misconfiguration".** The domain has its own member for that, `apmgr-no-misconfig`, rendered `None`. A dash means the release did not report the leaf.

**`filter_name` is empty unless a filter assigned the tags.** The controller sends an empty string where no filter exists — measured on 17.12.8, present on 3 of 3 records with no `ap-filter-configs` configured — so the cell reads `-` for the same reason an absent leaf does: there is no filter name to report.

**This view is the outcome and the three tag views are the definitions.** It reports which tags are in force on an access point and cannot report a tag that exists and is bound to nothing — which is the tag a delete is usually aimed at. `wnc show policy-tag`, `wnc show site-tag` and `wnc show rf-tag` read the three configuration lists that `wnc set` and `wnc delete` write.

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
