<p align="center">
  <img src="assets/entiqon_logo.png" alt="Entiqon Logo" width="150"/>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg)](https://pkg.go.dev/github.com/ialopezg/entiqon)

# Entiqon Library

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

## 🌱 Overview

Entiqon is a modular query engine designed for extensible data modeling, fluent query building, and structured execution.

---

## ✅ Supported Builders

- `SelectBuilder` with condition chaining, pagination
- `InsertBuilder` with multi-row insert, `RETURNING` support

---

## 🚀 Quick Start

```bash
go get github.com/ialopezg/entiqon

```go
package your-package

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

## ✍️ Usage Example (INSERT)

```go
sql, args, err := builder.NewInsert().
  Into("users").
  Columns("id", "name").
  Values(1, "Sherlock").
  Returning("id").
  Build()

// Result:
// INSERT INTO users (id, name) VALUES (?, ?) RETURNING id
```

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
