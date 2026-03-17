package dork

import (
	"fmt"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
)

func init() {
	register(generateFinancialDorks)
}

func generateFinancialDorks(c *caseinfo.Case) []Dork {
	name := c.Name
	if name == "" {
		return nil
	}

	var dorks []Dork

	// PayPal/Venmo/CashApp mentions
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" paypal OR venmo OR cashapp OR "cash app"`, name),
		Category: "financial",
		Region:   "global",
		Priority: 1,
		Label:    "PayPal/Venmo/CashApp mentions",
	})

	// Bank/account/transaction mentions
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" bank OR account OR transaction OR "wire transfer"`, name),
		Category: "financial",
		Region:   "global",
		Priority: 1,
		Label:    "Bank and transaction mentions",
	})

	// Loan/mortgage/credit mentions
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" loan OR mortgage OR credit OR "credit report"`, name),
		Category: "financial",
		Region:   "global",
		Priority: 1,
		Label:    "Loan/mortgage/credit mentions",
	})

	// Phone contact spreadsheets
	dorks = append(dorks, Dork{
		Query:    fmt.Sprintf(`"%s" filetype:xlsx OR filetype:csv phone OR contact OR directory`, name),
		Category: "financial",
		Region:   "global",
		Priority: 1,
		Label:    "Phone contact spreadsheets",
	})

	return dorks
}
