package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateAcademicDorks)
}

func generateAcademicDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	enc := url.QueryEscape(c.Name)
	return []Dork{
		{
			Query:    fmt.Sprintf("https://scholar.google.com/scholar?q=%s", enc),
			Category: "academic",
			Region:   "global",
			Priority: 2,
			Label:    "Google Scholar",
		},
		{
			Query:    fmt.Sprintf("https://orcid.org/orcid-search/search?searchQuery=%s", enc),
			Category: "academic",
			Region:   "global",
			Priority: 2,
			Label:    "ORCID researcher search",
		},
		{
			Query:    fmt.Sprintf("https://www.researchgate.net/search/researcher?q=%s", enc),
			Category: "academic",
			Region:   "global",
			Priority: 2,
			Label:    "ResearchGate researcher",
		},
		{
			Query:    fmt.Sprintf("https://www.semanticscholar.org/search?q=%s", enc),
			Category: "academic",
			Region:   "global",
			Priority: 1,
			Label:    "Semantic Scholar",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:pdf site:arxiv.org OR site:ssrn.com`, c.Name),
			Category: "academic",
			Region:   "global",
			Priority: 1,
			Label:    "ArXiv/SSRN papers",
		},
		{
			Query:    fmt.Sprintf(`"%s" "ph.d" OR "phd thesis" OR "dissertation"`, c.Name),
			Category: "academic",
			Region:   "global",
			Priority: 1,
			Label:    fmt.Sprintf(`PhD/dissertation search for %s`, c.Name),
		},
	}
}
