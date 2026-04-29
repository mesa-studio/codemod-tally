package detector

import (
	"testing"
)

func TestShellDetectorLines(t *testing.T) {
	d := &ShellDetector{
		Command: `printf "src/index.js:5:  console.log('hello');\nsrc/utils.js:12:console.log(x);\n"`,
		Parser:  "lines",
	}
	matches, err := d.Run(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d: %+v", len(matches), matches)
	}
	if matches[0].File != "src/index.js" || matches[0].Line != 5 {
		t.Errorf("unexpected first match: %+v", matches[0])
	}
}

func TestParseLines(t *testing.T) {
	input := []byte("src/index.js:5:  console.log('hello');\nsrc/utils.js:12:console.log(x);\n")
	matches, err := parseLines(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].Content != "  console.log('hello');" {
		t.Errorf("unexpected content: %q", matches[0].Content)
	}
}
