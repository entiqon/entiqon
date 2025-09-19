# 🗄️ Generic Dialect

> Part of [Entiqon](../../../) / [Database](../../) / [Dialect](../)

The **Generic Dialect** is the **ANSI SQL–compliant baseline** implementation for the Entiqon SQL builder.  
It provides a safe default when generating SQL queries against unknown or mixed databases.

---

## ✨ Features

- **ANSI-compliant quoting**  
  - Identifiers quoted with double quotes (`"users"`, `"UserData"`)  
  - Lowercase simple names left unquoted (`users`)  

- **Placeholders**  
  - Always `?` (non-positional, ANSI style)  

- **Pagination**  
  - Standard `LIMIT` and `OFFSET` syntax  

- **Literal quoting**  
  - Strings: escaped and single-quoted (`'O''Reilly'`)  
  - Booleans: `TRUE` / `FALSE`  
  - Numbers: plain decimal form  
  - `time.Time`: UTC `'YYYY-MM-DD HH:MM:SS'`  
  - `nil`: `NULL`  

- **Capabilities**  
  - ✅ CTE (`WITH ...`)  
  - ✅ Window functions (`OVER(...)`)  
  - ❌ RETURNING not supported  
  - ❌ MERGE not supported  
  - ❌ UPSERT not supported  

---

## 🚀 Usage

```go
import (
    "fmt"
    "entiqon/db/dialect/generic"
)

func main() {
    d := generic.New()

    // Quote identifiers
    fmt.Println(d.QuoteIdentifier("users"))     // → users
    fmt.Println(d.QuoteIdentifier("UserData"))  // → "UserData"

    // Quote literals
    fmt.Println(d.QuoteLiteral("O'Reilly"))     // → 'O''Reilly'
    fmt.Println(d.QuoteLiteral(true))           // → TRUE

    // Pagination
    sql := fmt.Sprintf("SELECT %s FROM %s%s",
        d.QuoteIdentifier("id"),
        d.QuoteIdentifier("users"),
        d.PaginationSyntax(10, 0),
    )
    fmt.Println(sql)
    // → SELECT "id" FROM "users" LIMIT 10
}
```

---

## 🔍 When to Use

- As a **fallback** when the target database dialect is unknown.  
- For **unit tests** to validate SQL builders in a dialect-agnostic way.  
- As a **template** when creating a new vendor-specific dialect (e.g. Postgres, MySQL).  

---

## 📂 Related

- [`dialect.Options`](../options.go) — shared capability matrix.  
- [`dialect.SQLDialect`](../options.go) — interface contract implemented by Generic and vendor dialects.  
- [`dialect/postgres`](../postgres) — Postgres-specific dialect.  
- [`dialect/mysql`](../mysql) — MySQL-specific dialect.  
