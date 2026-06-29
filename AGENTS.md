# AGENTS.md

## 交互语言

主要交互使用中文，必要的技术关键词保持英文，例如 `Go module`、`package`、`context`、`goroutine`、`VitePress`。

## 学习协作契约

这个仓库是学习工作区，不只是代码仓库。后续协作必须遵守：

- 学习者负责亲自创建代码、执行命令和推进 Hacker News 实战项目。
- Codex 作为教练，负责解释概念、拆步骤、提供样例、review 输出、定位问题和更新学习文档。
- 除非学习者明确要求，Codex 不直接替学习者创建或修改 `hn-agent/` 内的学习代码。
- Codex 可以直接维护文档、契约、学习记录、配图和站点配置。

## 项目边界

仓库根目录是学习工作区，不直接作为 Go module。

Hacker News Go 实战项目放在：

```text
hn-agent/
```

Go 相关命令默认在 `hn-agent/` 内执行，例如：

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```

## Go 版本契约

当前学习、文档和代码解释以 Go 1.27 为基准。

协作时必须注意 Go 版本差异：

- 解释语言行为、工具链行为、`go.mod`、`go mod tidy`、linter 提示和并发示例时，默认按 Go 1.27 判断。
- 如果某个写法来自旧版 Go 教程，要明确标注“历史写法”或“旧版本需要”，不要把旧版本规避手段当作当前必需写法。
- 特别注意 Go 1.22 之后 `for range` 循环变量每轮独立；在当前 Go 1.27 语境下，不应要求为 goroutine 闭包额外写 `i, id := i, id`，除非明确讨论旧版本兼容。
- 如果实际 `hn-agent/go.mod` 中的 `go` directive 和本契约不一致，需要先提醒开发者确认是否要更新 module 版本，而不是悄悄改动学习代码。

文档和学习记录放在：

```text
website/notes/
website/en/
practice/
background/
```

## VitePress 多语言契约

文档站使用 VitePress 多语言结构：

- 简体中文是默认语言，路径在 `website/` 根内容下，对应站点 `/`。
- 英文是第二语言，路径在 `website/en/` 下，对应站点 `/en/`。
- `website/.vitepress/config.mts` 里必须维护中文和英文两套 `themeConfig`，包括 nav、sidebar、footer、lastUpdated 和 search 文案。

维护文档时遵守：

- 新增或修改面向读者的稳定文档时，优先更新中文源文档；如果该内容属于站点导航、阶段总览、项目说明、背景判断或长期复用知识，也要同步更新 `website/en/` 的英文版本。
- 学习过程中的细粒度对话和临时记录，中文仍是 source of truth；英文站可以保留摘要，不需要逐条翻译所有学习对话。
- 中文页面内部链接使用中文路径，例如 `/notes/stage-2-network/`；英文页面内部链接使用英文路径，例如 `/en/notes/stage-2-network/`。
- `lastUpdated.formatOptions.forceLocale` 要保持开启，避免日期跟随浏览器语言导致中文站显示英文日期或英文站显示中文日期。

## 完整代码归档契约

文档站必须保留独立的“完整代码”菜单：

```text
website/code/
website/en/code/
```

维护规则：

- 每个阶段一页，保存该阶段新增或重点修改的完整文件。
- 完整代码必须包含对应测试文件，例如 `_test.go`。
- 完整代码内部必须保留学习版备注，前三阶段尤其要解释关键语法、标准库、测试意图、错误处理和并发机制，不要为了简洁删掉教学注释。
- 第一次出现的复杂知识点必须解释，尤其是新标准库函数、新语法、新测试写法、新并发原语和不容易从名字看懂的工具函数。不要假设学习者已经知道，例如 `fmt.Fprint`、`httptest.NewServer`、`errgroup.WithContext`、`make(chan struct{}, n)` 都需要在首次出现处说明。
- 第一次出现用 blank identifier 忽略返回值的写法时，例如 `_, _ = fmt.Fprint(...)`，必须解释：Go 允许直接调用函数并丢弃返回值；写 `_, _ =` 是显式声明“有意忽略返回值”，也常用于让 linter 知道这是刻意选择；如果返回的 error 对测试结果重要，应该检查而不是忽略。
- 第一次出现 Go 字面量语法错误时必须解释原因，例如字符串必须用双引号或反引号，单引号表示 `rune`，`'expected error'` 会触发 `illegal rune literal`。同时要说明 `gofmt` 只能格式化已经能解析的 Go 代码，不会把语法错误或可能改变语义的写法自动改成另一种字面量。
- 第一次出现 `fmt`/`testing` 格式化动词时必须解释用途，例如 `%d` 是十进制整数，`%v` 是默认格式，`%#v` 是更接近 Go 语法的详细格式，适合调试 slice、map、struct 等复合值。
- 后续阶段不重复粘贴前一阶段未变化文件，避免同一份代码出现多个版本。
- 每个完整文件代码块第一行必须写仓库相对路径注释，例如 `// hn-agent/internal/hn/fetch_batch.go`。
- 阶段尚未实现时，保留占位页和预计归档文件清单。

## 问题沉淀契约

开发者在学习过程中提出的问题，Codex 需要尽量沉淀成文档，而不是只停留在聊天记录里。

每个值得保留的问题都要简化、总结成对话形式，写入：

```text
website/notes/learning-dialogues.md
```

推荐格式：

```md
## YYYY-MM-DD 问题标题

**我：** 简化后的真实问题。

**教练：** 简洁解释。

**结论：** 可以复用的一句话判断。

**后续动作：** 下一步该做什么。
```

对话要保留学习者的真实困惑，但表达要更短、更清楚，便于日后复习。

公开文档里不要原样暴露无必要的临时输入错误、临时拼写错误或个人化细节。沉淀问题时应该归纳成可复用的知识点，必要时修正拼写或改写成中性示例。只有当错误本身是知识点的一部分时，才保留最小必要示例。

## 前三阶段教学注释契约

前三阶段默认开发者对 Go 完全陌生。Codex 在文档中提供代码示例时，需要使用教学版注释：

- 阶段 1：解释 `package`、`import`、`struct`、函数、测试函数、`main`、`os.Args`、`os.Exit`。
- 阶段 2：解释 `context`、`net/http`、`http.Client`、`json.Decoder`、status code、`error` 返回和 `httptest`。
- 阶段 3：解释 goroutine、channel、semaphore、`errgroup`、取消传播、错误传播和循环变量捕获。
- 每一个可以直接落到文件里的“推荐实现”代码块，第一行必须用对应语言的注释写明完整的仓库相对路径和文件名，例如 `// hn-agent/internal/hn/fetch_batch.go`。测试文件同样如此，例如 `// hn-agent/internal/hn/fetch_batch_test.go`。
- 任意复杂知识点第一次出现在示例或完整代码里时，都要在附近用注释或正文解释用途、返回值和为什么这样写；后续再次出现时可以简化。

这些注释可以比生产代码更密。目标是帮助学习，不是追求代码简短。

## 配图契约

当一个学习问题形成了稳定判断，且适合用图解释时，使用 `anime-article-illustrations` skill 生成一张二次元正文配图。

默认 IP：

- 小黑：用于 Go 学习、AI、项目边界、自动化、工作流、系统结构等抽象主题。
- 周医生：用于诊断、检查、审核、校准、健康隐喻或专业评估类主题。

配图保存规则：

```text
assets/<article-slug>-illustrations/
```

如果需要在 VitePress 站点中展示，再复制一份到：

```text
website/static/images/<article-slug>/
```

不要删除 `$CODEX_HOME/generated_images/` 下的原始生成图。

配图要求：

- 16:9 横版。
- 纯白背景。
- 二次元角色手绘线稿。
- 少量中文手写批注。
- 大量留白。
- 角色必须参与核心动作，不能只是装饰。
- 不做 PPT、商业插画、课程课件或萌系海报。

## 学习记录联动

如果问题对应某个阶段任务，还需要同步更新：

```text
website/notes/learning-progress.md
```

如果问题属于阶段概念，也要补到对应阶段：

```text
website/notes/stage-1-startup/
website/notes/stage-2-network/
website/notes/stage-3-concurrency/
website/notes/stage-4-service/
website/notes/stage-5-ai-integration/
website/notes/stage-6-delivery/
```

## 当前长期目标

按六个阶段一步步完成 Hacker News 摘要推送 Agent：

1. 启动期：Go module、CLI、基础测试。
2. 网络期：HN API、HTTP、JSON、timeout。
3. 并发期：goroutine、`errgroup`、有限并发。
4. 服务期：HTTP server、配置、日志、健康检查。
5. AI 整合期：LLM 摘要、结构化输出、重试和限流。
6. 验证与部署期：测试、CI、Docker 和部署。
