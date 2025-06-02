# 🧩 Creating a Custom SQL Dialect with BaseDriver

If your SQL engine is not natively supported by Entiqon, you can create your own dialect by embedding the `BaseDriver` and customizing only what you need.

---

## ✅ Step 1: Import the `core/driver` package

```go
import "github.com/entiqon/entiqon/core/driver"
```

---

## ✅ Step 2: Define your custom driver

```go
type MyAwesomeDriver struct {
	driver.BaseDriver
}

func NewMyAwesomeDriver() *MyAwesomeDriver {
	return &MyAwesomeDriver{
		BaseDriver: driver.BaseDriver{
			Name:      "awesome",
			Quotation: driver.QuoteDouble,
			Placeholder: func(i int) string {
				return fmt.Sprintf("#%d", i) // or whatever your engine uses
			},
			SupportsUpsert:    true,
			SupportsReturning: false,
		},
	}
}
```

---

## 🛠 Available Quote Styles (`driver.QuoteType`)

| Value           | Output        | Use for         |
|-----------------|---------------|-----------------|
| `QuoteNone`     | `users`       | Generic engines |
| `QuoteDouble`   | `"users"`     | PostgreSQL      |
| `QuoteBacktick` | `` `users` `` | MySQL, SQLite   |
| `QuoteBracket`  | `[users]`     | SQL Server      |

---

## 🧠 Placeholder Strategy

Define placeholders using a function:

```go
Placeholder: func(i int) string {
	return fmt.Sprintf("@p%d", i) // or "?" for generics
}
```

---

## ✅ Final Integration

```go
b := builder.NewSelectBuilder().UseDialect("awesome")
```

To register it globally:

```go
driver.RegisterDialect("awesome", NewMyAwesomeDriver())
```

---

## 🔒 BaseDriver Responsibilities

`BaseDriver` already provides:

- Quoting via `QuoteIdentifier`
- Placeholder resolution via `Placeholder()`
- LIMIT/OFFSET via `BuildLimitOffset()`
- FROM clause rendering via `RenderFrom()`

---

## ✅ Examples

- `core/driver/postgres.go`
- `core/driver/mysql.go`
- `core/driver/mssql.go`

---

## 📣 Need Help?

Contributions and questions welcome via GitHub Issues or Discussions!