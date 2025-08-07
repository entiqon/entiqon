<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🌱 Overview

The **common** package provides foundational shared utilities, reusable components, and common helpers used across the
Entiqon project.

It enables code reuse, collaboration, and shared resource management for multiple modules and packages.

---

## 🚀 Quick Start

```bash
go get github.com/entiqon/entiqon
```

---

## 📦 Package Functions

### 🧩 `common`

| Function    | Signature                                                 | Description                                                                                         |
|-------------|-----------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `BoolToStr` | `func BoolToStr(b bool, trueStr, falseStr string) string` | Converts a boolean to one of two strings depending on its value, e.g., `"enabled"` or `"disabled"`. |

#### Example Usage

```go
import (
    "fmt"
    "github.com/entiqon/entiqon/common"
)

func main() {
    status := common.BoolToStr(true, "enabled", "disabled")
    fmt.Println("status:", status)
}
```

---

### 🧩 `common/math/float`

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

### 🧩 `common/object`

The **`common/object`** package provides utilities to manipulate dynamic key-value maps (`map[string]any`) as flexible
objects.

| Function   | Signature                                                                         | Description                                                                                                                                                    |
|------------|-----------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Exists`   | `func Exists(object map[string]any, key string) bool`                             | Checks whether the given key exists in the object.                                                                                                             |
| `GetValue` | `func GetValue[T any](object map[string]any, key string, defaultVal T) T`         | Returns the value for the key cast to type `T`.If the key is missing or the cast fails, returns `defaultVal`.                                                  |
| `SetValue` | `func SetValue[T any](object map[string]any, key string, value T) map[string]any` | Sets `value` into `object[key]` only if the key is missing or the existing value differs (deep equality).Initializes the map if `nil`.Returns the updated map. |

#### Example Usage

```go
import "github.com/entiqon/entiqon/common/object"

func main() {
    m := map[string]any{"count": 10}
    
    if object.Exists(m, "count") {
        fmt.Println("count exists")
    }
    
    count := object.GetValue[int](m, "count", 0)
    fmt.Println("count:", count)
    
    m = object.SetValue(m, "count", 20)
}

```

---

### 🧩 `common/errors`

The **`common/errors`** package provides extended error types with cause and reason semantics, plus process stage error
tracking.

| Type / Function        | Signature / Description                                                                | Notes                                                                   |
|------------------------|----------------------------------------------------------------------------------------|-------------------------------------------------------------------------|
| `CausableError`        | `interface { error; Cause() string; Reason() string }`                                 | Extended error interface providing cause and reason details             |
| `NewCausableError`     | `func NewCausableError(cause, reason string) CausableError`                            | Creates a new `CausableError` with cause and reason                     |
| `ProcessStageError`    | `struct { stage ProcessStage; cause, reason string; err error }`                       | Error with process stage, cause, reason, and optional wrapped error     |
| `NewProcessStageError` | `func NewProcessStageError(stage ProcessStage, cause, reason string, err error) error` | Creates a new `ProcessStageError` wrapping an optional underlying error |

### Example Usage

```go
import (
    "errors"
    "fmt"
    
    "github.com/entiqon/entiqon/common"
    "github.com/entiqon/entiqon/common/errors"
)

func main() {
    cause := "Database"
    reason := "Connection timeout"
    stage := common.ProcessStage("Init")
    wrappedErr := errors.New("network unreachable")
    
    err := errors.NewProcessStageError(stage, cause, reason, wrappedErr)
    
    fmt.Println(err.Error()) // [Database] at stage Init: Connection timeout: network unreachable
    
    if errors.Is(err, wrappedErr) {
        fmt.Println("Wrapped error detected")
    }
    
    var ce errors.CausableError
    if errors.As(err, &ce) {
        fmt.Printf("Cause: %s Reason: %s", ce.Cause(), ce.Reason())
    }
}
```

---

## 📘 Developer Guides

### 📚 Common

- [ProcessStage](guides/ProcessStage_Developer_Guide.md) — Dialects, StageToken, ParamBinder integration

---

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project
