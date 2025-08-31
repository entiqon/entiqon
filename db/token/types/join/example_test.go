package join_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/types/join"
)

func ExampleType_String() {
	fmt.Println(join.Inner.String())
	fmt.Println(join.Left.String())
	fmt.Println(join.Cross.String())
	fmt.Println(join.Natural.String())
	// Output:
	// INNER JOIN
	// LEFT JOIN
	// CROSS JOIN
	// NATURAL JOIN
}

func ExampleParseFrom() {
	j := join.ParseFrom("right join")
	if j.IsValid() {
		fmt.Println(j)
	}

	j2 := join.ParseFrom("invalid join")
	fmt.Println(j2.IsValid())
	// Output:
	// RIGHT JOIN
	// false
}
