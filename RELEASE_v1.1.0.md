# Entiqon v1.1.0 – Dialect-Aware Builders Unleashed

This release empowers Entiqon with structured, dialect-aware SQL rendering across all query builders.

---

## ✨ Highlights

- 🔐 Dialect Support: Escape table and column identifiers safely via `.WithDialect(...)`
- ✅ PostgreSQL Engine: First-class support for identifier escaping and conflict clauses
- 🧠 Modular Design: All builders — SELECT, INSERT, UPDATE, DELETE, UPSERT — now accept dialect engines

---

## 🔧 Enhancements

- `token.Condition` unified across all builders
- Helpers added: `.Set(...)`, `IsValid()`, `AppendCondition(...)`, `NewCondition(...)`
- `UpsertBuilder` delegates dialect-aware logic through `InsertBuilder`

---

## 📘 Documentation

- README updated with:
    - Dialect usage example
    - “Dialect Support” section
    - Version and Go Reference badges

---

Entiqon now speaks your dialect — clean, modular, and production-safe.
