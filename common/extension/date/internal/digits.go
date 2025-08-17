package internal

import (
	"errors"
	"fmt"
	"math"
)

// AllDigits reports whether s consists only of ASCII decimal digits.
func AllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func DecimalDigits(x int64) int {
	if x == 0 {
		return 1
	}
	if x < 0 {
		x = -x
	}
	d := 0
	for x > 0 {
		x /= 10
		d++
	}
	return d
}

func LooksLikeEpochMillis(x int64) bool {
	d := DecimalDigits(x)
	return d == 12 || d == 13
}

func ToSecondsSigned(x int64) (int64, error) {
	if LooksLikeEpochMillis(x) {
		return 0, fmt.Errorf("date.ParseFrom: got %d (looks like milliseconds); integers are interpreted as seconds", x)
	}
	return x, nil
}

func ToSecondsUnsigned(x uint64) (int64, error) {
	if x > math.MaxInt64 {
		return 0, errors.New("date.ParseFrom: unsigned value overflows int64 seconds")
	}
	s := int64(x)
	if LooksLikeEpochMillis(s) {
		return 0, fmt.Errorf("date.ParseFrom: got %d (looks like milliseconds); integers are interpreted as seconds", s)
	}
	return s, nil
}
