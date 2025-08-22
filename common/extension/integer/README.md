# Integer Parser

Package **integer** provides utilities to convert arbitrary values into `int` values.

Unlike `float` or `decimal`, this parser **always truncates** toward zero.  
No fractional values are preserved.

---

## Features

- Supports:
  - Signed and unsigned integers
  - Floats (fractional part truncated)
  - Strings (parsed as float then truncated)
  - Bools (`true` → 1, `false` → 0)
  - Pointers and interfaces wrapping supported types
- Rejects unsupported types with clear errors
- Safe: never panics

---

## Installation

```go
import "github.com/entiqon/entiqon/common/extension/integer"
```

---

## Usage

```go
package main

import (
    "fmt"
    "github.com/entiqon/entiqon/common/extension/integer"
)

func main() {
    v1, _ := integer.ParseFrom(42)
    fmt.Println(v1) // 42

    v2, _ := integer.ParseFrom(3.99)
    fmt.Println(v2) // 3

    v3, _ := integer.ParseFrom("123")
    fmt.Println(v3) // 123

    v4, _ := integer.ParseFrom(true)
    fmt.Println(v4) // 1
}
```

---

## Notes

- **Truncation**: `3.99 → 3`, `-1.9 → -1`, `0.5 → 0`.
- Returns `(0, error)` if input is `nil`, empty string, or unsupported type.
