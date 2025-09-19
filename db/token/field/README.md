# Field Token

> Part of [Entiqon](../../../) / [Database](../../) / [Token](../)

## 🌱 Overview

The `token.Field` type represents a SQL field (column, expression, literal, function, or subquery) with an optional alias.  
It is built on top of the **BaseToken** contract and integrates with shared **identifier** helpers for strict validation and classification.  
`Field` is consumed by higher-level builders (e.g., `SelectBuilder`) to construct safe and expressive SQL statements.

---

## Construction Rules

Fields are created using `field.New(...)` or `field.NewWithTable(...)`:

1. **No argument**
   ```go
   f := field.New()
   // → errored token (empty input)
   ```

2. **Plain field**
   ```go
   f := field.New("id")
   // → id
   ```

3. **Aliased (inline)**
   ```go
   f := field.New("id user_id")
   // → id AS user_id

   f = field.New("id AS user_id")
   // → id AS user_id
   ```

4. **Aliased (explicit arguments)**
   ```go
   f := field.New("id", "user_id")
   // → id AS user_id
   ```
   - The second argument may also be any `fmt.Stringer`.
   - Aliases are validated via `identifier.IsValidAlias`.

5. **Wildcard**
   ```go
   f := field.New("*")
   // → *
   ```
   ⚠️ Wildcards (`*`) cannot be aliased. Using `field.New("* alias")` or `field.New("*", "alias")` produces an errored token.

6. **Subquery**
   ```go
   f := field.New("(SELECT COUNT(*) FROM users) AS total")
   // → (SELECT COUNT(* ) FROM users) AS total

   f = field.New(field.New("id"), "alias")
   // → id AS alias
   ```
   ⚠️ Subqueries **must have an alias**, otherwise the token is errored.

7. **Computed / Function / Literal**
   ```go
   f := field.New("price * quantity", "total") // Computed expression
   // → (price * quantity) AS total

   f = field.New("SUM(price)", "sum_price")    // Aggregate function
   // → SUM(price) AS sum_price

   f = field.New("'constant'", "label")        // Literal with alias
   // → 'constant' AS label
   ```

8. **Invalid cases**
   - Empty string → errored
   - Invalid alias (reserved keyword, bad format) → errored
   - Passing another token directly (e.g. `field.New(field.New("id"))`) → errored, with hint to use `Clone()`
   - Too many parts in input (e.g. `field.New("field alias extra")`) → errored
   - Wrong types (e.g. `field.New(123)`) → errored

---

## Contracts Implemented

- **BaseToken** → core identity (`Input()`, `Expr()`, `Alias()`, `IsAliased()`, `ExpressionKind()`)
- **Clonable** → `Clone()` (safe duplication)
- **Debuggable** → `Debug()` (developer diagnostics with flags)
- **Errorable** → `IsErrored()`, `Error()`
- **Rawable** → `Raw()` (generic SQL fragment), `IsRaw()`
- **Renderable** → `Render()` (dialect‑agnostic SQL form)
- **Stringable** → `String()` (human‑friendly logs)
- **Validable** → `IsValid()` (validity check via `identifier.Validate*`)

---

## Examples

### Example: New with plain field
```go
f := field.New("id")
fmt.Println(f.String())
// Output: field(id)
```

### Example: New with inline alias
```go
f := field.New("id AS user_id")
fmt.Println(f.String())
// Output: field(id AS user_id)
```

### Example: New with explicit alias
```go
f := field.New("id", "user_id")
fmt.Println(f.String())
// Output: field(id AS user_id)
```

### Example: Wildcard without alias
```go
f := field.New("*")
fmt.Println(f.String())
// Output: field(*)
```

### Example: Wildcard with alias (error)
```go
f := field.New("* AS alias")
fmt.Println(f.Error())
// Output: '* 'cannot be aliased or raw
```

### Example: Subquery with alias
```go
f := field.New("(SELECT COUNT(*) FROM users) AS t")
fmt.Println(f.Render(dialect.Postgres))
// Output: (SELECT COUNT(* ) FROM users) AS t
```

### Example: Computed expression
```go
f := field.New("price * quantity", "total")
fmt.Println(f.Render(dialect.Postgres))
// Output: (price * quantity) AS total
```

### Example: Function
```go
f := field.New("SUM(price)", "sum_price")
fmt.Println(f.Render(dialect.Postgres))
// Output: SUM(price) AS sum_price
```

### Example: Literal
```go
f := field.New("'hello'", "greeting")
fmt.Println(f.Render(dialect.Postgres))
// Output: 'hello' AS greeting
```

### Example: Invalid input
```go
f := field.New("id as user_id foo")
fmt.Println(f.String())
// Output: ❌ field("id as user_id foo"): invalid format "id as user_id foo"
```

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project

