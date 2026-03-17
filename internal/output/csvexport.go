package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// CSV writes dorks to w in CSV format with a header row.
// Columns: label, category, region, priority, query, url.
func CSV(w io.Writer, dorks []dork.Dork, engine string) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"label", "category", "region", "priority", "query", "url"}); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	for _, d := range dorks {
		row := []string{
			d.Label,
			d.Category,
			d.Region,
			fmt.Sprintf("%d", d.Priority),
			d.Query,
			d.URL(engine),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
