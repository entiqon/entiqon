# Release Notes — v1.13.0

## Database Package (builder/select)

### Added
- Extended `SelectBuilder` with full clause support:
    - **Conditions**: `Where`, `And`, `Or` (resets, appends, normalized, ignores empty, renders after `FROM`)
    - **Grouping**: `GroupBy`, `ThenGroupBy` (reset/append, ignore empty, renders between `WHERE` and `HAVING`)
    - **Having**: `Having`, `AndHaving`, `OrHaving` (reset/append, normalized, ignore empty, renders after `GROUP BY`)
    - **Ordering**: `OrderBy`, `ThenOrderBy` (reset/append, ignore empty, renders after `WHERE`/`GROUP BY`/`HAVING`)

### Enhanced
- **Field diagnostics**
    - Consistent error messages (e.g., `⛔️ Field("true"): input type unsupported: bool`)
    - `Debug()` shows structured state with ✅/⛔️ markers
    - `String()` enhanced for clarity and UX
- **SelectBuilder reporting**
    - `Build()` aggregates invalid field errors into a single block
    - Clearer error messages (`❌ No source specified`, `❌ Wrong initialization`)
    - `String()` shows ✅/❌ style status

---

## Database Package (token/field)

### Refactored
- Moved `Field` into `db/token/field` subpackage
    - API preserved (`field.New(...)`)
    - Updated builder imports and Dockerfile
    - Normalized structure with `token/table`

### Added
- Introduced **Field Token contract**:
    - Aggregates `BaseToken`, `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, `Stringable`
    - Ownership methods: `HasOwner()`, `Owner()`, `SetOwner()`
    - ⚠️ Contract only, implementation staged

### Enhanced
- **Clone**
    - Deep copy of `owner` (avoids aliasing between clones and originals)
    - Removed unreachable `nil` branch (constructors guarantee non-nil)
    - Docstring updated
- **Role separation**
    - `Render()` → final SQL (future dialect-aware)
    - `Raw()` → SQL-generic (loggers)
    - `String()` → UX-friendly, concise with ✅/⛔
    - `Debug()` → developer diagnostics
- Normalized error handling with consistent `SetError` usage

### Documentation
- `doc.go`: extended **Field Behavior**
- `example_test.go`: placeholder `ExampleField_owner`
- `README.md`: new **Contracts and Auditability** section
- Root-level docs: `db/token/README.md` and `db/token/doc.go`

---

## Database Package (token/table)

### Added
- Introduced `token.Table` type:
    - Constructors for tables, aliases, raw inputs
    - Consistent rendering across dialects
    - Validation of invalid/empty inputs
- Added tests (constructors, methods, edge cases, 100% coverage)
- Added `doc.go`, `example_test.go`, `README.md`

---

## Database Package (contract)

### Added
- **BaseToken** interface (`db/contract/base_token.go`):
    - Identity and validation (`Input()`, `Expr()`, `Alias()`, `IsAliased()`, `IsValid()`)
- Example: `ExampleBaseToken` in `example_test.go`
- `doc.go` and `README.md` updated with BaseToken
- Added **Consistency** principle in doctrine
- Extended **Errorable** with `SetError(err error)`
    - Implemented in `Field` and `Table`
    - Docs and examples updated

---

## Common Package (extension/integer)

### Added
- Integer parser:
    - `ParseFrom(any)` with strict type validation
    - `IntegerOr` shortcut with defaults
- Consistent with boolean, float, decimal parsers
- 100% test coverage
- Runnable examples
- Docs updated under `common/extension`

---

## Tests
- Full coverage across Database and Common:
    - Clause handling
    - Field validation, cloning, error aggregation
    - Table token construction & rendering
    - Integer parsing (valid, invalid, defaults)

---

## Documentation
- **Database README.md** refactored:
    - Added **Doctrine** section
    - Replaced bullet list with **Capabilities table** (Modules, Features, Status icons)
    - Removed redundant Quickstart and guides
- **Database doc.go** extended with field & table examples
- **Database example_test.go** updated with runnable examples:
    - Select (where, and/or, ordering, grouping, having)
    - Table usage
    - Clonable
- **Common/Extension README.md** updated with integer parser usage
- **Common example_test.go** extended with parser examples
