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

// OpenInBrowser opens all dork URLs in the default browser.
func OpenInBrowser(dorks []dork.Dork, engine string, delay time.Duration) {
	for _, d := range dorks {
		u := d.URL(engine)
		if err := openURL(u); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to open %s: %v\n", d.Label, err)
		}
		if delay > 0 {
			time.Sleep(delay)
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
		cmd = exec.Command("xdg-open", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	case "windows":
		// Use rundll32 instead of cmd /c start to avoid shell metacharacter injection.
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
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
