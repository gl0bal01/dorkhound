package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateDocumentDorks)
}

func generateDocumentDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" {
		return nil
	}
	name := c.Name
	return []Dork{
		{
			Query:    fmt.Sprintf(`"%s" filetype:pdf`, name),
			Category: "documents",
			Region:   "global",
			Priority: 2,
			Label:    "PDF documents",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:docx OR filetype:doc`, name),
			Category: "documents",
			Region:   "global",
			Priority: 2,
			Label:    "Word documents",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:xlsx OR filetype:xls`, name),
			Category: "documents",
			Region:   "global",
			Priority: 2,
			Label:    "Spreadsheet documents",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:pptx OR filetype:ppt`, name),
			Category: "documents",
			Region:   "global",
			Priority: 1,
			Label:    "Presentation documents",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:vcf`, name),
			Category: "documents",
			Region:   "global",
			Priority: 2,
			Label:    "vCard / address book",
		},
		{
			Query:    fmt.Sprintf(`"%s" filetype:txt site:pastebin.com`, name),
			Category: "documents",
			Region:   "global",
			Priority: 1,
			Label:    "Text files on Pastebin",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:scribd.com OR site:slideshare.net OR site:issuu.com OR site:academia.edu`, name),
			Category: "documents",
			Region:   "global",
			Priority: 2,
			Label:    "Document sharing platforms",
		},
		{
			Query:    fmt.Sprintf(`"%s" site:docdroid.net OR site:calameo.com`, name),
			Category: "documents",
			Region:   "global",
			Priority: 1,
			Label:    "Alternative document platforms",
		},
	}
}
