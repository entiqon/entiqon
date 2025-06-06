# Token System Guide

This guide explains how internal tokens like `token` are structured, parsed, validated, and consumed by query builders.

> All tokens (e.g., `token`, `Table`, `Condition`) embed `BaseToken`, making them both **`Errorable`** (via
> `GetError/IsErrored/SetError`) and **`Kindable`** (via `GetKind/SetKind`). All these methods are nil‐safe: calling any
> of them on a `nil` token will not panic and returns a safe default (e.g., `IsErrored()` → `false`, `GetKind()` →
> `UnknownKind`).

---

## 📜 Principles

- Tokens are internal, immutable, and self‐validating.
- Each token may carry an `Error` (via the **Errorable** contract).
    - Calling `IsErrored()` on a token with no error (or on a `nil` token) returns `false`.
    - Calling `GetError()` on a token with no error (or on a `nil` token) returns `nil`.
- Tokens also expose a `Kind` (via the **Kindable** contract).
    - Calling `GetKind()` on a token with no kind set (or on a `nil` token) returns `UnknownKind`.
    - Calling `SetKind(...)` on a `nil` token is a no‐op (nil‐safe).
- Tokens do not handle dialect quoting or full SQL rendering beyond their own `RenderName()` / `RenderAlias()`.
- `BaseToken` provides name, alias, error, and kind logic for all tokens, ensuring it is reusable and consistent.

## 🧩 Core Behaviors

- Input normalization
- Alias resolution
- Error tracking
- Dialect-safe rendering

## 🔍 Qualification & Validation

A generic token may become invalid under any of the following scenarios. All tokens embed `BaseToken`, so these apply
universally (token, Table, Condition, etc.):

- **Empty input**: An input string that is blank or only whitespace sets an error.
- **Comma-separated expression**: Inputs containing commas (`,`) are rejected (aliases must not be comma-separated).
- **Starts with `AS`**: Inputs beginning with `AS ` or whose base resolves to `AS` only are invalid (missing identifier
  before `AS`).
- **Alias conflict**: When both inline and explicit aliases are provided but do not match (e.g.,
  `"users.id AS uid", "other_uid"`), an error is set.
- **Qualifier conflict**: For tokens that parse qualifiers (e.g., `users.id`), if the qualifier does not match an
  attached context (e.g., a mismatched table alias), the token becomes invalid.
- **Reserved word misuse**: If the parsed base identifier itself is a reserved word (e.g., `"AS"` only), the token is
  invalid.

All of the above conditions will cause `IsErrored()` to return `true`, with details accessible via `GetError()`. You can
check `IsValid()` (equivalent to `!IsErrored() && Name != ""`) before including a token in any generated SQL.

---

## 🧱 Related Tokens

The following core token is the base for all SQL tokens in Entiqon:

|                         | Name          | Description                                                                                                          | Access    |
|-------------------------|:--------------|----------------------------------------------------------------------------------------------------------------------|-----------|
| **[📄](base_token.md)** | **BaseToken** | Is the abstract foundational structure used by all SQL tokens in Entiqon (such as `token`, `Table`, and `Condition`) | `Private` |

This modular reference allows isolated testing and future reuse for other token types.

---
2025 — **© Entiqon Project**
