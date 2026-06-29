# Stage 3 Coach Notes

These notes summarize the Chinese coaching material for Stage 3.

## Illustration Guide

These Xiaohei illustrations are Chinese-labeled article figures that help build the mental model before reading the code.

![JS event loop versus Go goroutines](/images/goroutine-js-mindset/01-js-event-loop-vs-go-goroutine.png)

JavaScript often feels like queued callbacks returning to one scheduling model. Go can start many lightweight tasks with `go`, but starting them is only the easy part; managing them is the real work.

![errgroup.WithContext coordinates goroutines](/images/goroutine-js-mindset/02-errgroup-context-cancel.png)

`errgroup.WithContext(ctx)` ties a group of goroutines together: `g.Go` starts work, `g.Wait` gathers the result, and one error can cancel the derived context.

![make creates result slots](/images/goroutine-js-mindset/03-make-items-slots.png)

`make([]Item, len(ids))` preallocates result slots. Each goroutine writes to its own `items[i]`, keeping order stable and avoiding concurrent `append`.

## Focus

- Start goroutines intentionally.
- Use channels or semaphores to limit concurrency.
- Use `errgroup` to coordinate errors.
- Propagate cancellation through `context`.
- Understand old loop-variable capture examples and why Go 1.25 does not need manual copies.

## Term Notes

The Chinese site keeps key English terms while adding Chinese mental models.

| English | Chinese wording | Working meaning |
| --- | --- | --- |
| goroutine | Go 协程 / 轻量并发任务 | Lightweight concurrent task managed by the Go runtime |
| concurrency | 并发 | Managing multiple tasks during the same period |
| parallelism | 并行 | Multiple tasks actually running at the same instant |
| channel | 通道 | Pipe for communication, synchronization, or concurrency control |
| semaphore | 信号量 / 并发许可证 | Limits how many tasks can enter a section at once |
| worker pool | 工作池 | Fixed workers pulling jobs from a queue |
| errgroup | 错误组 / goroutine 任务组 | Starts, waits, collects errors, and propagates cancellation |
| context | 上下文 | Carries cancellation, timeout, and request scope |
| cancellation | 取消 | Asks unfinished work to stop as soon as possible |

For this project, keep the English term in code discussions and use the Chinese phrase to explain the mental model.

## Invalid Concurrency

`FetchItems` returns `([]Item, error)`, so a failure path returns two values:

```go
if concurrency < 1 {
    return nil, fmt.Errorf("concurrency must be positive")
}
```

Read this as: there is no valid item list, and there is an error explaining why.

`fmt.Errorf` does not print. It creates an `error` value. The caller receives it like this:

```go
items, err := client.FetchItems(ctx, ids, concurrency)
if err != nil {
    return err
}

// use items only after err is nil
```

Go functions must also return on the success path. After the concurrent work succeeds, the normal return should look like:

```go
return items, nil
```

## What `errgroup` Is

This import:

```go
"golang.org/x/sync/errgroup"
```

brings in the `errgroup` package. It is not part of the Go standard library. It comes from the Go-maintained extension module `golang.org/x/sync`.

Because it is not a standard-library package, the module must record the dependency first:

```bash
go get golang.org/x/sync/errgroup
go mod tidy
```

This is similar to installing a package in JavaScript. Go records the dependency in `go.mod` and checksum data in `go.sum`.

After `go get`, you may see:

```go
module github.com/aaron/go-with-ai/hn-agent

go 1.25.1

require golang.org/x/sync v0.21.0 // indirect
```

Read it as:

- `require`: this module needs the dependency.
- `golang.org/x/sync`: the dependency module path.
- `v0.21.0`: the selected version.
- `// indirect`: Go currently marks it as an indirect dependency.

`// indirect` is not an error. It often means the dependency is only needed through another dependency, or the current source code has not yet imported and used it directly.

After the code really imports:

```go
import "golang.org/x/sync/errgroup"
```

and uses:

```go
g, ctx := errgroup.WithContext(ctx)
```

run:

```bash
go mod tidy
```

Go will rescan the source. If the package is directly used, the dependency is usually recorded without `// indirect`.

`go mod tidy` works by scanning current imports. If `fetch_batch.go` imports `golang.org/x/sync/errgroup` and actually uses `errgroup.WithContext`, Go treats `golang.org/x/sync` as a direct dependency and removes `// indirect`.

If the source does not import it anymore, `go mod tidy` may remove the dependency entirely.

If this has not been done, the editor or Go language server may not resolve the package, so completion may not work.

Use `errgroup` when you need to manage a group of goroutines that can return errors.

Without it, you would need to manually:

- Start multiple goroutines.
- Wait for all of them.
- Collect errors.
- Cancel the remaining work after one task fails.

## What `errgroup.WithContext(ctx)` Does

This line:

```go
g, ctx := errgroup.WithContext(ctx)
```

The capital `W` matters. The exported function is `WithContext`, not `withContext`. In Go, names that start with an uppercase letter are exported from a package.

returns two values:

- `g`: the group used to start goroutines and wait for errors.
- `ctx`: a derived context that is canceled when any goroutine returns an error.

You start work with:

```go
g.Go(func() error {
    return nil
})
```

Each function returns an `error`:

- `nil` means the task succeeded.
- non-`nil` means the task failed.

Then wait with:

```go
if err := g.Wait(); err != nil {
    return nil, err
}
```

`g.Wait()` waits for all goroutines and returns one of their errors if any task failed.

The derived `ctx` is the important part: if one task fails, the context is canceled, and other HTTP requests using that context can stop early.

Compared with `sync.WaitGroup`, `errgroup.WithContext` adds error collection and cancellation propagation.

If completion does not appear, check three things:

1. The dependency has been added with `go get`.
2. The code uses `errgroup.WithContext`, with uppercase `W`.
3. The editor's Go language server has refreshed after `go.mod` changed.

## What `sem <- struct{}{}` and `<-sem` Mean

This code:

```go
sem <- struct{}{}
defer func() {
    <-sem
}()
```

belongs with the buffered channel created earlier:

```go
sem := make(chan struct{}, concurrency)
```

`sem` acts like a semaphore, or a pool of concurrency permits.

`struct{}{}` is an empty struct value. It carries no business data. Here, the value itself does not matter; occupying one buffered channel slot is what matters.

This line acquires a permit:

```go
sem <- struct{}{}
```

If the channel is not full, the send succeeds and the goroutine continues. If the channel is full, the goroutine blocks until another goroutine releases a permit.

This line releases a permit:

```go
<-sem
```

The `defer` makes sure the permit is released when the current goroutine function returns, whether it succeeds or returns an error:

```go
sem <- struct{}{} // acquire
defer func() {
    <-sem // release
}()
```

For a more cancellation-aware version, the acquire step can use `select`:

```go
select {
case sem <- struct{}{}:
    defer func() {
        <-sem
    }()
case <-ctx.Done():
    return ctx.Err()
}
```

The simple version is enough for the first Stage 3 implementation. The `select` version is a later refinement once cancellation and channel operations feel more familiar.

## Why Preallocate `items`

This line creates a result slice with the same length as `ids`:

```go
items := make([]Item, len(ids))
```

If `ids` has three elements, `items` starts with three slots:

```text
items[0]
items[1]
items[2]
```

Each goroutine writes the fetched item back to its own index:

```go
items[i] = item
```

This keeps the output order aligned with the input order. Even if the request for `ids[2]` finishes first, it still writes to `items[2]`.

It also avoids concurrent `append` on the same slice:

```go
items = append(items, item) // unsafe from multiple goroutines without coordination
```

For Stage 3, remember the pattern: preallocate `make([]Item, len(ids))`, capture `i` and `id` in the loop, then write `items[i] = item`.

The current learning contract uses Go 1.27, so you do not need the old manual-copy pattern:

```go
i, id := i, id
```

Old Go examples used that line to avoid closure capture bugs with reused loop variables. Since Go 1.22, `for range` loop variables are per-iteration variables, so linters may correctly warn that the copy is unneeded.

## What `make` Does Here

`make` is a Go built-in for creating and initializing slices, maps, and channels.

Stage 3 uses it twice:

```go
items := make([]Item, len(ids))
sem := make(chan struct{}, concurrency)
```

The first line creates a `[]Item` with length `len(ids)`, so each goroutine can write directly to `items[i]`.

The second line creates a buffered channel. Its capacity is `concurrency`, so it can act like a pool of concurrency permits.

## What `nil` Can Mean

`nil` is not a universal fallback value. It only works for nil-able types:

| Type | Can be `nil` | Example |
| --- | --- | --- |
| slice | yes | `[]Item` |
| map | yes | `map[string]int` |
| channel | yes | `chan Item` |
| function | yes | `func() error` |
| pointer | yes | `*Client` |
| interface | yes | `error` |
| struct | no | `Item` |
| int / bool / string | no | `int64`, `bool`, `string` |

`FetchItems` can return `nil, err` because its result type is `[]Item`, and a slice can be `nil`.

A function returning `(Item, error)` cannot return `nil, err`; it should return `Item{}, err`.

For `FetchItems`, `return nil, nil` is syntactically valid, but it means “no item list and no error.” Use it only when that meaning is intentional. The normal success path should usually return:

```go
return items, nil
```

## Target Command

Stage 3 should be tested in layers.

First verify the package-level behavior:

```bash
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
```

The required `FetchItems` tests are:

1. `concurrency <= 0` returns an error.
2. Multiple items can be fetched successfully.
3. The returned item order matches the input `ids` order.
4. If one item request fails, `FetchItems` returns an error.

The order test is important because concurrency changes completion order. Even if `ids[2]` finishes first, it should still be written to `items[2]`.

## Common Assertion Inversion

Do not read only the failure message. Also check the `if` condition that triggers `t.Fatal`.

If the test expects an error, write:

```go
items, err := client.FetchItems(ctx, ids, 0)
if err == nil {
	t.Fatal("expected error, got nil")
}
if items != nil {
	t.Fatalf("expected nil items, got %#v", items)
}
```

This means:

- fail when `err == nil`, because an error was expected.
- fail when `items != nil`, because the error path should return nil items.

This version is inverted:

```go
if err != nil {
	t.Fatal("expected error, got nil")
}
```

It fails when an error is actually returned, while the message says “got nil.” That makes the test output misleading.

Useful advanced tests:

- Track the maximum number of active requests and assert it never exceeds `concurrency`.
- Use `context.WithTimeout` or cancellation and assert the error propagates.
- Run `go test -race ./internal/hn` once while working on the concurrency code.

After the CLI is wired to the batch fetch logic, run:

```bash
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

`--concurrency=5` should control the maximum number of simultaneous item fetches.

If the CLI has not been connected yet, do not block Stage 3 on this command. First pass the package tests, then add the CLI integration check.
