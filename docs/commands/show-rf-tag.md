# wnc show rf-tag

One row per RF tag: which RF profile it names on each band.

```bash
wnc show rf-tag
```

```plaintext
RF Tag          Description                   2.4 GHz Profile  5 GHz Profile        6 GHz Profile            Controller
default-rf-tag  Preconfigured default RF tag  default_rf_24gh  default_rf_5gh       default-rf-profile-6ghz  192.168.0.233
labo-inside     -                             labo-rf-24gh     labo-rf-5gh-inside   labo-rf-6gh              192.168.0.233
labo-outside    -                             labo-rf-24gh     labo-rf-5gh-outside  labo-rf-6gh              192.168.0.233
```

## Columns

| Field           | Meaning                                 |
| --------------- | --------------------------------------- |
| `rf_tag`        | RF tag name, the list key               |
| `description`   | Description the tag carries             |
| `profile_24ghz` | RF profile on 2.4 GHz, the 802.11b leaf |
| `profile_5ghz`  | RF profile on 5 GHz, the 802.11a leaf   |
| `profile_6ghz`  | RF profile on 6 GHz                     |
| `controller`    | The controller this row was read from   |

## Notes

**The device CLI calls these columns something else.** `show wireless tag rf detailed <name>` labels the same three values `2.4ghz RF Policy`, `5ghz RF Policy` and `6ghz RF Policy`, and the device's own configuration keyword is `24ghz-rf-policy`. These columns follow the YANG leaf and the write flag instead, so the name an operator passes to `wnc set rf-tag` is the name the column carries.

**2.4 GHz is the 802.11b leaf and 5 GHz is the 802.11a leaf.** That is the pairing the write path's own setters name, and it is the one thing this view can get wrong while still looking right. The three column keys are the `--profile-24ghz`, `--profile-5ghz` and `--profile-6ghz` flags of `wnc set rf-tag`.

**The read asks for the values in force, and here it is necessary.** Measured on 17.12, 17.15 and 17.18, a plain read of `rf-tags` omits all three per-band profile names on the built-in `default-rf-tag` while `report-all` returns them, so a plain answer would report a configured default as an absence. There is no fallback to a plain read.

**`wnc show overview` answers a different question.** Its RF Profile column is the profile the radio is actually on, chosen by the radio's band rather than by its slot — an XOR radio on slot 2 can be operating in 5 or 6 GHz, and it takes whichever of these columns matches. This view is the tag — that one is the outcome.

**The per-slot radio profile list is not shown.** The tag also carries `rf-tag-radio-profiles`, whose entries pair a slot with a radio profile name. `wnc set rf-tag` does not write it either, so the three band columns are the whole of what this CLI carries for an RF tag.

**A `-` means the controller sent nothing for that band.** An omitted description and an omitted profile name arrive absent rather than empty, so neither is rendered as a value.

**The whole view is one request.** No column comes from a second collection, so there is no secondary read to degrade. A failed read drops this controller's rows rather than printing a partial list.

## Examples

Every RF tag before deciding what to delete:

```bash
wnc show rf-tag --controller 192.168.0.231
```

Tags with no 6 GHz profile, which is where a Wi-Fi 6E radio falls back to the default:

```bash
wnc show rf-tag --format json | jq -r '.[] | select(.profile_6ghz == null) | .rf_tag'
```

Which tags name one profile on any band:

```bash
wnc show rf-tag --format json \
  | jq -r '.[] | select([.profile_24ghz, .profile_5ghz, .profile_6ghz] | index("labo-rf-5gh")) | .rf_tag'
```

Grouped by the 5 GHz profile:

```bash
wnc show rf-tag --sort-by profile_5ghz
```
