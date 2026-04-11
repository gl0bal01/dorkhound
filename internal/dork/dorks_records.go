package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateRecordsDorks)
}

func generateRecordsDorks(c *caseinfo.Case) []Dork {
	name := c.Name
	if name == "" {
		return nil
	}

	var dorks []Dork

	// Helper to add DOB/age narrowing to a query
	narrow := func(q string) string {
		if c.DOB != "" {
			q += fmt.Sprintf(` "%s"`, c.DOB)
		}
		if c.Age != 0 {
			q += fmt.Sprintf(` "age %d"`, c.Age)
		}
		return q
	}

	// Court/arrest/case records
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" court OR arrest OR "case number" OR warrant`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "Court and arrest records",
	})

	// Property/deed records
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" property OR deed OR "real estate" OR parcel`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "Property and deed records",
	})

	// Marriage/divorce records
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" marriage OR divorce OR "marriage license" OR "divorce decree"`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "Marriage and divorce records",
	})

	// PDF public records
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" filetype:pdf public record OR filing OR document`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "PDF public records",
	})

	// Obituaries
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" site:legacy.com OR site:findagrave.com`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "Obituaries (Legacy.com, FindAGrave)",
	})

	// News articles
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" missing OR found OR "last seen" OR disappear`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "News articles",
	})

	// Newspaper archives / Archive.org
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" site:newspapers.com OR site:archive.org OR site:newspaperarchive.com`, name)),
		Category: "records",
		Region:   "global",
		Priority: 1,
		Label:    "Newspaper archives and Archive.org",
	})

	// Inmate/corrections records
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" inmate OR corrections OR "department of corrections" OR prisoner`, name)),
		Category: "records",
		Region:   "global",
		Priority: 2,
		Label:    "Inmate and corrections records",
	})

	// US-specific: mugshots and arrests
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" site:mugshots.com OR site:arrests.org`, name)),
		Category: "records",
		Region:   "us",
		Priority: 2,
		Label:    "US mugshots and arrest records",
	})

	// US-specific: VINE Link (victim notification)
	dorks = append(dorks, Dork{
		Query:    narrow(fmt.Sprintf(`"%s" site:vinelink.com`, name)),
		Category: "records",
		Region:   "us",
		Priority: 2,
		Label:    "VINELink victim notification",
	})

	// US missing-person registries
	dorks = append(dorks,
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:namus.gov`, name)),
			Category: "records",
			Region:   "us",
			Priority: 3,
			Label:    "NamUs (US DoJ missing-person registry)",
		},
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:charleyproject.org`, name)),
			Category: "records",
			Region:   "us",
			Priority: 3,
			Label:    "Charley Project (US cold cases)",
		},
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:doenetwork.org`, name)),
			Category: "records",
			Region:   "global",
			Priority: 2,
			Label:    "Doe Network (international missing persons)",
		},
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:missingkids.org`, name)),
			Category: "records",
			Region:   "us",
			Priority: 3,
			Label:    "NCMEC MissingKids",
		},
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:missingpeople.org.uk`, name)),
			Category: "records",
			Region:   "uk",
			Priority: 3,
			Label:    "Missing People UK",
		},
		Dork{
			Query:    narrow(fmt.Sprintf(`"%s" site:interpol.int inurl:Notices`, name)),
			Category: "records",
			Region:   "global",
			Priority: 2,
			Label:    "Interpol public notices",
		},
	)

	return dorks
}
