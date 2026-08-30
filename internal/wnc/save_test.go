package wnc

import (
	"errors"
	"net/http"
	"testing"
)

// The envelope key and the result string are the ones a controller sent; only the whitespace is
// invented, and JSON reads the device's pretty-printed answer and this compact one alike.
const (
	savedResult = "Save running-config successful"
	saveReply   = `{"cisco-ia:output":{"result":"` + savedResult + `"}}`

	// The operation's own name. The SDK's route constants are internal, so the path the request
	// must land on is asserted here rather than restated from the call.
	saveConfigRPC = "cisco-ia:save-config"
)

// The RPC is the only one this CLI posts whose schema declares an output container, so the
// controller's account of the save is a value this layer carries rather than a 204 it infers
// from.
func TestSaveConfigReturnsTheControllersAccount(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{saveConfigRPC: {body: saveReply}})

	got, err := c.SaveConfig(t.Context())
	if err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	if got != savedResult {
		t.Errorf("result = %q, want %q", got, savedResult)
	}
}

// The RPC takes no input at all, and the SDK's nil-payload path sends a bodiless POST. An
// "input" object would be a payload the schema does not declare, and nothing in a successful
// answer would show it, so the request is asserted.
func TestSaveConfigPostsNoBody(t *testing.T) {
	t.Parallel()

	r := newRecorder(t, http.StatusOK, saveReply)

	if _, err := r.client.SaveConfig(t.Context()); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	if n := len(r.all()); n != 1 {
		t.Errorf("exchanges = %d, want the write alone", n)
	}

	got := r.last(t)
	if got.method != http.MethodPost {
		t.Errorf("method = %s, want POST", got.method)
	}

	if !contains(got.path, saveConfigRPC) {
		t.Errorf("path = %s, want the save operation", got.path)
	}

	if got.body != "" {
		t.Errorf("payload = %s, want none", got.body)
	}
}

// A save whose outcome the controller did not report must not be reported as a save. The
// result string is the controller's whole account of it, so an answer carrying none leaves
// this layer with nothing to stand on.
func TestSaveConfigRefusesAnAnswerWithNoResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{name: "no envelope at all", body: `{}`},
		{name: "the envelope with no result", body: `{"cisco-ia:output":{}}`},
		{name: "an empty result", body: `{"cisco-ia:output":{"result":""}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{saveConfigRPC: {body: tt.body}})

			got, err := c.SaveConfig(t.Context())
			if !errors.Is(err, ErrSaveNotReported) {
				t.Errorf("err = %v, want ErrSaveNotReported", err)
			}

			if got != "" {
				t.Errorf("result = %q, want empty", got)
			}
		})
	}
}

// A 204 is not a measured shape for this RPC, and it must not read as a save either. The SDK
// decodes an empty body to a nil Output with no error, so it is refused by the same guard as an
// empty result and kept as its own case because the wire shape differs.
func TestSaveConfigRefusesAnAnswerWithNoBody(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{saveConfigRPC: {status: http.StatusNoContent}})

	if _, err := c.SaveConfig(t.Context()); !errors.Is(err, ErrSaveNotReported) {
		t.Errorf("err = %v, want ErrSaveNotReported", err)
	}
}

// An unknown result passes through as the controller sent it, which is the rule the display
// tables follow for an enum spelling. A release wording the success differently must not turn
// a save that worked into a failure.
func TestSaveConfigPassesAnUnknownResultThrough(t *testing.T) {
	t.Parallel()

	const other = "Configuration saved"

	c := newClient(t, routes{saveConfigRPC: {body: `{"cisco-ia:output":{"result":"` + other + `"}}`}})

	got, err := c.SaveConfig(t.Context())
	if err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}

	if got != other {
		t.Errorf("result = %q, want %q", got, other)
	}
}

func TestSaveConfigReportsTheControllerRefusal(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{saveConfigRPC: {status: http.StatusBadRequest}})

	_, err := c.SaveConfig(t.Context())
	if err == nil {
		t.Fatal("a 400 did not produce an error")
	}

	if cause, status := Classify(err); cause != CauseHTTP || status != http.StatusBadRequest {
		t.Errorf("Classify = %q/%d, want http/400", cause, status)
	}
}
