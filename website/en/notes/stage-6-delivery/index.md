# Stage 6: Delivery

## Goal

Make the project testable, reproducible, and deployable.

## Tasks

- Keep `go test ./...` stable.
- Cover HTTP clients and handlers.
- Write a Dockerfile.
- Add GitHub Actions.
- Document deployment steps.

## Key Concepts

- Tests should protect behavior, not implementation trivia.
- CI should run the same checks developers run locally.
- A Docker build should be reproducible.
- Deployment docs should explain rollback and verification.

