# Release Notes — v1.13.0

## Database Package (builder/select)

### Added
- Extended `SelectBuilder` with full clause support:
    - **Conditions**
        - `Where(...string)`: resets conditions
        - `And(...string)`: appends with `AND`
        - `Or(...string)`: appends with `OR`
        - Multiple conditions in a single `Where` call are normalized with `AND`
        - Ignores empty or whitespace-only inputs
        - `Build()` renders conditions immediately after `FROM`
    - **Grouping**
        - `GroupBy(...string)`: resets grouping fields
        - `ThenGroupBy(...string)`: appends grouping fields
        - Graceful handling of nil collections
        - Ignores empty or whitespace-only values
        - `Build()` renders `GROUP BY` between `WHERE` and `HAVING`
    - **Having**
        - `Having(...string)`: resets HAVING conditions
        - `AndHaving(...string)`: appends with `AND`
        - `OrHaving(...string)`: appends with `OR`
        - Multiple conditions in a single `Having` call are normalized with `AND`
        - Ignores empty or whitespace-only inputs
        - `Build()` renders `HAVING` immediately after `GROUP BY`
    - **Ordering**
        - `OrderBy(...string)`: resets ordering fields
        - `ThenOrderBy(...string)`: appends ordering fields
        - Ignores empty or whitespace-only values
        - `Build()` renders `ORDER BY` after `WHERE` / `GROUP BY` / `HAVING`

### Enhanced
- **Field diagnostics**
    - Invalid fields now produce consistent error messages
    - Example: `⛔️ Field("true"): input type unsupported: bool`
    - `Debug()` method on `token/field.Field` shows structured state with ✅/⛔️ markers
    - `String()` enhanced for clarity and consistency with Debug output
- **SelectBuilder reporting**
    - `Build()` aggregates invalid field errors into a single descriptive block
    - Clearer error messages when:
        - No source specified (`❌ [Build] - No source specified`)
        - Nil receiver (`❌ [Build] - Wrong initialization. Cannot build on receiver nil`)
    - `String()` provides status-style output:
        - ✅ successful SQL string with params
        - ❌ error message when build fails

## Database Package (contract)

### Added
- Introduced **BaseToken** interface (`db/contract/base_token.go`):
    - Provides core identity and validation for all tokens (`Field`, `Table`, etc.)
    - Methods: `Input()`, `Expr()`, `Alias()`, `IsAliased()`, `IsValid()`
    - Guarantees consistent identity and validation across all tokens
- Added runnable example (`ExampleBaseToken`) in `example_test.go`
- Updated `doc.go` with BaseToken in the contract overview (normalized style)
- Updated `README.md`:
    - New BaseToken section with purpose, methods, usage
    - Streamlined documentation structure
    - Extended philosophy with **Consistency** principle: all tokens share BaseToken
- Extended **Errorable** contract with `SetError(err error)`:
    - Enables tokens/builders to mark themselves as errored after construction
    - Implemented in `Field` and `Table` tokens
    - Updated `doc.go`, `README.md`, and `example_test.go` with usage examples

## Database Package (token/table)

### Added
- Introduced `token.Table` type to represent SQL sources in builders:
    - Provides constructors for tables, aliases, and raw inputs
    - Consistent rendering across dialects
    - Validation of invalid/empty inputs with clear error reporting
- Added full unit test coverage (constructors, methods, edge cases)
- Added `doc.go` with package overview and usage guidelines
- Added `example_test.go` with runnable examples
- Added `README.md` describing purpose, design, and usage

## Database Package (token/field)

### Refactored
- Moved `Field` into a dedicated subpackage `db/token/field`:
    - API preserved (`field.New(...)`) with unchanged contracts and behavior
    - Updated `builder/select.go` and `select_test.go` to use new import path
    - Dockerfile updated to copy `db/token/field/README.md` for documentation
    - Normalized structure for consistency with `token/table`

### Added
- Introduced **Field Token contract** as a scaffold to decompose Field identity into auditable pieces:
    - Aggregates `BaseToken`, `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, and `Stringable`
    - Defines ownership methods: `HasOwner()`, `Owner()`, and `SetOwner()`
    - Intention: separate every identity aspect (expr, alias, owner, validity, raw state) into dedicated, auditable contracts
    - ⚠️ Contract only; implementation staged for later commits

### Documentation
- `doc.go`: extended **Field Behavior** with `HasOwner`, `Owner`, `SetOwner`
- `example_test.go`: added placeholder `ExampleField_owner` (commented until implemented)
- `README.md`: new **Contracts and Auditability** section in Developer Guide
- Root-level docs:
    - Added `db/token/README.md` with package purpose and subpackage table
    - Added `db/token/doc.go` with GoDoc overview, principles, subpackages, and roadmap

### Refactored
- Moved `Field` into a dedicated subpackage `db/token/field`:
    - API preserved (`field.New(...)`) with unchanged contracts and behavior
    - Updated `builder/select.go` and `select_test.go` to use new import path
    - Dockerfile updated to copy `db/token/field/README.md` for documentation
    - Normalized structure for consistency with `token/table`

## Common Package (extension/integer)

### Added
- Introduced integer parser with full support:
    - `ParseFrom(any)` converts generic input into integer values
    - Rejects non-integer inputs with clear error reporting
    - Consistent behavior with existing parsers (`boolean`, `float`, `decimal`)
    - Added parser shortcuts (`IntegerOr`) for default values
- Full test coverage in `integer/parser_test.go`
- Added runnable examples and integrated into `example_test.go`
- Updated parser documentation under `common/extension` README

## Tests
- Comprehensive unit tests across `db/builder/select`, `token/table`, `token/field`, and `common/extension/integer` ensuring 100% coverage:
    - Clause handling (nil, reset, append, overwrite, ignore-empty)
    - Field validation and error aggregation
    - Table token construction, aliasing, and rendering
    - Integer parsing (valid, invalid, defaults)
    - Field cloning, rendering, and error detection

## Documentation
- **Database README.md** updated with Conditions, Grouping, Having, Ordering, Field Rules, and Table token section
- **Common/Extension README.md** updated with integer parser, usage table, and shortcuts
- **Database doc.go** extended with clause usage and table/field examples
- **Database example_test.go** enhanced with runnable examples:
    - `ExampleSelectBuilder_where`
    - `ExampleSelectBuilder_andOr`
    - `ExampleSelectBuilder_ordering`
    - `ExampleSelectBuilder_grouping`
    - `ExampleSelectBuilder_having`
    - `ExampleTable_basic`
    - `ExampleClonable`
- **Common example_test.go** enhanced with runnable examples:
    - `ExampleIntegerParseFrom` and shortcut usage
