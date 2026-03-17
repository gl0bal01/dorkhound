package dork

import (
	"net/url"
	"strings"
)

// Dork represents a single search query (dork) to be executed against a search engine.
type Dork struct {
	Query    string // the search query string
	Category string // social, records, financial, location, forums, people-db
	Region   string // country code (us, ca, uk, au, ru, fr, de, at, nl) or "global"
	Priority int    // 1-3, higher = opens first
	Label    string // human-readable label
}

// Engines maps engine names to their search URL base.
var Engines = map[string]string{
	"google":     "https://www.google.com/search?q=",
	"bing":       "https://www.bing.com/search?q=",
	"duckduckgo": "https://duckduckgo.com/?q=",
	"yandex":     "https://yandex.com/search/?text=",
}

// URL returns the full search URL for this dork on the given engine.
// If the Query already starts with http:// or https://, it is returned as-is
// (used for people-db direct links). Otherwise, the Query is URL-encoded and
// prepended with the engine's base URL. Unknown engines default to google.
func (d Dork) URL(engine string) string {
	if strings.HasPrefix(d.Query, "http://") || strings.HasPrefix(d.Query, "https://") {
		return d.Query
	}

	base, ok := Engines[engine]
	if !ok {
		base = Engines["google"]
	}

	return base + url.QueryEscape(d.Query)
}
