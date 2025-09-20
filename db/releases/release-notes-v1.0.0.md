# Release Notes — db v1.0.0 (2025-09-19)

This is the **first stable release** of the `entiqon/db` package.  
It provides a robust SQL builder and token framework with strict validation, test coverage, and extensibility.

---

## ✨ Added

### Builders
- **SelectBuilder**:
  - Structured condition handling with consistent `Build() (string, []any, error)` signature.
  - Supports:
    - Conditions: `Where`, `WhereAnd`, `WhereOr`.
    - Grouping: `GroupBy`, `ThenGroupBy`.
    - Having: `Having`, `AndHaving`, `OrHaving`.
    - Ordering: `OrderBy`, `ThenOrderBy`.
    - Pagination: `Take`, `Skip`.
- **Planned**: `InsertBuilder`, `UpdateBuilder`, `DeleteBuilder`, and `UpsertBuilder` (to be included in future versions).

### Tokens
- **Condition**: semantic-aware constructor with validation and dialect-aware rendering.
- **Field**: deterministic constructors (`New`, `NewWithTable`), aliased support, encapsulated token type.
- **Table**: constructors and helpers to define tables, aliases, and raw inputs with validation.
- **Join**: support for INNER, LEFT, RIGHT, CROSS, and NATURAL joins.

### Contracts
- **BaseToken**: core identity and validation for all tokens.
- **Errorable**: tokens can mark themselves errored after construction.
- **Kindable**, **Identifiable**, **Aliasable**: unify token identity handling.
- **Token**: aggregate contract for ownership and auditability.
  - Currently implemented by Field, Table, Join, and Condition tokens.
  - Future adoption planned for Having, GroupBy, and OrderBy tokens.
- **Rawable**, **Renderable**, **Stringable**, **Validable**: standardized behaviors for token rendering, debugging, and validation.

### Dialect Support
- Introduced **generic dialect** (implemented) using `?` placeholder strategy.
- Other dialects are planned for future releases, including Postgres, MySQL, MariaDB, SQLite, MSSQL, Oracle, DB2, Firebird, Informix, CockroachDB, TiDB, HANA, Snowflake, Redshift, Teradata, and ClickHouse.

---

## 🔄 Changed
- Renamed `Column` → `Field`.
- `SelectBuilder` constructor renamed from `NewSelect` → `New`.
- Refactored Field into dedicated subpackage `db/token/field`.
- Normalized resolver helpers (`ResolveExpression`, `ResolveCondition`, `IsValidSlice`).

---

## 📚 Documentation
- Comprehensive `doc.go` and `README.md` across builders, tokens, and contracts.
- Dialect guide and developer best practices.
- Runnable examples in `example_test.go` for every major component.

---

## 🧪 Tests & Coverage
- 100% unit test coverage across builders, tokens, dialects, and helpers.
- Extensive diagnostics for invalid fields, nil receivers, and missing sources.

---

## 📦 Availability
```bash
go get github.com/entiqon/db@v1.0.0
```
