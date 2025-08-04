# Entiqon v1.8.0 – "Atlas" Release Notes

Release Date: 2025-08-02

---

## Added Changes

- Modularization of the `db` package as a standalone Go module `github.com/entiqon/db`.
- Import paths updated from `github.com/entiqon/entiqon/...` to `github.com/entiqon/db/...`.
- Added initial test coverage and CI integration for the `db` module.
- Preparation for incremental modularization by isolating the `db` module before adding others like `core`, `auth`, and `http`.

---

## Improvements

- Enhanced code maintainability and scalability by adopting Go modules per package.
- Streamlined development workflow with `go.work` workspace support.
- Added `run-tests.sh` script to facilitate testing with optional coverage.

---

## Migration Notes

- Update your imports to use the new module paths:

  ```go
  import "github.com/entiqon/db/builder"
  ```

- Adjust your build and dependency management accordingly to support multiple modules in the repository.
- Refer to the updated documentation and developer guides for usage.

---

Thank you for using Entiqon! This release lays the foundation for a modular and scalable future.

---

© 2025 Entiqon Project — Inspired by Mythology & Legends
