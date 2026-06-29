# JS to Go

This project is designed for a learner who already has JavaScript / Node.js experience.

The goal is not to translate JavaScript line by line into Go. The useful move is to change the engineering model:

| JavaScript / Node.js | Go |
| --- | --- |
| `package.json` and runtime scripts | `go.mod`, packages, and compiled binaries |
| Exceptions or rejected promises | Explicit `error` returns |
| `fetch` / HTTP clients | `net/http` and `http.Client` |
| `Promise.all` | goroutine plus `errgroup` |
| Middleware-heavy frameworks | Standard library first, framework later |

## Practical Rule

When learning Go from a Node.js background, ask:

> What resource is created, who owns it, and when is it closed or canceled?

That question naturally leads to `context`, `defer`, explicit errors, and bounded concurrency.

