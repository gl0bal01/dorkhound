// Package output: CaseBandit export.
//
// Writes a JSON document matching docs/casebandit-bridge.md (schema
// `dorkhound-casebandit-v1`). CaseBandit's import endpoint
// (POST /api/import/dorkhound) ingests this format and produces a fully
// populated Case + Entities + Captures inside the operator's workspace.
package output

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// caseBanditSchemaVersion is the wire-format slug. Bump on incompatible
// changes — see docs/casebandit-bridge.md "Versioning".
const caseBanditSchemaVersion = "dorkhound-casebandit-v1"

type cbDocument struct {
	SchemaVersion string      `json:"schema_version"`
	Generator     cbGenerator `json:"generator"`
	Case          cbCase      `json:"case"`
	Entities      []cbEntity  `json:"entities"`
	Captures      []cbCapture `json:"captures"`
}

type cbGenerator struct {
	Tool        string `json:"tool"`
	Version     string `json:"version"`
	GeneratedAt string `json:"generated_at"`
}

type cbCase struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Tags           []string `json:"tags"`
	Status         string   `json:"status"`
	Notes          string   `json:"notes"`
	ChainOfCustody string   `json:"chainOfCustody"`
}

type cbEntity struct {
	ID         string   `json:"id"`
	CaseID     string   `json:"caseId"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
	Notes      string   `json:"notes"`
	Source     string   `json:"source,omitempty"`
	Tags       []string `json:"tags"`
	CaptureIDs []string `json:"captureIds"`
	Status     string   `json:"status,omitempty"`
	Important  bool     `json:"important,omitempty"`
}

type cbCapture struct {
	ID        string        `json:"id"`
	CaseID    string        `json:"caseId"`
	Timestamp string        `json:"timestamp"`
	URL       string        `json:"url"`
	Title     string        `json:"title"`
	Source    string        `json:"source,omitempty"`
	Type      string        `json:"type"`
	Status    string        `json:"status,omitempty"`
	Tags      []string      `json:"tags"`
	Content   cbCaptureBody `json:"content"`
	EntityIDs []string      `json:"extractedEntities,omitempty"`
}

type cbCaptureBody struct {
	Text string `json:"text"`
}

// CaseBanditExportOptions configures the export. Zero-value is valid.
type CaseBanditExportOptions struct {
	Version     string    // dorkhound version string; "dev" if empty
	Now         time.Time // injectable clock for deterministic tests
	Engine      string    // search engine selected for this run
	GeneratedAt time.Time
}

// CaseBandit writes a v1 bridge document to w.
func CaseBandit(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, opts CaseBanditExportOptions) error {
	if opts.Version == "" {
		opts.Version = "dev"
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	if opts.GeneratedAt.IsZero() {
		opts.GeneratedAt = opts.Now
	}
	if opts.Engine == "" {
		opts.Engine = "google"
	}

	caseID := cbCaseID(c)
	timestamp := opts.GeneratedAt.UTC().Format(time.RFC3339)

	entities, entityByValue := buildEntities(c, caseID)
	captures := buildCaptures(c, dorks, caseID, opts.Engine, timestamp)

	// Cross-link: every capture whose URL/label/text contains an entity
	// value gets that entity in extractedEntities, and the entity gets the
	// capture ID in captureIds. Importer can rely on either direction.
	linkEntitiesAndCaptures(entities, entityByValue, captures)

	doc := cbDocument{
		SchemaVersion: caseBanditSchemaVersion,
		Generator: cbGenerator{
			Tool:        "dorkhound",
			Version:     opts.Version,
			GeneratedAt: timestamp,
		},
		Case: cbCase{
			ID:          caseID,
			Name:        c.Name,
			Description: c.Description,
			Tags:        []string{"dorkhound", "imported"},
			Status:      "active",
			Notes:       buildCaseNotes(c),
		},
		Entities: entities,
		Captures: captures,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func buildEntities(c *caseinfo.Case, caseID string) ([]cbEntity, map[string]int) {
	var entities []cbEntity
	byValue := make(map[string]int)
	add := func(label, etype, notes string, tags ...string) {
		key := etype + ":" + strings.ToLower(strings.TrimSpace(label))
		if _, exists := byValue[key]; exists {
			return
		}
		allTags := append([]string{"dorkhound", "imported"}, tags...)
		ent := cbEntity{
			ID:         cbEntityID(caseID, etype, label),
			CaseID:     caseID,
			Label:      label,
			Type:       etype,
			Notes:      notes,
			Source:     "dorkhound:case-file",
			Tags:       allTags,
			CaptureIDs: []string{},
			Status:     "unconfirmed",
		}
		byValue[key] = len(entities)
		entities = append(entities, ent)
	}

	if strings.TrimSpace(c.Name) != "" {
		notes := "Subject of case (from dorkhound case file)."
		add(c.Name, "person", notes, "subject")
	}
	for _, alias := range c.Aliases {
		add(alias, "username", "Alias of "+c.Name, "alias")
	}
	for _, assoc := range c.Associates {
		add(assoc, "person", "Known associate of "+c.Name, "associate")
	}
	for _, email := range c.Emails {
		add(strings.ToLower(strings.TrimSpace(email)), "email", "")
	}
	for _, phone := range c.Phones {
		add(strings.TrimSpace(phone), "phone", "")
	}
	for _, username := range c.Usernames {
		add(strings.TrimPrefix(strings.TrimSpace(username), "@"), "username", "")
	}
	if strings.TrimSpace(c.Location) != "" {
		add(c.Location, "location", "Last known location")
	}
	if strings.TrimSpace(c.PhotoURL) != "" {
		add(c.PhotoURL, "url", "Photo URL for reverse image search", "photo")
	}
	return entities, byValue
}

func buildCaptures(c *caseinfo.Case, dorks []dork.Dork, caseID, engine, timestamp string) []cbCapture {
	captures := make([]cbCapture, 0, len(dorks))
	for _, d := range dorks {
		url := d.URL(engine)
		tags := []string{
			"dorkhound",
			d.Category,
			"category:" + d.Category,
			"region:" + d.Region,
			"engine:" + engine,
		}
		captures = append(captures, cbCapture{
			ID:        cbCaptureID(d.Category, d.Label, url),
			CaseID:    caseID,
			Timestamp: timestamp,
			URL:       url,
			Title:     d.Category + ": " + d.Label,
			Source:    "dorkhound",
			Type:      "page",
			Tags:      tags,
			Content: cbCaptureBody{
				Text: d.Label + " — " + d.Query,
			},
		})
	}
	return captures
}

// linkEntitiesAndCaptures walks the entity/capture cross-product once and
// records bidirectional links where an entity's literal value appears in a
// capture's URL or query text. The linear pass is fine in practice — even a
// large run is ~500 captures × ~20 entities = 10k comparisons.
func linkEntitiesAndCaptures(entities []cbEntity, byValue map[string]int, captures []cbCapture) {
	if len(entities) == 0 {
		return
	}
	type entRef struct {
		label, etype string
		idx          int
	}
	refs := make([]entRef, 0, len(entities))
	for i, e := range entities {
		v := strings.ToLower(strings.TrimSpace(e.Label))
		if v == "" {
			continue
		}
		refs = append(refs, entRef{label: v, etype: e.Type, idx: i})
	}
	_ = byValue // reserved for future O(1) lookups; intentionally unused
	for ci := range captures {
		hay := strings.ToLower(captures[ci].URL + " " + captures[ci].Content.Text)
		for _, r := range refs {
			if strings.Contains(hay, r.label) {
				captures[ci].EntityIDs = append(captures[ci].EntityIDs, entities[r.idx].ID)
				entities[r.idx].CaptureIDs = append(entities[r.idx].CaptureIDs, captures[ci].ID)
			}
		}
	}
}

func buildCaseNotes(c *caseinfo.Case) string {
	var b strings.Builder
	b.WriteString("## Scope\n\n")
	b.WriteString("Imported from dorkhound — auto-generated lead set.\n\n")
	b.WriteString("## Objectives\n\n")
	if strings.TrimSpace(c.Description) != "" {
		b.WriteString(c.Description)
		b.WriteString("\n\n")
	} else {
		b.WriteString("Triage generated dorks, confirm identifiers, capture evidence.\n\n")
	}
	b.WriteString("## Open Questions\n\n")
	b.WriteString("- Which identifiers are corroborated by multiple independent sources?\n")
	b.WriteString("- Which dead-end leads can be archived?\n")
	return b.String()
}

// cbCaseID produces a stable case ID derived from name/dob/location so reruns
// of the same case file resolve to the same CaseBandit case on reimport.
func cbCaseID(c *caseinfo.Case) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		strings.TrimSpace(c.Name),
		strings.TrimSpace(c.DOB),
		strings.TrimSpace(c.Location),
	}, "\x00")))
	return "dh-" + hex.EncodeToString(sum[:12])
}

func cbEntityID(caseID, etype, label string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		caseID,
		etype,
		strings.ToLower(strings.TrimSpace(label)),
	}, "\x00")))
	return "dh-ent-" + hex.EncodeToString(sum[:12])
}

func cbCaptureID(category, label, url string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		category,
		label,
		url,
	}, "\x00")))
	return "dh-cap-" + hex.EncodeToString(sum[:12])
}
