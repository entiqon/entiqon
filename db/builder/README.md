# 📊 SQL Builder Package

> Part of [Entiqon](../../) / [Database](../)

Currently, provides the `SelectBuilder` for constructing SQL SELECT queries.  
Builders for inserts, updates, deletes, and merge (upsert) are planned.

---

## 🚀 Quick Start

```go
import "github.com/entiqon/db/builder/selects"

sb := selects.New(nil).
    Fields("id").
    AppendFields("name", "username"). // expr + alias
    From("users u").
    Where("u.active = true").
    OrderBy("created_at DESC").
    Take(10)

sql, args, err := sb.Build()
if err != nil {
    log.Fatal(err)
}
fmt.Println(sql, args)
```

Output:

```sql
SELECT id, name AS username
FROM users AS u
WHERE u.active = true
ORDER BY created_at DESC
LIMIT 10
```

---

## 🔍 Current & Planned Builders

- ✅ `selects` — SELECT queries (implemented & fully tested)
- 🚧 `inserts` — INSERT queries (planned)
- 🚧 `updates` — UPDATE queries (planned)
- 🚧 `deletes` — DELETE queries (planned)
- 🚧 `merge`   — MERGE / UPSERT queries (planned)

---

## 🛠 Diagnostics

- `String()` → concise human-readable summary  
- `Debug()` → detailed diagnostic dump  

---

## 📦 Roadmap

- ✅ `selects` available today  
- 🚧 `inserts`, `updates`, `deletes`, `merge` in active development  
- 📝 Extended dialect support (Postgres, MySQL, SQLite) planned  

---

## 📄 License

MIT © Entiqon
