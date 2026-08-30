package show

import (
	"reflect"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/render"
)

// A column is declared in three places: the key list --sort-by validates against, the Column list
// both renderers walk, and the row struct's json tags. This asserts the three agree, in order,
// which is what catches a copied struct line json/v2 would reject at run time.
func assertColumns[R any](t *testing.T, keys []string, cols []render.Column[R]) {
	t.Helper()

	if got := strings.Join(render.Keys(cols), ","); got != strings.Join(keys, ",") {
		t.Errorf("column keys\n  %s\nkey list\n  %s", got, strings.Join(keys, ","))
	}

	tags, opts := jsonTags(reflect.TypeFor[R]())

	if got := strings.Join(tags, ","); got != strings.Join(keys, ",") {
		t.Errorf("json tags\n  %s\nkey list\n  %s", got, strings.Join(keys, ","))
	}

	for _, c := range cols {
		if c.Cell == nil {
			t.Errorf("column %q has no Cell function", c.Key)
		}

		if c.Header == "" {
			t.Errorf("column %q has no header", c.Key)
		}
	}

	assertTagOptions(t, reflect.TypeFor[R](), opts)
}

// assertTagOptions enforces the row-struct rules. omitempty is banned outright: it
// drops a reported zero, a reported empty string and a reported false, each of which
// is a reading. omitzero is allowed only on a pointer, where the zero value is nil
// and therefore genuinely means "not reported". The json format option is banned
// because every one of its values is rejected at run time, after passing both the
// compiler and the linter.
func assertTagOptions(t *testing.T, typ reflect.Type, opts map[string][]string) {
	t.Helper()

	for i := range typ.NumField() {
		f := typ.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")

		for _, o := range opts[name] {
			switch {
			case o == "omitempty":
				t.Errorf("field %s carries omitempty; use omitzero on a pointer", f.Name)
			case o == "omitzero" && f.Type.Kind() != reflect.Pointer:
				t.Errorf("field %s carries omitzero but is %s, not a pointer", f.Name, f.Type.Kind())
			case strings.HasPrefix(o, "format:"):
				t.Errorf("field %s carries a json format option, which fails at run time", f.Name)
			}
		}
	}
}

// jsonTags reads the json field names and their options, in declaration order.
func jsonTags(typ reflect.Type) (names []string, opts map[string][]string) {
	opts = make(map[string][]string, typ.NumField())

	for i := range typ.NumField() {
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		name, rest, _ := strings.Cut(tag, ",")
		names = append(names, name)

		if rest != "" {
			opts[name] = strings.Split(rest, ",")
		}
	}

	return names, opts
}

func TestAPTagColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, APTagKeys(), APTagColumns())
}

func TestAPJoinColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, APJoinKeys(), APJoinColumns())
}

func TestAPColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, APKeys(), APColumns())
}

func TestOverviewColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, OverviewKeys(), OverviewColumns())
}

func TestClientColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, ClientKeys(), ClientColumns())
}

func TestWLANColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, WLANKeys(), WLANColumns())
}

func TestPolicyTagColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, PolicyTagKeys(), PolicyTagColumns())
}

func TestSiteTagColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, SiteTagKeys(), SiteTagColumns())
}

func TestRFTagColumnsMatch(t *testing.T) {
	t.Parallel()

	assertColumns(t, RFTagKeys(), RFTagColumns())
}

// Every command's default sort key must be one of its own columns, or the first run
// of that command fails on a key the flag itself offered.
func TestDefaultSortKeysExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		keys   []string
		defkey string
	}{
		{name: "overview", keys: OverviewKeys(), defkey: DefaultSortAPName},
		{name: "ap", keys: APKeys(), defkey: DefaultSortAPName},
		{name: "ap-join", keys: APJoinKeys(), defkey: DefaultSortAPName},
		{name: "ap-tag", keys: APTagKeys(), defkey: DefaultSortAPName},
		{name: "client", keys: ClientKeys(), defkey: DefaultSortMAC},
		{name: "wlan", keys: WLANKeys(), defkey: DefaultSortWLANID},
		{name: "policy-tag", keys: PolicyTagKeys(), defkey: DefaultSortPolicyTag},
		{name: "site-tag", keys: SiteTagKeys(), defkey: DefaultSortSiteTag},
		{name: "rf-tag", keys: RFTagKeys(), defkey: DefaultSortRFTag},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !contains(tt.keys, tt.defkey) {
				t.Errorf("default sort key %q is not a column of %s", tt.defkey, tt.name)
			}
		})
	}
}

// Two columns of one command sharing a key would make --sort-by ambiguous, and json/v2 fails
// the marshal outright on two struct fields sharing a json name.
func TestKeysAreUnique(t *testing.T) {
	t.Parallel()

	lists := map[string][]string{
		"overview":   OverviewKeys(),
		"ap":         APKeys(),
		"ap-join":    APJoinKeys(),
		"ap-tag":     APTagKeys(),
		"client":     ClientKeys(),
		"wlan":       WLANKeys(),
		"policy-tag": PolicyTagKeys(),
		"site-tag":   SiteTagKeys(),
		"rf-tag":     RFTagKeys(),
	}

	for name, keys := range lists {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			seen := make(map[string]bool, len(keys))
			for _, k := range keys {
				if seen[k] {
					t.Errorf("key %q appears twice", k)
				}

				seen[k] = true
			}
		})
	}
}

// Every command carries the controller column: two controllers in one invocation can
// host access points and WLANs with colliding names, so a row is only identifiable
// with it.
func TestEveryCommandCarriesTheController(t *testing.T) {
	t.Parallel()

	for name, keys := range map[string][]string{
		"overview":   OverviewKeys(),
		"ap":         APKeys(),
		"ap-join":    APJoinKeys(),
		"ap-tag":     APTagKeys(),
		"client":     ClientKeys(),
		"wlan":       WLANKeys(),
		"policy-tag": PolicyTagKeys(),
		"site-tag":   SiteTagKeys(),
		"rf-tag":     RFTagKeys(),
	} {
		if !contains(keys, "controller") {
			t.Errorf("%s has no controller column", name)
		}
	}
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}

	return false
}
