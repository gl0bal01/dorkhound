package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateGitHubDorks)
}

func generateGitHubDorks(c *caseinfo.Case) []Dork {
	if len(c.Emails) == 0 && len(c.Usernames) == 0 && c.Name == "" {
		return nil
	}
	var dorks []Dork

	for _, email := range c.Emails {
		enc := url.QueryEscape(email)
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://github.com/search?q=author-email:%s&type=commits", enc),
				Category: "github",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("GitHub commits authored by %s", email),
			},
			Dork{
				Query:    fmt.Sprintf("https://github.com/search?q=committer-email:%s&type=commits", enc),
				Category: "github",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("GitHub commits committed by %s", email),
			},
			Dork{
				Query:    fmt.Sprintf("https://api.github.com/search/users?q=%s+in:email", enc),
				Category: "github",
				Region:   "global",
				Priority: 2,
				Label:    "GitHub users by email (API)",
			},
		)
	}

	for _, username := range c.Usernames {
		esc := url.PathEscape(username)
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://github.com/%s", esc),
				Category: "github",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("GitHub profile for %s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://github.com/%s.keys", esc),
				Category: "github",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("GitHub SSH public keys for %s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://github.com/%s.gpg", esc),
				Category: "github",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("GitHub GPG keys for %s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://api.github.com/users/%s/events/public", esc),
				Category: "github",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("GitHub public events for %s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://gist.github.com/%s", esc),
				Category: "github",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Gists by %s", username),
			},
		)
	}

	if c.Name != "" {
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf(`"%s" site:github.com inurl:README OR inurl:readme`, c.Name),
			Category: "github",
			Region:   "global",
			Priority: 1,
			Label:    fmt.Sprintf("GitHub README mentions of %s", c.Name),
		})
	}

	return dorks
}
