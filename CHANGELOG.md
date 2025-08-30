# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## v1.14.0 - Upcoming

### Database (join)
- Introduced **Join token (`join.Token`)** to represent SQL JOIN clauses:
    - Safe constructors: `NewInner`, `NewLeft`, `NewRight`, `NewFull`.
    - Flexible constructor: `New(kind any, left, right, condition)` for advanced/DSL scenarios.
    - Explicit `join.Kind` enum (`InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin`) with `String()`, `IsValid()`, and `ParseJoinKindFrom()`.
    - Early-exit validation: invalid kind → `invalid join type (n)`, nil/errored tables, or empty condition → explicit error states.
    - Implements all core contracts: `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, `Stringable`, `Validable`.

### Token (resolver)
- Added **resolver** module:
    - `ValidateType` enforces input types:
        - `string` accepted.
        - Existing tokens (`Validable`) rejected with `unsupported type …; if you want to create a copy, use Clone() instead`.
        - All other types → `invalid format (type …)`.
    - `ResolveExpr` extended with:
        - Subquery detection (input wrapped in parentheses treated as one expression).
        - Strict identifier validation (must be a single token).
        - Explicit alias handling (`AS`, trailing identifiers).

### Token (ExpressionKind)
- Added `Invalid` kind for unrecognized expressions.
- Updated `String()` and `IsValid()` accordingly.
- Extended classification rules:
    - Aggregates (`COUNT`, `SUM`, `AVG`, `MIN`, `MAX`) now reported as `Aggregate`.
    - Computed expressions (`price * quantity`) reported as `Computed`.
    - Functions (`JSON_EXTRACT(...)`) remain `Function`.

### Token (helpers)
- Introduced **helpers** package for reusable validation utilities.
    - Identifier validation:
        - `IsValidIdentifier` / `ValidateIdentifier` with strict rules.
        - Non-ASCII identifiers (e.g. café, mañana, niño) explicitly rejected until dialect-specific rules are added.
    - Alias validation:
        - `IsValidAlias` / `ValidateAlias` to ensure aliases are valid identifiers and not reserved keywords.
        - `ValidateTrailingAlias` / `HasTrailingAlias` to detect and validate trailing aliases (when no `AS` is present).
        - `ReservedKeywords()` returns the dialect-agnostic set of disallowed aliases.
    - Wildcard validation:
        - `ValidateWildcard(expr, alias)` ensures that `*` is only used without alias or raw.
        - Rejects invalid cases such as `* AS total`.
    - Deterministic alias generation:
        - `GenerateAlias(prefix, expr string)` produces safe aliases by combining a short code with a SHA-1 hash of the expression.
    - Independent test files with exhaustive valid/invalid cases and runnable examples.
    - Includes `doc.go` and `README.md` documenting rules and the consistency pattern (`ValidateXxx`, `IsValidXxx`, `GenerateAlias`).

### Database (table/field)
- Constructors now delegate to `resolver.ValidateType` for type safety.
- Error states improved:
    - Passing tokens directly now suggests using `Clone()`.
    - Invalid literal/aggregate use as table sources rejected with clear error messages.
    - Invalid alias cases correctly rejected (including reserved keywords).

### Tests & Docs
- `doc.go` updated to include **resolver**, **ExpressionKind**, **join**, and **helpers** (identifiers, aliases, trailing alias detection, wildcard validation, alias generation).
- `README.md` files updated:
    - Root `token` README now lists `field`, `table`, `join`, `resolver`, `ExpressionKind`, and `helpers`.
    - `helpers` README extended with identifiers, aliases, trailing alias rules, reserved keywords, wildcard validation, and `GenerateAlias`.
    - `table` README documents stricter alias validation, Clone() guidance, and error handling.
    - Headings normalized (emoji removed from `# Token`).
- `example_test.go` updated:
    - Subquery examples uncommented and corrected.
    - Added examples for identifiers, aliases, trailing aliases, wildcards, and generated aliases.
    - Added examples for invalid types and Clone() hints.
    - Adjusted `IsRaw` examples (currently false, will later derive from `Kind()`).

## [v1.13.0](https://github.com/entiqon/entiqon/releases/tag/v1.13.0) - 2025-08-26

### Database (builder/select)
- Added full clause support:
    - **Conditions**: `Where`, `And`, `Or` (reset, append, normalize with AND, ignore empty).
    - **Grouping**: `GroupBy`, `ThenGroupBy` (reset/append, ignore empty, rendered between WHERE and HAVING).
    - **Having**: `Having`, `AndHaving`, `OrHaving` (reset/append, normalize with AND, rendered after GROUP BY).
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
- Extended **Errorable** contract with `SetError(err error)`:
    - Allows tokens/builders to mark themselves as errored after construction
    - Implemented in `Field` and `Table` tokens
    - Updated `doc.go`, `README.md`, and `example_test.go` to demonstrate usage

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
- Struct `Field` is now hidden (`field` unexported); only the `Token` interface is public.
- Constructors (`New`, `NewWithTable`) return `Token` instead of *field, ensuring encapsulation.
- `Clone()` updated to return `Token`.
- Example tests aligned:
    - Functions renamed from ExampleField_* → ExampleToken_*.
    - Debug/String outputs updated to lowercase `field(...)`.
    - Clone example simplified (no pointer comparison).

- Refactored **Field** into dedicated subpackage `db/token/field`:
    - Preserved API (`field.New(...)`) and contract implementations unchanged.
    - Updated `builder/select.go` and `select_test.go` to import from the new path.
    - Adjusted **Dockerfile-documentation** to copy `db/token/field/README.md`.
    - Normalized structure for consistency with `token/table`.
- Introduced **Token contract** as a scaffold to decompose Field identity into auditable pieces:
    - Aggregates `BaseToken`, `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, and `Stringable`.
    - Defines ownership methods: `HasOwner()`, `Owner()`, and `SetOwner()`.
    - Intention: sort every identity aspect (expr, alias, owner, validity, raw state) into dedicated contracts for auditability.
    - ⚠️ Currently only the contract and documentation are provided; implementation is staged for later commits.
- Documentation updates:
    - `doc.go`: added `HasOwner`, `Owner`, `SetOwner` under **Field Behavior**.
    - `example_test.go`: added placeholder `ExampleField_owner` (commented out until implemented).
    - `README.md`: new **Contracts and Auditability** section in Developer Guide.
- Root-level documentation:
    - Added `db/token/README.md` with purpose and subpackage table (field, table).
    - Added `db/token/doc.go` with GoDoc overview, principles, subpackages, and roadmap.

### Common (extension/integer)
- Introduced **integer parser**:
    - `ParseFrom(any)` converts input safely to int, rejects invalid types.
    - `IntegerOr` shortcut with defaults.
    - Consistent with `boolean`, `float`, and `decimal` parsers.

### Tests & Docs
- Comprehensive unit tests across Database and Common ensuring 100% coverage.
- **README.md**, **doc.go**, and **example_test.go** updated with Conditions, Grouping, Having, Ordering, and Integer parser sections.
- Added runnable examples in `example_test.go` demonstrating new Having clause.

... (rest of history unchanged)
