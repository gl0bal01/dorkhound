package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	Register(generateForumsDorks)
}

func generateForumsDorks(c *caseinfo.Case) []Dork {
	name := c.Name
	if name == "" {
		return nil
	}

	var dorks []Dork

	// Reddit
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:reddit.com`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Reddit mentions",
	})

	// Quora
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:quora.com`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Quora mentions",
	})

	// Forum profiles
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" inurl:forum OR inurl:thread OR inurl:member OR "forum profile"`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Forum profiles and threads",
	})

	// Education .edu sites
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:edu`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Education (.edu) mentions",
	})

	// Resume/CV PDFs
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" filetype:pdf resume OR curriculum OR "CV"`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Resume/CV PDFs",
	})

	// Contact cards vcf/vcard
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" filetype:vcf OR "vcard" OR "contact card"`, name),
		Category: "forums",
		Region:   "global",
		Priority: 1,
		Label:    "Contact cards (VCF/vCard)",
	})

	// Associates cross-reference
	for _, assoc := range c.Associates {
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf(`"%s" "%s"`, name, assoc),
			Category: "forums",
			Region:   "global",
			Priority: 2,
			Label:    fmt.Sprintf("Associate cross-reference: %s", assoc),
		})
	}

	return dorks
}
