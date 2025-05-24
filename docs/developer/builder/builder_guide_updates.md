## SelectBuilder Guide

### Overview
The `SelectBuilder` allows fluent construction of SELECT queries. It now fully supports dialect-aware placeholder resolution and validation tagging via `StageToken`.

### Clause Mapping
- `SELECT`: Fields projection (`StageSelect`)
- `FROM`: Source tables (`StageFrom`)
- `JOIN`: Join logic (`StageJoin`)
- `WHERE`: Condition logic (`StageWhere`)
- `GROUP BY`: Aggregation grouping (`StageGroup`)
- `HAVING`: Filter after grouping (`StageHaving`)
- `ORDER BY`: Sorting (`StageOrder`)
- `LIMIT` / `OFFSET`: Pagination (`StageLimit`, `StageOffset`)

### New Features
- Dialect-aware placeholder injection using `ParamBinder`
- StageToken tagging for precise error localization
- 94%+ test coverage

---

## InsertBuilder Guide

### Overview
`InsertBuilder` builds safe, parameterized INSERT queries. Now aligned with the unified builder pattern.

### Clause Mapping
- `INTO`: Target table (`StageInto`)
- `VALUES`: Inserted values (`StageValues`)
- `RETURNING`: Optional return clause (`StageReturning`)

### New Features
- Uses `ParamBinder` for dialect-specific placeholder formatting
- Error tagging via `StageInto`, `StageValues`, `StageReturning`
- Upsert support with dialect validation (e.g., PostgreSQL)

---

## UpdateBuilder Guide

### Overview
`UpdateBuilder` builds UPDATE statements with support for conditional updates and optional returns.

### Clause Mapping
- `SET`: Field assignments (`StageSet`)
- `WHERE`: Filtering logic (`StageWhere`)
- `RETURNING`: Return affected rows (`StageReturning`)

### New Features
- Strict field-to-value assignment (no aliasing)
- Full dialect placeholder formatting using `ParamBinder`
- `StageToken` tagging for granular diagnostics

---

## DeleteBuilder Guide

### Overview
Builds DELETE statements with conditional logic.

### Clause Mapping
- `FROM`: Source table (`StageFrom`)
- `WHERE`: Filtering (`StageWhere`)
- `RETURNING`: Optional response (`StageReturning`)

### New Features
- Safe condition resolution via `Condition`
- Dialect-specific parameter handling
- Stage-based validation for traceable feedback

---

## Common Builder Notes

### Dialect Exposure
All builders now rely on a shared `Dialect` interface exposed under `driver/`, which defines:
- `QuoteIdentifier`
- `Placeholder`
- `SupportsReturning`, `SupportsUpsert`
- `Validate()` for contract checks

### StageToken
`StageToken` is used to annotate clauses during build-time validation. It supports:
- Selective error reporting
- Declarative stage labeling for error propagation

### ParamBinder
Used across all builders to insert parameters according to the active Dialect.
- Example: `?` (MySQL), `$1` (Postgres)
- Resolves consistency and dialect abstraction

---

Builders are now modular, test-covered, and validation-aware – ensuring safety across SQL generation pipelines.
