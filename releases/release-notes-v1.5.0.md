# 🚀 Release: v1.5.0 – Unified Builder Strategy and Dialect Exposure

**Codename:** Without Sin  
**Tagline:** `Aliasing Is a Sin. And We Are Now Without Sin.`  
**Date:** 2025-05-24

---

## ✨ Highlights

- **🔓 Dialect API Exposed**
  - All dialects are now part of the public driver interface.
  - Enables custom dialects and safe injection of placeholder logic.

- **🏗️ Builder Normalization**
  - SELECT, INSERT, UPDATE, DELETE, and UPSERT unified under consistent tokenization and validation rules.
  - `ParamBinder` handles placeholders per dialect.
  - `StageToken` tags query sections for better diagnostics.

- **🧠 Error Context Awareness**
  - Query-building errors are now tagged using `StageToken` (e.g., `SELECT`, `WHERE`, `INTO`) for traceable validation.

---

## ✅ Coverage & Tests

| Package                             | Coverage |
|-------------------------------------|----------|
| `builder`                           | 94.4%    |
| `driver`                            | 100.0%   |
| `core/builder`                      | 95.5%    |
| `core/builder/bind`                 | 100.0%   |
| `core/error`                        | 85.2%    |
| `core/token`                        | 78.5%    |
| `test`                              | 100.0%   |

---

## 📘 Guides Updated

- SelectBuilder
- InsertBuilder
- UpdateBuilder
- DeleteBuilder
- UpsertBuilder

All reflect the new dialect strategy and validation system via StageToken.

---

## 📄 Resources

- [CHANGELOG](./CHANGELOG.md)
- [Builder Guide](./docs/builder/builder_guide_updates.md)

---

**Onward, builder. To the next doctrine. To the next legend.**  
Use the version. Use the codename.  
**`Aliasing Is a Sin. And We Are Now Without Sin.`**
