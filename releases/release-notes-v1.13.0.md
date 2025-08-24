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
    - `Build()` renders `GROUP BY` between `WHERE` and `ORDER BY`
  - **Ordering**
    - `OrderBy(...string)`: resets ordering fields
    - `ThenOrderBy(...string)`: appends ordering fields
    - Ignores empty or whitespace-only values
    - `Build()` renders `ORDER BY` after `WHERE` / `GROUP BY`

### Enhanced
- **Field diagnostics**
  - Invalid fields now produce consistent error messages
  - Example: `⛔️ Field("true"): input type unsupported: bool`
  - `Debug()` method on `token.Field` shows structured state with ✅/⛔️ markers
  - `String()` enhanced for clarity and consistency with Debug output
- **SelectBuilder reporting**
  - `Build()` aggregates invalid field errors into a single descriptive block
  - Clearer error messages when:
    - No source specified (`❌ [Build] - No source specified`)
    - Nil receiver (`❌ [Build] - Wrong initialization. Cannot build on receiver nil`)
  - `String()` provides status-style output:
    - ✅ successful SQL string with params
    - ❌ error message when build fails

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
- Comprehensive unit tests across `db/builder/select` and `common/extension/integer` ensuring 100% coverage:
  - Clause handling (nil, reset, append, overwrite, ignore-empty)
  - Field validation and error aggregation
  - Integer parsing (valid, invalid, defaults)

## Documentation
- **Database README.md** updated with Conditions, Grouping, Ordering, and Field Rules
- **Common/Extension README.md** updated with integer parser, usage table, and shortcuts
- **Database doc.go** extended with clause usage examples
- **Database example_test.go** enhanced with runnable examples:
  - `ExampleSelectBuilder_where`
  - `ExampleSelectBuilder_andOr`
  - `ExampleSelectBuilder_ordering`
  - `ExampleSelectBuilder_grouping`
- **Common example_test.go** enhanced with runnable examples:
  - `ExampleIntegerParseFrom` and shortcut usage
