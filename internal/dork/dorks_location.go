package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateLocationDorks)
}

func generateLocationDorks(c *caseinfo.Case) []Dork {
	name := c.Name
	if name == "" {
		return nil
	}

	var dorks []Dork

	// Location mentions with city (only when Location is set)
	if c.Location != "" {
		dorks = append(dorks, Dork{
			Query:    fmt.Sprintf(`"%s" "%s"`, name, c.Location),
			Category: "location",
			Region:   "global",
			Priority: 2,
			Label:    "Location mentions with city",
		})
	}

	// Google Maps reviews
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:google.com/maps OR "google review" OR "reviewed by"`, name),
		Category: "location",
		Region:   "global",
		Priority: 1,
		Label:    "Google Maps reviews",
	})

	// City/state/moved mentions
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" moved OR relocated OR "lives in" OR "based in"`, name),
		Category: "location",
		Region:   "global",
		Priority: 1,
		Label:    "City/state/moved mentions",
	})

	// Travel/booking/reservation
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" travel OR booking OR reservation OR itinerary OR flight`, name),
		Category: "location",
		Region:   "global",
		Priority: 1,
		Label:    "Travel/booking/reservation mentions",
	})

	// Passport/visa/travel
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" passport OR visa OR immigration OR "travel document"`, name),
		Category: "location",
		Region:   "global",
		Priority: 1,
		Label:    "Passport/visa/travel document mentions",
	})

	// Yelp/Foursquare
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" site:yelp.com OR site:foursquare.com`, name),
		Category: "location",
		Region:   "global",
		Priority: 1,
		Label:    "Yelp and Foursquare activity",
	})

	return dorks
}
