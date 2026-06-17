package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

func TestCaseBanditExport_SchemaShape(t *testing.T) {
	t.Parallel()
	c := &caseinfo.Case{
		Name:      "Jane Doe",
		Location:  "Seattle, WA",
		DOB:       "1990-01-01",
		Emails:    []string{"jane@example.com"},
		Phones:    []string{"+1-555-0100"},
		Usernames: []string{"@jdoe42"},
		Aliases:   []string{"JD"},
		PhotoURL:  "https://example.com/jane.jpg",
	}
	dorks := []dork.Dork{
		{Query: `"Jane Doe" site:linkedin.com`, Category: "social", Region: "global", Priority: 3, Label: "linkedin exact"},
		{Query: "https://haveibeenpwned.com/account/jane@example.com", Category: "email", Region: "global", Priority: 3, Label: "HIBP"},
	}

	var buf bytes.Buffer
	err := CaseBandit(&buf, c, dorks, CaseBanditExportOptions{
		Version: "test-1.0.0",
		Engine:  "google",
		Now:     time.Date(2026, 5, 27, 10, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CaseBandit: %v", err)
	}

	var doc cbDocument
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshalling generated JSON: %v", err)
	}

	if doc.SchemaVersion != caseBanditSchemaVersion {
		t.Errorf("SchemaVersion = %q, want %q", doc.SchemaVersion, caseBanditSchemaVersion)
	}
	if doc.Generator.Tool != "dorkhound" {
		t.Errorf("Generator.Tool = %q, want dorkhound", doc.Generator.Tool)
	}
	if doc.Generator.Version != "test-1.0.0" {
		t.Errorf("Generator.Version = %q, want test-1.0.0", doc.Generator.Version)
	}
	if doc.Generator.GeneratedAt != "2026-05-27T10:00:00Z" {
		t.Errorf("Generator.GeneratedAt = %q, want 2026-05-27T10:00:00Z", doc.Generator.GeneratedAt)
	}
	if doc.Case.Status != "active" {
		t.Errorf("Case.Status = %q, want active", doc.Case.Status)
	}
	if !strings.HasPrefix(doc.Case.ID, "dh-") {
		t.Errorf("Case.ID = %q, want dh-* prefix", doc.Case.ID)
	}
	if len(doc.Captures) != 2 {
		t.Errorf("Captures = %d, want 2", len(doc.Captures))
	}
}

func TestCaseBanditExport_EntityCoverage(t *testing.T) {
	t.Parallel()
	c := &caseinfo.Case{
		Name:       "Jane Doe",
		Emails:     []string{"jane@example.com", "JANE@EXAMPLE.COM"},
		Phones:     []string{"+1-555-0100"},
		Usernames:  []string{"@jdoe42", "jdoe"},
		Aliases:    []string{"JD"},
		Associates: []string{"John Smith"},
		Location:   "Seattle, WA",
		PhotoURL:   "https://example.com/jane.jpg",
	}
	var buf bytes.Buffer
	err := CaseBandit(&buf, c, nil, CaseBanditExportOptions{})
	if err != nil {
		t.Fatalf("CaseBandit: %v", err)
	}
	var doc cbDocument
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshalling generated JSON: %v", err)
	}

	// Expected entities:
	//   person     × 2 (Jane Doe, John Smith)
	//   username   × 3 (JD alias, jdoe42, jdoe)
	//   email      × 1 (deduped by lower-case)
	//   phone      × 1
	//   location   × 1
	//   url        × 1
	wantTypes := map[string]int{
		"person":   2,
		"username": 3,
		"email":    1,
		"phone":    1,
		"location": 1,
		"url":      1,
	}
	got := make(map[string]int)
	for _, e := range doc.Entities {
		got[e.Type]++
	}
	for typ, want := range wantTypes {
		if got[typ] != want {
			t.Errorf("entity type %q count = %d, want %d (entities=%v)", typ, got[typ], want, got)
		}
	}
}

func TestCaseBanditExport_StableIDs(t *testing.T) {
	t.Parallel()
	c := &caseinfo.Case{
		Name:     "Jane Doe",
		Location: "Seattle, WA",
		DOB:      "1990-01-01",
		Emails:   []string{"jane@example.com"},
	}
	dorks := []dork.Dork{
		{Query: `"Jane Doe"`, Category: "social", Region: "global", Priority: 3, Label: "exact"},
	}
	opts := CaseBanditExportOptions{
		Version: "v1",
		Engine:  "google",
		Now:     time.Date(2026, 5, 27, 10, 0, 0, 0, time.UTC),
	}

	var bufA, bufB bytes.Buffer
	if err := CaseBandit(&bufA, c, dorks, opts); err != nil {
		t.Fatalf("CaseBandit A: %v", err)
	}
	if err := CaseBandit(&bufB, c, dorks, opts); err != nil {
		t.Fatalf("CaseBandit B: %v", err)
	}
	if !bytes.Equal(bufA.Bytes(), bufB.Bytes()) {
		t.Error("two runs with the same input produced different output; IDs must be stable")
	}
}

func TestCaseBanditExport_LinksEntitiesToCaptures(t *testing.T) {
	t.Parallel()
	c := &caseinfo.Case{
		Name:   "Jane Doe",
		Emails: []string{"jane@example.com"},
	}
	dorks := []dork.Dork{
		{Query: `"jane@example.com" site:github.com`, Category: "email", Region: "global", Priority: 3, Label: "email on github"},
	}
	var buf bytes.Buffer
	if err := CaseBandit(&buf, c, dorks, CaseBanditExportOptions{Engine: "google"}); err != nil {
		t.Fatalf("CaseBandit: %v", err)
	}
	var doc cbDocument
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshalling: %v", err)
	}

	var emailEntID string
	for _, e := range doc.Entities {
		if e.Type == "email" {
			emailEntID = e.ID
			if len(e.CaptureIDs) == 0 {
				t.Error("email entity has no CaptureIDs but matching capture exists")
			}
		}
	}
	if emailEntID == "" {
		t.Fatal("no email entity emitted")
	}
	found := false
	for _, cap := range doc.Captures {
		for _, eid := range cap.EntityIDs {
			if eid == emailEntID {
				found = true
			}
		}
	}
	if !found {
		t.Error("no capture references the email entity")
	}
}
