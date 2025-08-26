<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96" alt="entiqon"> Field
</h1>
<h6 align="left">Part of the <a href="../../../README.md">Entiqon</a> / <a href="../../README.md">Database</a> / <a href="../README.md">Token</a> toolkit.</h6>


<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="96" width="96" alt="entiqon"> Field
</h1>
<h6 align="left">Part of the <a href="../../../README.md">Entiqon</a> / <a href="../../README.md">Database</a> / <a href="../README.md">Token</a> toolkit.</h6>

## 📜 User Guide

`token.Field` represents a **single column, computed expression, or subquery** in a SQL statement.
The builder provides multiple ways to instantiate it, depending on what you want to express.

### Instantiation Rules

1. **Single string** → one field

    * **Plain column**

      ```go
      f := field.New("id")
      // renders: id
      ```
    * **Aliased by space**

      ```go
      f := field.New("id user_id")
      // renders: id AS user_id
      ```
    * **Aliased by AS keyword**

      ```go
      f := field.New("id AS user_id")
      // renders: id AS user_id
      ```
    * **Computed expression (functions, arithmetic)**

      ```go
      f := field.New("SUM(qty * price) total")
      // renders: SUM(qty * price) AS total
 
      f = field.New("qty * price")
      // renders: qty * price AS expr_alias_xxxxx
      ```
    * **Subquery**

      ```go
      f := field.New("(SELECT id FROM users) u")
      // renders: (SELECT id FROM users) AS u
 
      f = field.New("(SELECT id FROM users)")
      // renders: (SELECT id FROM users) AS expr_alias_xxxxx
      ```

2. **Two arguments (string, string)** → expr + alias

   ```go
   f := field.New("id", "user_id")
   // renders: id AS user_id
   ```

3. **Three arguments (string, string, bool)** → expr + alias + isRaw

   ```go
   f := field.New("COUNT(*)", "total", true)
   // renders: COUNT(*) AS total
   ```

> ℹ️ **Note:** Comma-separated lists like `"id, name, email"` are supported at the **SelectBuilder** level, **not** by `Field`.

---

## 📚 Developer Guide

### Internal Representation

A `Field` preserves the original input and provides strict parsing into components.
All internal members are kept **unexported** to enforce immutability.
They are only accessible through contract methods (`Input()`, `Expr()`, `Alias()`, etc.).

---

## 🐞 Debugging and Logging

Two methods are provided for inspection:

* **`String()`** → concise log/SQL view.

    * ✅ valid field:

      ```
      id
      id AS user_id
      ```
    * ⛔️ invalid field:

      ```
      ⛔️ Field(""): empty expression is not allowed
      ```

* **`Debug()`** → detailed diagnostic view.

    * ✅ valid field:

      ```
      Field(Input="COUNT(*) AS total", Expr="COUNT(*)", Alias="total", Raw=true, Err=<nil>)
      ```
    * ⛔️ invalid field:

      ```
      ⛔️ Field("false"): [raw: false, aliased: false, errored: true] – input type unsupported: bool
      ```

---

## ✅ Contracts

`Field` implements the shared **token contracts**:

* `BaseToken` → input, expr, alias, validity
* `Errorable` → explicit error state
* `Clonable` → safe deep copies
* `Rawable` → SQL-generic rendering (expr, alias, owner)
* `Renderable` → dialect-agnostic SQL fragment (`String()`)
* `Stringable` → concise logging
* `Ownerable` → ownership binding (`HasOwner`, `Owner`, `SetOwner`)

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project