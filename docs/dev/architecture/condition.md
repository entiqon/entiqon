
# 🧩 token.Condition and token.ConditionType

This document describes the condition model used across all SQL builders in Entiqon. The `Condition` struct and related helpers enable safe and composable construction of SQL WHERE clauses.

---

## 📦 Overview

The `Condition` structure is used internally by:

- `SelectBuilder`
- `UpdateBuilder`
- `DeleteBuilder`
- (future) `Having`, `Join`, etc.

It is dialect-agnostic and resolved via `FormatConditions(...)`.

---

## 🔖 Location

```bash
internal/core/token/condition.go
internal/core/token/condition_ops.go
```

---

## 🧱 Struct

```go
type Condition struct {
	Type   ConditionType // e.g., SIMPLE, AND, OR
	Key    string        // SQL snippet like "id = ?"
	Params []any         // parameters associated with the condition
}
```

---

## 🔘 ConditionType

```go
type ConditionType string

const (
	ConditionSimple ConditionType = ""     // First/base WHERE
	ConditionAnd    ConditionType = "AND"  // AND chained
	ConditionOr     ConditionType = "OR"   // OR chained
)
```

---

## 🛠 Constructor

### `NewCondition(typ ConditionType, condition string, params ...any)`

Creates a new Condition object with the provided SQL snippet and parameters.

```go
token.NewCondition(token.ConditionAnd, "id = ?", 1)
```

---

## ➕ `AppendCondition`

```go
func AppendCondition(existing []Condition, c Condition) []Condition
```

Safely appends a new condition to an existing slice.

Used by:
- `AndWhere(...)`
- `OrWhere(...)`

---

## 🧪 `FormatConditions`

```go
func FormatConditions(dialect driver.Dialect, conditions []Condition) (string, []any, error)
```

Converts a slice of `Condition` into:
- a SQL string (e.g., `"id = ? AND status = ?"`)
- a slice of arguments (`[]any`)
- an error if any unknown ConditionType is used

---

## ✅ Behavior Notes

- `ConditionSimple` is treated as the initial `WHERE`
- All chained conditions require a type (`AND`, `OR`)
- The dialect is optional but enables future escaping

---

## 📂 Example

```go
conds := []Condition{
    NewCondition(ConditionSimple, "id = ?", 1),
    NewCondition(ConditionAnd, "status = ?", "active"),
}

sql, args, _ := FormatConditions(nil, conds)
// sql:  "id = ? AND status = ?"
// args: [1, "active"]
```

