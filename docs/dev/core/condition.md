# 🧩 Developer Guide: Condition System

This document defines how to use and construct `Condition` objects for SQL query builders in Entiqon.

---

## 🎯 Purpose

The `Condition` struct represents a normalized WHERE/AND/OR clause used across all query builders. It encapsulates:

- The **field expr**
- The **SQL operator**
- A **slice of values**
- A **dialect-safe raw expression**
- A `ConditionType` (e.g., SIMPLE, AND, OR)
- A `.IsValid()` method for validation

---

## 🔧 Construction Rules

### ✅ Basic Equality

```go
c := NewCondition(ConditionSimple, "status", "active")
// → status = :status
```

### ✅ Inline Parsing (one param)

```go
c := NewCondition(ConditionSimple, "status != active")
// → status != :status
```

Supported inline operators (resolved in order):  
`NOT IN`, `IN`, `BETWEEN`, `<>`, `!=`, `>=`, `<=`, `LIKE`, `=`, `>`, `<`

---

## 🔨 Explicit Operators

Use `NewConditionWithOperator` for full control:

```go
c := NewConditionWithOperator(ConditionAnd, "age", ">=", 21)
```

Or use specialized constructors:

```go
NewConditionIn(ConditionSimple, "id", 1, 2, 3)
NewConditionBetween(ConditionAnd, "created_at", start, end)
NewConditionLike(ConditionOr, "expr", "Joh%")
```

---

## 🔐 Type Compatibility (Multi-value Conditions)

Functions like `NewConditionIn(...)` and `NewConditionBetween(...)` require all values to be:

- At least 2 items
- Type-compatible (e.g., all strings, all numbers, all time.Time)

Enforced via `AreCompatibleTypes(...)`.

---

## ❌ Deprecated

- `Set()` → removed
- `setOperator()` → removed
- `FormatConditions()` → will be removed in favor of `ParamBinder`

---

## 🧪 Validation

Use `.IsValid()` on any `Condition` to verify it was constructed correctly.

```go
if !c.IsValid() {
    log.Println("invalid condition:", c.Error)
}
```

---

## 📦 Required in Builders

All builder methods (`Where`, `AndWhere`, etc.) must construct conditions using these factories. Builders must **never** split or infer operators from strings.

