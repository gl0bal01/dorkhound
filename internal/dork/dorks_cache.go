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
		// Google's cache: operator was deprecated in early 2024. We use
		// "webcache.googleusercontent.com" via a regular search instead so
		// the dork still surfaces cached pages where they exist.
		Dork{
			Query:    fmt.Sprintf(`"%s" site:webcache.googleusercontent.com`, name),
			Category: "cache",
			Region:   "global",
			Priority: 2,
			Label:    "Google webcache (post-cache: replacement)",
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
	)

	if c.Location != "" {
		// Wayback CDX's url= parameter pattern-matches archived URLs by host
		// and path, not by content. A person's name fits that field poorly,
		// so we use a generic URL search that surfaces archived pages that
		// reference the name's slug form. The location gate stays as a
		// signal that the case has enough specifics to justify deep search.
		slug := url.QueryEscape(name)
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("https://web.archive.org/web/*/%s", slug),
			Category: "cache",
			Region:   "global",
			Priority: 3,
			Label:    "Wayback archived URLs containing name slug",
		})
	}

	return dorks
}
