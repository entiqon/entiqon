# 🧭 Dialect Engine Guide

This guide explains how to implement and extend SQL dialects in Entiqon.

---

## 🔁 Dialect Interface

Every dialect must implement the following methods:

```go
type Dialect interface {
    Name() string
    QuoteIdentifier(identifier string)
    QuoteLiteral(value any)
    BuildLimitOffset(limit, offset int)
    SupportsUpsert() bool
    SupportsReturning() bool
    Placeholder(index int) string  // Since: v1.4.0
}
```

---

## 🔢 Placeholder Support (Since: v1.4.0)

Each dialect must provide a formatting strategy for placeholders:

| Dialect     | Example Output |
|-------------|----------------|
| PostgreSQL  | `$1`, `$2`, ... |
| MySQL       | `?`, `?`, ...   |
| Generic     | `?`, `?`, ...   |

This enables query builders to emit correct SQL syntax per engine.

---

## 🔧 Dialect Usage by Builder

| Builder         | QuoteIdentifier | QuoteLiteral | Placeholder | Requires Dialect? |
|----------------|------------------|----------------|--------------|--------------------|
| SelectBuilder  | ✅               | ⚠️ Debug only  | ✅            | Optional           |
| InsertBuilder  | ✅               | ⚠️ Debug only  | ✅            | Optional           |
| UpdateBuilder  | ✅               | ⚠️ Debug only  | ✅            | Optional           |
| DeleteBuilder  | ✅               | ❌ Not used    | ✅            | Optional           |
| UpsertBuilder  | ✅               | ⚠️ Debug only  | ✅            | Optional           |

---

## 🆕 Quoting Policy (Since: v1.2.0)

| Method             | Purpose                      | Example        |
|--------------------|------------------------------|----------------|
| `QuoteIdentifier`  | Escapes table/column names   | `"user_id"`    |
| `QuoteLiteral`     | Escapes literal values       | `'abc'`, `42`  |

⚠️ `QuoteLiteral` is not SQL-safe and used only for logging/debugging.

---

## 🔄 Migrating a Custom Dialect

Update your dialects to support:

- ✅ `QuoteIdentifier(...)` (since v1.2.0)
- ✅ `QuoteLiteral(...)` (since v1.2.0)
- ✅ `Placeholder(index int)` (since v1.4.0)

---

## ✅ Example: PostgresDialect

```go
type PostgresDialect struct {
	BaseDialect
}

func (d *PostgresDialect) QuoteIdentifier(identifier string) string {
	return `"` + identifier + `"`
}

func (d *PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

func (d *PostgresDialect) SupportsUpsert() bool {
	return true
}

func (d *PostgresDialect) SupportsReturning() bool {
	return true
}
```

---

## 🧰 Helper: GeneratePlaceholders

```go
func GeneratePlaceholders(values [][]any, dialect driver.Dialect) ([]string, []any)
```

✅ Since: v1.4.0

Generates placeholder strings and flattens arguments for multi-row operations.

---

## 🔨 Adding a New Dialect

To extend Entiqon with a custom SQL dialect:

1. Create a new file in `internal/core/driver`, e.g. `dialect_sqlite.go`

2. Define a struct embedding `BaseDialect`:

```go
type SQLiteDialect struct {
	BaseDialect
}

func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{
		BaseDialect: BaseDialect{DialectName: "sqlite"},
	}
}
```

3. Override any required methods:

```go
func (d *SQLiteDialect) Placeholder(index int) string {
	return "?"
}

func (d *SQLiteDialect) SupportsUpsert() bool {
	return true
}

func (d *SQLiteDialect) QuoteIdentifier(identifier string) string {
	return "`" + identifier + "`" // MySQL/SQLite-style quoting
}
```

4. Use it directly in builders or expose it through `ResolveDialect(...)`.

---

## 🗑️ Deprecated

| Method         | Status        | Notes                                    |
|----------------|---------------|------------------------------------------|
| `Escape(...)`  | ❌ Removed     | Use `QuoteLiteral(...)` instead          |
| `WithDialect`  | ⚠️ Deprecated | Use `UseDialect(...)`. Removed in v1.4.0 |

---

## 🧭 Version History

| Feature                      | Version   |
|------------------------------|-----------|
| Dialect interface            | v1.3.0    |
| PostgresDialect              | v1.4.0    |
| GeneratePlaceholders helper  | v1.4.0    |