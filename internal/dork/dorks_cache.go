package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateCacheDorks)
}

func generateCacheDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name
	var dorks []Dork

	dorks = append(dorks,
		Dork{
			Query:    fmt.Sprintf(`cache:"%s"`, name),
			Category: "cache",
			Region:   "global",
			Priority: 2,
			Label:    "Google cache lookup",
		},
		Dork{
			Query:    fmt.Sprintf("https://web.archive.org/web/*/*%s*", url.PathEscape(name)),
			Category: "cache",
			Region:   "global",
			Priority: 2,
			Label:    "Wayback Machine URL template",
		},
		Dork{
			Query:    fmt.Sprintf(`"%s" site:translate.google.com`, name),
			Category: "cache",
			Region:   "global",
			Priority: 1,
			Label:    "Google Translate cache trick",
		},
		Dork{
			Query:    fmt.Sprintf("https://archive.ph/%s", url.PathEscape(name)),
			Category: "cache",
			Region:   "global",
			Priority: 2,
			Label:    "archive.today lookup",
		},
		Dork{
			Query:    fmt.Sprintf(`cache:"%s"`, name),
			Category: "cache",
			Region:   "global",
			Priority: 1,
			Label:    "Bing cache lookup",
		},
	)

	if c.Location != "" {
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=*%s*&output=json&limit=50", url.QueryEscape(name)),
			Category: "cache",
			Region:   "global",
			Priority: 3,
			Label:    "Wayback CDX API",
		})
	}

	return dorks
}
