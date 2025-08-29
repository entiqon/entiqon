package helpers_test

import (
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
			})
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		invalidCases := []string{
			"",           // empty
			"123abc",     // starts with digit
			"9",          // single digit
			"-name",      // starts with dash
			"$amount",    // starts with $
			"first-name", // contains dash
			"user id",    // contains space
			"has.dot",    // contains dot
			"name!",      // contains punctuation
			"@tag",       // starts with @
			"abc#",       // contains #
			"abc?",       // contains ?
			" ",          // only whitespace
			"\t",         // tab
			"\n",         // newline
			"中",          // unicode non-ASCII
		}

		for _, tc := range invalidCases {
			tc := tc
			t.Run(tc, func(t *testing.T) {
				if helpers.IsValidIdentifier(tc) {
					t.Errorf("expected %q to be invalid", tc)
				}
			})
		}
	})
}
