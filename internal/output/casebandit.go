// Package output: CaseBandit export.
//
// Writes a JSON document matching docs/casebandit-bridge.md (schema
// `dorkhound-casebandit-v1`). CaseBandit's import endpoint
// (POST /api/import/dorkhound) ingests this format and produces a fully
// populated Case + Entities + Captures inside the operator's workspace.
package output

import (
	"encoding/json"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// minLinkableLabel is the floor for entity-label matching against capture
// content. Labels shorter than this (e.g. "JD") would create too many
// incidental cross-links. Two-character labels are still emitted as
// entities — they just aren't auto-linked to captures.
const minLinkableLabel = 3

// caseBanditSchemaVersion is the wire-format slug. Bump on incompatible
// changes — see docs/casebandit-bridge.md "Versioning".
const caseBanditSchemaVersion = "dorkhound-casebandit-v1"

// Entity type slugs mirroring CaseBandit's ENTITY_TYPES tuple in
// shared/types.ts. Keep this block in sync with that tuple — a wire-format
// drift here silently breaks CaseBandit's importer. The Go side currently
// emits a subset (no organization/ip/hash/credential/other) — those are
// reserved for future producers.
const (
	entTypePerson   = "person"
	entTypeUsername = "username"
	entTypeURL      = "url"
	entTypeEmail    = "email"
	entTypePhone    = "phone"
	entTypeLocation = "location"
)

// Capture/Entity status slugs mirroring CaseBandit's EntityStatus.
const entStatusUnconfirmed = "unconfirmed"

// Case status slug.
const caseStatusActive = "active"

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
			Status:      caseStatusActive,
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
	name := strings.TrimSpace(c.Name)

	aliasNote := func(kind string) string {
		if name != "" {
			return kind + " of " + name
		}
		return kind + " (case subject not named)"
	}

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
			Status:     entStatusUnconfirmed,
		})
	}

	if name != "" {
		add(name, entTypePerson, "Subject of case (from dorkhound case file).", "subject")
	}
	for _, alias := range c.Aliases {
		add(alias, entTypeUsername, aliasNote("Alias"), "alias")
	}
	for _, assoc := range c.Associates {
		add(assoc, entTypePerson, aliasNote("Known associate"), "associate")
	}
	for _, email := range c.Emails {
		add(strings.ToLower(email), entTypeEmail, "")
	}
	for _, phone := range c.Phones {
		add(phone, entTypePhone, "")
	}
	for _, username := range c.Usernames {
		add(strings.TrimPrefix(strings.TrimSpace(username), "@"), entTypeUsername, "")
	}
	if loc := strings.TrimSpace(c.Location); loc != "" {
		add(loc, entTypeLocation, "Last known location")
	}
	if photo := strings.TrimSpace(c.PhotoURL); photo != "" {
		add(photo, entTypeURL, "Photo URL for reverse image search", "photo")
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

// linkEntitiesAndCaptures records bidirectional links where every
// alphanumeric token (≥ minLinkableLabel chars) of an entity's label
// appears in a capture's URL or query text. Matching is case-insensitive
// token-set lookup — the capture haystack is tokenized once, then each
// entity is an O(label-tokens) lookup. Compound labels like
// "jane@example.com" tokenize to {jane, example, com} and match only when
// all three appear; single-token labels like "jdoe42" match when the
// token is in the haystack. Sub-floor tokens like "jd" are excluded so
// they can't drag a compound label through on a single weak hit.
//
// Each (entity, capture) pair is visited at most once per capture, so no
// per-link dedupe is needed for EntityIDs. CaptureIDs uses slices.Contains
// because the same capture ID can repeat across iterations only via
// hash-collision-truncation, which is allowed by the encoding but
// vanishingly rare; the guard is cheap and removes the theoretical bug.
func linkEntitiesAndCaptures(entities []cbEntity, captures []cbCapture) {
	if len(entities) == 0 || len(captures) == 0 {
		return
	}
	type entRef struct {
		tokens []string
		idx    int
	}
	var refs []entRef
	for i, e := range entities {
		toks := significantTokens(e.Label)
		if len(toks) == 0 {
			continue
		}
		refs = append(refs, entRef{tokens: toks, idx: i})
	}
	for ci := range captures {
		hay := tokenize(captures[ci].URL + " " + captures[ci].Content.Text)
		for _, r := range refs {
			if !allInSet(r.tokens, hay) {
				continue
			}
			captures[ci].EntityIDs = append(captures[ci].EntityIDs, entities[r.idx].ID)
			if !slices.Contains(entities[r.idx].CaptureIDs, captures[ci].ID) {
				entities[r.idx].CaptureIDs = append(entities[r.idx].CaptureIDs, captures[ci].ID)
			}
		}
	}
}

// significantTokens returns the ≥ minLinkableLabel alphanumeric tokens of
// label, lowercased. Returns nil when no token meets the floor — those
// labels are emitted as entities but skipped for linking.
func significantTokens(label string) []string {
	all := tokenize(label)
	var out []string
	for t := range all {
		if len(t) >= minLinkableLabel {
			out = append(out, t)
		}
	}
	return out
}

func allInSet(tokens []string, set map[string]bool) bool {
	for _, t := range tokens {
		if !set[t] {
			return false
		}
	}
	return true
}

// tokenize splits s on every byte outside [a-zA-Z0-9] and returns the set of
// lowercased non-empty tokens. Used as a one-shot index per capture so the
// per-entity lookup in linkEntitiesAndCaptures is O(1). ASCII-only by
// design — matches the labels we produce (lowercased emails, trimmed
// usernames, locations).
func tokenize(s string) map[string]bool {
	out := make(map[string]bool)
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		out[strings.ToLower(b.String())] = true
		b.Reset()
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			b.WriteByte(c)
			continue
		}
		flush()
	}
	flush()
	return out
}

// cbCaseID produces a stable case ID derived from name/dob/location so reruns
// of the same case file resolve to the same CaseBandit case on reimport.
func cbCaseID(c *caseinfo.Case) string {
	return "dh-" + stableID(12,
		strings.TrimSpace(c.Name),
		strings.TrimSpace(c.DOB),
		strings.TrimSpace(c.Location),
	)
}

func cbEntityID(caseID, etype, label string) string {
	return "dh-ent-" + stableID(12, caseID, etype, strings.ToLower(strings.TrimSpace(label)))
}

func cbCaptureID(category, label, url string) string {
	return "dh-cap-" + stableID(12, category, label, url)
}
