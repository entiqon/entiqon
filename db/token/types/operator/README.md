# Operator

> Part of [Entiqon](../../../../) / [Database](../../../) / [Token](../../) / [Types](../)

**Part of Entiqon / Database**

Typed enum and helpers for SQL operators used in `WHERE`, `ON`, and `HAVING` clauses.  
Provides parsing, validation, canonical string forms, aliases, and discovery for resolvers.

---

## Purpose

Unify how Entiqon represents and handles SQL operators across tokens and resolvers:

- Strongly-typed operator enum (`Type`)
- Canonical `String()` (e.g., `"NOT IN"`, `">="`, `"IS NULL"`)
- Short `Alias()` mnemonics (e.g., `nin`, `gte`, `isnull`)
- Robust `ParseFrom(any)` (symbols, words, and aliases)
- `GetKnownOperators()` for resolver scans (longest-first ordering)

---

## Supported Operators

| Constant              | `String()`              | `Alias()`      | Category            | Notes |
|----------------------|-------------------------|----------------|---------------------|------|
| `Equal`              | `=`                     | `eq`           | Comparison          |      |
| `NotEqual`           | `!=`                    | `neq`          | Comparison          | Accepts `<>` on parse |
| `GreaterThan`        | `>`                     | `gt`           | Comparison          |      |
| `GreaterThanOrEqual` | `>=`                    | `gte`          | Comparison          |      |
| `LessThan`           | `<`                     | `lt`           | Comparison          |      |
| `LessThanOrEqual`    | `<=`                    | `lte`          | Comparison          |      |
| `In`                 | `IN`                    | `in`           | Membership          |      |
| `NotIn`              | `NOT IN`                | `nin`          | Membership          |      |
| `Between`            | `BETWEEN`               | `between`      | Range               |      |
| `Like`               | `LIKE`                  | `like`         | Pattern             |      |
| `NotLike`            | `NOT LIKE`              | `nlike`        | Pattern             |      |
| `IsNull`             | `IS NULL`               | `isnull`       | Nullness            |      |
| `IsNotNull`          | `IS NOT NULL`           | `notnull`      | Nullness            |      |
| `IsDistinctFrom`     | `IS DISTINCT FROM`      | `isdistinct`   | Set Distinctness    | PostgreSQL |
| `NotIsDistinctFrom`  | `IS NOT DISTINCT FROM`  | `notdistinct`  | Set Distinctness    | PostgreSQL |

> `Invalid` exists for unrecognized inputs.

---

## API Summary

```go
type Type int

func (t Type) String() string
func (t Type) Alias() string
func (t Type) IsValid() bool

func ParseFrom(v any) Type
func GetKnownOperators() []string
```

- `String()` → canonical SQL spelling (used for rendering)
- `Alias()` → stable lowercase mnemonic (good for JSON, logs, flags)
- `IsValid()` → bounds check
- `ParseFrom()` → accepts `Type`, `string`, `[]byte` (case/space-insensitive; symbols, words, aliases)
- `GetKnownOperators()` → canonical strings, **sorted longest-first** (helps when scanning text so `"IS NOT DISTINCT FROM"` wins over `"IS"`)

---

## Usage

### Parse And Render

```go
op := operator.ParseFrom("  not   in ")
if !op.IsValid() {
    // handle unsupported operator
}
fmt.Println(op.String()) // "NOT IN"
fmt.Println(op.Alias())  // "nin"
```

### With A Condition Resolver

Many Entiqon resolvers normalize the LHS while carrying the real operator and typed RHS:

```go
// Input: "COUNT(id) > 0"
expr, op, rhs, err := helpers.ResolveCondition("COUNT(id) > 0")
if err != nil { /* ... */ }

// Canonicalized expression for binders (placeholder always lowercased):
// "COUNT(id) > :count_id"
fmt.Println(expr)        // "COUNT(id) > :count_id"
fmt.Println(op.String()) // ">"
fmt.Printf("%v\n", rhs)  // 0 (int)
```

Other examples:

- `id BETWEEN 1 AND 3` → `expr="id = :id"`, `op=BETWEEN`, `rhs=[1 3]`
- `lastname IN ('a','b')` → `expr="lastname = :lastname"`, `op=IN`, `rhs=["a" "b"]`
- `deleted_at IS NULL` → `expr="deleted_at = :deleted_at"`, `op=IS NULL`, `rhs=nil`

> Tip: use `operator.GetKnownOperators()` to scan raw input; match in that order for correct multi-word precedence.

---

## Testing

Follow the project pattern **file → methods → cases** with PascalCase test names:

- `type_test.go`:
  - `TestType/String`
  - `TestType/Alias`
  - `TestType/IsValid`
  - `TestParseFrom/Symbols`
  - `TestParseFrom/Words`
  - `TestParseFrom/Aliases`
  - `TestGetKnownOperators/LongestFirst`

Include runnable examples in `example_test.go` (GoDoc examples).

---

## Integration Notes

- Prefer carrying `Type` alongside your normalized expression and typed RHS; render time uses `op.String()`.
- When scanning free text, always favor **longest-first** matching (use `GetKnownOperators()` result).
- Accept both `!=` and `<>` for inequality on parse; emit `!=` on `String()`.

---

## License

This module is part of the Entiqon repository. See the root **LICENSE** file.
