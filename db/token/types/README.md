> Part of [Entiqon](https://github.com/entiqon/entiqon) / [Database](../../) / [Token](../)

# Token Types

The `types` package groups enum classifications that represent
SQL structures in a consistent, dependency-free way.  
These enums are consumed by higher-level tokens (`Field`, `Table`,
`Join`, `Condition`, …) and builders (`SelectBuilder`, …) to ensure
safe, auditable, and dialect-agnostic SQL generation.

---

## Purpose

- Centralize classification logic for tokens into strongly-typed enums.
- Keep enums dependency-free and reusable across token packages.
- Expose a consistent API (`IsValid`, `ParseFrom`, `String`) for validation
  and rendering.
- Provide clear separation between **type classification** and
  **token implementation**.

---

## Available Types

| Package                      | Purpose                                                                                                               |
|------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| [`identifier`](./identifier) | Classifies SQL expressions into categories: `Identifier`, `Subquery`, `Computed`, `Aggregate`, `Function`, `Literal`. |
| [`join`](./join)             | Classifies SQL JOIN operations: `Inner`, `Left`, `Right`, `Full`, `Cross`, `Natural`.                                 |
| [`condition`](./condition)   | Classifies SQL conditional expressions: `Single`, `And`, `Or` for WHERE, HAVING, and ON clauses.                      |

---

## Example

```go
package main

import (
    "fmt"
	
    "github.com/entiqon/entiqon/db/token/types/condition"
    "github.com/entiqon/entiqon/db/token/types/join"
    "github.com/entiqon/entiqon/db/token/types/identifier"
)

func main() {
    // Condition type
    c := condition.And
    fmt.Println(c, c.IsValid()) // AND true

    // Join type
    j := join.Cross
    fmt.Println(j.String()) // CROSS JOIN

    // Identifier type
    id := identifier.Function
    fmt.Println(id.String()) // Function
}
```

---

## License

Released under the [MIT License](../../../../LICENSE).  
Copyright © 2025 [Entiqon Contributors](https://entiqon.io)
