package decimal

import (
	"fmt"
	"math/big"
)

// ParseFrom attempts to convert any input into Decimal.
func ParseFrom(v any) (Decimal, error) {
	switch val := v.(type) {
	case Decimal:
		return val, nil
	case string:
		r, ok := new(big.Rat).SetString(val)
		if !ok {
			return Decimal{}, fmt.Errorf("invalid decimal string: %q", val)
		}
		return Decimal{r: r}, nil
	case int:
		return Decimal{r: big.NewRat(int64(val), 1)}, nil
	case int64:
		return Decimal{r: big.NewRat(val, 1)}, nil
	case float64:
		return Decimal{r: new(big.Rat).SetFloat64(val)}, nil
	default:
		return Decimal{}, fmt.Errorf("decimal.ParseFrom: unsupported type %T", v)
	}
}
