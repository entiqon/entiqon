# Token Helpers

The `helpers` package provides small, reusable validation utilities
for SQL tokens. These functions encapsulate common rules that are
used across multiple token types (Field, Table, Join, etc.).

## Why Helpers?

- **Single Responsibility** — keep low-level validation logic out of
  the main token constructors.
- **Reusability** — shared across Field, Table, Resolver, and future
  builders.
- **Testability** — each helper is tested independently with full
  coverage.
- **Consistency** — all helpers follow the same pattern:
    - `ValidateXxx(s string) error` → returns rich error messages.
    - `IsValidXxx(s string) bool` → convenience wrapper for quick checks.
    - `GenerateAlias(prefix, expr string)` → produces deterministic,
      safe aliases for non-identifier expressions.
    - `ResolveExpressionType(expr string)` → classifies SQL expressions
      into identifier categories.
- **Future-proof** — current rules are dialect-agnostic, but dialect
  packages will later override them with grammar-specific rules.

## Current Helpers

- `IsValidIdentifier` / `ValidateIdentifier`  
  Strict checks for valid SQL identifiers
    - Must not be empty.
    - Must start with a letter (A–Z, a–z) or underscore (`_`).
    - Remaining characters may be letters, digits, or underscores.
    - Non-ASCII identifiers (e.g. `café`, `mañana`) are rejected until
      dialect-specific rules are introduced.

- `IsValidAlias` / `ValidateAlias`  
  Checks if a string is a valid alias.
    - Must be a valid identifier.
    - Must not be a reserved keyword (case-insensitive).

- `HasTrailingAlias` / `ValidateTrailingAlias`  
  Detects and validates trailing aliases (e.g. `(price * qty) total`).
    - Ignores explicit `AS` (handled by the resolver).
    - Ensures last token is not part of the expression.
    - Rejects reserved keywords or invalid alias syntax.

- `ReservedKeywords`  
  Returns the current dialect-agnostic set of SQL keywords disallowed
  as aliases. Dialects may extend or override this list.

- `GenerateAlias`  
  Produces deterministic aliases for non-identifier expressions.
    - Combines a two-letter kind code (from `identifier.Type.Alias()`)
      with a SHA-1 hash of the expression.
    - Always returns a safe SQL identifier.
    - Example: `GenerateAlias("fn", "SUM(price)") → "fn_a1b2c3"`.

- `ValidateWildcard`  
  Ensures the `*` wildcard is used only in valid contexts.
    - Bare `*` without alias is allowed.
    - Rejects aliased or raw `*` (e.g. `* AS total`).
    - Future: may extend to handle qualified wildcards (e.g. `table.*`).

- `ResolveExpressionType`  
  Classifies raw SQL expressions into broad categories, returning an `identifier.Type`.
    - Categories: `Invalid`, `Subquery`, `Computed`, `Aggregate`, `Function`, `Literal`, `Identifier`.
    - Purely syntactic, not semantic (e.g. `SUM(qty)` is classified as Aggregate even if used in an invalid context).
    - Aliases must already be stripped before classification.
    - Example:
      ```go
      helpers.ResolveExpressionType("SUM(price)")        // → Aggregate
      helpers.ResolveExpressionType("(a+b)")             // → Computed
      helpers.ResolveExpressionType("(SELECT * FROM t)") // → Subquery
      helpers.ResolveExpressionType("users")             // → Identifier
      ```

## Roadmap

As more helpers are promoted (e.g. literal checks), they will be added
here with their own independent tests and will follow the same
**Validate/IsValid/Generate/Resolve** pattern for consistency.

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project
