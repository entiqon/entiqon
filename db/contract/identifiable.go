package contract

// Identifiable is implemented by tokens that expose their core identity
// without aliasing. It separates two perspectives:
//
//   - Input() – the raw user-supplied value as provided to the constructor
//     (e.g. "users.id", "COUNT(*) AS total"). This is useful for auditing,
//     debugging, and traceability.
//
//   - Expr() – the normalized SQL expression derived from the input, suitable
//     for rendering in a dialect. This should not include any alias decoration.
//
// Other capabilities such as aliasing, cloning, error handling, rendering, etc.
// are defined in separate contracts (see doc.go).
type Identifiable interface {
	// Input returns the raw user-provided input string, typically exactly as it
	// was passed when constructing the token.
	Input() string

	// Expr returns the parsed or normalized SQL expression derived from Input(),
	// suitable for dialect rendering. This must exclude any alias component.
	Expr() string
}
