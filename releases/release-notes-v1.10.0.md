
# Release vX.Y.Z — Object Package Enhancements and Test Improvements

## Highlights

- **Enhanced `Exists` function** to accept any input type:
  - Supports both `map[string]any` and structs (including pointers).
  - Performs case-insensitive key or field lookup.
  - Skips unexported struct fields for accurate existence checks.
- **Improved `GetValue` robustness:**
  - Handles edge cases like nil objects, unsupported kinds, pointer nils.
  - Case-insensitive struct field matching.
- **Expanded `SetValue` coverage and correctness:**
  - Supports setting values with convertible types (e.g., `int32` to `int`).
  - Returns errors on nil pointers, pointers to non-structs, non-assignable values, and maps with non-string keys.
- **Comprehensive test suite updates:**
  - Added tests covering maps, structs, pointers, nils, and error conditions for `Exists`, `GetValue`, and `SetValue`.
  - Increased code coverage to nearly 100% on all object package functions.
- **Documentation updates:**
  - Updated package overview and function descriptions for better clarity and accuracy.
  - Provided usage examples reflecting the enhanced flexibility and behavior.

## Impact

- These improvements provide a more consistent, robust, and flexible API for dynamic object manipulation.
- Backwards-compatible enhancements ensure existing code continues to work without changes.
- Increased test coverage improves reliability and maintainability.

## Upgrade Notes

- No breaking API changes; just upgrade and enjoy the new functionality.
- Review updated examples for better integration patterns.
