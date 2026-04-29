package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCodemodTallySkillFrontmatter(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "skills", "codemod-tally", "SKILL.md"))
	if err != nil {
		t.Fatalf("read skill: %v", err)
	}
	content := string(data)

	if !strings.HasPrefix(content, "---\n") {
		t.Fatal("skill must start with YAML frontmatter")
	}
	if !strings.Contains(content, "\nname: codemod-tally\n") {
		t.Fatal("skill frontmatter must set name: codemod-tally")
	}
	if !strings.Contains(content, "\ndescription: Use when") {
		t.Fatal("skill description must start with 'Use when'")
	}
}
