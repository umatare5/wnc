package render

import (
	"strconv"
	"strings"
	"time"
)

// Sizes of the IEC units, and the boundaries of the duration units below.
const (
	iecBase        = 1024
	secondsPerMin  = 60
	secondsPerHour = 3600
	secondsPerDay  = 86400
)

// iecUnits are the suffixes applied above one kibibyte. A controller counter is a
// 64-bit octet count, so the list stops where that stops mattering.
var iecUnits = [...]string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}

// Str renders a string cell, mapping the empty string to Absent. Almost every
// string leaf on this SDK's read paths sits in a non-pointer struct, so an omitted
// container and an omitted leaf both arrive as "" and neither is a value.
func Str(s string) string {
	if s == "" {
		return Absent
	}

	return s
}

func StrPtr(p *string) string {
	if p == nil {
		return Absent
	}

	return Str(*p)
}

func Int[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

// IntPtr renders an optional integer, which is the shape to prefer: it separates
// "the controller said zero" from "the controller said nothing".
func IntPtr[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](p *T) string {
	if p == nil {
		return Absent
	}

	return Int(*p)
}

// UnitPtr renders an optional integer with its unit glued to the number, so the cell stays one
// whitespace-delimited field. An unreported value is Absent with no unit, because "-dBm" would
// read as a measurement rather than the lack of one.
func UnitPtr[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](
	p *T, unit string,
) string {
	if p == nil {
		return Absent
	}

	return Int(*p) + unit
}

// Bool renders an optional boolean as Yes or No. A nil is Absent rather than No:
// the controller omits some of these leaves, and reading that as No states a fact
// the controller did not.
func Bool(p *bool) string {
	switch {
	case p == nil:
		return Absent
	case *p:
		return "Yes"
	default:
		return "No"
	}
}

// IEC renders an octet count in binary units. The JSON output carries the raw
// count, so this affects the table alone.
func IEC(p *uint64) string {
	if p == nil {
		return Absent
	}

	n := float64(*p)
	if n < iecBase {
		return strconv.FormatUint(*p, 10) + "B"
	}

	unit := ""

	for _, u := range iecUnits {
		n /= iecBase
		unit = u

		if n < iecBase {
			break
		}
	}

	return strconv.FormatFloat(n, 'f', 1, 64) + unit + "B"
}

// Duration renders a second count as the two largest non-zero units, which is what
// makes an uptime of weeks and an association of minutes both readable in one
// column. The JSON output carries the raw seconds.
func Duration(p *int64) string {
	if p == nil {
		return Absent
	}

	s := *p
	if s < 0 {
		return Absent
	}

	var b strings.Builder

	switch {
	case s >= secondsPerDay:
		writeUnit(&b, s/secondsPerDay, "d")
		writeUnit(&b, (s%secondsPerDay)/secondsPerHour, "h")
	case s >= secondsPerHour:
		writeUnit(&b, s/secondsPerHour, "h")
		writeUnit(&b, (s%secondsPerHour)/secondsPerMin, "m")
	case s >= secondsPerMin:
		writeUnit(&b, s/secondsPerMin, "m")
		writeUnit(&b, s%secondsPerMin, "s")
	default:
		writeUnit(&b, s, "s")
	}

	return b.String()
}

func writeUnit(b *strings.Builder, v int64, unit string) {
	if v == 0 && b.Len() > 0 {
		return
	}

	b.WriteString(strconv.FormatInt(v, 10))
	b.WriteString(unit)
}

// Join renders a list cell. An empty list is Absent rather than an empty cell.
func Join(items []string, sep string) string {
	if len(items) == 0 {
		return Absent
	}

	return strings.Join(items, sep)
}

// SecondsSince converts a controller timestamp into an age in seconds. A zero time
// and the Unix epoch both mean the controller reported no instant: this estate
// returns the epoch on several sibling timestamp leaves, and an age computed from
// it would read as fifty-six years.
func SecondsSince(now, t time.Time) *int64 {
	if t.IsZero() || t.Unix() <= 0 {
		return nil
	}

	s := int64(now.Sub(t).Seconds())
	if s < 0 {
		return nil
	}

	return &s
}
