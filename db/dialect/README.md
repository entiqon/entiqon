# 🧩 Dialect Package

The **Dialect package** defines the **contract (`SQLDialect`)** and the **capability matrix (`Options`)** 
that all SQL dialects in Entiqon must implement.  
A dialect encapsulates the vendor-specific rules needed to render portable, correct SQL.

---

## 📜 Core Types

### `SQLDialect`

The shared interface implemented by every dialect:

```go
type SQLDialect interface {
    Name() string
    Options() Options
    QuoteIdentifier(name string) string
    QuoteLiteral(literal any) string
    PaginationSyntax(limit, offset int) string
    Placeholder(index int) string
}
```

### `Options`

A struct advertising the capabilities of a dialect:

```go
type Options struct {
    Name                  string
    QuoteStyle            string
    PlaceholderStyle      string
    AllowMerge            bool
    AllowUpsert           bool
    ForcedAliasing        bool
    EnableReturning       bool
    SupportsCTE           bool
    SupportsWindowFunctions bool
    MaxPlaceholderIndex   int
}
```

---

## 📂 Dialects

| Dialect       | Status       | Description                                                                 |
|---------------|--------------|-----------------------------------------------------------------------------|
| [`generic`](./generic)   | ✅ Implemented | ANSI-compliant fallback, safe default                                 |
| [`postgres`](./postgres) | 🚧 Planned     | PostgreSQL-specific rules (RETURNING, `$` placeholders)               |
| [`mysql`](./mysql)       | 🚧 Planned     | MySQL rules (backtick quoting, LIMIT syntax)                          |
| [`mariadb`](./mariadb)   | 🚧 Planned     | MariaDB rules, mostly MySQL-compatible with some extensions           |
| [`sqlite`](./sqlite)     | 🚧 Planned     | SQLite rules (dynamic typing, `?` placeholders, `LIMIT`/`OFFSET`)     |
| [`mssql`](./mssql)       | 🚧 Planned     | Microsoft SQL Server rules (`[bracket]` quoting, `TOP`, `OFFSET FETCH`)|
| [`oracle`](./oracle)     | 🚧 Planned     | Oracle rules (`:v1` placeholders, `ROWNUM`, `RETURNING INTO`)         |
| [`db2`](./db2)           | 🚧 Planned     | IBM DB2 rules (positional `?`, common table expressions, MERGE)       |
| [`firebird`](./firebird) | 🚧 Planned     | Firebird rules (`FIRST`/`SKIP` instead of LIMIT/OFFSET)               |
| [`informix`](./informix) | 🚧 Planned     | Informix rules (first/skip, specific syntax quirks)                   |
| [`cockroach`](./cockroach) | 🚧 Planned   | CockroachDB (Postgres-compatible with distributed SQL extensions)     |
| [`tidb`](./tidb)         | 🚧 Planned     | TiDB (MySQL-compatible with clustering features)                      |
| [`hana`](./hana)         | 🚧 Planned     | SAP HANA SQL (specialized functions, column store quirks)             |
| [`snowflake`](./snowflake) | 🚧 Planned   | Snowflake cloud SQL dialect (semi-structured data, `QUALIFY`)         |
| [`redshift`](./redshift) | 🚧 Planned     | Amazon Redshift (Postgres-like but missing some features)             |
| [`teradata`](./teradata) | 🚧 Planned     | Teradata SQL dialect (large-scale analytics focus)                    |
| [`clickhouse`](./clickhouse) | 🚧 Planned | ClickHouse (analytics engine with SQL-like syntax, not ANSI complete) |

---

## 🚀 Usage

```go
import (
    "fmt"
    "entiqon/db/dialect"
    "entiqon/db/dialect/generic"
)

func main() {
    var d dialect.SQLDialect = generic.New()

    fmt.Println(d.Name())          // "generic"
    fmt.Println(d.Placeholder(1))  // "?"
    fmt.Println(d.PaginationSyntax(10, 5)) 
    // → " LIMIT 10 OFFSET 5"

    opts := d.Options()
    fmt.Printf("SupportsCTE=%v Returning=%v\n", opts.SupportsCTE, opts.EnableReturning)
}
```

---

## 🧭 Extending

To add a new dialect:

1. Create a subpackage (`dialect/postgres`, `dialect/oracle`, etc.).  
2. Implement `SQLDialect`, initializing with proper `Options`.  
3. Override methods like `Placeholder`, `QuoteIdentifier`, or `PaginationSyntax` if the vendor differs from ANSI.  
