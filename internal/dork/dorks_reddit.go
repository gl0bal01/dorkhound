package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateRedditDorks)
}

func generateRedditDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" && len(c.Usernames) == 0 {
		return nil
	}

	var dorks []Dork

	for _, username := range c.Usernames {
		esc := url.PathEscape(username)
		encQ := url.QueryEscape(username)

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://www.reddit.com/user/%s/submitted/", esc),
				Category: "reddit",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Reddit submissions: /user/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://www.reddit.com/user/%s/comments/", esc),
				Category: "reddit",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Reddit comments: /user/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://old.reddit.com/user/%s/comments/.rss", esc),
				Category: "reddit",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Reddit RSS feed (no rate limit): /user/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://www.reddit.com/user/%s/about.json", esc),
				Category: "reddit",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Reddit JSON profile (no auth): /user/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://redditsearch.io/?term=&author=%s", esc),
				Category: "reddit",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("RedditSearch.io: author=%s (often down)", username),
			},
			Dork{
				Query:    fmt.Sprintf(`https://camas.github.io/reddit-search/#%%7B%%22author%%22%%3A%%22%s%%22%%7D`, encQ),
				Category: "reddit",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("Camas Reddit search: author=%s", username),
			},
		)
	}

	if c.Name != "" {
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" site:reddit.com OR site:old.reddit.com`, c.Name),
				Category: "reddit",
				Region:   "global",
				Priority: 2,
				Label:    "Reddit name search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:redditsearch.io`, c.Name),
				Category: "reddit",
				Region:   "global",
				Priority: 1,
				Label:    "RedditSearch.io name search",
			},
		)
	}

	return dorks
}
