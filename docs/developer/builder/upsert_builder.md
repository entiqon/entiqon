# 🛠️ UpsertBuilder Developer Guide
**UPSERT (INSERT ON CONFLICT) builder with conflict resolution and dialect safety.**


The `UpsertBuilder` composes SQL `INSERT ... ON CONFLICT DO UPDATE/NOTHING` statements using a fluent API over an internal `InsertBuilder`. It ensures safe dialect-specific identifier handling and consistent structure across SQL dialects.

---

## ✅ Supported Features

* INSERT clause via embedded `InsertBuilder`
* ON CONFLICT clause with configurable conflict columns
* DO UPDATE SET clause using `Assignment{Column, Expr}`
* DO NOTHING fallback support
* RETURNING clause support
* Dialect-specific identifier quoting via `QuoteIdentifier`

---

## 🧱 Fluent API

```go
builder.NewUpsert().
    Into("users").
    Columns("id", "email").
    Values(1, "john@example.com").
    OnConflict("id").
    DoUpdateSet(
        builder.Assignment{Column: "email", Expr: "EXCLUDED.email"},
    ).
    Returning("id").
    UseDialect("postgres")
```

Produces:

```sql
INSERT INTO "users" ("id", "email")
VALUES (?, ?)
ON CONFLICT ("id")
DO UPDATE SET "email" = EXCLUDED.email
RETURNING "id"
```

---

## 📌 Dialect Quoting Strategy

We enforce **identifier-safe quoting** using `dialect.QuoteIdentifier(name string)`.

* ✅ Use for column and table names.
* ❌ Avoid using raw string formatting for identifiers.

For literal values (used only in **debug logs or test output**), use:

```go
dialect.QuoteLiteral(value any)
```

> ⚠️ This method is **NOT SQL-injection safe** and **MUST NOT** be used in query generation.

---

## 🔧 Internal Helpers

* `Assignment` struct is used to define DO UPDATE SET clauses.
* `UseDialect(...)` sets the dialect for identifier formatting.
* `WithDialect(...)` is **deprecated** — use `UseDialect(...)` instead.

---

## 💡 Naming Convention

As a project-wide guideline, we enforce:

* ✅ Use of **descriptive names** (e.g., `QuoteIdentifier`, `QuoteLiteral`)
* ❌ Avoid abbreviations like `QuoteIdent` or `EscapeVal`

> This improves readability, onboarding, and cross-dialect compatibility.

---

## 📂 Related Builders

* [InsertBuilder](./insert_builder.md)
* [SelectBuilder](./select_builder.md)
* [UpdateBuilder](./update_builder.md)

---

## ✅ Summary

The `UpsertBuilder` ensures:

* Dialect-safe identifier usage
* Flexible `ON CONFLICT` resolution
* Compatibility with Postgres' `RETURNING`
* Modern and consistent Go API design

Use it when your INSERT logic may result in duplicates and should either update or skip based on conflict rules.
---
## 🔧 Method Reference (Summary)
| Method        | Description                                    |
|---------------|------------------------------------------------|
| `UseDialect`  | Sets dialect for escaping identifiers          |
| `Into`        | Sets the target table                          |
| `Columns`     | Declares insert columns                        |
| `Values`      | Adds value rows for insertion                  |
| `OnConflict`  | Defines conflict detection columns             |
| `DoUpdateSet` | Assigns updates to apply on conflict           |
| `Returning`   | Specifies which columns to return              |
| `Build()`     | Compiles the full UPSERT SQL with placeholders |
---

## 🧪 Test Coverage

✅ **100% tested**, including deprecated behavior.

| Area                             | Coverage |
|----------------------------------|----------|
| `Into`, `Columns`, `Values`      | ✅        |
| `OnConflict`, `DoUpdateSet`      | ✅        |
| `RETURNING` clause               | ✅        |
| Dialect injection (`UseDialect`) | ✅        |
| Deprecated `WithDialect`         | ✅        |
| Validation rules                 | ✅        |
| Clause ordering                  | ✅        |
| Dialect-specific quoting         | ✅        |

> ✅ All builder methods and clause behaviors are tested.
> Even deprecated features like `WithDialect(...)` are covered for backward compatibility.