package config

import (
	"time"
)

// Dur decodes a Go duration string such as "30s", and going through TextUnmarshaler is also what
// rejects a bare integer a numeric decode would read as nanoseconds. A json format tag is no way
// out: encoding/json/v2 shipped without `format`, so a field carrying one fails the decode of its
// whole struct whatever the field's type.
type Dur time.Duration

// UnmarshalText decodes the duration. A non-string JSON value never reaches here, so a numeric
// form cannot be tolerated in this method.
func (d *Dur) UnmarshalText(b []byte) error {
	v, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}

	*d = Dur(v)

	return nil
}

func (d Dur) Duration() time.Duration {
	return time.Duration(d)
}
