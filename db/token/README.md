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

Key principles:
- **Immutability** — tokens are never mutated after construction; cloning is explicit.
- **Auditability** — identity aspects (input, expression, alias, owner, validation) are separated into contracts.
- **Consistency** — all tokens share common contracts like `BaseToken`, `Renderable`, `Errorable`.

---

## 📦 Sub-packages

| Package            | Purpose                                                                                               |
|--------------------|-------------------------------------------------------------------------------------------------------|
| [`field`](./field) | Represents a column or expression in a `SELECT` clause (with aliasing, raw expressions, validation).  |
| [`table`](./table) | Represents a SQL source (table or view) used in `FROM` / `JOIN` clauses with aliasing and validation. |

---

## 🚧 Roadmap

Future tokens will include:
- **conditions** (WHERE / HAVING)
- **joins** (INNER, LEFT, etc.)
- **functions** (aggregates, JSON, custom expressions)

Contracts will progressively enforce stricter auditability across all tokens.

---

## 📄 License

[MIT](../../LICENSE) — © Isidro Lopez / Entiqon Project