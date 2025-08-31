package join_test

import (
	"fmt"

	"github.com/entiqon/entiqon/db/token/join"
	"github.com/entiqon/entiqon/db/token/table"
	join2 "github.com/entiqon/entiqon/db/token/types/join"
)

// ExampleNew demonstrates constructing token tokens
// in valid and invalid scenarios.
func ExampleNew() {
	// Valid token: users LEFT JOIN accounts
	users := table.New("User U")
	accounts := table.New("Account A")
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.String())

	// Invalid token: missing right table
	j2 := join.New(join2.Inner, users, nil, "U.id = 1")
	fmt.Println(j2.String())

	// Invalid token: empty condition
	j3 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j3.String())

	// Output:
	// ✅ token("INNER JOIN Account AS A ON A.userId = U.id")
	// ⛔ token("INNER JOIN  ON U.id = 1"): token requires both left and right tables
	// ⛔ token("INNER JOIN Account AS A ON "): token condition is empty
}

// ExampleToken_Kind demonstrates using Type().
func ExampleToken_kind() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j.Kind())

	// Output:
	// INNER JOIN
}

// ExampleToken_Left demonstrates using Left().
func ExampleToken_Left() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j.Left().Raw())

	// Output:
	// User AS U
}

// ExampleToken_Right demonstrates using Right().
func ExampleToken_Right() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j.Right().Raw())

	// Output:
	// Account AS A
}

// ExampleToken_Condition demonstrates using Condition().
func ExampleToken_Condition() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j.Condition())

	// Output:
	// A.userId = U.id
}

// ExampleToken_Clone demonstrates using Clone().
func ExampleToken_clone() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")

	clone := j.Clone()
	fmt.Println("Original: ", j.String())
	fmt.Println("Copied  : ", clone.String())

	// Output:
	// Original:  ✅ token("INNER JOIN Account AS A ON A.userId = U.id")
	// Copied  :  ✅ token("INNER JOIN Account AS A ON A.userId = U.id")
}

// ExampleToken_Debug demonstrates using Debug() on valid and invalid joins.
func ExampleToken_debug() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.Debug())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.Debug())

	// Output:
	// Join{Type:"INNER JOIN", Left:"User AS U", Right:"Account AS A", Condition:"A.userId = U.id", Valid:true}
	// Join{Type:"INNER JOIN", Left:"User AS U", Right:"Account AS A", Condition:"", Valid:false, Err:token condition is empty}
}

// ExampleToken_Error demonstrates using Error() on valid and invalid joins.
func ExampleToken_error() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.Error())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.Error())

	// Output:
	// <nil>
	// token condition is empty
}

// ExampleToken_IsErrored demonstrates using IsErrored() on valid and invalid joins.
func ExampleToken_ssErrored() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.IsErrored())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.IsErrored())

	// Output:
	// false
	// true
}

// ExampleToken_SetError demonstrates using SetError() to force an error on a token.
func ExampleToken_setError() {
	users := table.New("User U")
	accounts := table.New("Account A")
	j := join.New(join2.Inner, users, accounts, "A.userId = U.id")

	// Force an error
	j = j.SetError(fmt.Errorf("forced error"))

	fmt.Println(j.IsErrored())
	fmt.Println(j.Error())

	// Output:
	// true
	// forced error
}

// ExampleToken_Raw demonstrates using Raw() on valid and invalid joins.
func ExampleToken_raw() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.Raw())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.Raw())

	// Output:
	// INNER JOIN Account AS A ON A.userId = U.id
	//
}

// ExampleToken_IsRaw demonstrates using IsRaw() on valid and invalid joins.
func ExampleToken_isRaw() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.IsRaw())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.IsRaw())

	// Output:
	// false
	// false
}

// ExampleToken_Render demonstrates using Render() on valid and invalid joins.
func ExampleToken_render() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.Render())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.Render())

	// Output:
	// INNER JOIN Account AS A ON A.userId = U.id
	//
}

// ExampleToken_String demonstrates using String() on valid and invalid joins.
func ExampleToken_string() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.String())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.String())

	// Output:
	// ✅ token("INNER JOIN Account AS A ON A.userId = U.id")
	// ⛔ token("INNER JOIN Account AS A ON "): token condition is empty
}

// ExampleToken_IsValid demonstrates using IsValid() on valid and invalid joins.
func ExampleToken_isValid() {
	users := table.New("User U")
	accounts := table.New("Account A")

	// Valid token
	j1 := join.New(join2.Inner, users, accounts, "A.userId = U.id")
	fmt.Println(j1.IsValid())

	// Invalid token: empty condition
	j2 := join.New(join2.Inner, users, accounts, "")
	fmt.Println(j2.IsValid())

	// Output:
	// true
	// false
}
