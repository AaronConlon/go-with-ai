# Go Variable Syntax

## Quick Distinction

| Syntax | Meaning |
| --- | --- |
| `var x int` | Declare a variable with an explicit type |
| `const version = "0.1.0"` | Declare a constant |
| `x := 1` | Short variable declaration inside a function |
| `x = 2` | Assign a new value to an existing variable |

## Practical Rule

Use `:=` inside functions when creating a new local variable. Use `var` when you need an explicit zero value, package-level declaration, or clearer type.

## Why Strings Use Double Quotes

In Go, strings use double quotes or backticks:

```go
t.Fatal("expected error, got nil")
```

or:

```go
json := `{"id":101,"title":"ok"}`
```

Single quotes do not mean string. They mean `rune`, which is one Unicode character:

```go
letter := 'a'
```

This is invalid:

```go
t.Fatal('expected error, got nil') // illegal rune literal
```

The single quotes contain many characters, but a rune literal can only represent one character.

`gofmt` will not automatically fix this. It formats Go code that already parses. `illegal rune literal` is a parse error, so `gofmt` does not have a valid AST to format.

Also, changing single quotes to double quotes can change the type:

```go
'a' // rune
"a" // string
```

Those are different Go types.
