# Release v1.12.0 — Documentation & Parser Shortcuts

## Highlights

- **Parser Shortcuts**: Added `Or` variants (`BooleanOr`, `NumberOr`, `FloatOr`, `DecimalOr`, `DateOr`) to allow explicit fallback values.  
- **Extension Documentation**: Each extension subpackage (`boolean`, `date`, `decimal`, `float`, `number`, `object`, `collection`) now ships with `README.md`, `doc.go`, and `example_test.go`.  
- **Object Helpers**: Normalized and moved from `common/object` → `common/extension/object`.  
- **Deprecations**:  
  - `BoolToStr` has been renamed to `BoolToString` and moved into `common/extension/boolean`.  
  - Existing usages of `BoolToStr` will continue to work but are discouraged.

## Documentation

- Added root-level `doc.go` for `common` describing purpose and structure.  
- Added `README.md` for `common`, linking subpackages `errors` and `extension`.  
- Added package-level READMEs and examples across all `common/extension` subpackages.  
- Dockerfile updated to copy package-level READMEs into site docs.  
- Navigation updated with alphabetized extension packages.

## Internal

- Unified test coverage for all parsing helpers and `Or` fallbacks.  
- Refactored Dockerfile to preserve base image but copy documentation automatically.  
- Clarified supported input formats in `date.ParseAndFormat` examples.  
