<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true" align="left" height="96" width="96"> Errors
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a>::<span>Common</span> toolkit.</h6>

## 🌱 Overview

The `errors` package provides extended error handling utilities for Entiqon.
It builds on Go's standard `error` interface by introducing structured
errors that carry additional context.

---

## ✨ Features

- **CausableError**  
  An error that has both:
  - `Cause()` → short identifier for the category of error (e.g., `"Database"`)
  - `Reason()` → detailed explanation for logs or user messages

- **ProcessStageError**  
  Extends `CausableError` by associating an error with a `ProcessStage`.
  Useful for tracing errors in multi-step processing pipelines.

---

## 🚀 Quick Start

```go
import "github.com/entiqon/entiqon/common/errors"

// Create a simple causable error
err := errors.NewCausableError("Database", "Connection failed")

// Create a process stage error
pse := errors.NewProcessStageError("Init", "Loader", "Failed to load resource", err)
```

---

## 📘 Guides

- Use `Cause()` to classify errors programmatically
- Use `Reason()` or `Error()` for logging and display
- Use `ProcessStageError.Stage()` to trace where the error occurred
- Errors support Go's `errors.Is` and `errors.As` for unwrapping

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Team
