package progress

import (
	"strings"
	"testing"
	"time"

	"github.com/mesa-studio/codemod-tally/internal/state"
)

func TestGenerateContainsRemaining(t *testing.T) {
	now := time.Now()
	s := &state.State{
		Recipe:   "console-to-logger",
		LastScan: now,
		Occurrences: []state.Occurrence{
			{Fingerprint: "aaa", File: "src/a.js", Line: 5, Content: "console.log('a');", Status: state.StatusTodo, FirstSeen: now},
			{Fingerprint: "bbb", File: "src/b.js", Line: 12, Content: "console.log('b');", Status: state.StatusDone, FirstSeen: now},
		},
	}

	output := Generate(s)

	if !strings.Contains(output, "## Remaining (1)") {
		t.Errorf("missing Remaining section:\n%s", output)
	}
	if !strings.Contains(output, "## Done (1)") {
		t.Errorf("missing Done section:\n%s", output)
	}
	if !strings.Contains(output, "- [ ] `src/a.js:5`") {
		t.Errorf("missing todo item:\n%s", output)
	}
	if !strings.Contains(output, "- [x] `src/b.js:12`") {
		t.Errorf("missing done item:\n%s", output)
	}
	if !strings.Contains(output, "Do not edit manually") {
		t.Errorf("missing warning:\n%s", output)
	}
}
