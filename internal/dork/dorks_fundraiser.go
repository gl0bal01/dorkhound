package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateFundraiserDorks)
}

func generateFundraiserDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name

	loc := func(q string) string {
		if c.Location != "" {
			q += fmt.Sprintf(` "%s"`, c.Location)
		}
		return q
	}

	return []Dork{
		{
			Query:    loc(fmt.Sprintf(`"%s" site:gofundme.com`, name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 3,
			Label:    "GoFundMe search",
		},
		{
			Query:    loc(fmt.Sprintf(`"%s" site:fundly.com OR site:youcaring.com OR site:justgiving.com`, name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 2,
			Label:    "Fundraiser platforms (Fundly/YouCaring/JustGiving)",
		},
		{
			Query:    loc(fmt.Sprintf(`"%s" site:caringbridge.org`, name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 3,
			Label:    "CaringBridge medical journey",
		},
		{
			Query:    loc(fmt.Sprintf(`"%s" site:tribute.co OR site:tributearchive.com OR site:ever-loved.com`, name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 3,
			Label:    "Tribute/memorial pages",
		},
		{
			Query:    loc(fmt.Sprintf(`"%s" "in memory" OR "memorial fund" OR "benefit for"`, name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 2,
			Label:    "Memorial fund mentions",
		},
		{
			Query:    fmt.Sprintf("https://www.gofundme.com/s?q=%s", url.QueryEscape(name)),
			Category: "fundraiser",
			Region:   "global",
			Priority: 3,
			Label:    "GoFundMe direct search",
		},
		{
			Query:    fmt.Sprintf("https://www.charleyproject.org/search?q=%s", url.QueryEscape(name)),
			Category: "fundraiser",
			Region:   "us",
			Priority: 3,
			Label:    "Charley Project (US missing persons) direct search",
		},
		{
			Query:    fmt.Sprintf("https://www.namus.gov/MissingPersons/Search?SearchQuery=%s", url.QueryEscape(name)),
			Category: "fundraiser",
			Region:   "us",
			Priority: 3,
			Label:    "NamUs (US DoJ) direct search",
		},
	}
}
