# 🌐 Developer Guide: SQL Dialects

This guide documents how Entiqon handles database-specific dialect behavior using a modular `Dialect` interface.

---

## 🧩 Interface

All dialects implement the following:

```go
type Dialect interface {
	Name() string
	Placeholder(position int) string
	QuoteTable(name string) string
}
```

---

## 🧱 Implementations

### ✅ BaseDialect (embedded)

Provides default behaviors:

- `Name()` → dialect name
- `Placeholder()` → returns `"?"`
- `QuoteIdentifier()` → returns identifier without quoting
- `QuoteLiteral(value)` → handles:
  - Strings → `'value'`
  - Numbers → `42`, `3.14`
  - Booleans → `true`, `false`
  - Other types → `fmt.Sprintf("'%v'", v)`
- `SupportsUpsert()` → `false`
- `SupportsReturning()` → `false`
- `BuildLimitOffset(limit, offset int)` → SQL LIMIT/OFFSET string

---

## 🚦 Limit/Offset Behavior

```go
BuildLimitOffset(limit, offset int)
```

Returns:

| Input                  | Output               |
|------------------------|----------------------|
| `limit=10, offset=20`  | `LIMIT 10 OFFSET 20` |
| `limit=5, offset=-1`   | `LIMIT 5`            |
| `limit=-1, offset=50`  | `OFFSET 50`          |
| `limit=-1, offset=-1`  | `""` (empty string)  |

✅ All conditions are fully tested.

---

## 🧠 Condition Resolution (Builder Integration)

- `Where("status = active")` is parsed inline and supported
- `Where("status", "active")` is structured key-value form
- Anything with more than 1 parameter is rejected for now
- Future support will allow:
  ```go
  Where("status = $1 AND age > $2", []any{true, 45})
  Where("status = :status AND age = :age", map[string]any{...})
  ```

### 🔒 Current Limitation
Expressions containing multiple logical parts (e.g. `AND`, `OR`) are not yet parsed or grouped.

### 🧠 TODO
- Grouped conditions will eventually be **wrapped in parentheses** to preserve precedence
- Example:
  ```sql
  WHERE (status = $1 AND age > $2) OR (country = $3)
  ```

---


## 🛠️ Provided Dialects

| DB Engine  | Dialect Name | Placeholder Style | Quote Style | Alias Style   | RETURNING Support | UPSERT Support | Since  |
|------------|--------------|-------------------|-------------|---------------|-------------------|----------------|--------|
| Generic    | `generic`    | `?`               | *(none)*    | ❌ Unsupported | ❌ None            | ❌ None         | v1.4.0 |
| PostgreSQL | `postgres`   | `$1, $2...`       | `"`         | `AS`          | ✅ Full            | ✅ Full         | v1.4.0 |
| MySQL      | `mysql`      | `?`               | `` ` ``     | `AS`          | ❌ None            | 🚫 Limited     | v1.4.0 |
| SQLite     | `sqlite`     | `?`, `:name`      | `"`         | Optional `AS` | ✅ v3.35+          | ✅ v3.24+       | v1.5.0 |
| SQL Server | `mssql`      | `@param`          | `[` `]`     | Optional `AS` | ✅ via OUTPUT      | 🚫 via MERGE   | v1.4.0 |
| Oracle     | `oracle`     | `:param`          | `"`         | Optional `AS` | ✅ Full            | ✅ via MERGE    | v1.6.0 |
| IBM DB2    | `db2`        | `:param`          | `"`         | Optional `AS` | ✅ Partial         | ✅ via MERGE    | v1.6.0 |
| Firebird   | `firebird`   | `?`, `:param`     | `"`         | Optional `AS` | ✅ Supported       | ✅ Limited      | v1.6.0 |
| Informix   | `informix`   | `?`               | `"`         | Optional `AS` | ✅ Supported       | ❌ Not native   | v1.6.0 |

Use:

```go
ResolveDialect("postgres") // returns PostgresDialect
ResolveDialect("unknown")  // returns BaseDialect named "generic"
```

## 🔗 Integration with ParamBinder

The `ParamBinder` uses `dialect.Placeholder(n)` to assign placeholders during query construction.

```go
pb := NewParamBinder(dialect)
pb.Bind("id")  // → $1 or ?
```

---

## ✅ Test Strategy

All dialects are tested via `TestDialectSuite` using shared assertions for:

- Placeholder generation
- Literal and identifier quoting
- Dialect fallback resolution
- Limit/Offset formatting

---

## 🚧 Notes

- Dialects are **non-configurable at runtime**
- If needed, create a `NewMySQLDialect()` or similar
- Do **not** hardcode placeholders inside builders — use dialect

---

## 📁 Files

- `engine.go` — interface
- `dialect_base.go` — base logic
- `dialect_generic.go` — `NewGenericDialect`
- `dialect_postgres.go` — PostgreSQL behavior
- `dialect_resolver.go` — dialect lookup map
- `dialect_test.go` — shared test coverage
