package dork

import (
	"crypto/md5" // #nosec G501 -- Gravatar requires MD5 of the normalized email address.
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateGravatarDorks)
}

func generateGravatarDorks(c *caseinfo.Case) []Dork {
	if len(c.Emails) == 0 {
		return nil
	}
	var dorks []Dork
	for _, email := range c.Emails {
		h := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email)))) // #nosec G401 -- Gravatar requires MD5 for lookup URLs.
		hash := hex.EncodeToString(h[:])
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=400&d=404", hash),
				Category: "gravatar",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Gravatar avatar for %s", email),
			},
			Dork{
				Query:    fmt.Sprintf("https://www.gravatar.com/%s.json", hash),
				Category: "gravatar",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Gravatar profile JSON for %s", email),
			},
			Dork{
				Query:    fmt.Sprintf("https://en.gravatar.com/%s", hash),
				Category: "gravatar",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Gravatar profile page for %s", email),
			},
		)
	}
	return dorks
}
