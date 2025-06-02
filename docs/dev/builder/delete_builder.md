# 🧱 DeleteBuilder Developer Guide
**DELETE builder with condition chaining, dialect integration, and error validation.**


The `DeleteBuilder` is a fluent query generator for composing SQL `DELETE` statements with argument binding and dialect-aware formatting. It is part of Entiqon's core builder module and provides composable support for:

- Table selection
- WHERE conditions (`WHERE`, `AND`, `OR`)
- RETURNING clause (Postgres)
- Dialect injection for identifier escaping

---

## 🚀 Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/entiqon/entiqon/internal/core/builder"
)

func main() {
	sql, args, err := builder.NewDelete().
		WithDialect(&dialect.PostgresEngine{}).
		From("users").
		Where("status = ?", "inactive").
		AndWhere("created_at < ?", "2024-01-01").
		Returning("id").
		Build()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)
}
```

---

## 🛠️ Constructor

```go
builder.NewDelete()
```

---

## 🧩 Fluent API

### `From(table string)`
Sets the table from which rows will be deleted.

---

### `Where(condition string, params ...any)`
Initializes the `WHERE` clause with a condition.
This will override any previous condition chain.

---

### `AndWhere(condition string, params ...any)`
Appends an additional `AND` condition to the current chain.

---

### `OrWhere(condition string, params ...any)`
Appends an additional `OR` condition to the current chain.

---

### `Returning(fields ...string)`
Appends a `RETURNING` clause to retrieve deleted rows (PostgreSQL only).

---

### `UseDialect(name string)`
Resolves and injects a dialect engine by name (e.g., `"postgres"`).
This enables dialect-aware table escaping.

---

### `WithDialect(d Dialect)` ⚠️ *Deprecated*
Allows injecting a dialect engine directly.
Use `UseDialect(name)` instead for consistency.

---

### `Build() (string, []any, error)`
Builds the final SQL statement and returns:
- the raw SQL string with `?` placeholders
- a slice of argument values
- an error if any part of the query is invalid

---

## 🔐 Clause Ordering

1. `DELETE FROM ...`
2. `WHERE ...` (optional)
3. `RETURNING ...` (optional)

---

## 🔄 Dialect Fallback Behavior

- If no dialect is set:
  - Identifiers are used as-is
  - Placeholders remain `?`
  - All validation rules still apply

---

## ⚠️ Validation Rules

- If `.From()` is missing, `Build()` returns an error.
- If any `ConditionType` is invalid, `Build()` will fail safely.

---

## ✅ Test Coverage

| Function        | Coverage | Status |
|-----------------|----------|--------|
| `NewDelete`     | 100%     | ✅     |
| `From`          | 100%     | ✅     |
| `Where`         | 100%     | ✅     |
| `AndWhere`      | 100%     | ✅     |
| `OrWhere`       | 100%     | ✅     |
| `Returning`     | 100%     | ✅     |
| `UseDialect`    | 100%     | ✅     |
| `WithDialect`   | 100%     | ✅ (Deprecated) |
| `Build`         | 100%     | ✅     |