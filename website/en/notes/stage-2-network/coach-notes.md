# Stage 2 Coach Notes

These notes summarize the Chinese coaching material for Stage 2.

## Focus

- Define API response structs.
- Understand exported fields and `json` tags.
- Build an HTTP client with `net/http`.
- Use `context` for cancellation and deadlines.
- Decode JSON with `json.Decoder`.
- Return explicit errors for non-2xx status codes.
- Use `httptest` to test without calling the real HN API.

## Error Handling Model

Go usually returns `(result, error)` instead of relying on a surrounding `catch`.

```go
ids, err := client.TopStories(ctx)
if err != nil {
    return err
}
```

Read this as: get the result and the error together, check the error first, and only then use the result.

## How to Tell Whether `Item` Is Valid

Do not use `Item{}` itself as the success/failure signal. Use `err`.

::: danger Check `err` First
For `(result, error)`, the rule is: **check `err` before using the result**.

When `err != nil`, the result value is not trustworthy. When `err == nil`, the result value is meaningful.
:::

```go
item, err := client.Item(ctx, id)
if err != nil {
    return err
}

// item is meaningful only after err is nil
fmt.Println(item.Title)
```

Rule:

```text
err != nil: the result value is not trustworthy.
err == nil: the function succeeded, and the result value is meaningful.
```

`Item{}` is only the zero value returned to satisfy the `(Item, error)` signature on failure. It is not a status flag. If the API later needs to represent “not found, but not an error,” use a clearer shape such as `(*Item, error)`, `(Item, bool, error)`, or a well-defined sentinel error.

## CLI Boundary

The target command is:

```bash
go run ./cmd/hnctl top --limit=10
```

This command only becomes meaningful after `Client.TopStories` is wired into `cmd/hnctl`. Before that, the acceptance check is:

```bash
go test ./internal/hn
```
