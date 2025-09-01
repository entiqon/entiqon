# Join Types

> Part of [Entiqon](https://github.com/entiqon/entiqon) / [Database](../../../) / [Token](../../) / [Types](../)

The `join` package classifies SQL `JOIN` clauses into canonical types.
This classification is purely **syntactic**, not semantic,
and is used internally by tokens and builders.

---

## Purpose

- Provide a dependency-free enum (`join.Type`) to represent
  supported SQL join operations.
- Enable higher-level tokens (`Table`, `Field`, …) and builders
  (`SelectBuilder`, …) to share consistent semantics without cycles.
- Support helpers such as `ParseFrom` to normalize input
  into one of the defined types.

---

## Types

The following join types are supported:

| Constant   | Description                             | Keyword        |
|------------|-----------------------------------------|----------------|
| `Invalid`  | Could not classify                      | `INVALID`      |
| `Inner`    | Standard inner join                     | `INNER JOIN`   |
| `Left`     | Left outer join                         | `LEFT JOIN`    |
| `Right`    | Right outer join                        | `RIGHT JOIN`   |
| `Full`     | Full outer join                         | `FULL JOIN`    |
| `Cross`    | Cartesian product                       | `CROSS JOIN`   |
| `Natural`  | Natural join (implicit column matching) | `NATURAL JOIN` |

---

## Example

```go
package main

import (
	"fmt"
	"entiqon/db/token/join"
)

func main() {
	// Direct usage
	j := join.Inner
	fmt.Println(j)

	// Parse from free-form string
	j2 := join.ParseFrom("cross join")
	if j2.IsValid() {
		fmt.Println(j2)
	}

	// Invalid or unknown types are handled safely
	fmt.Println(join.Invalid)
	fmt.Println(join.Type(99))

	// Output:
	// INNER JOIN
	// CROSS JOIN
	// INVALID
	// INVALID
}
```

---

## License

Released under the [MIT License](../../../../LICENSE).  
Copyright © 2025 [Entiqon Contributors](https://entiqon.io)

