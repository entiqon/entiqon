# 🚀 Release v1.14.0 – Contract Evolution

This release introduces a major refinement of the **contract layer** across `db/token`, ensuring better separation of concerns, type safety, and consistency.  
It also affects **builders** (notably `SelectBuilder`), since they directly depend on `table.Token`.

---

## ✨ Highlights

### New `Validable` contract
- Introduced **`Validable`** with `IsValid()` for structural validation.  
- Allows higher-level builders to quickly determine token validity without depending on `BaseToken`.  
- Adopted by `table.Token` and `field.Token`.  

### Generic `Errorable`
- **`Errorable`** is now generic: `Errorable[T any]`.  
- `SetError(err error)` now returns the concrete token type (`T`) for safe method chaining.  
- Applied consistently across all tokens (`Field`, `Table`, `Join`, …).  

### BaseToken cleanup
- **`BaseToken`** now embeds `Validable` instead of declaring `IsValid()` directly.  
- Keeps identity (input, expression, alias) and validation clearly separated.  

---

## 🛠️ Affected Tokens & Builders
- **`table.Token`** and **`field.Token`** updated to embed `Validable`.  
- Implementations adjusted to satisfy the generic `Errorable[T]`.  
- **Impact**: `table.Token` is consumed directly by `SelectBuilder` as a source.  
  - This means builders automatically gain structural validation through `IsValid()`.  
  - Invalid tables are now caught earlier, reducing runtime ambiguity.  

---

## 📚 Documentation & Examples
- **`doc.go`** rewritten with strict contract ordering:  
  `BaseToken → Clonable → Debuggable → Errorable → Rawable → Renderable → Stringable → Validable`.  
- **`README.md`** updated to reflect the new contract set and generic `Errorable`.  
- **`example_test.go`** revised:  
  - One method → one example.  
  - Valid and invalid cases shown consistently.  
  - `BaseToken` and `Validable` tested independently.  

---

## ✅ Why this matters
These changes make contracts:
- **More composable** → each contract handles a single concern.  
- **More type-safe** → generic `Errorable` avoids unsafe casts.  
- **More auditable** → `Validable` provides a universal way to check validity.  
- **Builder-aware** → `SelectBuilder` (and future builders) benefit automatically from `Validable` checks.  

This lays the foundation for more advanced tokens (`Join`, computed fields, subqueries) while ensuring builders remain safe and predictable.

---

## 🔗 New Join Token

This release also introduces a dedicated **`join.Token`** for representing SQL JOIN clauses:

- **Safe constructors**: `NewInner`, `NewLeft`, `NewRight`, `NewFull`.
- **Flexible constructor**: `New(kind any, left, right, condition)` for advanced scenarios (e.g. configuration, DSLs).
- **Kind enum**: `join.Kind` (`InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin`) with helpers:
  - `String()` → canonical SQL keyword or `invalid join type (n)` for invalid values.
  - `IsValid()` → structural validation.
  - `ParseJoinKindFrom()` → case-insensitive string parsing.
- **Validation rules**:
  - Early exit on invalid kind.
  - Left/Right tables must be present and valid.
  - Condition must not be empty.
- **Contracts**:
  - Implements all shared contracts: `Clonable`, `Debuggable`, `Errorable`, `Rawable`, `Renderable`, `Stringable`, `Validable`.

This ensures JOIN clauses are first-class citizens in the builder ecosystem, consistent with fields and tables.

---
