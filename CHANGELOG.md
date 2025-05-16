# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/)
and this project adheres to [Semantic Versioning](https://semver.org/).

---

## [v0.2.0] - 2024-05-15

### Added
- Full `SelectBuilder` implementation with fluent chaining
- Support for logical condition grouping (`Where`, `AndWhere`, `OrWhere`)
- Pagination via `Take` (LIMIT) and `Skip` (OFFSET)
- `ConditionToken` and `ConditionType` struct for WHERE clause abstraction

### Changed
- Method names normalized for readability and control
- `Build()` method added to generate validated SQL string
- All methods and types fully documented using GoDoc
- Updated `README.md` with example usage and supported features

### Tests
- Complete test suite for `SelectBuilder` with `testify/suite`

---

## [v0.1.0] - 2024-04-02

### Added
- Initial `SelectQueryBuilder` prototype with basic `Select`, `From`, `Where`, `OrderBy`, `Limit`, `Offset`
- Query construction using simple string chaining
- Basic test coverage for SELECT behavior

### Changed
- Early param binding experiments for WHERE clauses
- Initial logo and branding
- Prepared project for pkg.go.dev publication

### Known Gaps
- No structured condition typing
- No dialect interface
- No documentation or strong API separation
