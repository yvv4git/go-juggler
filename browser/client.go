package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"time"
)

// Client is an HTTP client for the camofox-browser REST API.
type Client struct {
	baseURL string
	hc      *http.Client
}

// ClientOption configures a Client.
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

// --- Tab actions ---

// Click clicks an element identified by ref (e.g. "e1") or CSS selector.
func (c *Client) Click(ctx context.Context, tabID, sessionKey, ref, selector string) error {
	return c.tabAction(ctx, tabID, sessionKey, "click", map[string]any{
		"ref":      ref,
		"selector": selector,
	})
}

// Type fills an input field identified by ref or selector with text.
func (c *Client) Type(ctx context.Context, tabID, sessionKey, ref, selector, text string) error {
	return c.tabAction(ctx, tabID, sessionKey, "type", map[string]any{
		"ref":      ref,
		"selector": selector,
		"text":     text,
	})
}

// Press presses a keyboard key (e.g. "Enter", "Tab", "Escape").
func (c *Client) Press(ctx context.Context, tabID, sessionKey, key string) error {
	return c.tabAction(ctx, tabID, sessionKey, "press", map[string]any{
		"key": key,
	})
}

// Scroll scrolls the page. direction is "up" or "down", amount is pixels.
func (c *Client) Scroll(ctx context.Context, tabID, sessionKey, direction string, amount int) error {
	return c.tabAction(ctx, tabID, sessionKey, "scroll", map[string]any{
		"direction": direction,
		"amount":    amount,
	})
}

// Back navigates back in history.
func (c *Client) Back(ctx context.Context, tabID, sessionKey string) error {
	return c.tabPOST(ctx, "/tabs/"+tabID+"/back", map[string]any{"userId": sessionKey}, nil)
}

// Forward navigates forward in history.
func (c *Client) Forward(ctx context.Context, tabID, sessionKey string) error {
	return c.tabPOST(ctx, "/tabs/"+tabID+"/forward", map[string]any{"userId": sessionKey}, nil)
}

// Refresh reloads the current page.
func (c *Client) Refresh(ctx context.Context, tabID, sessionKey string) error {
	return c.tabPOST(ctx, "/tabs/"+tabID+"/refresh", map[string]any{"userId": sessionKey}, nil)
}

// LinksResponse lists all links on the page.
type LinksResponse struct {
	Links      []LinkEntry `json:"links"`
	Pagination struct {
		Total   int  `json:"total"`
		Offset  int  `json:"offset"`
		Limit   int  `json:"limit"`
		HasMore bool `json:"hasMore"`
	} `json:"pagination"`
}

// LinkEntry is a single link on the page.
type LinkEntry struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

// Links returns all links on the page.
func (c *Client) Links(ctx context.Context, tabID, sessionKey string, limit, offset int) (*LinksResponse, error) {
	var r LinksResponse

	path := fmt.Sprintf("/tabs/%s/links?userId=%s&limit=%d&offset=%d", tabID, sessionKey, limit, offset)
	if err := c.get(ctx, path, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Screenshot returns the raw PNG bytes of the page screenshot.
func (c *Client) Screenshot(ctx context.Context, tabID, sessionKey string, fullPage bool) ([]byte, error) {
	path := fmt.Sprintf("/tabs/%s/screenshot?userId=%s&fullPage=%v", tabID, sessionKey, fullPage)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s: %d %s", path, resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

// TabStats describes tab state.
type TabStats struct {
	TabID       string   `json:"tabId"`
	URL         string   `json:"url"`
	VisitedURLs []string `json:"visitedUrls"`
	ToolCalls   int      `json:"toolCalls"`
	RefsCount   int      `json:"refsCount"`
}

// ResourceEntry describes one network request captured via the Performance API.
type ResourceEntry struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Size   int    `json:"size,omitempty"`
	Status int    `json:"status,omitempty"`
}

// NetworkRequests returns all resource requests loaded by the page
// (main document + subresources) in chronological order.
func (c *Client) NetworkRequests(ctx context.Context, tabID, sessionKey string) ([]ResourceEntry, error) {
	expr := `JSON.stringify([
		...performance.getEntriesByType("navigation").map(e => ({name:e.name, type:"navigation", size:e.transferSize||0})),
		...performance.getEntriesByType("resource").map(e => ({name:e.name, type:e.initiatorType, size:e.transferSize||0}))
	])`

	res, err := c.Evaluate(ctx, tabID, sessionKey, expr)
	if err != nil {
		return nil, err
	}

	raw, ok := res.Result.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", res.Result)
	}

	var entries []ResourceEntry
	if err := json.Unmarshal([]byte(raw), &entries); err != nil {
		return nil, fmt.Errorf("unmarshal entries: %w", err)
	}

	return entries, nil
}

// PollNetworkRequests reads resource entries repeatedly for the given duration
// and merges them with deduplication (first occurrence wins).
func (c *Client) PollNetworkRequests(ctx context.Context, tabID, sessionKey string, duration, interval time.Duration) ([]ResourceEntry, error) {
	deadline := time.Now().Add(duration)
	seen := map[string]bool{}

	var result []ResourceEntry

	for {
		entries, err := c.NetworkRequests(ctx, tabID, sessionKey)
		if err != nil {
			return nil, err
		}

		for _, e := range entries {
			if !seen[e.Name] {
				seen[e.Name] = true
				result = append(result, e)
			}
		}

		if time.Now().After(deadline) {
			break
		}

		time.Sleep(interval)
	}

	return result, nil
}

// EvaluateResponse is the result of a JavaScript evaluation.
type EvaluateResponse struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"`
}

// Evaluate runs an arbitrary JavaScript expression in the page.
func (c *Client) Evaluate(ctx context.Context, tabID, sessionKey, expression string) (*EvaluateResponse, error) {
	body := map[string]string{
		"userId":     sessionKey,
		"expression": expression,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/tabs/"+tabID+"/evaluate", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("evaluate endpoint not found — camofox-browser must be >= 1.4.0 (got %s, status 404)", resp.Status)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(raw))
	}

	var out EvaluateResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode: %w (body: %s)", err, string(raw))
	}

	return &out, nil
}

// TabInfo describes a single tab in a session.
type TabInfo struct {
	TargetID   string `json:"targetId"`
	TabID      string `json:"tabId"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	ListItemID string `json:"listItemId"`
}

// ListTabsResponse holds all tabs for a session.
type ListTabsResponse struct {
	Running bool      `json:"running"`
	Tabs    []TabInfo `json:"tabs"`
}

// ListTabs returns all open tabs in the given session.
func (c *Client) ListTabs(ctx context.Context, sessionKey string) (*ListTabsResponse, error) {
	var r ListTabsResponse
	if err := c.get(ctx, "/tabs?userId="+url.QueryEscape(sessionKey), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Stats returns the tab stats (URL, visited URLs, refs).
func (c *Client) Stats(ctx context.Context, tabID, sessionKey string) (*TabStats, error) {
	var r TabStats
	if err := c.get(ctx, fmt.Sprintf("/tabs/%s/stats?userId=%s", tabID, sessionKey), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// --- internal helpers ---

func (c *Client) tabAction(ctx context.Context, tabID, sessionKey, kind string, params map[string]any) error {
	payload := map[string]any{
		"kind":     kind,
		"userId":   sessionKey,
		"targetId": tabID,
	}
	maps.Copy(payload, params)

	return c.post(ctx, "/act", payload, nil)
}

func (c *Client) tabPOST(ctx context.Context, path string, payload any, out any) error {
	return c.post(ctx, path, payload, out)
}

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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: %d %s", req.Method, req.URL.Path, resp.StatusCode, body)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}
