package output

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// dashboardData is the JSON structure injected into the HTML template.
type dashboardData struct {
	CaseInfo dashboardCase     `json:"case_info"`
	Results  []dashboardResult `json:"results"`
}

type dashboardCase struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int    `json:"age"`
	DOB      string `json:"dob"`
}

type dashboardResult struct {
	Label    string `json:"label"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Priority int    `json:"priority"`
}

// ServeDashboard starts a local HTTP server serving an interactive dashboard.
func ServeDashboard(c *caseinfo.Case, dorks []dork.Dork, engine string, htmlTemplate string) error {
	// Build JSON data blob.
	results := make([]dashboardResult, len(dorks))
	for i, d := range dorks {
		results[i] = dashboardResult{
			Label:    d.Label,
			URL:      d.URL(engine),
			Category: d.Category,
			Priority: d.Priority,
		}
	}

	blob := dashboardData{
		CaseInfo: dashboardCase{
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

	// Replace the placeholder in the HTML template with actual data.
	html := strings.Replace(htmlTemplate, "/*DATA_PLACEHOLDER*/{}", string(jsonBytes), 1)

	// Listen on a random available port on localhost.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("starting listener: %w", err)
	}

	addr := listener.Addr().(*net.TCPAddr)
	url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)

	// Register handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})

	fmt.Fprintf(os.Stderr, "Dashboard running at %s\nPress Ctrl+C to stop.\n", url)

	// Open in browser (best-effort).
	_ = openURL(url)

	// Serve until interrupted.
	return http.Serve(listener, nil)
}
