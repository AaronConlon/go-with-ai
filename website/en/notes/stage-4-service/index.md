# Stage 4: Service

## Goal

Turn the CLI workflow into a long-running HTTP service.

::: tip Current Status
Stage 4 is complete. The service now has an entrypoint, health check, configuration loading, structured logging, graceful shutdown, and live reload guidance for development.
:::

## Tasks

- [x] Add `cmd/hn-agent`.
- [x] Implement `/healthz`.
- [x] Read configuration.
- [x] Add structured logs.
- [x] Handle graceful shutdown.
- [x] Document `context` and development live reload.

## Key Concepts

- A service needs clear startup and shutdown paths.
- Health checks should be cheap and reliable.
- Config should be explicit and easy to inspect.
- Logs should help diagnose production behavior.
- Digest APIs will be expanded later with the AI integration work.

## Notes

- [Stage 4 Coach Notes](/en/notes/stage-4-service/coach-notes)
- [Go Context Basics: How A Service Knows When To Stop](/en/notes/stage-4-service/go-context)
