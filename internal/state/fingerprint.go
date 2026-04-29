package state

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/mesa-studio/codemod-tally/internal/detector"
)

// Fingerprint is a stable identifier for a match location based on file path and surrounding context.
type Fingerprint string

// Compute returns a 16-char hex fingerprint for the match at file with the given context lines.
// Line number is intentionally excluded so fingerprints survive line shifts.
func Compute(file string, context []string) Fingerprint {
	return hashParts(file, strings.Join(context, "\n"))
}

// ComputeMatch returns a fingerprint for one detector match.
// It includes the matched content and its position within the context window so
// adjacent matches with the same surrounding context remain separate, while
// still avoiding absolute line numbers for ordinary line shifts.
func ComputeMatch(m detector.Match) Fingerprint {
	return hashParts(m.File, strings.Join(m.Context, "\n"), m.Content, strconv.Itoa(m.ContextLine))
}

func hashParts(parts ...string) Fingerprint {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil))[:16])
}
