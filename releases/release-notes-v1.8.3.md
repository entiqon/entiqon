# Release v1.8.3 — Safe Map Access and Mutation Utilities

**Release Date:** 2025-08-06

### Added

- **object** package:
  - Added generic utility **Exists** to check safely if a key exists in a map.
  - Added generic **GetValue** to retrieve values with type safety and a default fallback if the key is missing or type mismatched.
  - Added **SetValue** to safely initialize maps if nil and update keys only when the value changes.
- Added comprehensive tests covering `Exists`, `GetValue`, and `SetValue`, including edge cases for default values and type mismatches.

---

© 2025 Entiqon Project — Inspired by Mythology & Legends