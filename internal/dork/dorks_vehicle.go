package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateVehicleDorks)
}

func generateVehicleDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name

	dorks := []Dork{
		{
			Query:    fmt.Sprintf(`"%s" "VIN" OR "vehicle identification"`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "VIN / vehicle identification",
		},
		{
			Query:    fmt.Sprintf(`"%s" "license plate" OR "plate number"`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "License plate mention",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:carfax.com OR site:autocheck.com`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "CarFax / AutoCheck vehicle history",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:craigslist.org inurl:cars`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 2,
			Label:    "Craigslist cars listings",
		},
		{
			Query:    fmt.Sprintf(`"%s" "DMV" OR "department of motor vehicles"`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "DMV records mention",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:cars.com OR site:autotrader.com`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "Cars.com / AutoTrader seller listings",
		},
		{
			Query:    fmt.Sprintf(`"%s" "fleet vehicle" OR "commercial vehicle"`, name),
			Category: "vehicle",
			Region:   "global",
			Priority: 1,
			Label:    "Fleet / commercial vehicle mention",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:nicb.org`, name),
			Category: "vehicle",
			Region:   "us",
			Priority: 1,
			Label:    "NICB (National Insurance Crime Bureau) theft records",
		},
	}

	return dorks
}
