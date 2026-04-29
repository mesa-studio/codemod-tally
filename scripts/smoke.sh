#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

REPO="$TMPDIR/repo"
STATE="$TMPDIR/state"

mkdir -p "$REPO"
git -C "$REPO" init -q
printf "console.log('hello')\n" > "$REPO/app.js"

run_codemod_tally() {
	(
		cd "$ROOT"
		go run . --dir "$REPO" --recipe-dir "$ROOT/examples" --state-dir "$STATE" "$@"
	)
}

run_codemod_tally scan ripgrep-text >/dev/null
initial_status=$(run_codemod_tally status ripgrep-text)
printf '%s\n' "$initial_status" | grep -q "1 remaining"

run_codemod_tally prompt ripgrep-text >/dev/null

sed 's/console\.log/logger.info/' "$REPO/app.js" > "$REPO/app.js.new"
mv "$REPO/app.js.new" "$REPO/app.js"

run_codemod_tally scan ripgrep-text >/dev/null
final_status=$(run_codemod_tally status ripgrep-text)
printf '%s\n' "$final_status" | grep -q "0 remaining"

printf 'smoke ok\n'
