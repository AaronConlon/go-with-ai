# Go `if` Syntax

## Three Forms

### 1. Basic Condition

```go
if x > 0 {
    fmt.Println("positive")
}
```

- **No parentheses** around the condition (unlike JS)
- **Braces are required** (always, in all Go control flow)

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

- `else` **must be on the same line** as the closing `}`

### 3. With init statement (Go's unique feature)

```go
if err := g.Wait(); err != nil {
    return nil, err
}
```

The semicolon splits it into:

```
if   <init statement>  ;  <condition>  {
```

- Variables declared in the init statement are **scoped to the if/else block**

```go
// Shorter scope than separate declaration:
if err := fn(); err != nil {
    return nil, err
}
// err is not accessible here
```

## Most Common Pattern in Practice

```go
// Error check — the most frequent if usage in Go
if err := doSomething(); err != nil {
    return fmt.Errorf("do something: %w", err)
}

// nil check
if result == nil {
    return nil, errors.New("result is nil")
}
```
