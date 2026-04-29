package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// tracелабsCategoryOrder defines the preferred display order for TraceLabs submission.
var traceLabsCategoryOrder = []string{
	"image", "username", "email", "phone", "social", "people-db",
	"nuclei", "cache", "records", "documents", "forums", "financial",
	"location", "dating", "marketplace",
}

// TraceLabs writes a Markdown submission draft grouped by category, with
// each dork as a checkbox line: "- [ ] <label> — <url>". Intended to be
// pasted into the TraceLabs submission UI or a team scratchpad.
func TraceLabs(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, engine string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(w, "# TraceLabs Submission — %s\n", markdownText(c.Name))
	fmt.Fprintf(w, "Generated: %s\n\n", now)

	fmt.Fprintln(w, "## Case")
	tlField(w, "Name", c.Name)
	if len(c.Aliases) > 0 {
		tlField(w, "Aliases", strings.Join(c.Aliases, ", "))
	}
	tlField(w, "DOB", c.DOB)
	if c.Age != 0 {
		tlField(w, "Age", fmt.Sprintf("%d", c.Age))
	}
	tlField(w, "Location", c.Location)
	tlField(w, "Description", c.Description)
	if len(c.Associates) > 0 {
		tlField(w, "Associates", strings.Join(c.Associates, ", "))
	}
	if len(c.Emails) > 0 {
		tlField(w, "Emails", strings.Join(c.Emails, ", "))
	}
	if len(c.Phones) > 0 {
		tlField(w, "Phones", strings.Join(c.Phones, ", "))
	}
	if len(c.Usernames) > 0 {
		tlField(w, "Usernames", strings.Join(c.Usernames, ", "))
	}
	tlField(w, "Photo", c.PhotoSearchURL())

	// Group dorks by category.
	groups := groupByCategory(dorks)

	// Build ordered category list: known order first, then unknowns alphabetically.
	var orderedCats []string
	seen := make(map[string]bool)
	for _, cat := range traceLabsCategoryOrder {
		if _, ok := groups[cat]; ok {
			orderedCats = append(orderedCats, cat)
			seen[cat] = true
		}
	}
	var extra []string
	for cat := range groups {
		if !seen[cat] {
			extra = append(extra, cat)
		}
	}
	sort.Strings(extra)
	orderedCats = append(orderedCats, extra...)

	if len(orderedCats) > 0 {
		fmt.Fprintln(w, "\n## Leads by Category")
		for _, cat := range orderedCats {
			ds := groups[cat]
			if len(ds) == 0 {
				continue
			}
			// Sort by priority descending, then label ascending.
			sort.SliceStable(ds, func(i, j int) bool {
				if ds[i].Priority != ds[j].Priority {
					return ds[i].Priority > ds[j].Priority
				}
				return ds[i].Label < ds[j].Label
			})
			fmt.Fprintf(w, "\n### %s (%d)\n", markdownText(cat), len(ds))
			for _, d := range ds {
				fmt.Fprintf(w, "- [ ] %s — <%s>\n", markdownText(d.Label), neutralizeMentions(d.URL(engine)))
			}
		}
	}

	fmt.Fprintln(w, "\n## Notes")
	if c.Description != "" {
		fmt.Fprintln(w, markdownText(c.Description))
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "## Submission workflow")
	fmt.Fprintln(w, "1. Open each link in order of priority.")
	fmt.Fprintln(w, "2. Mark ☑ as you review.")
	fmt.Fprintln(w, "3. For each hit, capture: source URL, screenshot, extracted data.")
	fmt.Fprintln(w, "4. Cross-reference against TraceLabs flag categories before submission.")

	return nil
}

func tlField(w io.Writer, label, value string) {
	if value == "" {
		return
	}
	fmt.Fprintf(w, "- **%s:** %s\n", markdownText(label), markdownText(value))
}
