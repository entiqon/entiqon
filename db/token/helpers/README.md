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
- **Future-proof** — current rules are dialect-agnostic, but dialect
  packages will later override them with grammar-specific rules.

## Current Helpers

- `IsValidIdentifier` / `ValidateIdentifier` — strict checks for valid SQL identifiers
    - Must not be empty.
    - Must start with a letter (A–Z, a–z) or underscore (`_`).
    - Remaining characters may be letters, digits, or underscores.
    - Non-ASCII identifiers (e.g. `café`, `mañana`) are rejected until
      dialect-specific rules are introduced.

## Roadmap

As more helpers are promoted (alias validation, trailing alias
detection, literal checks, etc.), they will be added here with their
own independent tests and will follow the same **Validate/IsValid**
pattern for consistency.

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project
