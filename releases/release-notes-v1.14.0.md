# Release Notes

## v1.14.0 - Upcoming

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
