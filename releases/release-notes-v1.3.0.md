## 📦 Version 1.3.0— Dialect Integration + ParamBinder

> Released: **2025-05-xx**

---

### ✨ Highlights

- Introduced `ParamBinder` for dialect-aware argument binding
- PostgreSQL support with `$N` placeholders
- Added `driver.Dialect` interface with `QuoteIdentifier()` and `BuildLimitOffset()`
- Refactored condition handling using `ConditionType` enum
- Added `RenderConditions()` with AND/OR support

---

### 🧱 Improvements

- `SelectBuilder`, `UpdateBuilder`, and `DeleteBuilder` now support:
  - AND / OR where chaining
  - Quote-aware column and table rendering
- InsertBuilder:
  - `BuildInsertOnly()` introduced for Upsert reuse
- Condition validation added with `.IsValid()` and `.Error`

---

### 🧪 Testing

- Added integration tests for dialect placeholder rendering
- Partial coverage of core builders and renderers

---