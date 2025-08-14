// File: db/contract/renderable.go

package contract

// Renderable defines the behavior for types that can be rendered into a
// string representation.
//
// Implementations of Renderable are responsible for producing a valid,
// human-readable or machine-consumable string output that reflects the
// current state of the object.
//
// The specific rendering format is determined by the implementing type,
// and may vary depending on the context in which the type is used.
//
// Example usage:
//
//	type SQLExpression struct {
//	    Expr string
//	}
//
//	func (e SQLExpression) Render() string {
//	    return e.Expr
//	}
//
//	func main() {
//	    var r contract.Renderable = SQLExpression{Expr: "SELECT * FROM users"}
//	    fmt.Println(r.Render()) // Output: SELECT * FROM users
//	}
type Renderable interface {
	// Render returns the string representation of the implementing type.
	// The returned string should be a valid representation of the type's
	// state, suitable for the context in which the Renderable is used.
	Render() string
}
