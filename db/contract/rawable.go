/**
 * @Author: Isidro Lopez isidro.lopezg@live.com
 * @Date: 2025-08-24 05:42:00
 * @LastEditors: Isidro Lopez isidro.lopezg@live.com
 * @LastEditTime: 2025-08-24 05:42:04
 * @FilePath: db/contract/rawable.go
 * @Description: 这是默认设置,可以在设置》工具》File Description中进行配置
 */
// File: db/contract/rawable.go
//
// Rawable defines dialect-agnostic SQL fragments. See package-level
// documentation in doc.go for an overview of all contracts and
// their distinct purposes.

package contract

// Rawable defines the contract for objects that can expose a
// generic SQL representation of themselves.
//
// Raw() is dialect-agnostic and produces a normalized SQL fragment,
// including alias if present, but without dialect-specific quoting,
// escaping, or rewriting.
//
// IsRaw() reports whether the object was constructed in a "raw" form,
// meaning it bypassed parsing logic and represents direct user intent.
//
// Contrast with:
//   - Render(): dialect-aware canonical representation (machine-facing).
//   - String(): human-facing representation for logs/audits.
//   - Debug(): developer-facing diagnostic view.
//
// Example:
//
//	u := table.New("users", "u")
//	fmt.Println(u.Raw())    // "users AS u"
//	fmt.Println(u.IsRaw())  // true
//	fmt.Println(u.Render()) // e.g. `"users" AS "u"` depending on dialect
type Rawable interface {
	// Raw returns the generic SQL fragment of the object.
	// This representation may include aliases but does not
	// apply dialect-specific transformations.
	//
	// If the object is invalid or errored, Raw() should
	// return an empty string.
	Raw() string

	// IsRaw reports whether the object was explicitly
	// constructed as raw, bypassing parsing logic.
	IsRaw() bool
}
