# Release v1.9.0 — Float and Decimal Module Improvements

**Release Date:** 2025-08-07

## Highlights

- **Float module:**
    - Removed redundant recursive unwrapping of interface kinds in `ParseFrom`.
    - Simplified parsing logic by relying on Go’s `reflect.ValueOf` automatic interface unwrapping.
    - Improved maintainability and readability of the parsing code.
    - Added comprehensive tests with 100% coverage covering native types, pointers, and edge cases.

- **Decimal module:**
    - Enhanced parsing logic for decimal numbers.
    - Improved test coverage and validation.
    - Better handling of precision and rounding in parsing functions.

- **Number to Math migration:**
    - Migrated core parsing utilities and numeric logic from the legacy `number` package into the consolidated `math` package.
    - Updated import paths and dependencies across the codebase.
    - Streamlined numeric processing by unifying shared logic, reducing duplication, and improving maintainability.

## Benefits

- Cleaner, easier-to-maintain codebase for numeric parsing utilities.
- More reliable parsing behavior with thorough test coverage.
- Reduced complexity and eliminated unreachable code paths.
- Foundation for future improvements in number handling across modules.

---

© 2025 Entiqon Project — Inspired by Mythology & Legends