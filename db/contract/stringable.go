/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-24 05:42:00
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-24 05:42:04
 * @FilePath: db/contract/stringable.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: db/contract/stringable.go
//
// Stringable defines human-facing audit/log output. See package-level
// documentation in doc.go for an overview of all contracts and
// their distinct purposes.

package contract

// Stringable defines the contract for objects that can produce a human-facing
// representation of themselves.
//
// String() is intended for logging, auditing, or debugging. Implementations
// may include metadata, formatting, or contextual hints to improve
// observability.
//
// By convention, String() should integrate with Go’s fmt.Stringer so that
// using %s in fmt.Printf produces meaningful output.
//
// Example:
//
//	// Human-facing: suitable for logs and audits
//	fmt.Println(field) // Output: [Field expr="users.id" alias="uid"]
type Stringable interface {
	// String returns the human-facing representation of the object.
	// This is intended for logging, auditing, or debugging, and may
	// include additional metadata or formatting.
	String() string
}
