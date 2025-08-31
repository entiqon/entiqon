# Release Notes

## v1.14.0 - Upcoming

### Token (join)
We have introduced **join.Type** as the canonical enum to represent SQL JOIN clauses.  
This addition provides a type-safe, dialect-agnostic way to classify and render all supported JOIN operations.

#### Features
- **Supported join types**:
    - `Inner` → `INNER JOIN`
    - `Left` → `LEFT JOIN`
    - `Right` → `RIGHT JOIN`
    - `Full` → `FULL JOIN`
    - `Cross` → `CROSS JOIN`
    - `Natural` → `NATURAL JOIN`
- **Methods**:
    - `String()` returns canonical SQL92 keywords.
    - `IsValid()` ensures strict recognition of supported joins.
    - `ParseFrom(string)` safely normalizes user input (`"INNER"`, `"LEFT JOIN"`, `"CROSS"`, `"NATURAL JOIN"`, etc.).
- **Validation**:
    - Invalid or unrecognized join types are reported as `INVALID`.

#### Documentation
- Added `doc.go` with an overview, examples, and usage guidelines.
- Added `README.md` describing supported joins and philosophy.
- Added `example_test.go` with runnable examples for all join types.

---

### Token (identifier)
We have introduced a new **identifier.Type** enum to classify SQL expressions into broad syntactic categories. This addition provides a foundation for consistent parsing, validation, and auditability of expression inputs.

#### Features
- **Categories** supported:
    - `Invalid`: could not classify
    - `Subquery`: `(SELECT …)`
    - `Computed`: parenthesized expressions like `(a+b)`
    - `Aggregate`: `SUM`, `COUNT`, `AVG`, `MIN`, `MAX`
    - `Function`: function calls such as `JSON_EXTRACT(data, '$.id')`
    - `Literal`: quoted strings or numeric constants
    - `Identifier`: plain table or column names (default fallback)
- **Methods**:
    - `Alias()` provides deterministic short codes (`id`, `lt`, `fn`, `ag`, `cp`, `sq`, `ex`).
    - `IsValid()` ensures strict recognition of supported kinds.
    - `ParseFrom(any)` safely coerces values from `int`, `string`, or an existing `Type`.
    - `String()` returns capitalized labels (`Identifier`, `Function`, …) with `Unknown` fallback.

#### Documentation
- Added `doc.go` with overview, categories, and philosophy.
- Updated `README.md` to reflect new examples, philosophy, and license reference.
- Added `example_test.go` demonstrating usage and edge cases.

---

### Database (field)
- Updated **field.Token** documentation (`doc.go`):
    - Added `BaseToken` and `Validable` contracts to the list of implemented interfaces.
    - Expanded construction rules for plain fields, inline/explicit aliases, wildcards (with alias restriction), subqueries (alias required), computed expressions, functions, and literals.
    - Clarified invalid cases (empty input, too many tokens, invalid alias, direct token usage without `Clone()`, unsupported types).
    - Added detailed examples for `Render`, `String`, `Debug`, wildcards, subqueries, functions, literals, and invalid inputs.
    - Reinforced design principles: immutability, auditability, strict validation, and safe cloning.

---

### Token (helpers)
#### Refactor
- **ResolveExpression** in `helpers/identifier.go` has been streamlined:
    - Branches directly on `ResolveExpressionType`, eliminating redundant checks.
    - Unified alias handling for all expression types (`Identifier`, `Subquery`, `Computed`, `Aggregate`, `Function`, `Literal`).
    - Removed unreachable `default` branch, ensuring full coverage.
    - Simplified responsibility split: classification validates kind/shape, resolution only extracts alias.

The expression classifier has been normalized and renamed.

#### Changes
- `ClassifyExpression` has been **renamed** to `ResolveExpressionType` and now lives in `helpers/identifier.go`.
- Provides syntactic classification of raw expressions into `identifier.Type`.
- `resolver.ResolveExpression` remains temporarily in `resolver.go` with a different return type until migration is complete.
- All docs, examples, and tests have been updated to reference `ResolveExpressionType`.

---
