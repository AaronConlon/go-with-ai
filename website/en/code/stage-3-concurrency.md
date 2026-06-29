# Stage 3: Concurrency Full Code

Stage 3 adds bounded batch fetching: `FetchItems`, `errgroup`, semaphore, and batch fetch tests.

Stage 1 and Stage 2 files remain in place. This page keeps the full files added or substantially changed in Stage 3, with comments that explain goroutines, channels, semaphores, error propagation, and test intent.

## Files

```text
hn-agent/internal/hn/fetch_batch.go
hn-agent/internal/hn/fetch_batch_test.go
```

## Batch Fetch Implementation

### What `context` Does in Concurrency

In Stage 3, `context` does more than cancel one HTTP request. It connects a group of concurrent tasks.

This line:

```go
g, ctx := errgroup.WithContext(ctx)
```

derives a new `ctx` from the old one. If any goroutine returns an error, `errgroup` cancels the derived `ctx`. The HTTP requests then use that derived context:

```go
item, err := c.Item(ctx, id)
```

The effect is:

```text
one item request fails
-> that goroutine returns an error
-> errgroup cancels ctx
-> other unfinished HTTP requests can observe cancellation
-> FetchItems returns an error
```

So in Stage 3, you can read `context` as the signal that tells related concurrent work: “stop soon.”

```go
// hn-agent/internal/hn/fetch_batch.go
package hn

import (
	// context carries cancellation and timeout signals.
	// In Stage 3, it lets a group of goroutines be canceled together.
	"context"
	"fmt"

	// errgroup manages a group of goroutines that return errors.
	"golang.org/x/sync/errgroup"
)

// FetchItems fetches multiple HN items by ID.
// concurrency limits how many requests can run at the same time.
func (c *Client) FetchItems(ctx context.Context, ids []int64, concurrency int) ([]Item, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be positive")
	}

	// Preallocate the result slice so each goroutine writes to its own index.
	// This keeps the output order aligned with the input IDs.
	items := make([]Item, len(ids))

	// WithContext returns a group and a derived context.
	// If any goroutine returns an error, the derived context is canceled.
	g, ctx := errgroup.WithContext(ctx)

	// sem is a semaphore implemented with a buffered channel.
	// Its capacity is the maximum number of concurrent requests.
	sem := make(chan struct{}, concurrency)

	// Go 1.27 uses per-iteration range variables, so no old i, id := i, id copy is needed.
	for i, id := range ids {
		// g.Go starts a goroutine and records its returned error.
		g.Go(func() error {
			// Acquire one concurrency permit.
			sem <- struct{}{}
			defer func() {
				// Release the permit when this goroutine exits.
				<-sem
			}()

			// Use the errgroup-derived context so cancellation can propagate.
			item, err := c.Item(ctx, id)
			if err != nil {
				return err
			}

			items[i] = item
			return nil
		})
	}

	// Wait blocks until all goroutines finish.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return items, nil
}
```

## Batch Fetch Test

### What `fmt.Fprint` Does in Tests

`fmt.Fprint(w, "...")` means: write this string to `w`.

Here `w` is an `http.ResponseWriter`, so this does not print to the terminal. It writes the JSON into the HTTP response body returned by the test server. The function returns the number of bytes written and an error, so the example uses:

```go
_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
```

The two `_` values intentionally ignore those return values.

You could also write:

```go
fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
```

Go allows a function call to stand alone as a statement, so its return values can be discarded. That means `_, _ =` is not required by the language.

Here `_, _ =` makes the intent explicit: the function has return values, and this test intentionally ignores the byte count and the write error. Many linters also prefer explicit ignoring.

If the write error matters for the test, check it instead:

```go
if _, err := fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`); err != nil {
	t.Fatalf("write response: %v", err)
}
```

### Why `t.Fatal('expected error, got nil')` Is Invalid

Go strings use double quotes or backticks:

```go
t.Fatal("expected error, got nil")
```

Single quotes mean `rune`, which is one Unicode character:

```go
letter := 'a'
```

This is invalid:

```go
t.Fatal('expected error, got nil') // illegal rune literal
```

The single quotes contain many characters, but a rune literal can only represent one character.

`gofmt` cannot automatically fix this. It formats Go code that already parses. `illegal rune literal` is a parse error, so there is no valid AST to format. Also, single quotes and double quotes produce different types, so automatic replacement could change meaning.

### `%d` vs `%#v`

The first argument to `t.Fatalf` is a format string. Verbs such as `%d` and `%#v` tell Go how to print the following values.

```go
t.Fatalf("expected 3 items, got %d", len(items))
```

This uses `%d` because `len(items)` returns an integer. `%d` prints a decimal integer.

```go
t.Fatalf("expected nil items on error, got %#v", items)
```

This uses `%#v` because `items` is a slice. If the test fails, we want to see whether it is nil, empty, or contains real items:

```text
[]hn.Item(nil)
[]hn.Item{}
[]hn.Item{hn.Item{ID:101, ...}}
```

`%v` is the default format. `%#v` is more debug-oriented and tries to print a Go-syntax representation, which is useful for slices, maps, and structs.

```go
// hn-agent/internal/hn/fetch_batch_test.go
package hn

import (
	"context"
	// fmt creates errors in implementation code and writes response bodies in tests.
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchItemsInvalidConcurrency(t *testing.T) {
	// This test checks input validation only, so it does not need a real server.
	client := NewClient()

	items, err := client.FetchItems(context.Background(), []int64{101}, 0)
	if err == nil {
		// Test failure messages are strings, so they use double quotes.
		// Single quotes are for runes, such as 'a'.
		t.Fatal("expected error, got nil")
	}
	if items != nil {
		// %#v prints a Go-syntax-like representation of items.
		// items is a []Item slice, so %#v is more useful for debugging than %d.
		t.Fatalf("expected nil items, got %#v", items)
	}
}

func TestFetchItemsSuccessPreservesOrder(t *testing.T) {
	// Local test server that simulates item endpoints.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/item/101.json":
			// fmt.Fprint writes content to the writer passed as its first argument.
			// Here w is an http.ResponseWriter, so this writes JSON into the test
			// HTTP response body rather than printing to the terminal.
			// It returns the number of bytes written and an error.
			// Go allows calling fmt.Fprint(...) directly and discarding return values;
			// here _, _ = explicitly tells readers and linters that both values are intentionally ignored.
			_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"first"}`)
		case "/item/102.json":
			_, _ = fmt.Fprint(w, `{"id":102,"type":"story","title":"second"}`)
		case "/item/103.json":
			_, _ = fmt.Fprint(w, `{"id":103,"type":"story","title":"third"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	items, err := client.FetchItems(context.Background(), []int64{101, 102, 103}, 2)
	if err != nil {
		t.Fatalf("FetchItems returned error: %v", err)
	}

	if len(items) != 3 {
		// len(items) is an int, so %d prints it as a decimal integer.
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	// Completion order may vary, but result order must match input order.
	wantIDs := []int64{101, 102, 103}
	for i, wantID := range wantIDs {
		if items[i].ID != wantID {
			t.Fatalf("items[%d].ID = %d, want %d", i, items[i].ID, wantID)
		}
	}
}

func TestFetchItemsReturnsErrorWhenOneItemFails(t *testing.T) {
	// Simulate one failing item request.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/item/102.json" {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// This also writes JSON into the HTTP response body.
		_, _ = fmt.Fprint(w, `{"id":101,"type":"story","title":"ok"}`)
	}))
	defer server.Close()

	client := &Client{
		BaseURL: server.URL,
		HTTP:    server.Client(),
	}

	items, err := client.FetchItems(context.Background(), []int64{101, 102}, 2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if items != nil {
		// If items is not nil, %#v helps show what it contains.
		t.Fatalf("expected nil items on error, got %#v", items)
	}
}
```

## Acceptance

```bash
cd hn-agent
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
go test -race ./internal/hn
```
