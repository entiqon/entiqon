<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🧩 `common`

| Package | Utility     | Signature                                                 | Description                                                                                         |
|---------|-------------|-----------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| Common  | `BoolToStr` | `func BoolToStr(b bool, trueStr, falseStr string) string` | Converts a boolean to one of two strings depending on its value, e.g., `"enabled"` or `"disabled"`. |

### Example Usage

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

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project
