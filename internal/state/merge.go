package state

import (
	"math"
	"time"

	"github.com/mesa-studio/codemod-tally/internal/detector"
)

type currentMatch struct {
	fp    Fingerprint
	match detector.Match
	used  bool
}

// Merge reconciles existing state with current scan results.
// An occurrence is marked done when its fingerprint is absent from current matches.
// Already-done occurrences remain done even when absent.
// New fingerprints in current are added as todo.
func Merge(existing *State, current []detector.Match, now time.Time) *State {
	currentMatches := make([]currentMatch, 0, len(current))
	for _, m := range current {
		currentMatches = append(currentMatches, currentMatch{
			fp:    ComputeMatch(m),
			match: m,
		})
	}

	result := &State{
		Recipe:   existing.Recipe,
		RepoRoot: existing.RepoRoot,
		LastScan: now,
	}

	for _, occ := range existing.Occurrences {
		if idx := findExact(currentMatches, occ.Fingerprint); idx >= 0 {
			m := currentMatches[idx].match
			currentMatches[idx].used = true
			occ = updateTodo(occ, currentMatches[idx].fp, m)
		} else if occ.Status == StatusTodo {
			if idx := findNearestSameLine(currentMatches, occ); idx >= 0 {
				m := currentMatches[idx].match
				currentMatches[idx].used = true
				occ = updateTodo(occ, currentMatches[idx].fp, m)
			} else {
				t := now
				occ.Status = StatusDone
				occ.ResolvedAt = &t
			}
		}
		result.Occurrences = append(result.Occurrences, occ)
	}

	for _, current := range currentMatches {
		if !current.used {
			result.Occurrences = append(result.Occurrences, Occurrence{
				Fingerprint: current.fp,
				File:        current.match.File,
				Line:        current.match.Line,
				Content:     current.match.Content,
				Status:      StatusTodo,
				FirstSeen:   now,
			})
		}
	}

	return result
}

func updateTodo(occ Occurrence, fp Fingerprint, m detector.Match) Occurrence {
	occ.Fingerprint = fp
	occ.File = m.File
	occ.Line = m.Line
	occ.Content = m.Content
	occ.Status = StatusTodo
	occ.ResolvedAt = nil
	return occ
}

func findExact(matches []currentMatch, fp Fingerprint) int {
	for i, current := range matches {
		if !current.used && current.fp == fp {
			return i
		}
	}
	return -1
}

func findNearestSameLine(matches []currentMatch, occ Occurrence) int {
	bestIdx := -1
	bestDistance := math.MaxInt

	for i, current := range matches {
		if current.used {
			continue
		}
		if current.match.File != occ.File || current.match.Content != occ.Content {
			continue
		}
		distance := abs(current.match.Line - occ.Line)
		if distance < bestDistance {
			bestDistance = distance
			bestIdx = i
		}
	}

	return bestIdx
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
