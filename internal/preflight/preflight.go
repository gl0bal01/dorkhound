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
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// ErrBlockedTarget is returned by the SSRF guard (both validateProbeTarget and
// dialValidated) when a probe would hit a private/loopback/link-local/multicast
// address. Check with errors.Is so wrapping by http.Client/url.Error is
// transparent and resistant to future format changes in Go's stdlib.
var ErrBlockedTarget = errors.New("blocked private network target")

// Options configures a preflight run.
type Options struct {
	Concurrency         int           // max parallel HEAD requests (default 8)
	Timeout             time.Duration // per-request timeout (default 5s)
	RateLimit           time.Duration // min delay between requests per worker (default 250ms)
	UserAgent           string        // default "dorkhound-preflight/1.0"
	AllowPrivateNetwork bool          // allow loopback/private/link-local targets
}

// Disposition is the outcome of probing one dork's direct URL. Exactly one
// value applies to any given Status. The iota typing prevents the "multiple
// bool flags accidentally set true" footgun the previous 4-bool encoding
// allowed.
type Disposition int

const (
	DispositionAlive     Disposition = iota // 2xx/3xx response
	DispositionDead                         // 4xx/5xx response or transport/DNS/TLS error
	DispositionBlocked                      // SSRF guard refused (private/loopback/link-local target)
	DispositionUnchecked                    // probe was aborted before completion (e.g. ctx cancel)
)

// Status records the outcome of probing one dork's direct URL.
type Status struct {
	URL         string
	StatusCode  int    // HTTP status code; 0 if request failed
	Err         string // non-empty on transport/TLS/timeout errors
	Disposition Disposition
}

// Convenience predicates for callers and tests.
func (s Status) Alive() bool     { return s.Disposition == DispositionAlive }
func (s Status) Dead() bool      { return s.Disposition == DispositionDead }
func (s Status) Blocked() bool   { return s.Disposition == DispositionBlocked }
func (s Status) Unchecked() bool { return s.Disposition == DispositionUnchecked }

// Report summarizes a preflight pass.
type Report struct {
	Checked   int // how many dorks were actually probed (direct URLs only)
	Skipped   int // search-engine dorks passed through
	Alive     int
	Dead      int
	Blocked   int      // SSRF guard refused target
	Unchecked int      // probe aborted before completion
	Results   []Status // one entry per checked dork, in input order
}

// Run probes the direct-URL dorks and returns (survivors, report).
// Search-query dorks are always preserved. Direct-URL dorks are dropped when
// HEAD returns 4xx/5xx or transport errors. On context cancellation any
// direct URL that did not complete is preserved in survivors and counted in
// report.Unchecked, so the operator can tell which leads were not validated.
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

	// Custom transport whose DialContext re-validates the resolved IP at
	// connect time. This closes the check-then-use gap where a DNS-
	// rebinding attacker could pass validation and resolve to a private IP
	// at dial time.
	baseDialer := &net.Dialer{Timeout: opts.Timeout}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialValidated(ctx, baseDialer, network, addr, opts.AllowPrivateNetwork)
		},
		TLSHandshakeTimeout:   opts.Timeout,
		ResponseHeaderTimeout: opts.Timeout,
		DisableKeepAlives:     true,
	}
	client := &http.Client{
		Timeout:   opts.Timeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	type job struct {
		idx int
		d   dork.Dork
	}

	// statuses is sized to len(dorks); a nil slot means "non-direct dork,
	// pass-through". A non-nil slot means "direct URL we tried to probe" —
	// the Status itself records whether the probe completed and how.
	var report Report
	statuses := make([]*Status, len(dorks))

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
				st := probe(ctx, client, j.d.Query, opts.UserAgent, opts.AllowPrivateNetwork)
				mu.Lock()
				statuses[j.idx] = &st
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
		st := statuses[i]
		switch {
		case st == nil && isDirectURL(d.Query):
			// Probe never completed (context canceled before scheduling).
			// Preserve the dork as unchecked so the operator can see the lead
			// was not validated, rather than silently dropping it.
			report.Unchecked++
			report.Results = append(report.Results, Status{URL: d.Query, Disposition: DispositionUnchecked})
			survivors = append(survivors, d)
		case st == nil:
			// Non-direct dork: search-engine-wrapped, always passes through.
			survivors = append(survivors, d)
		default:
			report.Results = append(report.Results, *st)
			switch st.Disposition {
			case DispositionBlocked:
				report.Blocked++
			case DispositionDead:
				report.Dead++
			case DispositionUnchecked:
				report.Unchecked++
				survivors = append(survivors, d)
			default: // DispositionAlive
				report.Alive++
				survivors = append(survivors, d)
			}
		}
	}

	return survivors, report
}

func isDirectURL(q string) bool {
	return strings.HasPrefix(q, "http://") || strings.HasPrefix(q, "https://")
}

func probe(ctx context.Context, client *http.Client, target, ua string, allowPrivateNetwork bool) Status {
	st := Status{URL: target}
	classify := func(err error) Disposition {
		// errors.Is unwraps the *url.Error wrapper http.Client puts around
		// dial-time refusals, so the sentinel survives transparently.
		if errors.Is(err, ErrBlockedTarget) {
			return DispositionBlocked
		}
		return DispositionDead
	}

	if err := validateProbeTarget(ctx, target, allowPrivateNetwork); err != nil {
		st.Err = err.Error()
		st.Disposition = classify(err)
		return st
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, target, nil)
	if err != nil {
		st.Err = err.Error()
		st.Disposition = DispositionDead
		return st
	}
	req.Header.Set("User-Agent", ua)
	resp, err := client.Do(req)
	if err != nil {
		st.Err = err.Error()
		st.Disposition = classify(err)
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
		st.Disposition = DispositionDead
	}
	return st
}

func validateProbeTarget(ctx context.Context, target string, allowPrivateNetwork bool) error {
	u, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %q", u.Scheme)
	}
	if allowPrivateNetwork {
		return nil
	}
	host := strings.TrimSuffix(strings.ToLower(u.Hostname()), ".")
	if host == "" {
		return fmt.Errorf("missing URL hostname")
	}
	if host == "localhost" {
		return fmt.Errorf("%w: %s", ErrBlockedTarget, host)
	}
	if ip := net.ParseIP(host); ip != nil {
		if isRestrictedIP(ip) {
			return fmt.Errorf("%w: %s", ErrBlockedTarget, ip.String())
		}
		return nil
	}
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		if isRestrictedIP(addr.IP) {
			return fmt.Errorf("%w: %s resolves to %s", ErrBlockedTarget, host, addr.IP.String())
		}
	}
	return nil
}

// dialValidated wraps the dialer with a second IP check at connect time.
// This closes the DNS-rebinding gap where validateProbeTarget resolves
// safely but the transport's later resolution returns a private IP.
func dialValidated(ctx context.Context, dialer *net.Dialer, network, addr string, allowPrivateNetwork bool) (net.Conn, error) {
	if allowPrivateNetwork {
		return dialer.DialContext(ctx, network, addr)
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(host); ip != nil {
		if isRestrictedIP(ip) {
			return nil, fmt.Errorf("%w at dial: %s", ErrBlockedTarget, ip.String())
		}
		return dialer.DialContext(ctx, network, addr)
	}
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	for _, ip := range addrs {
		if isRestrictedIP(ip.IP) {
			return nil, fmt.Errorf("%w at dial: %s resolves to %s", ErrBlockedTarget, host, ip.IP.String())
		}
	}
	// Dial the first resolved IP to keep the validated address. This trades
	// resilience against multi-homed first-IP failures for TOCTOU-safe SSRF
	// guarantees; the host resolution is intentionally not retried.
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no IPs resolved for %s", host)
	}
	first := addrs[0].IP.String()
	if strings.Contains(first, ":") {
		first = "[" + first + "]"
	}
	return dialer.DialContext(ctx, network, net.JoinHostPort(first, port))
}

func isRestrictedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified()
}
