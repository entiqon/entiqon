# Decimal Parser 💯

Utilities to parse values into `float64` and apply decimal rounding.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- Convert values (`string`, `int`, `float`, `bool`) into `float64`
- Configurable precision:
  - Negative → no rounding, raw float returned
  - `0` → round to integer
  - Positive → round to N decimal places
- Uses standard `math.Round` rules
- Safe error handling for invalid/unsupported inputs
- Full test coverage

---

## 📑 API Reference

### `ParseFrom(value any, precision int) (float64, error)`

Convert supported inputs into `float64` with optional rounding.

#### Supported inputs:
- **int / int8 / int16 / int32 / int64**
- **uint / uint8 / uint16 / uint32 / uint64**
- **float32 / float64**
- **string** → parsed as float
- **bool** → `true → 1.0`, `false → 0.0`

#### Unsupported:
- `nil`, structs, maps, slices, complex, etc.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/entiqon/common/extension/decimal"
)

func main() {
    v1, _ := decimal.ParseFrom("123.456789", 2)
    fmt.Println(v1)

    v2, _ := decimal.ParseFrom(42, 0)
    fmt.Println(v2)

    v3, _ := decimal.ParseFrom(true, 3)
    fmt.Println(v3)

    _, err := decimal.ParseFrom("invalid", 2)
    fmt.Println(err != nil)
}
```

Output:
```
123.46
42
1
true
```

---

## 📌 Summary

- **Core:** `ParseFrom(any, precision int) (float64, error)`
- **Flexible:** supports multiple types, with or without rounding
- **Safe:** clear errors on invalid inputs
