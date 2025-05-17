# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [v1.1.0] - 2025-05-17

### ✨ Added
- Introduced dialect-aware escaping via `WithDialect(...)` in all builders
- Implemented `PostgresEngine` with support for:
  - Escaping table and column identifiers
  - Escaping conflict and returning fields in UPSERT
- Exposed `Dialect Engine` interface for future extensibility

### 🔧 Refactored
- Unified condition handling via `token.Condition` with `Set`, `IsValid`, `AppendCondition`
- Applied shared `NewCondition(...)` constructor across all builders
- Updated `Select`, `Insert`, `Update`, `Delete`, and `Upsert` to support dialect injection
- Improved `UpsertBuilder` to delegate properly and inject dialect into `InsertBuilder`

### 📘 Documentation
- Updated README with:
  - Dialect usage example
  - New “Dialect Support” section
  - Go module version badge

---

Entiqon now provides a consistent, safe foundation for dialect-specific SQL generation — ready for PostgreSQL, and future engines.

---

## [v1.0.0] - 2025-05-16

### Added
- `SelectBuilder` upgraded to support argument binding and structured condition handling
- Consistent `Build() (string, []any, error)` signature across all builders
- Enhanced `ConditionToken` to store key, params, and raw string for traceability
- Suite-based test coverage added for all Select use cases

### Changed
- All builder examples in `README.md` now include clear descriptions
- `SelectBuilder` now aligns with INSERT, UPDATE, DELETE, and UPSERT builder patterns

### Stable
- Entiqon is now considered **stable** and tagged as `v1.0.0`
- CRUD is fully supported and extendable
- Structure, usage, and public API are locked and safe for production

---

## [v0.6.0] - 2025-05-16

### Added
- `DeleteBuilder` with fluent methods:
  - `.From()`, `.Where(...)`, `.Returning(...)`
  - `.Build()` outputs DELETE query with WHERE and RETURNING support
- PostgreSQL-style RETURNING clause supported
- Fully test-covered and GoDoc-compliant

### Changed
- README updated with runnable DELETE example
- Supported builders list now includes DELETE

---

## [v0.5.0] - 2025-05-16

### Added
- `UpsertBuilder` with full support for PostgreSQL-style UPSERT:
  - `.Into()`, `.Columns()`, `.Values()`, `.Returning()`
  - `.OnConflict(...)` and `.DoUpdateSet(...)`
- Delegates core insert logic to `InsertBuilder`
- Safe and expressive: `Build()` renders SQL + args for `ON CONFLICT DO UPDATE SET`

### Changed
- README updated with `UpsertBuilder` runnable example
- Supported builders list now includes UPSERT

### Tests
- Test verifies SQL + arg ordering
- Covers full clause behavior

---

## [v0.4.0] - 2025-05-16

### Added
- `UpdateBuilder` with fluent support for:
  - `.Table(...)`, `.Set(...)`, `.Where(...)`
  - `Build()` to generate safe SQL and ordered args

### Changed
- README updated with runnable `UpdateBuilder` example
- Usage examples now include full `package main` and imports for clarity

### Tests
- Unit test for `UpdateBuilder` covering SET and WHERE handling
- Argument ordering and SQL structure validated

---

## [v0.3.0] - 2025-05-16

### Added
- `InsertBuilder` with:
  - `.Into()`, `.Columns()`, `.Values()`
  - `.Returning()` for PostgreSQL-style response control
- `Build()` method for safe SQL generation

### Tests
- Insert test suite including multi-row and RETURNING clause
- Debug-friendly output with Watson audit signature

---

## [v0.2.0] - 2024-05-15

### Added
- Full `SelectBuilder` implementation with fluent chaining
- Support for logical condition grouping (`Where`, `AndWhere`, `OrWhere`)
- Pagination via `Take` (LIMIT) and `Skip` (OFFSET)
- `ConditionToken` and `ConditionType` struct for WHERE clause abstraction

### Changed
- Method names normalized for readability and control
- `Build()` method added to generate validated SQL string
- All methods and types fully documented using GoDoc
- Updated `README.md` with example usage and supported features

### Tests
- Complete test suite for `SelectBuilder` with `testify/suite`

---

## [v0.1.0] - 2024-04-02

### Added
- Initial `SelectQueryBuilder` prototype with basic `Select`, `From`, `Where`, `OrderBy`, `Limit`, `Offset`
- Query construction using simple string chaining
- Basic test coverage for SELECT behavior

### Changed
- Early param binding experiments for WHERE clauses
- Initial logo and branding
- Prepared project for pkg.go.dev publication

### Known Gaps
- No structured condition typing
- No dialect interface
- No documentation or strong API separation
