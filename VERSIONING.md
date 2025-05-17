# Entiqon Versioning Strategy

Entiqon follows [Semantic Versioning 2.0.0](https://semver.org) to communicate changes, stability, and compatibility.

---

## 🧭 Version Format

```
MAJOR.MINOR.PATCH
```

Each version number reflects the type of change introduced.

---

## ✅ PATCH (v1.0.X)

Patch releases are for:
- Bug fixes (e.g., param resolution, condition logic)
- Internal refactors without changing the public API
- Documentation improvements
- Test suite adjustments

> **Safe to upgrade without breakage.**

---

## ✅ MINOR (v1.X.0)

Minor versions add **backward-compatible** features, such as:
- New builders (e.g., JoinBuilder, MergeBuilder)
- New methods (e.g., Having, Limit, Skip support)
- Promotion of internal utilities (e.g., ConditionToken reuse across all builders)
- Flexible parsing or dialect additions

> **Safe to upgrade. Feature-enhancing.**

---

## ✅ MAJOR (vX.0.0)

Major versions introduce **breaking changes**, including:
- Method signature changes
- Removal or renaming of builders
- Return structure changes
- Core design adjustments that require user adoption

> **Upgrade only after reviewing changelogs and migration notes.**

---

## 🔄 Pre-Release Tags

Entiqon supports pre-release phases for feature previews:

```
v1.2.0-alpha.1
v2.0.0-beta.2
v2.0.0-rc.1
```

Use these for experimentation and early validation. These versions do not appear as latest stable on `pkg.go.dev`.

---

## 🧘‍♂️ Discipline

Every change to Entiqon adheres to this protocol:

1. Implementation
2. Review & Audit
3. Testing
4. Documentation (incl. README, GoDoc)
5. CHANGELOG Update
6. Tagging & Release

---

> Entiqon stands in perfect order — and so must its versioning.

