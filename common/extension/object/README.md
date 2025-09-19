# Object Utilities 🧩

Reflection-based helpers for dynamic object access and manipulation.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- **Exists** → check if a property or key exists  
- **GetValue** → safe typed retrieval with defaults  
- **SetValue** → update fields/keys on maps and structs  
- Supports:
  - `map[string]any` (case-insensitive keys)
  - structs (exported fields + zero-arg methods)
  - pointers to structs
- Case-insensitive field/key lookup
- Handles `nil`, unsupported types gracefully with safe errors

---

## 📑 API Reference

### `Exists(object any, key string) bool`
Check if a key/field exists in `map[string]any` or struct (exported fields only).

### `GetValue[T any](object any, key string, defaultVal T) T`
Retrieve property `key` as type `T`. Returns `defaultVal` if missing or type mismatch.

### `SetValue[O any, T any](object O, key string, value T) (O, error)`
Set a property on `map[string]any` or struct pointer.  
Returns updated object or error.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/common/extension/object"
)

type Item struct {
    ID   int
    Name string
}

func main() {
    // Map example
    m := map[string]any{"Foo": 123}
    fmt.Println(object.Exists(m, "foo")) // true

    val := object.GetValue[int](m, "foo", 0)
    fmt.Println(val) // 123

    m, _ = object.SetValue(m, "Bar", "baz")
    fmt.Println(m["Bar"]) // baz

    // Struct example
    item := &Item{ID: 10, Name: "Book"}
    fmt.Println(object.Exists(item, "Name")) // true

    name := object.GetValue[string](item, "Name", "default")
    fmt.Println(name) // Book

    item, _ = object.SetValue(item, "Name", "Notebook")
    fmt.Println(item.Name) // Notebook
}
```

---

## 📌 Summary

- **Exists**: quick presence check  
- **GetValue**: safe retrieval with default fallback  
- **SetValue**: update fields/keys dynamically  
