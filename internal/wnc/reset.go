package wnc

import "context"

// ResetCAPWAP tears down and re-establishes one access point's CAPWAP session, which restarts
// the session and not the access point. The name arm is what goes on the wire, so no address the
// operator never typed appears in the payload.
func (c *Client) ResetCAPWAP(ctx context.Context, apName string) error {
	return c.sdk.AP().ResetCAPWAPByName(ctx, apName)
}

// ResetAP restarts one access point. The RPC declares no output container, so a 204 establishes
// that the instruction was accepted and nothing more, and reset-ap is absent from the payload
// because the schema declares it defaulting to true, which is the restart this performs.
func (c *Client) ResetAP(ctx context.Context, apName string) error {
	return c.sdk.AP().ResetAPByName(ctx, apName)
}
