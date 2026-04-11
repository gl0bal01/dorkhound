// Package nuclei shells out to the nuclei v2 binary to run OSINT templates
// against a target username and parses the JSONL output back into Dorks.
package nuclei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// DefaultTags is the default --nuclei-tags value.
const DefaultTags = "osint,osint-social"

// Available reports whether the nuclei binary is installed and in PATH.
func Available() bool {
	_, err := exec.LookPath("nuclei")
	return err == nil
}

// Run invokes nuclei once per username with -tags <tags> and parses JSONL
// matches into Dorks (Category "nuclei"). Respects ctx for cancellation.
// Returns dorks, a non-nil error only on hard failures (binary missing,
// context cancelled). Per-username errors are aggregated into the returned
// error but do not prevent partial results.
func Run(ctx context.Context, usernames []string, tags string, timeout time.Duration) ([]dork.Dork, error) {
	if !Available() {
		return nil, fmt.Errorf("nuclei binary not found in PATH; install with: go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest")
	}
	if len(usernames) == 0 {
		return nil, nil
	}
	if tags == "" {
		tags = DefaultTags
	}
	var dorks []dork.Dork
	var errs []string
	for _, user := range usernames {
		user = strings.TrimSpace(user)
		if user == "" {
			continue
		}
		ds, err := runOne(ctx, user, tags, timeout)
		if err != nil {
			errs = append(errs, fmt.Sprintf("username %q: %v", user, err))
			continue
		}
		dorks = append(dorks, ds...)
	}
	if len(errs) > 0 {
		return dorks, fmt.Errorf("nuclei errors: %s", strings.Join(errs, "; "))
	}
	return dorks, nil
}

// nucleiLine is the minimal subset of the nuclei JSONL schema we care about.
type nucleiLine struct {
	TemplateID string `json:"template-id"`
	Info       struct {
		Name     string   `json:"name"`
		Severity string   `json:"severity"`
		Tags     []string `json:"tags"`
	} `json:"info"`
	MatchedAt string `json:"matched-at"`
	Matched   string `json:"matched"` // fallback
	Host      string `json:"host"`
}

func runOne(ctx context.Context, username, tags string, timeout time.Duration) ([]dork.Dork, error) {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// nuclei v2 uses -var key=value to pass template variables; -jsonl for
	// one-line JSON; -silent to suppress progress; -duc to skip update
	// checks; -disable-update-check as alias on older builds.
	args := []string{
		"-tags", tags,
		"-var", "user=" + username,
		"-jsonl",
		"-silent",
		"-duc",
	}
	cmd := exec.CommandContext(cctx, "nuclei", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if cctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout after %s", timeout)
		}
		// non-zero exit with no matches is normal for nuclei; only treat as
		// error if stderr has content AND stdout is empty
		if stdout.Len() == 0 && stderr.Len() > 0 {
			return nil, fmt.Errorf("nuclei: %s", strings.TrimSpace(stderr.String()))
		}
	}
	return parseJSONL(stdout.Bytes(), username), nil
}

func parseJSONL(data []byte, username string) []dork.Dork {
	var dorks []dork.Dork
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var nl nucleiLine
		if err := json.Unmarshal(line, &nl); err != nil {
			continue
		}
		url := nl.MatchedAt
		if url == "" {
			url = nl.Matched
		}
		if url == "" {
			url = nl.Host
		}
		if url == "" {
			continue
		}
		priority := 3
		if nl.Info.Severity == "info" || nl.Info.Severity == "" {
			priority = 2
		}
		label := nl.Info.Name
		if label == "" {
			label = nl.TemplateID
		}
		dorks = append(dorks, dork.Dork{
			Query:    url, // direct URL -> Dork.URL() returns as-is
			Category: "nuclei",
			Region:   "global",
			Priority: priority,
			Label:    fmt.Sprintf("Nuclei[%s] %s → %s", username, label, nl.TemplateID),
		})
	}
	return dorks
}
