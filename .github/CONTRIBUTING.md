# Contributing to Entiqon

Welcome! 🎉 We’re excited that you’re interested in contributing to Entiqon.  
This guide explains the conventions we follow for commits, tests, documentation, and releases.  
Consistency is key — every contribution should match the project’s style and rigor.

---

## 🧪 Testing

We require **80% test coverage**.  
Follow this structure for all tests:

- **Pattern**: `file → methods → cases`
- **Naming**: Always **PascalCase**, no underscores.
- **Case organization**: Use `t.Run` blocks to express hierarchy clearly.

Example:
```go
t.Run("Methods", func(t *testing.T) {
    t.Run("Having", func(t *testing.T) {
        t.Run("ResetCollection", func(t *testing.T) {
            // ...
        })
    })
})
```

Edge cases are **mandatory** (nil receivers, empty collections, invalid input, etc.).

---

## 📑 Documentation

Every feature must be reflected in:
- `doc.go` → Package-level usage and examples.
- `README.md` → High-level feature documentation with runnable snippets.
- `example_test.go` → Go runnable examples (`ExampleXxx` functions).
- `CHANGELOG.md` → Notable changes, grouped by release.
- `release-notes-vX.Y.Z.md` → Detailed notes for GitHub releases.

If a feature is not in all of these, it’s **not complete**.

---

## 💬 Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) with semantic prefixes.

- **Format**:
  ```
  feat(scope): short description
  ```
- **Body**: Explain *what* and *why*. List sub-bullets for details (API, tests, docs).
- **Style**: Clear, consistent, and with personality ✨

Example (detailed commit):
```text
feat(builder/select): add Having clause support to SelectBuilder

- Introduced Having(), AndHaving(), OrHaving() for HAVING clause
- Integrated into Build() after GROUP BY
- Added tests (NilCollection, Single, Multiple, IgnoreEmpty, ResetCollection, And, Or)
- Updated README.md, doc.go, example_test.go, and CHANGELOG
- Release notes updated to document Having support

✨ Time for Having!
```

### Squash Commits
If squashing multiple commits:
- Keep the subject line short.
- Body = concise summary of changes.

Example (squash):
```text
feat(builder/select): add Having clause support

- New methods: Having, AndHaving, OrHaving
- Integrated into Build() after GROUP BY
- Tests, docs, examples, and release notes updated
```

---

## 🚀 Release Process

1. **Finish features** → Ensure tests/docs/README/example_test are all updated.  
2. **Update docs** → `README.md`, `doc.go`, examples.  
3. **Update CHANGELOG.md** → Add entries under the current “Upcoming” section.  
4. **Update release notes** → `release-notes-vX.Y.Z.md`.  
5. **Tag and publish** → When ready, create a GitHub release with `gh release create`.  

---

## ✅ Summary

- Write tests first (nil → valid → edge).
- Update **all docs** (README, doc.go, examples).
- Follow **commit conventions** (detailed vs squash).
- Keep CHANGELOG and release notes up to date.
- Add some flair: ✅/⛔️ markers, fun phrases like *“Time for Having!”* make the project memorable.
