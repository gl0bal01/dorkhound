package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/category"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// TestDiscordRendersEveryCatalogCategory guards against the regression that
// previously dropped 18 of 24 categories from Discord/clipboard output.
func TestDiscordRendersEveryCatalogCategory(t *testing.T) {
	var dorks []dork.Dork
	for _, e := range category.Catalog {
		dorks = append(dorks, dork.Dork{
			Query:    `"jane doe" site:example.com`,
			Category: e.Slug,
			Region:   "global",
			Priority: 1,
			Label:    "synthetic-" + e.Slug,
		})
	}

	var buf bytes.Buffer
	Discord(&buf, &caseinfo.Case{Name: "Jane Doe"}, dorks, "google")
	out := buf.String()

	for _, e := range category.Catalog {
		// Every label was unique per category so a missing slug means the
		// category section was not rendered at all.
		if !strings.Contains(out, "synthetic-"+e.Slug) {
			t.Errorf("Discord output missing category %q (label %q):\n%s",
				e.Slug, "synthetic-"+e.Slug, out)
		}
	}
}

// TestDiscordRendersUnknownCategoryAlphabetically ensures unrecognized
// category slugs still appear so the export never silently drops data.
func TestDiscordRendersUnknownCategoryAlphabetically(t *testing.T) {
	dorks := []dork.Dork{
		{Query: "https://example.com/a", Category: "social", Priority: 1, Label: "a"},
		{Query: "https://example.com/z", Category: "zzz-experimental", Priority: 1, Label: "z"},
	}
	var buf bytes.Buffer
	Discord(&buf, &caseinfo.Case{Name: "Test"}, dorks, "google")
	out := buf.String()
	if !strings.Contains(out, "zzz-experimental") && !strings.Contains(out, "Zzz-experimental") {
		t.Errorf("unknown category not rendered:\n%s", out)
	}
}
