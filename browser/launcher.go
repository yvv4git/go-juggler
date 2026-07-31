package browser

import (
	"context"
	"errors"
	"time"
)

// Option customises browser launching.
type Option func(*config)

// config holds the launcher options.
type config struct {
	execPath string
	timeout  time.Duration
	headless bool
}

// WithExecPath sets the browser executable to launch.
func WithExecPath(path string) Option {
	return func(c *config) { c.execPath = path }
}

// WithTimeout sets the launch timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *config) { c.timeout = d }
}

// WithHeadless toggles headless mode.
func WithHeadless(headless bool) Option {
	return func(c *config) { c.headless = headless }
}

func defaultConfig() config {
	return config{timeout: 30 * time.Second}
}

// Launch starts a new browser process and connects to it.
//
// TODO: spawn cfg.execPath with the inherited pipe fds
// (transport.ReadFD/transport.WriteFD) or a WebSocket endpoint, then
// Connect to the resulting transport.
func Launch(ctx context.Context, opts ...Option) (*Browser, error) {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.execPath == "" {
		return nil, errors.New("browser: exec path is required (WithExecPath)")
	}

	return nil, errors.New("browser: launching is not implemented yet")
}
