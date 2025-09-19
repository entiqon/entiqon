# Collection 📦

Generic, type-safe collections for Go.  
Part of the [`entiqon`](https://github.com/entiqon/entiqon) common/extension toolkit.

[![Go Reference](https://pkg.go.dev/badge/github.com/entiqon/common/extension/collection.svg)](https://pkg.go.dev/github.com/entiqon/common/extension/collection)
[![Go Report Card](https://goreportcard.com/badge/github.com/entiqon/common/extension/collection)](https://goreportcard.com/report/github.com/entiqon/common/extension/collection)
[![Tests](https://github.com/entiqon/entiqon/actions/workflows/test.yml/badge.svg)](https://github.com/entiqon/entiqon/actions)

---

## ✨ Features

- **Type-safe**: generic container for any `comparable` type  
- **Fluent API**: method chaining supported  
- **Functional helpers**: `Map`, `Filter`, `ForEach`  
- **Safe operations**: no panics on invalid indices (`At`, `RemoveAt`)  
- **Full test coverage**: each method verified with `file → method → cases`

---

## 📚 Quick Reference

| Category     | Method / Function                                          | Description                          |
|--------------|------------------------------------------------------------|--------------------------------------|
| Constructors | [`New`](#newt)                                             | Create a new empty collection        |
|              | [`FromSlice`](#fromslcet)                                  | Initialize from an existing slice    |
| Mutators     | [`Add`](#addvalues-t)                                      | Append values to the collection      |
|              | [`InsertAt`](#insertatidx-int-values-t)                    | Insert values at a specific index    |
|              | [`Remove`](#removevalue-t)                                 | Remove all occurrences of a value    |
|              | [`RemoveAt`](#removeatidx-int)                             | Remove value at a specific index     |
|              | [`Clear`](#clear)                                          | Remove all elements                  |
| Queries      | [`Contains`](#containsvalue-t)                             | Check if value exists                |
|              | [`IndexOf`](#indexofvalue-t)                               | Get index of first occurrence        |
|              | [`At`](#atidx-int-t-bool)                                  | Safe access by index                 |
|              | [`Length`](#length)                                        | Number of elements in collection     |
|              | [`Items`](#items-t)                                        | Copy of underlying slice             |
| Functional   | [`ForEach`](#foreachfn-funct)                              | Apply function to each element       |
|              | [`Filter`](#filterc-collectiont-fn-funct-bool-collectiont) | Create collection matching predicate |
|              | [`Map`](#mapc-collectiont-fn-funct-r-collectionr)          | Transform collection into new type   |
| Utilities    | [`Clone`](#clone)                                          | Create a deep copy of collection     |

---

## 📚 Table of Contents

- [Import](#-import)
- [Constructors](#-constructors)
- [Mutators](#-mutators)
- [Queries](#-queries)
- [Functional Helpers](#-functional-helpers)
- [Utilities](#-utilities)
- [Summary](#-summary)

---

## 🔹 Import

```go
import "github.com/entiqon/common/extension/collection"
```

---

## 🔹 Constructors

### `New[T]()`
Create a new empty collection.
```go
c := collection.New[int]()
fmt.Println(c.Length()) // 0
```

### `FromSlice[T](src []T)`
Initialize from an existing slice (copied internally).
```go
c := collection.FromSlice([]string{"a", "b", "c"})
fmt.Println(c.Items()) // ["a", "b", "c"]
```

---

## 🔹 Mutators

### `Add(values ...T)`
Append values to the end.
```go
c := collection.New[int]().Add(1, 2, 3)
fmt.Println(c.Items()) // [1 2 3]
```

### `InsertAt(idx int, values ...T)`
Insert values at a specific index.  
- Negative index → insert at start  
- Index beyond length → append at end
```go
c := collection.FromSlice([]int{1, 3})
c.InsertAt(1, 2)
fmt.Println(c.Items()) // [1 2 3]
```

### `Remove(value T)`
Remove **all occurrences** of a value.
```go
c := collection.FromSlice([]int{1, 2, 2, 3})
c.Remove(2)
fmt.Println(c.Items()) // [1 3]
```

### `RemoveAt(idx int)`
Remove element by index.  
Safe no-op if index is invalid.
```go
c := collection.FromSlice([]int{1, 2, 3})
c.RemoveAt(1)
fmt.Println(c.Items()) // [1 3]
```

### `Clear()`
Remove all elements.
```go
c := collection.FromSlice([]int{1, 2, 3})
c.Clear()
fmt.Println(c.Length()) // 0
```

---

## 🔹 Queries

### `Contains(value T)`
Check if value exists.
```go
c := collection.FromSlice([]string{"a", "b"})
fmt.Println(c.Contains("b")) // true
```

### `IndexOf(value T)`
Return index of first occurrence, or -1 if not found.
```go
c := collection.FromSlice([]int{10, 20, 30})
fmt.Println(c.IndexOf(20)) // 1
fmt.Println(c.IndexOf(99)) // -1
```

### `At(idx int) (T, bool)`
Safe access by index.  
Returns `(zero, false)` if index is invalid.
```go
c := collection.FromSlice([]int{10, 20})
val, ok := c.At(1)   // val=20, ok=true
val2, ok2 := c.At(5) // val2=0, ok2=false
```

### `Length()`
Return number of elements.  
```go
c := collection.FromSlice([]int{1, 2, 3})
fmt.Println(c.Length()) // 3
```

### `Items() []T`
Return a **copy** of underlying slice.  
Mutating result does not affect collection.
```go
c := collection.FromSlice([]int{1, 2})
items := c.Items()
items[0] = 99
fmt.Println(c.Items()) // [1 2] (unchanged)
```

---

## 🔹 Functional Helpers

### `ForEach(fn func(T))`
Apply function to each element in order.
```go
c := collection.FromSlice([]int{1, 2, 3})
sum := 0
c.ForEach(func(v int) { sum += v })
fmt.Println(sum) // 6
```

### `Filter(c *Collection[T], fn func(T) bool) *Collection[T]`
Create a new collection with only values matching predicate.
```go
c := collection.FromSlice([]int{1, 2, 3, 4})
evens := collection.Filter(c, func(x int) bool { return x%2 == 0 })
fmt.Println(evens.Items()) // [2 4]
```

### `Map(c *Collection[T], fn func(T) R) *Collection[R]`
Transform collection into another type.
```go
c := collection.FromSlice([]int{1, 2, 3})
strs := collection.Map(c, func(x int) string { return fmt.Sprintf("#%d", x) })
fmt.Println(strs.Items()) // ["#1" "#2" "#3"]
```

---

## 🔹 Utilities

### `Clone()`
Create a deep copy of collection (new slice).
```go
c1 := collection.FromSlice([]int{1, 2, 3})
c2 := c1.Clone().Add(4)

fmt.Println(c1.Items()) // [1 2 3]
fmt.Println(c2.Items()) // [1 2 3 4]
```

---

## 📌 Summary

- **Constructors** → `New`, `FromSlice`  
- **Mutators** → `Add`, `InsertAt`, `Remove`, `RemoveAt`, `Clear`  
- **Queries** → `Contains`, `IndexOf`, `At`, `Length`, `Items`  
- **Functional** → `ForEach`, `Filter`, `Map`  
- **Utilities** → `Clone`  

The package is **safe, type-checked, chainable, and functional** — ideal for expressive business/domain code without reinventing collection logic.
