# Date Parser 📅

Robust parsing and formatting utilities for Go `time.Time` values.

Part of the [`entiqon`](https://github.com/entiqon/entiqon) `common/extension` toolkit.

---

## ✨ Features

- Parse strings and interfaces into `time.Time`
- Clean and normalize date inputs before parsing
- Support for default layouts (e.g. `2006-01-02`, `20060102`, `02-01-2006`)
- Customizable formatters with `ParseAndFormat`
- Error handling for invalid and incomplete inputs
- Full test coverage (file → methods → cases)

---

## 📑 API Reference

### `ParseFrom(value any) (time.Time, error)`

Attempt to parse supported values into `time.Time`.

- **string** → parsed against known date layouts after cleaning
- **[]byte** → parsed as string
- **time.Time** → returned directly
- **unsupported/nil** → returns error

### `ParseAndFormat(value any, layout string) (string, error)`

Parse an input into a time, then format with a given Go layout.

### `Cleaning(raw string) string`

Normalize a raw date string by trimming, removing spaces, replacing separators, etc.

---

## 🔹 Usage Examples

```go
package main

import (
    "fmt"
    "github.com/entiqon/common/extension/date"
)

func main() {
    t, _ := date.ParseFrom("2025-08-21")
    fmt.Println(t.Format("2006-01-02"))

    out := date.ParseAndFormat("21/08/2025", "2006-01-02")
    fmt.Println(out)
}
```

Output:
```
2025-08-21
2025-08-21
```

---

## 📌 Summary

- **Core:** `ParseFrom(any) (time.Time, error)`  
- **Formatting:** `ParseAndFormat(any, layout string) (string, error)`  
- **Cleaning:** `Cleaning(string) string`  
