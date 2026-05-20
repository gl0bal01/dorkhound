package dork

import "testing"

func TestTruecallerCountryDerivesFromPrefix(t *testing.T) {
	// Use 999... digits as the unmatched-prefix fallback case — no real
	// dialing prefix starts with 999 so the case-region branch can be
	// exercised independently.
	cases := []struct {
		digits, caseRegion, want string
	}{
		{"15551234567", "", "us"},
		{"33123456789", "", "fr"},
		{"44123456789", "", "uk"},
		{"49123456789", "", "de"},
		{"9991234567", "global", "us"}, // fallback when no prefix matches
		{"9991234567", "fr", "fr"},     // case region used as fallback
		{"9991234567", "all", "us"},    // 'all' is not a valid country
		{"61412345678", "", "au"},
	}
	for _, c := range cases {
		if got := truecallerCountry(c.digits, c.caseRegion); got != c.want {
			t.Errorf("truecallerCountry(%q, %q) = %q, want %q", c.digits, c.caseRegion, got, c.want)
		}
	}
}
