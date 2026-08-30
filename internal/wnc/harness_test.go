package wnc

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// fakeToken looks like a Basic auth token and decodes to nothing real.
const fakeToken = "TestToken0123456789ABCDEF=="

// routes maps the last element of a RESTCONF path to the body the controller should
// answer with. A status may be given instead, for the failure paths.
type routes map[string]answer

// answer is one canned response. An empty body with a 200 stands for the 204 an empty
// collection produces, which the SDK turns into a zero value and a nil error.
type answer struct {
	status int
	body   string

	// query records this route's raw query when set. A stub cannot reproduce the
	// truncated 200 an undeclared fields node answers with, so a prune is asserted on
	// the request instead.
	query *string
}

// newClient starts a TLS server serving the given routes and returns a client aimed at it.
func newClient(t *testing.T, r routes) *Client {
	t.Helper()

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		a, ok := r[path.Base(req.URL.Path)]
		if !ok {
			t.Errorf("unexpected request for %s", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if a.query != nil {
			*a.query = req.URL.RawQuery
		}

		if a.status != 0 && a.status != http.StatusOK {
			w.WriteHeader(a.status)

			return
		}

		w.Header().Set("Content-Type", "application/yang-data+json")

		if a.body != "" {
			_, _ = w.Write([]byte(a.body))
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
		logger,
		"wnc/test",
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	return c
}

// newClientWithQuery serves one body for every path and records the RAW query, not url.Values: a
// fields expression is built from semicolons, which net/url has refused to parse since Go 1.17, so
// Query() would silently drop the parameter this asserts on.
func newClientWithQuery(t *testing.T, query *string, body string) *Client {
	t.Helper()

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*query = req.URL.RawQuery

		w.Header().Set("Content-Type", "application/yang-data+json")
		_, _ = w.Write([]byte(body))
	}))
	srv.StartTLS()

	logger, err := log.NewWithOutput(&discard{}, "error")
	if err != nil {
		t.Fatalf("building the logger: %v", err)
	}

	c, err := NewClient(
		config.Target{Name: "test", Host: srv.Listener.Addr().String(), Token: fakeToken},
		config.Settings{Timeout: 10 * time.Second, Insecure: true},
		logger,
		"wnc/test",
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	return c
}

// discard swallows the logger's output.
type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }

// The SDK pins its own DialContext, so no fake transport can be injected and the test
// server has to be a real TLS listener the client can reach. This asserts that much
// before any fixture-driven test relies on it.
func TestHarnessIsReachable(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[]}`}})

	tags, err := c.APTags(t.Context())
	if err != nil {
		t.Fatalf("APTags: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("got %d rows from an empty collection", len(tags))
	}
}

// logForTest builds a quiet logger.
func logForTest(t *testing.T) (*logrus.Logger, error) {
	t.Helper()

	return log.NewWithOutput(&discard{}, "error")
}

// newClientFor builds a client without a server, for the construction checks.
func newClientFor(t *testing.T, host, token string, logger *logrus.Logger) (*Client, error) {
	t.Helper()

	return NewClient(
		config.Target{Name: "test", Host: host, Token: token},
		config.Settings{Timeout: time.Second, Insecure: true},
		logger,
		"wnc/test",
	)
}

// newClientWithTimeout builds a client with an explicit timeout.
func newClientWithTimeout(t *testing.T, logger *logrus.Logger, d time.Duration) (*Client, error) {
	t.Helper()

	return NewClient(
		config.Target{Name: "test", Host: "192.0.2.1", Token: fakeToken},
		config.Settings{Timeout: d, Insecure: true},
		logger,
		"wnc/test",
	)
}

// contextWithDeadlinePassed returns a context whose deadline is already behind it, so
// the very first request fails as a timeout.
func contextWithDeadlinePassed(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	return context.WithTimeout(t.Context(), time.Nanosecond)
}

// marshalOf renders a value so a test can assert what is not in it.
func marshalOf(t *testing.T, v any) string {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	return string(b)
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
