package dork

import (
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestGenerateVehicleDorks_Empty(t *testing.T) {
	got := generateVehicleDorks(&caseinfo.Case{})
	if len(got) != 0 {
		t.Errorf("empty case should produce 0 dorks, got %d", len(got))
	}
}

func TestGenerateVehicleDorks_WithName(t *testing.T) {
	c := &caseinfo.Case{Name: "Jane Doe"}
	got := generateVehicleDorks(c)
	if len(got) == 0 {
		t.Fatal("expected dorks for name case")
	}
	for _, d := range got {
		if d.Category != "vehicle" {
			t.Errorf("expected category 'vehicle', got %q", d.Category)
		}
	}
	foundVIN := false
	for _, d := range got {
		if strings.Contains(d.Query, "VIN") {
			foundVIN = true
			break
		}
	}
	if !foundVIN {
		t.Error("expected at least one VIN dork")
	}
}
