package generic_test

import (
	"fmt"
	"time"

	"github.com/entiqon/db/dialect/generic"
)

func Example() {
	d := generic.New()
	fmt.Println(d.Name())
	// Output:
	// generic
}

func Example_options() {
	d := generic.New()
	opts := d.Options()
	fmt.Printf("Name=%s, Placeholder=%s\n", opts.Name, opts.PlaceholderStyle)
	// Output:
	// Name=generic, Placeholder=?
}

func Example_quoteIdentifier() {
	d := generic.New()
	fmt.Println(d.QuoteIdentifier("users"))
	fmt.Println(d.QuoteIdentifier("UserData"))
	fmt.Println(d.QuoteIdentifier("order items"))
	// Output:
	// users
	// "UserData"
	// "order items"
}

func Example_quoteLiteral() {
	d := generic.New()
	fmt.Println(d.QuoteLiteral("O'Reilly"))
	fmt.Println(d.QuoteLiteral(true))
	fmt.Println(d.QuoteLiteral(42))
	fmt.Println(d.QuoteLiteral(time.Date(2025, 9, 19, 3, 30, 0, 0, time.UTC)))
	// Output:
	// 'O''Reilly'
	// TRUE
	// 42
	// '2025-09-19 03:30:00'
}

func Example_paginationSyntax() {
	d := generic.New()
	fmt.Println(d.PaginationSyntax(10, 0))
	fmt.Println(d.PaginationSyntax(10, 20))
	fmt.Println(d.PaginationSyntax(0, 0))
	// Output:
	//  LIMIT 10
	//  LIMIT 10 OFFSET 20
	//
}

func Example_placeholder() {
	d := generic.New()
	fmt.Println(d.Placeholder(1))
	fmt.Println(d.Placeholder(99))
	// Output:
	// ?
	// ?
}
