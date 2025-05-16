<p align="center">
  <img src="assets/entiqon_logo.png" alt="Entiqon Logo" width="150"/>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg)](https://pkg.go.dev/github.com/ialopezg/entiqon)

# Entiqon Library

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

## 🌱 Overview

This project is built with [Entiqon](https://github.com/ialopezg/entiqon) — a modular, dialect-aware query engine designed for extensible data modeling, intelligent query building, and systemic evolution.

Whether you're modeling entities, resolving relationships, or generating dialect-specific queries, this library gives you structure, power, and growth potential.

---

## 📦 Structure

```
.
├── builder/       # Query builders (select, insert, etc.)
├── dialect/       # Dialect interfaces & engines (Postgres, MySQL)
├── entity/        # Metadata, resolution logic
├── docs/          # Markdown documentation
├── examples/      # Sample usage
├── assets/        # Logo, branding
```

---

## 🚀 Quick Start

```bash
go get github.com/ialopezg/entiqon
```

```go
import (
  "github.com/ialopezg/entiqon/builder"
  "github.com/ialopezg/entiqon/dialect"
)

pg := dialect.PostgresDialect{}
qb := builder.NewQuery(pg)

qb.Select().
Columns("id", "email").
From("users").
WhereNamed("email = :email", map[string]any{"email": "user@example.com"})

sql, args, err := qb.Build()

```

---

## 💡 Powered by Entiqon

Entiqon is more than a query builder — it’s a **structured intelligence layer**.  
It models entities, resolves relationships, and evolves with your data.

> 🤖 *Not a brain. A system. Supercharged.*

---

## 🧪 Run Tests

```bash
go test ./...
```

---

## ✅ Supported Query Builders

- `SelectBuilder` with fluent chaining
- Logical condition grouping via `Where`, `AndWhere`, `OrWhere`
- Pagination using `Take` and `Skip`
- Clause-safe `Build()` output

## 🚀 Usage Example (SELECT)

```go
sql, err := builder.NewSelect().
  Select("id", "name").
  From("users").
  Where("status = 'active'").
  AndWhere("created_at > '2023-01-01'").
  OrderBy("created_at DESC").
  Take(10).
  Skip(5).
  Build()

// Result:
// SELECT id, name FROM users WHERE status = 'active' AND created_at > '2023-01-01' ORDER BY created_at DESC LIMIT 10 OFFSET 5
```

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
