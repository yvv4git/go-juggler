package transport

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestPipeRoundTrip(t *testing.T) {
	peerR, clientW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	clientR, peerW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	client := NewPipe(clientR, clientW)
	peer := NewPipe(peerR, peerW)

	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()

	frame := `{"id":1,"method":"Browser.newPage","params":{"url":"https://example.com"}}`
	if err := client.Send(ctx, []byte(frame)); err != nil {
		t.Fatalf("send: %v", err)
	}

	got, err := peer.Receive(ctx)
	if err != nil {
		t.Fatalf("receive: %v", err)
	}

	if string(got) != frame {
		t.Fatalf("got %q, want %q", got, frame)
	}

	reply := `{"id":1,"result":{}}`
	if err := peer.Send(ctx, []byte(reply)); err != nil {
		t.Fatalf("send reply: %v", err)
	}

	got, err = client.Receive(ctx)
	if err != nil {
		t.Fatalf("receive reply: %v", err)
	}

	if string(got) != reply {
		t.Fatalf("got %q, want %q", got, reply)
	}
}

func TestPipeReceiveContextCancel(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	client := NewPipe(r, w)

	t.Cleanup(func() { _ = client.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := client.Receive(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
}
