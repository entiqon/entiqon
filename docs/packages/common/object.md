<h1><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="center" height="64" width="64"> Common Module</h1>

## 🧩 `common/object`

The **`common/object`** package provides utilities to manipulate dynamic key-value maps (`map[string]any`) and arbitrary structs as flexible objects.

| Function   | Signature                                                                        | Description                                                                                                                                                                                                                          |
|------------|---------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Exists`   | `func Exists(object map[string]any, key string) bool`                           | Checks if the given key exists in the object. Supports `map[string]any`.                                                                                                                                                           |
| `GetValue` | `func GetValue[T any](object any, key string, defaultVal T) T`                   | Returns the value for the given key cast to type `T`. Supports `map[string]any`, structs (including pointers to structs), and zero-argument getter methods. Returns `defaultVal` if the key is missing or cast fails.              |
| `SetValue` | `func SetValue[O any, T any](object O, key string, value T) (O, error)`         | Sets the value for the key in the object if missing or different from current value. Supports `map[string]any` and pointers to structs. For structs, finds the field (case-insensitive) and assigns safely via reflection with type conversion. Returns error if assignment fails. |

### Key Features

- Supports **dynamic maps** and **arbitrary structs** interchangeably for property access.
- Case-insensitive key and field name matching.
- Reflection-backed safe access to exported fields or zero-arg getter methods.
- Generic typed retrieval with fallback default values.
- Efficient map update logic with nil initialization.
- Detailed error handling for non-existent fields, unexported fields, incompatible types, and invalid objects.

### Example Usage

```go
import (
	"fmt"
	"github.com/entiqon/entiqon/common/object"
)

type ShipmentItem struct {
	LineNo int
	SKU    string
}

func main() {
	m := map[string]any{"LineNo": 10}
	s := ShipmentItem{LineNo: 20}

	if object.Exists(m, "LineNo") {
		fmt.Println("LineNo exists in map")
	}

	lineNoMap := object.GetValue[int](m, "LineNo", 0)
	fmt.Println("LineNo from map:", lineNoMap) // 10

	lineNoStruct := object.GetValue[int](s, "LineNo", 0)
	fmt.Println("LineNo from struct:", lineNoStruct) // 20

	// Set a new value in the map
	m, err := object.SetValue(m, "LineNo", 30)
	if err != nil {
		fmt.Println("Error setting value:", err)
	} else {
		fmt.Println("Updated LineNo in map:", m["LineNo"]) // 30
	}
}
```

---

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project