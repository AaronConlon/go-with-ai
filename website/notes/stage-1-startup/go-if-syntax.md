# Go `if` 语法

## 三种变体

### 1. 基本条件

```go
if x > 0 {
    fmt.Println("positive")
}
```

- 条件**不需要括号**（和 JS/Python 不同）
- 大括号**不可省略**（和 Go 所有流程控制一致）

### 2. if-else

```go
if x > 0 {
    fmt.Println("positive")
} else if x == 0 {
    fmt.Println("zero")
} else {
    fmt.Println("negative")
}
```

- `else` 必须和 `}` 在同一行，否则编译错误
- 因为 Go 会在 `}` 后自动插入 `;`，导致 `else` 单独成行变成非法语法

### 3. 带 init statement（Go 特色）

```go
if err := g.Wait(); err != nil {
    return nil, err
}
```

分号拆成两部分：

```
if   <初始化语句>  ;  <条件表达式>  {
```

- 初始化语句中的变量只在 `if / else if / else` 块内可见
- 等价于分开写，但作用域更小：

```go
// 分开写：err 会污染外层
err := g.Wait()
if err != nil {
    return nil, err
}

// init 写法：err 在 if 结束后不可见
if err := g.Wait(); err != nil {
    return nil, err
}
// 此处不能访问 err
```

## 与 JS 的关键区别

| | Go | JS |
|---|---|---|
| 括号 | 不用 | 必须 |
| 大括号 | 必须 | 可选（单行可省） |
| `else` 位置 | 与 `}` 同行 | 可折行 |
| init statement | 内置语法 | 无（可用 let 模仿但不常用） |
| 错误处理 | `if err := fn(); err != nil` | `try/catch` |

## 实战中最常见的模式

```go
// 错误检查（Go 里出现频率最高）
if err := doSomething(); err != nil {
    return fmt.Errorf("do something: %w", err)
}

// nil 检查
if result == nil {
    return nil, errors.New("result is nil")
}
```
