// Package transport provides the low-level channels between the client and
// a Juggler-enabled browser: an inherited OS pipe pair (fd 3/fd 4) or a
// WebSocket.
//
// A Transport carries opaque, newline-framed JSON messages. Message
// semantics live in the protocol package; transports only move bytes in
// both directions.
package transport

import "context"

// Transport is the bidirectional byte channel between the client and the
// browser. Implementations must be safe for concurrent use.
type Transport interface {
	// Send writes a single JSON frame to the browser. The frame must be
	// newline-terminated by the implementation.
	Send(ctx context.Context, data []byte) error
	// Receive reads the next JSON frame from the browser.
	Receive(ctx context.Context) ([]byte, error)
	// Close tears down the transport and releases associated resources.
	Close() error
}
