package main

import (
	"strings"
	"testing"
)

func TestValidateCategoriesRejectsTypos(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"all", false},
		{"username", false},
		{"username,email,phone", false},
		{"usernames", true}, // common typo
		{"emial", true},
		{"social,bogus", true},
	}
	for _, c := range cases {
		err := validateCategories(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("validateCategories(%q) err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestValidateRegionsRejectsTypos(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"global", false},
		{"us", false},
		{"us,fr", false},
		{"usa", true},
		{"united-states", true},
		{"us,wrong", true},
	}
	for _, c := range cases {
		err := validateRegions(c.input)
		if (err != nil) != c.wantErr {
			t.Errorf("validateRegions(%q) err=%v, wantErr=%v", c.input, err, c.wantErr)
		}
	}
}

func TestValidateEngineRejectsTypos(t *testing.T) {
	for _, good := range []string{"google", "bing", "duckduckgo", "yandex", "GOOGLE"} {
		if err := validateEngine(good); err != nil {
			t.Errorf("validateEngine(%q) err=%v, want nil", good, err)
		}
	}
	for _, bad := range []string{"googel", "duckgo", "kagi"} {
		err := validateEngine(bad)
		if err == nil {
			t.Errorf("validateEngine(%q) err=nil, want error", bad)
			continue
		}
		if !strings.Contains(err.Error(), "google") {
			t.Errorf("validateEngine(%q) error should list valid choices, got %q", bad, err.Error())
		}
	}
}
