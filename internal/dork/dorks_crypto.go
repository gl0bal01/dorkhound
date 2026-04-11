package dork

import (
	"fmt"
	"net/url"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateCryptoDorks)
}

func generateCryptoDorks(c *caseinfo.Case) []Dork {
	if c.Name == "" && len(c.Emails) == 0 && len(c.Usernames) == 0 {
		return nil
	}

	var dorks []Dork

	if c.Name != "" {
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" "bitcoin" OR "BTC" OR "ethereum" OR "crypto wallet"`, c.Name),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    "Crypto wallet / coin mention",
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" "ENS" OR ".eth domain"`, c.Name),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    "ENS / .eth domain mention",
			},
		)
	}

	for _, email := range c.Emails {
		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf(`"%s" site:etherscan.io`, email),
				Category: "crypto",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("Etherscan search: %s", email),
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:blockchain.com`, email),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("Blockchain.com search: %s", email),
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:blockchair.com`, email),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("Blockchair search: %s", email),
			},
		)
	}

	for _, username := range c.Usernames {
		esc := url.PathEscape(username)
		encQ := url.QueryEscape(username)

		dorks = append(dorks,
			Dork{
				Query:    fmt.Sprintf("https://app.ens.domains/%s.eth", esc),
				Category: "crypto",
				Region:   "global",
				Priority: 2,
				Label:    fmt.Sprintf("ENS lookup: %s.eth", username),
			},
			Dork{
				Query:    fmt.Sprintf(`"%s" site:opensea.io OR site:rarible.com`, username),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("NFT profile: %s", username),
			},
			Dork{
				Query:    fmt.Sprintf("https://etherscan.io/search?f=0&q=%s", encQ),
				Category: "crypto",
				Region:   "global",
				Priority: 1,
				Label:    fmt.Sprintf("Etherscan direct search: %s", username),
			},
		)
	}

	return dorks
}
