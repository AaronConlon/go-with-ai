# Full Code

This section keeps the complete code for each learning stage, including test files.

“Complete code” means the full files added or substantially changed in that stage. Files that stay unchanged from a previous stage are not repeated on every later page, so there is one clear place to review each stage change.

These code blocks use learning-oriented comments. They do not aim for minimal production comments. In the first three stages, comments explain syntax, standard library usage, test intent, error handling, and concurrency mechanics.

Complex concepts are explained where they first appear. For example, `fmt.Fprint` does not print to the terminal; it writes content to a writer, which is why tests can use it to write an HTTP response body.

Every full-file code block starts with a repository-relative path comment:

```go
// hn-agent/internal/hn/story.go
package hn
```

## Stage Index

- [Stage 1: Startup](/en/code/stage-1-startup): Go module, minimal CLI, domain model, and first tests.
- [Stage 2: Networking](/en/code/stage-2-network): HN API client, HTTP, JSON, and `httptest`.
- [Stage 3: Concurrency](/en/code/stage-3-concurrency): `FetchItems`, `errgroup`, semaphore, and batch fetch tests.
- [Stage 4: Service](/en/code/stage-4-service): placeholder until implementation starts.
- [Stage 5: AI Integration](/en/code/stage-5-ai-integration): placeholder until implementation starts.
- [Stage 6: Delivery](/en/code/stage-6-delivery): placeholder until implementation starts.
