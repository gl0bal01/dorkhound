package output

import (
	"strings"
	"testing"
)

// TestDashboardRowIDUniqueForSamePrefixURLs locks in the fix for the
// dashboard ID collision: search-engine wrapped URLs share a long base64
// prefix, so prefix-truncated IDs collided. dashboardRowID must yield a
// distinct ID for every (category, label, url, index) tuple.
func TestDashboardRowIDUniqueForSamePrefixURLs(t *testing.T) {
	pairs := []struct{ cat, label, url string }{
		{"social", "Google Jane site:facebook.com", "https://www.google.com/search?q=%22Jane+Doe%22+site%3Afacebook.com"},
		{"social", "Google Jane site:twitter.com", "https://www.google.com/search?q=%22Jane+Doe%22+site%3Atwitter.com"},
		{"social", "Google Jane site:instagram.com", "https://www.google.com/search?q=%22Jane+Doe%22+site%3Ainstagram.com"},
		{"social", "Bing Jane site:facebook.com", "https://www.bing.com/search?q=%22Jane+Doe%22+site%3Afacebook.com"},
	}

	seen := make(map[string]string)
	for i, p := range pairs {
		id := dashboardRowID(p.cat, p.label, p.url, i)
		if existing, ok := seen[id]; ok {
			t.Errorf("ID collision: %q collides with %q (id=%s)", p.url, existing, id)
		}
		seen[id] = p.url
		if !strings.HasPrefix(id, p.cat+":") {
			t.Errorf("ID %q missing category prefix %q", id, p.cat+":")
		}
	}
}

func TestDashboardRowIDStableForSameInput(t *testing.T) {
	a := dashboardRowID("social", "label", "https://example.com", 0)
	b := dashboardRowID("social", "label", "https://example.com", 0)
	if a != b {
		t.Errorf("dashboardRowID not deterministic: %s vs %s", a, b)
	}
}
