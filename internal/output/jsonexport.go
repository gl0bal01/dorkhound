package output

import (
	"encoding/json"
	"io"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// jsonExportSchemaVersion bumps when the JSON export shape changes so
// downstream tooling can detect incompatibilities.
const jsonExportSchemaVersion = 2

// jsonCase mirrors caseinfo.Case for export. Sensitive identifiers are
// included intentionally — JSON is the canonical machine-readable handoff
// format. Operators are warned in README.md to store exports securely.
type jsonCase struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases,omitempty"`
	DOB         string   `json:"dob,omitempty"`
	Age         int      `json:"age,omitempty"`
	Location    string   `json:"location,omitempty"`
	Description string   `json:"description,omitempty"`
	Associates  []string `json:"associates,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	Phones      []string `json:"phones,omitempty"`
	Usernames   []string `json:"usernames,omitempty"`
	PhotoURL    string   `json:"photo_url,omitempty"`
	PhotoPath   string   `json:"photo_path,omitempty"`
	Region      string   `json:"region,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Engine      string   `json:"engine,omitempty"`
}

// jsonResult is the JSON-serializable representation of a single dork result.
type jsonResult struct {
	Label    string `json:"label"`
	Query    string `json:"query"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Region   string `json:"region,omitempty"`
	Priority int    `json:"priority"`
}

// jsonOutput is the top-level JSON output structure.
type jsonOutput struct {
	SchemaVersion int          `json:"schema_version"`
	Engine        string       `json:"engine"`
	Case          jsonCase     `json:"case"`
	Results       []jsonResult `json:"results"`
}

// JSON writes dorks and full case metadata to w in formatted JSON. The
// schema is round-trippable: every caseinfo.Case field is preserved and
// each result includes both the raw query and the rendered engine URL.
func JSON(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, engine string) error {
	results := make([]jsonResult, len(dorks))
	for i, d := range dorks {
		results[i] = jsonResult{
			Label:    d.Label,
			Query:    d.Query,
			URL:      d.URL(engine),
			Category: d.Category,
			Region:   d.Region,
			Priority: d.Priority,
		}
	}

	out := jsonOutput{
		SchemaVersion: jsonExportSchemaVersion,
		Engine:        engine,
		Case: jsonCase{
			Name:        c.Name,
			Aliases:     c.Aliases,
			DOB:         c.DOB,
			Age:         c.Age,
			Location:    c.Location,
			Description: c.Description,
			Associates:  c.Associates,
			Emails:      c.Emails,
			Phones:      c.Phones,
			Usernames:   c.Usernames,
			PhotoURL:    c.PhotoURL,
			PhotoPath:   c.PhotoPath,
			Region:      c.Region,
			Categories:  c.Categories,
			Engine:      c.Engine,
		},
		Results: results,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
