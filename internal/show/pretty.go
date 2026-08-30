package show

import "github.com/umatare5/wnc/internal/render"

// Glyphs for the bordered table. Each must be one code point whose width does not move with the
// reader's locale: a U+FE0F selector or an East Asian ambiguous character asks for a two-column
// rendering tablewriter's measurement does not see, and every bordered row past that cell drifts.
const (
	glyphOK   = "✅"
	glyphBad  = "❌"
	glyphOff  = "⬜"
	glyphNo   = "✕"
	glyphWarn = "⚠"
)

// prettyBool maps an optional boolean onto two glyphs. A nil stays render.Absent, because a glyph
// there would assert a reading nothing gave.
func prettyBool(p *bool, whenTrue, whenFalse string) string {
	if p == nil {
		return render.Absent
	}

	if *p {
		return whenTrue
	}

	return whenFalse
}

// prettyState maps a two-member enum's display string onto the healthy glyph and one
// caller-chosen negative glyph. A spelling neither branch names passes through as text, because
// folding it into either glyph would report a state a future release invented as one of these two.
func prettyState(p *string, on, off, offGlyph string) string {
	if p == nil {
		return render.Absent
	}

	switch *p {
	case on:
		return glyphOK
	case off:
		return offGlyph
	default:
		return render.Str(*p)
	}
}

// prettyOtherwise maps one display string onto the healthy glyph and folds every other reported
// value onto a single glyph. show ap's State column takes this because only Registered is the
// serving member of its six; an unreported value stays render.Absent.
func prettyOtherwise(p *string, healthy, otherGlyph string) string {
	if p == nil || *p == "" {
		return render.Absent
	}

	if *p == healthy {
		return glyphOK
	}

	return otherGlyph
}
