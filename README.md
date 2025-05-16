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
- `UpdateBuilder` with SET + WHERE and param binding

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

```go
package main

import (
	"fmt"

	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, err := builder.NewSelect().
		Select("id", "name").
		From("users").
		Where("status = 'active'").
		AndWhere("created_at > '2023-01-01'").
		OrderBy("created_at DESC").
		Take(10).
		Skip(5).
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(sql)
	// Output:
	// SELECT id, name FROM users WHERE status = 'active' AND created_at > '2023-01-01' ORDER BY created_at DESC LIMIT 10 OFFSET 5
}
```

---

#### ✍️ Usage Example (INSERT)

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
	// Output:
	// INSERT INTO users (id, name) VALUES (?, ?) RETURNING id
	// [1 Sherlock]
}
```

---

#### 🔄 Usage Example (UPDATE)

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
	// Output:
	// UPDATE users SET status = ? WHERE id = ?
	// [active 42]
}
```

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
