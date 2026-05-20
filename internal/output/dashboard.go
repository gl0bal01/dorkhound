package output

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// dashboardData is the JSON structure injected into the HTML template.
type dashboardData struct {
	CaseInfo dashboardCase     `json:"case_info"`
	Results  []dashboardResult `json:"results"`
}

type dashboardCase struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int    `json:"age"`
	DOB      string `json:"dob"`
}

type dashboardResult struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Region   string `json:"region"`
	Priority int    `json:"priority"`
}

// ServeDashboard starts a local HTTP server serving an interactive dashboard.
func ServeDashboard(c *caseinfo.Case, dorks []dork.Dork, engine string, htmlTemplate string) error {
	// Build JSON data blob. Each row gets a collision-resistant ID derived
	// from category + label + URL + position so duplicate search-engine
	// URL prefixes never collide in the DOM.
	results := make([]dashboardResult, len(dorks))
	for i, d := range dorks {
		url := d.URL(engine)
		results[i] = dashboardResult{
			ID:       dashboardRowID(d.Category, d.Label, url, i),
			Label:    d.Label,
			URL:      url,
			Category: d.Category,
			Region:   d.Region,
			Priority: d.Priority,
		}
	}

	blob := dashboardData{
		CaseInfo: dashboardCase{
			ID:       dashboardCaseID(c),
			Name:     c.Name,
			Location: c.Location,
			Age:      c.Age,
			DOB:      c.DOB,
		},
		Results: results,
	}

	jsonBytes, err := json.Marshal(blob)
	if err != nil {
		return fmt.Errorf("marshalling dashboard data: %w", err)
	}

	token, err := randomURLToken()
	if err != nil {
		return fmt.Errorf("generating dashboard token: %w", err)
	}
	nonce, err := randomCSPNonce()
	if err != nil {
		return fmt.Errorf("generating CSP nonce: %w", err)
	}

	// Replace the placeholder in the HTML template with actual data.
	html := strings.Replace(htmlTemplate, "/*DATA_PLACEHOLDER*/{}", string(jsonBytes), 1)
	html = strings.ReplaceAll(html, "CSP_NONCE", nonce)

	// Listen on a random available port on localhost.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("starting listener: %w", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	dashboardPath := "/" + token
	url := fmt.Sprintf("http://127.0.0.1:%d%s", addr.Port, dashboardPath)

	// Use a dedicated mux instead of http.DefaultServeMux to avoid
	// route pollution from imported packages.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only serve the unguessable dashboard path.
		if r.URL.Path != dashboardPath {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csp := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s'; style-src 'nonce-%s'; object-src 'none'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'", nonce, nonce)
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, html)
	})

	fmt.Fprintf(os.Stderr, "Dashboard running at %s\nPress Ctrl+C to stop.\n", url)

	// Open in browser (best-effort).
	_ = openURL(url)

	// Serve until interrupted, with timeouts to prevent slowloris.
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	return srv.Serve(listener)
}

func randomURLToken() (string, error) {
	var b [18]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}

func randomCSPNonce() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b[:]), nil
}

// dashboardRowID returns a stable, collision-resistant ID for a dashboard
// row. The category prefix is sanitized so a slug containing characters
// invalid in CSS selectors (`:`, `.`, space) can never break querySelector
// calls in the dashboard JS.
func dashboardRowID(category, label, url string, index int) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		category,
		label,
		url,
		fmt.Sprintf("%d", index),
	}, "\x00")))
	return sanitizeIDPrefix(category) + ":" + hex.EncodeToString(sum[:12])
}

// sanitizeIDPrefix replaces every byte outside [A-Za-z0-9_-] with '_' so
// the prefix is safe to use in DOM IDs and CSS selectors.
func sanitizeIDPrefix(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '_':
			b.WriteByte(c)
		default:
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "cat"
	}
	return b.String()
}

func dashboardCaseID(c *caseinfo.Case) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		strings.TrimSpace(c.Name),
		strings.TrimSpace(c.DOB),
		strings.TrimSpace(c.Location),
	}, "\x00")))
	return hex.EncodeToString(sum[:16])
}
