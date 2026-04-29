---
name: codemod-tally
description: Use when handling large, finite code migrations tracked by Codemod Tally, or when the user asks an agent to continue or complete a Codemod Tally refactoring.
---

# Codemod Tally

Codemod Tally tracks mass refactoring progress with deterministic scans. The agent edits code; Codemod Tally decides what is done.

## Use When

- The task has a finite set of locations that a detector can find.
- The user asks for a large migration, API replacement, import update, or deprecation cleanup.
- A `~/.codemod-tally/state/.../progress.md` file or `codemod-tally prompt <name>` block is present.

Do not use Codemod Tally for vague cleanup, open-ended redesigns, or changes that cannot be verified by a detector.

## Workflow

1. Verify the CLI:
   ```bash
   command -v codemod-tally && codemod-tally doctor
   ```
2. If no recipe exists, create one:
   ```bash
   codemod-tally init <recipe-name> --template ripgrep-text
   ```
   Then edit `detector.yaml` and `recipe.md` before scanning.
3. Run:
   ```bash
   codemod-tally scan <recipe-name>
   codemod-tally prompt <recipe-name>
   ```
4. Work through `Remaining` items in `progress.md`, top to bottom.
5. Edit target repo files only. Never edit `progress.md` or `.scan-cache.json`.
6. If the recipe does not cover a case, write a short note to `journal.md` and skip it.
7. After edits, run `codemod-tally scan <recipe-name>` again. Continue until remaining is 0.

## Rules

- Trust `codemod-tally scan` over the agent's memory.
- Use `codemod-tally status <recipe-name>` for progress without rerunning detectors.
- Stop and ask when the detector output looks wrong, the recipe is ambiguous, or remaining items are not transformable.
