# 🔗 Developer Guide: ParamBinder

The `ParamBinder` is a utility for dialect-safe placeholder generation and parameter binding in SQL query builders.

---

## 🎯 Purpose

`ParamBinder` handles:

- Placeholder generation based on the target dialect (e.g., `$1`, `?`, `:field`)
- Safe tracking of bound arguments in the correct order
- Simplifies condition and query rendering logic without manual indexing

---

## 📦 Location

```go
/internal/core/driver/param.go
/internal/core/driver/param_ops.go
```

---

## 🔧 Constructor

```go
pb := driver.NewParamBinder(driver.NewPostgresDialect())
```

Starts with `position = 1` to align with 1-based dialects like PostgreSQL.

---

## 🔨 Usage

### Single value

```go
placeholder := pb.Bind("admin")  // → $1
args := pb.Args()                // → ["admin"]
```

### Multiple values

```go
placeholders := pb.BindMany(1, 2, 3)  // → ["$1", "$2", "$3"]
args := pb.Args()                     // → [1, 2, 3]
```

---

## 💡 Integration

Use `ParamBinder` in builders like `DeleteBuilder.Build()` or helpers like `FormatConditions` to ensure placeholder correctness without repeating position logic.

---

## ✅ Best Practices

- Always use a new `ParamBinder` per SQL statement
- Avoid mixing dialects mid-binding
- Let the builder own and resolve all placeholders before execution

---

## 🧪 Fully Tested

| Method      | Coverage |
|-------------|----------|
| `Bind()`    | ✅ 100%   |
| `BindMany()`| ✅ 100%   |
| `Args()`    | ✅ 100%   |
| Constructor | ✅ 100%   |

This makes `ParamBinder` a reliable and reusable backbone for all dialect-safe SQL binding operations.
