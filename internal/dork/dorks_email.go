package dork

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateEmailDorks)
}

func generateEmailDorks(c *caseinfo.Case) []Dork {
	if len(c.Emails) == 0 {
		return nil
	}
	var dorks []Dork
	for _, email := range c.Emails {
		parts := strings.SplitN(email, "@", 2)
		localPart := parts[0]

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s"`, email),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "Exact email mention",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:pastebin.com OR site:ghostbin.com OR site:dpaste.org OR site:rentry.co`, email),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "Email in paste sites",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:github.com`, email),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "Email on GitHub",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:gist.github.com`, email),
				Category: "email",
				Region:   "global",
				Priority: 2,
				Label:    "Email on GitHub Gists",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:twitter.com OR site:instagram.com OR site:reddit.com`, localPart),
				Category: "email",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Email local-part %q as social handle", localPart),
			},
			Dork{
				Query:    fmt.Sprintf("https://haveibeenpwned.com/account/%s", url.PathEscape(email)),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "HIBP breach lookup",
			},
			Dork{
				Query:    fmt.Sprintf("https://intelx.io/?s=%s", url.QueryEscape(email)),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "IntelX search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:scribd.com OR site:slideshare.net OR site:issuu.com`, email),
				Category: "email",
				Region:   "global",
				Priority: 2,
				Label:    "Email in document platforms",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:linkedin.com`, email),
				Category: "email",
				Region:   "global",
				Priority: 3,
				Label:    "Email on LinkedIn",
			},
		)
	}
	return dorks
}
