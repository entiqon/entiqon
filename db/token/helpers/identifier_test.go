package helpers_test

import (
	"strings"
	"testing"

	"github.com/entiqon/entiqon/db/token/helpers"
)

func TestIdentifier(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		validCases := []string{
			"id",
			"user_id",
			"_col123",
			"U1",
			"columnName",
			"Name123",
			"n",
			"__hidden",
			"a_b_c",
			"XYZ",
		}

		for _, tc := range validCases {
			tc := tc
			t.Run(tc, func(t *testing.T) {
				if !helpers.IsValidIdentifier(tc) {
					t.Errorf("expected %q to be valid", tc)
				}
				if err := helpers.ValidateIdentifier(tc); err != nil {
					t.Errorf("expected %q valid, got error: %v", tc, err)
				}
			})
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		invalidCases := map[string]string{
			"":           "empty",
			"123abc":     "digit",
			"9":          "digit",
			"-name":      "syntax",
			"$amount":    "syntax",
			"first-name": "syntax",
			"user id":    "syntax",
			"has.dot":    "syntax",
			"name!":      "syntax",
			"@tag":       "syntax",
			"abc#":       "syntax",
			"abc?":       "syntax",
			" ":          "syntax",
			"\t":         "syntax",
			"\n":         "syntax",
			"中":          "syntax",
			"café":       "syntax",
			"mañana":     "syntax",
			"niño":       "syntax",
		}

		for tc, expected := range invalidCases {
			tc := tc
			expected := expected
			t.Run(tc, func(t *testing.T) {
				if helpers.IsValidIdentifier(tc) {
					t.Errorf("expected %q to be invalid", tc)
				}
				err := helpers.ValidateIdentifier(tc)
				if err == nil {
					t.Errorf("expected error for %q, got nil", tc)
					return
				}
				if !strings.Contains(strings.ToLower(err.Error()), expected) {
					t.Errorf("expected error about %s, got %v", expected, err)
				}
			})
		}
	})
}
