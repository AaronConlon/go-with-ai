# Go Method Receiver Syntax

## Example

```go
func (c *Client) TopStories(ctx context.Context) ([]int64, error) {
    // ...
}
```

## Reading It

- `(c *Client)` is the receiver.
- The function is a method on `*Client`.
- Inside the method, `c` gives access to fields such as `c.BaseURL` and `c.HTTP`.
- `([]int64, error)` means the method returns story IDs and an error.

## Error Handling Compared With JavaScript

In JavaScript, you may write:

```js
try {
  const ids = await client.topStories()
} catch (err) {
  // handle the error
}
```

In Go, the usual shape is:

```go
ids, err := client.TopStories(ctx)
if err != nil {
    return err
}

// use ids only after err is nil
```

Go treats errors as explicit return values. The caller handles the error near the call site before using the successful result.

## Practical Rule

Use methods when behavior belongs to a type and needs its fields.
