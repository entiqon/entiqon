# Rendering Principles in Entiqon

## Why Dialect Is Injected into Token Rendering

Tokens like `Column`, `Table`, and `Condition` represent **structural SQL elements**.

They are intentionally designed to be **dialect-agnostic** — storing things like:
- Column expr
- Table alias
- Aliasing relationships

However, they do **not** store any information about how to render themselves in SQL syntax. That responsibility is delegated to the injected `Dialect`.

---

### ✨ Benefits of Injecting Dialect

1. ✅ Tokens can be reused across multiple database dialects (e.g., PostgreSQL, MySQL).
2. ✅ Rendering remains **stateless**, **testable**, and **composable**.
3. ✅ `Dialect` rules (quoting, placeholder style, aliasing) are centralized.

---

### 🔁 Example

```go
col := NewColumn("users.status AS s")

// In PostgreSQL
sql := col.Render(PostgresDialect{}) // → "users"."status" AS "s"

// In MySQL
sql := col.Render(MySQLDialect{})    // → `users`.`status` AS `s`
```

The same token structure yields different SQL output based on the rendering strategy injected at runtime.

---

### ❌ What Happens If Tokens Hold Dialect

- Tokens become tightly coupled to rendering logic
- Prevents reuse across dialects
- Breaks separation of concerns

---

### ✅ Conclusion

> Rendering is **always dialect-aware**, but tokens are **dialect-agnostic**.

This keeps Entiqon’s query system flexible, clean, and powerful.
