// Package table defines the token.Table type, which represents a SQL table
// or subquery source with optional alias. It provides strict construction
// rules, error preservation, and multiple forms of string output.
//
// A Table always preserves the original user input (for auditability),
// normalizes name and alias when possible, and never panics. Invalid inputs
// produce an errored Table that can be inspected via the Errorable contract.
//
// # Construction
//
// Tables are created using New(...):
//
//   - Plain table:
//     table.New("users") → users
//
//   - Aliased (inline):
//     table.New("users u")     → users AS u
//     table.New("users AS u")  → users AS u
//
//   - Aliased (explicit arguments):
//     table.New("users", "u")  → users AS u
//
//   - Subquery:
//     table.New("(SELECT COUNT(*) FROM users) AS t") → subquery with alias
//     table.New("(SELECT COUNT(*) FROM users)", "t") → subquery with alias
//
//   - Errors:
//     table.New("")            → errored
//     table.New("users AS")    → errored
//     table.New("users x y z") → errored
//
// # Contracts
//
// Table implements the following contracts from db/contract:
//
//   - Renderable → Render()
//   - Rawable    → Raw(), IsRaw()
//   - Stringable → String()
//   - Debuggable → Debug()
//   - Clonable   → Clone()
//   - Errorable  → IsErrored(), Error()
//
// Logging
//
//   - String() → concise, audit-friendly logs
//   - Debug()  → verbose, developer diagnostics
//
// Philosophy
//
//   - Never panic: always returns a *Table, even if errored.
//   - Auditability: preserves user input in all cases.
//   - Strict rules: invalid forms are rejected early.
//   - Delegation: parsing logic resides in table.New, not in builders.
package table
