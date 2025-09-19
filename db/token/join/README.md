# 🔗 Join

> Part of [Entiqon](../../../) / [Database](../../) / [Token](../)

## 📜 User Guide

`join.Token` represents a **SQL JOIN clause** between two tables or subqueries,
with an explicit join kind (`INNER`, `LEFT`, `RIGHT`, `FULL`) and a condition.

### Instantiation Rules

1. **Safe constructors** → explicit kind

   * **Inner Join**

     ```go
     j := join.NewInner("users", "orders", "users.id = orders.user_id")
     // renders: INNER JOIN orders ON users.id = orders.user_id
     ```

   * **Left Join**

     ```go
     j := join.NewLeft("products", "categories", "products.cat_id = categories.id")
     // renders: LEFT JOIN categories ON products.cat_id = categories.id
     ```

   * **Right Join**

     ```go
     j := join.NewRight("employees", "departments", "employees.dep_id = departments.id")
     // renders: RIGHT JOIN departments ON employees.dep_id = departments.id
     ```

   * **Full Join**

     ```go
     j := join.NewFull("a", "b", "a.id = b.id")
     // renders: FULL JOIN b ON a.id = b.id
     ```

2. **Flexible constructor** → dynamic kind (advanced usage)

   ```go
   j := join.New("LEFT", "users u", "orders o", "u.id = o.user_id")
   // renders: LEFT JOIN orders o ON u.id = o.user_id

   j = join.New(join.RightJoin, "x", "y", "x.id = y.id")
   // renders: RIGHT JOIN y ON x.id = y.id
   ```

   > ⚠️ If the kind is invalid (unsupported string or enum), `New` returns an errored join:
   >
   > ```
   > ⛔️ Join(): invalid join type (99)
   > ```

---

## 📚 Developer Guide

### Internal Representation

A `Join` preserves:

* **Kind** → `join.Kind` enum (`InnerJoin`, `LeftJoin`, …)
* **Left / Right** → `table.Token` (may be raw strings or parsed tables)
* **Condition** → ON expression

All members are **unexported** to enforce immutability and are exposed only
through contract methods (`Kind()`, `Left()`, `Right()`, `Condition()`, …).

### Validation Rules

* Kind must be valid (`IsValid()`).
* Both left and right tables must be present.
* Left/Right tables must not be errored.
* Condition must not be empty.

Violations produce explicit error states.

---

## 🐞 Debugging and Logging

* **`String()`** → concise log/SQL view.

  * ✅ valid join:

    ```
    INNER JOIN orders ON users.id = orders.user_id
    ```

  * ⛔️ invalid join:

    ```
    ⛔️ Join(): invalid join type (99)
    ```

* **`Debug()`** → detailed diagnostic view.

  * ✅ valid join:

    ```
    Join(Kind=INNER JOIN, Left=users, Right=orders, Condition="users.id = orders.user_id")
    ```

  * ⛔️ invalid join:

    ```
    ⛔️ Join(Kind=invalid join type (99)): left=nil, right=nil, condition="id = 1"
    ```

---

## ✅ Contracts

`Join` implements the shared **contracts**:

* `Clonable` → safe deep copies
* `Errorable` → explicit error state
* `Debuggable` → detailed diagnostics
* `Rawable` → SQL-generic rendering (`Raw()`)
* `Renderable` → dialect-agnostic SQL fragment (`Render()`)
* `Stringable` → concise logging (`String()`)
* `Validable` → explicit validity checks (`IsValid()`)

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
