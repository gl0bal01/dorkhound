package preflight

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/moved", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/alive", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/gone", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
	})
	mux.HandleFunc("/servererror", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/headnot", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return httptest.NewServer(mux)
}

func TestRun(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	defer srv.Close()

	dorks := []dork.Dork{
		{Query: srv.URL + "/alive", Label: "alive"},
		{Query: srv.URL + "/moved", Label: "moved"},
		{Query: srv.URL + "/gone", Label: "gone"},
		{Query: srv.URL + "/servererror", Label: "servererror"},
		{Query: srv.URL + "/headnot", Label: "headnot"},
		{Query: `"jane doe" site:example.com`, Label: "search1"},
		{Query: `jane doe facebook`, Label: "search2"},
	}

	survivors, report := Run(context.Background(), dorks, Options{
		Timeout:             500 * time.Millisecond,
		RateLimit:           0,
		AllowPrivateNetwork: true,
	})

	if report.Checked != 5 {
		t.Errorf("Checked = %d, want 5", report.Checked)
	}
	if report.Skipped != 2 {
		t.Errorf("Skipped = %d, want 2", report.Skipped)
	}
	// alive, moved, headnot survive; gone and servererror are dead
	if report.Alive != 3 {
		t.Errorf("Alive = %d, want 3", report.Alive)
	}
	if report.Dead != 2 {
		t.Errorf("Dead = %d, want 2", report.Dead)
	}

	// Search dorks must always be in survivors.
	foundSearch := 0
	for _, s := range survivors {
		if s.Label == "search1" || s.Label == "search2" {
			foundSearch++
		}
	}
	if foundSearch != 2 {
		t.Errorf("search dorks in survivors = %d, want 2", foundSearch)
	}

	// gone and servererror must not be in survivors.
	for _, s := range survivors {
		if s.Label == "gone" || s.Label == "servererror" {
			t.Errorf("dead dork %q should not be in survivors", s.Label)
		}
	}

	// headnot must survive (GET fallback).
	foundHeadnot := false
	for _, s := range survivors {
		if s.Label == "headnot" {
			foundHeadnot = true
		}
	}
	if !foundHeadnot {
		t.Error("headnot dork should survive via GET fallback")
	}
}

func TestRunSlowTimeout(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	defer srv.Close()

	dorks := []dork.Dork{
		{Query: srv.URL + "/slow", Label: "slow"},
	}

	start := time.Now()
	survivors, report := Run(context.Background(), dorks, Options{
		Timeout:             100 * time.Millisecond,
		RateLimit:           0,
		AllowPrivateNetwork: true,
	})
	elapsed := time.Since(start)

	if elapsed > time.Second {
		t.Errorf("Run took %v, expected < 1s with 100ms timeout", elapsed)
	}
	if report.Dead != 1 {
		t.Errorf("Dead = %d, want 1 (slow should timeout)", report.Dead)
	}
	if len(survivors) != 0 {
		t.Errorf("survivors = %d, want 0", len(survivors))
	}
}

func TestRunContextCancellation(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	defer srv.Close()

	dorks := []dork.Dork{
		{Query: srv.URL + "/slow", Label: "slow"},
		{Query: srv.URL + "/slow", Label: "slow2"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	// Must not panic.
	_, _ = Run(ctx, dorks, Options{
		Timeout:             5 * time.Second,
		RateLimit:           0,
		AllowPrivateNetwork: true,
	})
}

func TestRunNoDirectURLs(t *testing.T) {
	t.Parallel()
	dorks := []dork.Dork{
		{Query: `"john doe" site:facebook.com`, Label: "s1"},
		{Query: `john doe linkedin`, Label: "s2"},
	}
	survivors, report := Run(context.Background(), dorks, Options{})
	if report.Checked != 0 {
		t.Errorf("Checked = %d, want 0", report.Checked)
	}
	if report.Skipped != 2 {
		t.Errorf("Skipped = %d, want 2", report.Skipped)
	}
	if len(survivors) != 2 {
		t.Errorf("survivors = %d, want 2", len(survivors))
	}
}

func TestRunBlocksPrivateNetworkTargets(t *testing.T) {
	t.Parallel()

	dorks := []dork.Dork{
		{Query: "http://127.0.0.1:8080/alive", Label: "loopback"},
		{Query: "http://169.254.169.254/latest/meta-data/", Label: "metadata"},
		{Query: "http://localhost:8080/alive", Label: "localhost"},
	}

	survivors, report := Run(context.Background(), dorks, Options{
		Timeout:   100 * time.Millisecond,
		RateLimit: 0,
	})

	if report.Checked != 3 {
		t.Errorf("Checked = %d, want 3", report.Checked)
	}
	// SSRF-refused targets must be Blocked, not Dead, so operators can tell
	// "site refused with 404" apart from "we never let the dial happen".
	if report.Blocked != 3 {
		t.Errorf("Blocked = %d, want 3", report.Blocked)
	}
	if report.Dead != 0 {
		t.Errorf("Dead = %d, want 0 (SSRF refusals must be Blocked, not Dead)", report.Dead)
	}
	if len(survivors) != 0 {
		t.Errorf("survivors = %d, want 0", len(survivors))
	}
	for _, st := range report.Results {
		if !st.Blocked() {
			t.Errorf("Status.Blocked() = false for %s, want true", st.URL)
		}
	}
}

func TestValidateProbeTargetAllowsPrivateWhenOptedIn(t *testing.T) {
	t.Parallel()

	if err := validateProbeTarget(context.Background(), "http://127.0.0.1:8080/alive", true); err != nil {
		t.Fatalf("validateProbeTarget with AllowPrivateNetwork = %v, want nil", err)
	}
}

func TestIsDirectURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		q    string
		want bool
	}{
		{"http://example.com", true},
		{"https://example.com/path?q=1", true},
		{`"jane doe" site:x.com`, false},
		{"", false},
		{"example.com", false},
		{"ftp://example.com", false},
	}
	for _, tc := range cases {
		if got := isDirectURL(tc.q); got != tc.want {
			t.Errorf("isDirectURL(%q) = %v, want %v", tc.q, got, tc.want)
		}
	}
}
