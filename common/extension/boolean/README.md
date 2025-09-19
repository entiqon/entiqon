# Boolean Parser ✅

Utilities to parse and format boolean values.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- Parse values (`string`, `int`, `float`, `bool`) into `bool`
- Accepts extended string tokens: `true/false`, `1/0`, `yes/no`, `on/off`, `y/n`, `t/f`
- Case-insensitive, whitespace-tolerant
- Convert booleans into custom string representations with `BoolToString`
- Deprecated alias `BoolToStr` for backwards compatibility
- Safe error handling for unsupported inputs
- Full test coverage

---

## 📑 API Reference

### `ParseFrom(value any) (bool, error)`
Parse a supported value into a `bool`.

### `BoolToString(b bool, trueStr, falseStr string) string`
Return `trueStr` if `b` is true, otherwise `falseStr`.

### `BoolToStr(b bool, trueStr, falseStr string) string`
Deprecated alias of `BoolToString`.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/common/extension/boolean"
)

func main() {
    v1, _ := boolean.ParseFrom("yes")
    fmt.Println(v1)

    v2 := boolean.BoolToString(false, "enabled", "disabled")
    fmt.Println(v2)
}
```

Output:
```
true
disabled
```

---

## 📌 Summary

- **Parsing:** `ParseFrom(any) (bool, error)`
- **Formatting:** `BoolToString(b, trueStr, falseStr) string`
- **Deprecated:** `BoolToStr` (use `BoolToString`)
