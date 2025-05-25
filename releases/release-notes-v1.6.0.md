# Release v1.6.0 — Aliasable SQL expressions with table support

## ✨ Features

- ✅ **Aliasable Expressions**
  - SQL elements can now carry aliases and optional table prefixes
  - Supports parsing input like: `table.name AS alias`
  - Foundation for consistent token formatting and rendering

- ✅ Reusable Abstractions
  - Built on top of the `AliasableToken` core
  - Includes parsing helpers for alias and `table.column` resolution

## 🧪 Tests

- Covers common usage formats, table overrides, and malformed input handling

## 📅 Release Date

2025-05-25

**Codename:** Keystone
