package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// categoryOrder defines the display order for Discord output.
var categoryOrder = []string{
	"social",
	"records",
	"financial",
	"location",
	"forums",
	"people-db",
}

// categoryTitles maps category slugs to display titles.
var categoryTitles = map[string]string{
	"social":    "Social",
	"records":   "Records",
	"financial": "Financial",
	"location":  "Location",
	"forums":    "Forums",
	"people-db": "People-DB",
}

// Discord writes dorks to w in Discord-flavored Markdown format,
// grouped by category in a fixed order.
func Discord(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, engine string) {
	// Header
	fmt.Fprintf(w, "## OSINT Results: %s\n", c.Name)

	// Metadata line — only include non-empty fields
	var meta []string
	if c.Location != "" {
		meta = append(meta, fmt.Sprintf("**Location:** %s", c.Location))
	}
	if c.Age != 0 {
		meta = append(meta, fmt.Sprintf("**Age:** ~%d", c.Age))
	}
	if c.DOB != "" {
		meta = append(meta, fmt.Sprintf("**DOB:** %s", c.DOB))
	}
	if len(meta) > 0 {
		fmt.Fprintf(w, "%s\n", strings.Join(meta, " | "))
	}

	// Group dorks by category
	groups := groupByCategory(dorks)

	// Output in fixed category order
	for _, cat := range categoryOrder {
		ds, ok := groups[cat]
		if !ok || len(ds) == 0 {
			continue
		}
		title := categoryTitles[cat]
		if title == "" {
			title = cat
		}
		fmt.Fprintf(w, "\n### %s (%d links)\n", title, len(ds))
		for _, d := range ds {
			fmt.Fprintf(w, "- %s: %s\n", d.Label, d.URL(engine))
		}
	}
}

// groupByCategory groups a slice of dorks into a map keyed by Category.
func groupByCategory(dorks []dork.Dork) map[string][]dork.Dork {
	groups := make(map[string][]dork.Dork)
	for _, d := range dorks {
		groups[d.Category] = append(groups[d.Category], d)
	}
	return groups
}
