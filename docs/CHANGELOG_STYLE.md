# Changelog Style Guide

This project uses a **two-tier documentation system** for changes:

- **CHANGELOG.md** → concise, user-facing overview of *what* changed.
- **Release Notes** (per version) → detailed explanation of *why* and *how* changes were made.

---

## Format

We follow [Keep a Changelog](https://keepachangelog.com/) and [Semantic Versioning](https://semver.org/).

Each release section should use the following headings (include only if relevant):

```markdown
## vX.Y.Z - YYYY-MM-DD

### Added
- New features or modules (e.g. tokens, helpers, contracts).

### Changed
- Modifications or enhancements (validation stricter, API refinements).

### Fixed
- Bug fixes (panics resolved, alias parsing issues).

### Removed
- Deprecated or removed features.

### Docs & Tests
- Documentation updates and test changes.
```

---

## Rules

- **Concise** → one bullet point per change, one line max.
- **Consistent** → always use the same headings and tone.
- **User-facing** → describe what changed, not internal refactoring details.
- **Detailed context** belongs in release notes, not the CHANGELOG.

---

## Example

```markdown
## v1.14.0 - 2025-08-29

### Added
- `join.Token` with `JoinKind` enum and safe/flexible constructors.
- `resolver` module with `ValidateType` and `ResolveExpr`.
- `Invalid` kind in `ExpressionKind`.
- `helpers` package with `IsValidIdentifier`.

### Changed
- Extended classification: aggregates, computed expressions, functions.
- Stricter alias validation in `table.Token` and `field.Token`.

### Docs & Tests
- Updated docs and examples for resolver, ExpressionKind, join, and helpers.
```
