<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🧩 `common/errors`

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

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project
