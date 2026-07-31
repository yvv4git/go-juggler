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

// Browser is the re-exported browser.Browser type.
type (
	Browser = browser.Browser

	// Tab is the re-exported browser.Tab type.
	Tab = browser.Tab

	// Option is the re-exported browser.Option type.
	Option = browser.Option

	// Client is the re-exported browser.Client type.
	Client = browser.Client

	// HealthResponse is the re-exported browser.HealthResponse type.
	HealthResponse = browser.HealthResponse

	// TabResponse is the re-exported browser.TabResponse type.
	TabResponse = browser.TabResponse

	// SnapshotResponse is the re-exported browser.SnapshotResponse type.
	SnapshotResponse = browser.SnapshotResponse

	// EvaluateResponse is the re-exported browser.EvaluateResponse type.
	EvaluateResponse = browser.EvaluateResponse

	// ResourceEntry is the re-exported browser.ResourceEntry type.
	ResourceEntry = browser.ResourceEntry

	// TabInfo is the re-exported browser.TabInfo type.
	TabInfo = browser.TabInfo

	// ListTabsResponse is the re-exported browser.ListTabsResponse type.
	ListTabsResponse = browser.ListTabsResponse
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
