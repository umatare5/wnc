package wnc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// Cause labels a failed request, a write included. The set is closed so
// docs/TROUBLESHOOTING.md can index it, and it is what the log line carries as a field.
type Cause string

const (
	CauseAuth       Cause = "auth"
	CauseForbidden  Cause = "forbidden"
	CauseNotFound   Cause = "not-found"
	CauseTimeout    Cause = "timeout"
	CauseTLS        Cause = "tls"
	CauseConnection Cause = "connection"
	CauseHTTP       Cause = "http"
	CauseCancelled  Cause = "canceled"
	CauseInternal   Cause = "internal"
)

// Classify names the cause of a failed request and, where the controller answered, the status it
// answered with. The order is load-bearing: a deadline arrives as ErrRequestTimeout and never as
// an APIError, and a certificate fault is a transport fault too, so both are settled before the
// generic connection case can claim them.
func Classify(err error) (cause Cause, status int) {
	switch {
	case err == nil:
		return "", 0
	case errors.Is(err, context.Canceled):
		return CauseCancelled, 0
	case errors.Is(err, sdk.ErrRequestTimeout), errors.Is(err, context.DeadlineExceeded):
		return CauseTimeout, 0
	}

	var apiErr *sdk.APIError
	if errors.As(err, &apiErr) {
		return statusCause(apiErr.StatusCode), apiErr.StatusCode
	}

	if isTLS(err) {
		return CauseTLS, 0
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return CauseConnection, 0
	}

	return CauseInternal, 0
}

func statusCause(status int) Cause {
	switch status {
	case http.StatusUnauthorized:
		return CauseAuth
	case http.StatusForbidden:
		return CauseForbidden
	case http.StatusNotFound:
		return CauseNotFound
	default:
		return CauseHTTP
	}
}

func isTLS(err error) bool {
	var (
		verify   *tls.CertificateVerificationError
		unknown  x509.UnknownAuthorityError
		hostname x509.HostnameError
		expired  x509.CertificateInvalidError
		record   tls.RecordHeaderError
	)

	return errors.As(err, &verify) ||
		errors.As(err, &unknown) ||
		errors.As(err, &hostname) ||
		errors.As(err, &expired) ||
		errors.As(err, &record)
}

// Message is the one-line description of a failure that is safe to show an operator, and every
// class is written here rather than taken from the error. An *APIError carries up to 512 bytes of
// the controller's own error document, which on some paths echoes the configuration that was read,
// and a transport failure carries a *url.Error naming the full request URL.
func Message(err error) string {
	var apiErr *sdk.APIError
	if errors.As(err, &apiErr) {
		return fmt.Sprintf("the controller answered %d %s",
			apiErr.StatusCode, http.StatusText(apiErr.StatusCode))
	}

	// The timeout clause does not name the flag: the SDK pins net.Dialer.Timeout at 30s and
	// exports no option for it, so a connect-stage timeout is not --timeout elapsing.
	// CauseCancelled is absent because exitCode maps a canceled context to ExitSignal, at which
	// run.go prints nothing.
	switch cause, _ := Classify(err); cause {
	case CauseTimeout:
		return "the controller did not answer in time"
	case CauseTLS:
		return "the controller's certificate did not verify"
	case CauseConnection:
		return "the controller could not be reached"
	}

	// The residual class is the one place the error is quoted, and that is load-bearing:
	// absentOperation's re-wording of a pre-17.15 rejection is a bare errors.New and reaches the
	// operator only through here. What else lands here is a decode or envelope fault whose text is
	// the SDK's and can carry the request path.
	return err.Error()
}
