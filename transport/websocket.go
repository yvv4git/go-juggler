package transport

import (
	"context"
	"errors"
	"net/url"
)

// ErrNotImplemented reports that a transport is declared but not yet wired
// to a concrete implementation.
var ErrNotImplemented = errors.New("transport: not implemented")

// DialWebSocket establishes a WebSocket transport to the browser's
// automation endpoint.
//
// TODO: wire up a WebSocket client such as github.com/gorilla/websocket or
// nhooyr.io/websocket. The endpoint is retained here so the signature stays
// stable while the connection logic lands.
func DialWebSocket(ctx context.Context, endpoint *url.URL) (Transport, error) {
	return nil, ErrNotImplemented
}
