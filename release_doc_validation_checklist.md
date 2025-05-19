
# 📦 Entiqon Documentation Release Validation Checklist

This checklist ensures a safe, consistent, and dependency-aware rollout of builder and dialect documentation for Entiqon v1.2.0+.

---

## ✅ Phase 1: Builder Docs Audit

Verify each builder:
- [ ] Uses `QuoteIdentifier(...)` (not `Escape(...)`)
- [ ] Does **not** redefine or override dialect behavior
- [ ] Contains a reference or note pointing to the Dialect Guide
- [ ] Has ≥100% test coverage, including validation failures
- [ ] Has clear GoDoc comments for every method
- [ ] Is listed in `README.md` under `Developer Guides`

---

## ✅ Phase 2: Dialect Guide Review

Ensure the Dialect Guide:
- [ ] Fully documents the `Dialect` interface
- [ ] Clarifies the replacement of `Escape(...)` with `QuoteIdentifier(...)` and `QuoteLiteral(...)`
- [ ] Warns that `QuoteLiteral(...)` is **not** safe for production query building
- [ ] Explains builder quoting expectations (field, table, column)
- [ ] Provides implementation examples (`PostgresDialect`)
- [ ] Links back to builder guides for usage context

---

## ✅ Phase 3: Readme Index Validation

Ensure the root `README.md`:
- [ ] Has a `Developer Guides` section with all builder and dialect links
- [ ] Reflects correct paths under `docs/developer/builder/*` and `docs/developer/architecture/*`
- [ ] Has a `Principles & Best Practices` section
- [ ] Notes test coverage, naming clarity, and deprecation handling

---

## ✅ Phase 4: Commit & Tag Preparation

- [ ] Remove any deprecated files (`upsert.md`, `upsert_builder.md`, old dialect notes)
- [ ] Commit updated builder docs and dialect guide
- [ ] Confirm Markdown renders correctly in GitHub/GitHub Pages
- [ ] Tag release (`v1.2.x` or `v1.3.0`) if no breaking changes are present

---

## ✅ Notes

- `WithDialect(...)` is covered and deprecated — removal planned for `v1.4.0`
- All builders are now dialect-safe by contract via `QuoteIdentifier(...)`
- Dialect authors must implement all six methods to be compliant
