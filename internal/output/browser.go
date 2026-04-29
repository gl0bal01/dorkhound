package output

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gl0bal01/dorkhound/internal/dork"
)

// OpenInBrowser opens all dork URLs in the default browser, with optional
// batching. batchSize URLs are opened with delay between each; then batchPause
// is slept before the next batch. If batchSize <= 0, no batching is applied.
func OpenInBrowser(dorks []dork.Dork, engine string, delay, batchPause time.Duration, batchSize int) {
	for i, d := range dorks {
		u := d.URL(engine)
		if err := openURL(u); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to open %s: %v\n", d.Label, err)
		}
		if delay > 0 && i < len(dorks)-1 {
			time.Sleep(delay)
		}
		if batchSize > 0 && batchPause > 0 && (i+1)%batchSize == 0 && i < len(dorks)-1 {
			batchNum := (i + 1) / batchSize
			fmt.Fprintf(os.Stderr, "opened batch %d (%d/%d), sleeping %s...\n",
				batchNum, i+1, len(dorks), batchPause)
			time.Sleep(batchPause)
		}
	}
}

func openURL(rawURL string) error {
	// Validate URL to prevent command injection (especially on Windows
	// where cmd /c start interprets shell metacharacters).
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %q", parsed.Scheme)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", rawURL) // #nosec G204 -- rawURL is validated http(s) and exec.Command does not invoke a shell.
	case "darwin":
		cmd = exec.Command("open", rawURL) // #nosec G204 -- rawURL is validated http(s) and exec.Command does not invoke a shell.
	case "windows":
		// Use rundll32 instead of cmd /c start to avoid shell metacharacter injection.
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL) // #nosec G204 -- rawURL is validated http(s) and exec.Command does not invoke a shell.
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	// Reap the child process in the background to avoid zombies.
	go cmd.Wait()
	return nil
}
