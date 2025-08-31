# Release Notes v1.14.0

## Highlights

This release further refines the database token system, with a focus on **JOIN handling**, **expression resolution**, and **documentation improvements**. It also introduces several breaking changes around join types and struct naming for consistency and clarity.

---

## Database (join)

- Introduced **Join token (`join.Token`)** to represent SQL JOIN clauses:
  - Added safe constructors: `NewInner`, `NewLeft`, `NewRight`, `NewFull`.
  - Added flexible constructor: `New(kind any, left, right, condition)` for DSL scenarios.
  - Implemented new **join.Type** enum (`Inner`, `Left`, `Right`, `Full`, `Cross`, `Natural`).
  - Added early validation: invalid types → `invalid join type (n)`, errored tables, or missing condition produce clear error states.
  - Implements all core contracts: `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, `Stringable`, `Validable`.
- Added new join types:
  - `Cross` → renders as `CROSS JOIN`.
  - `Natural` → renders as `NATURAL JOIN`.
- **Breaking change**:
  - Removed legacy `join.Kind` in favor of `join.Type`.
  - Deleted `kind.go`.
  - Renamed struct from `join` → `token` for consistency with field/table tokens.
  - Updated `contract.go` and `token.go` (formerly `join.go`) accordingly.

---

## Database (field)

- Expanded **field.Token** construction rules:
  - Plain identifiers, inline/explicit aliases, wildcards (with alias restriction).
  - Subqueries (alias required), computed expressions, functions, literals.
- Added `BaseToken` and `Validable` to contracts.
- Clarified invalid cases (empty input, invalid alias, unsupported type, direct token without `Clone()`).
- Improved examples for `Render`, `String`, `Debug`, and error reporting.

---

## Token (resolver)

- Added new **resolver** module:
  - `ValidateType` rejects unsupported tokens and suggests `Clone()` for copies.
  - `ResolveExpr` extended with subquery detection, strict identifier validation, and explicit alias handling.

---

## Token (ExpressionKind)

- Added `Invalid` kind for unrecognized expressions.
- Updated classification rules:
  - Aggregates (`COUNT`, `SUM`, `AVG`, …) → `Aggregate`.
  - Computed (`price * qty`) → `Computed`.
  - Functions remain `Function`.

---

## Token (identifier)

- Introduced **identifier.Type** enum with categories: `Invalid`, `Subquery`, `Computed`, `Aggregate`, `Function`, `Literal`, `Identifier`.
- Provides short codes via `Alias()` (`id`, `fn`, `ag`, …).
- Added strict validation, parsing from `int|string|Type`, and safe `String()` output.

---

## Token (condition)

- Introduced **condition.Type** enum to classify SQL conditional expressions:
    - Supported values: `Invalid`, `Single`, `And`, `Or`.
    - Methods:
        - `IsValid()` validates recognized types.
        - `ParseFrom(any)` coerces from `Type`, `int`, or `string`.
        - `String()` returns canonical SQL keyword (`AND`, `OR`, or empty for `Single`).
    - Includes `normalize()` helper for case-insensitive parsing of strings.
- Added complete documentation:
    - `doc.go` with overview, categories, and usage philosophy.
    - `README.md` mirroring identifier/join structure with Purpose, Types, Example, Integration, License.
    - `example_test.go` demonstrating usage for `IsValid`, `String`, and `ParseFrom`.
    - `type_test.go` covering all constructors, branches, and edge cases with 100% coverage.

---

## Token (helpers)

- Refactored **ResolveExpression** to branch directly on `ResolveExpressionType`, unifying alias handling.
- Introduced **helpers** package:
  - Identifier and alias validation, reserved keywords.
  - Wildcard validation (`*` cannot be aliased).
  - Deterministic alias generation (`prefix + SHA-1`).
  - Expression classification via `ResolveExpressionType`.

---

## Database (table/field)

- Constructors now delegate to `resolver.ValidateType`.
- Clearer error messages for invalid literals, aggregates, or reserved aliases.
- Direct token usage now explicitly suggests `Clone()`.

---

## Tests & Documentation

- `doc.go` extended with resolver, ExpressionKind, join, and helpers.
- Updated README files for `token`, `helpers`, and `table` with stricter rules, validation guidance, and alias handling.
- Normalized headings.
- `example_test.go` updated:
  - Added examples for identifiers, aliases, wildcards, generated aliases, and expression classification.
  - Added invalid type examples and Clone() hints.
  - Adjusted IsRaw examples.

---

## Breaking Changes

- `join.Kind` → removed. Use `join.Type`.
- Struct `join` → renamed to `token`.
- `kind.go` → deleted.
- Contracts updated in `contract.go` and `token.go` (formerly `join.go`).

---

## Summary

This release consolidates the **join API** with a type-safe enum, removes outdated constructs, and improves expression resolution across the board. It also strengthens validation and expands test/documentation coverage, ensuring tokens remain immutable, auditable, and safe to use in builders.

