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

## 🛠️ Provided Dialects

| Dialect   | Placeholder | Quotes      | Supports RETURNING/UPSERT |
|-----------|-------------|-------------|----------------------------|
| `generic` | `?`         | none        | ❌                         |
| `postgres`| `$1`, `$2`  | `"column"`  | ✅                         |

Use:

```go
ResolveDialect("postgres") // returns PostgresDialect
ResolveDialect("unknown")  // returns BaseDialect named "generic"
```

---

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
