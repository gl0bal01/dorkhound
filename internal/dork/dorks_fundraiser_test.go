package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateFundraiserDorks_Empty(t *testing.T) {
	got := generateFundraiserDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateFundraiserDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateFundraiserDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name case")
	}
	for _, d := range got {
		if d.Category != "fundraiser" {
			t.Errorf("expected category 'fundraiser', got %q", d.Category)
		}
	}
	foundGofundme := false
	for _, d := range got {
		if strings.Contains(d.Query, "gofundme") {
			foundGofundme = true
			break
		}
	}
	if !foundGofundme {
		t.Error("expected at least one gofundme dork")
	}
}
