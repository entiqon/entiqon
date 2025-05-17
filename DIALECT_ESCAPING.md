# Dialect and Escaping Strategy (Design Notes)

---

## 🧠 Context

The `Condition` type in Entiqon serves as a resolver and collector — not a renderer.

This means:

- It structures logical condition tokens (`WHERE`, `AND`, `OR`)
- It stores raw condition strings (`Key`) and param values (`Params`)
- It **does not** render SQL or escape values

---

## 🔐 Escaping Responsibility

Escaping values for SQL (e.g., `'active'` → `'active'` or `'O'Reilly'` → `'O''Reilly'`) is **dialect-specific** and must be done at the point of SQL rendering.

This includes:
- Rendering `token.Condition.Raw` in a preview
- Producing literal SQL for engines that require inline formatting
- Supporting dialect abstractions (e.g., `MySQL`, `PostgreSQL`, `SQLite`)

---

## 📍 Location for Escaping

Escaping must be applied:

- At the **Build()** method of each builder
- Or during dialect-specific rendering (`DialectEngine.RenderCondition(...)`)

Never inside:
- `Condition.Set(...)`
- `NewCondition(...)`
- `AppendCondition(...)`

---

## 🛠️ Future Path

A `Dialect` interface may define:

```go
type Dialect interface {
	Escape(value any) string
}
```

Used like:
```go
for _, c := range sb.conditions {
	raw = strings.Replace(c.Key, "?", dialect.Escape(val), 1)
}
```

---

> Condition defines logic. Dialect defines output.
