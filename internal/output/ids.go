package output

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// stableID returns a hex-encoded sha256 prefix derived from parts joined by
// NUL. Both the dashboard (DOM IDs) and the CaseBandit exporter (Entity /
// Capture / Case IDs) need stable collision-resistant IDs over similar
// input tuples; this is the one primitive both share.
//
// `n` is the number of leading sha256 bytes to hex-encode. 12 bytes (24 hex
// chars) gives 2^96 collision resistance per category, which is plenty for
// hundreds of dorks per case.
func stableID(n int, parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:n])
}
