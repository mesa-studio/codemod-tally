# Release Process

Codemod Tally v0.x releases are source-only. Do not advertise prebuilt binaries,
package manager taps, or release archives until those artifacts are actually
produced by the project.

## v0.1 Checklist

1. Confirm the repository remote and module path are both
   `github.com/mesa-studio/codemod-tally`.
2. Start from a clean worktree.
3. Run:

   ```bash
   make check
   ```

4. Verify local installation:

   ```bash
   go install .
   command -v codemod-tally
   codemod-tally doctor
   ```

5. Create and push an annotated tag:

   ```bash
   git tag -a v0.1.0 -m "Codemod Tally v0.1.0"
   git push origin main
   git push origin v0.1.0
   ```

6. Create GitHub release notes that describe the source install path:

   ```bash
   go install github.com/mesa-studio/codemod-tally@v0.1.0
   ```

## Future Release Work

Add binary builds only after there is a real release workflow that produces and
checks those artifacts. Until then, keep README and release notes source-first.
