<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🧩 `common/math/float`

The **`common/math/float`** package provides utilities for parsing input values into `float64` without rounding.

| Function    | Signature                                            | Description                                                                                                                 |
|-------------|------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| `ParseFrom` | `func ParseFrom(value interface{}) (float64, error)` | Parses various input types (string, numeric, bool) into float64, returning errors on unsupported types or parsing failures. |

---

### 🧩 `common/math/decimal`

The **`common/math/decimal`** package builds upon `common/math/float` to provide decimal rounding utilities.

| Function    | Signature                                                           | Description                                                                                    |
|-------------|---------------------------------------------------------------------|------------------------------------------------------------------------------------------------|
| `ParseFrom` | `func ParseFrom(value interface{}, precision int) (float64, error)` | Parses input and rounds to specified decimal places; if precision < 0, no rounding is applied. |

---

### 🧩 `common/math/number`

The **`common/math/number`** package focuses on integer parsing from flexible input types.

| Function    | Signature                                                    | Description                                                                                                                   |
|-------------|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `ParseFrom` | `func ParseFrom(value interface{}, round bool) (int, error)` | Parses input into an int, with optional rounding for floats. When round=false, floats must be near integers or parsing fails. |

---

#### Example usage

```go
import (
    "fmt"
    
    "github.com/entiqon/entiqon/common/math/float"
    "github.com/entiqon/entiqon/common/math/decimal"
    "github.com/entiqon/entiqon/common/math/number"
)

func main() {
    f, err := float.ParseFrom("123.456")
    fmt.Println(f, err) // 123.456 <nil>
    
    d, err := decimal.ParseFrom("123.456789", 3)
    fmt.Println(d, err) // 123.457 <nil>
    
    i, err := number.ParseFrom("123.6", true)
    fmt.Println(i, err) // 124 <nil>
}
```

---

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project
