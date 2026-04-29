package detector_test

import (
	"testing"

	"github.com/mesa-studio/codemod-tally/internal/detector"
)

func TestMatchFields(t *testing.T) {
	m := detector.Match{
		File:    "src/index.js",
		Line:    5,
		Content: "  console.log('hello');",
		Context: []string{"function foo() {", "  console.log('hello');", "}"},
	}
	if m.File != "src/index.js" {
		t.Errorf("expected src/index.js, got %s", m.File)
	}
	if m.Line != 5 {
		t.Errorf("expected line 5, got %d", m.Line)
	}
}
