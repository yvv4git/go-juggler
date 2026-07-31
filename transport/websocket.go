package transport

import (
	"context"
	"net/url"
	"sync"

	"github.com/coder/websocket"
)

// WebSocketTransport carries messages over a single WebSocket connection.
type WebSocketTransport struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

// DialWebSocket establishes a WebSocket transport to the browser's
// automation endpoint (e.g. ws://localhost:9377).
func DialWebSocket(ctx context.Context, endpoint *url.URL) (*WebSocketTransport, error) {
	conn, _, err := websocket.Dial(ctx, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebSocketTransport{conn: conn}, nil
}

// Send writes a single JSON frame.
func (t *WebSocketTransport) Send(ctx context.Context, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.conn.Write(ctx, websocket.MessageText, data)
}

// Receive reads the next JSON frame.
func (t *WebSocketTransport) Receive(ctx context.Context) ([]byte, error) {
	_, data, err := t.conn.Read(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Close tears down the WebSocket connection.
func (t *WebSocketTransport) Close() error {
	return t.conn.CloseNow()
}
