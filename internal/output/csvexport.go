package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// CSV writes dorks to w in CSV format with a header row.
// Columns: label, category, priority, url.
func CSV(w io.Writer, dorks []dork.Dork, engine string) error {
	cw := csv.NewWriter(w)

	// Write header
	if err := cw.Write([]string{"label", "category", "priority", "url"}); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	// Write data rows
	for _, d := range dorks {
		row := []string{
			d.Label,
			d.Category,
			fmt.Sprintf("%d", d.Priority),
			d.URL(engine),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
