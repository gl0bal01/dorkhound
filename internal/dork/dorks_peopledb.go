package dork

import (
	"fmt"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	Register(generatePeopleDBDorks)
}

func generatePeopleDBDorks(c *caseinfo.Case) []Dork {
	first := strings.ToLower(c.FirstName)
	last := strings.ToLower(c.LastName)

	// We need at least a last name to build URLs.
	if last == "" {
		return nil
	}

	var dorks []Dork

	// For single-name people (no first name), use last name only in URL slugs.
	// Helper to build hyphenated slug: "first-last" or just "last".
	slug := func(sep string) string {
		if first == "" {
			return last
		}
		return first + sep + last
	}

	// Helper for "first+last" or just "last" (query params).
	plusName := func() string {
		if first == "" {
			return last
		}
		return first + "+" + last
	}

	// Original-case helpers for site: search queries
	displayName := c.Name

	// --- US region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.spokeo.com/%s", slug("-")),
		Category: "people-db",
		Region:   "us",
		Priority: 3,
		Label:    "Spokeo people search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.whitepages.com/name/%s", slug("-")),
		Category: "people-db",
		Region:   "us",
		Priority: 3,
		Label:    "Whitepages lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.truepeoplesearch.com/results?name=%s", plusName()),
		Category: "people-db",
		Region:   "us",
		Priority: 3,
		Label:    "TruePeopleSearch lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.fastpeoplesearch.com/name/%s", slug("-")),
		Category: "people-db",
		Region:   "us",
		Priority: 3,
		Label:    "FastPeopleSearch lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.beenverified.com/people/%s/", slug("-")),
		Category: "people-db",
		Region:   "us",
		Priority: 3,
		Label:    "BeenVerified people search",
	})

	// --- CA region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.canada411.ca/search/?stype=si&what=%s", plusName()),
		Category: "people-db",
		Region:   "ca",
		Priority: 3,
		Label:    "Canada411 lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:canadapeoplesearch.com`, displayName),
		Category: "people-db",
		Region:   "ca",
		Priority: 3,
		Label:    "CanadaPeopleSearch lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.whitepages.ca/name/%s", slug("-")),
		Category: "people-db",
		Region:   "ca",
		Priority: 3,
		Label:    "WhitePages.ca lookup",
	})

	// --- UK region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:192.com`, displayName),
		Category: "people-db",
		Region:   "uk",
		Priority: 3,
		Label:    "192.com people search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:findmypast.co.uk`, displayName),
		Category: "people-db",
		Region:   "uk",
		Priority: 3,
		Label:    "FindMyPast search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:thephonebook.bt.com`, displayName),
		Category: "people-db",
		Region:   "uk",
		Priority: 3,
		Label:    "BT Phone Book search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:ukelectoralroll.co.uk`, displayName),
		Category: "people-db",
		Region:   "uk",
		Priority: 3,
		Label:    "UK Electoral Roll search",
	})

	// --- AU region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:whitepages.com.au`, displayName),
		Category: "people-db",
		Region:   "au",
		Priority: 3,
		Label:    "WhitePages Australia",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:peoplefinder.com.au`, displayName),
		Category: "people-db",
		Region:   "au",
		Priority: 3,
		Label:    "PeopleFinder Australia",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:reverseaustralia.com`, displayName),
		Category: "people-db",
		Region:   "au",
		Priority: 3,
		Label:    "ReverseAustralia search",
	})

	// --- RU region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://vk.com/search?c%%5Bq%%5D=%s", plusName()),
		Category: "people-db",
		Region:   "ru",
		Priority: 3,
		Label:    "VK people search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://ok.ru/search?st.query=%s", plusName()),
		Category: "people-db",
		Region:   "ru",
		Priority: 3,
		Label:    "OK.ru people search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:yandex.ru/people`, displayName),
		Category: "people-db",
		Region:   "ru",
		Priority: 3,
		Label:    "Yandex People search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:numbuster.com`, displayName),
		Category: "people-db",
		Region:   "ru",
		Priority: 3,
		Label:    "NumBuster search",
	})

	// --- FR region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.pagesjaunes.fr/pagesblanches/recherche?quoiqui=%s", plusName()),
		Category: "people-db",
		Region:   "fr",
		Priority: 3,
		Label:    "PagesBlanches (PagesJaunes) lookup",
	})

	// --- DE region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf("https://www.dastelefonbuch.de/Suche/%s", slug("+")),
		Category: "people-db",
		Region:   "de",
		Priority: 3,
		Label:    "DasTelefonbuch lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:telefonbuch.de`, displayName),
		Category: "people-db",
		Region:   "de",
		Priority: 3,
		Label:    "Telefonbuch.de search",
	})

	// --- AT region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:herold.at`, displayName),
		Category: "people-db",
		Region:   "at",
		Priority: 3,
		Label:    "Herold.at search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:dastelefonbuch.at`, displayName),
		Category: "people-db",
		Region:   "at",
		Priority: 3,
		Label:    "DasTelefonbuch.at search",
	})

	// --- NL region ---
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:detelefoongids.nl`, displayName),
		Category: "people-db",
		Region:   "nl",
		Priority: 3,
		Label:    "DeTelefoongids lookup",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:whitepages.nl`, displayName),
		Category: "people-db",
		Region:   "nl",
		Priority: 3,
		Label:    "WhitePages.nl search",
	})
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:numberway.nl`, displayName),
		Category: "people-db",
		Region:   "nl",
		Priority: 3,
		Label:    "Numberway.nl search",
	})

	return dorks
}
