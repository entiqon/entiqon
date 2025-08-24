# 📝 Contributor Quick Reference

## ✅ Before Commit
- [ ] Tests written and passing (nil → valid → edge cases).  
- [ ] Tests follow `file → methods → cases` pattern with **PascalCase**.  
- [ ] Docs updated:
  - [ ] `doc.go`
  - [ ] `README.md`
  - [ ] `example_test.go`
  - [ ] `CHANGELOG.md`
  - [ ] `release-notes-vX.Y.Z.md`

## 💬 Commit Rules
- Use **Conventional Commits**:
  - `feat(scope): ...` → new feature  
  - `fix(scope): ...` → bug fix  
  - `docs(scope): ...` → docs only  
  - `refactor(scope): ...` → no behavior change  
- Detailed commits → list features, tests, docs.  
- Squash commits → short summary.  

Examples:  
```text
feat(builder/select): add Having clause support

- New methods: Having, AndHaving, OrHaving
- Integrated into Build() after GROUP BY
- Tests, docs, examples, and release notes updated
```

## 🚀 Release Flow
1. Ensure **100% coverage**.  
2. Update **docs + examples**.  
3. Update **CHANGELOG** (under "Upcoming").  
4. Update **release notes** file.  
5. Tag & release with `gh release create`.  

## ✨ Style
- ✅ / ⛔️ markers in errors and docs.  
- Optional: add a fun closing phrase in commits, e.g.:  
  ```
  ✨ Time for Having!
  ```

---

## ✨ Notes for Reviewers

Please describe any additional context, special considerations, or follow-up work here.
