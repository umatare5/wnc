package show

// optional turns a leaf into an absence-preserving field. Almost every string on this SDK's read
// paths sits in a non-pointer struct, so an omitted container and an omitted leaf both arrive as
// "" and there is no pointer left to test.
func optional(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// optionalStr collapses an optional string leaf the controller sent empty. A few of the SDK's
// configuration structs do declare a pointer here, and a pointer to "" would keep the JSON key
// while the table cell renders Absent.
func optionalStr(p *string) *string {
	if p == nil {
		return nil
	}

	return optional(*p)
}

// zeroAbsent turns a numeric leaf whose zero cannot be a real reading into an absence-preserving
// field. Use it only where zero is impossible: a joined access point has at least one radio slot,
// a serving radio has a channel, and no controller reports an RSSI of exactly 0 dBm. Never where
// zero is a reading, such as a client count or an SNR.
func zeroAbsent[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](
	v T,
) *T {
	if v == 0 {
		return nil
	}

	return &v
}

// ptr takes the address of a reported value, for a leaf whose zero is a reading.
func ptr[T any](v T) *T {
	return &v
}
