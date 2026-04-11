package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateTelegramDorks)
}

func generateTelegramDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" && len(c.Usernames) == 0 {
		return nil
	}

	var dorks []Dork

	for _, username := range c.Usernames {
		esc := url.PathEscape(username)
		encQ := url.QueryEscape(username)

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://t.me/%s", esc),
				Category: "telegram",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Telegram direct: t.me/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://t.me/s/%s", esc),
				Category: "telegram",
				Region:   "global",
				Priority: 3,
				Label:    fmt.Sprintf("Telegram public channel preview: t.me/s/%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://tgstat.com/en/channel/@%s", esc),
				Category: "telegram",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("TGStat analytics: @%s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://tgstat.com/search?q=%s", encQ),
				Category: "telegram",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("TGStat search: %s", username),
			},
		)
	}

	if c.Name != "" {
		nameEnc := url.QueryEscape(c.Name)

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" site:t.me`, c.Name),
				Category: "telegram",
				Region:   "global",
				Priority: 1,
				Label:    "Telegram name search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:tgstat.com`, c.Name),
				Category: "telegram",
				Region:   "global",
				Priority: 3,
				Label:    "TGStat name search",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:telegra.ph`, c.Name),
				Category: "telegram",
				Region:   "global",
				Priority: 2,
				Label:    "Telegraph (Telegram blog) name search",
			},
			Dork{
				Query:    fmt.Sprintf("https://tgstat.com/search?q=%s", nameEnc),
				Category: "telegram",
				Region:   "global",
				Priority: 3,
				Label:    "TGStat direct name search",
			},
		)
	}

	return dorks
}
