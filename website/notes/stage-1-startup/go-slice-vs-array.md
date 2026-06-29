# Go 切片与数组

## 核心区别

| | 数组 `[N]T` | 切片 `[]T` |
|---|---|---|
| 长度 | **固定**，编译期确定，属于类型的一部分 | **动态**，运行时可变 |
| 传参行为 | 拷贝**整个数组**（值语义） | 拷贝**描述子**（24 字节），底层数据共享 |
| 可否 `append` | ❌ 编译错误 | ✅ 可以追加元素 |
| 能否为 `nil` | ❌ 不能，永远是值 | ✅ 可以，零值是 `nil` |

## 直观对比

```go
// 数组：长度固定 3，类型是 [3]int
arr := [3]int{101, 102, 103}

// 切片：长度可变的 []int
sl := []int{101, 102, 103}

arr = append(arr, 104)  // 编译错误
sl = append(sl, 104)    // 正常，sl 变为 [101, 102, 103, 104]
```

## 传参差异

```go
func modifyArray(a [3]int) { a[0] = 999 }  // 修改的是副本
func modifySlice(s []int)   { s[0] = 999 }  // 修改的是原数据

arr := [3]int{1, 2, 3}
sl := []int{1, 2, 3}

modifyArray(arr)   // arr 不变：{1, 2, 3}
modifySlice(sl)    // sl 变了：{999, 2, 3}
```

因为切片传参只拷贝描述子（指向同一块底层数组），所以函数内修改会反映到原切片。

## 切片的底层表示

切片内部是一个三元组：

```
┌──────────┐
│  data ptr │ → 指向底层数组某位置
│  len      │ → 当前元素个数
│  cap      │ → 容量（无需扩容能塞多少）
└──────────┘
```

`make([]int, 3, 5)` 创建长度为 3、容量为 5 的切片。

## `make` 函数是什么

`make` 是 Go 的内置函数，专门用来创建并初始化这三类类型：

- slice
- map
- channel

现阶段先把它理解成：

```text
把某种“容器”创建到可以直接使用的状态。
```

### 创建 slice

```go
items := make([]Item, len(ids))
```

这行表示：

```text
创建一个 []Item。
长度是 len(ids)。
每个位置先放 Item{} 零值。
```

如果：

```go
ids := []int64{101, 102, 103}
```

那么：

```go
items := make([]Item, len(ids))
```

会得到一个长度为 3 的 slice：

```text
items[0] == Item{}
items[1] == Item{}
items[2] == Item{}
```

后面可以按下标写入：

```go
items[0] = item0
items[1] = item1
items[2] = item2
```

`make` 创建 slice 时有两种常见写法：

```go
make([]Item, length)
make([]Item, length, capacity)
```

例如：

```go
make([]int, 3)
```

表示长度是 3，容量也是 3。

```go
make([]int, 3, 5)
```

表示长度是 3，容量是 5。

`len` 是当前已经能访问的元素数量，`cap` 是在扩容前最多能容纳的数量。

```go
s := make([]int, 3, 5)

fmt.Println(len(s)) // 3
fmt.Println(cap(s)) // 5
```

### 创建 map

map 的零值是 `nil`，不能直接写入：

```go
var scores map[string]int
scores["go"] = 100 // 会 panic
```

要先用 `make` 初始化：

```go
scores := make(map[string]int)
scores["go"] = 100
```

### 创建 channel

阶段三会用 channel 做并发许可证：

```go
sem := make(chan struct{}, concurrency)
```

这行表示：

```text
创建一个 channel。
里面传递 struct{}。
缓冲区大小是 concurrency。
```

也就是最多能同时放入 `concurrency` 个许可证。

### `make` 和 `new` 有什么不同

先记一个粗略区别：

| 写法 | 主要用途 | 返回什么 |
| --- | --- | --- |
| `make([]Item, n)` | 创建 slice、map、channel | 返回可用的值本身 |
| `new(Item)` | 给任意类型分配零值内存 | 返回指针 `*Item` |

日常写 slice、map、channel 时，用 `make`。

```go
items := make([]Item, len(ids))
scores := make(map[string]int)
sem := make(chan struct{}, concurrency)
```

创建 struct 值时，通常用字面量：

```go
item := Item{}
client := &Client{}
```

## 什么时候用哪个

- **切片**：日常默认选择。用于函数参数、集合操作、API 返回值。
- **数组**：极少数场景。固定长度的 hash 键、控制内存布局、或明确不需要变长。

对比 JS：JS 的 `Array` 对应 Go 的 `[]T`（切片），Go 的 `[N]T`（数组）在 JS 中没有直接等价物。

## 相关阅读

- [Go 测试代码语法拆解](/notes/stage-1-startup/go-test-syntax)
- [Go 变量定义语法](/notes/stage-1-startup/go-variable-syntax)
