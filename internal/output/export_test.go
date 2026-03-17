package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

func TestJSONFormat(t *testing.T) {
	c := &caseinfo.Case{
		Name:     "John Doe",
		Location: "Seattle",
		Age:      34,
	}

	dorks := []dork.Dork{
		{
			Query:    `"John Doe" site:facebook.com`,
			Category: "social",
			Region:   "global",
			Priority: 2,
			Label:    "Facebook profile",
		},
		{
			Query:    `"John Doe" site:twitter.com`,
			Category: "social",
			Region:   "global",
			Priority: 1,
			Label:    "Twitter profile",
		},
		{
			Query:    `"John Doe" site:courtlistener.com`,
			Category: "records",
			Region:   "us",
			Priority: 1,
			Label:    "Court records",
		},
	}

	var buf bytes.Buffer
	if err := JSON(&buf, c, dorks, "google"); err != nil {
		t.Fatalf("JSON() returned error: %v", err)
	}

	out := buf.String()

	// Verify output is valid JSON
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, out)
	}

	// Verify output contains the case name
	if !strings.Contains(out, "John Doe") {
		t.Errorf("output missing case name 'John Doe'; got:\n%s", out)
	}

	// Verify the case section
	var caseSection struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Age      int    `json:"age"`
	}
	if err := json.Unmarshal(parsed["case"], &caseSection); err != nil {
		t.Fatalf("failed to unmarshal case section: %v", err)
	}
	if caseSection.Name != "John Doe" {
		t.Errorf("case.name = %q, want %q", caseSection.Name, "John Doe")
	}
	if caseSection.Location != "Seattle" {
		t.Errorf("case.location = %q, want %q", caseSection.Location, "Seattle")
	}
	if caseSection.Age != 34 {
		t.Errorf("case.age = %d, want %d", caseSection.Age, 34)
	}

	// Verify results count
	var results []struct {
		Label    string `json:"label"`
		URL      string `json:"url"`
		Category string `json:"category"`
		Priority int    `json:"priority"`
	}
	if err := json.Unmarshal(parsed["results"], &results); err != nil {
		t.Fatalf("failed to unmarshal results: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}

	// Verify each result has label, url, category, priority
	for i, r := range results {
		if r.Label == "" {
			t.Errorf("results[%d].label is empty", i)
		}
		if r.URL == "" {
			t.Errorf("results[%d].url is empty", i)
		}
		if r.Category == "" {
			t.Errorf("results[%d].category is empty", i)
		}
		if r.Priority == 0 {
			t.Errorf("results[%d].priority is zero", i)
		}
	}

	// Verify URLs are google URLs
	for i, r := range results {
		if !strings.Contains(r.URL, "google.com/search") {
			t.Errorf("results[%d].url missing google.com/search; got %q", i, r.URL)
		}
	}
}

func TestCSVFormat(t *testing.T) {
	dorks := []dork.Dork{
		{
			Query:    `"John Doe" site:facebook.com`,
			Category: "social",
			Region:   "global",
			Priority: 2,
			Label:    "Facebook profile",
		},
		{
			Query:    `"John Doe" site:twitter.com`,
			Category: "social",
			Region:   "global",
			Priority: 1,
			Label:    "Twitter profile",
		},
		{
			Query:    `"John Doe" site:courtlistener.com`,
			Category: "records",
			Region:   "us",
			Priority: 1,
			Label:    "Court records",
		},
	}

	var buf bytes.Buffer
	if err := CSV(&buf, dorks, "google"); err != nil {
		t.Fatalf("CSV() returned error: %v", err)
	}

	out := buf.String()

	// Parse CSV output
	reader := csv.NewReader(strings.NewReader(out))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v\noutput:\n%s", err, out)
	}

	// Verify header row
	if len(records) == 0 {
		t.Fatal("CSV output has no rows")
	}
	header := records[0]
	expectedHeader := []string{"label", "category", "priority", "url"}
	if len(header) != len(expectedHeader) {
		t.Fatalf("header has %d columns, want %d", len(header), len(expectedHeader))
	}
	for i, h := range expectedHeader {
		if header[i] != h {
			t.Errorf("header[%d] = %q, want %q", i, header[i], h)
		}
	}

	// Verify correct number of data rows (header + 3 dorks)
	if len(records) != 4 {
		t.Errorf("expected 4 rows (1 header + 3 data), got %d", len(records))
	}

	// Verify first data row
	row1 := records[1]
	if row1[0] != "Facebook profile" {
		t.Errorf("row1 label = %q, want %q", row1[0], "Facebook profile")
	}
	if row1[1] != "social" {
		t.Errorf("row1 category = %q, want %q", row1[1], "social")
	}
	if row1[2] != "2" {
		t.Errorf("row1 priority = %q, want %q", row1[2], "2")
	}
	if !strings.Contains(row1[3], "google.com/search") {
		t.Errorf("row1 url missing google.com/search; got %q", row1[3])
	}
}
