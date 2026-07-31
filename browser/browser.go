// Package browser provides high-level control of a Juggler-enabled browser:
// launching and shutting down a process, and managing tabs.
//
// A Browser wraps a transport.Transport with a message pump that turns
// protocol messages into typed calls, so callers work with Browser and Tab
// handles instead of raw JSON.
package browser

import (
	"context"
	"errors"
	"sync"

	"github.com/yvv4git/go-juggler/transport"
)

// Browser is a client handle to a running Juggler-enabled browser.
type Browser struct {
	tr transport.Transport

	mu   sync.Mutex
	tabs []*Tab
}

// Connect wraps an established transport with a high-level Browser handle.
func Connect(ctx context.Context, tr transport.Transport) (*Browser, error) {
	if tr == nil {
		return nil, errors.New("browser: nil transport")
	}
	// TODO: perform the Juggler session handshake (Browser.newSession) and
	// wait for the first event before returning a ready Browser.
	return &Browser{tr: tr}, nil
}

// Close shuts down the browser and releases the transport.
func (b *Browser) Close(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var errs []error
	for _, tab := range b.tabs {
		if err := tab.close(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	b.tabs = nil

	if err := b.tr.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// Tabs returns the tabs currently owned by the browser.
func (b *Browser) Tabs() []*Tab {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]*Tab(nil), b.tabs...)
}
