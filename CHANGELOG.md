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
    - Initial helper: `IsValidIdentifier` / `ValidateIdentifier` with strict rules.
    - Non-ASCII identifiers (e.g. café, mañana, niño) explicitly rejected until dialect-specific rules are added.
    - Independent test file with exhaustive valid/invalid cases and runnable examples.
    - Includes `doc.go` and `README.md` documenting the consistency rule: `ValidateXxx` + `IsValidXxx`.

### Database (table/field)
- Constructors now delegate to `resolver.ValidateType` for type safety.
- Error states improved:
    - Passing tokens directly now suggests using `Clone()`.
    - Invalid literal/aggregate use as table sources rejected with clear error messages.
    - Invalid alias cases correctly rejected (including reserved keywords).

### Tests & Docs
- `doc.go` updated to include **resolver**, **ExpressionKind**, **join**, and **helpers**.
- `README.md` files updated:
    - Root `token` README now lists `field`, `table`, `join`, `resolver`, `ExpressionKind`, and `helpers`.
    - `table` README documents stricter alias validation, Clone() guidance, and error handling.
    - Headings normalized (emoji removed from `# Token`).
- `example_test.go` updated:
    - Subquery examples uncommented and corrected.
    - Added examples for invalid types and Clone() hints.
    - Adjusted `IsRaw` examples (currently false, will later derive from `Kind()`).

