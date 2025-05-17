package token

// As sets the alias for the FieldToken.
func (f FieldToken) As(alias string) FieldToken {
	f.Alias = alias
	return f
}

// Field creates a standard identifier field (escaped by dialect).
func Field(name string) FieldToken {
	return FieldToken{Name: name}
}

// FieldExpr creates a raw SQL expression with an optional alias (unescaped).
func FieldExpr(expression string, alias string) FieldToken {
	return FieldToken{
		Name:  expression,
		Alias: alias,
		IsRaw: true,
	}
}
