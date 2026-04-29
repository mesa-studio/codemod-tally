package cmd

import "testing"

func TestScanSummaryLineUsesTrackedWording(t *testing.T) {
	got := scanSummaryLine(3, 1)
	want := "Tracking 3 occurrences (1 new)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
