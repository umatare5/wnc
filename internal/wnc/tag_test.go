package wnc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// exchange is one request the recorder saw.
type exchange struct {
	method string
	path   string
	body   string
}

// recorder answers every request with the given status and body and keeps what arrived. A
// write's correctness is in the request rather than in the answer: the controller reports a
// rejected payload as a bare 400, so the payload has to be asserted here instead.
type recorder struct {
	client *Client

	mu   sync.Mutex
	seen []exchange
}

func newRecorder(t *testing.T, status int, answer string) *recorder {
	t.Helper()

	r := &recorder{}

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Errorf("reading the request body: %v", err)
		}

		r.mu.Lock()
		// EscapedPath and not Path: net/url decodes a percent escape into Path, so the
		// escaping this asserts on would be invisible there.
		r.seen = append(r.seen, exchange{method: req.Method, path: req.URL.EscapedPath(), body: string(body)})
		r.mu.Unlock()

		if answer != "" {
			w.Header().Set("Content-Type", "application/yang-data+json")
		}

		w.WriteHeader(status)

		if answer != "" {
			_, _ = w.Write([]byte(answer))
		}
	}))
	srv.StartTLS()

	logger, err := log.NewWithOutput(&discard{}, "error")
	if err != nil {
		t.Fatalf("building the logger: %v", err)
	}

	c, err := NewClient(
		config.Target{Name: "test", Host: srv.Listener.Addr().String(), Token: fakeToken},
		config.Settings{Timeout: 10 * time.Second, Insecure: true},
		logger, "wnc/test",
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	r.client = c

	return r
}

func (r *recorder) last(t *testing.T) exchange {
	t.Helper()

	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.seen) == 0 {
		t.Fatal("no request was made")
	}

	return r.seen[len(r.seen)-1]
}

// all returns every exchange, for a call that is more than one request.
func (r *recorder) all() []exchange {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]exchange(nil), r.seen...)
}

// A create carries the key leaf plus the fields that were named, and nothing else. The
// controller accepts a body of the key alone — nothing in the three tag groupings is
// mandatory — so an unnamed field must be absent rather than sent as an empty value.
func TestTagCreatePayloads(t *testing.T) {
	t.Parallel()

	desc := "written by a test"
	profile := "lab-profile"

	tests := []struct {
		name   string
		call   func(*Client) error
		want   string
		method string
	}{
		{
			name:   "policy tag, key alone",
			call:   func(c *Client) error { return c.CreatePolicyTag(t.Context(), "t1", TagFields{}) },
			want:   `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entry":{"tag-name":"t1"}}`,
			method: http.MethodPost,
		},
		{
			name: "policy tag with a WLAN binding",
			call: func(c *Client) error {
				return c.CreatePolicyTag(t.Context(), "t1", TagFields{
					Description: &desc, WLAN: &profile, PolicyProfile: &profile,
				})
			},
			want: `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entry":{"tag-name":"t1",` +
				`"description":"written by a test","wlan-policies":{"wlan-policy":` +
				`[{"wlan-profile-name":"lab-profile","policy-profile-name":"lab-profile"}]}}}`,
			method: http.MethodPost,
		},
		{
			name:   "site tag, key alone",
			call:   func(c *Client) error { return c.CreateSiteTag(t.Context(), "t2", TagFields{}) },
			want:   `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-config":{"site-tag-name":"t2"}}`,
			method: http.MethodPost,
		},
		{
			name: "site tag with an AP join profile",
			call: func(c *Client) error {
				return c.CreateSiteTag(t.Context(), "t2", TagFields{APJoinProfile: &profile})
			},
			want: `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-config":` +
				`{"site-tag-name":"t2","ap-join-profile":"lab-profile"}}`,
			method: http.MethodPost,
		},
		{
			// The flex-profile leaf declares when "../is-local-site = 'false'" and
			// is-local-site defaults to TRUE, so a create naming a flex profile must send
			// the flag false or the controller answers 400 on the when — measured on 17.12.8.
			name: "site tag with a flex profile carries is-local-site false",
			call: func(c *Client) error {
				return c.CreateSiteTag(t.Context(), "t2", TagFields{FlexProfile: &profile})
			},
			want: `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-config":` +
				`{"site-tag-name":"t2","flex-profile":"lab-profile","is-local-site":false}}`,
			method: http.MethodPost,
		},
		{
			name:   "RF tag, key alone",
			call:   func(c *Client) error { return c.CreateRFTag(t.Context(), "t3", TagFields{}) },
			want:   `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":{"tag-name":"t3"}}`,
			method: http.MethodPost,
		},
		{
			name: "RF tag with a 5 GHz profile",
			call: func(c *Client) error {
				return c.CreateRFTag(t.Context(), "t3", TagFields{Profile5GHz: &profile})
			},
			want: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":` +
				`{"tag-name":"t3","dot11a-rf-profile-name":"lab-profile"}}`,
			method: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusCreated, "")

			if err := tt.call(r.client); err != nil {
				t.Fatalf("create: %v", err)
			}

			got := r.last(t)
			if got.method != tt.method {
				t.Errorf("method = %s, want %s", got.method, tt.method)
			}

			if got.body != tt.want {
				t.Errorf("payload =\n  %s\nwant\n  %s", got.body, tt.want)
			}
		})
	}
}

// An RF tag write must never carry rf-tag-radio-profiles: a body naming that container with a null
// list answers 400 "invalid value for: rf-tag-radio-profile". The SDK declares the field as a
// pointer with omitempty, so this guards the dependency rather than a hand-built body.
func TestRFTagWriteOmitsTheRadioProfileContainer(t *testing.T) {
	t.Parallel()

	desc := "written by a test"
	existing := `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":[{"tag-name":"t3"}]}`

	for _, tt := range []struct {
		name   string
		answer string
		call   func(*Client) error
	}{
		{name: "create", call: func(c *Client) error {
			return c.CreateRFTag(t.Context(), "t3", TagFields{Description: &desc})
		}},
		// The update reads the tag before writing it back, so the read must answer.
		{name: "update", answer: existing, call: func(c *Client) error {
			return c.UpdateRFTag(t.Context(), "t3", TagFields{Description: &desc})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			status := http.StatusNoContent
			if tt.answer != "" {
				status = http.StatusOK
			}

			r := newRecorder(t, status, tt.answer)

			if err := tt.call(r.client); err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}

			if body := r.last(t).body; strings.Contains(body, "rf-tag-radio-profile") {
				t.Errorf("payload names the radio-profile container: %s", body)
			}
		})
	}
}

// An RF tag update is a read of the record followed by a merge PATCH of it. The merge is
// what carries a field the command did not name, and the read is the SDK's own — it is why
// this kind still takes two requests where the site and policy kinds take one.
func TestRFTagUpdateReadsThenMergesOnTheKeyedURL(t *testing.T) {
	t.Parallel()

	profile := "lab-profile"
	r := newRecorder(t, http.StatusOK,
		`{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":[{"tag-name":"has space","description":"kept"}]}`)

	if err := r.client.UpdateRFTag(t.Context(), "has space", TagFields{Profile24GHz: &profile}); err != nil {
		t.Fatalf("UpdateRFTag: %v", err)
	}

	if n := len(r.all()); n != 2 {
		t.Fatalf("exchanges = %d, want a read then a write", n)
	}

	if first := r.all()[0].method; first != http.MethodGet {
		t.Errorf("first method = %s, want %s", first, http.MethodGet)
	}

	got := r.last(t)
	if got.method != http.MethodPatch {
		t.Errorf("method = %s, want %s", got.method, http.MethodPatch)
	}

	// The space is escaped rather than sent raw, or the controller reads a different node.
	if !strings.HasSuffix(got.path, "/rf-tags/rf-tag=has%20space") {
		t.Errorf("path = %s, want a percent-escaped list key", got.path)
	}

	// The description was not named and comes back from the read, so the PATCH carries it: the
	// SDK writes the whole record back rather than only the named field.
	if !strings.Contains(got.body, `"description":"kept"`) ||
		!strings.Contains(got.body, `"dot11b-rf-profile-name":"lab-profile"`) {
		t.Errorf("payload dropped an unnamed field or the named one: %s", got.body)
	}
}

// A name the controller does not hold answers 404 on the keyed URL, and that is an absence
// rather than a failure. Every other status stays a failure, or a delete would be sent at a
// controller that could not answer for it.
func TestTagExistsSeparatesAbsenceFromFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		answer     string
		wantExists bool
		wantErr    bool
	}{
		{
			name: "present", status: http.StatusOK,
			answer:     `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":[{"tag-name":"t3"}]}`,
			wantExists: true,
		},
		{name: "absent", status: http.StatusNotFound},
		{name: "unauthorized", status: http.StatusUnauthorized, wantErr: true},
		{name: "server fault", status: http.StatusInternalServerError, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, tt.status, tt.answer)

			got, err := r.client.RFTagExists(t.Context(), "t3")

			if tt.wantErr {
				if err == nil {
					t.Fatal("a failure was reported as an absence")
				}

				return
			}

			if err != nil {
				t.Fatalf("RFTagExists: %v", err)
			}

			if got != tt.wantExists {
				t.Errorf("exists = %v, want %v", got, tt.wantExists)
			}
		})
	}
}

func TestTagDeletesUseTheKeyedURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*Client) error
		want string
	}{
		{
			name: "policy tag", want: "/policy-list-entries/policy-list-entry=t1",
			call: func(c *Client) error { return c.DeletePolicyTag(t.Context(), "t1") },
		},
		{
			name: "site tag", want: "/site-tag-configs/site-tag-config=t2",
			call: func(c *Client) error { return c.DeleteSiteTag(t.Context(), "t2") },
		},
		{
			name: "RF tag", want: "/rf-tags/rf-tag=t3",
			call: func(c *Client) error { return c.DeleteRFTag(t.Context(), "t3") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusNoContent, "")

			if err := tt.call(r.client); err != nil {
				t.Fatalf("delete: %v", err)
			}

			got := r.last(t)
			if got.method != http.MethodDelete {
				t.Errorf("method = %s, want %s", got.method, http.MethodDelete)
			}

			if !strings.HasSuffix(got.path, tt.want) {
				t.Errorf("path = %s, want it to end with %s", got.path, tt.want)
			}

			if got.body != "" {
				t.Errorf("a delete carried a body: %q", got.body)
			}
		})
	}
}

// An update naming nothing must make no request at all: the CLI reports that there is
// nothing to change, and a request would contradict the report.
func TestTagUpdatesWithNoFieldSendNothing(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		call func(*Client) error
	}{
		{name: "policy tag", call: func(c *Client) error {
			return c.UpdatePolicyTag(t.Context(), "t1", TagFields{})
		}},
		{name: "site tag", call: func(c *Client) error {
			return c.UpdateSiteTag(t.Context(), "t2", TagFields{})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusNoContent, "")

			if err := tt.call(r.client); err != nil {
				t.Fatalf("update: %v", err)
			}

			r.mu.Lock()
			defer r.mu.Unlock()

			if len(r.seen) != 0 {
				t.Errorf("an update with no field made %d requests", len(r.seen))
			}
		})
	}
}

func TestTagFieldsEmpty(t *testing.T) {
	t.Parallel()

	if !(TagFields{}).Empty() {
		t.Error("a zero TagFields is not reported as empty")
	}

	v := "x"
	yes := true

	for name, f := range map[string]TagFields{
		"description":     {Description: &v},
		"2.4 GHz profile": {Profile24GHz: &v},
		"5 GHz profile":   {Profile5GHz: &v},
		"6 GHz profile":   {Profile6GHz: &v},
		"AP join profile": {APJoinProfile: &v},
		"flex profile":    {FlexProfile: &v},
		"local site":      {LocalSite: &yes},
		"WLAN":            {WLAN: &v},
		"policy profile":  {PolicyProfile: &v},
	} {
		if f.Empty() {
			t.Errorf("a TagFields carrying only the %s is reported as empty", name)
		}
	}
}
