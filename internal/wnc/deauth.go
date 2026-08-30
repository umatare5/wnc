package wnc

import "context"

// Neither arm of apf-ms-delete-all reports whether a client was there: the RPC declares no output
// container and answers 204 for an identifier matching nothing exactly as for a session it
// dropped, which is why each call has a resolve in front of it.

// DeauthenticateClientByMAC drops the session of the client at one address. The SDK lowercases the
// address and re-separates it with colons, which is the form common-oper-data already serves, so
// the resolved address this is given reaches the wire unchanged.
func (c *Client) DeauthenticateClientByMAC(ctx context.Context, clientMAC string) error {
	return c.sdk.Client().DeauthenticateByMAC(ctx, clientMAC)
}

// DeauthenticateClientByUsername drops the sessions authenticated under one username. The plural
// is deliberate: the leaf is a bare string with no cardinality stated, so the caller resolves the
// count and puts it in the prompt.
func (c *Client) DeauthenticateClientByUsername(ctx context.Context, username string) error {
	return c.sdk.Client().DeauthenticateByUsername(ctx, username)
}
