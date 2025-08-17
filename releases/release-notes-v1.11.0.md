# v1.11.0 --- Parsing shortcuts, deterministic date cleaning, and SQL token upgrades

## Highlights

-   **New parsing façade**: one-line helpers for `Boolean`, `Float`,
    `Decimal`, and `Date`.
-   **Deterministic date cleaning**: `CleanAndParse`,
    `CleanAndParseAsString`, strict `YYYYMMDD` prefix path, and 100%
    tests.
-   **Boolean parser++**: extended tokens (`on/off`, `y/n`, `t/f`) and
    explicit `nil` rejection.
-   **SQL Builder/Token overhaul**: `Column → Field`, `FieldCollection`,
    deterministic `NewField` inputs, and Postgres/Base dialects with
    tests.

------------------------------------------------------------------------

## What changed since v1.10.0

### ✨ Features

-   **common/extension**
    -   `b720f7f` feat: add parser shortcuts for **Boolean**, **Float**,
        **Decimal**, **Date**
-   **common/date**
    -   `dd173da` feat: introduce **Cleaning**, **ParseAndFormat**, and
        **ParseFrom** with full coverage
-   **contract / token / builder / dialect**
    -   `b2b47ca` feat(contract): add generic **Cloanable** interface;
        assert `Field` compliance
    -   `d140010` feat(token): **FieldCollection** with full ops + tests
    -   `d4c0905` feat(token): implement **Column.Render** and satisfy
        `Renderable` with tests
    -   `2ce91f1` feat(token): initial **Column** token with full unit
        tests
    -   `ca22e45` feat(builder): first basic **select** builder
    -   `ccdb48f` feat(dialect): **PostgresDialect** with PG-specific
        syntax + tests
    -   `ff00ca6` feat(dialect): **BaseDialect** baseline implementation

### ♻️ Refactors

-   `f45ff2a` refactor(common): rename **math package → extension**
-   `c8ca52f` refactor(token): make **Field** constructor inputs
    deterministic; improved raw-alias handling
-   `45da5dc` refactor(token): **rename Column → Field** and add tests
-   `9fc782f` refactor(builder): improve `builder.Select`, delegate
    field resolution to `token.Field`, add `select.GetFields()`
-   `e722fdc` refactor(builder): move `builder` to **internal/builder**

### 🛠️ Fixes

-   `3aabe8c` fix(entiqon): update real project logo

### 📚 Docs

-   `9907b5c` docs(contract): full **GoDoc** for `Renderable`
-   `7f94040` docs(entiqon): update real project logo & disposition

------------------------------------------------------------------------

## Breaking / Migration notes

-   **math** → **extension** package rename
-   **token.Column** → **token.Field**
-   `builder.Select` changes: use `select.Fields(...)` and
    `select.GetFields()`
-   Date parsing: use `CleanAndParse` with options, or
    `StrictYYYYMMDDOptions()` for strict feeds

------------------------------------------------------------------------

## Test coverage highlights

-   Full coverage for `CleanAndParse`, `CleanAndParseAsString`,
    `ParseFrom`
-   Boolean parser extended token cases covered
-   Field, FieldCollection, Dialects fully unit-tested

------------------------------------------------------------------------

## Upgrade checklist

-   Update imports: `common/math/...` → `common/extension/...`
-   Replace `token.Column` with `token.Field`
-   For date parsing: switch to `CleanAndParse` or strict options
