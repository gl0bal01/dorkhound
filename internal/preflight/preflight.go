// Package preflight concurrently probes dork URLs to filter out dead direct
// links (404, 410, 5xx, DNS failures, TLS errors) before the operator sees
// them. It only probes dorks whose Query is a direct HTTP(S) URL — search-
// engine-wrapped dorks are passed through untouched, because a HEAD on
// google.com/search always returns 200 and provides no signal about whether
// the underlying search has any hits.
package preflight

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// Options configures a preflight run.
type Options struct {
	Concurrency int           // max parallel HEAD requests (default 8)
	Timeout     time.Duration // per-request timeout (default 5s)
	RateLimit   time.Duration // min delay between requests per worker (default 250ms)
	UserAgent   string        // default "dorkhound-preflight/1.0"
}

// Status records the outcome of probing one dork's direct URL.
type Status struct {
	URL        string
	StatusCode int    // HTTP status code; 0 if request failed
	Err        string // non-empty on transport/TLS/timeout errors
	Dead       bool   // true if the dork should be dropped (4xx/5xx or error)
}

// Report summarizes a preflight pass.
type Report struct {
	Checked int // how many dorks were actually probed (direct URLs only)
	Skipped int // search-engine dorks passed through
	Alive   int
	Dead    int
	Results []Status // one entry per alive checked dork, in input order
}

// Run probes the direct-URL dorks and returns (survivors, report).
// Search-query dorks are always preserved. Direct-URL dorks are dropped when
// HEAD returns 4xx/5xx or transport errors. On context cancellation, all
// pending probes are aborted and partial results returned.
func Run(ctx context.Context, dorks []dork.Dork, opts Options) ([]dork.Dork, Report) {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 8
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 5 * time.Second
	}
	if opts.RateLimit < 0 {
		opts.RateLimit = 0
	}
	if opts.UserAgent == "" {
		opts.UserAgent = "dorkhound-preflight/1.0"
	}

	client := &http.Client{
		Timeout: opts.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	type job struct {
		idx int
		d   dork.Dork
	}

	var report Report
	deadSet := make(map[int]bool)
	statusByIdx := make(map[int]Status)

	var jobs []job
	for i, d := range dorks {
		if isDirectURL(d.Query) {
			jobs = append(jobs, job{i, d})
			report.Checked++
		} else {
			report.Skipped++
		}
	}

	if len(jobs) == 0 {
		return dorks, report
	}

	jobCh := make(chan job)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for w := 0; w < opts.Concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				st := probe(ctx, client, j.d.Query, opts.UserAgent)
				mu.Lock()
				statusByIdx[j.idx] = st
				if st.Dead {
					deadSet[j.idx] = true
				}
				mu.Unlock()
				if opts.RateLimit > 0 {
					select {
					case <-time.After(opts.RateLimit):
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobCh)
		for _, j := range jobs {
			select {
			case jobCh <- j:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()

	survivors := make([]dork.Dork, 0, len(dorks))
	for i, d := range dorks {
		if deadSet[i] {
			report.Dead++
			continue
		}
		if st, ok := statusByIdx[i]; ok {
			report.Alive++
			report.Results = append(report.Results, st)
		}
		survivors = append(survivors, d)
	}

	return survivors, report
}

func isDirectURL(q string) bool {
	return strings.HasPrefix(q, "http://") || strings.HasPrefix(q, "https://")
}

func probe(ctx context.Context, client *http.Client, target, ua string) Status {
	st := Status{URL: target}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target, nil)
	if err != nil {
		st.Err = err.Error()
		st.Dead = true
		return st
	}
	req.Header.Set("User-Agent", ua)
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			st.Err = err.Error()
			st.Dead = true
			return st
		}
		st.Err = err.Error()
		st.Dead = true
		return st
	}
	defer resp.Body.Close()
	st.StatusCode = resp.StatusCode
	if resp.StatusCode == http.StatusMethodNotAllowed {
		reqG, _ := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		reqG.Header.Set("User-Agent", ua)
		respG, errG := client.Do(reqG)
		if errG == nil {
			defer respG.Body.Close()
			st.StatusCode = respG.StatusCode
		}
	}
	if st.StatusCode >= 400 {
		st.Dead = true
	}
	return st
}
