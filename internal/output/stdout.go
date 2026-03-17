package output

import (
	"fmt"
	"io"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// Stdout writes dorks to w in a plain-text format, grouped by category.
// Each category gets a header line, and each dork is printed as "  Label: URL".
func Stdout(w io.Writer, dorks []dork.Dork, engine string) {
	groups := groupByCategory(dorks)

	// Collect categories in order of first appearance.
	seen := make(map[string]bool)
	var cats []string
	for _, d := range dorks {
		if !seen[d.Category] {
			seen[d.Category] = true
			cats = append(cats, d.Category)
		}
	}

	for _, cat := range cats {
		fmt.Fprintf(w, "\n== %s ==\n", cat)
		for _, d := range groups[cat] {
			fmt.Fprintf(w, "  %s: %s\n", d.Label, d.URL(engine))
		}
	}
}
