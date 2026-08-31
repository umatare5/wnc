package wnc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// The paths the harness routes on. Both RPC paths are spelt out because the CLI declares neither
// and the SDK's route constants are internal, so the operation each request must land on is
// asserted rather than restated.
const (
	apResetRPC     = "Cisco-IOS-XE-wireless-access-point-cmd-rpc:ap-reset"
	capwapResetRPC = "Cisco-IOS-XE-wireless-access-point-cmd-rpc:set-rad-capwap-reset"
	nameAP1        = "TEST-AP01"
	nameAP2        = "TEST-AP02"
	keyedName1     = "ap-name-mac-map=" + nameAP1
	keyedName2     = "ap-name-mac-map=" + nameAP2
	nameMapRow     = `{"Cisco-IOS-XE-wireless-access-point-oper:ap-name-mac-map":[` +
		`{"wtp-name":"` + nameAP1 + `","wtp-mac":"` + macAP1 + `","eth-mac":"00:00:5e:00:53:11"}]}`
	emptyNameMap = `{"Cisco-IOS-XE-wireless-access-point-oper:ap-name-mac-map":[]}`
)

// The row carries all three leaves, and only the radio address has a consumer: the keyed radio
// read is the one read this CLI makes that is not keyed on the name.
func TestAPRadioMACByName(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{keyedName1: {body: nameMapRow}})

	mac, found, err := c.APRadioMACByName(t.Context(), nameAP1)
	if err != nil {
		t.Fatalf("APRadioMACByName: %v", err)
	}

	if !found {
		t.Fatal("a row that came back was reported as an absence")
	}

	if mac != macAP1 {
		t.Errorf("mac = %q, want %q", mac, macAP1)
	}
}

// A controller holding no access point under the name answers 404, and the caller has to be
// able to tell that apart from a controller that failed.
func TestAPRadioMACByNameReportsAnAbsentNameAsNotFound(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{keyedName2: {status: http.StatusNotFound}})

	if _, _, err := c.APRadioMACByName(t.Context(), nameAP2); err == nil {
		t.Fatal("a 404 did not produce an error")
	} else if cause, status := Classify(err); cause != CauseNotFound || status != http.StatusNotFound {
		t.Errorf("Classify = %q/%d, want not-found/404", cause, status)
	}
}

// The bool is the whole reason this returns three values. A 200 carrying no row is not a
// measured shape for this list, so it must not decode to an access point with an empty
// address: that address is what the keyed radio read is keyed on.
func TestAPRadioMACByNameReportsAnAnswerWithNoRowAsAnAbsence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		served answer
	}{
		{name: "an empty list", served: answer{body: emptyNameMap}},
		{name: "a 204 with no body", served: answer{status: http.StatusNoContent}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{keyedName1: tt.served})

			mac, found, err := c.APRadioMACByName(t.Context(), nameAP1)
			if err != nil {
				t.Fatalf("an absence produced an error: %v", err)
			}

			if found {
				t.Error("an answer with no row was reported as an access point")
			}

			if mac != "" {
				t.Errorf("mac = %q, want empty", mac)
			}
		})
	}
}

// The read is deliberately unpruned: the row is three leaves, so a fields expression buys
// nothing, and one naming a node a release does not declare answers 200 with a body that stops
// mid-object. Nothing in the answer shows that, so the query is asserted.
func TestAPRadioMACByNameReadsTheWholeRowRatherThanPruningIt(t *testing.T) {
	t.Parallel()

	var got string

	c := newClientWithQuery(t, &got, nameMapRow)

	if _, _, err := c.APRadioMACByName(t.Context(), nameAP1); err != nil {
		t.Fatalf("APRadioMACByName: %v", err)
	}

	if got != "" {
		t.Errorf("query = %q, want none", got)
	}
}

// The RPC's input is a mandatory choice and the controller answers 400 when neither arm is sent,
// but nothing in the response distinguishes a well-formed payload from a lucky one, so the body is
// asserted. The name arm is what this posts, where the SDK's other path resolves an address first.
func TestResetAPSendsTheNameArm(t *testing.T) {
	t.Parallel()

	r := newRecorder(t, http.StatusNoContent, "")

	if err := r.client.ResetAP(t.Context(), nameAP1); err != nil {
		t.Fatalf("ResetAP: %v", err)
	}

	if n := len(r.all()); n != 1 {
		t.Errorf("exchanges = %d, want the write alone", n)
	}

	got := r.last(t)
	if !contains(got.path, apResetRPC) {
		t.Errorf("path = %s, want the reset operation", got.path)
	}

	if want := `{"input":{"ap-name":"` + nameAP1 + `"}}`; got.body != want {
		t.Errorf("payload = %s, want %s", got.body, want)
	}
}

// A 204 with no body is the whole success signal, because the RPC declares no output. The
// harness fails the test if an unrouted path is requested, which is what proves the RPC is
// reached without a collection read in front of it.
func TestResetAP(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{apResetRPC: {status: http.StatusNoContent}})

	if err := c.ResetAP(t.Context(), nameAP1); err != nil {
		t.Errorf("ResetAP: %v", err)
	}
}

// The RPC's input is a mandatory choice, and the controller answers 400 with
// "must provide one of: ap-name, mac-addr" when neither arm is sent. Nothing in the response
// distinguishes a well-formed payload from a lucky one, so the body is asserted.
func TestResetCAPWAPSendsTheNameArm(t *testing.T) {
	t.Parallel()

	var got string

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Errorf("reading the request body: %v", err)
		}

		got = string(body)

		w.WriteHeader(http.StatusNoContent)
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

	if err := c.ResetCAPWAP(t.Context(), nameAP1); err != nil {
		t.Fatalf("ResetCAPWAP: %v", err)
	}

	if want := `{"input":{"ap-name":"` + nameAP1 + `"}}`; got != want {
		t.Errorf("payload = %s, want %s", got, want)
	}
}

// A 204 with no body is the whole success signal, because the RPC declares no output.
func TestResetCAPWAP(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{capwapResetRPC: {status: http.StatusNoContent}})

	if err := c.ResetCAPWAP(t.Context(), nameAP1); err != nil {
		t.Errorf("ResetCAPWAP: %v", err)
	}
}

func TestResetCAPWAPReportsTheControllerRefusal(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{capwapResetRPC: {status: http.StatusBadRequest}})

	err := c.ResetCAPWAP(t.Context(), nameAP1)
	if err == nil {
		t.Fatal("a 400 did not produce an error")
	}

	if cause, status := Classify(err); cause != CauseHTTP || status != http.StatusBadRequest {
		t.Errorf("Classify = %q/%d, want http/400", cause, status)
	}
}

func TestResetAPReportsTheControllerRefusal(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{apResetRPC: {status: http.StatusBadRequest}})

	err := c.ResetAP(t.Context(), nameAP1)
	if err == nil {
		t.Fatal("a 400 did not produce an error")
	}

	if cause, status := Classify(err); cause != CauseHTTP || status != http.StatusBadRequest {
		t.Errorf("Classify = %q/%d, want http/400", cause, status)
	}
}
