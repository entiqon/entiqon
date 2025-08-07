# Release v1.8.4 — Common/number Module Enhancements

**Release Date:** 2025-08-07

## Highlights

- **Enhanced `ParseFrom` utility with improved float parsing from strings**  
  Now supports parsing long-form float strings like `"1.000000000000003"` seamlessly.

- **Introduced rounding mode flag to `ParseFrom`**  
  Added a `round bool` parameter to control float handling:  
  - `round=false` enforces strict validation allowing only floats close to integers within a small tolerance (1e-9).  
  - `round=true` enables lenient parsing by rounding floats to the nearest integer.

- **Updated string parsing logic**  
  When integer parsing fails, `ParseFrom` attempts float parsing with the configured rounding behavior, improving robustness for dynamic numeric inputs.

- **Expanded comprehensive test coverage**  
  Tests now cover all supported input types, rounding modes, and edge cases ensuring high reliability.

## Impact

These improvements enhance the flexibility and accuracy of dynamic numeric data parsing across Entiqon modules, reducing errors and supporting more varied input formats.

## Migration

- The signature of `ParseFrom` changed from `ParseFrom(value interface{}) (int, error)` to  
  `ParseFrom(value interface{}, round bool) (int, error)` — update all calls accordingly.

- Choose the rounding mode (`round` flag) based on your application’s tolerance for float-to-int conversion.

---

© 2025 Entiqon Project — Inspired by Mythology & Legends