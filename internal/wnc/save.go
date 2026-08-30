package wnc

import (
	"context"
	"errors"
)

// ErrSaveNotReported means the controller accepted the save and said nothing about it. The result
// string is the controller's whole account of a save, so a reply carrying none leaves the outcome
// unknown.
var ErrSaveNotReported = errors.New("the controller reported no result for the save")

// SaveConfig copies the controller's running configuration to its startup configuration and
// returns the controller's own account of it. The result is returned rather than matched, because
// a release wording it differently must not turn a save that worked into a failure; a bodiless
// answer arrives as a nil Output rather than as a decode fault and is refused rather than read.
func (c *Client) SaveConfig(ctx context.Context) (string, error) {
	resp, err := c.sdk.Controller().SaveConfig(ctx)
	if err != nil {
		return "", err
	}

	if resp.Output == nil || resp.Output.Result == "" {
		return "", ErrSaveNotReported
	}

	return resp.Output.Result, nil
}
