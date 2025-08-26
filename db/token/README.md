<h1 align="left">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true" align="left" height="82" width="82" alt="entiqon"> token
</h1>
<h6 align="left">Part of the <a href="../../README.md">Entiqon</a> / <a href="../README.md">Database</a> toolkit.</h6>

---

## 🎯 Purpose

The **token** module provides low-level primitives that represent SQL
fragments in a dialect-agnostic way.  
These tokens are consumed by higher-level builders (e.g. `SelectBuilder`)
to assemble safe, expressive, and auditable SQL statements.

### Doctrine

- **Never panic** — constructors always return a non-nil token, even if errored.
- **Auditability** — preserve original input for logs and debugging.
- **Strict validation** — invalid inputs are rejected immediately with explicit errors.
- **Delegation** — parsing rules live inside tokens, not in builders.
- **Clarity** — contracts split responsibilities (identity, error, clone, render).

---

## 📜 Contracts

All tokens implement a shared set of contracts:

- **BaseToken** — input, expression, alias, validity
- **Errorable** — explicit error state, never panic
- **Clonable** — safe duplication with preserved state
- **Rawable** — SQL-generic rendering (expr, alias, owner)
- **Renderable** — dialect-agnostic `String()` output
- **Stringable** — concise diagnostic/logging string
- **Ownerable** — ownership binding (`HasOwner`, `Owner`, `SetOwner`)

Each contract has its own test suite to ensure isolation and strict coverage.

---

## 📦 Sub-packages

| Package            | Purpose                                                                                               |
|--------------------|-------------------------------------------------------------------------------------------------------|
| [`field`](./field) | Represents a column, identifier, or computed expression (with aliasing, raw expressions, validation). |
| [`table`](./table) | Represents a SQL source (table or view) used in `FROM` / `JOIN` clauses with aliasing and validation. |

---

## 🚧 Roadmap

Planned tokens:
- **conditions** (WHERE / HAVING)
- **joins** (INNER, LEFT, etc.)
- **functions** (aggregates, JSON, custom expressions)

Contracts will progressively enforce stricter auditability across all tokens.

---

## 📄 License

[MIT](../../LICENSE) — © Entiqon Project
