package decimal

import (
	"fmt"
	"math/big"
)

// Decimal represents a high-precision decimal number.
// Backed by big.Rat for exact rational arithmetic.
type Decimal struct {
	r *big.Rat
}

// MustNew creates a Decimal from string or number and panics on error.
// Useful for constants and tests.
func MustNew(v any) Decimal {
	d, err := ParseFrom(v)
	if err != nil {
		panic(err)
	}
	return d
}

// String returns canonical string form (numerator/denominator).
func (d Decimal) String() string {
	return d.r.RatString()
}

// Float64 returns approximate float64 value.
func (d Decimal) Float64() float64 {
	f, _ := d.r.Float64()
	return f
}

// Add returns d + other.
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{r: new(big.Rat).Add(d.r, other.r)}
}

// Sub returns d - other.
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{r: new(big.Rat).Sub(d.r, other.r)}
}

// Mul returns d * other.
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{r: new(big.Rat).Mul(d.r, other.r)}
}

// Div returns d / other, or error if other == 0.
func (d Decimal) Div(other Decimal) (Decimal, error) {
	if other.r.Sign() == 0 {
		return Decimal{}, fmt.Errorf("division by zero")
	}
	return Decimal{r: new(big.Rat).Quo(d.r, other.r)}, nil
}
