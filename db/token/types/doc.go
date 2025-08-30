// Package types defines low-level, dependency-free classifications used by SQL tokens.
//
// These types provide common building blocks such as identifier
// classifications, enabling higher-level tokens (Field, Table, …)
// and builders (SelectBuilder, …) to share consistent semantics
// without cycles.
//
// Subpackages:
//   - identifier: classification of SQL expressions
package types
