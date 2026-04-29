package state

import (
	"testing"
)

func TestComputeFingerprint(t *testing.T) {
	context := []string{
		"function foo() {",
		"  console.log('hello');",
		"}",
	}
	fp1 := Compute("src/index.js", context)
	fp2 := Compute("src/index.js", context)

	if fp1 != fp2 {
		t.Error("same input must produce same fingerprint")
	}
	if len(string(fp1)) != 16 {
		t.Errorf("expected 16 chars, got %d: %s", len(string(fp1)), fp1)
	}
}

func TestComputeFingerprintDifferentFile(t *testing.T) {
	context := []string{"  console.log('hello');"}
	fp1 := Compute("src/a.js", context)
	fp2 := Compute("src/b.js", context)
	if fp1 == fp2 {
		t.Error("different files must produce different fingerprints")
	}
}

func TestComputeFingerprintLineShift(t *testing.T) {
	context := []string{
		"function foo() {",
		"  console.log('hello');",
		"}",
	}
	fp1 := Compute("src/index.js", context)
	fp2 := Compute("src/index.js", context)
	if fp1 != fp2 {
		t.Error("line shift must not change fingerprint")
	}
}
