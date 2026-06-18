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
	"regexp"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// minLinkableLabel is the floor for entity-label substring matching against
// capture content. Labels shorter than this (e.g. "JD") would create too many
// incidental cross-links inside URLs and English prose. Two-character labels
// are still emitted as entities — they just aren't auto-linked to captures.
const minLinkableLabel = 3

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
	ID        string `json:"id"`
	CaseID    string `json:"caseId"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Source    string `json:"source,omitempty"`
	Type      string `json:"type"`
	// Status mirrors CaseBandit's Capture.status traffic-light enum
	// (`blue` | `green` | `yellow` | `red`). The writer intentionally
	// leaves it unset — dorkhound generates leads, not graded captures.
	// CaseBandit's importer assigns the initial status.
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
	Engine      string    // search engine; "google" if empty
	GeneratedAt time.Time // timestamp stamped on the doc; time.Now().UTC() if zero
}

// CaseBandit writes a v1 bridge document to w.
func CaseBandit(w io.Writer, c *caseinfo.Case, dorks []dork.Dork, opts CaseBanditExportOptions) error {
	if opts.Version == "" {
		opts.Version = "dev"
	}
	if opts.GeneratedAt.IsZero() {
		opts.GeneratedAt = time.Now().UTC()
	}
	if opts.Engine == "" {
		opts.Engine = "google"
	}

	caseID := cbCaseID(c)
	timestamp := opts.GeneratedAt.UTC().Format(time.RFC3339)

	entities := buildEntities(c, caseID)
	captures := buildCaptures(c, dorks, caseID, opts.Engine, timestamp)

	// Cross-link: every capture whose URL/label/text contains an entity
	// value gets that entity in extractedEntities, and the entity gets the
	// capture ID in captureIds. Importer can rely on either direction.
	linkEntitiesAndCaptures(entities, captures)

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
			Notes:       c.Description,
		},
		Entities: entities,
		Captures: captures,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func buildEntities(c *caseinfo.Case, caseID string) []cbEntity {
	var entities []cbEntity
	seen := make(map[string]bool)
	add := func(label, etype, notes string, tags ...string) {
		label = strings.TrimSpace(label)
		if label == "" {
			return
		}
		key := etype + ":" + strings.ToLower(label)
		if seen[key] {
			return
		}
		seen[key] = true
		allTags := append([]string{"dorkhound", "imported"}, tags...)
		entities = append(entities, cbEntity{
			ID:         cbEntityID(caseID, etype, label),
			CaseID:     caseID,
			Label:      label,
			Type:       etype,
			Notes:      notes,
			Source:     "dorkhound:case-file",
			Tags:       allTags,
			CaptureIDs: []string{},
			Status:     "unconfirmed",
		})
	}

	// aliasNote / associateNote produce "Alias of <Name>" only when Name is set,
	// otherwise a clean fallback so the field doesn't read "Alias of ".
	aliasNote := func(kind string) string {
		if name := strings.TrimSpace(c.Name); name != "" {
			return kind + " of " + name
		}
		return kind + " (case subject not named)"
	}

	if name := strings.TrimSpace(c.Name); name != "" {
		add(name, "person", "Subject of case (from dorkhound case file).", "subject")
	}
	for _, alias := range c.Aliases {
		add(alias, "username", aliasNote("Alias"), "alias")
	}
	for _, assoc := range c.Associates {
		add(assoc, "person", aliasNote("Known associate"), "associate")
	}
	for _, email := range c.Emails {
		add(strings.ToLower(email), "email", "")
	}
	for _, phone := range c.Phones {
		add(phone, "phone", "")
	}
	for _, username := range c.Usernames {
		add(strings.TrimPrefix(strings.TrimSpace(username), "@"), "username", "")
	}
	if loc := strings.TrimSpace(c.Location); loc != "" {
		add(loc, "location", "Last known location")
	}
	if url := strings.TrimSpace(c.PhotoURL); url != "" {
		add(url, "url", "Photo URL for reverse image search", "photo")
	}
	return entities
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
// records bidirectional links where an entity's value matches a word in a
// capture's URL or query text. Matching is word-boundary, case-insensitive, and
// skips labels shorter than minLinkableLabel to avoid incidental matches like
// `"JD"` hitting `"adjusted"` or `"jdoe"` hitting `"jdoe42"`. The linear pass
// is fine in practice: a large run is ~500 captures × ~20 entities = 10k
// comparisons.
//
// Email addresses are split-and-matched on the local-part too, so an email
// entity links to a capture that searched for just the local part — that's
// the actual dork pattern dorks_email.go generates.
func linkEntitiesAndCaptures(entities []cbEntity, captures []cbCapture) {
	if len(entities) == 0 || len(captures) == 0 {
		return
	}
	type entRef struct {
		re  *regexp.Regexp
		idx int
	}
	var refs []entRef
	for i, e := range entities {
		label := strings.ToLower(strings.TrimSpace(e.Label))
		if len(label) < minLinkableLabel {
			continue
		}
		re, err := regexp.Compile(`(?i)\b` + regexp.QuoteMeta(label) + `\b`)
		if err != nil {
			continue
		}
		refs = append(refs, entRef{re: re, idx: i})
	}
	for ci := range captures {
		hay := captures[ci].URL + " " + captures[ci].Content.Text
		entSeen := make(map[string]bool)
		for _, r := range refs {
			if !r.re.MatchString(hay) {
				continue
			}
			entID := entities[r.idx].ID
			if entSeen[entID] {
				continue
			}
			entSeen[entID] = true
			captures[ci].EntityIDs = append(captures[ci].EntityIDs, entID)
			entities[r.idx].CaptureIDs = appendUnique(entities[r.idx].CaptureIDs, captures[ci].ID)
		}
	}
}

func appendUnique(s []string, v string) []string {
	for _, e := range s {
		if e == v {
			return s
		}
	}
	return append(s, v)
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
