package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateTwitterDorks)
}

func generateTwitterDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" && len(c.Usernames) == 0 {
		return nil
	}

	var dorks []Dork

	nitterInstances := []string{
		"nitter.net",
		"nitter.privacydev.net",
		"nitter.poast.org",
	}

	for _, username := range c.Usernames {
		esc := url.PathEscape(username)

		// Nitter mirror instances
		for _, instance := range nitterInstances {
			dorks = append(dorks, Dork{
				Query:    fmt.Sprintf("https://%s/%s", instance, esc),
				Category: "twitter",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Nitter: %s/%s (login-less mirror)", instance, username),
			})
		}

		// Twitter syndication endpoint (no auth required)
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("https://cdn.syndication.twimg.com/timeline/profile?screen_name=%s&suppress_response_codes=true", esc),
			Category: "twitter",
			Region:   "global",
			Priority: 3,
			Label:    "Twitter syndication profile API (public, no auth)",
		})

		// Wayback Machine archives
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://web.archive.org/web/*/twitter.com/%s", esc),
				Category: "twitter",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Wayback: historical twitter.com/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://web.archive.org/web/*/x.com/%s", esc),
				Category: "twitter",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Wayback: historical x.com/%s", username),
			},
		)

		// Google cache of twitter profile
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("cache:twitter.com/%s", esc),
			Category: "twitter",
			Region:   "global",
			Priority: 1,
			Label:    fmt.Sprintf("Google cache of twitter.com/%s", username),
		})

		// X advanced search — tweets from user
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("https://x.com/search?q=from%%3A%s&src=typed_query&f=live", esc),
			Category: "twitter",
			Region:   "global",
			Priority: 2,
			Label:    fmt.Sprintf("X advanced search: from:%s", username),
		})

		// X advanced search — tweets mentioning user
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf("https://x.com/search?q=%%40%s&src=typed_query&f=live", esc),
			Category: "twitter",
			Region:   "global",
			Priority: 2,
			Label:    fmt.Sprintf("X advanced search: @%s mentions", username),
		})
	}

	if c.Name != "" {
		nameEnc := url.QueryEscape(c.Name)

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" site:twitter.com OR site:x.com`, c.Name),
				Category: "twitter",
				Region:   "global",
				Priority: 2,
				Label:    "Twitter/X name search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:nitter.net OR site:nitter.poast.org OR site:nitter.privacydev.net`, c.Name),
				Category: "twitter",
				Region:   "global",
				Priority: 1,
				Label:    "Nitter name search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" "twitter.com" OR "x.com" site:web.archive.org`, c.Name),
				Category: "twitter",
				Region:   "global",
				Priority: 2,
				Label:    "Wayback twitter/x via search",
			},
			Dork{
				Query:    fmt.Sprintf("https://x.com/search?q=%%22%s%%22&src=typed_query", nameEnc),
				Category: "twitter",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("X advanced search: \"%s\"", c.Name),
			},
		)
	}

	return dorks
}
