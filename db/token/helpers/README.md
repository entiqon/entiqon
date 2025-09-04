# Token Helpers

The `helpers` package provides small, reusable validation and resolution
utilities for SQL tokens. These functions encapsulate common rules that are
used across multiple token types (Field, Table, Join, Condition, etc.).

## Why Helpers?

- **Single Responsibility** — keep low-level validation logic out of
  the main token constructors.
- **Reusability** — shared across Field, Table, Condition, Resolver, and future
  builders.
- **Testability** — each helper is tested independently with full
  coverage.
- **Consistency** — all helpers follow the same pattern:
    - `ValidateXxx(s string) error` → returns rich error messages
    - `IsValidXxx(s string) bool` → convenience wrapper for quick checks
    - `GenerateAlias(prefix, expr string)` → produces deterministic,
      safe aliases
    - `ResolveExpressionType(expr string)` → classifies SQL expressions
    - `ResolveCondition(expr string)` → parses condition expressions into field/operator/value
- **Future-proof** — current rules are dialect-agnostic, but dialect
  packages will later override them with grammar-specific rules.

## Current Helpers

- **Identifiers**
    - `IsValidIdentifier` / `ValidateIdentifier`  
      Strict checks for valid SQL identifiers.
        - Must not be empty.
        - Must start with a letter (A–Z, a–z) or underscore (`_`).
        - Remaining characters may be letters, digits, or underscores.
        - Non-ASCII identifiers are rejected until dialect-specific rules are introduced.

- **Aliases**
    - `IsValidAlias` / `ValidateAlias`  
      Checks if a string is a valid alias.
        - Must be a valid identifier.
        - Must not be a reserved keyword (case-insensitive).

    - `HasTrailingAlias` / `ValidateTrailingAlias`  
      Detects and validates trailing aliases (e.g. `(price * qty) total`).
        - Ignores explicit `AS` (handled by the resolver).
        - Rejects reserved keywords or invalid alias syntax.

    - `ReservedKeywords`  
      Returns the current dialect-agnostic set of SQL keywords disallowed
      as aliases. Dialects may extend or override this list.

    - `GenerateAlias`  
      Produces deterministic aliases for non-identifier expressions.
        - Combines a two-letter kind code (from `identifier.Type.Alias()`)
          with a SHA-1 hash of the expression.
        - Always returns a safe SQL identifier.
        - Example:
          ```go
          helpers.GenerateAlias("fn", "SUM(price)") // → "fn_a1b2c3"
          ```

- **Wildcards**
    - `ValidateWildcard`  
      Ensures the `*` wildcard is used only in valid contexts.
        - Bare `*` without alias is allowed.
        - Rejects aliased or raw `*` (e.g. `* AS total`).
        - Future: may extend to handle qualified wildcards (e.g. `table.*`).

- **Expression Resolution**
    - `ResolveExpressionType`  
      Classifies raw SQL expressions into broad categories: `Invalid`, `Subquery`,
      `Computed`, `Aggregate`, `Function`, `Literal`, `Identifier`.

    - `ResolveExpression`  
      Splits an expression into its kind, core expression, and optional alias.

- **Condition Resolution**
    - `ResolveCondition`  
      Parses SQL-like conditions into `(field, operator, value)` triples.
        - `"id = 1"` → field=`id`, op=`=`, value=`1`
        - `"price BETWEEN 1 AND 10"` → field=`price`, op=`BETWEEN`, value=`[1,10]`
        - `"lastname IN ('a','b')"` → field=`lastname`, op=`IN`, value=`["a","b"]`
        - `"deleted_at IS NULL"` → field=`deleted_at`, op=`IS NULL`, value=`nil`
        - `"id"` → defaults to `id = :id`, op=`=`, value=`nil`

    - `IsValidSlice`  
      Validates operator/value consistency:
        - `IN` / `NOT IN` → non-empty slice
        - `BETWEEN` → exactly two elements

- **Low-Level Parsers**
    - `parseBetween` → splits `"x AND y"` into `[x, y]`
    - `parseList` → parses CSV or `(a, b, c)` lists into `[]any`
    - `coerceScalar` → coerces tokens into `int`, `float64`, `nil`, or `string`
    - `ToParamKey` → converts identifiers into safe parameter keys (`users.id → users_id`)
    - `splitCSVRespectingQuotes` → splits on commas but preserves quoted text

## Roadmap

As more helpers are promoted (e.g. dialect-specific literals or operators),
they will be added here with their own independent tests and will follow the
same **Validate/IsValid/Generate/Resolve** pattern for consistency.

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project
