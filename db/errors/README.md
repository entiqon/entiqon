# 🛑 Entiqon DB Errors

> Part of [Entiqon](../../) / [Database](../)

The `db/errors` package defines **sentinel errors** used throughout Entiqon’s SQL
builder (`table`, `field`, `join`, `condition`).

These sentinels provide a consistent way to classify and detect common failure
modes across constructors and validators.

---

## 📌 Overview

Two error values are currently exported:

- **`UnsupportedTypeError`**  
  Returned when a constructor (e.g. `table.New`) receives a type that is not
  supported.  
  Example:
  ```go
  table.New(table.New("users"))
  // → error: unsupported type; if you want to create a copy, use Clone() instead
  ```

- **`InvalidIdentifierError`**  
  Returned when an identifier fails validation (e.g. contains invalid
  characters).  
  Example:
  ```go
  table.New("???")
  // → error: invalid table identifier
  ```

---

## 🔍 Usage with `errors.Is`

Sentinel errors are intended to be checked with the Go standard library
[`errors.Is`](https://pkg.go.dev/errors#Is):

```go
t := table.New("???")
if errors.Is(t.Error(), dberrors.InvalidIdentifierError) {
    log.Printf("invalid identifier: %v", t.Error())
}
```

This allows constructors (`table.New`, `field.New`) to wrap errors with
context-specific messages while still preserving the ability to detect the
underlying cause.

---

## ✅ Testing

The package includes `errors_test.go`, which provides both examples and
table-driven tests to ensure that wrapped errors can be matched with
`errors.Is`.

Run tests with:

```bash
go test ./db/errors
```

---

