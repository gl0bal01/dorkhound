package dork

import (
	"sort"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func TestDorkCount_FullInput(t *testing.T) {
	c := caseinfo.New("Jane Doe")
	c.Emails = []string{"jane@example.com"}
	c.Usernames = []string{"jdoe42"}
	c.PhotoURL = "https://example.com/j.jpg"

	generated := Generate(c)
	all := Filter(generated, []string{"all"}, []string{"all"})

	cats := map[string]int{}
	for _, d := range all {
		cats[d.Category]++
	}

	keys := make([]string, 0, len(cats))
	for k := range cats {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	t.Logf("Total dorks (all categories, all regions): %d", len(all))
	for _, cat := range keys {
		t.Logf("  %-20s %d", cat, cats[cat])
	}

	if len(all) < 200 {
		t.Logf("WARNING: total dork count %d is below expected 200+", len(all))
	}
}
