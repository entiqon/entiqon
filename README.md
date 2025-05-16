<p align="center">
  <img src="assets/entiqon_logo.png" alt="Entiqon Logo" width="150"/>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg)](https://pkg.go.dev/github.com/ialopezg/entiqon)

# Entiqon Library

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

## 🌱 Overview

Entiqon is a modular query engine designed for extensible data modeling, fluent query building, and structured
execution.

---

## ✅ Supported Builders

- `SelectBuilder` with condition chaining, pagination
- `InsertBuilder` with multi-row insert, `RETURNING` support
- `UpdateBuilder` with SET + WHERE and param binding
- `UpsertBuilder` with `ON CONFLICT ... DO UPDATE SET ...` support
- `DeleteBuilder` with WHERE clause and optional RETURNING support

---

## 🚀 Quick Start

---

### ↘️ Installation

```bash
go get github.com/ialopezg/entiqon
```

---

### 📘 Usage

---

#### 🚀 Usage Example (SELECT)

A SELECT query that retrieves user emails filtered by status, role, and signup date — with ordering, pagination, and parameter binding.

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewSelect().
		Select("id", "email").
		From("users").
		Where("status = ?", "active").
		AndWhere("role = ?", "admin").
		AndWhere("created_at > ? AND region = ?", "2024-01-01", "NA").
		OrderBy("last_login DESC").
		Take(50).
		Skip(0).
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
```

---

#### ✍️ Usage Example (INSERT)

Inserts a new user record and returns the inserted ID using PostgreSQL's RETURNING clause.

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewInsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Sherlock").
		Returning("id").
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
```

---

#### 🔄 Usage Example (UPDATE)

Updates a user's status using parameterized WHERE and SET clauses.

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewUpdate().
		Table("users").
		Set("status", "active").
		Where("id = ?", 42).
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
```

---

#### ♻️ Usage Example (UPSERT)

Performs an UPSERT — inserts or updates an existing user if a conflict on ID occurs.

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewUpsert().
		Into("users").
		Columns("id", "name").
		Values(1, "Watson").
		OnConflict("id").
		DoUpdateSet(map[string]string{
			"name": "EXCLUDED.name",
		}).
		Returning("id").
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
```

---

#### 🗑️ Usage Example (DELETE)

Deletes a user by ID and returns the deleted ID — supports PostgreSQL's RETURNING clause.

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewDelete().
		From("users").
		Where("id = ?", 42).
		Returning("id").
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	fmt.Println(args)
}
```

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
