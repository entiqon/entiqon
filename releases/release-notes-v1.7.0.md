# Release v1.7.0 - Strengthened Token Foundation

**Codename:** Forge  
**Tagline:** `Strengthened Token Foundation`  
**Date:** 2025-06-09

---

## Overview

This release marks the completion of **Task 1: Column Injection Logic** from the Global Normalization Plan. It finalizes
the core token abstraction (`BaseToken`), adds strong typing and error-handling contracts, and brings our test coverage,
documentation, and code consistency to 100%.

## ✨ Highlights

- **Column Injection Logic Complete**
    - Orphan, Qualified, Qualified+Aliased, and ambiguous qualifiers handled by `NewColumn(...)` + `BaseToken`
    - Alias parsing & conflict detection fully centralized in `BaseToken`

- **Errorable Contract**
    - Introduced `Errorable` interface (`GetError()`, `IsErrored()`, `SetError()`) in
      `internal/core/contract/errorable.go`
    - `BaseToken` implements `Errorable` with nil-safe methods
    - Deprecated `SetErrorWith()` and `HasError()` remain for backward compatibility

- **Kindable Contract**
    - Defined `Kind` enum in `internal/core/contract/kind.go`
    - Introduced `Kindable` interface (`GetKind()`, `SetKind()`) in `internal/core/contract/kindable.go`
    - `BaseToken` implements `Kindable`; `String()` now shows token’s kind

- **BaseToken Refactor & API Enhancements**
    - Renamed `HasError()`→`IsErrored()`, `SetErrorWith()`→`SetError()`
    - Added `GetError()` for direct error inspection
    - Added `SetKind()`/`GetKind()` for type classification
    - Nil-safe everywhere: calling any method on a `nil` token returns safe defaults

- **Documentation & Guides**
    - Created `base_token.md`—complete developer guide for `BaseToken`
    - Updated `token.md` to reference `Errorable` & `Kindable`, nil-safe behavior, and generic token cases
    - Enhanced readability with icons and refined headings

- **100% Test Coverage**
    - All new methods and branches covered by unit tests
    - Defensive tests for nil receivers, alias conflicts, error propagation, and kind classification

## Breaking Changes

- **Interface Renames**:
    - `GenericToken` now uses `IsErrored()` instead of `HasError()`
    - `Kinded` renamed to `Kindable`

- **Package Path Changes**:
    - Contract interfaces moved to `internal/core/contract`
