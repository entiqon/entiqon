# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## \[v1.3.0] - 2025-05-19

### ✨ Added

- `update_builder.md`: merged and normalized UpdateBuilder documentation
- Introduced `Dialect Guide` with version-tagged interfaces and test philosophy
- Added `Principles & Best Practices` section to README.md
- Linked all finalized builder docs under `docs/developer/builder/`
- Version tags added to all builder docs (`Since v1.2.0`)
- Explicit `Method Reference` and `Clause Ordering` sections added to guides

### 🧪 Coverage

- Achieved **100.0%** test coverage across all:
  - Builder methods
  - Token resolvers
  - Dialect interfaces (base, postgres, resolver)
- Deprecated methods (e.g., `WithDialect(...)`) remain tested until removal
- All helper methods, even unused, are now covered and documented

### 🧠 Philosophy

- Injected sarcastic validation quote into dialect guide:
  > “Even if necessary, tests will be tested.”
- All docs and builders adhere to strict validation, quoting, and fluent chaining principles

---

## \[v1.2.0] - 2025-05-18

### 📚 Documentation

* Moved all builder documentation into `/docs/builder/`
* Added centralized `/docs/index.md` with badges, overview, and links
* Integrated GitHub Pages deployment via Actions
* Updated README to offload examples and link each builder spec

### 🛠 Builders

* Finalized `UpsertBuilder` with clause-order enforcement
* Added `BuildInsertOnly()` to `InsertBuilder` for better delegation
* 100% test coverage including all validation branches and dialect fallback
* Strict enforcement of alias rules in `UpdateBuilder` and `UpsertBuilder`

### ⚙️ CI/CD

* Introduced `docs.yml` GitHub Action to auto-deploy docs on push to `main`
* Pages deploy pipeline ensures live site reflects every change

---

Entiqon is now fully documented and auto-published, with hardened query building and consistent structure across all SQL operations.

---

## \[v1.1.0] - 2025-05-17

### ✨ Added

* Introduced dialect-aware escaping via `WithDialect(...)` in all builders
* Implemented `PostgresEngine` with support for:

  * Escaping table and column identifiers
  * Escaping conflict and returning fields in UPSERT
* Exposed `Dialect Engine` interface for future extensibility

### 🔧 Refactored

* Unified condition handling via `token.Condition` with `Set`, `IsValid`, `AppendCondition`
* Applied shared `NewCondition(...)` constructor across all builders
* Updated `Select`, `Insert`, `Update`, `Delete`, and `Upsert` to support dialect injection
* Improved `UpsertBuilder` to delegate properly and inject dialect into `InsertBuilder`

### 📘 Documentation

* Updated README with:

  * Dialect usage example
  * New “Dialect Support” section
  * Go module version badge

---

Entiqon now provides a consistent, safe foundation for dialect-specific SQL generation — ready for PostgreSQL, and future engines.

---

## \[v1.0.0] - 2025-05-16

### Added

* `SelectBuilder` upgraded to support argument binding and structured condition handling
* Consistent `Build() (string, []any, error)` signature across all builders
* Enhanced \`ConditionToken
