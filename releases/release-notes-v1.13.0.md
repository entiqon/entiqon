# Release v1.13.0 — Builder Conditions & Field Diagnostics

## 🚀 Features

- **db/builder/select**
  - Introduced `Where()`, `And()`, and `Or()` methods for building `WHERE` clauses
  - Normalization of multiple conditions:  
    `Where("a", "b")` → `WHERE a AND b`
  - `Where()` clears existing conditions (like `Fields()`), while `And()`/`Or()` append
  - Defensive handling of manual initialization (`&SelectBuilder{}`)
  - Simplified `Build()` logic for condition rendering

- **token/field**
  - Added `Debug()` for compact diagnostic output:  
    `✅ Field("COUNT(*) AS total"): [raw: true, aliased: true, errored: false]`  
    `⛔️ Field("true"): [raw: false, aliased: false, errored: true] – input type unsupported: bool`
  - Enhanced `String()` with ✅/⛔️ icons and explicit *wrong initialization* handling

## 🧪 Tests

- 100% coverage for `SelectBuilder` and `token.Field`
- Added test cases for:
  - Defaults, single/multiple conditions, ignore empty
  - `Where` reset behavior
  - `And`/`Or` appends
  - Edge case: mixed `AND` + `OR`
  - Manual initialization (`&SelectBuilder{}`)

## 📖 Documentation

- Updated `doc.go` for `builder` with conditions section
- Added new `example_test.go` examples covering `Where`, `And`, `Or`, and reset semantics
- Updated `field.md` guide with new `String()`/`Debug()` documentation
- Updated `README.md` with:
  - Strict **Field Rules**
  - Examples of `Where`/`And`/`Or`
  - Debugging output (`String()`/`Debug()`)
  - Revised error cases
  - Updated Status section

---

## 📄 Summary

This release introduces **full WHERE clause support** with `Where`/`And`/`Or` methods,  
improves **diagnostics for fields** with ✅/⛔️ outputs, and ensures **100% test coverage**.  
All docs (README, guides, examples) have been updated accordingly.
