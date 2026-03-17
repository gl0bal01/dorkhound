package output

import (
	"bytes"

	"github.com/atotto/clipboard"
	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
)

// Clipboard renders dorks in Discord format and copies the result to the
// system clipboard using github.com/atotto/clipboard.
func Clipboard(c *caseinfo.Case, dorks []dork.Dork, engine string) error {
	var buf bytes.Buffer
	Discord(&buf, c, dorks, engine)
	return clipboard.WriteAll(buf.String())
}
