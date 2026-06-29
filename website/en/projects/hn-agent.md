# HN Digest Agent

## Project Positioning

This is the project used to learn Go end to end: build a backend service that periodically fetches Hacker News stories, asks an LLM to generate Chinese summaries, and sends the digest to a chosen channel.

It covers core backend skills:

- HTTP client.
- JSON encoding and decoding.
- Concurrent fetching.
- Deduplication.
- Storage.
- Background jobs.
- API server.
- AI API calls.
- Testing, logs, metrics, and deployment.

## MVP Capabilities

1. Fetch HN `topstories`.
2. Fetch story item details.
3. Deduplicate and store stories.
4. Generate Chinese summaries with an LLM.
5. Build a digest.
6. Expose `/healthz` and `/api/v1/digests/latest`.

## Suggested Package Layout

```text
hn-agent/
  cmd/hn-agent/
  internal/config/
  internal/hn/
  internal/summary/
  internal/digest/
  internal/store/
  internal/notifier/
  internal/httpapi/
  internal/obs/
```

## First Acceptance Bar

- `go test ./...` passes.
- The HN client is covered with `httptest`.
- Every external request has a timeout.
- Summary generation failure does not discard fetched stories.
- Each digest has an idempotency key to avoid duplicate delivery.

