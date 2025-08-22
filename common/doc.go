// Package common provides foundational building blocks shared across Entiqon.
//
// It contains reusable utilities, error handling extensions, reflection helpers,
// and strongly typed parsers for primitive and structured values.
//
// Subpackages include:
//
//   - errors: structured error types (CausableError, ProcessStageError).
//   - extension: strongly typed value parsers and helpers, such as:
//   - boolean    → flexible boolean parsing (true/false, yes/no, on/off…)
//   - date       → date parsing and normalization
//   - decimal    → decimal parsing with precision control
//   - float      → floating-point parsing
//   - number     → integer parsing and rounding
//   - object     → reflection helpers (Exists, GetValue, SetValue)
//   - collection → generic typed collections with rich helper methods
//
// The common module enables consistency and code reuse across Entiqon’s ecosystem.
package common
