# Go Context Basics: How A Service Knows When To Stop

Stage 4 uses `context` heavily in the service entrypoint:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

<-ctx.Done()
stop()

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
	...
}
```

This page focuses only on the mental model needed for Stage 4.

## The Short Version

`context` is not a business data container. Treat it as a control ticket that answers:

```text
Should this work continue?
When must it stop?
Has someone canceled it?
```

![Xiaohei pulls the context cancellation signal](/images/go-context/01-context-cancel-signal.png)

In Stage 4, `context` tells the service when to move from running to shutting down.

## `context.Background()`

```go
context.Background()
```

is the empty starting point:

```text
not canceled
no timeout
no extra values
```

Other contexts are derived from it:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

## `signal.NotifyContext`

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
```

means:

```text
create a context,
then cancel it when the process receives Ctrl+C or SIGTERM.
```

It returns:

- `ctx`: the context you wait on.
- `stop`: a cleanup function that stops signal notification.

## `ctx.Done()` Returns A Channel

```go
ctx.Done()
```

returns a channel. You can read its type as:

```go
<-chan struct{}
```

It does not carry business data. It signals that the context is done.

## `<-ctx.Done()` Waits

```go
<-ctx.Done()
```

means:

```text
wait here until ctx is canceled.
```

If no signal has arrived, the current goroutine blocks.

![Xiaohei waits at the Done door](/images/go-context/03-context-done-channel.png)

`ctx.Done()` only gets the channel.

`<-ctx.Done()` waits on the channel.

## Why Shutdown Uses A Fresh Context

After the stop signal arrives, the code creates a new context:

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := srv.Shutdown(shutdownCtx); err != nil {
	...
}
```

The earlier `ctx` is already canceled. Its job was to tell the main goroutine:

```text
someone asked the service to stop.
```

`Shutdown` needs a different context:

```text
give the server up to 5 seconds to finish.
```

The causal direction matters:

```text
shutdownCtx does not trigger shutdown.
calling srv.Shutdown(shutdownCtx) triggers shutdown.
shutdownCtx only limits how long that shutdown may wait.
```

![Xiaohei gives shutdown a 5-second hourglass](/images/go-context/02-context-timeout-shutdown.png)

## `context.WithTimeout`

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

creates a context that is automatically canceled after 5 seconds.

More concretely, when the 5-second deadline expires:

```text
shutdownCtx is marked canceled
-> shutdownCtx.Done() is closed
-> shutdownCtx.Err() returns context deadline exceeded
```

You can observe it with:

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

<-shutdownCtx.Done()

fmt.Println(shutdownCtx.Err())
// context deadline exceeded
```

If you call:

```go
cancel()
```

first, then `shutdownCtx.Err()` is usually:

```text
context canceled
```

So there are two common endings:

| Ending | `Done()` | `Err()` |
| --- | --- | --- |
| deadline expires | closed | `context deadline exceeded` |
| manual `cancel()` | closed | `context canceled` |

`srv.Shutdown(shutdownCtx)` uses it to:

```text
stop accepting new requests,
wait for in-flight requests,
return an error if the timeout expires.
```

More precisely:

```text
you call srv.Shutdown(shutdownCtx)
-> the server immediately enters shutdown
-> it stops accepting new connections
-> it closes idle connections
-> it waits for active requests to finish
-> if shutdownCtx times out or is canceled first, Shutdown returns an error
```

So `shutdownCtx` is not the thing that sends a shutdown message. It is a timer and cancellation switch that tells `Shutdown`: wait, but not forever.

`defer cancel()` releases timer resources.

## Do Not Use Context As Ordinary Storage

For Stage 4, keep this rule:

```text
context carries cancellation and timeout signals.
Do not treat it like a map or global variable.
```

Configuration should stay explicit:

```go
cfg := config.Load()

srv := &http.Server{
	Addr: cfg.Addr,
}
```

## Stage 4 Mental Model

There are two lines of execution:

```text
server goroutine:
ListenAndServe() runs until the server closes

main goroutine:
waits on ctx.Done()
receives stop signal
creates shutdownCtx
calls srv.Shutdown(shutdownCtx)
```

The two contexts have different jobs:

| Name | Created By | Job | Canceled When |
| --- | --- | --- | --- |
| `ctx` | `signal.NotifyContext` | Wait for stop signal | `Ctrl+C` or `SIGTERM` |
| `shutdownCtx` | `context.WithTimeout` | Limit shutdown time | timeout or `cancel()` |

## Beginner Checklist

When you see `context`, ask:

- Where did this `ctx` come from?
- When will it be canceled?
- Who waits on `Done()`?
- What happens after cancellation?

In Stage 4:

```text
ctx comes from signal.NotifyContext.
It is canceled by a stop signal.
main goroutine waits with <-ctx.Done().
then the service calls srv.Shutdown.
```
