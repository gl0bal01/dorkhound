package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateCryptoDorks_Empty(t *testing.T) {
	got := generateCryptoDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateCryptoDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateCryptoDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name case")
	}
	for _, d := range got {
		if d.Category != "crypto" {
			t.Errorf("expected category 'crypto', got %q", d.Category)
		}
	}
}

func TestGenerateCryptoDorks_WithUsername(t *testing.T) {
	c := &caseinfo.Case{Usernames: []string{"jdoe42"}}
	got := generateCryptoDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for username case")
	}
	foundENS := false
	for _, d := range got {
		if strings.Contains(d.Query, "ens.domains") || strings.Contains(d.Query, ".eth") {
			foundENS = true
			break
		}
	}
	if !foundENS {
		t.Error("expected at least one ENS dork")
	}
	foundHTTPS := false
	for _, d := range got {
		if strings.HasPrefix(d.Query, "https://") {
			foundHTTPS = true
			break
		}
	}
	if !foundHTTPS {
		t.Error("expected at least one https:// direct URL")
	}
}
