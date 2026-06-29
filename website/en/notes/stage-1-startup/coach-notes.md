# Stage 1 Coach Notes

These notes summarize the Chinese coaching material for Stage 1.

## Focus

- Understand the workspace boundary.
- Create `hn-agent/` as the actual Go module.
- Learn `package`, `import`, `struct`, functions, tests, `main`, `os.Args`, and `os.Exit`.
- Build a minimal CLI without adding a CLI framework too early.

## Current CLI Boundary

Stage 1 only needs:

```bash
go run ./cmd/hnctl
go run ./cmd/hnctl version
```

Future commands such as `top --limit=10` may already appear in the roadmap, but they are not Stage 1 behavior yet. At this stage, those arguments only demonstrate how values enter `run(args []string)`.

## Acceptance

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```

