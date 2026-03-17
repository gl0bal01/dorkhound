package output

import (
	"fmt"
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

func openURL(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", url).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
