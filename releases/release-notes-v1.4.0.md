## 📦 Version 1.4.0 — Complete SQL Builder Suite

> Released: **2025-05-22**

---

### ✨ Highlights

- 🔧 Fully integrated builder suite:
  - `SelectBuilder`
  - `InsertBuilder`
  - `UpdateBuilder`
  - `DeleteBuilder`
  - `UpsertBuilder`
- ✅ Semantic-aware `NewCondition(...)` with operator, type, and placeholder resolution
- 💡 Dialect-safe rendering with support for PostgreSQL, generic (`?`, `$1`, etc.)
- 🔐 Full stage validation using `AddStageError(...)`
- 🧪 **100% test coverage** across all builders and helpers

---

### 🆕 Features

- New `ParamBinder` with offset-aware placeholder generation
- Dialect interface includes `SupportsReturning()` check
- Condition inference: parses `status = active`, `id IN ?`, `BETWEEN ? AND ?`
- Field alias rejection in INSERT and UPSERT for correctness

---

### 🧩 Dialect-Aware Logic

- `InsertBuilder` and `UpsertBuilder` align on internal placeholder binding
- `UpsertBuilder` supports `ON CONFLICT DO UPDATE` and `DO NOTHING`
- `RETURNING` clause only enabled if `dialect.SupportsReturning()` returns `true`

---

### ✅ Test Coverage

> 100% test coverage in:

- `select.go`, `insert.go`, `update.go`, `delete.go`, `upsert.go`
- `condition_helpers.go`, `condition_renderer.go`, `param_binder.go`

---

### 📌 Tag

```
v1.4.0
```
