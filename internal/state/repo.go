package state

import (
	"crypto/sha256"
	"encoding/hex"
	"os/exec"
	"strings"
)

// RepoID returns an 8-char stable identifier for the repository containing dir.
// Uses git remote URL if available, falls back to git root path, then dir path.
func RepoID(dir string) (string, error) {
	key := dir

	if root, err := gitRoot(dir); err == nil {
		key = root
		if remote, err := gitRemote(root); err == nil && remote != "" {
			key = remote
		}
	}

	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])[:8], nil
}

func gitRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitRemote(root string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
