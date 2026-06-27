# Go With AI

这是一个用于学习 Go 开发、沉淀工程知识和推进项目实践的 Markdown 文档仓库。

当前结构分为两部分：

- `background/`：学习背景、研究结论、路线设计和项目方向沉淀。
- `practice/`：从零开始推进 HN 摘要推送 Agent 的实战目录。
- `website/`：基于 VitePress 的静态文档站点，用于持续记录 Go 学习笔记、知识卡片和项目实践。
- `hn-agent/`：Hacker News 摘要推送 Agent 的 Go 实战项目目录。这个目录由学习者在阶段一创建。

## 快速开始

进入静态站点目录：

```bash
cd website
npm install
npm run dev
```

构建静态站点：

```bash
cd website
npm run build
```

预览构建结果：

```bash
cd website
npm run preview
```

## 静态资源约定

图片等静态资源统一放在：

```text
website/static/images/
```

在 Markdown 中引用图片时，推荐使用：

```md
![图片说明](/images/example.png)
```

`website/static/` 会作为 VitePress 的静态资源目录，构建时复制到站点根路径。

## 学习记录约定

知识库按六个阶段组织：

1. 启动期。
2. 网络期。
3. 并发期。
4. 服务期。
5. AI 整合期。
6. 验证与部署期。

实战推进记录放在 `practice/`，阶段知识沉淀回填到 `website/notes/` 对应模块。

## Go 项目位置约定

仓库根目录是学习工作区，不直接作为 Go module。Hacker News 实战项目放在：

```text
hn-agent/
```

因此 Go 相关命令默认在 `hn-agent/` 内执行，例如：

```bash
cd hn-agent
go test ./...
go run ./cmd/hnctl
```
