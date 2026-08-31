# wnc show policy-tag

One row per WLAN binding a policy tag carries: what exists, and what each tag binds.

```bash
wnc show policy-tag
```

```plaintext
Policy Tag          Description                       WLAN                 Policy Profile         Controller
default-policy-tag  Preconfigured default policy-tag  -                    -                      WNC3
test-wlan-flex      -                                 test-wlan-profile01  test-policy-profile01  WNC3
test-wlan-flex      -                                 test-wlan-profile02  test-policy-profile01  WNC3
test-wlan-flex      -                                 test-wlan-profile03  test-policy-profile01  WNC3
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field            | Meaning                                      |
| :--------------- | :------------------------------------------- |
| `policy_tag`     | The list key                                 |
| `description`    | Description the tag carries                  |
| `wlan`           | WLAN profile the tag binds, not the SSID     |
| `policy_profile` | Policy profile that WLAN profile is bound to |
| `controller`     | The controller this row was read from        |

## Notes

- **A tag binding three WLANs is three rows** — one binding nothing keeps a row with both empty
- **`wlan` is the profile name** — [`show wlan`](show-wlan.md) carries both, so it matches the two up
- **Both binding columns are list keys** — which is why `set policy-tag` refuses one alone
- **A binding may name nothing that exists** — this view shows it, [`show wlan`](show-wlan.md) counts it

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Every tag, with its bindings, before deciding what to delete:

```bash
wnc show policy-tag --controller 192.168.0.1
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
wnc show policy-tag --format json | jq -r '.[] | select(.wlan == "test-wlan-profile01") | .policy_tag'
```
