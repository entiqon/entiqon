<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96"> token.Field
</h1>
<h6 align="left">Part of the <a href="https://github.com/entiqon/entiqon">Entiqon</a>::<span>Database</span> toolkit.</h6>


## 📜 User Guide

`token.Field` represents a **single column or expression** in a `SELECT` statement.  
The builder provides multiple ways to instantiate it, depending on what you want to express.

### Instantiation Rules

1. **Single string** → one or more fields
   - **Plain column**
     ```go
     .Fields("id")
     // SELECT id
     ```
   - **Comma-separated list**
     ```go
     .Fields("id, name, email")
     // SELECT id, name, email
     ```
   - **Aliased by space**
     ```go
     .Fields("id user_id")
     // SELECT id AS user_id
     ```
   - **Aliased by AS keyword**
     ```go
     .Fields("id AS user_id")
     // SELECT id AS user_id
     ```
   - **Function expression with alias**
     ```go
     .Fields("COUNT(id) AS row_count")
     // SELECT COUNT(id) AS row_count
     ```

2. **Two arguments (string, string)** → field + alias
   ```go
   .Fields("id", "user_id")
   // SELECT id AS user_id
   ```

3. **Three arguments (string, string, bool)** → raw expression with alias
   ```go
   .Fields("COUNT(*)", "total", true)
   // SELECT COUNT(*) AS total
   ```

4. **Multiple `Field` objects** → explicit choice
   ```go
   .Fields(
       *token.NewField("id"),
       *token.NewField("name"),
   )
   // SELECT id, name
   ```

### Things to Avoid

- Don’t pass `("id", "name")` unless you mean **alias** (`id AS name`).  
- Don’t try to concatenate raw SQL strings directly into `Fields` with multiple arguments unless you mean aliasing.  
- Prefer explicit `token.NewField(...)` if you want to avoid ambiguity.

---

## 📚 Developer Guide

### Internal Representation

A `Field` has the following structure:

```go
type Field struct {
    Input  string  // raw input as given by user
    Expr   string  // resolved expression (alias stripped)
    Alias  string  // optional alias
    IsRaw  bool    // true if raw expression
    Error  error   // set if instantiation failed
}
```

### Lifecycle

1. **Construction**  
   - Always done through `token.NewField(...)`.  
   - Handles all valid instantiation forms:
     - `NewField("expr")`
     - `NewField("expr alias")` or `NewField("expr AS alias")`
     - `NewField("expr", "alias")`
     - `NewField("expr", "alias", true)`
     - `NewField(*Field)` (not allowed, use `.Clone()` instead)

2. **Validation**  
   - If input type is unsupported, `Error` is set.  
   - If input is another `Field`, `Error` advises to use `.Clone()` instead.  

3. **Rendering**  
   - `Render()` → dialect-agnostic SQL fragment.  
   - `IsAliased()` → true if alias is set.  
   - `IsErrored()` → true if `Error` is non-nil.  
   - `IsValid()` → true if no error and Expr is set.  

4. **Cloning**  
   - `Clone()` returns a deep copy, preserving `nil` if original is nil.  
   - Prevents mutation of shared fields.

---

## 🐞 Debugging and Logging

Two methods are provided for inspection:

- **`String()`** → concise log view.  
  - ✅ valid field:  
    ```
    ✅ Field("id")
    ✅ Field("id AS user_id")
    ```
  - ⛔️ invalid field:  
    ```
    ⛔️ Field("true"): input type unsupported: bool
    ⛔️ Field(<nil>): wrong initialization
    ```

- **`Debug()`** → detailed diagnostic view.  
  - ✅ valid field:  
    ```
    ✅ Field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]
    ```
  - ⛔️ invalid field:  
    ```
    ⛔️ Field("false"): [raw: false, aliased: false, errored: true] – input type unsupported: bool
    ```

---

## ✅ Summary

- **Users**: treat `Field` as *one column or expression*. Use `.Fields(...)` safely with the instantiation rules.  
- **Contributors**: enforce immutability, strict parsing in `NewField`, and clear reporting through `String()` and `Debug()`.
