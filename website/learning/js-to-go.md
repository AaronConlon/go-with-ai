# JS 到 Go 心智迁移

## 不是逐句翻译

从 JavaScript / Node.js 迁移到 Go 时，最重要的是避免把已有写法逐句翻译。

例如：

- `Promise.all()` 不等于“随便开很多 goroutine”。
- `try/catch` 不等于“把错误集中到最外层处理”。
- Express middleware 不等于“任何 Go 服务都要先上框架”。

## 高价值对照

| 你熟悉的方式 | Go 的方式 | 迁移重点 |
| --- | --- | --- |
| `async/await` | goroutine + `context` | 关注取消、超时和资源释放。 |
| `Promise.all()` | `errgroup` | 关注错误传播和并发边界。 |
| `try/catch` | 显式 `error` 返回 | 错误是 API 契约的一部分。 |
| npm scripts | Go 官方命令 | `go test`、`go build`、`go fmt` 是主路径。 |
| Express/Fastify | `net/http` | 先用标准库理解底层模型。 |

## 一条实用规则

每当你想写一个抽象时，先问：

> 这是调用方真实需要的边界，还是我为了“看起来像框架”提前设计出来的？

Go 倾向于小接口、小包和直接的数据流。先让代码跑通、测稳，再抽象。

