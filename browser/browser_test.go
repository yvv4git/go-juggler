package browser

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type fakeTransport struct {
	mu       sync.Mutex
	closed   bool
	received [][]byte
}

func (f *fakeTransport) Send(_ context.Context, data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return errors.New("closed")
	}
	f.received = append(f.received, append([]byte(nil), data...))
	return nil
}

func (f *fakeTransport) Receive(context.Context) ([]byte, error) {
	return nil, errors.New("no data")
}

func (f *fakeTransport) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	return nil
}

func TestConnectNilTransport(t *testing.T) {
	if _, err := Connect(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil transport")
	}
}

func TestConnectAndClose(t *testing.T) {
	tr := &fakeTransport{}
	b, err := Connect(context.Background(), tr)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if !tr.closed {
		t.Fatal("expected transport to be closed")
	}
}

func TestLaunchRequiresExecPath(t *testing.T) {
	if _, err := Launch(context.Background()); err == nil {
		t.Fatal("expected error when exec path is missing")
	}
}
