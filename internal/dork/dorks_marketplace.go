package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateMarketplaceDorks)
}

func generateMarketplaceDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name
	return []Dork{
		{
			Query:    fmt.Sprintf(`"%s" site:ebay.com`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 2,
			Label:    "eBay seller/buyer",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:etsy.com`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 2,
			Label:    "Etsy shop/review",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:poshmark.com OR site:mercari.com OR site:depop.com`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 2,
			Label:    "Poshmark/Mercari/Depop",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:craigslist.org`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 1,
			Label:    "Craigslist listing",
		},
		{
			Query:    fmt.Sprintf(`"%s" "facebook marketplace"`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 1,
			Label:    "Facebook Marketplace",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:offerup.com OR site:letgo.com`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 1,
			Label:    "OfferUp/Letgo listing",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:nextdoor.com`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 2,
			Label:    "Nextdoor mention",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:amazon.com/hz/wishlist OR site:myregistry.com OR site:theknot.com/registry`, name),
			Category: "marketplace",
			Region:   "global",
			Priority: 2,
			Label:    "Gift registries",
		},
	}
}
