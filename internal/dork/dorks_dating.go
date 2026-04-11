package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateDatingDorks)
}

func generateDatingDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name
	loc := ""
	if c.Location != "" {
		loc = fmt.Sprintf(` "%s"`, c.Location)
	}

	queries := []struct {
		q        string
		priority int
		label    string
	}{
		{fmt.Sprintf(`"%s" site:pof.com OR site:okcupid.com OR site:match.com`, name), 2, "POF/OKCupid/Match profile"},
		{fmt.Sprintf(`"%s" site:bumble.com OR "bumble profile"`, name), 2, "Bumble profile"},
		{fmt.Sprintf(`"%s" site:hinge.co OR "hinge profile"`, name), 2, "Hinge profile"},
		{fmt.Sprintf(`"%s" "tinder" OR "tinder profile"`, name), 1, "Tinder mention"},
		{fmt.Sprintf(`"%s" site:seekingarrangement.com OR site:seeking.com`, name), 1, "SeekingArrangement profile"},
		{fmt.Sprintf(`"%s" site:adultfriendfinder.com`, name), 1, "AdultFriendFinder profile"},
		{fmt.Sprintf(`"%s" site:eharmony.com OR site:zoosk.com`, name), 1, "eHarmony/Zoosk profile"},
	}

	var dorks []Dork
	for _, q := range queries {
		dorks = append(dorks, Dork{
			Query:    q.q + loc,
			Category: "dating",
			Region:   "global",
			Priority: q.priority,
			Label:    q.label,
		})
	}

	for _, alias := range c.Aliases {
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf(`"%s" OR "%s" site:pof.com OR site:okcupid.com OR site:match.com%s`, name, alias, loc),
			Category: "dating",
			Region:   "global",
			Priority: 2,
			Label:    fmt.Sprintf("Alias %q on dating sites", alias),
		})
	}

	return dorks
}
