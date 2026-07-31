package transport

import (
	"bufio"
	"context"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrClosed is returned when operating on a closed transport.
var ErrClosed = errors.New("transport: closed")

// File descriptors of the pipes inherited from the browser process.
const (
	// ReadFD is the fd the browser writes to; the client reads from it.
	ReadFD = 3
	// WriteFD is the fd the client writes to; the browser reads from it.
	WriteFD = 4
)

// PipeTransport carries messages over a pair of inherited OS pipes, the
// in-process transport used when the browser is spawned by the client.
// Messages are framed as newline-terminated JSON.
type PipeTransport struct {
	rFile *os.File
	r     *bufio.Reader
	wFile *os.File
	w     *bufio.Writer

	rMu sync.Mutex
	wMu sync.Mutex

	closed bool
}

// NewPipe returns a PipeTransport reading from r and writing to w.
//
// The exact fd mapping is defined by the caller so the transport stays
// independent of the launcher. NewPipeFromFds builds the conventional
// fd 3/fd 4 pairing.
func NewPipe(r, w *os.File) *PipeTransport {
	return &PipeTransport{
		rFile: r,
		r:     bufio.NewReader(r),
		wFile: w,
		w:     bufio.NewWriter(w),
	}
}

// NewPipeFromFds builds a PipeTransport from the inherited file descriptors
// ReadFD and WriteFD.
func NewPipeFromFds() *PipeTransport {
	r := os.NewFile(ReadFD, "juggler-read")
	w := os.NewFile(WriteFD, "juggler-write")

	return NewPipe(r, w)
}

// Send writes a single newline-terminated JSON frame to the browser.
func (t *PipeTransport) Send(_ context.Context, data []byte) error {
	t.wMu.Lock()
	defer t.wMu.Unlock()

	if t.closed {
		return ErrClosed
	}

	if _, err := t.w.Write(data); err != nil {
		return err
	}

	if err := t.w.WriteByte('\n'); err != nil {
		return err
	}

	return t.w.Flush()
}

// Receive reads the next JSON frame from the browser. A blocking read is
// interrupted as soon as ctx is done.
func (t *PipeTransport) Receive(ctx context.Context) ([]byte, error) {
	t.rMu.Lock()
	defer t.rMu.Unlock()

	if t.closed {
		return nil, ErrClosed
	}

	stop := context.AfterFunc(ctx, func() { _ = t.rFile.SetReadDeadline(time.Now()) })
	defer stop()

	line, err := t.r.ReadBytes('\n')
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, ctx.Err()
		}

		return nil, err
	}

	return trimNewline(line), nil
}

// Close closes both pipe ends.
func (t *PipeTransport) Close() error {
	t.rMu.Lock()
	t.wMu.Lock()
	defer t.rMu.Unlock()
	defer t.wMu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	var errs []error
	if err := t.rFile.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := t.wFile.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func trimNewline(line []byte) []byte {
	if n := len(line); n > 0 && line[n-1] == '\n' {
		line = line[:n-1]
	}

	if n := len(line); n > 0 && line[n-1] == '\r' {
		line = line[:n-1]
	}

	return line
}
