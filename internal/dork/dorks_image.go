package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateImageDorks)
}

func generateImageDorks(c *caseinfo.Case) []Dork {
	if c.PhotoURL == "" && c.Name == "" {
		return nil
	}

	var dorks []Dork

	if c.PhotoURL != "" {
		p := url.QueryEscape(c.PhotoURL)
		pp := url.PathEscape(c.PhotoURL)
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://lens.google.com/uploadbyurl?url=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 3,
				Label:    "Google Lens reverse search",
			},
			Dork{
				Query:    fmt.Sprintf("https://www.google.com/searchbyimage?image_url=%s&safe=off", p),
				Category: "image",
				Region:   "global",
				Priority: 3,
				Label:    "Google Images reverse search",
			},
			Dork{
				Query:    fmt.Sprintf("https://yandex.com/images/search?url=%s&rpt=imageview", p),
				Category: "image",
				Region:   "global",
				Priority: 3,
				Label:    "Yandex reverse search (best for faces)",
			},
			Dork{
				Query:    fmt.Sprintf("https://tineye.com/search?url=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 3,
				Label:    "TinEye reverse search",
			},
			Dork{
				Query:    fmt.Sprintf("https://www.bing.com/images/searchbyimage?cbir=sbi&imgurl=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 2,
				Label:    "Bing Visual Search",
			},
			Dork{
				Query:    fmt.Sprintf("https://pimeyes.com/en/search?search=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 3,
				Label:    "PimEyes face recognition",
			},
			Dork{
				Query:    fmt.Sprintf("https://saucenao.com/search.php?url=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 2,
				Label:    "SauceNAO reverse search",
			},
			Dork{
				Query:    fmt.Sprintf("https://iqdb.org/?url=%s", p),
				Category: "image",
				Region:   "global",
				Priority: 1,
				Label:    "IQDB reverse search",
			},
			Dork{
				Query:    fmt.Sprintf("https://karmadecay.com/%s", pp),
				Category: "image",
				Region:   "global",
				Priority: 1,
				Label:    "KarmaDecay reddit reverse search",
			},
		)
	}

	if c.Name != "" {
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" site:flickr.com OR site:500px.com OR site:smugmug.com`, c.Name),
				Category: "image",
				Region:   "global",
				Priority: 2,
				Label:    "Named photo platforms",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" filetype:jpg OR filetype:png site:imgur.com`, c.Name),
				Category: "image",
				Region:   "global",
				Priority: 2,
				Label:    "Imgur named photo",
			},
		)
	}

	return dorks
}
