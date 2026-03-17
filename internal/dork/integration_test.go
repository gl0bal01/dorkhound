package dork

import (
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestFullPipeline_NameOnly(t *testing.T) {
	c := caseinfo.New("John Doe")
	dorks := Generate(c)
	filtered := Filter(dorks, []string{"all"}, []string{"global"})
	sorted := Sort(filtered)

	if len(sorted) == 0 {
		t.Fatal("pipeline produced no dorks")
	}
	// Verify priority ordering
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Priority > sorted[i-1].Priority {
			t.Errorf("not sorted: index %d priority %d > index %d priority %d",
				i, sorted[i].Priority, i-1, sorted[i-1].Priority)
		}
	}
	// Verify no region-specific dorks leaked through global filter
	for _, d := range sorted {
		if d.Region != "global" {
			t.Errorf("global filter let through region %q: %s", d.Region, d.Label)
		}
	}
}

func TestFullPipeline_WithRegions(t *testing.T) {
	c := caseinfo.New("John Doe")
	dorks := Generate(c)
	filtered := Filter(dorks, []string{"all"}, []string{"us", "ca"})

	regions := map[string]bool{}
	for _, d := range filtered {
		regions[d.Region] = true
	}
	if !regions["global"] {
		t.Error("missing global dorks")
	}
	if !regions["us"] {
		t.Error("missing us dorks")
	}
	if !regions["ca"] {
		t.Error("missing ca dorks")
	}
	for r := range regions {
		if r != "global" && r != "us" && r != "ca" {
			t.Errorf("unexpected region %q leaked through filter", r)
		}
	}
}

func TestFullPipeline_CategoryFilter(t *testing.T) {
	c := caseinfo.New("John Doe")
	dorks := Generate(c)
	filtered := Filter(dorks, []string{"social"}, []string{"all"})
	for _, d := range filtered {
		if d.Category != "social" {
			t.Errorf("category filter let through %q", d.Category)
		}
	}
}

func TestFullPipeline_WithAllCaseFields(t *testing.T) {
	c := &caseinfo.Case{
		Name: "John Doe", Location: "Seattle, WA", Age: 34,
		DOB: "1990-01-15", Aliases: []string{"JD", "Johnny"},
		Associates: []string{"Jane Smith"}, Description: "Red hair",
	}
	c.FirstName, c.LastName = caseinfo.ParseName(c.Name)

	fullDorks := Generate(c)
	nameOnlyDorks := Generate(caseinfo.New("John Doe"))

	if len(fullDorks) <= len(nameOnlyDorks) {
		t.Errorf("full case (%d dorks) should produce more than name-only (%d)",
			len(fullDorks), len(nameOnlyDorks))
	}
}
