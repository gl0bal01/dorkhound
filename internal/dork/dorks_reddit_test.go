package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateRedditDorks_Empty(t *testing.T) {
	got := generateRedditDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateRedditDorks_WithUsername(t *testing.T) {
	c := &caseinfo.Case{Usernames: []string{"jdoe42"}}
	got := generateRedditDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for username case")
	}
	for _, d := range got {
		if d.Category != "reddit" {
			t.Errorf("expected category 'reddit', got %q", d.Category)
		}
	}
	foundDirect := false
	for _, d := range got {
		if strings.HasPrefix(d.Query, "https://") {
			foundDirect = true
			break
		}
	}
	if !foundDirect {
		t.Error("expected at least one https:// direct URL")
	}
}

func TestGenerateRedditDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateRedditDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name-only case")
	}
	for _, d := range got {
		if d.Category != "reddit" {
			t.Errorf("expected category 'reddit', got %q", d.Category)
		}
	}
}
