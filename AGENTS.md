# Agent Instructions

Codemod Tally is a Go CLI for tracking finite, detector-backed refactoring work.

- Keep the tool generic. Do not bake in Codex-, Claude-, or project-specific
  assumptions outside optional skill documentation.
- Treat `codemod-tally scan` as the source of truth for progress. Never manually edit
  generated `progress.md` or `.scan-cache.json` files as a way to mark work
  done.
- Keep state under `~/.codemod-tally/` or caller-provided `--state-dir`; do not add
  artifacts to target repositories.
- Before claiming completion, run `make check` or explain exactly which check
  could not be run.
- For public-facing changes, keep README examples aligned with the module path
  `github.com/mesa-studio/codemod-tally`.
