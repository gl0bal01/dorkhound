package output

import (
	"encoding/json"
	"io"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// jsonCase is the JSON-serializable representation of case metadata.
type jsonCase struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int    `json:"age"`
	DOB      string `json:"dob"`
}

// jsonResult is the JSON-serializable representation of a single dork result.
type jsonResult struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Priority int    `json:"priority"`
}

// jsonOutput is the top-level JSON output structure.
type jsonOutput struct {
	Case    jsonCase     `json:"case"`
	Results []jsonResult `json:"results"`
}

// JSON writes dorks and case info to w in formatted JSON.
func JSON(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, engine string) error {
	results := make([]jsonResult, len(dorks))
	for i, d := range dorks {
		results[i] = jsonResult{
			Label:    d.Label,
			URL:      d.URL(engine),
			Category: d.Category,
			Priority: d.Priority,
		}
	}

	out := jsonOutput{
		Case: jsonCase{
			Name:     c.Name,
			Location: c.Location,
			Age:      c.Age,
			DOB:      c.DOB,
		},
		Results: results,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
