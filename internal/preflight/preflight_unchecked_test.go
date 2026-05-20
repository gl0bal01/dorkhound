package preflight

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// TestRunCancellationCountsUnchecked verifies that direct-URL dorks abandoned
// by context cancellation are reported as Unchecked, not silently flagged
// alive or dropped.
func TestRunCancellationCountsUnchecked(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	dorks := []dork.Dork{
		{Query: srv.URL + "/slow", Label: "slow1"},
		{Query: srv.URL + "/slow", Label: "slow2"},
		{Query: `"jane" site:example.com`, Label: "search"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancel

	survivors, report := Run(ctx, dorks, Options{
		Timeout:             5 * time.Second,
		RateLimit:           0,
		AllowPrivateNetwork: true,
	})

	if report.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1", report.Skipped)
	}
	// Two direct URLs: each is either alive, dead, or unchecked.
	if report.Alive+report.Dead+report.Unchecked != 2 {
		t.Errorf("alive=%d dead=%d unchecked=%d, sum=%d, want 2",
			report.Alive, report.Dead, report.Unchecked, report.Alive+report.Dead+report.Unchecked)
	}
	// Search dork must always survive.
	foundSearch := false
	for _, s := range survivors {
		if s.Label == "search" {
			foundSearch = true
		}
	}
	if !foundSearch {
		t.Error("search dork must survive cancellation")
	}
}
