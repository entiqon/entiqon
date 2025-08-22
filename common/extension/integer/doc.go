// Package integer provides parsing utilities to convert values into int.
//
// It reuses float.ParseFrom for all parsing, and truncates the result
// toward zero to return an integer value.
//
// Examples:
//
//	integer.ParseFrom(1.9)   -> 1
//	integer.ParseFrom(-1.9)  -> -1
//	integer.ParseFrom("42")  -> 42
//	integer.ParseFrom(true)  -> 1
//	integer.ParseFrom(false) -> 0
package integer
