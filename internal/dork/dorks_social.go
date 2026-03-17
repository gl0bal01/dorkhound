package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	Register(generateSocialDorks)
}

func generateSocialDorks(c *caseinfo.Case) []Dork {
	name := c.Name
	if name == "" {
		return nil
	}

	var dorks []Dork

	// Global site-specific searches
	globalSites := []struct {
		site  string
		label string
	}{
		{"facebook.com", "Facebook"},
		{"linkedin.com", "LinkedIn"},
		{"instagram.com", "Instagram"},
		{"twitter.com", "Twitter/X"},
		{"tiktok.com", "TikTok"},
		{"youtube.com", "YouTube"},
		{"github.com", "GitHub"},
		{"medium.com", "Medium"},
		{"reddit.com", "Reddit"},
	}

	for _, s := range globalSites {
		q := fmt.Sprintf(`"%s" site:%s`, name, s.site)
		if c.Location != "" {
			q += fmt.Sprintf(` "%s"`, c.Location)
		}
		dorks = append(dorks, Dork{
			Query:    q,
			Category: "social",
			Region:   "global",
			Priority: 3,
			Label:    fmt.Sprintf("%s profile search", s.label),
		})
	}

	// Generic profile URL patterns
	q := fmt.Sprintf(`"%s" inurl:profile OR inurl:about OR inurl:bio`, name)
	if c.Location != "" {
		q += fmt.Sprintf(` "%s"`, c.Location)
	}
	dorks = append(dorks, Dork{
		Query:    q,
		Category: "social",
		Region:   "global",
		Priority: 2,
		Label:    "Generic profile/about/bio pages",
	})

	// Alias searches on major social platforms
	majorSites := []string{"facebook.com", "linkedin.com", "instagram.com", "twitter.com", "tiktok.com"}
	for _, alias := range c.Aliases {
		for _, site := range majorSites {
			q := fmt.Sprintf(`"%s" OR "%s" site:%s`, name, alias, site)
			if c.Location != "" {
				q += fmt.Sprintf(` "%s"`, c.Location)
			}
			dorks = append(dorks, Dork{
				Query:    q,
				Category: "social",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Alias %q on %s", alias, site),
			})
		}
	}

	// Region-specific: Russia (VK, OK.ru)
	ruSites := []struct {
		site  string
		label string
	}{
		{"vk.com", "VK"},
		{"ok.ru", "Odnoklassniki"},
	}
	for _, s := range ruSites {
		q := fmt.Sprintf(`"%s" site:%s`, name, s.site)
		if c.Location != "" {
			q += fmt.Sprintf(` "%s"`, c.Location)
		}
		dorks = append(dorks, Dork{
			Query:    q,
			Category: "social",
			Region:   "ru",
			Priority: 3,
			Label:    fmt.Sprintf("%s profile search", s.label),
		})
	}

	// Region-specific: France (Copainsdavant)
	frQ := fmt.Sprintf(`"%s" site:copainsdavant.linternaute.com`, name)
	if c.Location != "" {
		frQ += fmt.Sprintf(` "%s"`, c.Location)
	}
	dorks = append(dorks, Dork{
		Query:    frQ,
		Category: "social",
		Region:   "fr",
		Priority: 2,
		Label:    "Copainsdavant profile search",
	})

	// Region-specific: Germany (StayFriends)
	deQ := fmt.Sprintf(`"%s" site:stayfriends.de`, name)
	if c.Location != "" {
		deQ += fmt.Sprintf(` "%s"`, c.Location)
	}
	dorks = append(dorks, Dork{
		Query:    deQ,
		Category: "social",
		Region:   "de",
		Priority: 2,
		Label:    "StayFriends profile search",
	})

	// Region-specific: Netherlands (Hyves archives)
	nlQ := fmt.Sprintf(`"%s" site:hyves.nl OR "%s" "hyves"`, name, name)
	if c.Location != "" {
		nlQ += fmt.Sprintf(` "%s"`, c.Location)
	}
	dorks = append(dorks, Dork{
		Query:    nlQ,
		Category: "social",
		Region:   "nl",
		Priority: 2,
		Label:    "Hyves archive search",
	})

	return dorks
}
