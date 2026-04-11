package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateTelegramDorks_Empty(t *testing.T) {
	got := generateTelegramDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateTelegramDorks_WithUsername(t *testing.T) {
	c := &caseinfo.Case{Usernames: []string{"jdoe42"}}
	got := generateTelegramDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for username case")
	}
	for _, d := range got {
		if d.Category != "telegram" {
			t.Errorf("expected category 'telegram', got %q", d.Category)
		}
	}
	foundTme := false
	for _, d := range got {
		if strings.HasPrefix(d.Query, "https://t.me/") {
			foundTme = true
			break
		}
	}
	if !foundTme {
		t.Error("expected at least one t.me direct URL")
	}
}

func TestGenerateTelegramDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateTelegramDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name-only case")
	}
	for _, d := range got {
		if d.Category != "telegram" {
			t.Errorf("expected category 'telegram', got %q", d.Category)
		}
	}
}
