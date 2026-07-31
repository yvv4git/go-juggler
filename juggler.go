// Package juggler is a Go client for the Juggler browser automation
// protocol. It brings the low-level transport and protocol packages
// together behind a single high-level API for driving Juggler-enabled
// browsers (Firefox/Camoufox).
//
// The package is split into layers:
//
//	transport/  low-level byte channels: pipe (fd3/fd4) and WebSocket
//	protocol/   Juggler message types and JSON encoding/decoding
//	browser/    high-level browser and tab control
//
// This root package re-exports the browser API so that most callers only
// need a single import.
package juggler

import (
	"context"

	"github.com/yvv4git/go-juggler/browser"
	"github.com/yvv4git/go-juggler/transport"
)

// Browser, Tab, Client and related types are re-exported for convenient
// single-package imports.
type (
	Browser = browser.Browser
	Tab     = browser.Tab
	Option  = browser.Option
	Client  = browser.Client

	HealthResponse  = browser.HealthResponse
	TabResponse     = browser.TabResponse
	SnapshotResponse = browser.SnapshotResponse
)

// Connect wraps an established transport with a high-level Browser handle.
func Connect(ctx context.Context, tr transport.Transport) (*Browser, error) {
	return browser.Connect(ctx, tr)
}

// Launch starts a new Juggler-enabled browser and returns a handle to it.
func Launch(ctx context.Context, opts ...Option) (*Browser, error) {
	return browser.Launch(ctx, opts...)
}

// NewClient creates a new HTTP client for the given camofox-browser address.
func NewClient(addr string, opts ...browser.ClientOption) *Client {
	return browser.NewClient(addr, opts...)
}

// WithExecPath sets the browser executable to launch.
var WithExecPath = browser.WithExecPath

// WithTimeout sets the launch timeout.
var WithTimeout = browser.WithTimeout

// WithHeadless toggles headless mode.
var WithHeadless = browser.WithHeadless
