# Stage 1: Startup

## Goal

Create, run, and test a minimal Go module.

## Tasks

- Create the independent `hn-agent/` project directory.
- Initialize the Go module.
- Create `cmd/hnctl/` and `internal/hn/`.
- Write the first domain function.
- Write the first table-driven test.
- Run the minimal CLI.

## Acceptance

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```

## Knowledge Docs

- [Go Slice vs Array](/en/notes/stage-1-startup/go-slice-vs-array)
- [Go if Syntax](/en/notes/stage-1-startup/go-if-syntax)

## Key Concepts

- `go.mod` defines the module boundary.
- `cmd/hnctl` is the CLI entry point.
- `internal/hn` contains Hacker News domain code.
- `os.Args[1:]` passes command-line arguments into `run(args []string)`.
- Stage 1 only needs a minimal CLI; future commands like `top --limit=10` are not implemented yet.

