package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

func TestDiscordFormat(t *testing.T) {
	c := &caseinfo.Case{
		Name:     "John Doe",
		Location: "Seattle, WA",
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
		{
			Query:    "https://www.spokeo.com/John-Doe",
			Category: "people-db",
			Region:   "us",
			Priority: 1,
			Label:    "Spokeo",
		},
	}

	var buf bytes.Buffer
	Discord(&buf, c, dorks, "google")
	out := buf.String()

	// Verify header with case name
	if !strings.Contains(out, "## OSINT Results: John Doe") {
		t.Errorf("missing header with name; got:\n%s", out)
	}

	// Verify metadata line with location
	if !strings.Contains(out, "**Location:** Seattle, WA") {
		t.Errorf("missing location metadata; got:\n%s", out)
	}

	// Verify metadata line with age
	if !strings.Contains(out, "**Age:** ~34") {
		t.Errorf("missing age metadata; got:\n%s", out)
	}

	// Verify Social category header
	if !strings.Contains(out, "### Social") {
		t.Errorf("missing Social category header; got:\n%s", out)
	}

	// Verify Records category header
	if !strings.Contains(out, "### Records") {
		t.Errorf("missing Records category header; got:\n%s", out)
	}

	// Verify labels appear
	if !strings.Contains(out, "Facebook profile") {
		t.Errorf("missing Facebook profile label; got:\n%s", out)
	}
	if !strings.Contains(out, "Court records") {
		t.Errorf("missing Court records label; got:\n%s", out)
	}

	// Verify google.com/search URLs for search-query dorks
	if !strings.Contains(out, "google.com/search") {
		t.Errorf("missing google.com/search URL; got:\n%s", out)
	}

	// Verify direct URL is preserved for people-db
	if !strings.Contains(out, "https://www.spokeo.com/John-Doe") {
		t.Errorf("missing direct Spokeo URL; got:\n%s", out)
	}

	// Verify link count in section headers
	if !strings.Contains(out, "(2 links)") {
		t.Errorf("missing '(2 links)' count for social section; got:\n%s", out)
	}
	if !strings.Contains(out, "(1 link)") {
		t.Errorf("missing '(1 link)' count for records section; got:\n%s", out)
	}
}

func TestDiscordMetadataPartialFields(t *testing.T) {
	// Only name and DOB set, no location or age
	c := &caseinfo.Case{
		Name: "Jane Smith",
		DOB:  "1990-05-15",
	}

	dorks := []dork.Dork{
		{
			Query:    `"Jane Smith" site:facebook.com`,
			Category: "social",
			Region:   "global",
			Priority: 1,
			Label:    "Facebook profile",
		},
	}

	var buf bytes.Buffer
	Discord(&buf, c, dorks, "google")
	out := buf.String()

	// Should contain DOB
	if !strings.Contains(out, "**DOB:** 1990-05-15") {
		t.Errorf("missing DOB metadata; got:\n%s", out)
	}

	// Should NOT contain Location or Age since they're empty/zero
	if strings.Contains(out, "**Location:**") {
		t.Errorf("should not contain Location when empty; got:\n%s", out)
	}
	if strings.Contains(out, "**Age:**") {
		t.Errorf("should not contain Age when zero; got:\n%s", out)
	}
}

func TestStdoutFormat(t *testing.T) {
	dorks := []dork.Dork{
		{
			Query:    `"John Doe" site:facebook.com`,
			Category: "social",
			Region:   "global",
			Priority: 2,
			Label:    "Facebook profile",
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
	Stdout(&buf, dorks, "google")
	out := buf.String()

	// Verify category headers
	if !strings.Contains(out, "== social ==") {
		t.Errorf("missing social category header; got:\n%s", out)
	}
	if !strings.Contains(out, "== records ==") {
		t.Errorf("missing records category header; got:\n%s", out)
	}

	// Verify label: URL format
	if !strings.Contains(out, "Facebook profile: ") {
		t.Errorf("missing label line; got:\n%s", out)
	}

	// Verify google URLs
	if !strings.Contains(out, "google.com/search") {
		t.Errorf("missing google.com/search URL; got:\n%s", out)
	}
}

func TestGroupByCategory(t *testing.T) {
	dorks := []dork.Dork{
		{Category: "social", Label: "a"},
		{Category: "social", Label: "b"},
		{Category: "records", Label: "c"},
	}

	groups := groupByCategory(dorks)

	if len(groups["social"]) != 2 {
		t.Errorf("expected 2 social dorks, got %d", len(groups["social"]))
	}
	if len(groups["records"]) != 1 {
		t.Errorf("expected 1 records dork, got %d", len(groups["records"]))
	}
}
