# Six-stage Roadmap

## Stage 1: Startup

Goal: create, run, and test a Go module.

Must do:

- Run `go mod init`.
- Run `go test ./...`.
- Build a minimal CLI.
- Write a table-driven test.

Output:

- `cmd/hnctl`.
- Initial `internal/hn` package.

## Stage 2: Networking

Goal: call the Hacker News API reliably.

Must do:

- Build requests with `http.NewRequestWithContext`.
- Decode responses with `json.Decoder`.
- Return clear errors for non-2xx responses.
- Give every external request a timeout.

Output:

- `TopStories(ctx)`.
- `Item(ctx, id)`.

## Stage 3: Concurrency

Goal: fetch story details with bounded concurrency.

Must do:

- Use `errgroup.WithContext` to manage concurrent work.
- Limit concurrency with a channel or semaphore.
- Define the failure policy clearly.

Output:

- `FetchItems(ctx, ids, concurrency)`.

## Stage 4: Service

Goal: turn CLI capability into a long-running service.

Must do:

- `/healthz`.
- `/api/v1/digests/latest`.
- Graceful shutdown.
- Structured logs with `slog`.

Output:

- `cmd/hn-agent`.

## Stage 5: AI Integration

Goal: generate stable Chinese summaries.

Must do:

- Read the API key from backend config.
- Add timeouts to every request.
- Retry 429 and 5xx responses with limits.
- Validate JSON output from the model.

Output:

- `internal/summary`.

## Stage 6: Delivery

Goal: validate and deploy through CI.

Must do:

- Cover HTTP clients with `httptest`.
- Test database logic.
- Write a Docker multi-stage build.
- Run `go test ./...` in GitHub Actions.

Output:

- Deployable service.
- Reproducible build flow.

