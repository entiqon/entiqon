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

Example usage:

```go
import "github.com/entiqon/entiqon/common"

func example() {
    status := common.BoolToStr(true, "enabled", "disabled")
    fmt.Println(status) // Output: enabled
}
```

---

## 📘 Developer Guides

### 📚 Common

- [ProcessStage](guides/ProcessStage_Developer_Guide.md) — Dialects, StageToken, ParamBinder integration


---

## 📄 License

[MIT](../../../LICENSE) — © Isidro Lopez / Entiqon Project
