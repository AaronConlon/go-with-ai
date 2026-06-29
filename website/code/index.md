# 完整代码

这个栏目单独保存每一阶段的完整代码，包含测试文件。

这里的“完整代码”指本阶段新增或重点修改的完整文件。上一阶段已经稳定且未变化的文件，不在后续阶段重复粘贴，避免同一份代码出现多个版本。

这些代码块采用学习版注释，不追求生产代码的最少注释。前三阶段会在代码内部解释关键语法、标准库、测试意图、错误处理和并发机制。

第一次出现的复杂知识点会在出现处解释。比如 `fmt.Fprint` 不是“打印到终端”，而是“把内容写入某个 writer”，在测试里常用来写 HTTP response body。

每个完整代码块第一行都会标注仓库相对路径，例如：

```go
// hn-agent/internal/hn/story.go
package hn
```

## 阶段索引

- [阶段 1：启动期](/code/stage-1-startup)：Go module、最小 CLI、领域模型和第一组测试。
- [阶段 2：网络期](/code/stage-2-network)：HN API client、HTTP、JSON、`httptest`。
- [阶段 3：并发期](/code/stage-3-concurrency)：`FetchItems`、`errgroup`、semaphore、批量抓取测试。
- [阶段 4：服务期](/code/stage-4-service)：待实现后归档。
- [阶段 5：AI 整合期](/code/stage-5-ai-integration)：待实现后归档。
- [阶段 6：验证与部署期](/code/stage-6-delivery)：待实现后归档。

## 使用方式

如果你正在学习某个阶段，建议顺序是：

1. 先读对应阶段的教练笔记。
2. 自己在 `hn-agent/` 里写代码。
3. 卡住时再打开这里对照完整代码。
4. 最后运行该阶段页面里的测试命令。
