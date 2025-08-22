# Number Parser 🔢

Utilities to parse values into `int` with optional rounding.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- Convert values (`string`, `int`, `float`, `bool`) into `int`
- Overflow check for `uint64` → safe downcast to `int`
- Configurable float handling:
  - `round = true` → round to nearest integer
  - `round = false` → must be within `1e-9` tolerance of an integer
- Supports string parsing of integers and floats
- Safe error handling for invalid/unsupported inputs
- Full test coverage (file → methods → cases)

---

## 📑 API Reference

### `ParseFrom(value any, round bool) (int, error)`

Convert supported inputs into `int`, controlling float behavior.

#### Supported inputs:
- **int / int8 / int16 / int32 / int64**
- **uint / uint8 / uint16 / uint32 / uint64** (with overflow check)
- **float32 / float64**  
  - strict: require integer within tolerance (`round=false`)  
  - lenient: round to nearest integer (`round=true`)
- **string** → parsed as int, or float (with above rules)
- **bool** → `true → 1`, `false → 0`

#### Unsupported:
- `nil`, structs, maps, slices, complex, etc.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/entiqon/common/extension/number"
)

func main() {
    n1, _ := number.ParseFrom("123.6", true)
    fmt.Println(n1)

    n2, err := number.ParseFrom(123.4, false)
    fmt.Println(n2, err != nil)

    n3, _ := number.ParseFrom(true, false)
    fmt.Println(n3)
}
```

Output:
```
124
0 true
1
```

---

## 📌 Summary

- **Core:** `ParseFrom(any, round bool) (int, error)`  
- **Strict or lenient:** control float tolerance  
- **Safe:** overflow checks, clear errors
