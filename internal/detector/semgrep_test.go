package detector

import (
	"testing"
)

func TestParseSemgrepJSON(t *testing.T) {
	input := `{
  "results": [
    {
      "check_id": "console-log",
      "path": "src/index.js",
      "start": {"line": 5, "col": 3, "offset": 42},
      "end": {"line": 5, "col": 25, "offset": 64},
      "extra": {
        "lines": "  console.log('hello');"
      }
    },
    {
      "check_id": "console-log",
      "path": "src/utils.js",
      "start": {"line": 12, "col": 1, "offset": 100},
      "end": {"line": 12, "col": 20, "offset": 119},
      "extra": {
        "lines": "console.log(x);"
      }
    }
  ],
  "errors": []
}`

	matches, err := parseSemgrepJSON([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].File != "src/index.js" || matches[0].Line != 5 {
		t.Errorf("unexpected first match: %+v", matches[0])
	}
	if matches[0].Content != "  console.log('hello');" {
		t.Errorf("unexpected content: %q", matches[0].Content)
	}
	if matches[1].File != "src/utils.js" || matches[1].Line != 12 {
		t.Errorf("unexpected second match: %+v", matches[1])
	}
}
