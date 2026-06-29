# Stage 2: Networking Full Code

Stage 2 adds the HN API client: HTTP requests, JSON decoding, timeout, and `httptest`.

Stage 1 files remain in place. This page keeps the full files added or substantially changed in Stage 2, with learning-oriented comments.

## Files

```text
hn-agent/internal/hn/client.go
hn-agent/internal/hn/client_test.go
```

## What `context` Does

`context` is Go's way to pass control signals through a call chain.

In Stage 2, the two most important signals are:

- Cancellation: the caller no longer wants to wait, so ongoing work should stop soon.
- Timeout: if work takes too long, it is canceled automatically.

It is not a replacement for normal business parameters. A `context.Context` is usually the first parameter of a function and tells that function:

```text
Which request does this work belong to?
Has it timed out?
Has it been canceled?
```

For example:

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
```

This creates an HTTP request and binds it to `ctx`. If `ctx` is canceled, the HTTP request is canceled too.

Tests often use:

```go
ids, err := client.TopStories(context.Background())
```

`context.Background()` is an empty root context. It has no timeout and is not canceled by itself. It is common in `main`, tests, or outer entry points. You can derive a timeout context from it:

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```

`defer cancel()` cleans up resources associated with that timeout context.

## HN Client

```go
// hn-agent/internal/hn/client.go
package hn

import (
	// context carries cancellation, timeout, and request-scope signals.
	"context"
	// encoding/json decodes JSON responses into Go values.
	"encoding/json"
	// fmt creates formatted strings and errors.
	"fmt"
	// net/http provides the standard HTTP client and server APIs.
	"net/http"
	// time represents durations such as 10 * time.Second.
	"time"
)

// Item mirrors the Hacker News API item shape.
// Fields must be exported, meaning uppercase, so encoding/json can set them.
type Item struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Descendants int    `json:"descendants"`
}

// Client stores the information needed to call the HN API.
type Client struct {
	// BaseURL is replaceable so tests can point the client at a local server.
	BaseURL string
	// HTTP is the standard client that performs requests.
	HTTP *http.Client
}

// NewClient builds the default production HN client.
func NewClient() *Client {
	return &Client{
		BaseURL: "https://hacker-news.firebaseio.com/v0",
		HTTP: &http.Client{
			// Timeout prevents external I/O from hanging forever.
			Timeout: 10 * time.Second,
		},
	}
}

// TopStories fetches the list of top story IDs.
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
	// NewRequestWithContext connects the HTTP request to ctx.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/topstories.json", nil)
	if err != nil {
		return nil, err
	}

	// Do sends the request.
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	// Always close response bodies.
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ids []int64
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// Item fetches one HN item by ID.
func (c *Client) Item(ctx context.Context, id int64) (Item, error) {
	url := fmt.Sprintf("%s/item/%d.json", c.BaseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// Item is a struct, so it cannot be nil.
		return Item{}, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Item{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
```

## HN Client Test

```go
// hn-agent/internal/hn/client_test.go
package hn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTopStories(t *testing.T) {
	// httptest.NewServer starts a local HTTP server for this test.
	// It replaces the real Hacker News API.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/topstories.json" {
			t.Fatalf("expected path /topstories.json, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[1, 2, 3]`))
	}))
	// Close the test server when this test ends.
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}

	ids, err := client.TopStories(context.Background())
	if err != nil {
		t.Fatalf("TopStories returned error: %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(ids))
	}
}
```

## Acceptance

```bash
cd hn-agent
go test ./internal/hn
```
