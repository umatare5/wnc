package wnc

import (
	"net/http"
	"testing"
)

// The client address is an RFC 7042 documentation MAC, which no interface can hold, and the route
// is keyed the way the controller keys the list. rowClientMAC is the same address in the other
// case, so a read that echoed its argument instead of returning the row cannot pass — every
// client-mac a controller serves is lowercase, which is the fixture's own device here.
const (
	macClient    = "00:00:5e:00:53:a1"
	rowClientMAC = "00:00:5E:00:53:A1"
	keyedClient  = "common-oper-data=" + macClient
	clientRow    = `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[` +
		`{"client-mac":"` + rowClientMAC + `","ap-name":"` + nameAP1 + `","co-state":"client-status-run"}]}`
	emptyClients = `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[]}`
)

// deauthRPC is the operation's own name. The SDK's route constants are internal, so the route a
// request must land on is asserted here rather than restated from the call.
const deauthRPC = "Cisco-IOS-XE-wireless-client-rpc:apf-ms-delete-all"

// The username two clients share, and one no client carries. A username is a bare string in the
// schema, so a space in it is legal and is here on purpose.
const (
	sharedUsername = "test user"
	absentUsername = "test-nobody"
	usernameRows   = `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[` +
		`{"client-mac":"` + rowClientMAC + `","username":"` + sharedUsername + `"},` +
		`{"client-mac":"00:00:5E:00:53:A2","username":"` + sharedUsername + `"},` +
		`{"client-mac":"00:00:5E:00:53:A3","username":"test-other"},` +
		`{"client-mac":"00:00:5E:00:53:A4","username":""}]}`
)

// The address the row carries is what goes on the wire, so the read has to return it rather
// than echo what the caller passed.
func TestClientByMACReturnsTheRowsAddress(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{keyedClient: {body: clientRow}})

	mac, found, err := c.ClientByMAC(t.Context(), macClient)
	if err != nil {
		t.Fatalf("ClientByMAC: %v", err)
	}

	if !found {
		t.Fatal("a row that came back was reported as an absence")
	}

	if mac != rowClientMAC {
		t.Errorf("mac = %q, want the row's %q", mac, rowClientMAC)
	}
}

// A controller holding no client at the address answers 404, and the caller has to tell that
// apart from a controller that failed.
func TestClientByMACReportsAnAbsentClientAsNotFound(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{keyedClient: {status: http.StatusNotFound}})

	if _, _, err := c.ClientByMAC(t.Context(), macClient); err == nil {
		t.Fatal("a 404 did not produce an error")
	} else if cause, status := Classify(err); cause != CauseNotFound || status != http.StatusNotFound {
		t.Errorf("Classify = %q/%d, want not-found/404", cause, status)
	}
}

// A 200 carrying no row is an absence and not a client with an empty address, which would go
// on the wire as the deauthentication's target.
func TestClientByMACReportsAnEmptyListAsAnAbsence(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{keyedClient: {body: emptyClients}})

	mac, found, err := c.ClientByMAC(t.Context(), macClient)
	if err != nil {
		t.Fatalf("ClientByMAC: %v", err)
	}

	if found {
		t.Error("an empty list was reported as a client")
	}

	if mac != "" {
		t.Errorf("mac = %q, want empty", mac)
	}
}

// The username resolve counts sessions and nothing else, so what it must get right is which rows
// it counts. The fixture holds the three cases that can be got wrong at once: two rows sharing
// the username, one carrying a different one, and one carrying the empty string most clients
// carry.
func TestClientsByUsernameCountsOnlyExactMatches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		want     int
	}{
		{name: "two sessions share it", username: sharedUsername, want: 2},
		{name: "one session carries it", username: "test-other", want: 1},
		{name: "no session carries it", username: absentUsername, want: 0},
		// The guard against this is in the CLI, and this pins what would happen without it.
		{name: "the empty value is not a wildcard", username: "", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{"common-oper-data": {body: usernameRows}})

			got, err := c.ClientsByUsername(t.Context(), tt.username)
			if err != nil {
				t.Fatalf("ClientsByUsername: %v", err)
			}

			if got != tt.want {
				t.Errorf("sessions = %d, want %d", got, tt.want)
			}
		})
	}
}

// The resolve is pruned to the two leaves it joins on, out of the twenty-two each record carries.
// A pruned read's correctness is in the query rather than in the answer, because a stub cannot
// reproduce what the controller does with a wrong expression.
func TestClientsByUsernamePrunesTheReadToTheTwoLeavesItNeeds(t *testing.T) {
	t.Parallel()

	var query string

	c := newClientWithQuery(t, &query, usernameRows)

	if _, err := c.ClientsByUsername(t.Context(), sharedUsername); err != nil {
		t.Fatalf("ClientsByUsername: %v", err)
	}

	// Spelt out rather than built from clientUsernameFields, for the reason
	// TestAPTagsPrunesTheRequestToTheTagContainer gives: comparing the constant against
	// itself would pass whatever the constant became. Equality and not a substring, so a
	// second parameter appearing beside it is a failure too.
	if want := "fields=client-mac;username"; query != want {
		t.Errorf("query = %q, want %q", query, want)
	}
}

// A controller with no client associated answers the collection empty, which is no sessions and
// not a failure.
func TestClientsByUsernameReportsAnEmptyCollectionAsNoSessions(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"common-oper-data": {body: emptyClients}})

	got, err := c.ClientsByUsername(t.Context(), sharedUsername)
	if err != nil {
		t.Fatalf("ClientsByUsername: %v", err)
	}

	if got != 0 {
		t.Errorf("sessions = %d, want 0", got)
	}
}

// The RPC's input is a mandatory choice of three arms and only one may be sent, so each arm is
// asserted whole: the leaf it fills, and that the body carries nothing beside it. The address
// arm is given rowClientMAC, whose case no controller measured has served, so what it pins is the
// SDK's normalization and not a transform any real row goes through.
func TestDeauthenticateClientSendsOneArmInItsCanonicalForm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*Client) error
		want string
	}{
		{
			name: "the address arm",
			call: func(c *Client) error { return c.DeauthenticateClientByMAC(t.Context(), rowClientMAC) },
			want: `{"input":{"mac-addr":"` + macClient + `"}}`,
		},
		{
			name: "the username arm",
			call: func(c *Client) error {
				return c.DeauthenticateClientByUsername(t.Context(), sharedUsername)
			},
			want: `{"input":{"user-name":"` + sharedUsername + `"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusNoContent, "")

			if err := tt.call(r.client); err != nil {
				t.Fatalf("the post failed: %v", err)
			}

			got := r.last(t)
			if got.method != http.MethodPost {
				t.Errorf("method = %s, want POST", got.method)
			}

			if !contains(got.path, deauthRPC) {
				t.Errorf("path = %s, want the client delete operation", got.path)
			}

			if got.body != tt.want {
				t.Errorf("payload = %s, want %s", got.body, tt.want)
			}
		})
	}
}

// The RPC declares no output container and answers 204, so an empty body is the success shape
// here rather than the failure it is for save-config.
func TestDeauthenticateClientAcceptsAnAnswerWithNoBody(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{deauthRPC: {status: http.StatusNoContent}})

	if err := c.DeauthenticateClientByMAC(t.Context(), macClient); err != nil {
		t.Errorf("a 204 was reported as a failure: %v", err)
	}

	if err := c.DeauthenticateClientByUsername(t.Context(), sharedUsername); err != nil {
		t.Errorf("a 204 was reported as a failure: %v", err)
	}
}

// A release that does not serve the operation answers 400 rather than 404 — measured on
// 17.12.8, carrying error-tag malformed-message and error-message "invalid path". The CLI
// re-words that status, so this pins the classification the re-wording keys on, for both arms.
func TestDeauthenticateClientReportsTheRejectedPath(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*Client) error{
		"the address arm": func(c *Client) error {
			return c.DeauthenticateClientByMAC(t.Context(), macClient)
		},
		"the username arm": func(c *Client) error {
			return c.DeauthenticateClientByUsername(t.Context(), sharedUsername)
		},
	}

	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{deauthRPC: {status: http.StatusBadRequest}})

			err := call(c)
			if err == nil {
				t.Fatal("a 400 did not produce an error")
			}

			if cause, status := Classify(err); cause != CauseHTTP || status != http.StatusBadRequest {
				t.Errorf("Classify = %q/%d, want http/400", cause, status)
			}
		})
	}
}
