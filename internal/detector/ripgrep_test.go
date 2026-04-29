package detector

import (
	"os/exec"
	"testing"
)

func TestParseRipgrepJSON(t *testing.T) {
	input := `{"type":"begin","data":{"path":{"text":"src/index.js"}}}
{"type":"match","data":{"path":{"text":"src/index.js"},"lines":{"text":"  console.log('hello');\n"},"line_number":5,"absolute_offset":42,"submatches":[]}}
{"type":"end","data":{"path":{"text":"src/index.js"},"binary_offset":null,"stats":{}}}
{"type":"match","data":{"path":{"text":"src/utils.js"},"lines":{"text":"console.log(x);\n"},"line_number":12,"absolute_offset":100,"submatches":[]}}`

	matches, err := parseRipgrepJSON([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].File != "src/index.js" {
		t.Errorf("expected src/index.js, got %s", matches[0].File)
	}
	if matches[0].Line != 5 {
		t.Errorf("expected line 5, got %d", matches[0].Line)
	}
	if matches[0].Content != "  console.log('hello');" {
		t.Errorf("unexpected content: %q", matches[0].Content)
	}
	if matches[1].File != "src/utils.js" || matches[1].Line != 12 {
		t.Errorf("unexpected second match: %+v", matches[1])
	}
}

func TestRipgrepDetectorRun(t *testing.T) {
	if _, err := exec.LookPath("rg"); err != nil {
		t.Skip("rg not in PATH")
	}
	d := &RipgrepDetector{Pattern: `console\.log\(`}
	_, err := d.Run(".")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
}
