# Codemod Tally

[![CI](https://github.com/mesa-studio/codemod-tally/actions/workflows/ci.yml/badge.svg)](https://github.com/mesa-studio/codemod-tally/actions/workflows/ci.yml)

AI-assisted migrations without losing progress.

Project status: early v0.x. Codemod Tally is ready for source-based installation and
day-to-day local use, but it does not publish prebuilt binaries yet.

LLM agents are useful for "replace this pattern across the repo" work, but they often miss locations, lose context, or mark work complete too early. Codemod Tally keeps the search and completion check deterministic: detectors find every location, the agent rewrites code, and `codemod-tally scan` decides what is done.

Codemod Tally is not a codemod engine. Use ripgrep, semgrep, ast-grep, or a shell command as the detector; use your agent for the judgment-heavy rewrite.

## Requirements

- Go 1.26.1 or newer for source installation.
- `git` for repository identity and stable state paths.
- `rg` from ripgrep for the default text-search workflow.
- `semgrep` and `ast-grep` are optional unless a recipe uses those detector types.

## Quick Start

```bash
# Install
go install github.com/mesa-studio/codemod-tally@latest

# Check required tools
codemod-tally doctor

# Create a tracked migration recipe
codemod-tally init console-to-logger --template ripgrep-text

# Edit the generated detector pattern and agent instructions
$EDITOR ~/.codemod-tally/recipes/console-to-logger/detector.yaml
$EDITOR ~/.codemod-tally/recipes/console-to-logger/recipe.md

# Scan the current repository
codemod-tally scan console-to-logger

# Give this block to Codex, Claude Code, Cursor, or another agent
codemod-tally prompt console-to-logger

# After the agent edits files, rescan. Repeat until remaining is 0.
codemod-tally scan console-to-logger
codemod-tally status console-to-logger
```

All Codemod Tally artifacts live in `~/.codemod-tally/`. Nothing is added to the target repository.

To try the full loop without touching a real repository:

```bash
tmp=$(mktemp -d)
repo="$tmp/repo"
recipes="$tmp/recipes"
state="$tmp/state"

mkdir -p "$repo"
git -C "$repo" init -q
printf "console.log('hello')\n" > "$repo/app.js"

codemod-tally --dir "$repo" --recipe-dir "$recipes" --state-dir "$state" \
  init console-to-logger --template ripgrep-text

cat > "$recipes/console-to-logger/detector.yaml" <<'YAML'
type: ripgrep
pattern: 'console\.log\('
flags: []
YAML

codemod-tally --dir "$repo" --recipe-dir "$recipes" --state-dir "$state" \
  scan console-to-logger
codemod-tally --dir "$repo" --recipe-dir "$recipes" --state-dir "$state" \
  status console-to-logger

# Simulate the agent edit.
printf "logger.info('hello')\n" > "$repo/app.js"

codemod-tally --dir "$repo" --recipe-dir "$recipes" --state-dir "$state" \
  scan console-to-logger
codemod-tally --dir "$repo" --recipe-dir "$recipes" --state-dir "$state" \
  status console-to-logger
```

## How it works

1. Write a recipe: detector config, scope, agent instructions, and examples.
2. Run `codemod-tally scan <name>` to create `progress.md`.
3. Run `codemod-tally prompt <name>` and paste the block into an agent.
4. The agent works through `Remaining` items.
5. Run `codemod-tally scan <name>` again. Locations where the detector no longer matches are marked done.
6. Continue until `codemod-tally status <name>` shows 0 remaining.

The detector marks items done, not the agent.

## Install

```bash
go install github.com/mesa-studio/codemod-tally@latest
```

Or build from source:

```bash
git clone https://github.com/mesa-studio/codemod-tally
cd codemod-tally
go build -o codemod-tally .
```

To verify the install:

```bash
command -v codemod-tally
codemod-tally doctor
```

## Commands

| Command | Description |
|---------|-------------|
| `codemod-tally init <name>` | Create a recipe skeleton |
| `codemod-tally init <name> --template <template>` | Create from a built-in template |
| `codemod-tally init --list-templates` | List built-in templates |
| `codemod-tally doctor` | Check local tools and directories |
| `codemod-tally doctor <name>` | Check a recipe and its detector dependency |
| `codemod-tally list` | List available recipes |
| `codemod-tally scan <name>` | Run detector and update `progress.md` |
| `codemod-tally status <name>` | Show progress from cache without rerunning detector |
| `codemod-tally prompt <name>` | Print an agent prompt block |
| `codemod-tally clean <name>` | Delete state for the current repository |

Global flags: `--dir`, `--recipe-dir`, `--state-dir`.

## Templates

```bash
codemod-tally init --list-templates
codemod-tally init api-migration --template ripgrep-text
codemod-tally init js-ast-migration --template semgrep-js
codemod-tally init js-astgrep-migration --template astgrep-js
```

Templates are starting points. Edit the generated `detector.yaml` and `recipe.md` before scanning.

## Recipe structure

```
~/.codemod-tally/recipes/<name>/
  config.yaml      # metadata and scope settings
  detector.yaml    # detector configuration
  recipe.md        # instructions for the agent
  examples/        # reference diffs
```

### config.yaml

```yaml
name: console-to-logger
description: Replace console.log with logger.info
detector: detector.yaml
recipe: recipe.md
examples_dir: examples/
scope:
  include: ["**/*.js", "**/*.ts"]
  exclude: ["**/*.test.js", "**/node_modules/**"]
```

### detector.yaml

Ripgrep:

```yaml
type: ripgrep
pattern: 'console\.log\('
flags: []
```

Shell command:

```yaml
type: shell
command: "rg -n 'console\\.log\\(' --json"
parser: ripgrep
```

Semgrep:

```yaml
type: semgrep
rules:
  - id: console-log
    pattern: console.log(...)
    languages: [javascript, typescript]
```

Ast-grep:

```yaml
type: astgrep
language: JavaScript
rule:
  pattern: console.log($$$ARGS)
```

See `examples/` for complete recipes.

## Agent skill

Codemod Tally includes an optional skill for agents that support installable skills:

```bash
mkdir -p ~/.codex/skills
cp -R skills/codemod-tally ~/.codex/skills/
```

For Claude Code, copy the same `skills/codemod-tally` directory into your Claude skills directory.

The skill is intentionally thin. It teaches the agent when to use Codemod Tally, how to run `doctor`, `scan`, and `prompt`, and that it must never edit `progress.md`.

## Troubleshooting

`codemod-tally doctor` reports missing `rg`: install ripgrep and make sure it is on
`PATH`. The default `ripgrep-text` workflow depends on it.

`codemod-tally doctor` reports missing `semgrep` or `ast-grep`: install the missing tool
only if the recipe uses that detector type. These are optional dependencies.

`codemod-tally doctor` reports missing recipe or state directories: this is normal on a
fresh install. `codemod-tally init <name>` creates recipes, and `codemod-tally scan <name>`
creates state.

`codemod-tally status <name>` reports no state: run `codemod-tally scan <name>` first from the
target repository, or pass the same `--dir` and `--state-dir` values used during
the scan.

## Development

```bash
make check
```

This runs `gofmt` verification, unit tests, `go vet`, and a smoke test of the
scan/prompt/rescan workflow. Release steps are documented in
[`docs/RELEASE.md`](docs/RELEASE.md).

## When Codemod Tally fits

A task is a good fit if it is:

- Findable: remaining locations can be found automatically.
- Verifiable: the detector can confirm when a location is handled.
- Transformable: the rewrite can be described with instructions and examples.

Good fits: API migrations, import changes, language syntax upgrades, library replacements, style unification, deprecations.

Poor fits: vague cleanup, broad redesigns, tasks requiring deep semantic understanding of the whole project, and work with no finite list of locations.

## Design principles

- Search and verification are deterministic.
- Rewriting is the agent's job.
- State lives in files under `~/.codemod-tally/state/`.
- Codemod Tally does not run or orchestrate agents.
- Target repositories stay clean.
