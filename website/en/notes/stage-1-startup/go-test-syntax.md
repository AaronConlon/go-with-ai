# Go Test Syntax

## Minimal Test Shape

```go
func TestSomething(t *testing.T) {
    got := Something()
    if got != want {
        t.Fatalf("got %v, want %v", got, want)
    }
}
```

## Key Concepts

- Test files end with `_test.go`.
- Test functions start with `Test`.
- `*testing.T` lets the test report failure.
- `t.Fatalf` stops the current test immediately.
- Table-driven tests keep multiple cases in one test function.

## Function Calls

Go function calls must match the function signature. If a function defines one argument, the call must pass one argument:

```go
got := IsValidStory(tt.story)
```

Go does not turn missing arguments into `undefined`, and it does not automatically fill missing values with `nil`.

Too few or too many arguments fail at compile time:

```go
got := IsValidStory()
// not enough arguments in call to IsValidStory

got := IsValidStory(tt.story, true)
// too many arguments in call to IsValidStory
```

When a value is intentionally absent, pass an explicit value that matches the parameter type. `nil` is only valid for nil-capable types such as slices, maps, channels, functions, pointers, and interfaces.

## Common Format Verbs

`t.Fatalf`, `fmt.Printf`, and `fmt.Errorf` use format verbs. A verb starts with `%`.

| Verb | Good for | Example |
| --- | --- | --- |
| `%s` | strings | `"hello"` |
| `%d` | decimal integers | `3`, `101` |
| `%v` | default format | `true`, `some error` |
| `%#v` | Go-syntax-like debug format | `[]hn.Item(nil)`, `hn.Item{ID:101}` |

These two messages use different verbs because they print different kinds of values:

```go
t.Fatalf("expected 3 items, got %d", len(items))
t.Fatalf("expected nil items on error, got %#v", items)
```

`len(items)` is an integer, so `%d` is right.

`items` is a slice, so `%#v` is better for debugging whether it is nil, empty, or contains elements.

## Stage 1 Rule

Prefer clear, heavily commented examples. The goal is learning the mechanics, not minimizing code.
