package condition_test

import (
	"fmt"

	"github.com/entiqon/db/token/condition"
	ct "github.com/entiqon/db/token/types/condition"
)

//func ExampleNew() {
//	c := condition.New(99)
//	fmt.Println(fmt.Sprintf("[%q]: %s", c.Kind(), c.Error().Error()))
//
//	c = condition.New(ct.And)
//	fmt.Println(fmt.Sprintf("[%q]: %s", c.Kind(), c.Error().Error()))
//
//	c = condition.New(ct.And, 123456)
//	fmt.Println(fmt.Sprintf("[%q]: %s", c.Kind(), c.Error().Error()))
//
//	c = condition.New(ct.And, "")
//	fmt.Println(fmt.Sprintf("[%q]: %s", c.Kind(), c.Error().Error()))
//
//	c = condition.New(ct.And, "id = 1")
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	c = condition.New(ct.And, "id = :id", 1)
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	c = condition.New(ct.And, "id = ?", 1)
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	c = condition.New(ct.And, "id", 1)
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	c = condition.New(ct.And, "id != 1")
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	c = condition.New(ct.And, "id IN (1, 2, 3)")
//	fmt.Println(fmt.Sprintf("[%q]: expr=%q, value=%v", c.Kind(), c.Expr(), c.Value()))
//
//	// Output:
//	// ["Invalid"]: invalid condition type: Invalid
//	// ["AND"]: invalid condition input: []
//	// ["AND"]: expr must be string, got int
//	// ["AND"]: empty expression
//	// ["AND"]: expr="id = :id", value=1
//	// ["AND"]: expr="id = :id", value=1
//	// ["AND"]: expr="id = :id", value=1
//	// ["AND"]: expr="id = :id", value=1
//	// ["AND"]: expr="id != :id", value=1
//	// ["AND"]: expr="id IN (:id)", value=[1 2 3]
//}

func ExampleNew() {
	c := condition.New(ct.And, "COUNT(id) > 0")
	fmt.Println(fmt.Sprintf("[%s]: input=%s, name=%q, expr=%q, operator=%q, value=%v", c.Kind(), c.Input(), c.Name(), c.Expr(), c.Operator(), c.Value()))

	// Output:
	// [AND]: input=COUNT(id) > 0, name="count_id", expr="COUNT(id) > :count_id", operator=">", value=0
}
