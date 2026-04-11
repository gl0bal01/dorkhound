package dork

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generatePhoneDorks)
}

func generatePhoneDorks(c *caseinfo.Case) []Dork {
	if len(c.Phones) == 0 {
		return nil
	}
	var dorks []Dork
	for _, phone := range c.Phones {
		variants := phoneVariants(phone)
		quoted := make([]string, len(variants))
		for i, v := range variants {
			quoted[i] = fmt.Sprintf(`"%s"`, v)
		}
		orQuery := strings.Join(quoted, " OR ")
		digits := digitsOnly(phone)

		dorks = append(dorks,
			Dork{
				Query:    orQuery,
				Category: "phone",
				Region:   "global",
				Priority: 3,
				Label:    "Phone number variants",
			},
			Dork{
				Query:    fmt.Sprintf(`%s site:pastebin.com OR site:ghostbin.com OR site:dpaste.org`, orQuery),
				Category: "phone",
				Region:   "global",
				Priority: 2,
				Label:    "Phone in paste sites",
			},
			Dork{
				Query:    fmt.Sprintf(`%s site:craigslist.org OR site:facebook.com/marketplace`, orQuery),
				Category: "phone",
				Region:   "global",
				Priority: 2,
				Label:    "Phone in classifieds/marketplaces",
			},
		)

		if c.Name != "" {
			dorks = append(dorks, Dork{
				Query:    fmt.Sprintf(`"%s" %s`, c.Name, orQuery),
				Category: "phone",
				Region:   "global",
				Priority: 3,
				Label:    "Phone + name combo",
			})
		}

		if digits != "" {
			dorks = append(dorks,
				Dork{
					Query:    fmt.Sprintf("https://www.truecaller.com/search/us/%s", digits),
					Category: "phone",
					Region:   "global",
					Priority: 3,
					Label:    "TrueCaller lookup",
				},
				Dork{
					Query:    fmt.Sprintf("https://www.whocalled.us/lookup/%s", digits),
					Category: "phone",
					Region:   "global",
					Priority: 2,
					Label:    "WhoCallsMe lookup",
				},
			)
		}
	}
	return dorks
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func phoneVariants(phone string) []string {
	raw := phone
	digits := digitsOnly(phone)
	seen := map[string]bool{raw: true}
	variants := []string{raw}

	add := func(v string) {
		if v != "" && !seen[v] {
			seen[v] = true
			variants = append(variants, v)
		}
	}

	if digits != "" && digits != raw {
		add(digits)
	}

	switch len(digits) {
	case 10:
		add(fmt.Sprintf("+1%s", digits))
		add(fmt.Sprintf("(%s) %s-%s", digits[0:3], digits[3:6], digits[6:10]))
		add(fmt.Sprintf("%s-%s-%s", digits[0:3], digits[3:6], digits[6:10]))
	case 11:
		if digits[0] == '1' {
			d := digits[1:]
			add(fmt.Sprintf("+%s", digits))
			add(fmt.Sprintf("(%s) %s-%s", d[0:3], d[3:6], d[6:10]))
			add(fmt.Sprintf("%s-%s-%s", d[0:3], d[3:6], d[6:10]))
		}
	}
	return variants
}
