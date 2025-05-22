## 📦 Version 1.2.0 — Core Builder Foundations

> Released: **2025-04-xx**

---

### ✨ Highlights

- Introduced initial builder architecture for `SelectBuilder`, `InsertBuilder`, and `UpdateBuilder`
- Early support for `.Where()`, `.From()`, `.Columns()`, and `.Values()` chaining
- Basic placeholder binding using inline `?` support
- Started modular token handling for column rendering

---

### 🧱 Components

- `SelectBuilder`: columns, table, where, limit/offset
- `InsertBuilder`: basic insert with value enforcement
- `UpdateBuilder`: set and where clauses only
- `BaseBuilder`: common error and dialect handling logic

---

### ❗ Limitations

- No semantic condition resolution
- No type enforcement or dialect-safe binding
- No alias validation
- No `RETURNING` support or upsert logic

---