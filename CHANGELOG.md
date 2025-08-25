# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v1.13.0] - Upcoming

### Database (builder/select)
- Added full clause support:
    - **Conditions**: `Where`, `And`, `Or` (reset, append, normalize with AND, ignore empty).
    - **Grouping**: `GroupBy`, `ThenGroupBy` (reset/append, ignore empty, rendered between WHERE and HAVING).
    - **Having**: `Having`, `AndHaving`, `OrHaving` (reset/append, normalize with AND, ignore empty, rendered after GROUP BY).
    - **Ordering**: `OrderBy`, `ThenOrderBy` (reset/append, ignore empty, rendered after GROUP BY/HAVING).
- Enhanced diagnostics & reporting:
    - Invalid fields produce consistent `⛔️ Field("<expr>"): input type unsupported: <type>` errors.
    - `Debug()` and `String()` improved with ✅/⛔️ status markers.
    - `Build()` aggregates invalid fields, detects nil receiver and missing source, with clear ❌ messages.

### Database (contract)
- Introduced **BaseToken** contract (`db/contract/base_token.go`):
    - Provides core identity and validation for all tokens
    - Methods: `Input()`, `Expr()`, `Alias()`, `IsAliased()`, `IsValid()`
    - Ensures `Field`, `Table`, and future tokens expose consistent state
- Added runnable example (`ExampleBaseToken`) in `example_test.go`
- Updated `doc.go` to include BaseToken in contract overview with normalized style
- Updated `README.md`:
    - Added BaseToken section with purpose, methods, and usage
    - Streamlined layout, removed redundancy
    - Extended philosophy with **Consistency** principle: all tokens share BaseToken
  
### Database (token/table)
- Introduced **Table token** to represent SQL sources in builders:
    - Provides constructors and helpers to define tables, aliases, and raw inputs.
    - Supports consistent rendering across dialects.
    - Ensures validation of invalid/empty inputs with clear error reporting.
- Added **unit tests** covering constructors, methods, and edge cases with 100% coverage.
- Added **doc.go** with package overview and usage guidelines.
- Added **example_test.go** with runnable examples.
- Added **README.md** documenting purpose, design, and usage of token.Table.

### Database (token/field)
- Refactored **Field** into dedicated subpackage `db/token/field`:
    - Preserved API (`field.New(...)`) and contract implementations unchanged.
    - Updated `builder/select.go` and `select_test.go` to import from the new path.
    - Adjusted **Dockerfile-documentation** to copy `db/token/field/README.md`.
    - Normalized structure for consistency with `token/table`.

### Common (extension/integer)
- Introduced **integer parser**:
    - `ParseFrom(any)` converts input safely to int, rejects invalid types.
    - `IntegerOr` shortcut with defaults.
    - Consistent with `boolean`, `float`, and `decimal` parsers.

### Tests & Docs
- Comprehensive unit tests across Database and Common ensuring 100% coverage.
- **README.md**, **doc.go**, and **example_test.go** updated with Conditions, Grouping, Having, Ordering, and Integer parser sections.
- Added runnable examples in `example_test.go` demonstrating new Having clause.

## [v1.12.0](https://github.com/entiqon/entiqon/releases/tag/v1.12.0) - 2025-08-22

### Highlights

- **Parser Shortcuts**: Added `Or` variants (`BooleanOr`, `NumberOr`, `FloatOr`, `DecimalOr`, `DateOr`) to allow explicit fallback values.
- **Extension Documentation**: Each extension subpackage (`boolean`, `date`, `decimal`, `float`, `number`, `object`, `collection`) now ships with `README.md`, `doc.go`, and `example_test.go`.
- **Object Helpers**: Normalized and moved from `common/object` → `common/extension/object`.
- **Deprecations**:
    - `BoolToStr` has been renamed to `BoolToString` and moved into `common/extension/boolean`.
    - Existing usages of `BoolToStr` will continue to work but are discouraged.

### Documentation

- Added root-level `doc.go` for `common` describing purpose and structure.
- Added `README.md` for `common`, linking subpackages `errors` and `extension`.
- Added package-level READMEs and examples across all `common/extension` subpackages.
- Dockerfile updated to copy package-level READMEs into site docs.
- Navigation updated with alphabetized extension packages.

### Internal

- Unified test coverage for all parsing helpers and `Or` fallbacks.
- Refactored Dockerfile to preserve base image but copy documentation automatically.
- Clarified supported input formats in `date.ParseAndFormat` examples.

## [v1.11.0](https://github.com/entiqon/entiqon/releases/tag/v1.11.0) - 2025-08-17

### Highlights

-   **New parsing façade**: one-line helpers for `Boolean`, `Float`,
    `Decimal`, and `Date`.
-   **Deterministic date cleaning**: `CleanAndParse`,
    `CleanAndParseAsString`, strict `YYYYMMDD` prefix path, and 100%
    tests.
-   **Boolean parser++**: extended tokens (`on/off`, `y/n`, `t/f`) and
    explicit `nil` rejection.
-   **SQL Builder/Token overhaul**: `Column → Field`, `FieldCollection`,
    deterministic `NewField` inputs, and Postgres/Base dialects with
    tests.

## [v1.10.0](https://github.com/entiqon/entiqon/releases/tag/v1.10.1) - 2025-08-07

- dce6cf7 feat(object): enhance Exists, GetValue, SetValue with flexible types; add extensive tests; update docs (Isidro Lopez)

## [v1.9.0](https://github.com/entiqon/entiqon/releases/tag/v1.9.0) - 2025-08-07

### Added

- 09e6632 feat(common/math): add float and decimal parsing packages and update documentation (Isidro Lopez)

## [v1.8.4](https://github.com/entiqon/entiqon/releases/tag/v1.8.4) - 2025-08-07

### Added

- bb3af79 docs(common/number): add ParseFrom utility details with rounding flag and float string parsing (Isidro
  Lopez)
- 7c6982e feat(common/number): enhance ParseFrom with float string parsing and rounding flag (Isidro Lopez)
- 4b204c3 feat(number): extend ParseFrom to support bool values (Isidro Lopez)

## [v1.8.3](https://github.com/entiqon/entiqon/releases/tag/v1.8.3) - 2025-08-06

### Added

- 55d24c3 feat(common): add generic object utilities and enhanced error types (Isidro Lopez)

## [v1.8.2](https://github.com/entiqon/entiqon/releases/tag/v1.8.2)

### Added

- 4cd3d28 feat(project): add ProcessStageError, ProcessStage guide, deployment automation, and docs table update (
  Isidro Lopez)

## [v1.8.1](https://github.com/entiqon/entiqon/releases/tag/v1.8.1)

### Added

- 0b32dbd feat(entiqon): return to monorepo technique (Isidro Lopez)

## [v1.8.0] - 2025-08-02

### Added

- Modularization of `db` package as a standalone Go module: `github.com/entiqon/db`.
- Updated import paths from `github.com/entiqon/entiqon/...` to `github.com/entiqon/db/...`.
- Added initial test coverage and CI integration for `db` module.
- Future modules like `core` will be added modularly following this pattern.
- Icon and documentation support for:
    - Category 1: **Entiqon Sharicon** (orange shared folder)
    - Category 2: **Entiqon Corecon** (green gear)
    - Category 3: **Entiqon Datacon** (green DB cylinder)
    - Category 4: **Entiqon Commicon** (teal antenna)
    - Category 5: **Entiqon Toolicon** (red wrench & hammer)

**Note:** This is a breaking change; downstream users must update import paths accordingly.

**Codename:** Atlas

## [v1.7.0] - 2025-06-09

- refactor(token): rename HasError → IsErrored(), SetErrorWith → SetError()
- feat(token): add GetError(), SetError(), IsErrored() (Errorable contract)
- feat(token): introduce Kind enum and Kindable interface
- feat(token): add SetKind(), GetKind() to BaseToken
- chore(token): update BaseToken.String() to use Kindable & Errorable
- docs: add base_token.md and update token.md to document new contracts
- test: add tests for Errorable and Kindable methods (nil-safe coverage)

**Codename:** Forge

---

## [v1.6.0] - 2025-05-25

### Added

- Support for aliasable SQL expressions with optional table qualification
- Introduces reusable token abstraction for handling elements like `table.column AS alias`
- Provides helper methods: `Raw()`, `String()`, `WithTable()`, `IsValid()`, etc.

### Notes

This version enables standardized handling of SQL identifiers across the builder, including aliases and table scoping.
It's the first practical application of the previously introduced `AliasableToken`.

**Codename:** Keystone

---

# Changelog

## [v1.5.0] – 2025-05-24

### ✨ Features

- `StageToken`: Standardized clause tagging for error validation
- Public Dialect API: Enables custom dialect implementations
- Enhanced builder error tagging with `StageToken`

### 🛠️ Refactors

- Normalized ParamBinder and Condition to use injected Dialect
- All builders updated to unified token/dialect handling strategy
- Centralized placeholder resolution across WHERE, SET, VALUES

### ✅ Test Coverage

- builder: 94.4%
- driver: 100.0%
- core/builder: 95.5%
- bind, test packages: 100%
- error/token: >75%

### 📄 Docs

- Updated all builder guides
- Documented new dialect interface and StageToken strategy

---

## [v1.4.0] - 2025-05-22

### ✨ Added

- `NewCondition`: semantic-aware condition constructor with support for:
    - Inline and placeholder syntax (e.g., `"status = active"`, `"status = ?"` + value)
    - Type inference (`int`, `bool`, `float64`, `string`)
    - Operator support: `=`, `!=`, `<>`, `<`, `>`, `<=`, `>=`, `IN`, `NOT IN`, `BETWEEN`, `LIKE`, `NOT LIKE`
- `ParamBinder`: unified argument binding for dialects (`?`, `$N`, `:field`)
- `condition_helpers.go`: includes `InferLiteralType`, `ParsePlaceholderPattern`, `AllSameType`, and
  `ContainsUnboundPlaceholder`

### 🧱 Builders Implemented

- `SelectBuilder`
- `InsertBuilder`
- `UpdateBuilder`
- `DeleteBuilder`
- `UpsertBuilder` (with `ON CONFLICT`, `DO UPDATE`, and `DO NOTHING`)

### ✅ Validation

- All builders enforce:
    - Table presence
    - Column count/value alignment
    - Alias rejection in INSERT/UPSERT
    - Invalid or unsupported condition rejection via `AddStageError(...)`
- `UpsertBuilder` supports `RETURNING` only if the dialect allows it (`SupportsReturning()`)

### 🧪 Test Coverage

- Achieved **100% coverage** on:
    - `select.go`, `insert.go`, `update.go`, `delete.go`, `upsert.go`
    - `condition_renderer.go`, `condition_helpers.go`, `param_binder.go`

---

## [v1.3.0] - 2025-05-19

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

## [v1.2.0] - 2025-05-18

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

Entiqon is now fully documented and auto-published, with hardened query building and consistent structure across all SQL
operations.

---

## [v1.1.0] - 2025-05-17

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

Entiqon now provides a consistent, safe foundation for dialect-specific SQL generation — ready for PostgreSQL, and
future engines.

---

## [v1.0.0] - 2025-05-16

### Added

* `SelectBuilder` upgraded to support argument binding and structured condition handling
* Consistent `Build() (string, []any, error)` signature across all builders
* Enhanced `ConditionToken
