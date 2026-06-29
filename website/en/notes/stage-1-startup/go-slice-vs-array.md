# Go: Slice vs Array

## Key Difference

| | Array `[N]T` | Slice `[]T` |
|---|---|---|
| Length | **Fixed** at compile time, part of the type | **Dynamic**, grows at runtime |
| Pass by | **Copies the whole array** (value semantics) | **Copies the descriptor** (24 bytes), shares underlying data |
| `append` | ❌ Compile error | ✅ Yes |
| Can be `nil` | ❌ No, always a value | ✅ Yes, zero value is `nil` |

## Visual Comparison

```go
arr := [3]int{101, 102, 103}  // array, type is [3]int
sl  := []int{101, 102, 103}   // slice, type is []int

arr = append(arr, 104)  // compile error
sl  = append(sl, 104)   // OK, sl becomes [101, 102, 103, 104]
```

## Passing to Functions

```go
func modifyArray(a [3]int) { a[0] = 999 }  // modifies a copy
func modifySlice(s []int)   { s[0] = 999 }  // modifies original data

arr := [3]int{1, 2, 3}
sl  := []int{1, 2, 3}

modifyArray(arr)   // arr unchanged: {1, 2, 3}
modifySlice(sl)    // sl changed: {999, 2, 3}
```

Because a slice descriptor is small (pointer + len + cap), passing it is cheap—but modifications inside a function affect the caller's data.

## Slice Internals

A slice is a 3-field descriptor:

```
┌──────────┐
│ data ptr │ → pointer to an element of the backing array
│ len      │ → number of elements
│ cap      │ → capacity (max elements without reallocation)
└──────────┘
```

`make([]int, 3, 5)` creates a slice with len=3, cap=5.

## What `make` Does

`make` is a Go built-in used to create and initialize:

- slice
- map
- channel

For a slice:

```go
items := make([]Item, len(ids))
```

Read it as:

```text
Create a []Item.
Its length is len(ids).
Each slot starts as Item{}.
```

Common slice forms:

```go
make([]Item, length)
make([]Item, length, capacity)
```

For a map:

```go
scores := make(map[string]int)
scores["go"] = 100
```

For a channel:

```go
sem := make(chan struct{}, concurrency)
```

Rough distinction:

| Form | Use | Returns |
| --- | --- | --- |
| `make([]Item, n)` | slice, map, channel | usable value |
| `new(Item)` | allocate zero value memory | pointer `*Item` |

## When to Use Which

- **Slice** — default choice. Use for function parameters, collections, API return values.
- **Array** — rare. Fixed-size hash keys, explicit memory layout, or when you really want value semantics.

JS comparison: JavaScript `Array` maps to Go `[]T` (slice). Go `[N]T` (array) has no direct JS equivalent.
