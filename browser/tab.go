package browser

import (
	"context"
	"errors"
)

// Tab is a single page owned by a Browser.
type Tab struct {
	b  *Browser
	id string
}

// ID returns the tab's protocol identifier.
func (t *Tab) ID() string { return t.id }

// Navigate navigates the tab to the given URL.
//
// TODO: send the Page.navigate request and wait for its response.
func (t *Tab) Navigate(ctx context.Context, url string) error {
	return errors.New("browser: not implemented yet")
}

// Reload reloads the current page.
//
// TODO: send the Page.reload request and wait for its response.
func (t *Tab) Reload(ctx context.Context) error {
	return errors.New("browser: not implemented yet")
}

// close tears down the tab.
//
// TODO: send the Browser.closePage request.
func (t *Tab) close(ctx context.Context) error {
	return nil
}
