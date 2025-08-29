# Table Token

## 🌱 Overview

The `token.Table` type represents a SQL table (or subquery) token with optional alias.  
It delegates parsing and classification to the shared **resolver** and **expression_kind** modules,  
enforcing strict construction rules and exposing multiple forms for logging, debugging, and SQL rendering.

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
   - Alias may also be any `fmt.Stringer` implementation  
   - Always validated via resolver

4. **Subquery**
   - `table.New("(SELECT COUNT(*) FROM users) AS t")` → subquery with alias  
   - `table.New("(SELECT COUNT(*) FROM users)", "t")` → subquery with alias  
   ⚠️ Subqueries **must have an alias**, otherwise the token is errored.

5. **Errors**
   - Empty input → errored
   - Invalid alias (including reserved keywords such as `AS`, `FROM`, `SELECT`) → errored
   - Passing a token directly (e.g. `table.New(table.New("users"))`) → errored, with hint to use `Clone()`
   - Invalid types (e.g. `table.New(123)`) → errored
   - Too many tokens or malformed input → errored

---

## Contracts Implemented

- **TableToken** → 
- **Clonable** → `Clone()` (safe duplication)
- **Debuggable** → `Debug()` (developer diagnostics with flags)
- **Errorable** → `IsErrored()`, `Error()`
- **Rawable** → `Raw()` (generic SQL fragment), `IsRaw()`
- **Renderable** → `Render()` (canonical SQL form)
- **Stringable** → `String()` (human-facing logs)
- **Validable** → `IsValid()` (validity check based on resolver rules)

---

## Logging

- **String()**  
  Concise, log-friendly:  
  ```text
  ✅ table(users AS u)
  ❌ table("users AS"): invalid format "users AS"
  ```

- **Debug()**  
  Verbose developer output with flags:  
  ```text
  ✅ table("users AS u"): [raw:false, aliased:true, errored:false]
  ❌ table("users AS"): [raw:false, aliased:false, errored:true] {err=invalid format "users AS"}
  ```

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project

