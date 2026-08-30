// Package wnc is the only place that imports the RESTCONF SDK. Every read a command
// makes is a method here, returning a shape of this package's own, so the SDK's
// structs and their absence quirks stay on one side of a single seam.
package wnc

import (
	"fmt"

	"github.com/sirupsen/logrus"
	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// Client is the CLI's handle on one controller.
type Client struct {
	sdk *sdk.Client
}

// NewClient builds the client for one controller. All three timeouts are set from
// the one flag: the SDK's own response-header and handshake defaults are a few
// seconds, short enough that a whole-container read of a busy controller trips them
// long before the request timeout is spent.
func NewClient(t config.Target, s config.Settings, logger *logrus.Logger, userAgent string) (*Client, error) {
	c, err := sdk.NewClient(t.Host, t.Token,
		sdk.WithTimeout(s.Timeout),
		sdk.WithResponseHeaderTimeout(s.Timeout),
		sdk.WithTLSHandshakeTimeout(s.Timeout),
		sdk.WithInsecureSkipVerify(s.Insecure),
		sdk.WithUserAgent(userAgent),
		sdk.WithLogger(log.SDKLogger(logger)),
	)
	if err != nil {
		// The SDK's message names the option that failed and never the token.
		return nil, fmt.Errorf("building the client: %w", err)
	}

	return &Client{sdk: c}, nil
}
