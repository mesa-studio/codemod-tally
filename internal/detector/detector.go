package detector

// Match is a single location where a detector found the pattern.
type Match struct {
	File        string
	Line        int
	Content     string   // the matched line
	Context     []string // populated by scanner: lines [line-N .. line+N]
	ContextLine int      // zero-based index of the matched line in Context
}

// Detector runs a search against a directory and returns all matches.
type Detector interface {
	Run(dir string) ([]Match, error)
}
