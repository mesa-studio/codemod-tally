# Contributing

Codemod Tally is a small Go CLI. Keep changes focused, deterministic, and easy to
verify from a clean checkout.

## Local Setup

```bash
go mod download
make check
```

`make check` runs formatting verification, unit tests, `go vet`, and a smoke
test of the scan/prompt/rescan workflow.

## Development Notes

- Keep detector behavior deterministic. Codemod Tally tracks progress by rerunning
  detectors; agent-written status is not trusted.
- Do not hand-edit `progress.md` or `.scan-cache.json` in target repositories.
  They are generated from `codemod-tally scan`.
- Add or update tests when changing detector parsing, state merging, scanner
  behavior, or command output that users rely on.
- New recipe templates should include `config.yaml`, `detector.yaml`, and
  `recipe.md`, plus a smokeable example when possible.

## Pull Request Checklist

Before opening a PR, run:

```bash
make check
```

If a change intentionally alters public CLI behavior, update `README.md` and
`docs/RELEASE.md` where relevant.
