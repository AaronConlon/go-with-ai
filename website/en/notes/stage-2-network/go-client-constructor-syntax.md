# Go Client Constructor Syntax

## Example

```go
func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        HTTP: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}
```

## Reading It

- `NewClient` is a constructor-style function.
- `*Client` means it returns a pointer to `Client`.
- `&Client{...}` creates a `Client` value and returns its address.
- `BaseURL: baseURL` sets a struct field.
- `10 * time.Second` builds a duration.

## Practical Rule

Use a constructor when setup needs defaults, such as timeout values or base URLs.

