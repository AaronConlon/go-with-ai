# Go `httptest` Usage

## What It Does

`httptest` starts a temporary local HTTP server inside a test. Your client calls that server instead of the real external API.

## Minimal Shape

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _, _ = w.Write([]byte(`[1, 2, 3]`))
}))
defer server.Close()
```

## Key Points

- Put `httptest` in `_test.go` files.
- Point `Client.BaseURL` at `server.URL`.
- Use `defer server.Close()` to clean up the temporary server.
- Check request paths inside the handler.
- Return controlled JSON to test success and failure cases.

