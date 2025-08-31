# Identifier Types

The `identifier` package classifies SQL expressions into broad syntactic
categories. This classification is purely **syntactic**, not semantic,
and is used internally by token resolvers and helpers.

---

## Purpose

- Provide a dependency-free enum (`identifier.Type`) to represent
  the form of a SQL expression.
- Enable higher-level tokens (`Field`, `Table`, …) and builders
  (`SelectBuilder`, …) to share consistent semantics without cycles.
- Support helpers such as `ClassifyExpression` to normalize input
  into one of the defined categories.

---

## Categories

The following categories are supported:

| Constant     | Description                                                      | Example                        |
|--------------|------------------------------------------------------------------|--------------------------------|
| `Invalid`    | Could not classify                                               | `""`                           |
| `Subquery`   | Parenthesized SELECT                                             | `(SELECT * FROM users)`        |
| `Computed`   | Other parenthesized expression                                   | `(a + b)`                      |
| `Aggregate`  | Aggregate function                                               | `SUM(qty)`, `COUNT(*)`         |
| `Function`   | Any other function or call                                       | `JSON_EXTRACT(data, '$.id')`   |
| `Literal`    | Quoted string or numeric constant                                | `'abc'`, `"xyz"`, `42`         |
| `Identifier` | Plain table or column name (default fallback)                    | `users`, `id`                  |

---

## Example

```go
package main

import (
	"fmt"

	"entiqon/db/token/types/identifier"
)

func main() {
	var t identifier.Type

	// Functions render as capitalized names
	t = identifier.Function
	fmt.Println(t)

	// Subqueries render as capitalized names
	t = identifier.Subquery
	fmt.Println(t)

	// Invalid or unknown types are handled safely
	fmt.Println(identifier.Invalid)
	fmt.Println(identifier.Type(99))

	// Output:
	// Function
	// Subquery
	// Invalid
	// Unknown
}
```

---

## License

Released under the [MIT License](../../../../LICENSE).  
Copyright © 2025 [Entiqon Contributors](https://entiqon.io)
