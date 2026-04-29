package detector

import (
	"testing"
)

func TestParseAstGrepJSON(t *testing.T) {
	input := `[
  {
    "text": "console.log('hello')",
    "range": {
      "start": {"line": 4, "column": 2},
      "end":   {"line": 4, "column": 22}
    },
    "file": "src/index.js",
    "language": "JavaScript"
  },
  {
    "text": "console.log(x)",
    "range": {
      "start": {"line": 11, "column": 0},
      "end":   {"line": 11, "column": 14}
    },
    "file": "src/utils.js",
    "language": "JavaScript"
  }
]`

	matches, err := parseAstGrepJSON([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].File != "src/index.js" || matches[0].Line != 5 {
		t.Errorf("unexpected first match: %+v", matches[0])
	}
	if matches[0].Content != "console.log('hello')" {
		t.Errorf("unexpected content: %q", matches[0].Content)
	}
	if matches[1].File != "src/utils.js" || matches[1].Line != 12 {
		t.Errorf("unexpected second match: %+v", matches[1])
	}
}

func TestAstGrepArgsIncludePatternAndLanguage(t *testing.T) {
	d := &AstGrepDetector{
		Rule:     map[string]any{"pattern": "console.log($$$ARGS)"},
		Language: "JavaScript",
	}

	args, err := d.args("/repo")
	if err != nil {
		t.Fatalf("args failed: %v", err)
	}

	want := []string{"run", "--json=compact", "--pattern", "console.log($$$ARGS)", "--lang", "JavaScript", "/repo"}
	if len(args) != len(want) {
		t.Fatalf("expected %v, got %v", want, args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, args)
		}
	}
}
