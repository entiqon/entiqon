# Float Parser 🌊

Utilities to parse various types into `float64`.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- Convert values (`string`, `int`, `float`, `bool`) into `float64`
- Support for pointers and interface wrappers
- No rounding — values are preserved as-is
- Safe error handling for unsupported/invalid inputs
- Full test coverage (file → methods → cases)

---

## 📑 API Reference

### `ParseFrom(value any) (float64, error)`

Convert supported inputs into `float64`.

#### Supported inputs:
- **int / int8 / int16 / int32 / int64**
- **uint / uint8 / uint16 / uint32 / uint64 / uintptr**
- **float32 / float64**
- **string** → parsed as float
- **bool** → `true → 1.0`, `false → 0.0`
- **pointers/interfaces** wrapping the above

#### Unsupported:
- `nil`
- unsupported structs, slices, maps, complex, etc.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/entiqon/common/extension/float"
)

func main() {
    f1, _ := float.ParseFrom("123.456")
    fmt.Println(f1)

    f2, _ := float.ParseFrom(42)
    fmt.Println(f2)

    f3, _ := float.ParseFrom(true)
    fmt.Println(f3)
}
```

Output:
```
123.456
42
1
```

---

## 📌 Summary

- **Core:** `ParseFrom(any) (float64, error)`  
- **No rounding:** values preserved  
- **Safe:** clear errors on invalid/unsupported inputs
