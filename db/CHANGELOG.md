# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)  
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [v1.0.0] - 2025-09-19

### Added

- **Builders**
    - `SelectBuilder` with structured condition handling and consistent `Build() (string, []any, error)` signature.
    - Planned: `InsertBuilder`, `UpdateBuilder`, `DeleteBuilder`, and `UpsertBuilder` (to include full clause ordering
      and validation rules in future versions).
    - Support for conditions:
        - `Where`, `WhereAnd`, and `WhereOr` with validation and normalization.
    - Support for grouping and having clauses:
        - `GroupBy`, `ThenGroupBy`.
        - `Having`, `AndHaving`, `OrHaving`.
    - Support for ordering and pagination:
        - `OrderBy`, `ThenOrderBy`.
        - Pagination with `Take` and `Skip`.
- **Tokens**
    - `Condition`: semantic-aware constructor with validation and dialect-aware rendering.
    - `Field`: deterministic constructors (`New`, `NewWithTable`), aliased support, encapsulated token type.
    - `Table`: constructors and helpers to define tables, aliases, and raw inputs with validation.
    - `Join`: support for INNER, LEFT, RIGHT, CROSS, and NATURAL joins.
- **Contracts**
    - `BaseToken`: core identity and validation for all tokens.
    - `Errorable`: tokens can mark themselves errored after construction.
    - `Kindable`, `Identifiable`, `Aliasable`: unify token identity handling.
    - `Rawable`, `Renderable`, `Stringable`, `Validable`: standardized behaviors for token rendering, debugging, and
      validation.
    - `Token`: aggregate contract for ownership and auditability.
        - Currently implemented by **Field**, **Table**, **Join**, and **Condition** tokens.
        - Future adoption planned for **Having**, **GroupBy**, and **OrderBy** tokens.
- **Dialect Support**
    - Introduced **generic dialect** (implemented) with `?` placeholder strategy.
    - Other dialects are **planned** for future releases, including Postgres, MySQL, MariaDB, SQLite, MSSQL, Oracle,
      DB2, Firebird, Informix, CockroachDB, TiDB, HANA, Snowflake, Redshift, Teradata, and ClickHouse.

### Changed

- Renamed `Column` → `Field`.
- `SelectBuilder` constructor renamed from `NewSelect` → `New`.
- Refactored `Field` into dedicated subpackage `db/token/field`.
- Normalized resolver helpers (`ResolveExpression`, `ResolveCondition`, `IsValidSlice`).

### Documentation

- Comprehensive `doc.go` and `README.md` across builders, tokens, and contracts.
- Dialect guide and developer best practices.
- Runnable examples in `example_test.go` for every major component.

### Tests & Coverage

- 100% unit test coverage across builders, tokens, dialects, and helpers.
- Extensive diagnostics for invalid fields, nil receivers, and missing sources.
