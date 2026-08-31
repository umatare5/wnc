# wnc show rf-tag

One row per RF tag: which RF profile it names on each band.

```bash
wnc show rf-tag
```

```plaintext
RF Tag          Description                   2.4 GHz Profile    5 GHz Profile      6 GHz Profile            Controller
default-rf-tag  Preconfigured default RF tag  default_rf_24gh    default_rf_5gh     default-rf-profile-6ghz  WNC3
test-inside     -                             test-rf-profile01  test-rf-profile03  test-rf-profile05        WNC3
test-outside    -                             test-rf-profile01  test-rf-profile04  test-rf-profile05        WNC3
```

## Flags

The shared flags are in [`configuration.md`](../configuration.md#flags).

## Columns

| Field           | Meaning                                 |
| :-------------- | :-------------------------------------- |
| `rf_tag`        | The list key                            |
| `description`   | Description the tag carries             |
| `profile_24ghz` | RF profile on 2.4 GHz, the 802.11b leaf |
| `profile_5ghz`  | RF profile on 5 GHz, the 802.11a leaf   |
| `profile_6ghz`  | RF profile on 6 GHz                     |
| `controller`    | The controller this row was read from   |

## Notes

- **The b leaf is 2.4 and the a leaf is 5** — the one pairing this view can invert unnoticed
- **[`show overview`](show-overview.md) is the outcome** — this view is the tag that produced it
- **The per-slot radio profile list is out of reach** — the three band columns are all there is

The readings behind these sit in [`measurements.md`](../measurements.md).

## Examples

Every RF tag before deciding what to delete:

```bash
wnc show rf-tag --controller 192.168.0.1
```

Tags with no 6 GHz profile, which is where a Wi-Fi 6E radio falls back to the default:

```bash
wnc show rf-tag --format json | jq -r '.[] | select(.profile_6ghz == null) | .rf_tag'
```

Which tags name one profile on any band:

```bash
wnc show rf-tag --format json \
  | jq -r '.[] | select([.profile_24ghz, .profile_5ghz, .profile_6ghz] | index("test-rf-profile02")) | .rf_tag'
```

Grouped by the 5 GHz profile:

```bash
wnc show rf-tag --sort-by profile_5ghz
```
