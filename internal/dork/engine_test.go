package dork

import (
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestDorkURL_SearchQuery(t *testing.T) {
	d := Dork{Query: `"John Doe" site:facebook.com`}
	got := d.URL("google")
	want := "https://www.google.com/search?q=%22John+Doe%22+site%3Afacebook.com"
	if got != want {
		t.Errorf("URL(google) = %q, want %q", got, want)
	}
}

func TestDorkURL_DirectURL(t *testing.T) {
	d := Dork{Query: "https://www.example.com/profile/123"}
	got := d.URL("google")
	want := "https://www.example.com/profile/123"
	if got != want {
		t.Errorf("URL(google) = %q, want %q", got, want)
	}

	// Also test http://
	d2 := Dork{Query: "http://www.example.com/profile/123"}
	got2 := d2.URL("bing")
	want2 := "http://www.example.com/profile/123"
	if got2 != want2 {
		t.Errorf("URL(bing) = %q, want %q", got2, want2)
	}
}

func TestDorkURL_DefaultEngine(t *testing.T) {
	d := Dork{Query: "test query"}
	got := d.URL("nonexistent_engine")
	want := "https://www.google.com/search?q=test+query"
	if got != want {
		t.Errorf("URL(nonexistent_engine) = %q, want %q", got, want)
	}
}

func TestFilter_ByCategory(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Category: "social", Region: "global"},
		{Query: "q2", Category: "records", Region: "global"},
		{Query: "q3", Category: "social", Region: "global"},
		{Query: "q4", Category: "financial", Region: "global"},
	}

	filtered := Filter(dorks, []string{"social"}, []string{"all"})
	if len(filtered) != 2 {
		t.Fatalf("Filter by social: got %d dorks, want 2", len(filtered))
	}
	for _, d := range filtered {
		if d.Category != "social" {
			t.Errorf("expected category social, got %q", d.Category)
		}
	}
}

func TestFilter_ByCategory_All(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Category: "social", Region: "global"},
		{Query: "q2", Category: "records", Region: "global"},
		{Query: "q3", Category: "financial", Region: "global"},
	}

	filtered := Filter(dorks, []string{"all"}, []string{"all"})
	if len(filtered) != 3 {
		t.Fatalf("Filter by all categories: got %d dorks, want 3", len(filtered))
	}
}

func TestFilter_ByRegion_Global(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Category: "social", Region: "global"},
		{Query: "q2", Category: "social", Region: "us"},
		{Query: "q3", Category: "social", Region: "ca"},
	}

	filtered := Filter(dorks, []string{"all"}, []string{"global"})
	if len(filtered) != 1 {
		t.Fatalf("Filter by global region: got %d dorks, want 1", len(filtered))
	}
	if filtered[0].Region != "global" {
		t.Errorf("expected region global, got %q", filtered[0].Region)
	}
}

func TestFilter_ByRegion_Specific(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Category: "social", Region: "global"},
		{Query: "q2", Category: "social", Region: "us"},
		{Query: "q3", Category: "social", Region: "ca"},
		{Query: "q4", Category: "social", Region: "uk"},
	}

	filtered := Filter(dorks, []string{"all"}, []string{"us"})
	if len(filtered) != 2 {
		t.Fatalf("Filter by us region: got %d dorks, want 2 (global + us)", len(filtered))
	}
	for _, d := range filtered {
		if d.Region != "global" && d.Region != "us" {
			t.Errorf("expected region global or us, got %q", d.Region)
		}
	}
}

func TestFilter_ByRegion_All(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Category: "social", Region: "global"},
		{Query: "q2", Category: "social", Region: "us"},
		{Query: "q3", Category: "social", Region: "ca"},
		{Query: "q4", Category: "social", Region: "uk"},
	}

	filtered := Filter(dorks, []string{"all"}, []string{"all"})
	if len(filtered) != 4 {
		t.Fatalf("Filter by all regions: got %d dorks, want 4", len(filtered))
	}
}

func TestSort_ByPriority(t *testing.T) {
	dorks := []Dork{
		{Query: "q1", Priority: 1},
		{Query: "q2", Priority: 3},
		{Query: "q3", Priority: 2},
		{Query: "q4", Priority: 3},
	}

	sorted := Sort(dorks)

	// Should be sorted by priority descending (3, 3, 2, 1)
	if sorted[0].Priority != 3 || sorted[1].Priority != 3 {
		t.Errorf("expected first two to have priority 3, got %d and %d", sorted[0].Priority, sorted[1].Priority)
	}
	if sorted[2].Priority != 2 {
		t.Errorf("expected third to have priority 2, got %d", sorted[2].Priority)
	}
	if sorted[3].Priority != 1 {
		t.Errorf("expected fourth to have priority 1, got %d", sorted[3].Priority)
	}

	// Stable sort: among priority-3 dorks, q2 should come before q4
	if sorted[0].Query != "q2" || sorted[1].Query != "q4" {
		t.Errorf("stable sort violated: got %q, %q; want q2, q4", sorted[0].Query, sorted[1].Query)
	}

	// Original slice should not be modified
	if dorks[0].Query != "q1" {
		t.Error("Sort modified the original slice")
	}
}

func TestGenerate_Empty(t *testing.T) {
	// Save and restore the global registry to isolate this test
	saved := registry
	registry = nil
	defer func() { registry = saved }()

	c := caseinfo.New("Test Person")
	result := Generate(c)
	if len(result) != 0 {
		t.Errorf("Generate with no generators: got %d dorks, want 0", len(result))
	}
}
