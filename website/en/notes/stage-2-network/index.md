# Stage 2: Networking

## Goal

Call the Hacker News API reliably.

## Tasks

- Define the HN item model.
- Create an HTTP client type.
- Implement `TopStories(ctx)`.
- Implement `Item(ctx, id)`.
- Add timeouts.
- Test the client with `httptest`.

## Acceptance

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl top --limit=10
```

`top --limit=10` is the target CLI shape after wiring the HN client into `hnctl`. Before that integration exists, use `go test ./internal/hn` as the client-level acceptance check.

## Key Concepts

- `context` carries cancellation and deadlines.
- `net/http` is enough for the first HTTP client.
- `json.Decoder` reads API responses into Go structs.
- Non-2xx responses should become explicit errors.
- `httptest.NewServer` replaces the real external API in tests.

