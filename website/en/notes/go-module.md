# Go Module

## What It Solves

A Go module defines the boundary of a Go project. It records the module path and dependencies in `go.mod`.

In this workspace, the repository root is not the Go module. The real Go project lives in:

```text
hn-agent/
```

Run Go commands from there:

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```

## Key Point

Keep the learning workspace and the Go module separate. This keeps docs, practice notes, and project code from being confused with one another.

## What `go mod tidy` Does

`go mod tidy` recalculates module dependencies from the current source code.

Read it as:

```text
scan Go imports
-> add missing dependencies
-> remove unused dependencies
-> update go.sum checksums
-> reclassify direct and indirect dependencies
```

It is not just formatting `go.mod`. It makes `go.mod` and `go.sum` match the code that actually exists.

For example, after:

```bash
go get golang.org/x/sync/errgroup
```

`go.mod` may contain:

```go
require golang.org/x/sync v0.21.0 // indirect
```

After the source code really imports and uses:

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(ctx)
```

running:

```bash
go mod tidy
```

lets Go rescan the source. Since the dependency is now directly imported, `// indirect` is usually removed:

```go
require golang.org/x/sync v0.21.0
```

Run it after adding or removing third-party imports, after `go get`, and before committing Go dependency changes.
