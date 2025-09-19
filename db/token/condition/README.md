# Condition Token

> Part of [Entiqon](../../../) / [Database](../../) / [Token](../)

## 🌱 Overview

The `condition.Token` type represents a SQL condition (predicate) such as  
`id = 1`, `age > 18`, or `status IN ('active', 'pending')`.  
It provides a structured, validated, and renderable abstraction for WHERE and HAVING clauses.

Built on top of Entiqon’s shared **contracts** (e.g., `Renderable`, `Validable`, `Debuggable`),  
`Token` ensures safe construction, immutability, and diagnostic visibility of SQL conditions.

---

## Construction Rules

Conditions are created using `condition.New(...)` or convenience helpers:

1. **Basic expression**
   ```go
   c := condition.New(ct.Single, "id = 1")
   // → id = :id, value=1
   ```

2. **Expression with parameter**
   ```go
   c := condition.New(ct.Single, "id = ?", 123)
   // → id = :id, value=123
   ```

3. **Named parameter**
   ```go
   c := condition.New(ct.Single, "id = :id", 123)
   // → id = :id, value=123
   ```

4. **Operator with value**
   ```go
   c := condition.New(ct.Single, "age", operator.GreaterThan, 18)
   // → age > :age, value=18
   ```

5. **Logical AND / OR**
   ```go
   c := condition.NewAnd("status = ?", "active")
   // Kind=And

   c = condition.NewOr("deleted = ?", false)
   // Kind=Or
   ```

6. **Inline collections**
   ```go
   c := condition.New(ct.Single, "id IN (1, 2, 3)")
   // → id IN (:id), value=[1 2 3]
   ```

7. **Special operators**
   ```go
   c := condition.New(ct.Single, "id IS NULL")
   // → id IS NULL, value=nil
   ```

8. **Invalid cases**
   - Empty kind → errored (`invalid condition type`)
   - Empty expression → errored
   - Wrong type (non-string expr) → errored
   - Invalid operator/value list (e.g., `IN` with empty slice) → errored

---

## Contracts Implemented

- **Kindable** → classification (`Kind()`, `SetKind()`)
- **Identifiable** → core identity (`Input()`, `Expr()`, `Name()`)
- **Errorable** → `Error()`, `IsErrored()`, `SetError()`
- **Debuggable** → `Debug()` (full diagnostic output)
- **Rawable** → `Raw()`, `IsRaw()`
- **Renderable** → `Render()` (SQL form)
- **Stringable** → `String()` (concise log form)
- **Validable** → `IsValid()` (inverse of `IsErrored()`)

---

## Examples

### Example: Simple condition
```go
c := condition.New(ct.Single, "id = ?", 42)
fmt.Println(c.Render())
// Output: id = :id
fmt.Println(c.Value())
// Output: 42
```

### Example: Named parameter
```go
c := condition.New(ct.Single, "id = :id", 42)
fmt.Println(c.String())
// Output: Condition("id = :id"): name="id", value=42, errored=false
```

### Example: With operator
```go
c := condition.New(ct.Single, "age", operator.GreaterThan, 18)
fmt.Println(c.Render())
// Output: age > :age
```

### Example: Logical AND
```go
c := condition.NewAnd("status = ?", "active")
fmt.Println(c.Kind())
// Output: And
```

### Example: Debug output
```go
c := condition.New(ct.Single, "id", operator.In, []int{1, 2, 3})
fmt.Println(c.Debug())
// Output: Condition{Input="id IN [1 2 3]", Type:"Single", Expression="id IN (:id)", Value=[1 2 3], Error=<nil>}
```

### Example: Invalid expression
```go
c := condition.New(ct.Single, "")
fmt.Println(c.Error())
// Output: empty expression
```

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project
