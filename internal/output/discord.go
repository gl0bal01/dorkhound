package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/category"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// Discord writes dorks to w in Discord-flavored Markdown format,
// grouped by category. Known categories are emitted in catalog order;
// any unknown categories are appended alphabetically so a new category
// added to the engine flows through without code changes here.
func Discord(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, engine string) {
	// Header
	fmt.Fprintf(w, "## OSINT Results: %s\n", markdownText(c.Name))

	// Metadata line — only include non-empty fields
	var meta []string
	if c.Location != "" {
		meta = append(meta, fmt.Sprintf("**Location:** %s", markdownText(c.Location)))
	}
	if c.Age != 0 {
		meta = append(meta, fmt.Sprintf("**Age:** ~%d", c.Age))
	}
	if c.DOB != "" {
		meta = append(meta, fmt.Sprintf("**DOB:** %s", markdownText(c.DOB)))
	}
	if len(meta) > 0 {
		fmt.Fprintf(w, "%s\n", strings.Join(meta, " | "))
	}

	// Group dorks by category.
	groups := groupByCategory(dorks)
	if len(groups) == 0 {
		return
	}

	present := make([]string, 0, len(groups))
	for cat := range groups {
		present = append(present, cat)
	}
	titles := category.Titles()

	for _, cat := range category.OrderForExport(present) {
		ds := groups[cat]
		if len(ds) == 0 {
			continue
		}
		title := titles[cat]
		if title == "" {
			title = capitalizeSlug(cat)
		}
		noun := "links"
		if len(ds) == 1 {
			noun = "link"
		}
		fmt.Fprintf(w, "\n### %s (%d %s)\n", markdownText(title), len(ds), noun)
		for _, d := range ds {
			fmt.Fprintf(w, "- %s: <%s>\n", markdownText(d.Label), neutralizeMentions(d.URL(engine)))
		}
	}
}

// capitalizeSlug uppercases the first ASCII letter of a slug for display.
func capitalizeSlug(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// groupByCategory groups a slice of dorks into a map keyed by Category.
func groupByCategory(dorks []dork.Dork) map[string][]dork.Dork {
	groups := make(map[string][]dork.Dork)
	for _, d := range dorks {
		groups[d.Category] = append(groups[d.Category], d)
	}
	return groups
}
