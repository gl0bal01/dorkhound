package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

func TestJSONIncludesFullCase(t *testing.T) {
	c := &caseinfo.Case{
		Name:       "Jane Doe",
		Aliases:    []string{"JD"},
		DOB:        "1990-01-15",
		Age:        34,
		Location:   "Seattle, WA",
		Emails:     []string{"jane@example.com"},
		Phones:     []string{"+1-555-0000"},
		Usernames:  []string{"jdoe42"},
		PhotoURL:   "https://example.com/j.jpg",
		Region:     "us",
		Categories: []string{"username", "email"},
		Engine:     "google",
	}
	dorks := []dork.Dork{
		{Query: `"Jane Doe"`, Category: "social", Region: "global", Priority: 2, Label: "Bare"},
	}

	var buf bytes.Buffer
	if err := JSON(&buf, c, dorks, "duckduckgo"); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if out["schema_version"] == nil {
		t.Error("missing schema_version")
	}
	if got := out["engine"]; got != "duckduckgo" {
		t.Errorf("engine = %v, want duckduckgo", got)
	}

	caseMap, ok := out["case"].(map[string]any)
	if !ok {
		t.Fatalf("case is not a map: %T", out["case"])
	}
	for _, want := range []string{"name", "aliases", "dob", "age", "location", "emails", "phones", "usernames", "photo_url", "region", "categories", "engine"} {
		if _, ok := caseMap[want]; !ok {
			t.Errorf("case missing field %q", want)
		}
	}

	results, ok := out["results"].([]any)
	if !ok || len(results) != 1 {
		t.Fatalf("results malformed: %v", out["results"])
	}
	r := results[0].(map[string]any)
	for _, want := range []string{"label", "query", "url", "category", "priority"} {
		if _, ok := r[want]; !ok {
			t.Errorf("result missing field %q", want)
		}
	}
}
