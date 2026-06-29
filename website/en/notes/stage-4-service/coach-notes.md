# Stage 4 Coach Notes

Stage 4 turns the previous package and CLI work into a long-running HTTP service.

The main shift is:

> The program does not finish immediately. It listens for requests, handles them, logs what happens, and shuts down cleanly.

## Mental Model

Use the Go standard library first:

```text
config -> dependencies -> routes -> HTTP server -> requests -> graceful shutdown
```

Keep boundaries clear:

- `cmd/hn-agent` starts the service.
- `internal/server` owns HTTP handlers and routes.
- `internal/config` loads runtime configuration.
- domain logic stays in packages such as `internal/hn`.

## Target Files

Start with:

```text
hn-agent/cmd/hn-agent/main.go
hn-agent/internal/config/config.go
hn-agent/internal/server/server.go
hn-agent/internal/server/server_test.go
```

The first API should be small:

```text
GET /healthz
```

Later, add:

```text
GET /api/v1/digests/latest
POST /api/v1/jobs/digest
```

## Terms

| Term | Meaning |
| --- | --- |
| HTTP server | A long-running process that listens for HTTP requests |
| handler | Code that reads a request and writes a response |
| route | A path mapped to a handler |
| mux / router | Dispatches requests to handlers |
| health check | A cheap endpoint used to verify the service is alive |
| graceful shutdown | Stop accepting new work while allowing in-flight requests to finish |
| config | Runtime settings such as address and timeout |
| structured logging | Logs written as key/value pairs |

## What `net/http` Does

Stage 4 uses `net/http` on the server side:

- `http.Server`
- `http.ServeMux`
- `http.HandlerFunc`

And in tests:

- `httptest.NewRecorder`
- `httptest.NewRequest`

The basic flow is:

```text
client -> HTTP request -> Go HTTP server -> handler -> HTTP response
```

## Handler Shape

```go
// hn-agent/internal/server/server.go
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
```

- `w http.ResponseWriter`: writes status, headers, and body.
- `r *http.Request`: the current request.

Unlike some JavaScript frameworks, handlers do not return a response object. They write the response through `w`.

## Recommended Server Shape

The `/healthz` path is a good candidate for a constant. This is similar to JavaScript:

```js
const HEALTHZ_PATH = "/healthz"
```

In Go, a lowercase package-level constant such as `healthzPath` is unexported, so it can be used inside the `server` package. If another package needs to use it later, rename it to `HealthzPath`.

```go
// hn-agent/internal/server/server.go
package server

import (
	"net/http"
)

// healthzPath is the fixed path for the health check endpoint.
// The lowercase name keeps it internal to the server package.
const healthzPath = "/healthz"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(healthzPath, healthz)
	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}
```

`http.ServeMux` is the standard-library router. It is enough for the first Stage 4 version.

## Handler Test

```go
// hn-agent/internal/server/server_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodGet, healthzPath, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
```

## Go Function Calls Do Not Omit Arguments

This call has three arguments:

```go
req := httptest.NewRequest(http.MethodGet, healthzPath, nil)
```

You can read the `httptest.NewRequest` signature as:

```go
func NewRequest(method string, target string, body io.Reader) *http.Request
```

So the call must provide exactly three arguments:

- `method`: `http.MethodGet`.
- `target`: `healthzPath`.
- `body`: `nil`, because this request has no body.

This differs from JavaScript. In JS, a missing argument usually becomes `undefined`:

```js
function request(method, path, body) {}

request("GET", "/healthz") // body is undefined
```

Go does not do that. Too few or too many arguments fail at compile time:

```go
httptest.NewRequest(http.MethodGet, healthzPath)
// not enough arguments in call to httptest.NewRequest

httptest.NewRequest(http.MethodGet, healthzPath, nil, nil)
// too many arguments in call to httptest.NewRequest
```

`nil` is not an automatic default. It is an explicit value you pass, and the parameter type must allow `nil`. Here, `body` is an `io.Reader` interface, so `nil` means “no request body”. For value types such as `string`, `int`, and `bool`, you cannot pass `nil`.

Go also does not have default parameter values for ordinary functions. Common alternatives are explicit zero values, config structs, or variadic parameters such as `fmt.Println(a ...any)`.

## Why Tests Can Call `NewMux` Directly

`server_test.go` starts with:

```go
package server
```

That means the test file belongs to the same package as `server.go`. When you run:

```bash
go test ./internal/server
```

Go compiles the package's ordinary `.go` files and `_test.go` files together into a temporary test binary:

```text
server.go      -> production code
server_test.go -> test code
go test        -> compile together, then run TestXxx
```

So the test can call:

```go
mux := NewMux()
```

This is not an accidental import of non-test code. It is the point of the test: call the real production code and verify its behavior.

If the test rebuilt its own mux instead of calling `NewMux`, it could accidentally test a fake copy rather than the real service routing.

## `package server` vs `package server_test`

Go tests commonly use one of two package names:

```go
package server
```

This is an internal package test. It can access names in the same package. It is easier for Stage 4.

The other style is:

```go
package server_test
```

This is an external package test. It imports the package like a real user and can only access exported names, such as `server.NewMux()`.

For Stage 4, start with:

```go
package server
```

Learn external package tests later when package boundaries feel clearer.

`httptest.NewRecorder()` creates a fake `http.ResponseWriter`. It records what the handler writes, so tests can assert the response without opening a real port.

## Config

The service entrypoint will call `config.Load()` first, so define the first version of the config package before reading `cmd/hn-agent/main.go`.

Start small: read the listen address from an environment variable.

```go
// hn-agent/internal/config/config.go
package config

import "os"

// Config stores runtime settings needed when the service starts.
// Stage 4 starts with Addr. Later stages can add timeouts or API URLs.
type Config struct {
	Addr string
}

// Load reads configuration from environment variables.
// If HN_AGENT_ADDR is not set, the service listens on :8080.
func Load() Config {
	addr := os.Getenv("HN_AGENT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		Addr: addr,
	}
}
```

### Why `Load()` Returns `Config`, Not `*Config`

This first version uses:

```go
func Load() Config
```

instead of:

```go
func Load() *Config
```

because the first `Config` is small and mostly read-only after startup:

- It currently has one field, `Addr string`, so copying it is cheap.
- `Load()` can always return a usable config. If the environment variable is missing, it uses `:8080`.
- A value return avoids unnecessary shared mutable state.
- Callers do not need to ask whether `cfg == nil`.

Use a pointer when the struct is large, the function needs to express “no value”, multiple places must share and mutate the same object, or methods need to mutate the receiver.

If config loading can fail later, for example because `OPENAI_API_KEY` is required, prefer:

```go
func Load() (Config, error)
```

On success return `cfg, nil`; on failure return `Config{}, err`. This matches the explicit error handling pattern from earlier stages.

Avoid `.env`, config files, or third-party config libraries for the first slice. The Stage 4 goal is the service lifecycle, not a full configuration system.

In `main.go`, this line:

```go
cfg := config.Load()
```

means: load runtime config from the `internal/config` package, then use `cfg.Addr` as the HTTP server listen address.

### Why The Function Must Be `Load`, Not `load`

Go uses the first letter to control package visibility:

- `Load`: exported, so other packages can call it.
- `load`: unexported, so only the `config` package can call it.

`hn-agent/internal/config/config.go` belongs to:

```go
package config
```

The service entrypoint `hn-agent/cmd/hn-agent/main.go` belongs to:

```go
package main
```

They are different packages. So `main.go` must import the package and call the exported function:

```go
import "github.com/aaron/go-with-ai/hn-agent/internal/config"

func main() {
	cfg := config.Load()
	_ = cfg
}
```

If you write:

```go
func load() Config
```

an editor warning such as `function "load" is unused` is expected. It is an unexported function, nothing inside the `config` package calls it, and external packages cannot call it.

Also note that `cfg := config.Load()` only reads the config. If `cfg` is not used afterward, Go will report `cfg` as unused. The complete entrypoint uses `cfg.Addr`:

```go
srv := &http.Server{
	Addr: cfg.Addr,
}
```

### Beginner Mistake: `load`, `config`, and `cfg` Warnings Together

When Stage 4 first splits code into packages, these warnings often appear together:

```text
function "load" is unused
undefined: config
declared and not used: cfg
```

They usually come from an incomplete version like this.

`hn-agent/internal/config/config.go`:

```go
package config

type Config struct {
	Addr string
}

func load() Config {
	return Config{Addr: ":8080"}
}
```

`hn-agent/cmd/hn-agent/main.go`:

```go
package main

func main() {
	cfg := config
}
```

There are three issues:

- `load` is lowercase, so only the `config` package can call it.
- `main.go` does not import `internal/config`, so the name `config` does not exist.
- `cfg` is declared but not used.

The corrected direction is:

```go
// hn-agent/internal/config/config.go
package config

import "os"

type Config struct {
	Addr string
}

func Load() Config {
	addr := os.Getenv("HN_AGENT_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{Addr: addr}
}
```

```go
// hn-agent/cmd/hn-agent/main.go
package main

import (
	"net/http"

	"github.com/aaron/go-with-ai/hn-agent/internal/config"
)

func main() {
	cfg := config.Load()

	srv := &http.Server{
		Addr: cfg.Addr,
	}

	_ = srv
}
```

Debug it in this order:

1. If another package must call the function, use exported `Load`.
2. Check that the caller imports the package.
3. Check that the variable is actually used.

## Service Entrypoint

```go
// hn-agent/cmd/hn-agent/main.go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaron/go-with-ai/hn-agent/internal/config"
	"github.com/aaron/go-with-ai/hn-agent/internal/server"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := server.NewMux()

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "addr", cfg.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
```

## `main()` Walkthrough

The entrypoint combines many new Go ideas at once. Read it as a service lifecycle:

```text
load config
create logger
create mux
create http.Server
listen for shutdown signals
start the server in a goroutine
wait in the main goroutine
create a shutdown timeout
call srv.Shutdown
log that the service stopped
```

`cfg := config.Load()` loads runtime settings from `internal/config`. The first config only has `Addr`, which becomes the HTTP server listen address.

`logger := slog.New(slog.NewTextHandler(os.Stdout, nil))` creates a structured logger that writes text logs to standard output. The `nil` option means default handler options.

If you want a custom timestamp format and a fixed Shanghai timezone, replace it with:

```go
logger := newLogger()
```

and see the custom logger section below.

`mux := server.NewMux()` creates the routing table. `main.go` does not need to know handler details; it only gives the mux to the HTTP server.

```go
srv := &http.Server{
	Addr:              cfg.Addr,
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
}
```

`&http.Server{...}` creates a server struct and takes its address, so `srv` is a `*http.Server`. `Addr` decides where to listen, `Handler` routes requests, and `ReadHeaderTimeout` avoids waiting forever for incomplete request headers.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
```

`signal.NotifyContext` creates a context that is canceled when the process receives `Ctrl+C` (`os.Interrupt`) or `SIGTERM`. It returns `ctx` for waiting and `stop` for cleanup.

```go
go func() {
	...
}()
```

starts an anonymous function in a new goroutine. This is necessary because `srv.ListenAndServe()` blocks while the server is running. If it ran in the main goroutine, the shutdown waiting code would never run.

```go
if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
```

means: treat the error as fatal only if it is not the expected `http.ErrServerClosed` returned during graceful shutdown.

A common beginner mistake is to forget the `!`:

```go
if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
	logger.Error("server failed", "error", err)
	os.Exit(1)
}
```

That reverses the meaning: the expected shutdown error is logged as a failure. After pressing `Ctrl+C`, you may see:

```text
level=INFO msg="server stopped"
level=ERROR msg="server failed" error="http: Server closed"
exit status 1
```

The correct condition is:

```go
if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	logger.Error("server failed", "error", err)
	os.Exit(1)
}
```

Read it as: only log a fatal error when it is not the normal server-closed error.

```go
<-ctx.Done()
stop()
```

waits until the signal context is canceled, then stops signal notification for this context.

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

creates a fresh timeout context for shutdown. The earlier `ctx` is already canceled, so `Shutdown` needs a new context that gives the server a limited amount of time to finish in-flight requests.

```go
if err := srv.Shutdown(shutdownCtx); err != nil {
```

stops accepting new requests and waits for in-flight requests to finish until the timeout expires.

If execution reaches:

```go
logger.Info("server stopped")
```

the service completed the shutdown path.

If `context` still feels unfamiliar, read the standalone page: [Go Context Basics: How A Service Knows When To Stop](/en/notes/stage-4-service/go-context). It explains `context.Background()`, `signal.NotifyContext`, `context.WithTimeout`, and `<-ctx.Done()` with illustrations.

`ListenAndServe` blocks, so it runs in a goroutine. The main goroutine waits for a shutdown signal and then calls `Shutdown`.

## What `<-ctx.Done()` Means

This line is very Go-specific:

```go
<-ctx.Done()
```

Read it as:

```text
wait here until ctx is canceled.
```

`ctx.Done()` returns a channel. You can think of its type as:

```go
<-chan struct{}
```

That means “a receive-only channel used as a signal”. It is not carrying business data. It signals that the context is done.

The leading `<-` is a receive operation:

```go
<-someChannel
```

If no value is available and the channel is not closed, the current goroutine blocks.

Here, `signal.NotifyContext(...)` creates a context that is canceled when the process receives `os.Interrupt` or `syscall.SIGTERM`. When that happens, the `ctx.Done()` channel is closed, so:

```go
<-ctx.Done()
```

unblocks and the code continues to `srv.Shutdown(...)`.

If you wrote only:

```go
ctx.Done()
```

you would merely get the channel. The program would not wait. The `<-` is what makes this line wait.

## `log/slog`

`log/slog` is the standard-library structured logging package:

```go
logger.Info("server starting", "addr", cfg.Addr)
logger.Error("server failed", "error", err)
```

The message is paired with key/value fields, which are easier to search later.

## Custom Log Time Format And Shanghai Timezone

The default text handler prints timestamps like:

```text
time=2026-06-29T16:18:29.075+08:00 level=INFO msg="server starting" addr=:8080
```

To force a custom format and the Shanghai timezone, use `slog.HandlerOptions.ReplaceAttr`:

```go
// hn-agent/cmd/hn-agent/main.go
func newLogger() *slog.Logger {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		shanghai = time.FixedZone("Asia/Shanghai", 8*60*60)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				t := a.Value.Time().In(shanghai)
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
			}

			return a
		},
	})

	return slog.New(handler)
}
```

Then call:

```go
logger := newLogger()
```

The output becomes:

```text
time="2026-06-29 16:18:29" level=INFO msg="server starting" addr=:8080
```

Go time layouts use the fixed reference time `Mon Jan 2 15:04:05 MST 2006`, not `YYYY-MM-DD HH:mm:ss`. So the layout for `2026-06-29 16:18:29` is:

```go
"2006-01-02 15:04:05"
```

## Development Live Reload

Go does not include a Vite-like dev server with built-in hot reload. The common development workflow is to use a watcher that:

```text
stops the old process
rebuilds the binary
starts the new process
```

This is live reload, not production hot deploy.

For Stage 4, use `air`:

```bash
go install github.com/air-verse/air@latest
```

If `air` is not found after installation, make sure your Go bin directory is in `PATH`. It is commonly:

```text
$HOME/go/bin
```

On macOS, Homebrew also provides:

```bash
brew install go-air
```

Because this project's service entrypoint is:

```text
hn-agent/cmd/hn-agent
```

create `hn-agent/.air.toml`:

```toml
# hn-agent/.air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/hn-agent ./cmd/hn-agent"
bin = "./tmp/hn-agent"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
```

Then run inside `hn-agent/`:

```bash
air
```

When you save a `.go` file, Air rebuilds and restarts the service. If behavior looks confusing, fall back to:

```bash
go test ./...
go run ./cmd/hn-agent
```

to separate Go code issues from watcher configuration issues.

## Graceful Shutdown

Graceful shutdown means:

```text
receive stop signal
-> stop accepting new requests
-> give in-flight requests time to finish
-> exit after timeout
```

The standard-library method is:

```go
srv.Shutdown(ctx)
```

Use a timeout context so shutdown does not hang forever.

## Recommended Steps

1. Add `internal/server` with `NewMux()` and `/healthz`.
2. Test `/healthz` with `httptest.NewRecorder`.
3. Add `internal/config` with `HN_AGENT_ADDR`, defaulting to `:8080`.
4. Add `cmd/hn-agent/main.go`.
5. Run `curl http://localhost:8080/healthz`.
6. Add graceful shutdown.
7. Add digest APIs later.

## Test Direction

Start with:

- `GET /healthz` returns 200.
- `POST /healthz` returns 405.
- response body is `ok\n`.
- default config address is `:8080`.
- `HN_AGENT_ADDR` overrides the default.

Commands:

```bash
cd hn-agent
go test ./internal/server
go test ./internal/config
go test ./...
```

## Acceptance

```bash
cd hn-agent
go test ./...
go run ./cmd/hn-agent
```

In another terminal:

```bash
curl -i http://localhost:8080/healthz
```

Expected:

```text
HTTP/1.1 200 OK
...
ok
```

## Do Not Start With

Avoid these in the first Stage 4 slice:

- a web framework
- a database
- complex scheduling
- AI summary logic inside handlers

First make the service skeleton testable.
