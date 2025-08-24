<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="82" width="82" alt=""> Table
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a> / <a href="https://github.com/entiqon/entiqon/tree/main/db">Database</a> / <a href="https://github.com/entiqon/entiqon/tree/main/db/token">Token</a> toolkit.</h6>

## 🌱 Overview

The `token.Table` type represents a SQL table (or subquery) token with optional alias.
It enforces strict rules for construction and exposes multiple forms for logging,
debugging, and SQL rendering.

---

## Construction Rules

Tables are created using `table.New(...)`:

1. **Plain table**
   - `table.New("users")` → `users`

2. **Aliased (inline)**
   - `table.New("users u")` → `users AS u`
   - `table.New("users AS u")` → `users AS u`

3. **Aliased (explicit arguments)**
   - `table.New("users", "u")` → `users AS u`
   - Always marked as `isRaw=true`.

4. **Subquery**
   - `table.New("(SELECT COUNT(*) FROM users) AS t")` → subquery with alias  
   - `table.New("(SELECT COUNT(*) FROM users)", "t")` → subquery with alias, `isRaw=true`  
   ⚠️ Subqueries **must have an alias**, otherwise the token is errored.

5. **Errors**
   - Empty input → errored
   - Invalid formats (e.g. `"users AS"`) → errored
   - Too many tokens → errored

---

## Contracts Implemented

- **Renderable** → `Render()` (canonical SQL form)
- **Rawable** → `Raw()` (generic SQL fragment), `IsRaw()`
- **Stringable** → `String()` (human-facing logs)
- **Debuggable** → `Debug()` (developer diagnostics with flags)
- **Clonable** → `Clone()` (safe duplication)
- **Errorable** → `IsErrored()`, `Error()`

---

## Logging

- **String()**  
  Concise, log-friendly:  
  ```
  ✅ Table(users AS u)
  ❌ Table("users AS"): invalid format "users AS"
  ```

- **Debug()**  
  Verbose developer output with flags:  
  ```
  ✅ Table("users AS u"): [raw:false, aliased:true, errored:false]
  ❌ Table("users AS"): [raw:false, aliased:false, errored:true] {err=invalid format "users AS"}
  ```

---

## Philosophy

- **Never panic** — always returns a `*Table`, even if errored.
- **Auditability** — preserves original input for logs.
- **Strict enforcement** — invalid inputs are rejected immediately.
- **Delegation** — parsing rules live in `table.New`, not in builders.

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project