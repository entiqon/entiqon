<p align="center">
  <img src="assets/entiqon_logo.png" alt="Entiqon Logo" width="150"/>
</p>

<p align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/ialopezg/entiqon.svg)](https://pkg.go.dev/github.com/ialopezg/entiqon)
[![Go Report Card](https://goreportcard.com/badge/github.com/ialopezg/entiqon)](https://goreportcard.com/report/github.com/ialopezg/entiqon)
[![License](https://img.shields.io/github/license/ialopezg/entiqon)](https://github.com/ialopezg/entiqon/blob/main/LICENSE)
[![Build](https://github.com/ialopezg/entiqon/actions/workflows/test-and-coverage.yml/badge.svg)](https://github.com/ialopezg/entiqon/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ialopezg/entiqon/branch/main/graph/badge.svg)](https://codecov.io/gh/ialopezg/entiqon)
[![Latest Release](https://img.shields.io/github/v/release/ialopezg/entiqon)](https://github.com/ialopezg/entiqon/releases)

</p>

# Entiqon Library

> ⚙️ A structured, intelligent foundation for building queryable, entity-aware Go systems.

---

## 🌱 Overview

Entiqon is a modular query engine designed for extensible data modeling, fluent query building, and structured execution.

---

## ✅ Supported Builders

- `SelectBuilder` with condition chaining, aliasing, ordering, pagination
- `InsertBuilder` with multi-row insert, `RETURNING` support
- `UpdateBuilder` with strict value assignment and column validation

---

## 🚀 Quick Start

### ↘️ Installation

```bash
go get github.com/ialopezg/entiqon
```

---

### 🔍 Select Example

```go
package main

import (
	"fmt"
	"github.com/ialopezg/entiqon/builder"
)

func main() {
	sql, args, err := builder.NewSelect().
		Select("id", "name", "email AS contact").
		From("users").
		Where("status = ?", "active").
		OrderBy("created_at DESC").
		Take(10).
		Skip(5).
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)
}
```

---

### 🧾 Insert Example

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
		Values(1, "Watson").
		Returning("id").
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)
}
```

---

### 🛠 Update Example

> ❗️ Column aliasing is **not allowed** in `UPDATE` queries and will be rejected at build time.

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

	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)
}
```

---

## 📄 License

[MIT](LICENSE) — © Isidro Lopez / Entiqon Project
