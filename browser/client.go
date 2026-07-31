package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for the camofox-browser REST API.
type Client struct {
	baseURL string
	hc      *http.Client
}

// Option configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.hc = hc }
}

// HealthResponse describes the browser status.
type HealthResponse struct {
	OK               bool   `json:"ok"`
	Engine           string `json:"engine"`
	Sessions         int    `json:"sessions"`
	BrowserConnected bool   `json:"browserConnected"`
}

// TabResponse is returned when a tab is opened.
type TabResponse struct {
	TabID    string `json:"tabId"`
	TargetID string `json:"targetId"`
	URL      string `json:"url"`
	Title    string `json:"title"`
}

// SnapshotResponse is returned by the snapshot endpoint.
type SnapshotResponse struct {
	URL       string `json:"url"`
	Snapshot  string `json:"snapshot"`
	RefsCount int    `json:"refsCount"`
}

// NewClient creates a new HTTP client for the given camofox-browser address.
func NewClient(addr string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: addr,
		hc: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Health checks the browser status.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var r HealthResponse
	if err := c.get(ctx, "/health", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// OpenTab opens a new tab and navigates to the given URL.
func (c *Client) OpenTab(ctx context.Context, sessionKey, url string) (*TabResponse, error) {
	var r TabResponse
	if err := c.post(ctx, "/tabs/open", map[string]any{
		"userId":     sessionKey,
		"sessionKey": sessionKey,
		"url":        url,
	}, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Navigate navigates an existing tab to the given URL.
func (c *Client) Navigate(ctx context.Context, tabID, sessionKey, url string) error {
	return c.post(ctx, "/tabs/"+tabID+"/navigate", map[string]any{
		"userId": sessionKey,
		"url":    url,
	}, nil)
}

// Snapshot returns the ARIA snapshot of the tab.
func (c *Client) Snapshot(ctx context.Context, tabID, sessionKey string) (*SnapshotResponse, error) {
	var r SnapshotResponse
	if err := c.get(ctx, fmt.Sprintf("/tabs/%s/snapshot?userId=%s", tabID, sessionKey), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// CloseTab closes a tab.
func (c *Client) CloseTab(ctx context.Context, tabID, sessionKey string) error {
	return c.delete(ctx, "/tabs/"+tabID, map[string]any{"userId": sessionKey})
}

// CloseSession destroys the entire session and all its tabs.
func (c *Client) CloseSession(ctx context.Context, sessionKey string) error {
	return c.delete(ctx, "/sessions/"+sessionKey, nil)
}

// --- internal helpers ---

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, payload any, out any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, out)
}

func (c *Client) delete(ctx context.Context, path string, payload any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+path, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, nil)
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: %d %s", req.Method, req.URL.Path, resp.StatusCode, body)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
