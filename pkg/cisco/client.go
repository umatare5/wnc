// Package cisco provides a wrapper around the cisco-ios-xe-wireless-go library
package cisco

import (
	"errors"
	"time"

	wnc "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// Client represents a WNC client
type Client = wnc.Client

// ClientOption represents a client configuration option
type ClientOption = wnc.ClientOption

// NewClient creates a new WNC client with the specified parameters
func NewClient(controller, apikey string, isSecure *bool) (*Client, error) {
	config := wnc.Config{
		Controller:  controller,
		AccessToken: apikey,
	}

	if isSecure != nil && !*isSecure {
		config.InsecureSkipVerify = true
	}

	return wnc.NewClient(config)
}

// NewClientWithTimeout creates a new WNC client with timeout
func NewClientWithTimeout(controller, apikey string, timeout time.Duration, isSecure *bool) (*Client, error) {
	// Validate timeout - zero or negative timeouts are not allowed
	if timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}

	config := wnc.Config{
		Controller:  controller,
		AccessToken: apikey,
		Timeout:     timeout,
	}

	if isSecure != nil && !*isSecure {
		config.InsecureSkipVerify = true
	}

	return wnc.NewClient(config)
}

// NewClientWithOptions creates a new WNC client with custom options
// Note: This function is now deprecated. Use NewClientWithConfig instead.
func NewClientWithOptions(controller, apikey string, options ...ClientOption) (*Client, error) {
	config := wnc.Config{
		Controller:  controller,
		AccessToken: apikey,
	}

	// Since the new API uses Config directly, this function is limited
	// For more advanced configuration, use NewClientWithConfig directly
	return wnc.NewClient(config)
}

// NewClientWithConfig creates a new WNC client with a Config struct
func NewClientWithConfig(config wnc.Config) (*Client, error) {
	return wnc.NewClient(config)
}
