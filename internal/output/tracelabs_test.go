package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

func TestTraceLabs_Basic(t *testing.T) {
	c := &caseinfo.Case{
		Name:     "Jane Doe",
		Location: "New York, NY",
		Emails:   []string{"jane@example.com"},
	}
	dorks := []dork.Dork{
		{Query: "https://lens.google.com/uploadbyurl?url=X", Category: "image", Priority: 3, Label: "Google Lens reverse search"},
		{Query: `"Jane Doe" site:twitter.com`, Category: "social", Priority: 2, Label: "Twitter search"},
		{Query: `"Jane Doe" site:whitepages.com`, Category: "records", Priority: 1, Label: "Whitepages"},
	}

	var buf bytes.Buffer
	if err := TraceLabs(&buf, c, dorks, "google"); err != nil {
		t.Fatalf("TraceLabs returned error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "# TraceLabs Submission — Jane Doe") {
		t.Errorf("missing header line; got:\n%s", out)
	}
	if !strings.Contains(out, "### image (1)") {
		t.Errorf("missing image category heading; got:\n%s", out)
	}
	if !strings.Contains(out, "### social (1)") {
		t.Errorf("missing social category heading; got:\n%s", out)
	}
	if !strings.Contains(out, "### records (1)") {
		t.Errorf("missing records category heading; got:\n%s", out)
	}
	if !strings.Contains(out, "- [ ] Google Lens reverse search — https://lens.google.com/uploadbyurl?url=X") {
		t.Errorf("missing image dork checkbox line; got:\n%s", out)
	}
	if !strings.Contains(out, "- [ ] Twitter search —") {
		t.Errorf("missing social dork checkbox line; got:\n%s", out)
	}
}

func TestTraceLabs_OmitsEmptyFields(t *testing.T) {
	c := &caseinfo.Case{Name: "Solo Name"}

	var buf bytes.Buffer
	if err := TraceLabs(&buf, c, nil, "google"); err != nil {
		t.Fatalf("TraceLabs returned error: %v", err)
	}
	out := buf.String()

	for _, forbidden := range []string{"**Aliases:**", "**DOB:**", "**Age:**", "**Location:**", "**Description:**", "**Associates:**", "**Emails:**", "**Phones:**", "**Usernames:**", "**Photo:**"} {
		if strings.Contains(out, forbidden) {
			t.Errorf("output should not contain %q but does; got:\n%s", forbidden, out)
		}
	}
	if !strings.Contains(out, "**Name:** Solo Name") {
		t.Errorf("output should contain name field; got:\n%s", out)
	}
}

func TestTraceLabs_CategoryOrdering(t *testing.T) {
	c := &caseinfo.Case{Name: "Test Person"}
	dorks := []dork.Dork{
		{Query: `"Test Person" records`, Category: "records", Priority: 1, Label: "Records dork"},
		{Query: `"Test Person" social`, Category: "social", Priority: 1, Label: "Social dork"},
		{Query: "https://lens.google.com/x", Category: "image", Priority: 3, Label: "Image dork"},
		{Query: `"Test Person" weird`, Category: "weird", Priority: 1, Label: "Weird dork"},
	}

	var buf bytes.Buffer
	if err := TraceLabs(&buf, c, dorks, "google"); err != nil {
		t.Fatalf("TraceLabs returned error: %v", err)
	}
	out := buf.String()

	imagePos := strings.Index(out, "### image")
	socialPos := strings.Index(out, "### social")
	recordsPos := strings.Index(out, "### records")
	weirdPos := strings.Index(out, "### weird")

	if imagePos < 0 || socialPos < 0 || recordsPos < 0 || weirdPos < 0 {
		t.Fatalf("one or more category headings missing; got:\n%s", out)
	}
	if !(imagePos < socialPos) {
		t.Errorf("image (%d) should come before social (%d)", imagePos, socialPos)
	}
	if !(socialPos < recordsPos) {
		t.Errorf("social (%d) should come before records (%d)", socialPos, recordsPos)
	}
	if !(recordsPos < weirdPos) {
		t.Errorf("records (%d) should come before weird (%d)", recordsPos, weirdPos)
	}
}

func TestTraceLabs_UsesEngineURL(t *testing.T) {
	c := &caseinfo.Case{Name: "Bing Person"}
	dorks := []dork.Dork{
		{Query: `"Bing Person" site:twitter.com`, Category: "social", Priority: 2, Label: "Twitter search"},
		{Query: "https://lens.google.com/uploadbyurl?url=X", Category: "image", Priority: 3, Label: "Google Lens reverse search"},
	}

	var buf bytes.Buffer
	if err := TraceLabs(&buf, c, dorks, "bing"); err != nil {
		t.Fatalf("TraceLabs returned error: %v", err)
	}
	out := buf.String()

	// Search dorks should render as bing URLs.
	if !strings.Contains(out, "bing.com/search") {
		t.Errorf("search dork should use bing engine; got:\n%s", out)
	}
	// Direct URLs pass through unchanged.
	if !strings.Contains(out, "https://lens.google.com/uploadbyurl?url=X") {
		t.Errorf("direct URL dork should pass through unchanged; got:\n%s", out)
	}
}
