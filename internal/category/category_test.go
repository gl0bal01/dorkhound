package category

import (
	"sort"
	"testing"
)

func TestSlugsAreUnique(t *testing.T) {
	seen := make(map[string]bool)
	for _, e := range Catalog {
		if seen[e.Slug] {
			t.Fatalf("duplicate slug %q in catalog", e.Slug)
		}
		seen[e.Slug] = true
		if e.Title == "" {
			t.Errorf("slug %q has empty title", e.Slug)
		}
	}
}

func TestIsKnown(t *testing.T) {
	if !IsKnown("all") {
		t.Error("'all' should be known sentinel")
	}
	if !IsKnown("username") {
		t.Error("'username' must be known")
	}
	if IsKnown("usernames") {
		t.Error("misspelled 'usernames' must not be known")
	}
	if IsKnown("") {
		t.Error("empty string must not be known")
	}
}

func TestOrderForExportPutsCatalogFirst(t *testing.T) {
	// Catalog: image comes before username, both before crypto.
	got := OrderForExport([]string{"crypto", "username", "image"})
	want := []string{"image", "username", "crypto"}
	if !equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestOrderForExportAppendsUnknownAlphabetically(t *testing.T) {
	got := OrderForExport([]string{"zebra", "alpha", "username"})
	// username first (catalog), then alpha, zebra alphabetically.
	want := []string{"username", "alpha", "zebra"}
	if !equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAllSlugsIncludesAllAndCatalog(t *testing.T) {
	all := AllSlugs()
	if len(all) != len(Catalog)+1 {
		t.Errorf("AllSlugs len = %d, want %d", len(all), len(Catalog)+1)
	}
	if all[0] != "all" {
		t.Errorf("AllSlugs[0] = %q, want 'all'", all[0])
	}
	// Every catalog slug appears.
	got := make(map[string]bool)
	for _, s := range all {
		got[s] = true
	}
	for _, e := range Catalog {
		if !got[e.Slug] {
			t.Errorf("AllSlugs missing %q", e.Slug)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCatalogOrderStable(t *testing.T) {
	// Sanity: catalog isn't accidentally alphabetized — the order encodes
	// signal strength priority for export.
	slugs := Slugs()
	sorted := append([]string(nil), slugs...)
	sort.Strings(sorted)
	if equal(slugs, sorted) {
		t.Error("catalog appears alphabetized; order should encode signal strength")
	}
}
