> Part of [Entiqon](https://github.com/entiqon/entiqon) / [Database](../../../) / [Token](../../) / [Types](../)

# Condition Types

The `condition` package classifies SQL conditional expressions
(`WHERE`, `HAVING`, `ON`) into canonical types.
This classification is purely **syntactic**, not semantic,
and is used internally by tokens and builders.

---

## Purpose

- Provide a dependency-free enum (`condition.Type`) to represent
  supported SQL condition kinds.
- Enable higher-level tokens (`Table`, `Field`, `Join`, …) and builders
  (`SelectBuilder`, …) to share consistent semantics without cycles.
- Support helpers such as `ParseFrom` to normalize input
  into one of the defined types.

---

## Types

The following condition types are supported:

| Constant  | Description                        | Keyword   |
|-----------|------------------------------------|-----------|
| `Invalid` | Could not classify                 | `INVALID` |
| `Single`  | Single expression (default)        | *(none)*  |
| `And`     | Logical conjunction of expressions | `AND`     |
| `Or`      | Logical disjunction of expressions | `OR`      |

---

## Example

```go
package main

import (
    "fmt"
	
    "github.com/entiqon/entiqon/db/token/condition"
)

func main() {
    // Direct usage
    c := condition.Single
    fmt.Println(c)

    // Parse from free-form string
    c2 := condition.ParseFrom("and")
    if c2.IsValid() {
        fmt.Println(c2)
    }

	// Parse int
	c3 := condition.ParseFrom(3)
	if c3.IsValid() {
		fmt.Println(c2)
	}

    // Invalid or unknown types are handled safely
    fmt.Println(condition.Invalid)
    fmt.Println(condition.Type(99))

    // Output:
    //
    // AND
	// OR
    // Invalid
    // Invalid
}
```

---

## Integration

Condition types are consumed by higher-level builders to render SQL clauses:

- `SelectBuilder.Where` for row filtering  
- `SelectBuilder.Having` for group filtering  
- `Join.On` for join predicates  

Future versions will extend conditions to support structured construction,
parameter binding, and logical composition (`condition.And(c1, c2)`).

---

## License

Released under the [MIT License](../../../../LICENSE).  
Copyright © 2025 [Entiqon Contributors](https://entiqon.io)
