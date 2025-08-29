# 🚀 Release v1.14.0 – Contract & Token Evolution

This release refines the **contract layer** across `db/token` and introduces key supporting modules (`resolver`, `ExpressionKind`) along with the new **Join token**.  
The changes strengthen type safety, alias validation, and parsing, ensuring higher reliability for builders (notably `SelectBuilder`) which directly depend on tokens.

---

## ✨ Highlights

### New `Validable` contract
- Introduced **`Validable`** with `IsValid()` for structural validation.  
- Adopted by `table.Token`, `field.Token`, and `join.Token`.  
- Allows builders to validate tokens early and consistently.

### Generic `Errorable`
- `Errorable` is now generic: `Errorable[T any]`.  
- `SetError(err error)` returns the concrete type `T` for safe chaining.  
- Adopted across all tokens.

### BaseToken cleanup
- `BaseToken` now embeds `Validable`.  
- Keeps identity (input, expression, alias) and validation separate.  

---

## 🔧 Supporting Modules

### Token (resolver)
- Centralizes type validation and expression resolution.  
- `ValidateType` rules:
  - `string` → accepted.
  - Existing tokens (`Validable`) → rejected with **Clone()** hint.
  - All other types → `invalid format (type …)`.
- `ResolveExpr` enhancements:
  - Subquery detection: `(SELECT …)` treated as one expression.
  - Strict identifier validation (must be a single token).
  - Alias parsing via inline, `AS`, or trailing identifier.

### Token (ExpressionKind)
- Added `Invalid` kind for unrecognized expressions.  
- `String()` and `IsValid()` updated accordingly.  
- Expression classification improved:
  - Aggregates (`COUNT`, `SUM`, `AVG`, `MIN`, `MAX`) now reported as `Aggregate`.
  - Computed expressions (`price * quantity`) reported as `Computed`.
  - Functions (`JSON_EXTRACT(...)`) remain `Function`.

### Token (helpers)
- Introduced **helpers** package for reusable validation utilities.  
  - Identifier validation:
    - `IsValidIdentifier` / `ValidateIdentifier` with strict ASCII-only rules.
    - Non-ASCII identifiers (e.g. `café`, `mañana`, `niño`) explicitly rejected until dialect-specific rules are added.  
  - Alias validation:
    - `IsValidAlias` / `ValidateAlias` ensure aliases are valid identifiers and not reserved keywords.  
    - `ValidateTrailingAlias` / `HasTrailingAlias` detect and validate trailing aliases (when no `AS` is present).  
    - `ReservedKeywords()` returns the dialect-agnostic set of disallowed aliases.  
  - Deterministic alias generation:
    - `GenerateAlias(prefix, expr string)` produces safe aliases by combining a short code (e.g. `fn`, `sq`, `cp`) with a SHA-1 hash of the expression.  
  - Independent test files with exhaustive valid/invalid cases and runnable examples.  
  - Includes `doc.go` and `README.md` documenting rules and the consistency pattern (`ValidateXxx`, `IsValidXxx`, `GenerateAlias`).

---

## 🔗 New Join Token

A dedicated **`join.Token`** was added to represent SQL JOIN clauses:

- **Safe constructors**: `NewInner`, `NewLeft`, `NewRight`, `NewFull`.  
- **Flexible constructor**: `New(kind any, left, right, condition)` for DSLs/config use.  
- **JoinKind enum**: `InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin` with helpers:
  - `String()` → canonical SQL keyword or `invalid join type (n)`.
  - `IsValid()` → structural validation.
  - `ParseJoinKindFrom()` → case-insensitive parsing.  
- **Validation**:
  - Invalid kind rejected early.
  - Left/Right tables must be valid.
  - Join condition must not be empty.  
- **Contracts**: Implements all shared contracts (`Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, `Stringable`, `Validable`).

---

## 🛠️ Affected Tokens & Builders
- **`table.Token`** and **`field.Token`** now use `resolver.ValidateType`.  
- Invalid states improved:
  - Passing tokens directly → rejected with **Clone()** hint.
  - Invalid alias (including reserved words) → rejected.  
  - Literals/aggregates rejected as table sources.  
- **Impact**: Builders such as `SelectBuilder` now automatically benefit from strict validation and error reporting.

---

## 📚 Documentation & Examples
- `doc.go` updated to mention **resolver**, **ExpressionKind**, **helpers**, and **join**.  
- `README.md`:
  - Root token README lists `field`, `table`, `join`, `resolver`, `ExpressionKind`, and `helpers`.  
  - `helpers` README updated with identifiers, aliases, trailing alias detection, reserved keywords, and deterministic alias generation.  
  - `table` README updated with stricter alias validation, Clone() guidance, and error cases.  
  - Headings normalized (removed emoji from `# Token`).  
- `example_test.go`:
  - Subquery examples uncommented.  
  - New examples for identifiers, aliases, trailing aliases, and generated aliases.  
  - New examples for invalid input and Clone() hints.  
  - Identifier examples added showing non-ASCII rejection.  
  - `IsRaw` examples updated (currently false, will later derive from `Kind()`).  

---

## ✅ Why this matters
- **Consistency**: All helpers and tokens now share strict validation with a `ValidateXxx` + `IsValidXxx` + `GenerateAlias` pattern.  
- **Safety**: Builders detect invalid tokens earlier (reserved aliases, unsupported types, non-ASCII identifiers, literals in FROM).  
- **Extensibility**: Foundation laid for literal helpers, trailing alias refinements, and dialect-specific rules.  
- **Clarity**: Documentation and examples aligned with real behavior.
