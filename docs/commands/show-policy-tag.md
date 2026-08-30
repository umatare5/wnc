# wnc show policy-tag

One row per WLAN binding a policy tag carries: what exists, and what each tag binds.

```bash
wnc show policy-tag
```

```plaintext
Policy Tag          Description                       WLAN         Policy Profile     Controller
default-policy-tag  Preconfigured default policy-tag  -            -                  192.168.0.233
labo-wlan-flex      -                                 labo-p736b2  labo-wlan-profile  192.168.0.233
labo-wlan-flex      -                                 labo-p736b5  labo-wlan-profile  192.168.0.233
labo-wlan-flex      -                                 labo-t6c73d  labo-wlan-profile  192.168.0.233
```

## Columns

| Field            | Meaning                                      |
| ---------------- | -------------------------------------------- |
| `policy_tag`     | Policy tag name, the list key                |
| `description`    | Description the tag carries                  |
| `wlan`           | WLAN profile the tag binds, not the SSID     |
| `policy_profile` | Policy profile that WLAN profile is bound to |
| `controller`     | The controller this row was read from        |

## Notes

**A tag binding three WLANs is three rows.** That is how `wnc show wlan` pairs a WLAN with each policy profile bound to it, read from the other side. A tag binding nothing keeps a row of its own with the two binding columns empty — it exists, it is a candidate for `wnc delete policy-tag`, and dropping it would hide exactly the tag most likely to be deleted.

**The device CLI names the second binding column differently.** `show wireless tag policy detailed <name>` prints `WLAN Profile Name` and `Policy Name` over the same pair. This view keeps `policy_profile`, which is the leaf's own spelling and the one `wnc set policy-tag --policy-profile` takes.

**`wlan` is the WLAN profile name and not the SSID.** The binding keys on the profile, and the two differ on any WLAN whose profile was not named after its SSID. `wnc show wlan` carries both, so it is where the two are matched up.

**The two binding columns never occur one without the other.** Both are keys of the `wlan-policy` list, which is why `wnc set policy-tag` refuses `--wlan` without `--policy-profile` and vice versa.

**A binding may name a WLAN profile or a policy profile that does not exist.** The three tag modules declare no leafref on any release measured, so the controller accepts and keeps such a binding. `wnc show wlan` reports how many bindings name a WLAN profile with no configuration entry — this view shows the binding itself.

**The read asks for the values in force.** This is a configuration read, so `with-defaults=report-all` is the right request: it is what separates a leaf left at its default from one the controller withheld. There is no fallback to a plain read — a plain answer would report a default as an absence.

**The whole view is one request.** No column comes from a second collection, so there is no secondary read to degrade. A failed read drops this controller's rows rather than printing a partial list.

## Examples

Every tag, with its bindings, before deciding what to delete:

```bash
wnc show policy-tag --controller 192.168.0.231
```

Tags that bind nothing, which is what a delete is usually aimed at:

```bash
wnc show policy-tag --format json | jq -r '.[] | select(.wlan == null) | .policy_tag'
```

Grouped by the policy profile the bindings point at:

```bash
wnc show policy-tag --sort-by policy_profile
```

Which tags bind one WLAN profile:

```bash
wnc show policy-tag --format json | jq -r '.[] | select(.wlan == "labo-p736b2") | .policy_tag'
```
