package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateTwitterDorks_Empty(t *testing.T) {
	got := generateTwitterDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateTwitterDorks_WithUsername(t *testing.T) {
	c := &caseinfo.Case{Usernames: []string{"jdoe42"}}
	got := generateTwitterDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for username case")
	}
	for _, d := range got {
		if d.Category != "twitter" {
			t.Errorf("expected category 'twitter', got %q", d.Category)
		}
	}
	foundNitter := false
	for _, d := range got {
		if strings.Contains(d.Query, "nitter") {
			foundNitter = true
			break
		}
	}
	if !foundNitter {
		t.Error("expected at least one nitter direct URL")
	}
	foundHTTPS := false
	for _, d := range got {
		if strings.HasPrefix(d.Query, "https://") {
			foundHTTPS = true
			break
		}
	}
	if !foundHTTPS {
		t.Error("expected at least one dork with https:// URL")
	}
}

func TestGenerateTwitterDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateTwitterDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name-only case")
	}
	for _, d := range got {
		if d.Category != "twitter" {
			t.Errorf("expected category 'twitter', got %q", d.Category)
		}
	}
}
