# Stage 3: Concurrency

## Goal

Fetch story details with bounded concurrency.

::: tip Current Status
Stage 3 has been completed and synced to the remote repository by the learner. The project now moves into Stage 4: Service.
:::

## Tasks

- Use goroutines for parallel fetches.
- Limit concurrency with a semaphore-like channel.
- Use `errgroup` to collect errors.
- Propagate cancellation through `context`.
- Understand why old Go examples copied loop variables, and why this project does not need that pattern on Go 1.27.

## Acceptance

First verify the package behavior:

```bash
cd hn-agent
go test ./internal/hn
go test -run TestFetchItems -v ./internal/hn
```

Optionally run the race detector:

```bash
go test -race ./internal/hn
```

After the CLI is wired to the batch fetch logic, run the integration check:

```bash
go run ./cmd/hnctl top --limit=20 --concurrency=5
```

Stage 3 acceptance means bounded concurrency, stable result order, error propagation, and cancellation propagation. The CLI command is the final integration check, not the only proof.

## Key Concepts

- A goroutine is lightweight concurrent work.
- A channel can coordinate or limit work.
- `errgroup.WithContext` helps connect errors and cancellation.
- Bounded concurrency keeps external APIs and local resources under control.

## Term Notes

| English | Chinese mental model |
| --- | --- |
| goroutine | Go 的轻量并发任务 |
| concurrency | 并发：managing multiple tasks during the same period |
| parallelism | 并行：tasks actually running at the same instant |
| channel | 通道：communication or synchronization pipe |
| semaphore | 信号量 / 并发许可证 |
| errgroup | 错误组 / goroutine 任务组 |
| cancellation | 取消：ask unfinished work to stop |
| context | 上下文：carries cancellation, timeout, and request scope |
