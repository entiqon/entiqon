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
- **Future-proof** — current rules are dialect-agnostic, but dialect
  packages will later override them with grammar-specific rules.

## Current Helpers

- `IsValidIdentifier` — strict check for valid SQL identifiers

## Roadmap

As more helpers are promoted (alias validation, trailing alias
detection, literal checks, etc.), they will be added here with their
own independent tests.

---

## 📄 License

[MIT](../../../LICENSE) — © Entiqon Project