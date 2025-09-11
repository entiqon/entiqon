# Release Notes v1.14.0

## Highlights

This release refines the database token system with a focus on **JOIN handling**, **expression resolution**, and **operator classification**. It introduces a type registry for operators, restructures type enums, and expands documentation and tests.

---

## Contracts

- Introduced **Kindable**, **Identifiable**, **Aliasable** contracts.
- **BaseToken** composes Identifiable + Aliasable to avoid duplication.
- Documentation updated (`doc.go`, `README.md`, examples).

---

## Types

- **join.Type**: new enum (`Inner`, `Left`, `Right`, `Full`, `Cross`, `Natural`); removed legacy `join.Kind`.
- **condition.Type**: new enum (`Invalid`, `Single`, `And`, `Or`).
- **identifier.Type**: new enum (`Identifier`, `Subquery`, `Literal`, `Aggregate`, `Function`, `Computed`).
- **ExpressionKind**: added `Invalid` classification; extended for aggregates, computed, functions.
- **operator.Type**: refactored to a typed **registry** `{String, Alias, Position, Synonyms}`:
    - Deterministic `GetKnownOperators()` ordering by Position.
    - O(1) `ParseFrom()` via reverse index; accepts symbols (`!=`, `>=`), words (`NOT IN`, `IS NULL`), aliases (`nin`, `gte`, `isnull`).
    - Simplified `String()`, `Alias()`, `IsValid()` with registry lookups.
    - Added full GoDoc, `doc.go`, README, and 100% tests.

---

## Tokens

- **Condition token (`condition.Token`)**:
    - Introduced `Token` interface and concrete implementation.
    - New constructors: `New`, `NewAnd`, `NewOr`.
    - Supports inline expressions, named parameters, and operator/value inputs.
    - Enforces operator/value validation (`IN`, `NOT IN`, `BETWEEN` rules).
    - Implements all shared contracts (Kindable, Identifiable, Errorable, Debuggable, Rawable, Renderable, Stringable, Validable).
    - Full unit tests and examples with 100% coverage.

- **Join token (`join.Token`)**:
    - New safe constructors: `NewInner`, `NewLeft`, `NewRight`, `NewFull`.
    - Flexible DSL constructor: `New(kind, left, right, condition)`.
    - Added `Cross`, `Natural` join support.
    - Implements all core contracts.
    - **Breaking**: struct renamed `join` → `token`; `kind.go` deleted.

- **Field/Table tokens**:
    - Constructors delegate to `resolver.ValidateType`.
    - Clearer error messages and Clone() hints.
    - Validation for invalid/empty inputs, wildcards, subqueries.

---

## Helpers

- New **helpers** package: identifier/alias validation, wildcard restrictions, reserved keywords, expression and condition resolution.
- **Alias generation**:
    - `GenerateAlias(prefix, expr)` produces deterministic, safe aliases by combining a two-letter prefix with a SHA-1 hash of the expression.
- **Expression resolution**:
    - `ResolveExpressionType` classifies identifiers, literals, aggregates, functions, and computed expressions.
    - `ResolveExpression` simplified default branch with unified invalid error handling.
- **Condition resolution**:
    - `ResolveCondition` parses conditions into `(field, operator, value)`.
    - Bare identifiers default to `=` (e.g. `"id"` → `id = :id`).
    - Invalid expressions (e.g. `"id ++ 1"`) return clear errors.
- **Operator/value validation**:
    - `IsValidSlice` enforces consistency:
        - `IN` / `NOT IN` require non-empty slices.
        - `BETWEEN` requires exactly two values.
- **Utility helpers**:
    - `parseBetween`, `parseList`, `coerceScalar`, `ToParamKey`, `splitCSVRespectingQuotes`.

---

## Documentation & Tests

- Updated all `doc.go` and `README.md` across tokens and helpers.
- Expanded `example_test.go` with identifiers, aliases, wildcards, subqueries, generated aliases.
- 100% coverage for new types (`condition`, `identifier`, `operator`) and helpers.

---

## Breaking Changes

- Removed `join.Kind` → use `join.Type`.
- Renamed struct `join` → `token`.
- Deleted `kind.go`.
- Contracts updated in `contract.go` and `token.go`.

---

## Summary

This release consolidates **type safety** with dedicated enums, introduces a robust **operator registry**, and improves **validation and documentation** across tokens. Builders now have deterministic operator resolution, cleaner join APIs, and fully tested helpers.
